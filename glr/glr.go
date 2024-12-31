package glr

import (
	"fmt"
	"sort"
)

// Debug flags
// var glrDebug = true

var glrDebug = false

// SetDebug toggles debug logging
func SetDebug(enabled bool) {
	glrDebug = enabled
}

// debugf prints debug messages if debug is enabled
func debugf(format string, args ...interface{}) {
	if glrDebug {
		fmt.Printf(format, args...)
	}
}

// debugf prints debug messages if debug is enabled
func debugln(args ...interface{}) {
	if glrDebug {
		fmt.Println(args...)
	}
}

func printNodeTree(n *ParseNode, spaces string) {
	debugf("%s%d: [%d, %d]: symbol: %q, value: %#v\n", spaces, n.numTerms, n.startPos, n.endPos, n.symbol, n.value)
	for _, child := range n.children {
		printNodeTree(child, spaces+"  ")
	}
}

// printActiveParser prints all active parsers
func printAllParsers(s *GLRState, p *StackNode) {
	if p != nil {
		debugf("  current parser with state: %v\n", p.state)
		printParser(p, "  ", nil)
	} else {
		debugf("  no current parser\n")
	}
	printParsers("active parsers", s.activeParsers, nil)
	printParsers("parsers to act", s.parsersToAct, nil)
	printParsers("parsers to shift", s.parsersToShift, s.statesToShiftTo)
	printParsers("accepting parser", s.acceptingParsers, nil)
}

func printParsers(label string, ps []*StackNode, states []int) {
	debugf("  %d %s with states: %v\n", len(ps), label, mapStates(ps))
	for i, p := range ps {
		var state *int
		if states != nil {
			state = &(states[i])
		}
		printParser(p, "  ", state)
	}
}

// printParser prints the parser state and backlinks recursively
func printParser(p *StackNode, parsersAfter string, state *int) {
	shift := ""
	if state != nil {
		shift = fmt.Sprintf("-> %d", *state)
	}

	if len(p.backlinks) == 0 {
		debugf("    %d %s%s\n", p.state, parsersAfter, shift)
		return
	}
	for _, backlink := range p.backlinks {
		backStackNode := backlink.stackNode
		if p.state == backStackNode.state {
			debugf("* %d %s%s\n", p.state, parsersAfter, shift)
		} else {
			printParser(backStackNode, fmt.Sprintf("- %s - %d %s%s", backlink.node.symbol, p.state, parsersAfter, shift), nil)
		}
	}
}

// mapStates returns a slice of parser states
func mapStates(ps []*StackNode) []int {
	states := make([]int, len(ps))
	for i, p := range ps {
		states[i] = p.state
	}
	return states
}

// StackNode represents a node in the GLR parsing stack
type StackNode struct {
	state     int
	backlinks []*StackLink
}

// StackLink represents a link between stack nodes
type StackLink struct {
	stackNode *StackNode
	node      *ParseNode
}

// GLRState maintains the state of GLR parsing
type GLRState struct {
	activeParser     *StackNode
	activeParsers    []*StackNode
	initialParsers   []*StackNode
	parsersToAct     []*StackNode
	parsersToShift   []*StackNode
	statesToShiftTo  []int
	acceptingParsers []*StackNode
	ruleNodes        map[string]*ParseNode
	symbolNodes      []*ParseNode
	lookahead        *ParseNode
	debug            bool
}

// Parse implements GLR parsing algorithm
func Parse(rls []Rule, sts []ParseState, input string) ([]*ParseNode, error) {
	debugf("\nstarting GLR parse\n")

	lexer := newLexer(input)
	s := &GLRState{
		ruleNodes: make(map[string]*ParseNode),
		debug:     false,
	}

	// Initialize with start state
	firstParser := &StackNode{state: 0, backlinks: nil}
	s.activeParsers = []*StackNode{firstParser}
	debugf("initialized with start state 0\n")

	lval := &yySymType{}
	var token *ParseNode
	pos := 0

	for {
		// Get next token
		tokenType := lexer.Lex(lval)

		// Create parse node for token
		token = &ParseNode{
			symbol:   lval.token,
			value:    tokenType,
			startPos: pos,
			endPos:   pos + 1,
			numTerms: 1,
		}
		pos++

		if !parseSymbol(rls, sts, s, token) {
			continue // Skip invalid tokens
		}

		if tokenType < 0 {
			break
		}
	}

	// Handle end of input
	if len(s.acceptingParsers) == 0 {
		return nil, fmt.Errorf("parsing failed")
	}

	var results []*ParseNode
	for _, parser := range s.acceptingParsers {
		if len(parser.backlinks) > 0 {
			results = append(results, parser.backlinks[0].node)
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].numTerms > results[j].numTerms
	})

	for i, result := range results {
		debugf("result[%d]\n", i)
		printNodeTree(result, "")
	}
	return results, nil
}

