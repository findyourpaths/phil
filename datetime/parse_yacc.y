%{
package datetime
%}

%token ILLEGAL

%token A
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
%token DATES
%token DEC
%token FROM
%token GOOGLE
%token ICS
%token IN
%token LPAREN
%token OF
%token ON
%token ORD_IND
%token PART
%token P
%token PM
%token PERIOD
%token QUO
%token RPAREN
%token SAVE
%token SEMICOLON
%token SUB
%token THROUGH
%token T
%token TH
%token THE
%token TILL
%token TIME
%token TO
%token UNTIL
%token WHEN
%token Z

%token <string> IDENT
%token <string> INT
%token <string> MONTH_NAME
%token <string> TIME_NAME
%token <string> TIME_ZONE
%token <string> TIME_ZONE_ABBREV
%token <string> RELATIVE_DAY
%token <string> WEEKDAY_NAME
%token <string> YEAR


 /* Type of each nonterminal. */
%type <DateTimeRanges> root
%type <DateTimeRanges> DateTimeRanges

%type <DateTimeRange> DateTimeRange

%type <DateTime> DateTime
%type <DateTime> RFC3339DateTime

%type <Date> Date
%type <Date> RFC3339Date

%type <Time> Time
%type <Time> RFC3339Time

%type <TimeZone> TimeZone
%type <TimeZone> TimeZoneOpt
%type <TimeZone> RFC3339TimeZone

%type <string> Day
%type <string> Month
%type <string> Weekday
%type <string> WeekdayOpt
%type <string> Year

%type <strings> DayPlus
%type <strings> DayPlus1


%start root

 /* Type of each nonterminal. */
%union {
    DateTimeRanges *DateTimeRanges
    DateTimeRange *DateTimeRange
    DateTime *DateTime
    Date *Date
    Time *Time
    TimeZone  *TimeZone
    string string
    strings []string
    }

%%

root:
  DateTimeRanges {$$ = $1}
| root RootSuffixPlus {$$ = $1}
| RootPrefixPlus root {$$ = $2}
;


RootPrefixPlus:
  RootPrefix
| RootPrefixPlus RootPrefix
;
RootPrefix:
  PART INT ColonOpt
| SAVE THE DATES ColonOpt
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


ColonOpt:

| COLON
;


DateTimeRanges:
  DateTimeRange {$$ = &DateTimeRanges{Items: []*DateTimeRange{$1}}}

| RangesPrefixPlus DateTimeRanges {$$ = $2}
| DateTimeRanges RangesSepPlus DateTimeRange {$$ = AppendDateTimeRanges($1, $3)}

  // "Feb 3, 4"
| Month DayPlus1 {$$ = NewRangesWithStartDates(NewRawDateFromMDsYs($1, $2, nil)...)}
  // "3, 4 Feb"
| DayPlus1 Month {$$ = NewRangesWithStartDates(NewRawDateFromDsMYs($1, $2, nil)...)}

  // "Feb 3, 4 2023"
| Month DayPlus1 Year {$$ = NewRangesWithStartDates(NewRawDateFromMDsYs($1, $2, $3)...)}
  // "Feb 3, 4 and Mar 5 2023"
| Month DayPlus AND Month DayPlus Year {$$ = NewRangesWithStartDates(append(NewRawDateFromMDsYs($1, $2, $6), NewRawDateFromMDsYs($4, $5, $6)...)...)}
  // "3, 4 Feb 2023"
| DayPlus1 Month Year {$$ = NewRangesWithStartDates(NewRawDateFromDsMYs($1, $2, $3)...)}
  // "3, 4 Feb and 5 Mar 2023"
| DayPlus Month AND DayPlus Month Year {$$ = NewRangesWithStartDates(append(NewRawDateFromDsMYs($1, $2, $6), NewRawDateFromDsMYs($4, $5, $6)...)...)}

  // "Feb 1-2, 3-4"
| Month Day RangeSep Day Day RangeSep Day {$$ = NewRanges(NewRangeWithStartEndDates(NewRawDateFromMDY($1, $2, nil), NewRawDateFromMDY($1, $4, nil)), NewRangeWithStartEndDates(NewRawDateFromMDY($1, $5, nil), NewRawDateFromMDY($1, $7, nil)))}
  // "1-2, 3-4 Feb"
| Day RangeSep Day Day RangeSep Day Month {$$ = NewRanges(NewRangeWithStartEndDates(NewRawDateFromDMY($1, $7, nil), NewRawDateFromDMY($3, $7, nil)), NewRangeWithStartEndDates(NewRawDateFromDMY($4, $7, nil), NewRawDateFromDMY($6, $7, nil)))}

  // "Feb 1-2, Mar 3-4"
