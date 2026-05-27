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
%token BULLET
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
%token NEXT
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
%token STARTS
%token SUB
%token THROUGH
%token T
%token TH
%token THE
%token THIS
%token TILL
%token TIME
%token TO
%token UNTIL
%token WHEN
%token ENDS
%token Z

%token <string> IDENT
%token <string> INT
%token <string> MONTH_NAME
%token <string> TIME_NAME
%token <string> TIME_ZONE
%token <string> TIME_ZONE_ABBREV
%token <string> RELATIVE_DAY
%token <string> WEEKDAY_NAME
%token <string> WEEKDAY_SHORT_NAME
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
%type <string> WeekdayName
%type <string> WeekdayShortName
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
  DateTimeRange {$$ = NewRanges($1)}

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

  // "Feb 3, Mar 4" (MD)
| Month Day Month Day {$$ = NewRanges(NewRangeWithStartDate(NewRawDateFromMDY($1, $2, nil)), NewRangeWithStartDate(NewRawDateFromMDY($3, $4, nil)))}
  // "3 Feb, 4 Mar" (DM)
| Day Month Day Month {$$ = NewRanges(NewRangeWithStartDate(NewRawDateFromDMY($1, $2, nil)), NewRangeWithStartDate(NewRawDateFromDMY($3, $4, nil)))}
  // "Feb 3, Mar 4 2023" (MD)
| Month Day Month Day Year {$$ = NewRanges(NewRangeWithStartDate(NewRawDateFromMDY($1, $2, $5)), NewRangeWithStartDate(NewRawDateFromMDY($3, $4, $5)))}
  // "3 Feb, 4 Mar 2023" (DM)
| Day Month Day Month Year {$$ = NewRanges(NewRangeWithStartDate(NewRawDateFromDMY($1, $2, nil)), NewRangeWithStartDate(NewRawDateFromDMY($3, $4, $5)))}

  // "Wednesdays February 1st & 8th 12:00p-3:00p" (MD)
| WeekdayOpt Month Day DateSepOpt Day Time RangeSepOpt Time {$$ = NewRanges(NewRange(NewDateTime(NewRawDateFromWMDY($1, $2, $3, nil), $6, nil), NewDateTime(NewRawDateFromWMDY($1, $2, $3, nil), $8, nil)), NewRange(NewDateTime(NewRawDateFromWMDY($1, $2, $5, nil), $6, nil), NewDateTime(NewRawDateFromWMDY($1, $2, $5, nil), $8, nil)))}
  // "Wednesdays 1st & 8th February 12:00p-3:00p" (DM)
| WeekdayOpt Day DateSepOpt Day Month Time RangeSepOpt Time {$$ = NewRanges(NewRange(NewDateTime(NewRawDateFromWDMY($1, $2, $5, nil), $6, nil), NewDateTime(NewRawDateFromWDMY($1, $2, $5, nil), $8, nil)), NewRange(NewDateTime(NewRawDateFromWDMY($1, $4, $5, nil), $6, nil), NewDateTime(NewRawDateFromWDMY($1, $4, $5, nil), $8, nil)))}
  // "Wednesdays February 1, 8 9:00 AM - 12:00 PM ET" (MD)
| WeekdayOpt Month Day DateSepOpt Day Time RangeSepOpt Time TimeZone {$$ = NewRanges(NewRange(NewDateTime(NewRawDateFromWMDY($1, $2, $3, nil), $6, $9), NewDateTime(NewRawDateFromWMDY($1, $2, $3, nil), $8, $9)), NewRange(NewDateTime(NewRawDateFromWMDY($1, $2, $5, nil), $6, $9), NewDateTime(NewRawDateFromWMDY($1, $2, $5, nil), $8, $9)))}
  // "Wednesdays 1, 8 February 9:00 AM - 12:00 PM ET" (DM)
| WeekdayOpt Day DateSepOpt Day Month Time RangeSepOpt Time TimeZone {$$ = NewRanges(NewRange(NewDateTime(NewRawDateFromWDMY($1, $2, $5, nil), $6, $9), NewDateTime(NewRawDateFromWDMY($1, $2, $5, nil), $8, $9)), NewRange(NewDateTime(NewRawDateFromWDMY($1, $4, $5, nil), $6, $9), NewDateTime(NewRawDateFromWDMY($1, $4, $5, nil), $8, $9)))}
  // "February 3rd & 4th, 9:00 am - noon Eastern time" (MD)
