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

func loadGrammarRules(grammarFile string) ([]Rule, error) {
	file, err := os.Open(grammarFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open grammar file: %v", err)
	}
	defer file.Close()

	var rules []Rule
	scanner := bufio.NewScanner(file)

	inRules := false
	currentRule := Rule{}

	nontermRE := regexp.MustCompile(`^(.*):$`)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

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

		nontermMatch := nontermRE.FindStringSubmatch(line)
		if len(nontermMatch) > 1 {
			currentRule = Rule{
				Nonterminal: strings.TrimSpace(nontermMatch[1]),
			}
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
			debugln("line contains |", "currentRule.NonTerminal", currentRule.Nonterminal)
			// Alternative production for current rule
			if currentRule.Nonterminal == "" {
				return nil, fmt.Errorf("alternative production without rule at line %d: %s", lineNum, line)
			}
			parts := strings.SplitN(line, "|", 2)
			rhsPart := strings.TrimSpace(parts[1])
			rhs := parseRHS(rhsPart)
			if len(rhs) > 0 {
				rule := Rule{
					Nonterminal: currentRule.Nonterminal,
					RHS:         rhs,
				}
				rules = append(rules, rule)
			}
		} else if !strings.HasPrefix(line, "%") && !strings.HasPrefix(line, "{") {
			// Regular production
			if currentRule.Nonterminal == "" {
				return nil, fmt.Errorf("production without rule at line %d: %s", lineNum, line)
			}
			rhs := parseRHS(line)
			if len(rhs) > 0 {
				rule := Rule{
					Nonterminal: currentRule.Nonterminal,
					RHS:         rhs,
				}
				rules = append(rules, rule)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading grammar file: %v", err)
	}

	if len(rules) == 0 {
		return nil, fmt.Errorf("no valid rules found in grammar file")
	}

	for i, rule := range rules {
		debugln("i", i, "rule", fmt.Sprintf("%#v", rule))
	}
	return rules, nil
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

func loadStates(statesFile string) ([]ParseState, error) {
	file, err := os.Open(statesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open states file: %v", err)
	}
	defer file.Close()

	var states []ParseState
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
			for len(states) <= currentState {
				states = append(states, ParseState{
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
			states[currentState].Actions[symbol] = append(states[currentState].Actions[symbol], StateAction{
				Action: "shift",
				State:  state,
			})
		case "reduce":
			rule, err := parseNumber(fields[2])
			if err != nil {
				return nil, fmt.Errorf("invalid shift target at line %d: %s", lineNum, line)
			}
			states[currentState].Actions[symbol] = append(states[currentState].Actions[symbol], StateAction{
				Action: "reduce",
				Rule:   rule - 1,
			})
		case "goto":
			target, err := parseNumber(fields[2])
			if err != nil {
				return nil, fmt.Errorf("invalid goto target at line %d: %s", lineNum, line)
			}
			states[currentState].Gotos[symbol] = target
		case "accept":
			states[currentState].Actions[symbol] = append(states[currentState].Actions[symbol], StateAction{
				Action: "accept",
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading states file: %v", err)
	}

	if len(states) == 0 {
		return nil, fmt.Errorf("no valid states found in states file")
	}

	for i, state := range states {
		debugln("i", i, "state", fmt.Sprintf("%#v", state))
	}
	return states, nil
}

func loadGrammarRulesAndStates(grammarFile string, statesFile string) ([]Rule, []ParseState, error) {
	rules, err := loadGrammarRules(grammarFile)
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

	states, err := loadStates(statesFile)
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
