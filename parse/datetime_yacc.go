// Code generated by goyacc -l -o datetime_yacc.go -v datetime_yacc.states.txt datetime_yacc.y. DO NOT EDIT.
package parse

import __yyfmt__ "fmt"

import (
	"cloud.google.com/go/civil"
	"time"
)

func setResult(l yyLexer, root *datetime_ranges) {
	l.(*datetimeLexer).ast = &ast{
		root: root,
	}
}

var ambiguousDateMode string

func constructDate(first, second, year int) civil.Date {
	if ambiguousDateMode == "us" {
		return civil.Date{Month: time.Month(first), Day: second, Year: year}
	}
	return civil.Date{Day: first, Month: time.Month(second), Year: year}
}

type yySymType struct {
	yys             int
	op              string
	label           string
	string          string
	int             int
	date            *civil.Date
	datetime        *civil.DateTime
	datetime_range  *datetime_range
	datetime_ranges *datetime_ranges
	time            *civil.Time
}

const ILLEGAL = 57346
const AM = 57347
const AMP = 57348
const CALENDAR = 57349
const COLON = 57350
const GOOGLE = 57351
const ICS = 57352
const PM = 57353
const QUO = 57354
const SEMICOLON = 57355
const SUB = 57356
const THROUGH = 57357
const TO = 57358
const IDENT = 57359
const MONTH = 57360
const INT = 57361

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"ILLEGAL",
	"AM",
	"AMP",
	"CALENDAR",
	"COLON",
	"GOOGLE",
	"ICS",
	"PM",
	"QUO",
	"SEMICOLON",
	"SUB",
	"THROUGH",
	"TO",
	"IDENT",
	"MONTH",
	"INT",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

var yyExca = [...]int8{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 54,
	12, 37,
	18, 37,
	-2, 12,
}

const yyPrivate = 57344

const yyLast = 93

var yyAct = [...]int8{
	24, 45, 15, 7, 27, 46, 4, 27, 19, 54,
	27, 13, 14, 19, 20, 7, 34, 26, 40, 23,
	26, 33, 19, 57, 60, 39, 41, 5, 6, 52,
	37, 48, 49, 20, 35, 36, 27, 38, 22, 55,
	59, 56, 58, 53, 16, 17, 18, 40, 40, 26,
	16, 17, 18, 62, 19, 20, 27, 63, 29, 64,
	65, 12, 3, 68, 19, 69, 61, 40, 19, 57,
	16, 17, 18, 42, 66, 32, 43, 50, 8, 44,
	67, 11, 9, 31, 25, 47, 51, 21, 30, 10,
	2, 28, 1,
}

var yyPact = [...]int16{
	9, -1000, 72, 48, 36, 19, 56, -2, 14, 46,
	-1000, 76, 9, -1000, -5, -2, -1000, -1000, -1000, -1000,
	-1000, 56, 56, -5, 30, -2, 68, -1000, -14, 12,
	-14, 67, -1000, -1000, 31, 10, -10, -1000, 4, 50,
	-1000, 31, -1000, 5, -1000, -1000, -1000, 54, -1000, -1000,
	-1000, -14, -1000, -1000, -1000, 31, -2, 68, -2, 1,
	69, -14, -1000, 31, 1, 31, -1000, -1000, -1000, 31,
}

var yyPgo = [...]int8{
	0, 92, 90, 62, 6, 2, 0, 89, 87, 12,
	86, 1, 78, 82, 85, 84,
}

var yyR1 = [...]int8{
	0, 1, 7, 7, 2, 2, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 9, 9,
	9, 4, 4, 5, 5, 5, 5, 6, 6, 6,
	6, 6, 6, 6, 15, 12, 14, 13, 8, 10,
	11,
}

var yyR2 = [...]int8{
	0, 2, 0, 3, 1, 3, 1, 5, 4, 2,
	3, 3, 4, 4, 5, 3, 6, 5, 1, 1,
	1, 1, 2, 3, 5, 2, 3, 2, 2, 2,
	3, 4, 2, 4, 1, 1, 1, 1, 1, 1,
	1,
}

