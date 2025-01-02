package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/findyourpaths/phil/glr"
)

func main() {
	glrPkg := "glr."
	if os.Args[1] == "--in-glr-pkg" {
		glrPkg = ""
	}
	grammarPath := os.Args[len(os.Args)-2]
	statesPath := os.Args[len(os.Args)-1]

	pkg, rules, err := readGrammarRules(grammarPath)
	if err != nil {
		panic(fmt.Sprintf("error reading grammar rules: %v", err))
	}

	states, err := readStates(statesPath)
	if err != nil {
		panic(fmt.Sprintf("error reading states: %v", err))
	}

	// Validate that rule references in states are valid
	for stateNum, state := range states.Items {
		for _, actions := range state.Actions {
			for _, action := range actions {
				if action.Type == "reduce" && action.Rule >= len(rules.Items) {
					panic(fmt.Sprintf("state %d references invalid rule number %d", stateNum, action.Rule))
				}
			}
		}
	}

	r := fmt.Sprintf("package %s\n\n", pkg)
	if glrPkg != "" {
		r += "import \"github.com/findyourpaths/phil/glr\"\n\n"
	}

	// Print rules in YACC format.
	r += "/*\nRules\n\n"
	for _, rule := range rules.Items {
		if rule.Nonterminal == "" {
			continue
		}
		rhs := "<empty>"
		if len(rule.RHS) > 0 {
			rhs = strings.Join(rule.RHS, " ")
		}
		r += fmt.Sprintf("%s:\n  %s\n", rule.Nonterminal, rhs)
	}
	r += "*/\n\n"

	r += fmt.Sprintf("var %sRules = &%sRules{Items:[]%sRule{", pkg, glrPkg, glrPkg)
	for i, rule := range rules.Items {
		r += fmt.Sprintf("\n  /* %3d */ %#v,", i, rule)
		if i == 0 {
			r += " // ignored because rule-numbering starts at 1"
		}
	}
	r += "\n}}\n\n"
	if glrPkg == "" {
		r = strings.Replace(r, "glr.Rule", "Rule", -1)
	}

	r += fmt.Sprintf("var %sStates = &%sParseStates{Items:[]%sParseState{", pkg, glrPkg, glrPkg)
	for i, state := range states.Items {
		r += fmt.Sprintf("\n  /* %3d */ %#v,", i, state)
	}
	r += "\n}}\n\n"
	if glrPkg == "" {
		r = strings.Replace(r, "glr.ParseState", "ParseState", -1)
		r = strings.Replace(r, "glr.Action", "Action", -1)
	}

	outPath := strings.TrimSuffix(grammarPath, "_yacc.y") + "_glr.go"
	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(r)
	if err != nil {
		panic(err)
	}
}

func readGrammarRules(grammarFile string) (string, *glr.Rules, error) {
	file, err := os.Open(grammarFile)
	if err != nil {
		return "", nil, fmt.Errorf("failed to open grammar file: %v", err)
	}
	defer file.Close()

	// Rules are numbered starting at 1.
	rs := &glr.Rules{Items: []glr.Rule{{}}}
	scanner := bufio.NewScanner(file)

	inRules := false
	currentRule := &glr.Rule{}
	expectingRHS := false

	nontermRE := regexp.MustCompile(`^(.*):$`)

	pkg := ""
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "package") && pkg == "" {
			pkg = strings.Fields(line)[1]
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
			currentRule = &glr.Rule{
				Nonterminal: strings.TrimSpace(nontermMatch[1]),
			}
			expectingRHS = true
			// slog.Debug("line contains :", "currentRule.NonTerminal", currentRule.Nonterminal)

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
			// slog.Debug("line contains | RHS", "currentRule.NonTerminal", currentRule.Nonterminal)
			// Alternative production for current rule
			if currentRule.Nonterminal == "" {
				return "", nil, fmt.Errorf("alternative production without rule at line %d: %s", lineNum, line)
			}
			parts := strings.SplitN(line, "|", 2)
			rhsPart := strings.TrimSpace(parts[1])
			rhs := parseRHS(rhsPart)
			rs.Items = append(rs.Items, glr.Rule{
				Nonterminal: currentRule.Nonterminal,
				RHS:         rhs,
			})
			expectingRHS = false
		} else {
			// slog.Debug("line contains bare RHS", "currentRule.NonTerminal", currentRule.Nonterminal)
			// Regular production
			if currentRule.Nonterminal == "" {
				return "", nil, fmt.Errorf("production without rule at line %d: %s", lineNum, line)
			}
			rhs := parseRHS(line)
			rs.Items = append(rs.Items, glr.Rule{
				Nonterminal: currentRule.Nonterminal,
				RHS:         rhs,
			})
			expectingRHS = false
		}
	}

	if err := scanner.Err(); err != nil {
		return "", nil, fmt.Errorf("error reading grammar file: %v", err)
	}

	if len(rs.Items) == 0 {
		return "", nil, fmt.Errorf("no valid rules found in grammar file")
	}

	// Print rules.
	for i, rule := range rs.Items {
		if rule.Nonterminal == "" {
			continue
		}
		slog.Debug("i", i, "rule", fmt.Sprintf("%#v", rule))
	}

	return pkg, rs, nil
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

func readStates(statesFile string) (*glr.ParseStates, error) {
	file, err := os.Open(statesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open states file: %v", err)
	}
	defer file.Close()

	rs := &glr.ParseStates{}
	currentState := -1
	scanner := bufio.NewScanner(file)
	lineNum := 0

	srConflictRE := regexp.MustCompile(`^\d+: shift\/reduce conflict \(shift \d+\(\d+\), red'n (\d+)\(\d+\)\) on \S+$`)
	rrConflictRE := regexp.MustCompile(`^\s*\d+: reduce\/reduce conflict  \(red'ns \d+ and (\d+)\) on \S+$`)

	actions := map[string][]glr.Action{}

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
			slog.Debug("conflict reduce rule", rule)
			actions["."] = append(actions["."], glr.Action{
				Type: "reduce",
				Rule: rule,
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
			for len(rs.Items) <= currentState {
				rs.Items = append(rs.Items, glr.ParseState{
					Actions: actions,
					Gotos:   make(map[string]int),
				})
			}
			actions = map[string][]glr.Action{}
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
			rs.Items[currentState].Actions[sym] = append(rs.Items[currentState].Actions[sym], glr.Action{
				Type:  "shift",
				State: state,
			})
		case "reduce":
			rule, err := parseTarget(fields[2])
			if err != nil {
				return nil, fmt.Errorf("invalid reduce target at line %d: %q", lineNum, line)
			}
			rs.Items[currentState].Actions[sym] = append(rs.Items[currentState].Actions[sym], glr.Action{
				Type: "reduce",
				Rule: rule,
			})
		case "goto":
			target, err := parseTarget(fields[2])
			if err != nil {
				return nil, fmt.Errorf("invalid goto target at line %d: %q", lineNum, line)
			}
			rs.Items[currentState].Gotos[sym] = target
		case "accept":
			rs.Items[currentState].Actions[sym] = append(rs.Items[currentState].Actions[sym], glr.Action{
				Type: "accept",
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading states file: %v", err)
	}

	if len(rs.Items) == 0 {
		return nil, fmt.Errorf("no valid states found in states file")
	}

	for i, state := range rs.Items {
		slog.Debug("i", i, "state", fmt.Sprintf("%#v", state))
	}
	return rs, nil
}

func parseTarget(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return num, nil
}
