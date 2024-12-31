package glr

import (
	"errors"

	"go/scanner"
	"go/token"
)

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
	tok := yyTokname(val)
	return l.lval.token, tok, val >= 0
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
			lval.token = "$end"
			return -1
		}

		// Handle identifiers and special tokens
		if tok == token.IDENT {
			lval.token = lit
			switch lit {
			case "A":
				return A
			case "B":
				return B
			case "C":
				return C
			case "D":
				return D
			case "X":
				return X
			case "Y":
				return Y
			}
		}

		// For any other token, treat as illegal but keep parsing
		lval.token = lit
		return ILLEGAL
	}
}
