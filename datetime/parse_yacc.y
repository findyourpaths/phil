%{
package datetime

import "cloud.google.com/go/civil"
%}

%token ILLEGAL

%token ADD
%token AM
%token AMP
%token AND
%token AT
%token BEGINNING
%token CALENDAR
%token COLON
%token COMMA
%token DATE
%token DEC
%token FROM
%token GOOGLE
%token ICS
%token IN
%token LPAREN
%token OF
%token ON
%token ORD_IND
%token PM
%token PERIOD
%token QUO
%token RPAREN
%token SEMICOLON
%token SUB
%token THROUGH
%token T
%token TH
%token TILL
%token TIME
%token TO
%token UNTIL
%token WHEN
%token Z

%token <string> IDENT
%token <string> INT
%token <string> MONTH_NAME
%token <string> TIME_ZONE
%token <string> TIME_ZONE_ABBREV
%token <string> WEEKDAY_NAME
%token <string> YEAR


 /* Type of each nonterminal. */
%type <DateTimeTZRanges> root
%type <DateTimeTZRanges> DateTimeTZRanges
%type <DateTimeTZRange> DateTimeTZRange
%type <DateTimeTZ> DateTimeTZ
%type <DateTimeTZ> RFC3339DateTimeTZ
%type <Date> Date
%type <Date> RFC3339Date
%type <Time> Time
%type <Time> RFC3339Time
%type <TimeZone> TimeZone
%type <TimeZone> TimeZoneOpt
%type <TimeZone> RFC3339TimeZone

%type <string> Day
%type <string> Month
%type <string> Year
%type <strings> DayPlus
%type <strings> DayPlus1




%start root


 /* Type of each nonterminal. */
%union {
    DateTimeTZRanges *DateTimeTZRanges
    DateTimeTZRange *DateTimeTZRange
    DateTimeTZ *DateTimeTZ
    Date civil.Date
    Time civil.Time
    TimeZone  *TimeZone
    string string
    strings []string
    }

%%

root:
  DateTimeTZRanges {$$ = $1}
  DateTimeTZRanges RootSuffixPlus {$$ = $1}
;


RootSuffixPlus:
  RootSuffix
| RootSuffixPlus RootSuffix
;
RootSuffix:
  GOOGLE
| CALENDAR
| ICS
;


DateTimeTZRanges:
  DateTimeTZRange {$$ = &DateTimeTZRanges{Items: []*DateTimeTZRange{$1}}}

| RangesPrefixPlus DateTimeTZRanges {$$ = $2}
| DateTimeTZRanges RangesSepPlus DateTimeTZRange {$$ = AppendDateTimeTZRanges($1, $3)}

  // "Feb 3, 4"
| Month DayPlus1 {$$ = NewRangesWithStartDates(NewMDsYDates($1, $2, nil)...)}
  // "3, 4 Feb"
| DayPlus1 Month {$$ = NewRangesWithStartDates(NewDsMYDates($1, $2, nil)...)}

  // "Feb 3, 4 2023"
| Month DayPlus1 YEAR {$$ = NewRangesWithStartDates(NewMDsYDates($1, $2, $3)...)}
  // "Feb 3, 4 and Mar 5 2023"
| Month DayPlus AND Month DayPlus YEAR {$$ = NewRangesWithStartDates(append(NewMDsYDates($1, $2, $6), NewMDsYDates($4, $5, $6)...)...)}
  // "3, 4 Feb 2023"
| DayPlus1 Month YEAR {$$ = NewRangesWithStartDates(NewDsMYDates($1, $2, $3)...)}
  // "3, 4 Feb and 5 Mar 2023"
| DayPlus Month AND DayPlus Month YEAR {$$ = NewRangesWithStartDates(append(NewDsMYDates($1, $2, $6), NewDsMYDates($4, $5, $6)...)...)}

  // "Feb 1-2, 3-4"
| Month Day RangeSep Day Day RangeSep Day {$$ = NewRanges(NewRangeWithStartEndDates(NewMDYDate($1, $2, nil), NewMDYDate($1, $4, nil)), NewRangeWithStartEndDates(NewMDYDate($1, $5, nil), NewMDYDate($1, $7, nil)))}
  // "1-2, 3-4 Feb"
| Day RangeSep Day Day RangeSep Day Month {$$ = NewRanges(NewRangeWithStartEndDates(NewDMYDate($1, $7, nil), NewDMYDate($3, $7, nil)), NewRangeWithStartEndDates(NewDMYDate($4, $7, nil), NewDMYDate($6, $7, nil)))}

  // "Feb 1-2, Mar 3-4"
