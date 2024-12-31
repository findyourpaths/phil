%{
package parse

import (
"time"
"cloud.google.com/go/civil"
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

%}

/* %token OPEN */
/* %token CLOSE */
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
 /* %token <op> OP */
 /* %token <label> LABEL */


 /* Type of each nonterminal. */
%type <datetime_ranges> root
%type <datetime_ranges> datetime_ranges
%type <datetime_range> datetime_range
%type <datetime> datetime
%type <date> date
%type <time> time

 /* %left OP */

%start root


 /* Type of each nonterminal. */
%union{
    op    string
    label string
    string string
    int int
    date  *civil.Date
    datetime  *civil.DateTime
    datetime_range  *datetime_range
    datetime_ranges  *datetime_ranges
    time  *civil.Time
    }

%%
/* formula: */ /* OPEN expr CLOSE { */ /* setResult(yylex, $2) */ /* } */

root:
  datetime_ranges datetime_ranges_suffixes { setResult(yylex, $1) }
;

datetime_ranges_suffixes:

| GOOGLE CALENDAR ICS
;

datetime_ranges:
  datetime_range { $$ = &datetime_ranges{ items: []*datetime_range{$1}}}
// "Feb 3-6; Mar 4-7"
| datetime_range SEMICOLON datetime_range { $$ = &datetime_ranges{ items: []*datetime_range{$1, $3}}}
;

datetime_range:
  datetime { $$ = &datetime_range{ start: $1}}
// "February 3-6, 2023"
| MONTH start_day separator end_day year { $$ = nil }
// "February 3 - March 3, 2023"
| MONTH INT separator datetime { $$ = nil }
// "Feb 3 Mar 4"
| datetime datetime { $$ = nil }
// "Feb 3, 2023 - Feb 26, 2023"
| datetime separator datetime { $$ = nil }
// "3 - 6 February 2023"
| INT separator datetime { $$ = nil }
// "February 3-6"
| MONTH INT separator INT { $$ = nil }
// "Feb 3 2023 9am - 12pm"
| date time separator time { $$ = nil }
// "Feb 3 2023 9am - Feb 3 2023 12pm"
| date time separator date time { $$ = nil }
// "Feb 3 2023 9:00 AM 3:00 PM"
| date time time { $$ = nil }
// "Feb 3 2023 9:00 AM 09:00 Feb 3 2023 3:00 PM 15:00"
| date time time date time time { $$ = nil }
// "Feb 3 2023 9:00 AM 09:00 3:00 PM 15:00"
| date time time time time { $$ = nil }
;

separator:
  SUB
| THROUGH
| TO
;

/* // "Feb 3 Mar 4" */ /* | MONTH INT MONTH INT { */ /* $$ = &datetime_range{ */ /* start: &datetime{month: $1, day: $2}, */ /* end: &datetime{month: $3, day: $4}}} */

datetime:
/* LABEL { */ /* $$ = &node{typ: label, val: $1} */ /* } */ /* | OPEN expr CLOSE { */ /* $$ = $2 */ /* } */ /* | expr OP expr { */ /* $$ = &node{typ: op, val: $2, left: $1, right: $3} */ /* } */
  date { $$ = nil }
// "Feb 3 2023 12 pm PST"
| date time { $$ = nil }
;

date:
// "Feb 3 2023"
  month_name day year  { $$ = nil }
// "2/3/2023", but ambiguous between US (month-first) and European/British (day-first) style.
| day QUO month QUO year { $$ = nil }
// "Feb 3"
| month_name day { $$ = nil }
// "3 Feb 2023"
| day month_name year { $$ = nil }
;

time:
  time_prefix time { $$ = nil }
// "Feb 3 2023 11am PST"
| time IDENT { $$ = nil }
// "11am"
| INT AM { $$ = nil }
// "15:00"
| INT COLON INT { $$ = nil }
// "12:00 AM"
| INT COLON INT AM { $$ = nil }
// "11am"
| INT PM { $$ = nil }
// "12pm"
| INT COLON INT PM { $$ = nil }
;

time_prefix:
  AMP
;


month_name:
  MONTH
;

/* start_month_name: */
/*   MONTH */
/* ; */

month:
  INT
;

day:
  INT
;

start_day:
  INT
;

end_day:
  INT
;

year:
  INT
;
