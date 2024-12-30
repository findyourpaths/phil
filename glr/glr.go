package glr

import (
	"fmt"
)

// GLRParser represents a Generalized LR parser
type GLRParser struct {
	rules  []Rule
	states []ParseState
}

// ParseStack represents a parsing configuration
type ParseStack struct {
	stateStack []int
	nodeStack  []*ParseNode
	pos        int // Current position in input
}

// stackKey uniquely identifies a parser configuration
type stackKey struct {
	state int
	pos   int
}

// NewGLRParser creates a new GLR parser with the given rules and states
func NewGLRParser(rules []Rule, states []ParseState) (*GLRParser, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("no grammar rules provided")
	}
	if len(states) == 0 {
		return nil, fmt.Errorf("no parser states provided")
	}
	return &GLRParser{
		rules:  rules,
		states: states,
	}, nil
}

// mergeStacks combines equivalent stacks to prevent exponential growth
func mergeStacks(stacks []*ParseStack) []*ParseStack {
	stackMap := make(map[stackKey]*ParseStack)

	for _, stack := range stacks {
		if len(stack.stateStack) == 0 {
			continue
		}

		key := stackKey{
			state: stack.stateStack[len(stack.stateStack)-1],
			pos:   stack.pos,
		}

		if existing, ok := stackMap[key]; ok {
			// Merge parse trees if they represent the same derivation
			if len(stack.nodeStack) == len(existing.nodeStack) {
				continue // Keep existing stack
			}
		}
		stackMap[key] = stack
	}

	merged := make([]*ParseStack, 0, len(stackMap))
	for _, stack := range stackMap {
		merged = append(merged, stack)
	}
	return merged
}

// Parse parses the input string and returns the parse tree nodes
func (p *GLRParser) Parse(input string) ([]*ParseNode, error) {
	lexer := newLexer(input)
	var lval yySymType

	// Initialize with start state
	stacks := []*ParseStack{{
		stateStack: []int{0},
		nodeStack:  []*ParseNode{},
		pos:        0,
	}}

	// Track successful parses
	var successfulParses []*ParseNode

	for {
		// Get next token
		token := lexer.Lex(&lval)
		symbol := lexer.TokenVal
		isEOF := token < 0
		fmt.Println("token", token, "symbol", symbol, "isEOF", isEOF)

		var newStacks []*ParseStack

		// Process each stack configuration
		for i, stack := range stacks {
			if len(stack.stateStack) == 0 {
				continue
			}

			state := stack.stateStack[len(stack.stateStack)-1]
			fmt.Println("  reducing stack", "i", i, "state", state)

			// Process all possible reductions first
			for {
				dotActions := p.states[state].Actions["."]
				if len(dotActions) == 0 {
					break
				}

				var madeReduction bool
				for _, action := range dotActions {
					if action.Action != Reduce {
						continue
					}
					rule := p.rules[action.Target]
					rhsLen := len(rule.RHS)

					// Skip invalid reductions
					if rhsLen > len(stack.nodeStack) {
						continue
					}

					// Get children for reduction
					var children []*ParseNode
					if rhsLen > 0 {
						children = stack.nodeStack[len(stack.nodeStack)-rhsLen:]
					}

					// Create reduced node
					node := &ParseNode{
						symbol:   rule.Nonterminal,
						children: children,
					}
					if len(children) > 0 {
						node.startPos = children[0].startPos
						node.endPos = children[len(children)-1].endPos
					} else {
						node.startPos = stack.pos
						node.endPos = stack.pos
					}
					fmt.Println("    reduced", "node", fmt.Sprintf("%#v", node))

					// Create new stack after reduction
					newStateStack := append([]int{}, stack.stateStack[:len(stack.stateStack)-rhsLen]...)
					newNodeStack := append([]*ParseNode{}, stack.nodeStack[:len(stack.nodeStack)-rhsLen]...)

					// Get goto state
					gotoState := p.states[newStateStack[len(newStateStack)-1]].Gotos[rule.Nonterminal]
					newStateStack = append(newStateStack, gotoState)
					newNodeStack = append(newNodeStack, node)

					// Update current stack
					stack.stateStack = newStateStack
					stack.nodeStack = newNodeStack
					state = gotoState
					madeReduction = true
				}
				if !madeReduction {
					break
				}
			}

			state = stack.stateStack[len(stack.stateStack)-1]
			var actions []StateAction
			if isEOF {
				actions = p.states[state].Actions["$end"]
			} else {
				actions = p.states[state].Actions[symbol]
			}
			fmt.Println("  other actions with stack", "i", i, "state", state, "actions", actions)

			// Handle no actions - try error recovery
			if len(actions) == 0 {
				// Skip erroneous token and continue with same stack
				if !isEOF {
					newStacks = append(newStacks, &ParseStack{
						stateStack: append([]int{}, stack.stateStack...),
						nodeStack:  append([]*ParseNode{}, stack.nodeStack...),
						pos:        stack.pos,
					})
				}
				continue
			}

			// Process all possible actions
			for _, action := range actions {
				switch action.Action {
				case Shift:
					if isEOF {
						continue
					}

					// Create shifted token node
					node := &ParseNode{
						symbol:   symbol,
						value:    lval.token,
						startPos: stack.pos,
						endPos:   stack.pos + 1,
					}
					fmt.Println("    shifted", "node", fmt.Sprintf("%#v", node))

					// Create new stack after shift
					newStack := &ParseStack{
						stateStack: append(append([]int{}, stack.stateStack...), action.Target),
						nodeStack:  append(append([]*ParseNode{}, stack.nodeStack...), node),
						pos:        stack.pos + 1,
					}
					newStacks = append(newStacks, newStack)

				case Accept:
					fmt.Println("    accepting", "len(stack.nodeStack)", len(stack.nodeStack))
					if len(stack.nodeStack) > 0 {
						successfulParses = append(successfulParses, stack.nodeStack[len(stack.nodeStack)-1])
					}
				}
			}
		}

		if isEOF {
			break
		}

		// Merge similar stacks to control growth
		if len(newStacks) > 0 {
			stacks = mergeStacks(newStacks)
		}
	}

	for i, stack := range stacks {
		fmt.Println("stack", "i", i, "pos", stack.pos)
		for j, state := range stack.stateStack {
			if j < len(stack.nodeStack) {
				fmt.Println("  j", j, "state", state, "node", fmt.Sprintf("%#v", stack.nodeStack[j]))
			} else {
				fmt.Println("  j", j, "state", state)
			}
		}
	}

	// Return most complete successful parse
	if len(successfulParses) > 0 {
		var bestParse *ParseNode
		maxEnd := -1
		for _, parse := range successfulParses {
			if parse.endPos > maxEnd {
				maxEnd = parse.endPos
				bestParse = parse
			}
		}
		fmt.Println("bestParse", fmt.Sprintf("%#v", bestParse))
		return []*ParseNode{bestParse.children[0]}, nil
	}

	return nil, fmt.Errorf("failed to parse input")
}