| Month Day RangeSep Day Month Day RangeSep Day {$$ = NewRanges(NewRangeWithStartEndDates(NewMDYDate($1, $2, nil), NewMDYDate($1, $4, nil)), NewRangeWithStartEndDates(NewMDYDate($5, $6, nil), NewMDYDate($5, $8, nil)))}
  // "1-2 Feb, 3-4 Mar"
| Day RangeSep Day Month Day RangeSep Day Month {$$ = NewRanges(NewRangeWithStartEndDates(NewDMYDate($1, $4, nil), NewDMYDate($3, $4, nil)), NewRangeWithStartEndDates(NewDMYDate($5, $8, nil), NewDMYDate($7, $8, nil)))}

  // "Feb 1-2, 3-4 2023"
| Month Day RangeSep Day Day RangeSep Day YEAR {$$ = NewRanges(NewRangeWithStartEndDates(NewMDYDate($1, $2, $8), NewMDYDate($1, $4, $8)), NewRangeWithStartEndDates(NewMDYDate($1, $5, $8), NewMDYDate($1, $7, $8)))}
  // "1-2, 3-4 Feb 2023"
| Day RangeSep Day Day RangeSep Day Month YEAR {$$ = NewRanges(NewRangeWithStartEndDates(NewDMYDate($1, $7, $8), NewDMYDate($3, $7, $8)), NewRangeWithStartEndDates(NewDMYDate($4, $7, $8), NewDMYDate($6, $7, $8)))}

  // "Feb 1-2, Mar 3-4 2023"
| Month Day RangeSep Day Month Day RangeSep Day YEAR {$$ = NewRanges(NewRangeWithStartEndDates(NewMDYDate($1, $2, $9), NewMDYDate($1, $4, $9)), NewRangeWithStartEndDates(NewMDYDate($5, $6, $9), NewMDYDate($5, $8, $9)))}
  // "1-2 Feb, 3-4 Mar 2023"
| Day RangeSep Day Month Day RangeSep Day Month YEAR {$$ = NewRanges(NewRangeWithStartEndDates(NewDMYDate($1, $4, $9), NewDMYDate($3, $4, $9)), NewRangeWithStartEndDates(NewDMYDate($5, $8, $9), NewDMYDate($7, $8, $9)))}

  // "Feb 3, Mar 4"
| Month Day Month Day {$$ = NewRanges(NewRangeWithStart(NewMDYDate($1, $2, nil)), NewRangeWithStart(NewMDYDate($3, $4, nil)))}
  // "Feb 3, Mar 4 2023"
| Month Day Month Day YEAR {$$ = NewRanges(NewRangeWithStart(NewMDYDate($1, $2, $5)), NewRangeWithStart(NewMDYDate($3, $4, $5)))}
;


RangesPrefixPlus:
  RangesPrefix
| RangesPrefixPlus RangesPrefix
;
RangesPrefix:
  WHEN
;


RangesSepPlus:
  RangesSep
| RangesSepPlus RangesSep
;
RangesSep:
  AND
| COMMA
;


DayPlus:
  Day {$$ = []string{$1}}
| DayPlus1
;


DayPlus1:
  Day Day {$$ = []string{$1, $2}}
| DayPlus1 Day {$$ = append($1, $2)}
;


//
// Date Time Range
//

DateTimeTZRange:
  DateTimeTZ {$$ = &DateTimeTZRange{Start: $1}}

| RangePrefixPlus DateTimeTZRange {$$ = $2}
| DateTimeTZ RangeSepPlus Time {$$ = NewRangeWithStartEndDateTimes($1, NewDateTime($1.Date, $3, $1.TimeZone))}
| DateTimeTZ RangeSepPlus Time TimeZone {$$ = NewRangeWithStartEndDateTimes(NewDateTime($1.Date, $1.Time, $4), NewDateTime($1.Date, $3, $4))}
| Time RangeSepPlus DateTimeTZ {$$ = NewRangeWithStartEndDateTimes(NewDateTime($3.Date, $1, $3.TimeZone), $3)}

  // "Feb 3, 2023 - Feb 4, 2023"
| DateTimeTZ RangeSepPlus DateTimeTZ {$$ = &DateTimeTZRange{Start: $1, End: $3}}

  // "Feb 3-4"
| Month Day RangeSepPlus Day {$$ = NewRangeWithStartEndDates(NewMDYDate($1, $2, nil), NewMDYDate($1, $4, nil))}
  // "3-4 Feb"
| Day RangeSepPlus Day Month {$$ = NewRangeWithStartEndDates(NewDMYDate($1, $4, nil), NewDMYDate($3, $4, nil))}

  // "Feb 3-4, 2023"
| Month Day RangeSepPlus Day YEAR {$$ = NewRangeWithStartEndDates(NewMDYDate($1, $2, $5), NewMDYDate($1, $4, $5))}
  // "3-4 Feb 2023"
