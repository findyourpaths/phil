%{
package parse

import "cloud.google.com/go/civil"
%}

%token ILLEGAL

%token AM
%token AMP
%token AND
%token AT
%token CALENDAR
%token COLON
%token COMMA
%token GOOGLE
%token ICS
%token PM
%token PERIOD
%token QUO
%token SEMICOLON
%token SUB
%token THROUGH
%token TO
%token T
%token WHEN

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
%type <Time> Time

%start root


 /* Type of each nonterminal. */
%union {
    DateTimeTZRanges  *DateTimeTZRanges
    DateTimeTZRange  *DateTimeTZRange
    DateTimeTZ  *DateTimeTZ
    Date  civil.Date
    Time  civil.Time
    string string
    }

%%

root:
  PrefixOpt DateTimeTZRanges SuffixOpt {$$ = $2}
;

DateTimeTZRanges:
  DateTimeTZRange {$$ = &DateTimeTZRanges{Items: []*DateTimeTZRange{$1}}}
| DateTimeTZRanges AndOpt DateTimeTZRange {$$ = AppendDateTimeTZRanges($1, $3)}

  // "Feb 3, 4"
| MONTH_NAME INT INT {$$ = NewRangesWithStartDates(NewMDYDate($1, $2, ""), NewMDYDate($1, $3, ""))}
| MONTH_NAME INT INT INT {$$ = NewRangesWithStartDates(NewMDYDate($1, $2, ""), NewMDYDate($1, $3, ""), NewMDYDate($1, $4, ""))}
| MONTH_NAME INT INT INT INT {$$ = NewRangesWithStartDates(NewMDYDate($1, $2, ""), NewMDYDate($1, $3, ""), NewMDYDate($1, $4, ""), NewMDYDate($1, $5, ""))}
| MONTH_NAME INT INT INT INT INT {$$ = NewRangesWithStartDates(NewMDYDate($1, $2, ""), NewMDYDate($1, $3, ""), NewMDYDate($1, $4, ""), NewMDYDate($1, $5, ""), NewMDYDate($1, $6, ""))}
  // "3, 4 Feb"
| INT INT MONTH_NAME {$$ = NewRangesWithStartDates(NewDMYDate($1, $3, ""), NewDMYDate($2, $3, ""))}
| INT INT INT MONTH_NAME {$$ = NewRangesWithStartDates(NewDMYDate($1, $4, ""), NewDMYDate($2, $4, ""), NewDMYDate($3, $4, ""))}
| INT INT INT INT MONTH_NAME {$$ = NewRangesWithStartDates(NewDMYDate($1, $5, ""), NewDMYDate($2, $5, ""), NewDMYDate($3, $5, ""), NewDMYDate($4, $5, ""))}
| INT INT INT INT INT MONTH_NAME {$$ = NewRangesWithStartDates(NewDMYDate($1, $6, ""), NewDMYDate($2, $6, ""), NewDMYDate($3, $6, ""), NewDMYDate($4, $6, ""), NewDMYDate($5, $6, ""))}
/*   // "3, 4 Feb" */
/* | INT Date {$$ = NewRangesWithStartDates(NewDMYDate($1, $2.Month, ""), $2)} */
/* | INT INT Date {$$ = NewRangesWithStartDates(NewDMYDate($1, $3.Month, ""), NewDMYDate($2, $3.Month, ""), $3)} */
/* | INT INT INT Date {$$ = NewRangesWithStartDates(NewDMYDate($1, $4.Month, ""), NewDMYDate($2, $4.Month, ""), NewDMYDate($3, $4.Month, ""), $4)} */
/* | INT INT INT INT Date {$$ = NewRangesWithStartDates(NewDMYDate($1, $5.Month, ""), NewDMYDate($2, $5.Month, ""), NewDMYDate($3, $5.Month, ""), NewDMYDate($4, $5.Month, ""), $5)} */

  // "Feb 3, 4 2023"
| MONTH_NAME INT INT YEAR {$$ = NewRangesWithStartDates(NewMDYDate($1, $2, $4), NewMDYDate($1, $3, $4))}
| MONTH_NAME INT INT INT YEAR {$$ = NewRangesWithStartDates(NewMDYDate($1, $2, $5), NewMDYDate($1, $3, $5), NewMDYDate($1, $4, $5))}
| MONTH_NAME INT INT INT INT YEAR {$$ = NewRangesWithStartDates(NewMDYDate($1, $2, $6), NewMDYDate($1, $3, $6), NewMDYDate($1, $4, $6), NewMDYDate($1, $5, $6))}
| MONTH_NAME INT INT INT INT INT YEAR {$$ = NewRangesWithStartDates(NewMDYDate($1, $2, $7), NewMDYDate($1, $3, $7), NewMDYDate($1, $4, $7), NewMDYDate($1, $5, $7), NewMDYDate($1, $6, $7))}
  // "3, 4 Feb 2023"
