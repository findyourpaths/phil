%{
package parse

import "cloud.google.com/go/civil"
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
%token <string> MONTH_NAME
%token <string> WEEKDAY_NAME
%token <string> YEAR
%token <string> INT


 /* Type of each nonterminal. */
%type <DateTimeTZRanges> root
%type <DateTimeTZRanges> DateTimeTZRanges
%type <DateTimeTZRange> DateTimeTZRange
%type <DateTimeTZ> DateTimeTZ
%type <Date> Date

%start root


 /* Type of each nonterminal. */
%union {
    DateTimeTZRanges  *DateTimeTZRanges
    DateTimeTZRange  *DateTimeTZRange
    DateTimeTZ  *DateTimeTZ
    Date  civil.Date
    string string
    }

%%

root:
  DateTimeTZRanges {$$ = $1}
;

DateTimeTZRanges:
  DateTimeTZRange {$$ = &DateTimeTZRanges{Items: []*DateTimeTZRange{$1}}}
| DateTimeTZRanges DateTimeTZRange {$$ = AppendDateTimeTZRanges($1, $2)}

  // "Feb 3, 4"
| MONTH_NAME INT INT {$$ = NewRangesWithStarts(NewMDYDate($1, $2, ""), NewMDYDate($1, $3, ""))}
| MONTH_NAME INT INT INT {$$ = NewRangesWithStarts(NewMDYDate($1, $2, ""), NewMDYDate($1, $3, ""), NewMDYDate($1, $4, ""))}
| MONTH_NAME INT INT INT INT {$$ = NewRangesWithStarts(NewMDYDate($1, $2, ""), NewMDYDate($1, $3, ""), NewMDYDate($1, $4, ""), NewMDYDate($1, $5, ""))}
| MONTH_NAME INT INT INT INT INT {$$ = NewRangesWithStarts(NewMDYDate($1, $2, ""), NewMDYDate($1, $3, ""), NewMDYDate($1, $4, ""), NewMDYDate($1, $5, ""), NewMDYDate($1, $6, ""))}
  // "3, 4 Feb"
| INT INT MONTH_NAME {$$ = NewRangesWithStarts(NewDMYDate($1, $3, ""), NewDMYDate($2, $3, ""))}
| INT INT INT MONTH_NAME {$$ = NewRangesWithStarts(NewDMYDate($1, $4, ""), NewDMYDate($2, $4, ""), NewDMYDate($3, $4, ""))}
| INT INT INT INT MONTH_NAME {$$ = NewRangesWithStarts(NewDMYDate($1, $5, ""), NewDMYDate($2, $5, ""), NewDMYDate($3, $5, ""), NewDMYDate($4, $5, ""))}
| INT INT INT INT INT MONTH_NAME {$$ = NewRangesWithStarts(NewDMYDate($1, $6, ""), NewDMYDate($2, $6, ""), NewDMYDate($3, $6, ""), NewDMYDate($4, $6, ""), NewDMYDate($5, $6, ""))}

  // "Feb 3, 4 2023"
| MONTH_NAME INT INT YEAR {$$ = NewRangesWithStarts(NewMDYDate($1, $2, $4), NewMDYDate($1, $3, $4))}
| MONTH_NAME INT INT INT YEAR {$$ = NewRangesWithStarts(NewMDYDate($1, $2, $5), NewMDYDate($1, $3, $5), NewMDYDate($1, $4, $5))}
| MONTH_NAME INT INT INT INT YEAR {$$ = NewRangesWithStarts(NewMDYDate($1, $2, $6), NewMDYDate($1, $3, $6), NewMDYDate($1, $4, $6), NewMDYDate($1, $5, $6))}
| MONTH_NAME INT INT INT INT INT YEAR {$$ = NewRangesWithStarts(NewMDYDate($1, $2, $7), NewMDYDate($1, $3, $7), NewMDYDate($1, $4, $7), NewMDYDate($1, $5, $7), NewMDYDate($1, $6, $7))}
  // "3, 4 Feb 2023"
| INT INT MONTH_NAME YEAR {$$ = NewRangesWithStarts(NewDMYDate($1, $3, $4), NewDMYDate($2, $3, $4))}
| INT INT INT MONTH_NAME YEAR {$$ = NewRangesWithStarts(NewDMYDate($1, $4, $5), NewDMYDate($2, $4, $5), NewDMYDate($3, $4, $5))}
| INT INT INT INT MONTH_NAME YEAR {$$ = NewRangesWithStarts(NewDMYDate($1, $5, $6), NewDMYDate($2, $5, $6), NewDMYDate($3, $5, $6), NewDMYDate($4, $5, $6))}
| INT INT INT INT INT MONTH_NAME YEAR {$$ = NewRangesWithStarts(NewDMYDate($1, $6, $7), NewDMYDate($2, $6, $7), NewDMYDate($3, $6, $7), NewDMYDate($4, $6, $7), NewDMYDate($5, $6, $7))}

  // "Feb 1-2, 3-4"
| MONTH_NAME INT Sep INT INT Sep INT {$$ = NewRanges(NewRangeWithStartEnd(NewMDYDate($1, $2, ""), NewMDYDate($1, $4, "")), NewRangeWithStartEnd(NewMDYDate($1, $5, ""), NewMDYDate($1, $7, "")))}
  // "1-2, 3-4 Feb"
| INT Sep INT INT Sep INT MONTH_NAME {$$ = NewRanges(NewRangeWithStartEnd(NewDMYDate($1, $7, ""), NewDMYDate($3, $7, "")), NewRangeWithStartEnd(NewDMYDate($4, $7, ""), NewDMYDate($6, $7, "")))}

  // "Feb 1-2, Mar 3-4"
| MONTH_NAME INT Sep INT MONTH_NAME INT Sep INT {$$ = NewRanges(NewRangeWithStartEnd(NewMDYDate($1, $2, ""), NewMDYDate($1, $4, "")), NewRangeWithStartEnd(NewMDYDate($5, $6, ""), NewMDYDate($5, $8, "")))}
  // "1-2 Feb, 3-4 Mar"
| INT Sep INT MONTH_NAME INT Sep INT MONTH_NAME {$$ = NewRanges(NewRangeWithStartEnd(NewDMYDate($1, $4, ""), NewDMYDate($3, $4, "")), NewRangeWithStartEnd(NewDMYDate($5, $8, ""), NewDMYDate($7, $8, "")))}

  // "Feb 1-2, 3-4 2023"
| MONTH_NAME INT Sep INT INT Sep INT YEAR {$$ = NewRanges(NewRangeWithStartEnd(NewMDYDate($1, $2, $8), NewMDYDate($1, $4, $8)), NewRangeWithStartEnd(NewMDYDate($1, $5, $8), NewMDYDate($1, $7, $8)))}
  // "1-2, 3-4 Feb 2023"
| INT Sep INT INT Sep INT MONTH_NAME YEAR {$$ = NewRanges(NewRangeWithStartEnd(NewDMYDate($1, $7, $8), NewDMYDate($3, $7, $8)), NewRangeWithStartEnd(NewDMYDate($4, $7, $8), NewDMYDate($6, $7, $8)))}

  // "Feb 1-2, Mar 3-4 2023"
| MONTH_NAME INT Sep INT MONTH_NAME INT Sep INT YEAR {$$ = NewRanges(NewRangeWithStartEnd(NewMDYDate($1, $2, $9), NewMDYDate($1, $4, $9)), NewRangeWithStartEnd(NewMDYDate($5, $6, $9), NewMDYDate($5, $8, $9)))}
  // "1-2 Feb, 3-4 Mar 2023"
| INT Sep INT MONTH_NAME INT Sep INT MONTH_NAME YEAR {$$ = NewRanges(NewRangeWithStartEnd(NewDMYDate($1, $4, $9), NewDMYDate($3, $4, $9)), NewRangeWithStartEnd(NewDMYDate($5, $8, $9), NewDMYDate($7, $8, $9)))}

  // "Feb 3, Mar 4 2023"
| MONTH_NAME INT MONTH_NAME INT YEAR {$$ = NewRanges(NewRangeWithStart(NewMDYDate($1, $2, $5)), NewRangeWithStart(NewMDYDate($3, $4, $5)))}
;


DateTimeTZRange:
  DateTimeTZ {$$ = &DateTimeTZRange{Start: $1}}

  // "Feb 3-4"
| MONTH_NAME INT Sep INT {$$ = NewRangeWithStartEnd(NewMDYDate($1, $2, ""), NewMDYDate($1, $4, ""))}
  // "3-4 Feb"
| INT Sep INT MONTH_NAME {$$ = NewRangeWithStartEnd(NewDMYDate($1, $4, ""), NewDMYDate($3, $4, ""))}

  // "Feb 3-4, 2023"
| MONTH_NAME INT Sep INT YEAR {$$ = NewRangeWithStartEnd(NewMDYDate($1, $2, $5), NewMDYDate($1, $4, $5))}
  // "3-4 Feb 2023"
| INT Sep INT MONTH_NAME YEAR {$$ = NewRangeWithStartEnd(NewDMYDate($1, $4, $5), NewDMYDate($3, $4, $5))}

  // "Feb 3, 2023 - Feb 4, 2023"
| DateTimeTZ Sep DateTimeTZ {$$ = &DateTimeTZRange{Start: $1, End: $3}}

  // "Feb 3 - Mar 4, 2023"
| MONTH_NAME INT Sep MONTH_NAME INT YEAR {$$ = NewRangeWithStartEnd(NewMDYDate($1, $2, $6), NewMDYDate($4, $5, $6))}
;

DateTimeTZ:
  Date {$$ = &DateTimeTZ{DateTime: civil.DateTime{Date: $1}}}
;


Date:
  // "Feb 3 2023"
  MONTH_NAME INT YEAR {$$ = NewMDYDate($1, $2, $3)}

  // "3 Feb 2023"
| INT MONTH_NAME YEAR {$$ = NewDMYDate($1, $2, $3)}

  // "2/3/2023", but ambiguous between North America (month-day-year) and other (day-month-year) styles.
| INT QUO INT QUO YEAR {$$ = NewAmbiguousDate(ambiguousDateMode, $1, $3, $5)}

  // "Feb 3"
  // "Thu Feb 3"
| WeekDayNameOpt MONTH_NAME INT {$$ = NewMDYDate($2, $3, "")}

  // "3 Feb"
| WeekDayNameOpt INT MONTH_NAME {$$ = NewDMYDate($2, $3, "")}
;


WeekDayNameOpt:

| WEEKDAY_NAME
;


Sep:
  SUB
| THROUGH
| TO
