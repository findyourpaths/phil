package glr

import (
	"fmt"
	"go/scanner"
	"go/token"
	"math"
	"reflect"
	"sort"
	"strconv"
)

// Debug flags
// var DoDebug = true

var DoDebug = false

// SetDebug toggles debug logging
func SetDebug(enabled bool) {
	DoDebug = enabled
}

// debugf prints debug messages if debug is enabled
func debugf(format string, args ...any) {
	if DoDebug {
		fmt.Printf(format, args...)
	}
}

// debugf prints debug messages if debug is enabled
func debugln(args ...any) {
	if DoDebug {
		fmt.Println(args...)
	}
}

type Grammar struct {
	Rules   *Rules
	Actions *SemanticActions
	States  *ParseStates
}

type Rules struct {
	Items []Rule
}

type Rule struct {
	Nonterminal string
	RHS         []string
	Type        string
}

type SemanticActions struct {
	Items []any
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
	Term     string
	Children []*ParseNode

	startPos int
	endPos   int
	// numNonterms int
	numTerms int
	isAlt    bool
	ruleID   int
}

func (n ParseNode) size() int {
	if n.isAlt {
		r := math.MaxInt32
		for _, c := range n.Children {
			r = min(r, c.size())
		}
		return r
	}

	r := 1
	for _, c := range n.Children {
		r += c.size()
	}
	return r
}

func (n ParseNode) String() string {
	alt := ""
	if n.isAlt {
		alt = "ALT "
	}
	var rule string
	if n.Term != "" {
		rule = fmt.Sprintf("%q", n.Term)
	} else {
		rule = strconv.Itoa(n.ruleID)
	}
	return fmt.Sprintf("%d.%d: [%d, %d]: %s%s: %s", n.size(), n.numTerms, n.startPos, n.endPos, alt, n.Symbol, rule) //n.Value)
}

func GetParseNodeValue(g *Grammar, n *ParseNode, spaces string) any {
	debugf("%sgetting value for: %s\n", spaces, n)
	if n.isAlt {
		r := GetParseNodeValue(g, n.Children[0], spaces+"  ")
		debugf("%sreturning first alternative\n", spaces)
		return r
	}
	if n.Term != "" {
		debugf("%sreturning term: %q\n", spaces, n.Term)
		return n.Term
	}
	if n.ruleID == 0 {
		debugf("%sreturning empty rule result: %q\n", spaces, n.Term)
		return ""
	}

	fn := reflect.ValueOf(g.Actions.Items[n.ruleID])
	cs := n.Children
	args := make([]reflect.Value, len(cs))
	for i, c := range cs {
		args[i] = reflect.ValueOf(GetParseNodeValue(g, c, spaces+"  "))
	}
	r := fn.Call(args)[0].Interface()
	debugf("%sreturning computed value: %#v\n", spaces, r)
	return r
}

func printNodes(ns []*ParseNode, spaces string) {
	if !DoDebug {
		return
	}
	for _, n := range ns {
		debugf("%s%s\n", spaces, n)
	}
}

func printNodeTree(n *ParseNode, spaces string) {
	if !DoDebug {
		return
	}
	debugf("%s%s\n", spaces, n)
	for _, child := range n.Children {
		printNodeTree(child, spaces+"  ")
	}
}