| Month Day RangeSep Day Month Day RangeSep Day {$$ = NewRanges(NewRangeWithStartEndDates(NewRawDateFromMDY($1, $2, nil), NewRawDateFromMDY($1, $4, nil)), NewRangeWithStartEndDates(NewRawDateFromMDY($5, $6, nil), NewRawDateFromMDY($5, $8, nil)))}
  // "1-2 Feb, 3-4 Mar"
| Day RangeSep Day Month Day RangeSep Day Month {$$ = NewRanges(NewRangeWithStartEndDates(NewRawDateFromDMY($1, $4, nil), NewRawDateFromDMY($3, $4, nil)), NewRangeWithStartEndDates(NewRawDateFromDMY($5, $8, nil), NewRawDateFromDMY($7, $8, nil)))}

  // "Feb 1-2, 3-4 2023"
| Month Day RangeSep Day Day RangeSep Day Year {$$ = NewRanges(NewRangeWithStartEndDates(NewRawDateFromMDY($1, $2, $8), NewRawDateFromMDY($1, $4, $8)), NewRangeWithStartEndDates(NewRawDateFromMDY($1, $5, $8), NewRawDateFromMDY($1, $7, $8)))}
  // "1-2, 3-4 Feb 2023"
| Day RangeSep Day Day RangeSep Day Month Year {$$ = NewRanges(NewRangeWithStartEndDates(NewRawDateFromDMY($1, $7, $8), NewRawDateFromDMY($3, $7, $8)), NewRangeWithStartEndDates(NewRawDateFromDMY($4, $7, $8), NewRawDateFromDMY($6, $7, $8)))}

  // "Feb 1-2, Mar 3-4 2023"
| Month Day RangeSep Day Month Day RangeSep Day Year {$$ = NewRanges(NewRangeWithStartEndDates(NewRawDateFromMDY($1, $2, $9), NewRawDateFromMDY($1, $4, $9)), NewRangeWithStartEndDates(NewRawDateFromMDY($5, $6, $9), NewRawDateFromMDY($5, $8, $9)))}
  // "1-2 Feb, 3-4 Mar 2023"
| Day RangeSep Day Month Day RangeSep Day Month Year {$$ = NewRanges(NewRangeWithStartEndDates(NewRawDateFromDMY($1, $4, $9), NewRawDateFromDMY($3, $4, $9)), NewRangeWithStartEndDates(NewRawDateFromDMY($5, $8, $9), NewRawDateFromDMY($7, $8, $9)))}

  // "Feb 3, Mar 4"
| Month Day Month Day {$$ = NewRanges(NewRangeWithStart(NewRawDateFromMDY($1, $2, nil)), NewRangeWithStart(NewRawDateFromMDY($3, $4, nil)))}
  // "Feb 3, Mar 4 2023"
| Month Day Month Day Year {$$ = NewRanges(NewRangeWithStart(NewRawDateFromMDY($1, $2, $5)), NewRangeWithStart(NewRawDateFromMDY($3, $4, $5)))}

  // "Wednesdays February 1st & 8th 12:00p-3:00p"
| WeekdayOpt Month Day DateSepOpt Day Time DateTimeSepOpt Time {$$ = NewRanges(NewRange(NewDateTime(NewRawDateFromWMDY($1, $2, $3, nil), $6, nil), NewDateTime(NewRawDateFromWMDY($1, $2, $3, nil), $8, nil)), NewRange(NewDateTime(NewRawDateFromWMDY($1, $2, $5, nil), $6, nil), NewDateTime(NewRawDateFromWMDY($1, $2, $5, nil), $8, nil)))}
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

DateTimeRange:
  DateTime {$$ = &DateTimeRange{Start: $1}}

| RangePrefixPlus DateTimeRange {$$ = $2}
| DateTime RangeSepPlus Time {$$ = NewRange($1, NewDateTime($1.Date, $3, $1.TimeZone))}
| DateTime RangeSepPlus Time TimeZone {$$ = NewRange(NewDateTime($1.Date, $1.Time, $4), NewDateTime($1.Date, $3, $4))}
| Time RangeSepPlus DateTime {$$ = NewRange(NewDateTime($3.Date, $1, $3.TimeZone), $3)}

  // "Feb 3, 2023 - Feb 4, 2023"
| DateTime RangeSepPlus DateTime {$$ = &DateTimeRange{Start: $1, End: $3}}

  // "Feb 3-4"
| Month Day RangeSepPlus Day {$$ = NewRangeWithStartEndDates(NewRawDateFromMDY($1, $2, nil), NewRawDateFromMDY($1, $4, nil))}
  // "3-4 Feb"