func parseSymbol(rls []Rule, sts []ParseState, s *GLRState, tok *ParseNode) bool {
	debugf(fmt.Sprintf("\nparsing lookahead symbol: %q with value: %#v\n", tok.symbol, tok.value))
	printAllParsers(s, nil)

	s.lookahead = tok
	s.initialParsers = s.activeParsers
	s.parsersToAct = s.activeParsers
	s.parsersToShift = nil
	s.statesToShiftTo = nil
	s.ruleNodes = make(map[string]*ParseNode)
	s.symbolNodes = nil

	// Process all parsers
	for len(s.parsersToAct) > 0 {
		debugf("processing parsers\n")
		p := s.parsersToAct[0]
		s.parsersToAct = s.parsersToAct[1:]
		s.activeParser = p
		printAllParsers(s, p)
		actor(rls, sts, s, p)
	}

	// Perform shifts if any available
	if len(s.parsersToShift) == 0 && (len(s.acceptingParsers) == 0 || tok.symbol != "$end") {
		// No valid actions found
		debugf("  accepting with no valid actions found\n")
		s.acceptingParsers = append(s.acceptingParsers, s.initialParsers...)
		return false
	}
	s.activeParsers = shifter(s)
	return true
}

func actor(rls []Rule, sts []ParseState, s *GLRState, p *StackNode) {
	as := append(sts[p.state].Actions[s.lookahead.symbol], sts[p.state].Actions["."]...)
	debugf("found %d actions for p.state: %d, s.lookahead.symbol: %q\n", len(as), p.state, s.lookahead.symbol)
	for i, a := range as {
		debugf("  looking at action for p.state: %d, s.lookahead.symbol: %q, actions[%d]: %#v\n", p.state, s.lookahead.symbol, i, a)

		switch a.Action {
		case "shift":
			debugf("  shifting to state: %d\n", a.State)
			s.parsersToShift = append(s.parsersToShift, p)
			s.statesToShiftTo = append(s.statesToShiftTo, a.State)

		case "reduce":
			r := rls[a.Rule]
			debugf("  reducing with rule: %#v\n", r)
			doReductions(rls, sts, s, p, r, len(r.RHS), nil, nil, true)

		case "accept":
			debugf("  accepting\n")
			s.acceptingParsers = append(s.acceptingParsers, p)
		}
	}
}

func doReductions(rls []Rule, sts []ParseState, s *GLRState, p *StackNode, r Rule, length int, kids []*ParseNode, linkToSee *StackLink, linkSeen bool) {
	if length == 0 {
		if linkSeen {
			reducer(rls, sts, s, p, r, kids)
		}
		return
	}

	for _, link := range p.backlinks {
		newLinkSeen := linkSeen || link == linkToSee
		doReductions(rls, sts, s, link.stackNode, r, length-1, append([]*ParseNode{link.node}, kids...), linkToSee, newLinkSeen)
	}
}

func reducer(rls []Rule, sts []ParseState, s *GLRState, p *StackNode, r Rule, kids []*ParseNode) {
	printAllParsers(s, p)

	gotoState, ok := sts[p.state].Gotos[r.Nonterminal]
	if !ok {
		return
	}

	ruleNode := getRuleNode(s, r, kids)
	stackNode := getStackNode(s.activeParsers, gotoState)
	fmt.Println("  gotoState", gotoState, "stackNode", stackNode)
	printNodeTree(ruleNode, "  ")

	if stackNode == nil {
		// Create new parser state
		nonterminal := getSymbolNode(s, ruleNode)
		stackNode = &StackNode{
			state:     gotoState,
			backlinks: []*StackLink{{stackNode: p, node: nonterminal}}}
		s.activeParsers = append(s.activeParsers, stackNode)
		s.parsersToAct = append(s.parsersToAct, stackNode)
		return
	}

	// Check for existing path
	for _, link := range stackNode.backlinks {
		if link.stackNode == p {
			link.node = addAlternative(s, link.node, ruleNode)
			return
		}
	}

	// Add new path
	nonterminal := getSymbolNode(s, ruleNode)
	newLink := &StackLink{stackNode: p, node: nonterminal}
	stackNode.backlinks = append(stackNode.backlinks, newLink)

	// Process additional reductions
	for _, otherParser := range s.activeParsers {
		if !contains(s.parsersToAct, otherParser) {
			actions := sts[otherParser.state].Actions[s.lookahead.symbol]
			for _, action := range actions {
				if action.Action == "reduce" {
					otherRule := rls[action.Rule]
					doReductions(rls, sts, s, otherParser, otherRule, len(otherRule.RHS), nil, newLink, false)
				}
			}
		}
	}
}

