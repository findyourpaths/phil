package glr

import (
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/k0kubun/pp/v3"
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

func (r Rule) String() string {
	return fmt.Sprintf("%s -> %q", r.Nonterminal, strings.Join(r.RHS, " "))
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
	Type    string
	StateID int
	RuleID  int
}

type ParseNode struct {
	Symbol   string
	Term     string
	Children []*ParseNode

	startPos int
	endPos   int
	isAlt    bool
	ruleID   int
	score    *parseNodeScore
}

type parseNodeScore struct {
	numTerms int
	size     int
	depths   int
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
		rule = "(rule " + strconv.Itoa(n.ruleID) + ")"
	}
	return fmt.Sprintf("%s [%d:%d] %s%s %s", n.score, n.startPos, n.endPos, alt, n.Symbol, rule) //n.Value)
}

var space = "\u2758 "

func GetParseNodeValue(g *Grammar, n *ParseNode, spaces string) (any, error) {
	debugf("%sgetting value for: %s\n", spaces, n)
	if n.isAlt {
		for i, c := range n.Children {
			r, err := GetParseNodeValue(g, c, spaces+space)
			if err == nil {
				debugf("%sreturning alternative %d\n", spaces, i)
				return r, nil
			}
		}
		debugf("%sfailed to find alternative\n", spaces)
		return nil, nil // fmt.Errorf("failed to find alternative")
	}
	if n.Term != "" {
		debugf("%sreturning term: %q\n", spaces, n.Term)
		return n.Term, nil
	}
	if n.ruleID == 0 {
		debugf("%sreturning empty rule result: %q\n", spaces, n.Term)
		return "", nil
	}

	fn := reflect.ValueOf(g.Actions.Items[n.ruleID])
	args := make([]reflect.Value, len(n.Children))
	for i, c := range n.Children {
		val, err := GetParseNodeValue(g, c, spaces+space)
		if err != nil {
			return nil, err
		}
		if val == nil {
			return nil, nil
		}
		args[i] = reflect.ValueOf(val)
		debugf("%sgot rule %d args[%d]: %#v\n", spaces, n.ruleID, i, args[i])
	}

	debugf("%scalling rule %d fn with %d args\n", spaces, n.ruleID, len(args))
	var err error
	var r any
	func() {
		defer func() {
			if e := recover(); e != nil {
				err = errors.New(e.(string))
				debugf("%sgot err: %v", spaces, err)
			}
		}()
		r = fn.Call(args)[0].Interface()
	}()
	if spaces == "" {
		pp.Default.SetColoringEnabled(false)
		debugf("%sreturning computed value:\n%s\nerr: %v\n", spaces, pp.Sprint(r), err)
	}
	return r, err
}

func setNodeChildrenAndScore(n *ParseNode, children []*ParseNode) {
	if n.isAlt {
		n.Children = sortAlternatives(children)
		n.score = n.Children[0].score
		for _, c := range n.Children {
			if n.startPos > c.startPos {
				n.startPos = c.startPos
			}
			if n.endPos < c.endPos {
				n.endPos = c.endPos
			}
		}
		return
	}

	n.Children = children
	for _, c := range n.Children {
		if n.startPos > c.startPos {
			n.startPos = c.startPos
		}
		if n.endPos < c.endPos {
			n.endPos = c.endPos
		}
	}

	sc := &parseNodeScore{}
	sc.size = 1
	// sc.depths = depth
	if n.Term != "" {
		sc.numTerms = 1
	} else {
		for _, c := range n.Children {
			csc := c.score
			sc.numTerms += csc.numTerms
			sc.size += csc.size
			sc.depths += csc.depths
		}
	}
	sc.depths += sc.size
	n.score = sc
}

