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
		name       string
		input      string
		wantSymbol string
		wantErr    bool
	}{
		{
			name:       "Simple ABC",
			input:      "a b c",
			wantSymbol: "ABC",
			wantErr:    false,
		},
		{
			name:       "Simple BCD",
			input:      "b c d",
			wantSymbol: "BCD",
			wantErr:    false,
		},
		{
			name:       "Short BCD",
			input:      "b c",
			wantSymbol: "BCD",
			wantErr:    false,
		},
		{
			name:       "Simple ABCD",
			input:      "a b c d",
			wantSymbol: "ABCD",
			wantErr:    false,
		},
		{
			name:       "ABC with extra A",
			input:      "a a b c",
			wantSymbol: "ABC",
			wantErr:    false,
		},
		{
			name:       "ABC with noise",
			input:      "a b x c",
			wantSymbol: "ABC",
			wantErr:    false,
		},
		{
			name:       "Long BCD with noise",
			input:      "x b y c d",
			wantSymbol: "BCD",
			wantErr:    false,
		},
		{
			name:       "Long BCD with noise after",
			input:      "x b y c d y",
			wantSymbol: "BCD",
			wantErr:    false,
		},
		{
			name:       "Short BCD with noise",
			input:      "x b y c",
			wantSymbol: "BCD",
			wantErr:    false,
		},
		{
			name:       "Short BCD with noise after",
			input:      "x b y c y",
			wantSymbol: "BCD",
			wantErr:    false,
		},
		{
			name:       "ABCD with noise",
			input:      "x a y b c d x",
			wantSymbol: "ABCD",
			wantErr:    false,
		},
		{
			name:    "Invalid input",
			input:   "x y",
			wantErr: false,
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
			if len(results) == 0 && tt.wantSymbol == "" {
				return
			}

			// Get the root node (last node in result)
			root := results[0].Children[0]
			if root.Symbol != tt.wantSymbol {
				t.Errorf("Parse() got rule = %v, want %v", root.Symbol, tt.wantSymbol)
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
	if len(node.Children) > 0 {
		firstChild := node.Children[0]
		lastChild := node.Children[len(node.Children)-1]

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
	for _, child := range node.Children {
		verifyParseTreeHelper(t, child, visited)
	}
}
