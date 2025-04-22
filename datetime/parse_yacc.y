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
%type <string> Weekday
%type <string> WeekdayOpt
%type <string> Year

%type <strings> DayPlus
%type <strings> DayPlus1


%start root

 /* Type of each nonterminal. */
%union {
    DateTimeTZRanges *DateTimeTZRanges
    DateTimeTZRange *DateTimeTZRange
    DateTimeTZ *DateTimeTZ
    Date *Date
    Time *Time
    TimeZone  *TimeZone
    string string
    strings []string
    }

%%

root:
  DateTimeTZRanges {$$ = $1}
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


DateTimeTZRanges:
  DateTimeTZRange {$$ = &DateTimeTZRanges{Items: []*DateTimeTZRange{$1}}}

| RangesPrefixPlus DateTimeTZRanges {$$ = $2}
| DateTimeTZRanges RangesSepPlus DateTimeTZRange {$$ = AppendDateTimeTZRanges($1, $3)}

  // "Feb 3, 4"
| Month DayPlus1 {$$ = NewRangesWithStartDates(NewMDsYDates($1, $2, nil)...)}
  // "3, 4 Feb"
| DayPlus1 Month {$$ = NewRangesWithStartDates(NewDsMYDates($1, $2, nil)...)}

  // "Feb 3, 4 2023"
| Month DayPlus1 Year {$$ = NewRangesWithStartDates(NewMDsYDates($1, $2, $3)...)}
  // "Feb 3, 4 and Mar 5 2023"
| Month DayPlus AND Month DayPlus Year {$$ = NewRangesWithStartDates(append(NewMDsYDates($1, $2, $6), NewMDsYDates($4, $5, $6)...)...)}
  // "3, 4 Feb 2023"
| DayPlus1 Month Year {$$ = NewRangesWithStartDates(NewDsMYDates($1, $2, $3)...)}
  // "3, 4 Feb and 5 Mar 2023"
| DayPlus Month AND DayPlus Month Year {$$ = NewRangesWithStartDates(append(NewDsMYDates($1, $2, $6), NewDsMYDates($4, $5, $6)...)...)}

  // "Feb 1-2, 3-4"
| Month Day RangeSep Day Day RangeSep Day {$$ = NewRanges(NewRangeWithStartEndDates(NewMDYDate($1, $2, nil), NewMDYDate($1, $4, nil)), NewRangeWithStartEndDates(NewMDYDate($1, $5, nil), NewMDYDate($1, $7, nil)))}
  // "1-2, 3-4 Feb"
| Day RangeSep Day Day RangeSep Day Month {$$ = NewRanges(NewRangeWithStartEndDates(NewDMYDate($1, $7, nil), NewDMYDate($3, $7, nil)), NewRangeWithStartEndDates(NewDMYDate($4, $7, nil), NewDMYDate($6, $7, nil)))}

  // "Feb 1-2, Mar 3-4"
| Month Day RangeSep Day Month Day RangeSep Day {$$ = NewRanges(NewRangeWithStartEndDates(NewMDYDate($1, $2, nil), NewMDYDate($1, $4, nil)), NewRangeWithStartEndDates(NewMDYDate($5, $6, nil), NewMDYDate($5, $8, nil)))}
  // "1-2 Feb, 3-4 Mar"
| Day RangeSep Day Month Day RangeSep Day Month {$$ = NewRanges(NewRangeWithStartEndDates(NewDMYDate($1, $4, nil), NewDMYDate($3, $4, nil)), NewRangeWithStartEndDates(NewDMYDate($5, $8, nil), NewDMYDate($7, $8, nil)))}

  // "Feb 1-2, 3-4 2023"
| Month Day RangeSep Day Day RangeSep Day Year {$$ = NewRanges(NewRangeWithStartEndDates(NewMDYDate($1, $2, $8), NewMDYDate($1, $4, $8)), NewRangeWithStartEndDates(NewMDYDate($1, $5, $8), NewMDYDate($1, $7, $8)))}
  // "1-2, 3-4 Feb 2023"
| Day RangeSep Day Day RangeSep Day Month Year {$$ = NewRanges(NewRangeWithStartEndDates(NewDMYDate($1, $7, $8), NewDMYDate($3, $7, $8)), NewRangeWithStartEndDates(NewDMYDate($4, $7, $8), NewDMYDate($6, $7, $8)))}

  // "Feb 1-2, Mar 3-4 2023"
