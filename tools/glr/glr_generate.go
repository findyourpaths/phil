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
	logLevel := slog.LevelDebug
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	glrPkg := "glr."
	if os.Args[1] == "--in-glr-pkg" {
		glrPkg = ""
	}
	grammarPath := os.Args[len(os.Args)-2]
	statesPath := os.Args[len(os.Args)-1]

	pkg, rls, sas, err := readGrammarRules(grammarPath)
	if err != nil {
		panic(fmt.Sprintf("error reading grammar rules: %v", err))
	}

	states, err := readStates(statesPath)
	if err != nil {
		panic(fmt.Sprintf("error reading states: %v", err))
	}

	r := fmt.Sprintf("package %s\n\n", pkg)
	if glrPkg != "" {
		r += "import \"github.com/findyourpaths/phil/glr\"\n\n"
	}

	// Print rules in YACC format.
	r += "/*\nRules\n\n"
	for _, rule := range rls.Items {
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
	for i, rule := range rls.Items {
		r += fmt.Sprintf("\n  /* %3d */ %#v,", i, rule)
		if i == 0 {
			r += " // ignored because rule-numbering starts at 1"
		}
	}
	r += "\n}}\n\n"
	if glrPkg == "" {
		r = strings.Replace(r, "glr.Rule", "Rule", -1)
	}

	// Generate semantic action functions
	r += "// Semantic action functions\n\n"
	r += fmt.Sprintf("var %sActions = &%sSemanticActions{Items:[]any{", pkg, glrPkg)
	for i, action := range sas.Items {
		slog.Debug("", "i", i, "action", action)
		rule := rls.Items[i]
		if rule.Type == "" {
			r += fmt.Sprintf("\n  /* %3d */ nil, // empty type", i)
			continue
		}
		r += fmt.Sprintf("\n  /* %3d */ %s,", i, action.(string))
		// r += fmt.Sprintf("func %sSemanticAction%d(node *%sParseNode) %s {\n", pkg, i, glrPkg, rule.Type)
		// r += fmt.Sprintf("  return %s\n", action.Action)
		// r += "}\n\n"
	}
	r += "\n}}\n\n"

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
	_, err = f.WriteString(r)
	if err != nil {
		panic(err)
	}
	f.Close()

	// Validate that rule references in states are valid
	for stateNum, state := range states.Items {
		for _, actions := range state.Actions {
			for _, action := range actions {
				if action.Type == "reduce" && action.Rule >= len(rls.Items) {
					panic(fmt.Sprintf("state %d references invalid rule number %d", stateNum, action.Rule))
				}
			}
		}
	}
}