| Day RangeSepPlus Day Month YEAR {$$ = NewRangeWithStartEndDates(NewDMYDate($1, $4, $5), NewDMYDate($3, $4, $5))}

  // "Feb 3 - Mar 4, 2023"
| Month Day RangeSepPlus Month Day YEAR {$$ = NewRangeWithStartEndDates(NewMDYDate($1, $2, $6), NewMDYDate($4, $5, $6))}

  // "Thu Feb 3 - Sat Mar 4, 2023"
| WeekDay Month Day RangeSepPlus WeekDay Month Day YEAR {$$ = NewRangeWithStartEndDates(NewMDYDate($2, $3, $8), NewMDYDate($6, $7, $8))}

  // "Thu Feb 3 - Sat 4 Mar, 2023"
| WeekDay Month Day RangeSepPlus WeekDay Day Month YEAR {$$ = NewRangeWithStartEndDates(NewMDYDate($2, $3, $8), NewDMYDate($6, $7, $8))}

  // "Thu Feb 3 - Sat 4 Mar, 2023"
| WeekDay Month Day RangeSepPlus Day WeekDay Month YEAR {$$ = NewRangeWithStartEndDates(NewMDYDate($2, $3, $8), NewDMYDate($5, $7, $8))}

  // "Thu 3 Feb - Sat Mar 4, 2023"
| WeekDay Day Month RangeSepPlus WeekDay Month Day YEAR {$$ = NewRangeWithStartEndDates(NewDMYDate($2, $3, $8), NewMDYDate($6, $7, $8))}

  // "Thu 3 Feb - Sat 4 Mar, 2023"
| WeekDay Day Month RangeSepPlus WeekDay Day Month YEAR {$$ = NewRangeWithStartEndDates(NewDMYDate($2, $3, $8), NewDMYDate($6, $7, $8))}

  // "Thu 3 Feb - Sat 4 Mar, 2023"
| WeekDay Day Month RangeSepPlus Day WeekDay Month YEAR {$$ = NewRangeWithStartEndDates(NewDMYDate($2, $3, $8), NewDMYDate($5, $7, $8))}

  // "9:00am 3rd Feb - 4th Feb 3:00pm 2023"
| Time TimeZoneOpt DateTimeSepOpt Day DateSepOpt Month RangeSepPlus Day DateSepOpt Month DateTimeSepOpt Time TimeZoneOpt YEAR {$$ = NewRangeWithStartEndDateTimes(NewDateTime(NewDMYDate($4, $6, $14), $1, $2), NewDateTime(NewDMYDate($8, $10, $14), $12, $13))}

  // "Feb 3 2023 9:00 AM 09:00"
  // "Feb 3 2023 3:00 PM 15:00"
| Date Time TimeZoneOpt Time TimeZoneOpt {$$ = NewRangeWithStartEndDateTimes(NewDateTime($1, $2, $3), NewDateTime($1, $4, $5))}
;


RangePrefixPlus:
  RangePrefix
| RangePrefixPlus RangePrefix
;
RangePrefix:
  BEGINNING
| FROM
;


RangeSepPlus:
  RangeSep
| RangeSepPlus RangeSep
;
RangeSep:
  DEC
| SUB
| THROUGH
| TILL
| TO
| UNTIL
;


//
// Date Time
//

DateTimeTZ:
  Date {$$ = NewDateTimeWithDate($1, nil)}

| Date Time TimeZoneOpt {$$ = NewDateTime($1, $2, $3)}
| Date DateTimeSepPlus Time TimeZoneOpt {$$ = NewDateTime($1, $3, $4)}
| Time TimeZoneOpt Date {$$ = NewDateTime($3, $1, $2)}
| Time TimeZoneOpt DateTimeSepPlus Date {$$ = NewDateTime($4, $1, $2)}
| RFC3339DateTimeTZ
;


DateTimeSepOpt:

| DateTimeSepPlus
;
DateTimeSepPlus:
  DateTimeSep
| DateTimeSepPlus DateTimeSep
;
DateTimeSep:
  DEC
| SUB
| T
;


TimeZoneOpt:
  {$$ = nil}
| TimeZone
| TimeZonePrefix TimeZone TimeZoneSuffix {$$ = $2}
| TimeZoneSep TimeZone {$$ = $2}
| Z {$$ = nil}
;
TimeZonePrefix:
  LPAREN
;
TimeZoneSuffix:
  RPAREN
;
TimeZoneSep:
  IN
| SUB
;
TimeZone:
  TIME_ZONE {$$ = NewTimeZone($1, nil, nil)}
| TIME_ZONE_ABBREV {$$ = NewTimeZone(nil, $1, nil)}
;


//
// RDC3339 DateTimeTZ
//