// We prefer parse trees covering the most terms, with the fewest total
// children, but most deep nodes.
//
// For the last one, we use depth-based weighting to prefer "bushy" parse trees
// at the lower levels and more streamlined structures near the root.
//
// Assign weights to nodes based on their depth in the tree. Nodes closer to the
// root receive lower weights, while nodes further down receive higher weights.
//
// This encourages branching at lower levels because each new branch deeper in
// the tree contributes more to the overall score.
func sortAlternatives(ns []*ParseNode) []*ParseNode {
	sort.Slice(ns, func(i, j int) bool {
		return (ns[i].score.numTerms > ns[j].score.numTerms) ||

			(ns[i].score.numTerms == ns[j].score.numTerms &&
				ns[i].score.size < ns[j].score.size) ||

			(ns[i].score.numTerms == ns[j].score.numTerms &&
				ns[i].score.size == ns[j].score.size &&
				ns[i].score.depths > ns[j].score.depths)
	})
	return ns
}

func (s parseNodeScore) String() string {
	return fmt.Sprintf("[%2d, %2d, %3d]", s.numTerms, s.size, s.depths)
}

func printNodes(ns []*ParseNode, spaces string) {
	if !DoDebug {
		return
	}
	for _, n := range ns {
		debugf("%s%s\n", spaces, n)
	}
}

func printNodeRoot(n *ParseNode) {
	if !DoDebug {
		return
	}
	printNodeTree(n, map[*ParseNode]bool{}, "  ")
}

// var hideCycles = true
var hideCycles = false