| Day RangeSepPlus Day Month {$$ = NewRangeWithStartEndDates(NewRawDateFromDMY($1, $4, nil), NewRawDateFromDMY($3, $4, nil))}

  // "Feb 3-4, 2023"
| Month Day RangeSepPlus Day Year {$$ = NewRangeWithStartEndDates(NewRawDateFromMDY($1, $2, $5), NewRawDateFromMDY($1, $4, $5))}
  // "3-4 Feb 2023"
| Day RangeSepPlus Day Month Year {$$ = NewRangeWithStartEndDates(NewRawDateFromDMY($1, $4, $5), NewRawDateFromDMY($3, $4, $5))}

  // "Feb 3 - Mar 4, 2023"
| Month Day RangeSepPlus Month Day Year {$$ = NewRangeWithStartEndDates(NewRawDateFromMDY($1, $2, $6), NewRawDateFromMDY($4, $5, $6))}

  // "Thu Feb 3 - Sat Mar 4, 2023"
| WeekdayOpt Month Day RangeSepPlus WeekdayOpt Month Day Year {$$ = NewRangeWithStartEndDates(NewRawDateFromWMDY($1, $2, $3, $8), NewRawDateFromWMDY($5, $6, $7, $8))}

  // "Thu Feb 3 - Sat 4 Mar, 2023"
| WeekdayOpt Month Day RangeSepPlus WeekdayOpt Day Month Year {$$ = NewRangeWithStartEndDates(NewRawDateFromWMDY($1, $2, $3, $8), NewRawDateFromWDMY($5, $6, $7, $8))}

  // "Thu Feb 3 - 4 Sat Mar, 2023"
| WeekdayOpt Month Day RangeSepPlus Day WeekdayOpt Month Year {$$ = NewRangeWithStartEndDates(NewRawDateFromWMDY($1, $2, $3, $8), NewRawDateFromWDMY($6, $5, $7, $8))}

  // "Thu 3 Feb - Sat Mar 4, 2023"
| WeekdayOpt Day Month RangeSepPlus WeekdayOpt Month Day Year {$$ = NewRangeWithStartEndDates(NewRawDateFromWDMY($1, $2, $3, $8), NewRawDateFromWMDY($5, $6, $7, $8))}

  // "Thu 3 Feb - Sat 4 Mar, 2023"
| WeekdayOpt Day Month RangeSepPlus WeekdayOpt Day Month Year {$$ = NewRangeWithStartEndDates(NewRawDateFromWDMY($1, $2, $3, $8), NewRawDateFromWDMY($5, $6, $7, $8))}

  // "Thu 3 Feb - Sat 4 Mar, 2023"
| WeekdayOpt Day Month RangeSepPlus Day WeekdayOpt Month Year {$$ = NewRangeWithStartEndDates(NewRawDateFromWDMY($1, $2, $3, $8), NewRawDateFromWDMY($6, $5, $7, $8))}

  // "9:00am 3rd Feb - 4th Feb 3:00pm 2023"
| Time TimeZoneOpt DateTimeSepOpt Day DateSepOpt Month RangeSepPlus Day DateSepOpt Month DateTimeSepOpt Time TimeZoneOpt Year {$$ = NewRange(NewDateTime(NewRawDateFromDMY($4, $6, $14), $1, $2), NewDateTime(NewRawDateFromDMY($8, $10, $14), $12, $13))}

  // "Feb 3 2023 9:00 AM 09:00"
  // "Feb 3 2023 3:00 PM 15:00"
  // "February 3rd, 9-12pm ET"
| Date Time TimeZoneOpt RangeSepOpt Time TimeZoneOpt {$$ = NewRange(NewDateTime($1, $2, $3), NewDateTime($1, $5, $6))}
;


RangePrefixPlus:
  RangePrefix
| RangePrefixPlus RangePrefix
;
RangePrefix:
  BEGINNING
| FROM
;


RangeSepOpt:

| RangeSepPlus
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

DateTime:
  Date {$$ = NewDateTimeWithDate($1, nil)}

| Date Time TimeZoneOpt {$$ = NewDateTime($1, $2, $3)}
| Date DateTimeSepPlus Time TimeZoneOpt {$$ = NewDateTime($1, $3, $4)}
| Time TimeZoneOpt Date {$$ = NewDateTime($3, $1, $2)}
| Time TimeZoneOpt DateTimeSepPlus Date {$$ = NewDateTime($4, $1, $2)}
| RFC3339DateTime
;


DateTimeSepOpt:

| DateTimeSepPlus
;
DateTimeSepPlus:
  DateTimeSep
| DateTimeSepPlus DateTimeSep
;
DateTimeSep:
  COLON
| COMMA
| DEC
| QUO
| SUB
| T
;


TimeZoneOpt:
  {$$ = nil}
