// glr_test.go
package glr

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestGLRParser(t *testing.T) {
	if os.Getenv("DEBUG") == "true" {
		DoDebug = true
	}

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

	for i, tc := range tests {
		// for i, tc := range tests[11:12] {
		t.Run(fmt.Sprintf("%03d__%s", i, tc.name), func(t *testing.T) {
			g := &Grammar{Rules: glrRules, Actions: glrActions, States: glrStates}
			results, err := Parse(g, NewSimpleLexer(tc.input))
			if (err != nil) != tc.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err != nil {
				return
			}
			if len(results) == 0 && tc.want == nil {
				return
			}

			// Get the root node (last node in result)
			got, err := GetParseNodeValue(g, results[0], "")
			if err != nil {
				t.Errorf("error getting parse node value: %v", err)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Parse() got rule = %#v, want %#v", got, tc.want)
			}

			// Verify the parse tree structure
			verifyParseTree(t, results[0])
		})
	}
}

func verifyParseTree(t *testing.T, node *ParseNode) {
	visited := make(map[*ParseNode]bool)
	verifyParseTreeHelper(t, node, visited)
}

func verifyParseTreeHelper(t *testing.T, n *ParseNode, visited map[*ParseNode]bool) {
	if n == nil || visited[n] {
		return
	}
	visited[n] = true

	if n.isAlt {
		// Recursively verify children
		for _, child := range n.Children {
			verifyParseTreeHelper(t, child, visited)
		}
		return
	}

	// Verify node positions are consistent
	if len(n.Children) > 0 {
		first := n.Children[0]
		last := n.Children[len(n.Children)-1]

		if n.startPos != first.startPos {
			t.Errorf("Node start position inconsistent: node=%d, firstChild=%d",
				n.startPos, first.startPos)
		}

		if n.endPos != last.endPos {
			t.Errorf("Node end position inconsistent: node=%d, lastChild=%d",
				n.endPos, last.endPos)
		}
	}

	// Recursively verify children
	for _, child := range n.Children {
		verifyParseTreeHelper(t, child, visited)
	}
}
