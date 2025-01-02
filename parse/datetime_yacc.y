%{
package parse

import (
"time"
"cloud.google.com/go/civil"
)

func setResult(l yyLexer, root *DateTimeTZRanges) {
  l.(*datetimeLexer).root = root
}

var ambiguousDateMode string

func NewCivilDate(first, second, year int) civil.Date {
  if ambiguousDateMode == "us" {
      return civil.Date{Month: time.Month(first), Day: second, Year: year}
}
return civil.Date{Day: first, Month: time.Month(second), Year: year}
}

%}

%token ILLEGAL

%token AM
%token AMP
%token CALENDAR
%token COLON
%token GOOGLE
%token ICS
%token PM
%token QUO
%token SEMICOLON
%token SUB
%token THROUGH
%token TO

%token <string> IDENT
%token <int> MONTH
%token <int> INT


 /* Type of each nonterminal. */
%type <DateTimeTZRanges> root
%type <DateTimeTZRanges> DateTimeTZRanges
%type <DateTimeTZRange> DateTimeTZRange
%type <DateTimeTZ> DateTimeTZ
%type <Date> Date
%type <int> Day
%type <int> Year

%start root


 /* Type of each nonterminal. */
%union{
    op    string
    label string
    string string
    int int
    Date  civil.Date
    DateTimeTZ  *DateTimeTZ
    DateTimeTZRange  *DateTimeTZRange
    DateTimeTZRanges  *DateTimeTZRanges
    }

%%

root:
  DateTimeTZRanges { setResult(yylex, $1) }
;

DateTimeTZRanges:
  DateTimeTZRange { $$ = &DateTimeTZRanges{Items: []*DateTimeTZRange{$1}}}
;

DateTimeTZRange:
  DateTimeTZ { $$ = &DateTimeTZRange{Start: $1}}
;

DateTimeTZ:
  Date { $$ = &DateTimeTZ{DateTime: civil.DateTime{Date: $1}}}
;

Date:
// "Feb 3 2023"
  MonthName Day Year { $$ = civil.Date{Month: 1, Day: $2, Year: $3}}
;

MonthName:
  MONTH
;

Day:
  INT
;

Year:
  INT
;