RFC3339DateTimeTZ:
  RFC3339Date {$$ = NewDateTimeWithDate($1, nil)}
| RFC3339Date RFC3339Time {$$ = NewDateTime($1, $2, nil)}
| RFC3339Date RFC3339Time RFC3339TimeZone {$$ = NewDateTime($1, $2, $3)}
;

RFC3339Date:
  Year SUB INT SUB INT {$$ = NewDMYDate($5, $3, $1)}
;

RFC3339Time:
  T INT COLON INT COLON INT {$$ = NewTime($2, $4, $6, nil)}
;

RFC3339TimeZone:
  Z {$$ = nil}
| ADD INT COLON INT {$$ = NewTimeZone(nil, nil, "+" + $2 + ":" + $4)}
| SUB INT COLON INT {$$ = NewTimeZone(nil, nil, "-" + $2 + ":" + $4)}
;


//
// Date
//

Date:
  DatePrefixPlus Date {$$ = $2}

  // "02.03", but ambiguous between North America (month-day-year) and other (day-month-year) styles.
| Day DateSepPlus Day {$$ = NewAmbiguousDate($1, $3, nil)}

| Year {$$ = NewDMYDate(nil, nil, $1)}
| Year DateSepPlus Day {$$ = NewDMYDate(nil, $3, $1)}
| Year DateSepPlus Day DateSepPlus Day {$$ = NewDMYDate($5, $3, $1)}

  // "Feb"
| Month {$$ = NewMDYDate($1, nil, nil)}

  // "Feb 2023"
| Month Year {$$ = NewMDYDate($1, nil, $2)}

  // "Feb 3 2023"
| Month Day Year {$$ = NewMDYDate($1, $2, $3)}

  // "3 Feb 2023"
| Day Month Year {$$ = NewDMYDate($1, $2, $3)}

  // "2/3/2023", but ambiguous between North America (month-day-year) and other (day-month-year) styles.
| Day DateSepPlus Day DateSepPlus Year {$$ = NewAmbiguousDate($1, $3, $5)}

  // "Feb 3"
| Month Day {$$ = NewMDYDate($1, $2, nil)}

  // "3 Feb"
| Day Month {$$ = NewDMYDate($1, $2, nil)}

  // "2023 Feb 3"
| Year Month Day {$$ = NewDMYDate($3, $2, $1)}
;


DatePrefixPlus:
  DatePrefix
| DatePrefixPlus DatePrefix
;
DatePrefix:
  DATE
| WeekDay
| COLON
| COMMA
| ON
| TIME
;


DateSepOpt:

| DateSepPlus
;
DateSepPlus:
  DateSep
| DateSepPlus DateSep
;
DateSep:
  COMMA
| DEC
| PERIOD
| SUB
| QUO
;


Day:
  INT
| INT DaySuffixPlus
;
DaySuffixPlus:
  DaySuffix
| DaySuffixPlus DaySuffix
;
DaySuffix:
  COMMA
| ORD_IND
| PERIOD
| TH
;


Month:
  MONTH_NAME
| MONTH_NAME MonthSuffixPlus
;
MonthSuffixPlus:
  MonthSuffix
| MonthSuffixPlus MonthSuffix
;
MonthSuffix:
  COMMA
;


Year:
  YEAR
| INT
| Year YearSuffixPlus
;
YearSuffixPlus:
  YearSuffix
| YearSuffixPlus YearSuffix
;
YearSuffix:
  COMMA
;


WeekDay:
  TH
| WEEKDAY_NAME
;


//
// Time
//

Time:
  TimePrefixPlus Time {$$ = $2}
  /* INT {$$ = NewTime($1, nil, nil, nil)} */
  // "11am"
| INT AM {$$ = NewAMTime($1, nil, nil, nil)}
| INT PM {$$ = NewPMTime($1, nil, nil, nil)}

  // "12:00"
| INT TimeSep INT {$$ = NewTime($1, $3, nil, nil)}

  // "12:00:00"
| INT TimeSep INT TimeSep INT {$$ = NewTime($1, $3, $5, nil)}

  // "9:00 AM"
| INT TimeSep INT AM {$$ = NewAMTime($1, $3, nil, nil)}

  // "12:00 PM"
| INT TimeSep INT PM {$$ = NewPMTime($1, $3, nil, nil)}

/* // "Feb 3 2023 11am PST" */
/* |   time TimeZoneOpt { */
/*        $$ = $1} //civil.Time{Hour: $1}} //, ampm: $5, timezone: $6}} */
;


TimePrefixPlus:
  TimePrefix
| TimePrefixPlus TimePrefix
;
TimePrefix:
  AT
| COLON
| FROM
| ON
| TIME
;


TimeSep:
  COLON
| PERIOD
;