// printActiveParser prints all active parsers
func printAllParsers(s *GLRState, p *StackNode) {
	if !DoDebug {
		return
	}
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
	NextToken(int) (string, string, bool)
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
func Parse(g *Grammar, l Lexer) ([]*ParseNode, error) {
	debugf("\nstarting GLR parse\n")

	s := &GLRState{
		ruleNodes: map[string]*ParseNode{},
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
			Term:     val,
			startPos: pos,
			endPos:   pos + 1,
			numTerms: 1,
		}
		pos++
		if sym == "" {
			break
		}

		if !parseSymbol(g, s, term) {
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
	results = sortParseNodes(results)

	debugf("found %d results\n", len(results))
	for i, result := range results {
		debugf("result[%d]\n", i)
		printNodeTree(result, "")
	}
	return results, nil
}

func parseSymbol(g *Grammar, s *GLRState, term *ParseNode) bool {
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
		debugf("symbol nodes\n")
		printNodes(s.symbolNodes, "  ")
		actor(g, s, p)
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

func actor(g *Grammar, s *GLRState, p *StackNode) {
	as := append(g.States.Items[p.state].Actions[s.lookahead.Symbol], g.States.Items[p.state].Actions["."]...)
	debugf("found %d actions for p.state: %d, s.lookahead.symbol: %q\n", len(as), p.state, s.lookahead.Symbol)
	for i, a := range as {
		debugf("  looking at action for p.state: %d, s.lookahead.symbol: %q, actions[%d]: %#v\n", p.state, s.lookahead.Symbol, i, a)

		switch a.Type {
		case "shift":
			debugf("  shifting to state: %d\n", a.State)
			s.parsersToShift = append(s.parsersToShift, p)
			s.statesToShiftTo = append(s.statesToShiftTo, a.State)

		case "reduce":
			r := &(g.Rules.Items[a.Rule])
			debugf("  reducing with rule: %#v\n", r)
			doReductions(g, s, p, a.Rule, len(r.RHS), nil, nil, true)

		case "accept":
			debugf("  accepting\n")
			s.acceptingParsers = append(s.acceptingParsers, p)
		}
	}
}

func doReductions(g *Grammar, s *GLRState, p *StackNode, rlID int, length int, kids []*ParseNode, linkToSee *StackLink, linkSeen bool) {
	if length == 0 {
		if linkSeen {
			reducer(g, s, p, rlID, kids)
		}
		return
	}

	for _, link := range p.backlinks {
		newLinkSeen := linkSeen || link == linkToSee
		doReductions(g, s, link.stackNode, rlID, length-1, append([]*ParseNode{link.node}, kids...), linkToSee, newLinkSeen)
	}
}

func reducer(g *Grammar, s *GLRState, p *StackNode, rlID int, kids []*ParseNode) {
	printAllParsers(s, p)

	rl := &(g.Rules.Items[rlID])
	gotoState, ok := g.States.Items[p.state].Gotos[rl.Nonterminal]
	if !ok {
		debugf("sts[%v].Gotos[%q]: %v, %t\n", p.state, rl.Nonterminal, gotoState, ok)
		return
	}

	ruleNode := getRuleNode(g, s, rlID, kids)
	stackNode := getStackNode(s.activeParsers, gotoState)
	debugf("  gotoState: %d, stackNode: %v\n", gotoState, stackNode)
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
			actions := g.States.Items[otherParser.state].Actions[s.lookahead.Symbol]
			for _, action := range actions {
				if action.Type == "reduce" {
					otherRule := &(g.Rules.Items[action.Rule])
					doReductions(g, s, otherParser, action.Rule, len(otherRule.RHS), nil, newLink, false)
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

func getRuleNode(g *Grammar, s *GLRState, rlID int, kids []*ParseNode) *ParseNode {
	debugf("getting rule node for rule: %d and %d kids\n", rlID, len(kids))
	rl := &(g.Rules.Items[rlID])
	key := fmt.Sprintf("%s:%v", rl.Nonterminal, kids)
	if node, exists := s.ruleNodes[key]; exists {
		return node
	}

	// numNonterms := 1
	numTerms := 0
	for _, kid := range kids {
		// numNonterms += kid.numNonterms
		numTerms += kid.numTerms
	}

	node := &ParseNode{
		Symbol: rl.Nonterminal,
		// Value:    val,
		Children: kids,
		// numNonterms: numNonterms,
		numTerms: numTerms,
		ruleID:   rlID,
	}
	if len(kids) > 0 {
		node.startPos = kids[0].startPos
		node.endPos = kids[len(kids)-1].endPos
	}
	s.ruleNodes[key] = node
	return node
}

func getSymbolNode(s *GLRState, n *ParseNode) *ParseNode {
	debugf("maybe getting old symbol node for: %s\n", n)
	for _, old := range s.symbolNodes {
		if old.Symbol == n.Symbol && old.startPos == n.startPos && old.endPos == n.endPos {
			debugf("returning old: %s\n", old)
			return old
		}
	}
	debugf("returning new\n")
	s.symbolNodes = append(s.symbolNodes, n)
	return n
}

func addAlternative(s *GLRState, old *ParseNode, new *ParseNode) *ParseNode {
	debugf("maybe adding alternative with old.symbol: %q, old.isAlt: %t, new.symbol: %q, new.isAlt: %t\n", old.Symbol, old.isAlt, new.Symbol, new.isAlt)
	debugf("old\n")
	printNodeTree(old, "")
	debugf("new\n")
	printNodeTree(new, "")

	if reflect.DeepEqual(old, new) {
		debugf("skipping adding alternative because new == old\n")
		return old
	}

	if old.isAlt {
		for i, alt := range old.Children {
			if reflect.DeepEqual(alt, new) {
				debugf("skipping adding alternative because new == old.Children[%d]\n", i)
				return old
			}
		}
		debugf("adding new as alternative to old\n")
		old.Children = sortParseNodes(append(old.Children, new))
		debugf("returning\n")
		printNodeTree(old, "")
		return old
	}

	// This is ugly and not functional programming, but we ran into a bug where we
	// didn't see the node that already incorporated old. So leave old in the tree
	// but promote it to an alternative and create "newOld", a copy of old as the
	// new child of the promoted old, along with new.
	newOld := &ParseNode{
		Symbol:   old.Symbol,
		Children: old.Children,
		startPos: old.startPos,
		endPos:   old.endPos,
		numTerms: old.numTerms,
		ruleID:   old.ruleID,
	}

	old.Children = sortParseNodes([]*ParseNode{newOld, new})
	old.isAlt = true
	old.ruleID = 0

	debugf("adding alternative by merging new and old\n")
	debugf("returning\n")
	printNodeTree(old, "")
	return old
}

func sortParseNodes(ns []*ParseNode) []*ParseNode {
	sort.Slice(ns, func(i, j int) bool {
		return ns[i].numTerms > ns[j].numTerms ||
			(ns[i].numTerms == ns[j].numTerms &&
				ns[i].size() < ns[j].size())
	})
	return ns
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
