package datetime

import (
	"errors"
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
		debugf("tok: %q, lit: %q, tok == token.ILLEGAL: %t, tok == token.IDENT: %t\n", tok, lit, tok == token.ILLEGAL, tok == token.IDENT)

		// Skip whitespace and semicolons
		if tok == token.SEMICOLON {
			continue
		}
		if tok == token.EOF {
			lval.string = "$end"
			return -1
		}

		lval.string = lit
		if lit == "" {
			lval.string = tok.String()
		}
		// fmt.Printf("lval.string: %#v\n", lval.string)

		switch tok {

		case token.ILLEGAL:
			litl := strings.ToLower(lit)
			switch litl {
			case "@":
				return AMP

			case "–":
				return SUB

			case "—":
				return SUB

			default:
				return ILLEGAL
			}

		case token.IDENT:
			upLit := strings.ToUpper(lit)
			// tz, err := tzTimezone.GetTzAbbreviationInfo(upLit)
			// fmt.Println("upLit", upLit, "tz", tz, "err", err)
			if tz, _ := timezoneTZ.GetTzAbbreviationInfo(upLit); tz != nil {
				return TIME_ZONE_ABBREV
			}
			if tz, _ := timezoneTZ.GetTzInfo(upLit); tz != nil {
				return TIME_ZONE
			}

			lowLit := strings.ToLower(lit)
			if _, found := monthsByNames[lowLit]; found {
				return MONTH_NAME
			}

			if _, found := weekdaysByNames[lowLit]; found {
				return WEEKDAY_NAME
			}

			if ordinals[lowLit] {
				return ORD_IND
			}

			switch lowLit {
			case "am":
				return AM
			case "and":
				return AND
			case "at":
				return AT
			case "beginning":
				return BEGINNING
			case "calendar":
				return CALENDAR
			case "date":
				return DATE
			case "from":
				return FROM
			case "google":
				return GOOGLE
			case "ics":
				return ICS
			case "in":
				return IN
			case "of":
				return OF
			case "on":
				return ON
			case "pm":
				return PM
			case "through":
				return THROUGH
			case "t":
				return T
			case "th":
				return TH
			case "till":
				return TILL
			case "time":
				return TIME
			case "to":
				return TO
			case "until":
				return UNTIL
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
		case token.DEC:
			return DEC
		case token.LPAREN:
			return LPAREN
		case token.PERIOD:
			return PERIOD
		case token.QUO:
			return QUO
		case token.RPAREN:
			return RPAREN
		case token.SEMICOLON:
			return SEMICOLON
		case token.SUB:
			return SUB
		default:
			return ILLEGAL
		}
	}
}