| Month Day RangeSep Day Month Day RangeSep Day Year {$$ = NewRanges(NewRangeWithStartEndDates(NewMDYDate($1, $2, $9), NewMDYDate($1, $4, $9)), NewRangeWithStartEndDates(NewMDYDate($5, $6, $9), NewMDYDate($5, $8, $9)))}
  // "1-2 Feb, 3-4 Mar 2023"
| Day RangeSep Day Month Day RangeSep Day Month Year {$$ = NewRanges(NewRangeWithStartEndDates(NewDMYDate($1, $4, $9), NewDMYDate($3, $4, $9)), NewRangeWithStartEndDates(NewDMYDate($5, $8, $9), NewDMYDate($7, $8, $9)))}

  // "Feb 3, Mar 4"
| Month Day Month Day {$$ = NewRanges(NewRangeWithStart(NewMDYDate($1, $2, nil)), NewRangeWithStart(NewMDYDate($3, $4, nil)))}
  // "Feb 3, Mar 4 2023"
| Month Day Month Day Year {$$ = NewRanges(NewRangeWithStart(NewMDYDate($1, $2, $5)), NewRangeWithStart(NewMDYDate($3, $4, $5)))}

  // "Wednesdays February 1st & 8th 12:00p-3:00p"
| WeekdayOpt Month Day DateSepOpt Day Time DateTimeSepOpt Time {$$ = NewRanges(NewRange(NewDateTimeTZ(NewWMDYDate($1, $2, $3, nil), $6, nil), NewDateTimeTZ(NewWMDYDate($1, $2, $3, nil), $8, nil)), NewRange(NewDateTimeTZ(NewWMDYDate($1, $2, $5, nil), $6, nil), NewDateTimeTZ(NewWMDYDate($1, $2, $5, nil), $8, nil)))}
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
| DateTimeTZ RangeSepPlus Time {$$ = NewRange($1, NewDateTimeTZ($1.Date, $3, $1.TimeZone))}
| DateTimeTZ RangeSepPlus Time TimeZone {$$ = NewRange(NewDateTimeTZ($1.Date, $1.Time, $4), NewDateTimeTZ($1.Date, $3, $4))}
| Time RangeSepPlus DateTimeTZ {$$ = NewRange(NewDateTimeTZ($3.Date, $1, $3.TimeZone), $3)}

  // "Feb 3, 2023 - Feb 4, 2023"
| DateTimeTZ RangeSepPlus DateTimeTZ {$$ = &DateTimeTZRange{Start: $1, End: $3}}

  // "Feb 3-4"
| Month Day RangeSepPlus Day {$$ = NewRangeWithStartEndDates(NewMDYDate($1, $2, nil), NewMDYDate($1, $4, nil))}
  // "3-4 Feb"
| Day RangeSepPlus Day Month {$$ = NewRangeWithStartEndDates(NewDMYDate($1, $4, nil), NewDMYDate($3, $4, nil))}

  // "Feb 3-4, 2023"
| Month Day RangeSepPlus Day Year {$$ = NewRangeWithStartEndDates(NewMDYDate($1, $2, $5), NewMDYDate($1, $4, $5))}
  // "3-4 Feb 2023"
| Day RangeSepPlus Day Month Year {$$ = NewRangeWithStartEndDates(NewDMYDate($1, $4, $5), NewDMYDate($3, $4, $5))}

  // "Feb 3 - Mar 4, 2023"
| Month Day RangeSepPlus Month Day Year {$$ = NewRangeWithStartEndDates(NewMDYDate($1, $2, $6), NewMDYDate($4, $5, $6))}

  // "Thu Feb 3 - Sat Mar 4, 2023"
| WeekdayOpt Month Day RangeSepPlus WeekdayOpt Month Day Year {$$ = NewRangeWithStartEndDates(NewWMDYDate($1, $2, $3, $8), NewWMDYDate($5, $6, $7, $8))}

  // "Thu Feb 3 - Sat 4 Mar, 2023"
| WeekdayOpt Month Day RangeSepPlus WeekdayOpt Day Month Year {$$ = NewRangeWithStartEndDates(NewWMDYDate($1, $2, $3, $8), NewWDMYDate($5, $6, $7, $8))}

  // "Thu Feb 3 - 4 Sat Mar, 2023"
