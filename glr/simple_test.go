// glr_test.go
package glr

import (
	"reflect"
	"testing"
)

func TestGLRParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Alphabet
		wantErr bool
	}{
		{
			name:    "Simple ABC",
			input:   "a b c",
			want:    &Alphabet{ABC: &ABC{A: "a", B: "b", C: "c"}},
			wantErr: false,
		},
		{
			name:    "Simple BCD",
			input:   "b c d",
			want:    &Alphabet{BCD: &BCD{B: "b", C: "c", D: "d"}},
			wantErr: false,
		},
		{
			name:    "Short BCD",
			input:   "b c",
			want:    &Alphabet{BCD: &BCD{B: "b", C: "c"}},
			wantErr: false,
		},
		{
			name:    "Simple ABCD",
			input:   "a b c d",
			want:    &Alphabet{ABCD: &ABCD{A: "a", B: "b", C: "c", D: "d"}},
			wantErr: false,
		},
		{
			name:    "ABC with extra A",
			input:   "a a b c",
			want:    &Alphabet{ABC: &ABC{A: "a", B: "b", C: "c"}},
			wantErr: false,
		},
		{
			name:    "ABC with noise",
			input:   "a b x c",
			want:    &Alphabet{ABC: &ABC{A: "a", B: "b", C: "c"}},
			wantErr: false,
		},
		{
			name:    "Long BCD with noise",
			input:   "x b y c d",
			want:    &Alphabet{BCD: &BCD{B: "b", C: "c", D: "d"}},
			wantErr: false,
		},
		{
			name:    "Long BCD with noise after",
			input:   "x b y c d y",
			want:    &Alphabet{BCD: &BCD{B: "b", C: "c", D: "d"}},
			wantErr: false,
		},
		{
			name:    "Short BCD with noise",
			input:   "x b y c",
			want:    &Alphabet{BCD: &BCD{B: "b", C: "c"}},
			wantErr: false,
		},
		{
			name:    "Short BCD with noise after",
			input:   "x b y c y",
			want:    &Alphabet{BCD: &BCD{B: "b", C: "c"}},
			wantErr: false,
		},
		{
			name:    "ABCD with noise",
			input:   "x a y b c d x",
			want:    &Alphabet{ABCD: &ABCD{A: "a", B: "b", C: "c", D: "d"}},
			wantErr: false,
		},
		{
			name:    "Invalid input",
			input:   "x y",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := Parse(&Grammar{Rules: glrRules, Actions: glrActions, States: glrStates}, NewSimpleLexer(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if len(results) == 0 && tt.want == nil {
				return
			}

			// Get the root node (last node in result)
			root := results[0]
			if !reflect.DeepEqual(root.Value, tt.want) {
				t.Errorf("Parse() got rule = %#v, want %#v", root.Value, tt.want)
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
