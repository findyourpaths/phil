%{
package parse

import "cloud.google.com/go/civil"
%}

%token ILLEGAL

%token AM
%token AMP
%token AND
%token AT
%token BEGINNING
%token CALENDAR
%token COLON
%token COMMA
%token DEC
%token FROM
%token GOOGLE
%token ICS
%token OF
%token ORD_IND
%token PM
%token PERIOD
%token QUO
%token SEMICOLON
%token SUB
%token THROUGH
%token T
%token TH
%token TILL
%token TO
%token UNTIL
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

%type <string> Day

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
  RangePrefix DateTimeTZRange {$$ = $2}
| RangePrefixOpt DateTimeTZ {$$ = &DateTimeTZRange{Start: $2}}
| RangePrefixOpt DateTimeTZ RangeSep Time {$$ = NewRangeWithStartEndDateTimes($2, NewDateTime($2.Date, $4, ""))}

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

  // "Thu Feb 3 - Sat Mar 4, 2023"
| WeekDay MONTH_NAME INT RangeSep WeekDay MONTH_NAME INT YEAR {$$ = NewRangeWithStartEndDates(NewMDYDate($2, $3, $8), NewMDYDate($6, $7, $8))}

  // "Thu 3 Feb - Sat 4 Mar, 2023"
| WeekDay INT MONTH_NAME RangeSep WeekDay INT MONTH_NAME YEAR {$$ = NewRangeWithStartEndDates(NewDMYDate($2, $3, $8), NewDMYDate($6, $7, $8))}

  // "Thu Feb 3 - Sat 4 Mar, 2023"
| WeekDay MONTH_NAME INT RangeSep WeekDay INT MONTH_NAME YEAR {$$ = NewRangeWithStartEndDates(NewMDYDate($2, $3, $8), NewDMYDate($6, $7, $8))}

  // "Thu 3 Feb - Sat Mar 4, 2023"
| WeekDay MONTH_NAME INT RangeSep WeekDay MONTH_NAME INT YEAR {$$ = NewRangeWithStartEndDates(NewMDYDate($2, $3, $8), NewMDYDate($6, $7, $8))}

  // "Feb 3 2023 9:00 AM 09:00"
  // "Feb 3 2023 3:00 PM 15:00"
/* | Date Time Time {$$ = NewRangeWithStartEndDateTimes(NewDateTime($1, $3, ""), NewDateTime($1, $3, ""))} */
;


RangePrefixOpt:

| RangePrefix
;


RangePrefix:
  BEGINNING
| FROM
;


DateTimeTZ:
  Date {$$ = NewDateTimeWithDate($1)}
| Date DateTimeSepOpt Time {$$ = NewDateTime($1, $3, "")}
;

Date:
  WeekDay CommaOpt Date {$$ = $3}

  // "2006-01-02 T 15:04:05Z07:00"
| Date T

  // "02.03", but ambiguous between North America (month-day-year) and other (day-month-year) styles.
| Day DateSep Day {$$ = NewAmbiguousDate($1, $3, "")}

| YEAR {$$ = NewDMYDate("", "", $1)}
| YEAR DateSep INT {$$ = NewDMYDate("", $3, $1)}
| YEAR DateSep INT DateSep Day {$$ = NewDMYDate($5, $3, $1)}

  // "Feb 3 2023"
| MONTH_NAME Day YEAR {$$ = NewMDYDate($1, $2, $3)}

  // "3 Feb 2023"
| Day MONTH_NAME YEAR {$$ = NewDMYDate($1, $2, $3)}

  // "2/3/2023", but ambiguous between North America (month-day-year) and other (day-month-year) styles.
| INT DateSep INT DateSep YEAR {$$ = NewAmbiguousDate($1, $3, $5)}

  // "Feb 3"
| MONTH_NAME Day {$$ = NewMDYDate($1, $2, "")}

  // "3 Feb"
| Day MONTH_NAME {$$ = NewDMYDate($1, $2, "")}
;


Day:
  INT OrdinalIndicatorOpt OfOpt
;


OfOpt:

| OF
;


OrdinalIndicatorOpt:

| ORD_IND
| TH
;


WeekDay:
  TH
| WEEKDAY_NAME
;

/* Time */
Time:
  /* INT {$$ = NewTime($1, "", "", "")} */
  // "11am"
  INT AM {$$ = NewTime($1, "", "", "")}
| INT PM {$$ = NewTime((mustAtoi($1) % 12) + 12, "", "", "")}

  // "12:00"
| INT TimeSepOpt INT {$$ = NewTime($1, $3, "", "")}

  // "12:00:00"
| INT TimeSepOpt INT TimeSepOpt INT {$$ = NewTime((mustAtoi($1) % 12) + 12, $3, $5, "")}

  // "9:00 AM"
| INT TimeSepOpt INT AM {$$ = NewTime($1, $3, "", "")}

  // "12:00 PM"
| INT TimeSepOpt INT PM {$$ = NewTime((mustAtoi($1) % 12) + 12, $3, "", "")}

/* // "Feb 3 2023 11am PST" */
/* |   time TimeZoneOpt { */
/*        $$ = $1} //civil.Time{Hour: $1}} //, ampm: $5, timezone: $6}} */
;


/* Separators */

DateSep:
  DEC
| PERIOD
| SUB
| QUO
;


DateTimeSepOpt:

| AT
| DEC
| SUB
;


RangeSep:
  DEC
| SUB
| THROUGH
| TILL
| TO
| UNTIL
;


TimeSepOpt:

| COLON
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