func printNodeTree(n *ParseNode, seen map[*ParseNode]bool, spaces string) {
	if hideCycles && seen[n] {
		debugf("%s-> %s\n", spaces, n.Symbol)
		return
	}
	seen[n] = true
	debugf("%s%s\n", spaces, n)
	for _, child := range n.Children {
		printNodeTree(child, seen, spaces+space)
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
			printParser(backStackNode, fmt.Sprintf("-- %s -- %d %s%s", backlink.node, p.state, parsersAfter, shift), nil)
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
	// ruleNodes        map[string]*ParseNode
	ruleNodes      []*ParseNode
	symbolNodes    []*ParseNode
	partialResults []*ParseNode
	lookahead      *ParseNode
	debug          bool
	position       int
}

type Lexer interface {
	NextToken(int) (string, string, bool)
	Error(string)
}

func NewLexerScanner(l Lexer, input string) scanner.Scanner {
	input = strings.Replace(input, "//", "/ /", -1)
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
		// ruleNodes: map[string]*ParseNode{},
		// ruleNodes: []*ParseNode{},
		debug: false,
	}

	// Initialize with start state
	firstParser := &StackNode{state: 0, backlinks: nil}
	s.activeParsers = []*StackNode{firstParser}
	debugf("initialized with start state 0\n")

	for {
		debugf("\n")

		// Create parse node for token
		sym, val, hasMore := l.NextToken(s.position)
		term := &ParseNode{
			Symbol:   sym,
			Term:     val,
			startPos: s.position,
			endPos:   s.position + 1,
		}
		setNodeChildrenAndScore(term, nil)
		s.position++
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

	aps := []*ParseNode{}
	for _, parser := range s.acceptingParsers {
		if len(parser.backlinks) > 0 {
			aps = append(aps, parser.backlinks[0].node)
		}
	}
	prs := []*ParseNode{}
	for _, pr := range s.partialResults {
		found := false
		for _, ap := range aps {
			if pr == ap {
				found = true
				break
			}
			for _, apc := range ap.Children {
				if pr == apc {
					found = true
					break
				}
			}
		}
		if !found {
			prs = append(prs, pr)
		}
	}

	debugf("ended with %d accepting parsers and %d partial results\n", len(aps), len(prs))
	alts := append(aps, prs...)
	if len(alts) == 0 {
		debugf("found no results\n")
		return nil, nil
	}
	all := &ParseNode{isAlt: true}
	setNodeChildrenAndScore(all, alts)
	rs := all.Children
	debugf("found %d results\n", len(rs))
	for i, r := range rs {
		debugf("result[%d]\n", i)
		printNodeRoot(r)
	}
	return rs, nil
}

func parseSymbol(g *Grammar, s *GLRState, term *ParseNode) bool {
	debugf(fmt.Sprintf("parsing term: %#v\n", term))
	printAllParsers(s, nil)

	s.lookahead = term
	s.initialParsers = s.activeParsers
	s.parsersToAct = s.activeParsers
	s.parsersToShift = nil
	s.statesToShiftTo = nil
	// s.ruleNodes = make(map[string]*ParseNode)
	s.ruleNodes = nil
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
	if len(s.parsersToShift) == 0 && term.Symbol != "$end" {
		debugf("  ending with no valid actions found\n")
		return false
	}

	s.activeParsers = shifter(s)
	return true
}

func actor(g *Grammar, s *GLRState, p *StackNode) {
	as := getActions(g, p.state, s.lookahead.Symbol)
	debugf("found %d actions for p.state: %d, s.lookahead.symbol: %q\n", len(as), p.state, s.lookahead.Symbol)
	for i, a := range as {
		debugf("  looking at action for p.state: %d, s.lookahead.symbol: %q, actions[%d]: %#v\n", p.state, s.lookahead.Symbol, i, a)

		switch a.Type {
		case "shift":
			debugf("  shifting to state: %d\n", a.StateID)
			s.parsersToShift = append(s.parsersToShift, p)
			s.statesToShiftTo = append(s.statesToShiftTo, a.StateID)

		case "reduce":
			r := &(g.Rules.Items[a.RuleID])
			debugf("  reducing with rule %d: %s\n", a.RuleID, r)
			doReductions(g, s, p, a.RuleID, len(r.RHS), nil, nil, true)

		case "accept":
			debugf("  accepting\n")
			s.acceptingParsers = append(s.acceptingParsers, p)
		}
	}
}

func getActions(g *Grammar, state int, sym string) []Action {
	return append(g.States.Items[state].Actions[sym], g.States.Items[state].Actions["."]...)
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
	rl := &(g.Rules.Items[rlID])
	debugf("reducer with rule %d: %s\n", rlID, rl)
	printAllParsers(s, p)

	gotoState, ok := g.States.Items[p.state].Gotos[rl.Nonterminal]
	if !ok {
		debugf("sts[%v].Gotos[%q]: %v, %t\n", p.state, rl.Nonterminal, gotoState, ok)
		return
	}

	ruleNode := getRuleNode(g, s, rlID, kids)
	stackNode := getStackNode(s.activeParsers, gotoState)
	// debugf("  gotoState: %d, stackNode: %v\n", gotoState, stackNode)

	if stackNode == nil {
		// Create new parser state
		nonterm := getSymbolNode(s, ruleNode)
		debugf("  reducing to symbol node nonterm and going to new state: %d\n", gotoState)
		printNodeRoot(nonterm)
		stackNode = &StackNode{
			state:     gotoState,
			backlinks: []*StackLink{{stackNode: p, node: nonterm}}}
		s.activeParsers = append(s.activeParsers, stackNode)
		s.parsersToAct = append(s.parsersToAct, stackNode)
		if rl.Type == g.Rules.Items[1].Type {
			debugf("  saving current partial result with node of same type as root: %s\n", rl.Type)
			s.partialResults = append(s.partialResults, nonterm)
		}
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
	nonterm := getSymbolNode(s, ruleNode)
	debugf("  reducing new path nonterm and going to existing state: %d\n", gotoState)
	printNodeRoot(nonterm)

	newLink := &StackLink{stackNode: p, node: nonterm}
	stackNode.backlinks = append(stackNode.backlinks, newLink)
	for i, bl := range stackNode.backlinks {
		debugf("  stack node backlink %d to: %s\n", i, bl.node)
	}

	// Process additional reductions
	debugf("  looking at other parsers\n")
	for _, otherp := range s.activeParsers {
		if contains(s.parsersToAct, otherp) {
			continue
		}
		for _, a := range getActions(g, otherp.state, s.lookahead.Symbol) {
			if a.Type != "reduce" {
				continue
			}
			rhsLen := len(g.Rules.Items[a.RuleID].RHS)
			doReductions(g, s, otherp, a.RuleID, rhsLen, nil, newLink, false)
		}
	}
}

func shifter(s *GLRState) []*StackNode {
	var rs []*StackNode

	for i, p := range s.parsersToShift {
		newState := s.statesToShiftTo[i]
		newLink := &StackLink{stackNode: p, node: s.lookahead}
		r := getStackNode(rs, newState)

		if r != nil {
			r.backlinks = append(r.backlinks, newLink)
			continue
		}

		r = &StackNode{state: newState, backlinks: []*StackLink{newLink}}
		rs = append(rs, r)
	}

	return rs
}

func getRuleNode(g *Grammar, s *GLRState, rlID int, children []*ParseNode) *ParseNode {
	debugf("  getting rule node for rule: %d and %d kids\n", rlID, len(children))
	rl := &(g.Rules.Items[rlID])
	for _, old := range s.ruleNodes {
		if rl.Nonterminal != old.Symbol {
			continue
		}
		if len(old.Children) != len(children) {
			continue
		}
		sameCs := true
		for i, oldC := range old.Children {
			if oldC != children[i] {
				sameCs = false
				break
			}
		}
		if sameCs {
			debugf("  got old rule nonterm: %s\n", old)
			return old
		}
	}

	// key := fmt.Sprintf("%s:%v", rl.Nonterminal, kids)
	// debugf("  key: %q\n", key)
	// if node, exists := s.ruleNodes[key]; exists {
	// 	return node
	// }

	r := &ParseNode{
		Symbol: rl.Nonterminal,
		ruleID: rlID,
	}
	if len(children) > 0 {
		r.startPos = children[0].startPos
		r.endPos = children[len(children)-1].endPos
	} else {
		r.startPos = s.position
		r.endPos = s.position
	}
	setNodeChildrenAndScore(r, children)
	// s.ruleNodes[key] = r

	debugf("  got new rule node nonterm: %s\n", r)
	printNodeRoot(r)
	return r
}

func getSymbolNode(s *GLRState, n *ParseNode) *ParseNode {
	debugf("  maybe getting old symbol node for: %s\n", n)
	for _, old := range s.symbolNodes {
		if old.Symbol == n.Symbol && old.startPos == n.startPos && old.endPos == n.endPos {
			debugf("  returning old: %s\n", old)
			return old
		}
	}
	debugf("  returning new\n")
	s.symbolNodes = append(s.symbolNodes, n)
	return n
}

func addAlternative(s *GLRState, old *ParseNode, new *ParseNode) *ParseNode {
	debugf("maybe adding alternative with old.symbol: %q, old.isAlt: %t, new.symbol: %q, new.isAlt: %t\n", old.Symbol, old.isAlt, new.Symbol, new.isAlt)
	debugf("old\n")
	printNodeRoot(old)
	debugf("new\n")
	printNodeRoot(new)

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
		setNodeChildrenAndScore(old, append(old.Children, new))
		debugf("returning\n")
		printNodeRoot(old)
		return old
	}

	// This is ugly and not functional programming, but we ran into a bug where we
	// didn't see the node that already incorporated old. So leave old in the tree
	// but promote it to an alternative and create "newOld", a copy of old as the
	// new child of the promoted old, along with new.
	newOld := &ParseNode{
		Symbol:   old.Symbol,
		startPos: old.startPos,
		endPos:   old.endPos,
		ruleID:   old.ruleID,
	}
	setNodeChildrenAndScore(newOld, old.Children)

	old.isAlt = true
	old.ruleID = 0
	setNodeChildrenAndScore(old, []*ParseNode{newOld, new})

	debugf("adding alternative by merging new and old\n")
	debugf("returning\n")
	printNodeRoot(old)
	return old
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