| INT INT MONTH_NAME YEAR {$$ = NewRangesWithStartDates(NewDMYDate($1, $3, $4), NewDMYDate($2, $3, $4))}
| INT INT INT MONTH_NAME YEAR {$$ = NewRangesWithStartDates(NewDMYDate($1, $4, $5), NewDMYDate($2, $4, $5), NewDMYDate($3, $4, $5))}
| INT INT INT INT MONTH_NAME YEAR {$$ = NewRangesWithStartDates(NewDMYDate($1, $5, $6), NewDMYDate($2, $5, $6), NewDMYDate($3, $5, $6), NewDMYDate($4, $5, $6))}
| INT INT INT INT INT MONTH_NAME YEAR {$$ = NewRangesWithStartDates(NewDMYDate($1, $6, $7), NewDMYDate($2, $6, $7), NewDMYDate($3, $6, $7), NewDMYDate($4, $6, $7), NewDMYDate($5, $6, $7))}
/*   // "3, 4 Feb 2023" */
/* | INT Date {$$ = NewRangesWithStartDates(NewDMYDate($1, $2.Month, $2.Year), $2)} */
/* | INT INT Date {$$ = NewRangesWithStartDates(NewDMYDate($1, $3.Month, $3.Year), NewDMYDate($2, $3.Month, $3.Year), $3)} */
/* | INT INT INT Date {$$ = NewRangesWithStartDates(NewDMYDate($1, $4.Month, $4.Year), NewDMYDate($2, $4.Month, $4.Year), NewDMYDate($3, $4.Month, $4.Year), $4)} */
/* | INT INT INT INT Date {$$ = NewRangesWithStartDates(NewDMYDate($1, $5.Month, $5.Year), NewDMYDate($2, $5.Month, $5.Year), NewDMYDate($3, $5.Month, $5.Year), NewDMYDate($4, $5.Month, $5.Year), $5)} */

  // "Feb 1-2, 3-4"
| MONTH_NAME INT RangeSep INT INT RangeSep INT {$$ = NewRanges(NewRangeWithStartEndDates(NewMDYDate($1, $2, ""), NewMDYDate($1, $4, "")), NewRangeWithStartEndDates(NewMDYDate($1, $5, ""), NewMDYDate($1, $7, "")))}
  // "1-2, 3-4 Feb"
| INT RangeSep INT INT RangeSep INT MONTH_NAME {$$ = NewRanges(NewRangeWithStartEndDates(NewDMYDate($1, $7, ""), NewDMYDate($3, $7, "")), NewRangeWithStartEndDates(NewDMYDate($4, $7, ""), NewDMYDate($6, $7, "")))}

  // "Feb 1-2, Mar 3-4"
| MONTH_NAME INT RangeSep INT MONTH_NAME INT RangeSep INT {$$ = NewRanges(NewRangeWithStartEndDates(NewMDYDate($1, $2, ""), NewMDYDate($1, $4, "")), NewRangeWithStartEndDates(NewMDYDate($5, $6, ""), NewMDYDate($5, $8, "")))}
  // "1-2 Feb, 3-4 Mar"
| INT RangeSep INT MONTH_NAME INT RangeSep INT MONTH_NAME {$$ = NewRanges(NewRangeWithStartEndDates(NewDMYDate($1, $4, ""), NewDMYDate($3, $4, "")), NewRangeWithStartEndDates(NewDMYDate($5, $8, ""), NewDMYDate($7, $8, "")))}

  // "Feb 1-2, 3-4 2023"
| MONTH_NAME INT RangeSep INT INT RangeSep INT YEAR {$$ = NewRanges(NewRangeWithStartEndDates(NewMDYDate($1, $2, $8), NewMDYDate($1, $4, $8)), NewRangeWithStartEndDates(NewMDYDate($1, $5, $8), NewMDYDate($1, $7, $8)))}
  // "1-2, 3-4 Feb 2023"
| INT RangeSep INT INT RangeSep INT MONTH_NAME YEAR {$$ = NewRanges(NewRangeWithStartEndDates(NewDMYDate($1, $7, $8), NewDMYDate($3, $7, $8)), NewRangeWithStartEndDates(NewDMYDate($4, $7, $8), NewDMYDate($6, $7, $8)))}

  // "Feb 1-2, Mar 3-4 2023"
| MONTH_NAME INT RangeSep INT MONTH_NAME INT RangeSep INT YEAR {$$ = NewRanges(NewRangeWithStartEndDates(NewMDYDate($1, $2, $9), NewMDYDate($1, $4, $9)), NewRangeWithStartEndDates(NewMDYDate($5, $6, $9), NewMDYDate($5, $8, $9)))}
  // "1-2 Feb, 3-4 Mar 2023"
| INT RangeSep INT MONTH_NAME INT RangeSep INT MONTH_NAME YEAR {$$ = NewRanges(NewRangeWithStartEndDates(NewDMYDate($1, $4, $9), NewDMYDate($3, $4, $9)), NewRangeWithStartEndDates(NewDMYDate($5, $8, $9), NewDMYDate($7, $8, $9)))}

  // "Feb 3, Mar 4 2023"
