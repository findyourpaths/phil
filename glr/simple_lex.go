package glr

import (
	"errors"

	"go/scanner"
	"go/token"
)

type Alphabet struct {
	ABCD *ABCD
	ABC  *ABC
	BCD  *BCD
}

type ABCD struct {
	A string
	B string
	C string
	D string
}

type ABC struct {
	A string
	B string
	C string
}

type BCD struct {
	B string
	C string
	D string
}

// simpleLexer implements yyLexer interface for the parser generated by goyacc,
// and also GLR's Lexer interface.
type simpleLexer struct {
	lval    *yySymType
	scanner scanner.Scanner
	err     error
}

func NewSimpleLexer(input string) *simpleLexer {
	l := &simpleLexer{
		lval: &yySymType{},
	}
	l.scanner = NewLexerScanner(l, input)
	return l
}

func (l *simpleLexer) NextToken(pos int) (string, any, bool) {
	val := l.Lex(l.lval)
	sym := l.lval.string
	if val >= 57343 {
		sym = yyToknames[val-57343]
	}
	debugf("sym: %q, l.lval.string: %q, val: %d\n", sym, l.lval.string, val)
	return sym, l.lval.string, val >= 0
}

func (l *simpleLexer) Error(msg string) {
	l.err = errors.New(msg)
}

// yySymType is generated by goyacc
func (l *simpleLexer) Lex(lval *yySymType) int {
	for {
		_, tok, lit := l.scanner.Scan()

		// Skip whitespace and semicolons
		if tok == token.SEMICOLON {
			continue
		}
		if tok == token.EOF {
			lval.string = "$end"
			return -1
		}

		// Handle identifiers and special tokens
		if tok == token.IDENT {
			lval.string = lit
			switch lit {
			case "a":
				return A
			case "b":
				return B
			case "c":
				return C
			case "d":
				return D
			case "x":
				return X
			case "y":
				return Y
			}
		}

		// For any other token, treat as illegal but keep parsing
		lval.string = lit
		return ILLEGAL
	}
}
