package glr

import (
	"fmt"
	"sort"
)

// Debug flags
var glrDebug = true
var glrVerbose = true

// SetDebug toggles debug logging
func SetDebug(enabled bool) {
	glrDebug = enabled
}

// SetVerbose toggles verbose logging
func SetVerbose(enabled bool) {
	glrVerbose = enabled
}

// debugLog prints debug messages if debug is enabled
func debugLog(format string, args ...interface{}) {
	if glrDebug {
		fmt.Printf(format, args...)
	}
}

// verboseLog prints verbose messages if verbose is enabled
func verboseLog(format string, args ...interface{}) {
	if glrVerbose {
		fmt.Printf(format, args...)
	}
}

func printNodeTree(n *ParseNode, spaces string) {
	fmt.Printf("%s%d: [%d, %d]: symbol: %q, value: %#v\n", spaces, n.numTerms, n.startPos, n.endPos, n.symbol, n.value)
	for _, child := range n.children {
		printNodeTree(child, spaces+"  ")
	}
}

// outputParserName returns a string description of a parser
func outputParserName(p *StackNode) string {
	return fmt.Sprintf("parser with state %d", p.state)
}

// printParser prints the parser state and backlinks recursively
func printParser(p *StackNode, parsersAfter string) {
	if len(p.backlinks) > 0 {
		for _, backlink := range p.backlinks {
			backStackNode := backlink.stackNode
			if p.state == backStackNode.state {
				fmt.Printf("* %d %s\n", p.state, parsersAfter)
			} else {
				printParser(backStackNode, fmt.Sprintf("- %s - %d %s",
					backlink.node.symbol, p.state, parsersAfter))
			}
		}
	} else {
		fmt.Printf("%d %s\n", p.state, parsersAfter)
	}
}

// printActiveParser prints all active parsers
func printActiveParsers(ps []*StackNode) {
	debugLog("%d active parsers with states: %v\n",
		len(ps), mapStates(ps))
	for _, parser := range ps {
		printParser(parser, "")
		debugLog("\n")
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
	parsersToShift   [][2]interface{}
	acceptingParsers []*StackNode
	ruleNodes        map[string]*ParseNode
	symbolNodes      []*ParseNode
	lookahead        *ParseNode
	debug            bool
}

// Parse implements GLR parsing algorithm
func Parse(rls []Rule, sts []ParseState, input string) ([]*ParseNode, error) {
	debugLog("\nstarting GLR parse\n")

	lexer := newLexer(input)
	s := &GLRState{
		ruleNodes: make(map[string]*ParseNode),
		debug:     false,
	}

	// Initialize with start state
	firstParser := &StackNode{state: 0, backlinks: nil}
	s.activeParsers = []*StackNode{firstParser}
	debugLog("initialized with start state 0\n")

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
		fmt.Printf("result[%d]\n", i)
		printNodeTree(result, "")
	}
	return results, nil
}

func parseSymbol(rls []Rule, sts []ParseState, s *GLRState, tok *ParseNode) bool {
	debugLog(fmt.Sprintf("parsing lookahead symbol: %q with value: %#v\n", tok.symbol, tok.value))

	if glrDebug {
		printActiveParsers(s.activeParsers)
	}

	s.lookahead = tok
	s.initialParsers = s.activeParsers
	s.parsersToAct = s.activeParsers
	s.parsersToShift = nil
	s.ruleNodes = make(map[string]*ParseNode)
	s.symbolNodes = nil

	// Process all parsers
	for len(s.parsersToAct) > 0 {
		p := s.parsersToAct[0]
		s.parsersToAct = s.parsersToAct[1:]
		s.activeParser = p
		actor(rls, sts, s, p)
	}

	// Perform shifts if any available
	if len(s.parsersToShift) == 0 && (len(s.acceptingParsers) == 0 || tok.symbol != "$end") {
		// No valid actions found
		fmt.Printf("accepting with no valid actions found\n")
		s.acceptingParsers = append(s.acceptingParsers, s.initialParsers...)
		return false
	}
	s.activeParsers = shifter(s)
	return true
}

func actor(rls []Rule, sts []ParseState, s *GLRState, p *StackNode) {
	as := append(sts[p.state].Actions[s.lookahead.symbol], sts[p.state].Actions["."]...)
	for i, a := range as {
		fmt.Printf("looking at action for p.state: %d, s.lookahead.symbol: %q, actions[%d]: %#v\n", p.state, s.lookahead.symbol, i, a)

		switch a.Action {
		case "shift":
			fmt.Printf("shifting to state: %d\n", a.State)
			s.parsersToShift = append(s.parsersToShift, [2]interface{}{p, a.State})

		case "reduce":
			r := rls[a.Rule]
			fmt.Printf("reducing with rule: %#v\n", r)
			doReductions(rls, sts, s, p, r, len(r.RHS), nil, nil, true)

		case "accept":
			fmt.Printf("accepting\n")
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
	printParser(p, "")

	gotoState, ok := sts[p.state].Gotos[r.Nonterminal]
	if !ok {
		return
	}

	ruleNode := getRuleNode(s, r, kids)
	printNodeTree(ruleNode, "")
	stackNode := getStackNode(s.activeParsers, gotoState)

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
	fmt.Printf("shifter(s)\n")
	var newParsers []*StackNode

	for _, pair := range s.parsersToShift {
		parser := pair[0].(*StackNode)
		newState := pair[1].(int)

		newLink := &StackLink{stackNode: parser, node: s.lookahead}
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
	if old == new {
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