var yyChk = [...]int16{
	-1000, -1, -2, -3, -4, 18, 19, -5, -12, -13,
	-7, 9, 13, -4, -9, -5, 14, 15, 16, 18,
	19, -8, 19, -9, -6, -15, 19, 6, -13, 12,
	-12, 7, -3, -4, -6, -9, -9, -4, -9, -6,
	17, -6, 5, 8, 11, -11, 19, -14, 19, -11,
	10, -10, 19, -4, 19, -6, -5, 19, -5, -6,
	19, 12, -11, -6, -6, -6, 5, 11, -11, -6,
}

var yyDef = [...]int8{
	0, -2, 2, 4, 6, 0, 37, 21, 0, 0,
	1, 0, 0, 9, 0, 21, 18, 19, 20, 35,
	37, 0, 0, 0, 22, 0, 0, 34, 25, 0,
	0, 0, 5, 10, 22, 0, 0, 11, 0, 15,
	28, 27, 29, 0, 32, 23, 40, 0, 36, 26,
	3, 0, 39, 8, -2, 13, 0, 37, 0, 0,
	30, 0, 7, 14, 0, 17, 31, 33, 24, 16,
}

var yyTok1 = [...]int8{
	1,
}

var yyTok2 = [...]int8{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19,
}

var yyTok3 = [...]int8{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := int(yyPact[state])
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && int(yyChk[int(yyAct[n])]) == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || int(yyExca[i+1]) != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := int(yyExca[i])
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = int(yyTok1[0])
		goto out
	}
	if char < len(yyTok1) {
		token = int(yyTok1[char])
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = int(yyTok2[char-yyPrivate])
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = int(yyTok3[i+0])
		if token == char {
			token = int(yyTok3[i+1])
			goto out
		}
	}

out:
	if token == 0 {
		token = int(yyTok2[1]) /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = int(yyPact[yystate])
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = int(yyAct[yyn])
	if int(yyChk[yyn]) == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = int(yyDef[yystate])
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && int(yyExca[xi+1]) == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = int(yyExca[xi+0])
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = int(yyExca[xi+1])
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = int(yyPact[yyS[yyp].yys]) + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = int(yyAct[yyn]) /* simulate a shift of "error" */
					if int(yyChk[yystate]) == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= int(yyR2[yyn])
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = int(yyR1[yyn])
	yyg := int(yyPgo[yyn])
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = int(yyAct[yyg])
	} else {
		yystate = int(yyAct[yyj])
		if int(yyChk[yystate]) != -yyn {
			yystate = int(yyAct[yyg])
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			setResult(yylex, yyDollar[1].datetime_ranges)
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.datetime_ranges = &datetime_ranges{items: []*datetime_range{yyDollar[1].datetime_range}}
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.datetime_ranges = &datetime_ranges{items: []*datetime_range{yyDollar[1].datetime_range, yyDollar[3].datetime_range}}
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.datetime_range = &datetime_range{start: yyDollar[1].datetime}
		}
	case 7:
		yyDollar = yyS[yypt-5 : yypt+1]
		{
			yyVAL.datetime_range = nil
		}
	case 8:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.datetime_range = nil
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yyVAL.datetime_range = nil
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.datetime_range = nil
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.datetime_range = nil
		}
	case 12:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.datetime_range = nil
		}
	case 13:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.datetime_range = nil
		}
	case 14:
		yyDollar = yyS[yypt-5 : yypt+1]
		{
			yyVAL.datetime_range = nil
		}
	case 15:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.datetime_range = nil
		}
	case 16:
		yyDollar = yyS[yypt-6 : yypt+1]
		{
			yyVAL.datetime_range = nil
		}
	case 17:
		yyDollar = yyS[yypt-5 : yypt+1]
		{
			yyVAL.datetime_range = nil
		}
	case 21:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			yyVAL.datetime = nil
		}
	case 22:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yyVAL.datetime = nil
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.date = nil
		}
	case 24:
		yyDollar = yyS[yypt-5 : yypt+1]
		{
			yyVAL.date = nil
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yyVAL.date = nil
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.date = nil
		}
	case 27:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yyVAL.time = nil
		}
	case 28:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yyVAL.time = nil
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yyVAL.time = nil
		}
	case 30:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			yyVAL.time = nil
		}
	case 31:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.time = nil
		}
	case 32:
		yyDollar = yyS[yypt-2 : yypt+1]
		{
			yyVAL.time = nil
		}
	case 33:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			yyVAL.time = nil
		}
	}
	goto yystack /* stack new state and value */
}
