package glr

import (
	"fmt"
	"go/scanner"
	"go/token"
	"sort"
)

// Debug flags
var glrDebug = true

// var glrDebug = false

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

type Rules struct {
	Items []Rule
}

type Rule struct {
	Nonterminal string
	RHS         []string
}

type ParseStates struct {
	Items []ParseState
}

type ParseState struct {
	Actions map[string][]Action
	Gotos   map[string]int
}

type Action struct {
	Type  string
	State int
	Rule  int
}

type ParseNode struct {
	Symbol   string
	Children []*ParseNode

	value    interface{}
	startPos int
	endPos   int
	numTerms int
	isAlt    bool
}

func printNodeTree(n *ParseNode, spaces string) {
	alt := ""
	if n.isAlt {
		alt = "ALT "
	}
	debugf("%s%d: [%d, %d]: %ssymbol: %q, value: %#v\n", spaces, n.numTerms, n.startPos, n.endPos, alt, n.Symbol, n.value)
	for _, child := range n.Children {
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
			printParser(backStackNode, fmt.Sprintf("- %s - %d %s%s", backlink.node.Symbol, p.state, parsersAfter, shift), nil)
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

type Lexer interface {
	NextToken(int) (string, any, bool)
	Error(string)
}

func NewLexerScanner(l Lexer, input string) scanner.Scanner {
	inputBs := []byte(input)
	var s scanner.Scanner
	fset := token.NewFileSet()                          // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(inputBs)) // register input "file"
	errHandler := func(pos token.Position, msg string) { l.Error(fmt.Sprintf("error at position %v: %v", pos, msg)) }
	s.Init(file, inputBs, errHandler, 0) //, scanner.ScanComments)
	return s
}

// Parse implements GLR parsing algorithm
func Parse(rls *Rules, sts *ParseStates, l Lexer) ([]*ParseNode, error) {
	debugf("\nstarting GLR parse\n")

	s := &GLRState{
		ruleNodes: make(map[string]*ParseNode),
		debug:     false,
	}

	// Initialize with start state
	firstParser := &StackNode{state: 0, backlinks: nil}
	s.activeParsers = []*StackNode{firstParser}
	debugf("initialized with start state 0\n")

	pos := 0
	for {
		// Create parse node for token
		sym, val, hasMore := l.NextToken(pos)
		term := &ParseNode{
			Symbol:   sym,
			value:    val,
			startPos: pos,
			endPos:   pos + 1,
			numTerms: 1,
		}
		pos++
		if sym == "" {
			break
		}

		if !parseSymbol(rls, sts, s, term) {
			continue // Skip invalid tokens
		}

		if !hasMore {
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

	debugf("results\n")
	for i, result := range results {
		debugf("result[%d]\n", i)
		printNodeTree(result, "")
	}
	return results, nil
}

func parseSymbol(rls *Rules, sts *ParseStates, s *GLRState, term *ParseNode) bool {
	debugf(fmt.Sprintf("\nparsing term: %#v\n", term))
	printAllParsers(s, nil)

	s.lookahead = term
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
	if len(s.parsersToShift) == 0 && (len(s.acceptingParsers) == 0 || term.Symbol != "$end") {
		// No valid actions found
		debugf("  accepting with no valid actions found\n")
		s.acceptingParsers = append(s.acceptingParsers, s.initialParsers...)
		return false
	}
	s.activeParsers = shifter(s)
	return true
}

func actor(rls *Rules, sts *ParseStates, s *GLRState, p *StackNode) {
	as := append(sts.Items[p.state].Actions[s.lookahead.Symbol], sts.Items[p.state].Actions["."]...)
	debugf("found %d actions for p.state: %d, s.lookahead.symbol: %q\n", len(as), p.state, s.lookahead.Symbol)
	for i, a := range as {
		debugf("  looking at action for p.state: %d, s.lookahead.symbol: %q, actions[%d]: %#v\n", p.state, s.lookahead.Symbol, i, a)

		switch a.Type {
		case "shift":
			debugf("  shifting to state: %d\n", a.State)
			s.parsersToShift = append(s.parsersToShift, p)
			s.statesToShiftTo = append(s.statesToShiftTo, a.State)

		case "reduce":
			r := &(rls.Items[a.Rule])
			debugf("  reducing with rule: %#v\n", r)
			doReductions(rls, sts, s, p, r, len(r.RHS), nil, nil, true)

		case "accept":
			debugf("  accepting\n")
			s.acceptingParsers = append(s.acceptingParsers, p)
		}
	}
}

func doReductions(rls *Rules, sts *ParseStates, s *GLRState, p *StackNode, r *Rule, length int, kids []*ParseNode, linkToSee *StackLink, linkSeen bool) {
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

func reducer(rls *Rules, sts *ParseStates, s *GLRState, p *StackNode, r *Rule, kids []*ParseNode) {
	printAllParsers(s, p)

	gotoState, ok := sts.Items[p.state].Gotos[r.Nonterminal]
	if !ok {
		debugf("sts[%v].Gotos[%q]: %v, %t\n", p.state, r.Nonterminal, gotoState, ok)
		return
	}

	ruleNode := getRuleNode(s, r, kids)
	stackNode := getStackNode(s.activeParsers, gotoState)
	debugln("  gotoState", gotoState, "stackNode", stackNode)
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
			actions := sts.Items[otherParser.state].Actions[s.lookahead.Symbol]
			for _, action := range actions {
				if action.Type == "reduce" {
					otherRule := &(rls.Items[action.Rule])
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

func getRuleNode(s *GLRState, rl *Rule, kids []*ParseNode) *ParseNode {
	debugf("getRuleNode(s, rl %#v, len(kids): %d)\n", rl, len(kids))
	key := fmt.Sprintf("%s:%v", rl.Nonterminal, kids)
	if node, exists := s.ruleNodes[key]; exists {
		return node
	}

	numTerms := 0
	for _, kid := range kids {
		numTerms += kid.numTerms
	}
	node := &ParseNode{
		Symbol:   rl.Nonterminal,
		Children: kids,
		numTerms: numTerms,
	}
	if len(kids) > 0 {
		node.startPos = kids[0].startPos
		node.endPos = kids[len(kids)-1].endPos
	}
	s.ruleNodes[key] = node
	return node
}

func getSymbolNode(s *GLRState, n *ParseNode) *ParseNode {
	for _, node := range s.symbolNodes {
		if node.Symbol == n.Symbol && node.startPos == n.startPos && node.endPos == n.endPos {
			return node
		}
	}
	s.symbolNodes = append(s.symbolNodes, n)
	return n
}

func addAlternative(s *GLRState, old *ParseNode, new *ParseNode) *ParseNode {
	debugf("adding alternative with old.symbol: %q, old.isAlt: %t, new.symbol: %q, new.isAlt: %t, old == new: %t\n", old.Symbol, old.isAlt, new.Symbol, new.isAlt, old == new)
	if parseNodesEqual(old, new) {
		return old
	}

	if old.isAlt {
		old.Children = append(old.Children, new)
		return old
	}

	ambiguous := &ParseNode{
		Symbol:   old.Symbol,
		Children: []*ParseNode{old, new},
		startPos: old.startPos,
		endPos:   old.endPos,
		numTerms: max(old.numTerms, new.numTerms),
		isAlt:    true,
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
	if n1.Symbol != n2.Symbol ||
		n1.value != n2.value ||
		len(n1.Children) != len(n2.Children) ||
		n1.startPos != n2.startPos ||
		n1.endPos != n2.endPos ||
		n1.numTerms != n2.numTerms {
		return false
	}
	for i, n1C := range n1.Children {
		if !parseNodesEqual(n1C, n2.Children[i]) {
			return false
		}
	}
	return true
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
