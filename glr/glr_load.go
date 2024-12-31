package glr

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Rule struct {
	Nonterminal string
	RHS         []string
}

type ParseState struct {
	Actions map[string][]StateAction
	Gotos   map[string]int
}

type StateAction struct {
	Action string
	State  int
	Rule   int
}

type ParseNode struct {
	symbol   string
	value    interface{}
	children []*ParseNode
	startPos int
	endPos   int
	numTerms int
	isAlt    bool
}

func LoadGrammarRules(grammarFile string) ([]*Rule, error) {
	file, err := os.Open(grammarFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open grammar file: %v", err)
	}
	defer file.Close()

	// Rules are numbered starting at 1.
	rs := []*Rule{nil}
	scanner := bufio.NewScanner(file)

	inRules := false
	currentRule := &Rule{}
	expectingRHS := false

	nontermRE := regexp.MustCompile(`^(.*):$`)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Start of rules section
		if line == "%%" {
			if !inRules {
				inRules = true
			} else {
				break
			}
			continue
		}
		if !inRules {
			continue
		}

		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/* ") {
			// Skip comments.
			continue
		} else if strings.HasPrefix(line, "%") || strings.HasPrefix(line, "{") {
			// Skip code.
			continue
		} else if strings.HasPrefix(line, ";") {
			currentRule = nil
		} else if line == "" && expectingRHS == false {
			// Skip empty lines only if there's no current rule.
			continue
		} else if nontermMatch := nontermRE.FindStringSubmatch(line); len(nontermMatch) > 1 {
			currentRule = &Rule{
				Nonterminal: strings.TrimSpace(nontermMatch[1]),
			}
			expectingRHS = true
			// debugln("line contains :", "currentRule.NonTerminal", currentRule.Nonterminal)

			// Handle case where RHS is on same line as colon
			// rhsPart := strings.TrimSpace(parts[1])
			// if rhsPart != "" && !strings.HasPrefix(rhsPart, "{") {
			// 	rhs := parseRHS(rhsPart)
			// 	if len(rhs) > 0 {
			// 		rule := Rule{
			// 			NonTerminal: currentRule.NonTerminal,
			// 			RHS:         rhs,
			// 			Action:      createRuleAction(currentRule.NonTerminal, rhs),
			// 		}
			// 		rules = append(rules, rule)
			// 	}
			// }
		} else if strings.Contains(line, "|") {
			// debugln("line contains | RHS", "currentRule.NonTerminal", currentRule.Nonterminal)
			// Alternative production for current rule
			if currentRule.Nonterminal == "" {
				return nil, fmt.Errorf("alternative production without rule at line %d: %s", lineNum, line)
			}
			parts := strings.SplitN(line, "|", 2)
			rhsPart := strings.TrimSpace(parts[1])
			rhs := parseRHS(rhsPart)
			rs = append(rs, &Rule{
				Nonterminal: currentRule.Nonterminal,
				RHS:         rhs,
			})
			expectingRHS = false
		} else {
			// debugln("line contains bare RHS", "currentRule.NonTerminal", currentRule.Nonterminal)
			// Regular production
			if currentRule.Nonterminal == "" {
				return nil, fmt.Errorf("production without rule at line %d: %s", lineNum, line)
			}
			rhs := parseRHS(line)
			rs = append(rs, &Rule{
				Nonterminal: currentRule.Nonterminal,
				RHS:         rhs,
			})
			expectingRHS = false
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading grammar file: %v", err)
	}

	if len(rs) == 0 {
		return nil, fmt.Errorf("no valid rules found in grammar file")
	}

	// Print rules.
	for i, rule := range rs {
		if rule == nil {
			continue
		}
		debugln("i", i, "rule", fmt.Sprintf("%#v", rule))
	}

	// Print rules in YACC format.
	for _, rule := range rs {
		if rule == nil {
			continue
		}
		debugf("%s:\n", rule.Nonterminal)
		prefix := ""
		if rule.RHS != nil {
			debugf("  %s%s\n", prefix, strings.Join(rule.RHS, " "))
		}
	}

	return rs, nil
}

func parseRHS(line string) []string {
	// Find the position of the first opening brace
	braceIndex := strings.Index(line, "{")
	if braceIndex != -1 {
		// Only take the part before the brace
		line = line[:braceIndex]
	}

	line = strings.TrimSpace(line)
	if strings.HasSuffix(line, ";") {
		line = line[:len(line)-1]
	}

	var rhs []string
	for _, token := range strings.Fields(line) {
		if token != "|" && token != ";" && token != "" {
			rhs = append(rhs, strings.TrimSpace(token))
		}
	}
	return rhs
}

func LoadStates(statesFile string) ([]*ParseState, error) {
	file, err := os.Open(statesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open states file: %v", err)
	}
	defer file.Close()

	var rs []*ParseState
	currentState := -1
	scanner := bufio.NewScanner(file)
	lineNum := 0

	srConflictRE := regexp.MustCompile(`^\d+: shift\/reduce conflict \(shift \d+\(\d+\), red'n (\d+)\(\d+\)\) on \S+$`)
	rrConflictRE := regexp.MustCompile(`^\s*\d+: reduce\/reduce conflict  \(red'ns \d+ and (\d+)\) on \S+$`)

	actions := map[string][]StateAction{}

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		match := srConflictRE.FindStringSubmatch(line)
		if len(match) == 0 {
			match = rrConflictRE.FindStringSubmatch(line)
		}

		if len(match) > 1 {
			rule, err := parseTarget(match[1])
			if err != nil {
				return nil, fmt.Errorf("invalid conflict target at line %d: %q", lineNum, line)
			}
			debugln("conflict reduce rule", rule)
			actions["."] = append(actions["."], StateAction{
				Action: "reduce",
				Rule:   rule,
			})
		}

		if strings.HasPrefix(line, "state ") {
			stateStr := strings.TrimPrefix(line, "state ")
			newState, err := strconv.Atoi(stateStr)
			if err != nil {
				return nil, fmt.Errorf("invalid state number at line %d: %s", lineNum, line)
			}
			currentState = newState
			// Ensure states slice has enough capacity
			for len(rs) <= currentState {
				rs = append(rs, &ParseState{
					Actions: actions,
					Gotos:   make(map[string]int),
				})
			}
			actions = map[string][]StateAction{}
			continue
		}

		if currentState < 0 {
			return nil, fmt.Errorf("action or goto before state declaration at line %d: %s", lineNum, line)
		}

		// Parse actions and gotos
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue // Skip malformed lines
		}

		// Special handling for summary lines
		if strings.Contains(line, "entries") || strings.Contains(line, "reductions") {
			continue
		}

		sym := fields[0]
		action := fields[1]

		switch action {
		case "shift":
			// Handle lines like "A shift entries, B exceptions"
			if strings.Contains(line, "entries") {
				continue
			}
			state, err := parseTarget(fields[2])
			if err != nil {
				return nil, fmt.Errorf("invalid shift target at line %d: %q", lineNum, line)
			}
			rs[currentState].Actions[sym] = append(rs[currentState].Actions[sym], StateAction{
				Action: "shift",
				State:  state,
			})
		case "reduce":
			rule, err := parseTarget(fields[2])
			if err != nil {
				return nil, fmt.Errorf("invalid reduce target at line %d: %q", lineNum, line)
			}
			rs[currentState].Actions[sym] = append(rs[currentState].Actions[sym], StateAction{
				Action: "reduce",
				Rule:   rule,
			})
		case "goto":
			target, err := parseTarget(fields[2])
			if err != nil {
				return nil, fmt.Errorf("invalid goto target at line %d: %q", lineNum, line)
			}
			rs[currentState].Gotos[sym] = target
		case "accept":
			rs[currentState].Actions[sym] = append(rs[currentState].Actions[sym], StateAction{
				Action: "accept",
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading states file: %v", err)
	}

	if len(rs) == 0 {
		return nil, fmt.Errorf("no valid states found in states file")
	}

	for i, state := range rs {
		debugln("i", i, "state", fmt.Sprintf("%#v", state))
	}
	return rs, nil
}

func LoadGrammarRulesAndStates(grammarFile string, statesFile string) ([]*Rule, []*ParseState, error) {
	rules, err := LoadGrammarRules(grammarFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading grammar rules: %v", err)
	}

	states, err := LoadStates(statesFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading states: %v", err)
	}

	// Validate that rule references in states are valid
	for stateNum, state := range states {
		for _, actions := range state.Actions {
			for _, action := range actions {
				if action.Action == "reduce" && action.Rule >= len(rules) {
					return nil, nil, fmt.Errorf("state %d references invalid rule number %d", stateNum, action.Rule)
				}
			}
		}
	}

	return rules, states, nil
}

func parseTarget(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return num, nil
}
