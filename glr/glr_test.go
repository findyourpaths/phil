// glr_test.go
package glr

import (
	"testing"
)

func TestGLRParser(t *testing.T) {
	rules, states, err := LoadGrammarRulesAndStates("simple_yacc.y", "simple_yacc.states.txt")
	if err != nil {
		t.Fatalf("Failed to load rules and states: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		wantRule string
		wantErr  bool
	}{
		{
			name:     "Simple ABC",
			input:    "A B C",
			wantRule: "ABC",
			wantErr:  false,
		},
		{
			name:     "Simple BCD",
			input:    "B C D",
			wantRule: "BCD",
			wantErr:  false,
		},
		{
			name:     "Simple ABCD",
			input:    "A B C D",
			wantRule: "ABCD",
			wantErr:  false,
		},
		{
			name:     "ABC with extra A",
			input:    "A A B C",
			wantRule: "ABC",
			wantErr:  false,
		},
		{
			name:     "ABC with noise",
			input:    "A B X C",
			wantRule: "ABC",
			wantErr:  false,
		},
		{
			name:     "BCD with noise",
			input:    "X B Y C D",
			wantRule: "BCD",
			wantErr:  false,
		},
		{
			name:     "ABCD with noise",
			input:    "X A Y B C D X",
			wantRule: "ABCD",
			wantErr:  false,
		},
		{
			name:  "Invalid input",
			input: "X Y",
			// wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(rules, states, NewSimpleLexer(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if len(results) == 0 && tt.wantRule == "" {
				return
			}

			// Get the root node (last node in result)
			root := results[0].children[0]
			if root.symbol != tt.wantRule {
				t.Errorf("Parse() got rule = %v, want %v", root.symbol, tt.wantRule)
			}

			// Verify the parse tree structure
			verifyParseTree(t, root)
		})
	}
}

func verifyParseTree(t *testing.T, node *ParseNode) {
	visited := make(map[*ParseNode]bool)
	verifyParseTreeHelper(t, node, visited)
}

func verifyParseTreeHelper(t *testing.T, node *ParseNode, visited map[*ParseNode]bool) {
	if node == nil || visited[node] {
		return
	}
	visited[node] = true

	// Verify node positions are consistent
	if len(node.children) > 0 {
		firstChild := node.children[0]
		lastChild := node.children[len(node.children)-1]

		if node.startPos != firstChild.startPos {
			t.Errorf("Node start position inconsistent: node=%d, firstChild=%d",
				node.startPos, firstChild.startPos)
		}

		if node.endPos != lastChild.endPos {
			t.Errorf("Node end position inconsistent: node=%d, lastChild=%d",
				node.endPos, lastChild.endPos)
		}
	}

	// Recursively verify children
	for _, child := range node.children {
		verifyParseTreeHelper(t, child, visited)
	}
}