func readGrammarRules(p string) (string, *glr.Rules, *glr.SemanticActions, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to open grammar file: %v", err)
	}
	defer f.Close()

	// Rules are numbered starting at 1.
	rls := &glr.Rules{Items: []glr.Rule{{}}}
	sas := &glr.SemanticActions{Items: []any{nil}}
	scanner := bufio.NewScanner(f)

	inRules := false
	inUnion := false
	currentNonterm := ""
	expectingRHS := false

	nontermRE := regexp.MustCompile(`^(.*):$`)

	// Map to store type declarations
	typeFieldsByNonterm := map[string]string{}
	typesByField := map[string]string{}
	typesByNonterm := map[string]string{}

	pkg := ""
	i := 0
	for scanner.Scan() {
		i++
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		slog.Debug("in readGrammarRules()", "i", i, "inUnion", inUnion, "inRules", inRules)
		slog.Debug("in readGrammarRules()", "i", i, "line", line)
		slog.Debug("in readGrammarRules()", "i", i, "fields", fields)

		// Parse %type declarations, which specifies the type field for each nonterminals.
		if len(fields) == 3 && fields[0] == "%type" {
			tField := strings.Trim(fields[1], "<>")
			ntSym := fields[2]
			typeFieldsByNonterm[ntSym] = tField
			slog.Debug("in readGrammarRules()", "ntSym", ntSym, "tField", tField)
			continue
		}

		// Parse %union declaration, which specifies the type of each field.
		if len(fields) > 0 && fields[0] == "%union" {
			inUnion = true
			continue
		}
		if len(fields) == 1 && fields[0] == "}" && inUnion {
			inUnion = false
			for ntSym, tField := range typeFieldsByNonterm {
				typesByNonterm[ntSym] = typesByField[tField]
			}
			slog.Debug("in readGrammarRules()", "typesByNonterm", typesByNonterm)
			continue
		}
		if inUnion {
			slog.Debug("in readGrammarRules()", "fields", fields)
			if len(fields) == 2 {
				tField := fields[0]
				tType := fields[1]
				typesByField[tField] = tType
				slog.Debug("in readGrammarRules()", "tField", tField, "tType", tType)
			}
			continue
		}

		if strings.HasPrefix(line, "package") && pkg == "" {
			pkg = strings.Fields(line)[1]
		}

		// Start of rules section
		if line == "%%" {
			if !inRules {
				inRules = true
				slog.Debug("in readGrammarRules()", "typesByNonterm", typeFieldsByNonterm)
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
		} else if strings.HasPrefix(line, "%") {
			// Skip code.
			continue
		} else if strings.HasPrefix(line, ";") {
			currentNonterm = ""
		} else if (line == "" || strings.HasPrefix(strings.TrimSpace(line), "{")) && expectingRHS == false {
			// Skip empty lines only if there's no current rule.
			continue
		} else if nontermMatch := nontermRE.FindStringSubmatch(line); len(nontermMatch) > 1 {
			ntSym := strings.TrimSpace(nontermMatch[1])
			currentNonterm = ntSym
			expectingRHS = true
			slog.Debug("line contains :", "currentNonterm", currentNonterm)
		} else {
			slog.Debug("line contains RHS", "currentNonterm", currentNonterm)
			if currentNonterm == "" {
				return "", nil, nil, fmt.Errorf("production without rule at line %d: %s", i, line)
			}
			rhsLine := line
			if fields[0] == "|" {
				rhsLine = strings.Join(fields[1:], " ")
			}
			rhs, action := parseRHS(rhsLine)
			rule := glr.Rule{
				Nonterminal: currentNonterm,
				RHS:         rhs,
				Type:        typesByNonterm[currentNonterm],
			}
			rls.Items = append(rls.Items, rule)
			sas.Items = append(sas.Items, newSemanticAction(typesByNonterm, rule, action))
			slog.Debug("bare RHS", "rule", rls.Items[len(rls.Items)-1], "semantic action", sas.Items[len(sas.Items)-1])
			expectingRHS = false
		}
	}

	if err := scanner.Err(); err != nil {
		return "", nil, nil, fmt.Errorf("error reading grammar file: %v", err)
	}

	if len(rls.Items) == 0 {
		return "", nil, nil, fmt.Errorf("no valid rules found in grammar file")
	}

	return pkg, rls, sas, nil
}

var symbolRE = regexp.MustCompile(`\$(\d+)`)

func newSemanticAction(typesByNonterm map[string]string, rule glr.Rule, action string) string {
	symCounts := map[string]int{}
	symByPos := map[int]string{}
	for i, sym := range rule.RHS {
		symCounts[sym]++
		symByPos[i+1] = sym + strconv.Itoa(symCounts[sym])
	}

	// Replace $n variables with symbolN names
	if action == "" {
		action = "$1"
	}
	act := symbolRE.ReplaceAllStringFunc(action, func(match string) string {
		pos, _ := strconv.Atoi(match[1:])
		if sym, ok := symByPos[pos]; ok {
			return sym
		}
		return match
	})

	// Build arguments string
	var args []string
	for i := 1; i <= len(rule.RHS); i++ {
		sym := symByPos[i]
		t, found := typesByNonterm[rule.RHS[i-1]]
		if !found {
			t = "string"
		}
		args = append(args, fmt.Sprintf("%s %s", sym, t))
	}
	return fmt.Sprintf("func (%s) %s {return %s}", strings.Join(args, ", "), rule.Type, act)
}

func parseRHS(line string) ([]string, string) {
	line = strings.TrimSpace(line)
	if strings.HasSuffix(line, ";") {
		line = line[:len(line)-1]
	}

	// Find semantic action between braces
	var action string
	braceStart := strings.Index(line, "{")
	if braceStart != -1 {
		braceEnd := strings.LastIndex(line, "}")
		if braceEnd != -1 {
			action = strings.TrimPrefix(strings.TrimSpace(line[braceStart+1:braceEnd]), "$$ = ")
			line = strings.TrimSpace(line[:braceStart])
		}
	}

	slog.Debug("in parseRHS", "line", line, "strings.Fields(line)", strings.Fields(line))
	var rhs []string
	for _, token := range strings.Fields(line) {
		if token != "|" && token != ";" && token != "" {
			rhs = append(rhs, strings.TrimSpace(token))
		}
	}
	slog.Debug("in parseRHS", "rhs", rhs, "action", action)
	return rhs, action
}

func readStates(p string) (*glr.ParseStates, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("failed to open states file: %v", err)
	}
	defer f.Close()

	rs := &glr.ParseStates{}
	currentState := -1
	scanner := bufio.NewScanner(f)
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
			slog.Debug("conflict reduce rule", "rule", rule)
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
	return rs, nil
}

func parseTarget(s string) (int, error) {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return num, nil
}