| Month Day AND Day DateTimeSepPlus Time RangeSepOpt Time TimeZone {$$ = NewRanges(NewRange(NewDateTime(NewRawDateFromMDY($1, $2, nil), $6, $9), NewDateTime(NewRawDateFromMDY($1, $2, nil), $8, $9)), NewRange(NewDateTime(NewRawDateFromMDY($1, $4, nil), $6, $9), NewDateTime(NewRawDateFromMDY($1, $4, nil), $8, $9)))}
  // "3rd & 4th February, 9:00 am - noon Eastern time" (DM)
| Day AND Day Month DateTimeSepPlus Time RangeSepOpt Time TimeZone {$$ = NewRanges(NewRange(NewDateTime(NewRawDateFromDMY($1, $4, nil), $6, $9), NewDateTime(NewRawDateFromDMY($1, $4, nil), $8, $9)), NewRange(NewDateTime(NewRawDateFromDMY($3, $4, nil), $6, $9), NewDateTime(NewRawDateFromDMY($3, $4, nil), $8, $9)))}
  // "February 1, 8, and 15 9:00am - 12:00pm (ET)" — multi-day with time range and timezone (MD)
| Month DayPlus1 Time RangeSepOpt Time TimeZone {$$ = NewRangesFromDatesTimeRange(NewRawDateFromMDsYs($1, $2, nil), $3, $5, $6)}
  // "1, 8, and 15 February 9:00am - 12:00pm (ET)" (DM)
| DayPlus1 Month Time RangeSepOpt Time TimeZone {$$ = NewRangesFromDatesTimeRange(NewRawDateFromDsMYs($1, $2, nil), $3, $5, $6)}
  // "February 1, 8, 10 2023 7:00am - 11:00am (SAST)" — multi-day with year and time range (MD)
| Month DayPlus1 Year Time RangeSepOpt Time TimeZoneOpt {$$ = NewRangesFromDatesTimeRange(NewRawDateFromMDsYs($1, $2, $3), $4, $6, $7)}
  // "1, 8, 10 February 2023 7:00am - 11:00am (SAST)" (DM)
| DayPlus1 Month Year Time RangeSepOpt Time TimeZoneOpt {$$ = NewRangesFromDatesTimeRange(NewRawDateFromDsMYs($1, $2, $3), $4, $6, $7)}
  // "1, 3 Feb and 2 Mar 2023 5pm - 9pm (SAST)" — multi-month day list with year and time range (DM)
| DayPlus Month AND DayPlus Month Year Time RangeSepOpt Time TimeZoneOpt {$$ = NewRangesFromDatesTimeRange(append(NewRawDateFromDsMYs($1, $2, $6), NewRawDateFromDsMYs($4, $5, $6)...), $7, $9, $10)}
  // "Feb 1, 3 and Mar 2 2023 5pm - 9pm (SAST)" (MD)
| Month DayPlus AND Month DayPlus Year Time RangeSepOpt Time TimeZoneOpt {$$ = NewRangesFromDatesTimeRange(append(NewRawDateFromMDsYs($1, $2, $6), NewRawDateFromMDsYs($4, $5, $6)...), $7, $9, $10)}
  // "Feb. 3, 2023 12:00pm, 3:00pm" — comma-separated times sharing a date
  // Use COLON (not TimeSep which includes PERIOD) to avoid matching "Thu, 02.03.2023"
| DateTime COMMA INT COLON INT Am {$$ = NewRanges(NewRangeWithStart($1), NewRangeWithStart(NewDateTime($1.Date, NewAMTime($3, $5, nil, nil), $1.TimeZone)))}
| DateTime COMMA INT COLON INT Pm {$$ = NewRanges(NewRangeWithStart($1), NewRangeWithStart(NewDateTime($1.Date, NewPMTime($3, $5, nil, nil), $1.TimeZone)))}
| DateTime COMMA INT Am {$$ = NewRanges(NewRangeWithStart($1), NewRangeWithStart(NewDateTime($1.Date, NewAMTime($3, nil, nil, nil), $1.TimeZone)))}
| DateTime COMMA INT Pm {$$ = NewRanges(NewRangeWithStart($1), NewRangeWithStart(NewDateTime($1.Date, NewPMTime($3, nil, nil, nil), $1.TimeZone)))}
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
| DayPlus1 AND Day {$$ = append($1, $3)}
;


//
// Date Time Range
//