| MONTH_NAME INT MONTH_NAME INT YEAR {$$ = NewRanges(NewRangeWithStart(NewMDYDate($1, $2, $5)), NewRangeWithStart(NewMDYDate($3, $4, $5)))}
;


DateTimeTZRange:
  DateTimeTZ {$$ = &DateTimeTZRange{Start: $1}}
| DateTimeTZ RangeSep Time {$$ = NewRangeWithStartEndDateTimes($1, NewDateTime($1.Date, $3, ""))}

  // "Feb 3-4"
| MONTH_NAME INT RangeSep INT {$$ = NewRangeWithStartEndDates(NewMDYDate($1, $2, ""), NewMDYDate($1, $4, ""))}
  // "3-4 Feb"
| INT RangeSep INT MONTH_NAME {$$ = NewRangeWithStartEndDates(NewDMYDate($1, $4, ""), NewDMYDate($3, $4, ""))}

  // "Feb 3-4, 2023"
| MONTH_NAME INT RangeSep INT YEAR {$$ = NewRangeWithStartEndDates(NewMDYDate($1, $2, $5), NewMDYDate($1, $4, $5))}
  // "3-4 Feb 2023"
| INT RangeSep INT MONTH_NAME YEAR {$$ = NewRangeWithStartEndDates(NewDMYDate($1, $4, $5), NewDMYDate($3, $4, $5))}

  // "Feb 3, 2023 - Feb 4, 2023"
| DateTimeTZ RangeSep DateTimeTZ {$$ = &DateTimeTZRange{Start: $1, End: $3}}

  // "Feb 3 - Mar 4, 2023"
| MONTH_NAME INT RangeSep MONTH_NAME INT YEAR {$$ = NewRangeWithStartEndDates(NewMDYDate($1, $2, $6), NewMDYDate($4, $5, $6))}

  // "Feb 3 2023 9:00 AM 09:00"
  // "Feb 3 2023 3:00 PM 15:00"
/* | Date Time Time {$$ = NewRangeWithStartEndDateTimes(NewDateTime($1, $3, ""), NewDateTime($1, $3, ""))} */
;

DateTimeTZ:
  Date {$$ = NewDateTimeWithDate($1)}
| Date DateTimeSep Time {$$ = NewDateTime($1, $3, "")}
;

Date:
  WEEKDAY_NAME CommaOpt Date {$$ = $3}

  // "2006-01-02 T 15:04:05Z07:00"
| Date T

  // "02.03", but ambiguous between North America (month-day-year) and other (day-month-year) styles.
| INT DateSep INT {$$ = NewAmbiguousDate($1, $3, "")}

| YEAR {$$ = NewDMYDate("", "", $1)}
| YEAR SUB INT {$$ = NewDMYDate("", $3, $1)}
| YEAR SUB INT SUB INT {$$ = NewDMYDate($5, $3, $1)}

  // "Feb 3 2023"
| MONTH_NAME INT YEAR {$$ = NewMDYDate($1, $2, $3)}

  // "3 Feb 2023"
| INT MONTH_NAME YEAR {$$ = NewDMYDate($1, $2, $3)}

  // "2/3/2023", but ambiguous between North America (month-day-year) and other (day-month-year) styles.
| INT DateSep INT DateSep YEAR {$$ = NewAmbiguousDate($1, $3, $5)}

  // "Feb 3"
| MONTH_NAME INT {$$ = NewMDYDate($1, $2, "")}

  // "3 Feb"
| INT MONTH_NAME {$$ = NewDMYDate($1, $2, "")}
;


Time:
  INT {$$ = NewTime($1, "", "", "")}

  // "11am"
| INT AM {$$ = NewTime($1, "", "", "")}
| INT PM {$$ = NewTime((mustAtoi($1) % 12) + 12, "", "", "")}

  // "12:00"
| INT COLON INT {$$ = NewTime($1, $3, "", "")}

  // "12:00:00"
| INT COLON INT COLON INT {$$ = NewTime((mustAtoi($1) % 12) + 12, $3, $5, "")}

  // "9:00 AM"
| INT COLON INT AM {$$ = NewTime($1, $3, "", "")}

  // "12:00 PM"
| INT COLON INT PM {$$ = NewTime((mustAtoi($1) % 12) + 12, $3, "", "")}

/* // "Feb 3 2023 11am PST" */
/* |   time TimeZoneOpt { */
/*        $$ = $1} //civil.Time{Hour: $1}} //, ampm: $5, timezone: $6}} */
;


DateSep:
  QUO
| PERIOD
;


DateTimeSep:

| AT
| SUB
;


RangeSep:
  SUB
| THROUGH
| TO
;


PrefixOpt:
  WhenOpt
;


SuffixOpt:
  GoogleOpt CalendarOpt ICSOpt
;


/* All optional terminals */

AndOpt:

| AND
;


CalendarOpt:

| CALENDAR
;


CommaOpt:

| COMMA
;


GoogleOpt:

| GOOGLE
;


ICSOpt:

| ICS
;


WhenOpt:

| WHEN
;