| WeekdayOpt Month Day RangeSepPlus Day WeekdayOpt Month Year {$$ = NewRangeWithStartEndDates(NewWMDYDate($1, $2, $3, $8), NewWDMYDate($6, $5, $7, $8))}

  // "Thu 3 Feb - Sat Mar 4, 2023"
| WeekdayOpt Day Month RangeSepPlus WeekdayOpt Month Day Year {$$ = NewRangeWithStartEndDates(NewWDMYDate($1, $2, $3, $8), NewWMDYDate($5, $6, $7, $8))}

  // "Thu 3 Feb - Sat 4 Mar, 2023"
| WeekdayOpt Day Month RangeSepPlus WeekdayOpt Day Month Year {$$ = NewRangeWithStartEndDates(NewWDMYDate($1, $2, $3, $8), NewWDMYDate($5, $6, $7, $8))}

  // "Thu 3 Feb - Sat 4 Mar, 2023"
| WeekdayOpt Day Month RangeSepPlus Day WeekdayOpt Month Year {$$ = NewRangeWithStartEndDates(NewWDMYDate($1, $2, $3, $8), NewWDMYDate($6, $5, $7, $8))}

  // "9:00am 3rd Feb - 4th Feb 3:00pm 2023"
| Time TimeZoneOpt DateTimeSepOpt Day DateSepOpt Month RangeSepPlus Day DateSepOpt Month DateTimeSepOpt Time TimeZoneOpt Year {$$ = NewRange(NewDateTimeTZ(NewDMYDate($4, $6, $14), $1, $2), NewDateTimeTZ(NewDMYDate($8, $10, $14), $12, $13))}

  // "Feb 3 2023 9:00 AM 09:00"
  // "Feb 3 2023 3:00 PM 15:00"
  // "February 3rd, 9-12pm ET"
| Date Time TimeZoneOpt RangeSepOpt Time TimeZoneOpt {$$ = NewRange(NewDateTimeTZ($1, $2, $3), NewDateTimeTZ($1, $5, $6))}
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

DateTimeTZ:
  Date {$$ = NewDateTimeTZWithDate($1, nil)}

| Date Time TimeZoneOpt {$$ = NewDateTimeTZ($1, $2, $3)}
| Date DateTimeSepPlus Time TimeZoneOpt {$$ = NewDateTimeTZ($1, $3, $4)}
| Time TimeZoneOpt Date {$$ = NewDateTimeTZ($3, $1, $2)}
| Time TimeZoneOpt DateTimeSepPlus Date {$$ = NewDateTimeTZ($4, $1, $2)}
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
  COLON
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
// RDC3339 DateTimeTZ
//

RFC3339DateTimeTZ:
  RFC3339Date {$$ = NewDateTimeTZWithDate($1, nil)}
| RFC3339Date RFC3339Time {$$ = NewDateTimeTZ($1, $2, nil)}
| RFC3339Date RFC3339Time RFC3339TimeZone {$$ = NewDateTimeTZ($1, $2, $3)}
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
| WeekdayOpt Day DateSepPlus Day {$$ = NewAmbiguousDate($1, $2, $4, nil)}

| Year {$$ = NewDMYDate(nil, nil, $1)}
| Year DateSepPlus Day {$$ = NewDMYDate(nil, $3, $1)}
| Year DateSepPlus Day DateSepPlus Day {$$ = NewDMYDate($5, $3, $1)}

  // "Feb"
| Month {$$ = NewMDYDate($1, nil, nil)}

  // "Feb 2023"
| Month Year {$$ = NewMDYDate($1, nil, $2)}

  // "Feb 3 2023"
| WeekdayOpt Month Day Year {$$ = NewWMDYDate($1, $2, $3, $4)}

  // "3 Feb 2023"
| WeekdayOpt Day Month Year {$$ = NewWDMYDate($1, $2, $3, $4)}

  // "2/3/2023", but ambiguous between North America (month-day-year) and other (day-month-year) styles.
| WeekdayOpt Day DateSepPlus Day DateSepPlus Year {$$ = NewAmbiguousDate($1, $2, $4, $6)}

  // "Feb 3"
| WeekdayOpt Month Day {$$ = NewWMDYDate($1, $2, $3, nil)}

  // "3 Feb"
| WeekdayOpt Day Month {$$ = NewWDMYDate($1, $2, $3, nil)}

  // "2023 Feb 3"
| WeekdayOpt Year Month Day {$$ = NewWDMYDate($1, $4, $3, $2)}
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