DateTimeRange:
  DateTime {$$ = NewRangeWithStart($1)}

| DatePrefixPlus DateTimeRange {$$ = $2}
| STARTS ColonOpt DateTime ENDS ColonOpt DateTime {$$ = NewRange($3, $6)}
| STARTS ColonOpt DateTime {$$ = NewRangeWithStart($3)}

// TODO: we should handle semantics of this weekday, but not clear how.
| DateTimeRange DateTimeSepOpt ON Weekday {$$ = $1}

| RangePrefixPlus DateTimeRange {$$ = $2}

  // Time-only DateTimeRange: for inputs reduced to bare times after preprocessing
  // (e.g. "8PM" after stripping "Doors:", or "6:30pm ET" after stripping "Wednesdays at").
  // These rules live at DateTimeRange level (not DateTime) to avoid GLR ambiguity
  // with Time→DateTime paths that need a following Date.
| INT Am TimeZoneOpt {$$ = NewRangeWithStart(NewDateTime(nil, NewAMTime($1, nil, nil, nil), $3))}
| INT Pm TimeZoneOpt {$$ = NewRangeWithStart(NewDateTime(nil, NewPMTime($1, nil, nil, nil), $3))}
| INT TimeSep INT Am TimeZoneOpt {$$ = NewRangeWithStart(NewDateTime(nil, NewAMTime($1, $3, nil, nil), $5))}
| INT TimeSep INT Pm TimeZoneOpt {$$ = NewRangeWithStart(NewDateTime(nil, NewPMTime($1, $3, nil, nil), $5))}
| TIME_NAME TimeZoneOpt {$$ = NewRangeWithStart(NewDateTime(nil, NewTime($1, nil, nil, nil), $2))}

  // Time-only ranges: "5pm to 2am", "6:30pm - 9:30pm ET"
  // Note: INT TimeSep INT Am/Pm variants are omitted — they create GLR conflicts
  // with Time→DateTime reductions. Inputs like "6:30pm - 9:30pm" need preprocessing
  // to reach these rules (e.g. by stripping weekday prefix that leaves only the time).
| INT Am TimeZoneOpt RangeSepPlus Time TimeZoneOpt {$$ = NewRange(NewDateTime(nil, NewAMTime($1, nil, nil, nil), $3), NewDateTime(nil, $5, $6))}
| INT Pm TimeZoneOpt RangeSepPlus Time TimeZoneOpt {$$ = NewRange(NewDateTime(nil, NewPMTime($1, nil, nil, nil), $3), NewDateTime(nil, $5, $6))}
| TIME_NAME TimeZoneOpt RangeSepPlus Time TimeZoneOpt {$$ = NewRange(NewDateTime(nil, NewTime($1, nil, nil, nil), $2), NewDateTime(nil, $4, $5))}

  // "May 28, Wednesday • 16:15 – 16:45"
| Month Day DateTimeSepPlus Weekday DateTimeSepPlus Time RangeSepOpt Time TimeZoneOpt {$$ = NewRange(NewDateTime(NewRawDateFromWMDY($4, $1, $2, nil), $6, $9), NewDateTime(NewRawDateFromWMDY($4, $1, $2, nil), $8, $9))}
| Day Month DateTimeSepPlus Weekday DateTimeSepPlus Time RangeSepOpt Time TimeZoneOpt {$$ = NewRange(NewDateTime(NewRawDateFromWDMY($4, $1, $2, nil), $6, $9), NewDateTime(NewRawDateFromWDMY($4, $1, $2, nil), $8, $9))}
| Month Day Year DateTimeSepPlus Weekday DateTimeSepPlus Time RangeSepOpt Time TimeZoneOpt {$$ = NewRange(NewDateTime(NewRawDateFromWMDY($5, $1, $2, $3), $7, $10), NewDateTime(NewRawDateFromWMDY($5, $1, $2, $3), $9, $10))}
| Day Month Year DateTimeSepPlus Weekday DateTimeSepPlus Time RangeSepOpt Time TimeZoneOpt {$$ = NewRange(NewDateTime(NewRawDateFromWDMY($5, $1, $2, $3), $7, $10), NewDateTime(NewRawDateFromWDMY($5, $1, $2, $3), $9, $10))}