| TimeZone TimeZoneSuffixOpt {$$ = $1}
| LPAREN TimeZone RPAREN {$$ = $2}
| TimeZoneSep TimeZone {$$ = $2}
| Z {$$ = nil}
;

TimeZoneSuffixOpt:

| TIME
;
TimeZoneSep:
  IN
| SUB
;
TimeZone:
  TIME_ZONE {$$ = NewTimeZone($1, nil, nil)}
| TIME_ZONE_ABBREV {$$ = NewTimeZone(nil, $1, nil)}
| P {$$ = NewTimeZone(nil, "P", nil)}
;


//
// RDC3339 DateTime
//

RFC3339DateTime:
  RFC3339Date {$$ = NewDateTimeWithDate($1, nil)}
| RFC3339Date RFC3339Time {$$ = NewDateTime($1, $2, nil)}
| RFC3339Date RFC3339Time RFC3339TimeZone {$$ = NewDateTime($1, $2, $3)}
;

RFC3339Date:
  Year SUB INT SUB INT {$$ = NewRawDateFromDMY($5, $3, $1)}
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

| RELATIVE_DAY {$$ = NewRawDateFromRelative($1)}
| RELATIVE_DAY WeekdayOpt {$$ = NewRawDateFromRelative($1)}
| RELATIVE_DAY Date {$$ = $2}

  // "02.03", but ambiguous between North America (month-day-year) and other (day-month-year) styles.
| WeekdayOpt Day DateSepPlus Day {$$ = NewRawDateFromAmbiguous($1, $2, $4, nil)}

| Year {$$ = NewRawDateFromDMY(nil, nil, $1)}
| Year DateSepPlus Day {$$ = NewRawDateFromDMY(nil, $3, $1)}
| Year DateSepPlus Day DateSepPlus Day {$$ = NewRawDateFromDMY($5, $3, $1)}

  // "Feb"
| Month {$$ = NewRawDateFromMDY($1, nil, nil)}

  // "Feb 2023"
| Month Year {$$ = NewRawDateFromMDY($1, nil, $2)}

  // "Feb 3 2023"
| WeekdayOpt Month Day Year {$$ = NewRawDateFromWMDY($1, $2, $3, $4)}

  // "3 Feb 2023"
| WeekdayOpt Day Month Year {$$ = NewRawDateFromWDMY($1, $2, $3, $4)}

  // "2/3/2023", but ambiguous between North America (month-day-year) and other (day-month-year) styles.
| WeekdayOpt Day DateSepPlus Day DateSepPlus Year {$$ = NewRawDateFromAmbiguous($1, $2, $4, $6)}

  // "Feb 3"
| WeekdayOpt Month Day {$$ = NewRawDateFromWMDY($1, $2, $3, nil)}

  // "3 Feb"
| WeekdayOpt Day Month {$$ = NewRawDateFromWDMY($1, $2, $3, nil)}

  // "2023 Feb 3"
| WeekdayOpt Year Month Day {$$ = NewRawDateFromWDMY($1, $4, $3, $2)}
;


DatePrefixPlus:
  DatePrefix
| DatePrefixPlus DatePrefix
;
DatePrefix:
  DATE
| Weekday
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
  AND
| COMMA
| DEC
| PERIOD
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
| Year YearSuffixPlus
;
YearSuffixPlus:
  YearSuffix
| YearSuffixPlus YearSuffix
;
YearSuffix:
  COMMA
;


WeekdayOpt:
  {$$ = ""}
| Weekday
;
Weekday:
  TH {$$ = "TH"}
| WEEKDAY_NAME
;


//
// Time
//

Time:
  TimePrefixPlus Time {$$ = $2}
  /* INT {$$ = NewTime($1, nil, nil, nil)} */
  // "11am"
| INT Am {$$ = NewAMTime($1, nil, nil, nil)}
| INT Pm {$$ = NewPMTime($1, nil, nil, nil)}

  // "12"
| INT {$$ = NewTime($1, nil, nil, nil)}

  // "12:00"
| INT TimeSep INT {$$ = NewTime($1, $3, nil, nil)}

  // "12:00:00"
| INT TimeSep INT TimeSep INT {$$ = NewTime($1, $3, $5, nil)}

  // "9:00 AM"
| INT TimeSep INT Am {$$ = NewAMTime($1, $3, nil, nil)}

  // "12:00 PM"
| INT TimeSep INT Pm {$$ = NewPMTime($1, $3, nil, nil)}

| TIME_NAME {$$ = NewTime($1, nil, nil, nil)}

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
| TIME COLON
| FROM
| ON
| TIME
;


TimeSep:
  COLON
| PERIOD
;


Am:
  A
| AM
;


Pm:
  P
| PM
;