func shifter(s *GLRState) []*StackNode {
	var newParsers []*StackNode

	for i, p := range s.parsersToShift {
		newState := s.statesToShiftTo[i]
		newLink := &StackLink{stackNode: p, node: s.lookahead}
		stackNode := getStackNode(newParsers, newState)

		if stackNode != nil {
			stackNode.backlinks = append(stackNode.backlinks, newLink)
			continue
		}

		stackNode = &StackNode{state: newState, backlinks: []*StackLink{newLink}}
		newParsers = append(newParsers, stackNode)
	}

	return newParsers
}

func getRuleNode(s *GLRState, rl Rule, kids []*ParseNode) *ParseNode {
	key := fmt.Sprintf("%s:%v", rl.Nonterminal, kids)
	if node, exists := s.ruleNodes[key]; exists {
		return node
	}

	numTerms := 0
	for _, kid := range kids {
		numTerms += kid.numTerms
	}
	node := &ParseNode{
		symbol:   rl.Nonterminal,
		children: kids,
		startPos: kids[0].startPos,
		endPos:   kids[len(kids)-1].endPos,
		numTerms: numTerms,
	}
	s.ruleNodes[key] = node
	return node
}

func getSymbolNode(s *GLRState, n *ParseNode) *ParseNode {
	for _, node := range s.symbolNodes {
		if node.symbol == n.symbol && node.startPos == n.startPos && node.endPos == n.endPos {
			return node
		}
	}
	s.symbolNodes = append(s.symbolNodes, n)
	return n
}

func addAlternative(s *GLRState, old *ParseNode, new *ParseNode) *ParseNode {
	debugf("adding alternative with old.symbol: %q, new.symbol: %q, old == new: %t\n", old.symbol, new.symbol, old == new)
	if parseNodesEqual(old, new) {
		return old
	}

	// Create or update ambiguous node
	var ambiguous *ParseNode
	if old.symbol != new.symbol || old.startPos != new.startPos || old.endPos != new.endPos {
		return old
	}

	ambiguous = &ParseNode{
		symbol:   old.symbol,
		children: append([]*ParseNode{old}, new),
		startPos: old.startPos,
		endPos:   old.endPos,
		numTerms: max(old.numTerms, new.numTerms),
	}

	// Update references
	for i, node := range s.symbolNodes {
		if node == old {
			s.symbolNodes[i] = ambiguous
		}
	}
	return ambiguous
}

func parseNodesEqual(n1, n2 *ParseNode) bool {
	return n1.symbol == n2.symbol &&
		n1.value == n2.value &&
		len(n1.children) == len(n2.children) &&
		n1.startPos == n2.startPos &&
		n1.endPos == n2.endPos &&
		n1.numTerms == n2.numTerms
}

func getStackNode(ps []*StackNode, state int) *StackNode {
	for _, parser := range ps {
		if parser.state == state {
			return parser
		}
	}
	return nil
}

func contains(slice []*StackNode, item *StackNode) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// func parseTreeEqual(n1, n2 *ParseNode) {
// 	visited := make(map[*ParseNode]bool)
// 	parseTreeEqualHelper(n1, n2, visited)
// }

// func parseTreeEqualHelper(n1, n2 *ParseNode, visited map[*ParseNode]bool) {
// 	if (n1 == nil && n2 == nil) || (visited[n1] && visited[n2]) {
// 		return
// 	}
// 	visited[n1] = true
// 	visited[n2] = true

// 	// Verify node positions are consistent
// 	if len(node.children) > 0 {
// 		firstChild := node.children[0]
// 		lastChild := node.children[len(node.children)-1]

// 		if node.startPos != firstChild.startPos {
// 			t.Errorf("Node start position inconsistent: node=%d, firstChild=%d",
// 				node.startPos, firstChild.startPos)
// 		}

// 		if node.endPos != lastChild.endPos {
// 			t.Errorf("Node end position inconsistent: node=%d, lastChild=%d",
// 				node.endPos, lastChild.endPos)
// 		}
// 	}

// 	// Recursively verify children
// 	for _, child := range node.children {
// 		verifyParseTreeHelper(t, child, visited)
// 	}
// }