| DateTime RangeSepPlus Time {$$ = NewRange($1, NewDateTime($1.Date, $3, $1.TimeZone))}
| DateTime RangeSepPlus Time TimeZone {$$ = NewRange(NewDateTime($1.Date, $1.Time, $4), NewDateTime($1.Date, $3, $4))}
| Time RangeSepPlus DateTime {$$ = NewRange(NewDateTime($3.Date, $1, $3.TimeZone), $3)}

  // "Feb 3, 2023 - Feb 4, 2023"
| DateTime RangeSepPlus DateTime {$$ = NewRange($1, $3)}
  // "February 3rd - February 4th, 2023 ET"
| DateTime RangeSepPlus DateTime TimeZone {$$ = NewRange(NewDateTime($1.Date, $1.Time, $4), NewDateTime($3.Date, $3.Time, $4))}

  // "Thu Feb 3 - Sat Mar 4, 2023"
| Date RangeSepPlus Day {$$ = NewRangeWithStartEndDates($1, NewRawDateFromDMY($3, nil, nil))}

  // "Thu Feb 3 - Sat Mar 4, 2023"
| Day RangeSepPlus Date {$$ = NewRangeWithStartEndDates(NewRawDateFromDMY($1, nil, nil), $3)}

  // "Feb 3-4, 2023"
| Date RangeSepPlus Day Year {$$ = NewRangeWithStartEndDates($1, NewRawDateFromMDY(nil, $3, $4))}

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
  Date {$$ = NewDateTimeWithDate($1)}

| Date DateTimeSepPlus TimeZoneOpt {$$ = NewDateTime($1, nil, $3)}
| Date Time TimeZoneOpt {$$ = NewDateTime($1, $2, $3)}
| Date DateTimeSepPlus Time TimeZoneOpt {$$ = NewDateTime($1, $3, $4)}
  // Explicit Date + INT Am/Pm: needed because after Date, yacc reduces INT to Day
  // (state 76), blocking Time: INT Am/Pm paths. These inline the Time production
  // to create the right state machine transitions.
| Date INT Am TimeZoneOpt {$$ = NewDateTime($1, NewAMTime($2, nil, nil, nil), $4)}
| Date INT Pm TimeZoneOpt {$$ = NewDateTime($1, NewPMTime($2, nil, nil, nil), $4)}
| Date INT TimeSep INT Am TimeZoneOpt {$$ = NewDateTime($1, NewAMTime($2, $4, nil, nil), $6)}
| Date INT TimeSep INT Pm TimeZoneOpt {$$ = NewDateTime($1, NewPMTime($2, $4, nil, nil), $6)}
| Time TimeZoneOpt Date {$$ = NewDateTime($3, $1, $2)}
| Time TimeZoneOpt DateTimeSepPlus Date {$$ = NewDateTime($4, $1, $2)}
  // Time-only DateTime rules are at DateTimeRange level (not here) to avoid GLR
  // ambiguity with Time→DateTime paths. See DateTimeRange section below.
// 9:00am 3rd Feb - 4th Feb 3:00pm 2023
| Date Time Year TimeZoneOpt {$$ = NewDateTime(NewRawDateFromDMY($1.Day, $1.Month, $3), $2, $4)}
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
| BULLET
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
  RFC3339Date {$$ = NewDateTimeWithDate($1)}
| RFC3339Date RFC3339Time {$$ = NewDateTime($1, $2, nil)}
| RFC3339Date RFC3339Time RFC3339TimeZone {$$ = NewDateTime($1, $2, $3)}
;

RFC3339Date:
  Year SUB INT SUB INT {$$ = NewRawDateFromDMY($5, $3, $1)}
;

RFC3339Time:
  T INT COLON INT COLON INT {$$ = NewTime($2, $4, $6, nil)}
  // "T12:00" (HH:MM without seconds)
| T INT COLON INT {$$ = NewTime($2, $4, nil, nil)}
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

| WeekdayName {$$ = NewRawDateFromRelative($1)}

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
  ADD
| AND
| COMMA
| DEC
| PERIOD
| QUO
;


WeekdayOpt:
  {$$ = ""}
| Weekday WeekdaySepOpt
;
WeekdaySepOpt:

| WeekdaySep
;
WeekdaySep:
  COMMA
| SUB
;
Weekday:
  WeekdayPrefix Weekday {$$ = $2}
| WeekdayName
| WeekdayShortName
;
WeekdayPrefix:
  NEXT
| THIS
;
WeekdayName:
  WEEKDAY_NAME
;
WeekdayShortName:
  TH {$$ = "TH"}
| WEEKDAY_SHORT_NAME
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
