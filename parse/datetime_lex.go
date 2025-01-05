package parse

import (
	"errors"
	"fmt"
	"strings"

	"go/scanner"
	"go/token"

	"github.com/findyourpaths/phil/glr"
)

// datetimeLexer implements yyLexer interface for the parser generated by goyacc,
// and also GLR's Lexer interface.
type datetimeLexer struct {
	lval    *yySymType
	scanner scanner.Scanner
	err     error
	root    *DateTimeTZRanges
}

// // // ast is the abstract syntax tree for a label formula
// type ast struct {
// 	root *atetime_ranges
// }

// // node represents a node in abstract syntax tree
// type datetime_ranges struct {
// 	items []*datetime_range
// }

// // node represents a node in abstract syntax tree
// type datetime_range struct {
// 	start *civil.DateTime
// 	end   *civil.DateTime
// }

func NewDatetimeLexer(input string) *datetimeLexer {
	l := &datetimeLexer{
		lval: &yySymType{},
	}
	l.scanner = glr.NewLexerScanner(l, input)
	return l
}

func (l *datetimeLexer) NextToken(pos int) (string, string, bool) {
	val := l.Lex(l.lval)
	sym := l.lval.string
	if val >= 57343 {
		sym = yyToknames[val-57343]
	}
	// fmt.Printf("sym: %q, l.lval.string: %q, val: %d\n", sym, l.lval.string, val)
	return sym, l.lval.string, val >= 0
}

func (l *datetimeLexer) Error(msg string) {
	l.err = errors.New(msg)
}

// yySymType is generated by goyacc
func (l *datetimeLexer) Lex(lval *yySymType) int {
	for {
		_, tok, lit := l.scanner.Scan()
		lval.string = lit

		if yyDebug == 3 {
			fmt.Printf("tok: %q, lit: %q\n", tok, lit)
		}

		// Skip whitespace and semicolons
		if tok == token.SEMICOLON {
			continue
		}
		if tok == token.EOF {
			lval.string = "$end"
			return -1
		}
		// if tok == token.EOF || (tok == token.SEMICOLON && lit == "\n") {
		// 	lval.string = "$end"
		// 	return -1
		// }

		switch tok {

		case token.ILLEGAL:
			litl := strings.ToLower(lit)
			switch litl {
			case "@":
				return AMP

			default:
				return ILLEGAL
			}

		case token.IDENT:
			lowLit := strings.ToLower(lit)
			if _, found := monthsByNames[lowLit]; found {
				return MONTH_NAME
			}

			if _, found := weekdaysByNames[lowLit]; found {
				return WEEKDAY_NAME
			}

			if ordinals[lowLit] {
				// Skip this by re-running.
				return l.Lex(lval)
			}

			switch lowLit {
			case "am":
				return AM
			case "and":
				return AND
			case "at":
				return AT
			case "calendar":
				return CALENDAR
			case "google":
				return GOOGLE
			case "ics":
				return ICS
			case "pm":
				return PM
			case "through":
				return THROUGH
			case "t":
				return T
			case "to":
				return TO
			case "when":
				return WHEN

			default:
				return IDENT
			}

		case token.INT:
			if len(lit) == 4 &&
				(strings.HasPrefix(lit, "17") ||
					strings.HasPrefix(lit, "18") ||
					strings.HasPrefix(lit, "19") ||
					strings.HasPrefix(lit, "20") ||
					strings.HasPrefix(lit, "21")) {
				return YEAR
			}
			return INT

		case token.COLON:
			return COLON
		case token.COMMA:
			return COMMA
		case token.PERIOD:
			return PERIOD
		case token.QUO:
			return QUO
		case token.SEMICOLON:
			return SEMICOLON
		case token.SUB:
			return SUB

		default:
			return ILLEGAL
		}
	}
}

// 	// fmt.Printf("lexeme before: %q", lexeme)
//   // lexeme := strings.ToLower(l.s.TokenText())
// 	// fmt.Printf("lexeme before: %q", lexeme)

// 	// if isInt(lexeme) {
// 	// 	return NUMBER
// 	// }

//   switch lexeme {
// 	case ` `:
// 		return SPC
// 	case `january`, `february`, `march`, `april`, `may`, `june`, `july`, `august`, `september`, `october`, `november`, `december`,
// 		`jan`, `feb`, `mar`, `apr`, `jun`, `jul`, `aug`, `sep`, `oct`, `nov`, `dec`:
// 		return MONTH
// 	case "(":
//     return OPEN // generated by goyacc
//   case ")":
//     return CLOSE // generated by goyacc
//   case ",":
//     lval.op = andOp
//     return OP // generated by goyacc
//   case ";":
//     lval.op = orOp
//     return OP // generated by goyacc
//   default:
//     lval.label = lexeme
//     return LABEL // generated by goyacc
//   }
// }

// func isInt(s string) bool {
//     for _, c := range s {
//         if !unicode.IsDigit(c) {
//             return false
//         }
//     }
//     return true
// }
