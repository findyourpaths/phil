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
}

func LoadGrammarRules(grammarFile string) ([]*Rule, error) {
	file, err := os.Open(grammarFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open grammar file: %v", err)
	}
	defer file.Close()

	var rs []*Rule
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
			debugln("line contains :", "currentRule.NonTerminal", currentRule.Nonterminal)
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
			debugln("line contains | RHS", "currentRule.NonTerminal", currentRule.Nonterminal)
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
			debugln("line contains bare RHS", "currentRule.NonTerminal", currentRule.Nonterminal)
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

	for i, rule := range rs {
		debugln("i", i, "rule", fmt.Sprintf("%#v", rule))
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
	debugln("in parseRHS()", "rhs", rhs)
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

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "//") {
			continue
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
					Actions: make(map[string][]StateAction),
					Gotos:   make(map[string]int),
				})
			}
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

		symbol := fields[0]
		actionType := fields[1]

		switch actionType {
		case "shift":
			// Handle lines like "A shift entries, B exceptions"
			if strings.Contains(line, "entries") {
				continue
			}
			state, err := parseNumber(fields[2])
			if err != nil {
				return nil, fmt.Errorf("invalid shift target at line %d: %s", lineNum, line)
			}
			rs[currentState].Actions[symbol] = append(rs[currentState].Actions[symbol], StateAction{
				Action: "shift",
				State:  state,
			})
		case "reduce":
			rule, err := parseNumber(fields[2])
			if err != nil {
				return nil, fmt.Errorf("invalid shift target at line %d: %s", lineNum, line)
			}
			rs[currentState].Actions[symbol] = append(rs[currentState].Actions[symbol], StateAction{
				Action: "reduce",
				Rule:   rule - 1,
			})
		case "goto":
			target, err := parseNumber(fields[2])
			if err != nil {
				return nil, fmt.Errorf("invalid goto target at line %d: %s", lineNum, line)
			}
			rs[currentState].Gotos[symbol] = target
		case "accept":
			rs[currentState].Actions[symbol] = append(rs[currentState].Actions[symbol], StateAction{
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

	for _, rule := range rules {
		debugf("%s:\n", rule.Nonterminal)
		prefix := ""
		if rule.RHS != nil {
			debugf("  %s%s\n", prefix, strings.Join(rule.RHS, " "))
		}
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

func parseNumber(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return num, nil
}
