package datetime

import "github.com/findyourpaths/phil/glr"
import "cloud.google.com/go/civil"

/*
Rules

root:
  DateTimeTZRanges
root:
  DateTimeTZRanges RootSuffixPlus
RootSuffixPlus:
  RootSuffix
RootSuffixPlus:
  RootSuffixPlus RootSuffix
RootSuffix:
  GOOGLE
RootSuffix:
  CALENDAR
RootSuffix:
  ICS
DateTimeTZRanges:
  RangesPrefixPlus DateTimeTZRanges
DateTimeTZRanges:
  DateTimeTZRange
DateTimeTZRanges:
  DateTimeTZRanges AndOpt DateTimeTZRange
DateTimeTZRanges:
  Month DayPlus1
DateTimeTZRanges:
  DayPlus1 Month
DateTimeTZRanges:
  Month DayPlus1 YEAR
DateTimeTZRanges:
  Month DayPlus AND Month DayPlus YEAR
DateTimeTZRanges:
  DayPlus1 Month YEAR
DateTimeTZRanges:
  DayPlus Month AND DayPlus Month YEAR
DateTimeTZRanges:
  Month Day RangeSep Day Day RangeSep Day
DateTimeTZRanges:
  Day RangeSep Day Day RangeSep Day Month
DateTimeTZRanges:
  Month Day RangeSep Day Month Day RangeSep Day
DateTimeTZRanges:
  Day RangeSep Day Month Day RangeSep Day Month
DateTimeTZRanges:
  Month Day RangeSep Day Day RangeSep Day YEAR
DateTimeTZRanges:
  Day RangeSep Day Day RangeSep Day Month YEAR
DateTimeTZRanges:
  Month Day RangeSep Day Month Day RangeSep Day YEAR
DateTimeTZRanges:
  Day RangeSep Day Month Day RangeSep Day Month YEAR
DateTimeTZRanges:
  Month Day Month Day YEAR
RangesPrefixPlus:
  RangesPrefix
RangesPrefixPlus:
  RangesPrefixPlus RangesPrefix
RangesPrefix:
  WHEN
DayPlus:
  Day
DayPlus:
  DayPlus1
DayPlus1:
  Day Day
DayPlus1:
  DayPlus1 Day
DateTimeTZRange:
  RangePrefixPlus DateTimeTZRange
DateTimeTZRange:
  DateTimeTZ
DateTimeTZRange:
  DateTimeTZ RangeSepPlus Time
DateTimeTZRange:
  DateTimeTZ RangeSepPlus DateTimeTZ
DateTimeTZRange:
  Month Day RangeSepPlus Day
DateTimeTZRange:
  Day RangeSepPlus Day Month
DateTimeTZRange:
  Month Day RangeSepPlus Day YEAR
DateTimeTZRange:
  Day RangeSepPlus Day Month YEAR
DateTimeTZRange:
  Month Day RangeSepPlus Month Day YEAR
DateTimeTZRange:
  WeekDay Month Day RangeSepPlus WeekDay Month Day YEAR
DateTimeTZRange:
  WeekDay Month Day RangeSepPlus WeekDay Day Month YEAR
DateTimeTZRange:
  WeekDay Month Day RangeSepPlus Day WeekDay Month YEAR
DateTimeTZRange:
  WeekDay Day Month RangeSepPlus WeekDay Month Day YEAR
DateTimeTZRange:
  WeekDay Day Month RangeSepPlus WeekDay Day Month YEAR
DateTimeTZRange:
  WeekDay Day Month RangeSepPlus Day WeekDay Month YEAR
DateTimeTZRange:
  Time DateTimeSepOpt Day Month RangeSepPlus Day Month DateTimeSepOpt Time YEAR
DateTimeTZRange:
  Date Time Time
RangePrefixPlus:
  RangePrefix
RangePrefixPlus:
  RangePrefixPlus RangePrefix
RangePrefix:
  BEGINNING
RangePrefix:
  FROM
RangeSepPlus:
  RangeSep
RangeSepPlus:
  RangeSepPlus RangeSep
RangeSep:
  DEC
RangeSep:
  SUB
RangeSep:
  THROUGH
RangeSep:
  TILL
RangeSep:
  TO
RangeSep:
  UNTIL
DateTimeTZ:
  Date
DateTimeTZ:
  Date Time
DateTimeTZ:
  Date DateTimeSepPlus Time
DateTimeTZ:
  Time Date
DateTimeTZ:
  Time DateTimeSepPlus Date
DateTimeSepOpt:
  <empty>
DateTimeSepOpt:
  DateTimeSepPlus
DateTimeSepPlus:
  DateTimeSep
DateTimeSepPlus:
  DateTimeSepPlus DateTimeSep
DateTimeSep:
  DEC
DateTimeSep:
  SUB
Date:
  DatePrefixPlus Date
Date:
  Date DateSuffixPlus
Date:
  Day DateSepPlus Day
Date:
  Year
Date:
  Year DateSepPlus Day
Date:
  Year DateSepPlus Day DateSepPlus Day
Date:
  Month
Date:
  Month Year
Date:
  Month Day Year
Date:
  Day Month Year
Date:
  Day DateSepPlus Day DateSepPlus Year
Date:
  Month Day
Date:
  Day Month
Date:
  Year Month Day
DatePrefixPlus:
  DatePrefix
DatePrefixPlus:
  DatePrefixPlus DatePrefix
DatePrefix:
  DATE
DatePrefix:
  WeekDay
DatePrefix:
  COLON
DatePrefix:
  COMMA
DatePrefix:
  TIME
DateSepPlus:
  DateSep
DateSepPlus:
  DateSepPlus DateSep
DateSep:
  COMMA
DateSep:
  DEC
DateSep:
  PERIOD
DateSep:
  SUB
DateSep:
  QUO
DateSuffixPlus:
  DateSuffix
DateSuffixPlus:
  DateSuffixPlus DateSuffix
DateSuffix:
  T
Day:
  INT
Day:
  INT DaySuffixPlus
DaySuffixPlus:
  DaySuffix
DaySuffixPlus:
  DaySuffixPlus DaySuffix
DaySuffix:
  COMMA
DaySuffix:
  ORD_IND
DaySuffix:
  PERIOD
DaySuffix:
  TH
Month:
  MONTH_NAME
Month:
  MONTH_NAME MonthSuffixPlus
MonthSuffixPlus:
  MonthSuffix
MonthSuffixPlus:
  MonthSuffixPlus MonthSuffix
MonthSuffix:
  COMMA
Year:
  YEAR
Year:
  YEAR YearSuffixPlus
YearSuffixPlus:
  YearSuffix
YearSuffixPlus:
  YearSuffixPlus YearSuffix
YearSuffix:
  COMMA
WeekDay:
  TH
WeekDay:
  WEEKDAY_NAME
Time:
  TimePrefixPlus Time
Time:
  INT AM
Time:
  INT PM
Time:
  INT TimeSep INT
Time:
  INT TimeSep INT TimeSep INT
Time:
  INT TimeSep INT AM
Time:
  INT TimeSep INT PM
TimePrefixPlus:
  TimePrefix
TimePrefixPlus:
  TimePrefixPlus TimePrefix
TimePrefix:
  AT
TimePrefix:
  COLON
TimePrefix:
  FROM
TimePrefix:
  ON
TimePrefix:
  TIME
TimeSep:
  COLON
TimeSep:
  PERIOD
AndOpt:
  <empty>
AndOpt:
  AND
OfOpt:
  <empty>
OfOpt:
  OF
*/

var parseRules = &glr.Rules{Items:[]glr.Rule{
  /*   0 */ glr.Rule{Nonterminal:"", RHS:[]string(nil), Type:""}, // ignored because rule-numbering starts at 1
  /*   1 */ glr.Rule{Nonterminal:"root", RHS:[]string{"DateTimeTZRanges"}, Type:"*DateTimeTZRanges"},
  /*   2 */ glr.Rule{Nonterminal:"root", RHS:[]string{"DateTimeTZRanges", "RootSuffixPlus"}, Type:"*DateTimeTZRanges"},
  /*   3 */ glr.Rule{Nonterminal:"RootSuffixPlus", RHS:[]string{"RootSuffix"}, Type:""},
  /*   4 */ glr.Rule{Nonterminal:"RootSuffixPlus", RHS:[]string{"RootSuffixPlus", "RootSuffix"}, Type:""},
  /*   5 */ glr.Rule{Nonterminal:"RootSuffix", RHS:[]string{"GOOGLE"}, Type:""},
  /*   6 */ glr.Rule{Nonterminal:"RootSuffix", RHS:[]string{"CALENDAR"}, Type:""},
  /*   7 */ glr.Rule{Nonterminal:"RootSuffix", RHS:[]string{"ICS"}, Type:""},
  /*   8 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"RangesPrefixPlus", "DateTimeTZRanges"}, Type:"*DateTimeTZRanges"},
  /*   9 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"DateTimeTZRange"}, Type:"*DateTimeTZRanges"},
  /*  10 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"DateTimeTZRanges", "AndOpt", "DateTimeTZRange"}, Type:"*DateTimeTZRanges"},
  /*  11 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"Month", "DayPlus1"}, Type:"*DateTimeTZRanges"},
  /*  12 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"DayPlus1", "Month"}, Type:"*DateTimeTZRanges"},
  /*  13 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"Month", "DayPlus1", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  14 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"Month", "DayPlus", "AND", "Month", "DayPlus", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  15 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"DayPlus1", "Month", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  16 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"DayPlus", "Month", "AND", "DayPlus", "Month", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  17 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"Month", "Day", "RangeSep", "Day", "Day", "RangeSep", "Day"}, Type:"*DateTimeTZRanges"},
  /*  18 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"Day", "RangeSep", "Day", "Day", "RangeSep", "Day", "Month"}, Type:"*DateTimeTZRanges"},
  /*  19 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"Month", "Day", "RangeSep", "Day", "Month", "Day", "RangeSep", "Day"}, Type:"*DateTimeTZRanges"},
  /*  20 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"Day", "RangeSep", "Day", "Month", "Day", "RangeSep", "Day", "Month"}, Type:"*DateTimeTZRanges"},
  /*  21 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"Month", "Day", "RangeSep", "Day", "Day", "RangeSep", "Day", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  22 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"Day", "RangeSep", "Day", "Day", "RangeSep", "Day", "Month", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  23 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"Month", "Day", "RangeSep", "Day", "Month", "Day", "RangeSep", "Day", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  24 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"Day", "RangeSep", "Day", "Month", "Day", "RangeSep", "Day", "Month", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  25 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"Month", "Day", "Month", "Day", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  26 */ glr.Rule{Nonterminal:"RangesPrefixPlus", RHS:[]string{"RangesPrefix"}, Type:""},
  /*  27 */ glr.Rule{Nonterminal:"RangesPrefixPlus", RHS:[]string{"RangesPrefixPlus", "RangesPrefix"}, Type:""},
  /*  28 */ glr.Rule{Nonterminal:"RangesPrefix", RHS:[]string{"WHEN"}, Type:""},
  /*  29 */ glr.Rule{Nonterminal:"DayPlus", RHS:[]string{"Day"}, Type:"[]string"},
  /*  30 */ glr.Rule{Nonterminal:"DayPlus", RHS:[]string{"DayPlus1"}, Type:"[]string"},
  /*  31 */ glr.Rule{Nonterminal:"DayPlus1", RHS:[]string{"Day", "Day"}, Type:"[]string"},
  /*  32 */ glr.Rule{Nonterminal:"DayPlus1", RHS:[]string{"DayPlus1", "Day"}, Type:"[]string"},
  /*  33 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"RangePrefixPlus", "DateTimeTZRange"}, Type:"*DateTimeTZRange"},
  /*  34 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"DateTimeTZ"}, Type:"*DateTimeTZRange"},
  /*  35 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"DateTimeTZ", "RangeSepPlus", "Time"}, Type:"*DateTimeTZRange"},
  /*  36 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"DateTimeTZ", "RangeSepPlus", "DateTimeTZ"}, Type:"*DateTimeTZRange"},
  /*  37 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"Month", "Day", "RangeSepPlus", "Day"}, Type:"*DateTimeTZRange"},
  /*  38 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"Day", "RangeSepPlus", "Day", "Month"}, Type:"*DateTimeTZRange"},
  /*  39 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"Month", "Day", "RangeSepPlus", "Day", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  40 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"Day", "RangeSepPlus", "Day", "Month", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  41 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"Month", "Day", "RangeSepPlus", "Month", "Day", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  42 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"WeekDay", "Month", "Day", "RangeSepPlus", "WeekDay", "Month", "Day", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  43 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"WeekDay", "Month", "Day", "RangeSepPlus", "WeekDay", "Day", "Month", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  44 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"WeekDay", "Month", "Day", "RangeSepPlus", "Day", "WeekDay", "Month", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  45 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"WeekDay", "Day", "Month", "RangeSepPlus", "WeekDay", "Month", "Day", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  46 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"WeekDay", "Day", "Month", "RangeSepPlus", "WeekDay", "Day", "Month", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  47 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"WeekDay", "Day", "Month", "RangeSepPlus", "Day", "WeekDay", "Month", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  48 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"Time", "DateTimeSepOpt", "Day", "Month", "RangeSepPlus", "Day", "Month", "DateTimeSepOpt", "Time", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  49 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"Date", "Time", "Time"}, Type:"*DateTimeTZRange"},
  /*  50 */ glr.Rule{Nonterminal:"RangePrefixPlus", RHS:[]string{"RangePrefix"}, Type:""},
  /*  51 */ glr.Rule{Nonterminal:"RangePrefixPlus", RHS:[]string{"RangePrefixPlus", "RangePrefix"}, Type:""},
  /*  52 */ glr.Rule{Nonterminal:"RangePrefix", RHS:[]string{"BEGINNING"}, Type:""},
  /*  53 */ glr.Rule{Nonterminal:"RangePrefix", RHS:[]string{"FROM"}, Type:""},
  /*  54 */ glr.Rule{Nonterminal:"RangeSepPlus", RHS:[]string{"RangeSep"}, Type:""},
  /*  55 */ glr.Rule{Nonterminal:"RangeSepPlus", RHS:[]string{"RangeSepPlus", "RangeSep"}, Type:""},
  /*  56 */ glr.Rule{Nonterminal:"RangeSep", RHS:[]string{"DEC"}, Type:""},
  /*  57 */ glr.Rule{Nonterminal:"RangeSep", RHS:[]string{"SUB"}, Type:""},
  /*  58 */ glr.Rule{Nonterminal:"RangeSep", RHS:[]string{"THROUGH"}, Type:""},
  /*  59 */ glr.Rule{Nonterminal:"RangeSep", RHS:[]string{"TILL"}, Type:""},
  /*  60 */ glr.Rule{Nonterminal:"RangeSep", RHS:[]string{"TO"}, Type:""},
  /*  61 */ glr.Rule{Nonterminal:"RangeSep", RHS:[]string{"UNTIL"}, Type:""},
  /*  62 */ glr.Rule{Nonterminal:"DateTimeTZ", RHS:[]string{"Date"}, Type:"*DateTimeTZ"},
  /*  63 */ glr.Rule{Nonterminal:"DateTimeTZ", RHS:[]string{"Date", "Time"}, Type:"*DateTimeTZ"},
  /*  64 */ glr.Rule{Nonterminal:"DateTimeTZ", RHS:[]string{"Date", "DateTimeSepPlus", "Time"}, Type:"*DateTimeTZ"},
  /*  65 */ glr.Rule{Nonterminal:"DateTimeTZ", RHS:[]string{"Time", "Date"}, Type:"*DateTimeTZ"},
  /*  66 */ glr.Rule{Nonterminal:"DateTimeTZ", RHS:[]string{"Time", "DateTimeSepPlus", "Date"}, Type:"*DateTimeTZ"},
  /*  67 */ glr.Rule{Nonterminal:"DateTimeSepOpt", RHS:[]string(nil), Type:""},
  /*  68 */ glr.Rule{Nonterminal:"DateTimeSepOpt", RHS:[]string{"DateTimeSepPlus"}, Type:""},
  /*  69 */ glr.Rule{Nonterminal:"DateTimeSepPlus", RHS:[]string{"DateTimeSep"}, Type:""},
  /*  70 */ glr.Rule{Nonterminal:"DateTimeSepPlus", RHS:[]string{"DateTimeSepPlus", "DateTimeSep"}, Type:""},
  /*  71 */ glr.Rule{Nonterminal:"DateTimeSep", RHS:[]string{"DEC"}, Type:""},
  /*  72 */ glr.Rule{Nonterminal:"DateTimeSep", RHS:[]string{"SUB"}, Type:""},
  /*  73 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"DatePrefixPlus", "Date"}, Type:"civil.Date"},
  /*  74 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Date", "DateSuffixPlus"}, Type:"civil.Date"},
  /*  75 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Day", "DateSepPlus", "Day"}, Type:"civil.Date"},
  /*  76 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Year"}, Type:"civil.Date"},
  /*  77 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Year", "DateSepPlus", "Day"}, Type:"civil.Date"},
  /*  78 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Year", "DateSepPlus", "Day", "DateSepPlus", "Day"}, Type:"civil.Date"},
  /*  79 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Month"}, Type:"civil.Date"},
  /*  80 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Month", "Year"}, Type:"civil.Date"},
  /*  81 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Month", "Day", "Year"}, Type:"civil.Date"},
  /*  82 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Day", "Month", "Year"}, Type:"civil.Date"},
  /*  83 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Day", "DateSepPlus", "Day", "DateSepPlus", "Year"}, Type:"civil.Date"},
  /*  84 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Month", "Day"}, Type:"civil.Date"},
  /*  85 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Day", "Month"}, Type:"civil.Date"},
  /*  86 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Year", "Month", "Day"}, Type:"civil.Date"},
  /*  87 */ glr.Rule{Nonterminal:"DatePrefixPlus", RHS:[]string{"DatePrefix"}, Type:""},
  /*  88 */ glr.Rule{Nonterminal:"DatePrefixPlus", RHS:[]string{"DatePrefixPlus", "DatePrefix"}, Type:""},
  /*  89 */ glr.Rule{Nonterminal:"DatePrefix", RHS:[]string{"DATE"}, Type:""},
  /*  90 */ glr.Rule{Nonterminal:"DatePrefix", RHS:[]string{"WeekDay"}, Type:""},
  /*  91 */ glr.Rule{Nonterminal:"DatePrefix", RHS:[]string{"COLON"}, Type:""},
  /*  92 */ glr.Rule{Nonterminal:"DatePrefix", RHS:[]string{"COMMA"}, Type:""},
  /*  93 */ glr.Rule{Nonterminal:"DatePrefix", RHS:[]string{"TIME"}, Type:""},
  /*  94 */ glr.Rule{Nonterminal:"DateSepPlus", RHS:[]string{"DateSep"}, Type:""},
  /*  95 */ glr.Rule{Nonterminal:"DateSepPlus", RHS:[]string{"DateSepPlus", "DateSep"}, Type:""},
  /*  96 */ glr.Rule{Nonterminal:"DateSep", RHS:[]string{"COMMA"}, Type:""},
  /*  97 */ glr.Rule{Nonterminal:"DateSep", RHS:[]string{"DEC"}, Type:""},
  /*  98 */ glr.Rule{Nonterminal:"DateSep", RHS:[]string{"PERIOD"}, Type:""},
  /*  99 */ glr.Rule{Nonterminal:"DateSep", RHS:[]string{"SUB"}, Type:""},
  /* 100 */ glr.Rule{Nonterminal:"DateSep", RHS:[]string{"QUO"}, Type:""},
  /* 101 */ glr.Rule{Nonterminal:"DateSuffixPlus", RHS:[]string{"DateSuffix"}, Type:""},
  /* 102 */ glr.Rule{Nonterminal:"DateSuffixPlus", RHS:[]string{"DateSuffixPlus", "DateSuffix"}, Type:""},
  /* 103 */ glr.Rule{Nonterminal:"DateSuffix", RHS:[]string{"T"}, Type:""},
  /* 104 */ glr.Rule{Nonterminal:"Day", RHS:[]string{"INT"}, Type:"string"},
  /* 105 */ glr.Rule{Nonterminal:"Day", RHS:[]string{"INT", "DaySuffixPlus"}, Type:"string"},
  /* 106 */ glr.Rule{Nonterminal:"DaySuffixPlus", RHS:[]string{"DaySuffix"}, Type:""},
  /* 107 */ glr.Rule{Nonterminal:"DaySuffixPlus", RHS:[]string{"DaySuffixPlus", "DaySuffix"}, Type:""},
  /* 108 */ glr.Rule{Nonterminal:"DaySuffix", RHS:[]string{"COMMA"}, Type:""},
  /* 109 */ glr.Rule{Nonterminal:"DaySuffix", RHS:[]string{"ORD_IND"}, Type:""},
  /* 110 */ glr.Rule{Nonterminal:"DaySuffix", RHS:[]string{"PERIOD"}, Type:""},
  /* 111 */ glr.Rule{Nonterminal:"DaySuffix", RHS:[]string{"TH"}, Type:""},
  /* 112 */ glr.Rule{Nonterminal:"Month", RHS:[]string{"MONTH_NAME"}, Type:"string"},
  /* 113 */ glr.Rule{Nonterminal:"Month", RHS:[]string{"MONTH_NAME", "MonthSuffixPlus"}, Type:"string"},
  /* 114 */ glr.Rule{Nonterminal:"MonthSuffixPlus", RHS:[]string{"MonthSuffix"}, Type:""},
  /* 115 */ glr.Rule{Nonterminal:"MonthSuffixPlus", RHS:[]string{"MonthSuffixPlus", "MonthSuffix"}, Type:""},
  /* 116 */ glr.Rule{Nonterminal:"MonthSuffix", RHS:[]string{"COMMA"}, Type:""},
  /* 117 */ glr.Rule{Nonterminal:"Year", RHS:[]string{"YEAR"}, Type:"string"},
  /* 118 */ glr.Rule{Nonterminal:"Year", RHS:[]string{"YEAR", "YearSuffixPlus"}, Type:"string"},
  /* 119 */ glr.Rule{Nonterminal:"YearSuffixPlus", RHS:[]string{"YearSuffix"}, Type:""},
  /* 120 */ glr.Rule{Nonterminal:"YearSuffixPlus", RHS:[]string{"YearSuffixPlus", "YearSuffix"}, Type:""},
  /* 121 */ glr.Rule{Nonterminal:"YearSuffix", RHS:[]string{"COMMA"}, Type:""},
  /* 122 */ glr.Rule{Nonterminal:"WeekDay", RHS:[]string{"TH"}, Type:""},
  /* 123 */ glr.Rule{Nonterminal:"WeekDay", RHS:[]string{"WEEKDAY_NAME"}, Type:""},
  /* 124 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"TimePrefixPlus", "Time"}, Type:"civil.Time"},
  /* 125 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "AM"}, Type:"civil.Time"},
  /* 126 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "PM"}, Type:"civil.Time"},
  /* 127 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "TimeSep", "INT"}, Type:"civil.Time"},
  /* 128 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "TimeSep", "INT", "TimeSep", "INT"}, Type:"civil.Time"},
  /* 129 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "TimeSep", "INT", "AM"}, Type:"civil.Time"},
  /* 130 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "TimeSep", "INT", "PM"}, Type:"civil.Time"},
  /* 131 */ glr.Rule{Nonterminal:"TimePrefixPlus", RHS:[]string{"TimePrefix"}, Type:""},
  /* 132 */ glr.Rule{Nonterminal:"TimePrefixPlus", RHS:[]string{"TimePrefixPlus", "TimePrefix"}, Type:""},
  /* 133 */ glr.Rule{Nonterminal:"TimePrefix", RHS:[]string{"AT"}, Type:""},
  /* 134 */ glr.Rule{Nonterminal:"TimePrefix", RHS:[]string{"COLON"}, Type:""},
  /* 135 */ glr.Rule{Nonterminal:"TimePrefix", RHS:[]string{"FROM"}, Type:""},
  /* 136 */ glr.Rule{Nonterminal:"TimePrefix", RHS:[]string{"ON"}, Type:""},
  /* 137 */ glr.Rule{Nonterminal:"TimePrefix", RHS:[]string{"TIME"}, Type:""},
  /* 138 */ glr.Rule{Nonterminal:"TimeSep", RHS:[]string{"COLON"}, Type:""},
  /* 139 */ glr.Rule{Nonterminal:"TimeSep", RHS:[]string{"PERIOD"}, Type:""},
  /* 140 */ glr.Rule{Nonterminal:"AndOpt", RHS:[]string(nil), Type:""},
  /* 141 */ glr.Rule{Nonterminal:"AndOpt", RHS:[]string{"AND"}, Type:""},
  /* 142 */ glr.Rule{Nonterminal:"OfOpt", RHS:[]string(nil), Type:""},
  /* 143 */ glr.Rule{Nonterminal:"OfOpt", RHS:[]string{"OF"}, Type:""},
}}

// Semantic action functions

var parseActions = &glr.SemanticActions{Items:[]any{
  /*   0 */ nil, // empty action
  /*   1 */ func (DateTimeTZRanges1 *DateTimeTZRanges) *DateTimeTZRanges {return DateTimeTZRanges1},
  /*   2 */ func (DateTimeTZRanges1 *DateTimeTZRanges, RootSuffixPlus1 string) *DateTimeTZRanges {return DateTimeTZRanges1},
  /*   3 */ func (RootSuffix1 string) string {return RootSuffix1},
  /*   4 */ func (RootSuffixPlus1 string, RootSuffix1 string) string {return RootSuffixPlus1},
  /*   5 */ func (GOOGLE1 string) string {return GOOGLE1},
  /*   6 */ func (CALENDAR1 string) string {return CALENDAR1},
  /*   7 */ func (ICS1 string) string {return ICS1},
  /*   8 */ func (RangesPrefixPlus1 string, DateTimeTZRanges1 *DateTimeTZRanges) *DateTimeTZRanges {return DateTimeTZRanges1},
  /*   9 */ func (DateTimeTZRange1 *DateTimeTZRange) *DateTimeTZRanges {return &DateTimeTZRanges{Items: []*DateTimeTZRange{DateTimeTZRange1}}},
  /*  10 */ func (DateTimeTZRanges1 *DateTimeTZRanges, AndOpt1 string, DateTimeTZRange1 *DateTimeTZRange) *DateTimeTZRanges {return AppendDateTimeTZRanges(DateTimeTZRanges1, DateTimeTZRange1)},
  /*  11 */ func (Month1 string, DayPlus11 []string) *DateTimeTZRanges {return NewRangesWithStartDates(NewMDsYDates(Month1, DayPlus11, "")...)},
  /*  12 */ func (DayPlus11 []string, Month1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewDsMYDates(DayPlus11, Month1, "")...)},
  /*  13 */ func (Month1 string, DayPlus11 []string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewMDsYDates(Month1, DayPlus11, YEAR1)...)},
  /*  14 */ func (Month1 string, DayPlus1 []string, AND1 string, Month2 string, DayPlus2 []string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStartDates(append(NewMDsYDates(Month1, DayPlus1, YEAR1), NewMDsYDates(Month2, DayPlus2, YEAR1)...)...)},
  /*  15 */ func (DayPlus11 []string, Month1 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewDsMYDates(DayPlus11, Month1, YEAR1)...)},
  /*  16 */ func (DayPlus1 []string, Month1 string, AND1 string, DayPlus2 []string, Month2 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStartDates(append(NewDsMYDates(DayPlus1, Month1, YEAR1), NewDsMYDates(DayPlus2, Month2, YEAR1)...)...)},
  /*  17 */ func (Month1 string, Day1 string, RangeSep1 string, Day2 string, Day3 string, RangeSep2 string, Day4 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewMDYDate(Month1, Day1, ""), NewMDYDate(Month1, Day2, "")), NewRangeWithStartEndDates(NewMDYDate(Month1, Day3, ""), NewMDYDate(Month1, Day4, "")))},
  /*  18 */ func (Day1 string, RangeSep1 string, Day2 string, Day3 string, RangeSep2 string, Day4 string, Month1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewDMYDate(Day1, Month1, ""), NewDMYDate(Day2, Month1, "")), NewRangeWithStartEndDates(NewDMYDate(Day3, Month1, ""), NewDMYDate(Day4, Month1, "")))},
  /*  19 */ func (Month1 string, Day1 string, RangeSep1 string, Day2 string, Month2 string, Day3 string, RangeSep2 string, Day4 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewMDYDate(Month1, Day1, ""), NewMDYDate(Month1, Day2, "")), NewRangeWithStartEndDates(NewMDYDate(Month2, Day3, ""), NewMDYDate(Month2, Day4, "")))},
  /*  20 */ func (Day1 string, RangeSep1 string, Day2 string, Month1 string, Day3 string, RangeSep2 string, Day4 string, Month2 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewDMYDate(Day1, Month1, ""), NewDMYDate(Day2, Month1, "")), NewRangeWithStartEndDates(NewDMYDate(Day3, Month2, ""), NewDMYDate(Day4, Month2, "")))},
  /*  21 */ func (Month1 string, Day1 string, RangeSep1 string, Day2 string, Day3 string, RangeSep2 string, Day4 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewMDYDate(Month1, Day1, YEAR1), NewMDYDate(Month1, Day2, YEAR1)), NewRangeWithStartEndDates(NewMDYDate(Month1, Day3, YEAR1), NewMDYDate(Month1, Day4, YEAR1)))},
  /*  22 */ func (Day1 string, RangeSep1 string, Day2 string, Day3 string, RangeSep2 string, Day4 string, Month1 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewDMYDate(Day1, Month1, YEAR1), NewDMYDate(Day2, Month1, YEAR1)), NewRangeWithStartEndDates(NewDMYDate(Day3, Month1, YEAR1), NewDMYDate(Day4, Month1, YEAR1)))},
  /*  23 */ func (Month1 string, Day1 string, RangeSep1 string, Day2 string, Month2 string, Day3 string, RangeSep2 string, Day4 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewMDYDate(Month1, Day1, YEAR1), NewMDYDate(Month1, Day2, YEAR1)), NewRangeWithStartEndDates(NewMDYDate(Month2, Day3, YEAR1), NewMDYDate(Month2, Day4, YEAR1)))},
  /*  24 */ func (Day1 string, RangeSep1 string, Day2 string, Month1 string, Day3 string, RangeSep2 string, Day4 string, Month2 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewDMYDate(Day1, Month1, YEAR1), NewDMYDate(Day2, Month1, YEAR1)), NewRangeWithStartEndDates(NewDMYDate(Day3, Month2, YEAR1), NewDMYDate(Day4, Month2, YEAR1)))},
  /*  25 */ func (Month1 string, Day1 string, Month2 string, Day2 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStart(NewMDYDate(Month1, Day1, YEAR1)), NewRangeWithStart(NewMDYDate(Month2, Day2, YEAR1)))},
  /*  26 */ func (RangesPrefix1 string) string {return RangesPrefix1},
  /*  27 */ func (RangesPrefixPlus1 string, RangesPrefix1 string) string {return RangesPrefixPlus1},
  /*  28 */ func (WHEN1 string) string {return WHEN1},
  /*  29 */ func (Day1 string) []string {return []string{Day1}},
  /*  30 */ func (DayPlus11 []string) []string {return DayPlus11},
  /*  31 */ func (Day1 string, Day2 string) []string {return []string{Day1, Day2}},
  /*  32 */ func (DayPlus11 []string, Day1 string) []string {return append(DayPlus11, Day1)},
  /*  33 */ func (RangePrefixPlus1 string, DateTimeTZRange1 *DateTimeTZRange) *DateTimeTZRange {return DateTimeTZRange1},
  /*  34 */ func (DateTimeTZ1 *DateTimeTZ) *DateTimeTZRange {return &DateTimeTZRange{Start: DateTimeTZ1}},
  /*  35 */ func (DateTimeTZ1 *DateTimeTZ, RangeSepPlus1 string, Time1 civil.Time) *DateTimeTZRange {return NewRangeWithStartEndDateTimes(DateTimeTZ1, NewDateTime(DateTimeTZ1.Date, Time1, ""))},
  /*  36 */ func (DateTimeTZ1 *DateTimeTZ, RangeSepPlus1 string, DateTimeTZ2 *DateTimeTZ) *DateTimeTZRange {return &DateTimeTZRange{Start: DateTimeTZ1, End: DateTimeTZ2}},
  /*  37 */ func (Month1 string, Day1 string, RangeSepPlus1 string, Day2 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewMDYDate(Month1, Day1, ""), NewMDYDate(Month1, Day2, ""))},
  /*  38 */ func (Day1 string, RangeSepPlus1 string, Day2 string, Month1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewDMYDate(Day1, Month1, ""), NewDMYDate(Day2, Month1, ""))},
  /*  39 */ func (Month1 string, Day1 string, RangeSepPlus1 string, Day2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewMDYDate(Month1, Day1, YEAR1), NewMDYDate(Month1, Day2, YEAR1))},
  /*  40 */ func (Day1 string, RangeSepPlus1 string, Day2 string, Month1 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewDMYDate(Day1, Month1, YEAR1), NewDMYDate(Day2, Month1, YEAR1))},
  /*  41 */ func (Month1 string, Day1 string, RangeSepPlus1 string, Month2 string, Day2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewMDYDate(Month1, Day1, YEAR1), NewMDYDate(Month2, Day2, YEAR1))},
  /*  42 */ func (WeekDay1 string, Month1 string, Day1 string, RangeSepPlus1 string, WeekDay2 string, Month2 string, Day2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewMDYDate(Month1, Day1, YEAR1), NewMDYDate(Month2, Day2, YEAR1))},
  /*  43 */ func (WeekDay1 string, Month1 string, Day1 string, RangeSepPlus1 string, WeekDay2 string, Day2 string, Month2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewMDYDate(Month1, Day1, YEAR1), NewDMYDate(Day2, Month2, YEAR1))},
  /*  44 */ func (WeekDay1 string, Month1 string, Day1 string, RangeSepPlus1 string, Day2 string, WeekDay2 string, Month2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewMDYDate(Month1, Day1, YEAR1), NewDMYDate(Day2, Month2, YEAR1))},
  /*  45 */ func (WeekDay1 string, Day1 string, Month1 string, RangeSepPlus1 string, WeekDay2 string, Month2 string, Day2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewDMYDate(Day1, Month1, YEAR1), NewMDYDate(Month2, Day2, YEAR1))},
  /*  46 */ func (WeekDay1 string, Day1 string, Month1 string, RangeSepPlus1 string, WeekDay2 string, Day2 string, Month2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewDMYDate(Day1, Month1, YEAR1), NewDMYDate(Day2, Month2, YEAR1))},
  /*  47 */ func (WeekDay1 string, Day1 string, Month1 string, RangeSepPlus1 string, Day2 string, WeekDay2 string, Month2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewDMYDate(Day1, Month1, YEAR1), NewDMYDate(Day2, Month2, YEAR1))},
  /*  48 */ func (Time1 civil.Time, DateTimeSepOpt1 string, Day1 string, Month1 string, RangeSepPlus1 string, Day2 string, Month2 string, DateTimeSepOpt2 string, Time2 civil.Time, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDateTimes(NewDateTime(NewDMYDate(Day1, Month1, YEAR1), Time1, ""), NewDateTime(NewDMYDate(Day2, Month2, YEAR1), Time2, ""))},
  /*  49 */ func (Date1 civil.Date, Time1 civil.Time, Time2 civil.Time) *DateTimeTZRange {return NewRangeWithStartEndDateTimes(NewDateTime(Date1, Time1, ""), NewDateTime(Date1, Time2, ""))},
  /*  50 */ func (RangePrefix1 string) string {return RangePrefix1},
  /*  51 */ func (RangePrefixPlus1 string, RangePrefix1 string) string {return RangePrefixPlus1},
  /*  52 */ func (BEGINNING1 string) string {return BEGINNING1},
  /*  53 */ func (FROM1 string) string {return FROM1},
  /*  54 */ func (RangeSep1 string) string {return RangeSep1},
  /*  55 */ func (RangeSepPlus1 string, RangeSep1 string) string {return RangeSepPlus1},
  /*  56 */ func (DEC1 string) string {return DEC1},
  /*  57 */ func (SUB1 string) string {return SUB1},
  /*  58 */ func (THROUGH1 string) string {return THROUGH1},
  /*  59 */ func (TILL1 string) string {return TILL1},
  /*  60 */ func (TO1 string) string {return TO1},
  /*  61 */ func (UNTIL1 string) string {return UNTIL1},
  /*  62 */ func (Date1 civil.Date) *DateTimeTZ {return NewDateTimeWithDate(Date1)},
  /*  63 */ func (Date1 civil.Date, Time1 civil.Time) *DateTimeTZ {return NewDateTime(Date1, Time1, "")},
  /*  64 */ func (Date1 civil.Date, DateTimeSepPlus1 string, Time1 civil.Time) *DateTimeTZ {return NewDateTime(Date1, Time1, "")},
  /*  65 */ func (Time1 civil.Time, Date1 civil.Date) *DateTimeTZ {return NewDateTime(Date1, Time1, "")},
  /*  66 */ func (Time1 civil.Time, DateTimeSepPlus1 string, Date1 civil.Date) *DateTimeTZ {return NewDateTime(Date1, Time1, "")},
  /*  67 */ func () string {return ""},
  /*  68 */ func (DateTimeSepPlus1 string) string {return DateTimeSepPlus1},
  /*  69 */ func (DateTimeSep1 string) string {return DateTimeSep1},
  /*  70 */ func (DateTimeSepPlus1 string, DateTimeSep1 string) string {return DateTimeSepPlus1},
  /*  71 */ func (DEC1 string) string {return DEC1},
  /*  72 */ func (SUB1 string) string {return SUB1},
  /*  73 */ func (DatePrefixPlus1 string, Date1 civil.Date) civil.Date {return Date1},
  /*  74 */ func (Date1 civil.Date, DateSuffixPlus1 string) civil.Date {return Date1},
  /*  75 */ func (Day1 string, DateSepPlus1 string, Day2 string) civil.Date {return NewAmbiguousDate(Day1, Day2, "")},
  /*  76 */ func (Year1 string) civil.Date {return NewDMYDate("", "", Year1)},
  /*  77 */ func (Year1 string, DateSepPlus1 string, Day1 string) civil.Date {return NewDMYDate("", Day1, Year1)},
  /*  78 */ func (Year1 string, DateSepPlus1 string, Day1 string, DateSepPlus2 string, Day2 string) civil.Date {return NewDMYDate(Day2, Day1, Year1)},
  /*  79 */ func (Month1 string) civil.Date {return NewMDYDate(Month1, "", "")},
  /*  80 */ func (Month1 string, Year1 string) civil.Date {return NewMDYDate(Month1, "", Year1)},
  /*  81 */ func (Month1 string, Day1 string, Year1 string) civil.Date {return NewMDYDate(Month1, Day1, Year1)},
  /*  82 */ func (Day1 string, Month1 string, Year1 string) civil.Date {return NewDMYDate(Day1, Month1, Year1)},
  /*  83 */ func (Day1 string, DateSepPlus1 string, Day2 string, DateSepPlus2 string, Year1 string) civil.Date {return NewAmbiguousDate(Day1, Day2, Year1)},
  /*  84 */ func (Month1 string, Day1 string) civil.Date {return NewMDYDate(Month1, Day1, "")},
  /*  85 */ func (Day1 string, Month1 string) civil.Date {return NewDMYDate(Day1, Month1, "")},
  /*  86 */ func (Year1 string, Month1 string, Day1 string) civil.Date {return NewDMYDate(Day1, Month1, Year1)},
  /*  87 */ func (DatePrefix1 string) string {return DatePrefix1},
  /*  88 */ func (DatePrefixPlus1 string, DatePrefix1 string) string {return DatePrefixPlus1},
  /*  89 */ func (DATE1 string) string {return DATE1},
  /*  90 */ func (WeekDay1 string) string {return WeekDay1},
  /*  91 */ func (COLON1 string) string {return COLON1},
  /*  92 */ func (COMMA1 string) string {return COMMA1},
  /*  93 */ func (TIME1 string) string {return TIME1},
  /*  94 */ func (DateSep1 string) string {return DateSep1},
  /*  95 */ func (DateSepPlus1 string, DateSep1 string) string {return DateSepPlus1},
  /*  96 */ func (COMMA1 string) string {return COMMA1},
  /*  97 */ func (DEC1 string) string {return DEC1},
  /*  98 */ func (PERIOD1 string) string {return PERIOD1},
  /*  99 */ func (SUB1 string) string {return SUB1},
  /* 100 */ func (QUO1 string) string {return QUO1},
  /* 101 */ func (DateSuffix1 string) string {return DateSuffix1},
  /* 102 */ func (DateSuffixPlus1 string, DateSuffix1 string) string {return DateSuffixPlus1},
  /* 103 */ func (T1 string) string {return T1},
  /* 104 */ func (INT1 string) string {return INT1},
  /* 105 */ func (INT1 string, DaySuffixPlus1 string) string {return INT1},
  /* 106 */ func (DaySuffix1 string) string {return DaySuffix1},
  /* 107 */ func (DaySuffixPlus1 string, DaySuffix1 string) string {return DaySuffixPlus1},
  /* 108 */ func (COMMA1 string) string {return COMMA1},
  /* 109 */ func (ORD_IND1 string) string {return ORD_IND1},
  /* 110 */ func (PERIOD1 string) string {return PERIOD1},
  /* 111 */ func (TH1 string) string {return TH1},
  /* 112 */ func (MONTH_NAME1 string) string {return MONTH_NAME1},
  /* 113 */ func (MONTH_NAME1 string, MonthSuffixPlus1 string) string {return MONTH_NAME1},
  /* 114 */ func (MonthSuffix1 string) string {return MonthSuffix1},
  /* 115 */ func (MonthSuffixPlus1 string, MonthSuffix1 string) string {return MonthSuffixPlus1},
  /* 116 */ func (COMMA1 string) string {return COMMA1},
  /* 117 */ func (YEAR1 string) string {return YEAR1},
  /* 118 */ func (YEAR1 string, YearSuffixPlus1 string) string {return YEAR1},
  /* 119 */ func (YearSuffix1 string) string {return YearSuffix1},
  /* 120 */ func (YearSuffixPlus1 string, YearSuffix1 string) string {return YearSuffixPlus1},
  /* 121 */ func (COMMA1 string) string {return COMMA1},
  /* 122 */ func (TH1 string) string {return TH1},
  /* 123 */ func (WEEKDAY_NAME1 string) string {return WEEKDAY_NAME1},
  /* 124 */ func (TimePrefixPlus1 string, Time1 civil.Time) civil.Time {return Time1},
  /* 125 */ func (INT1 string, AM1 string) civil.Time {return NewTime(INT1, "", "", "")},
  /* 126 */ func (INT1 string, PM1 string) civil.Time {return NewTime((mustAtoi(INT1) % 12) + 12, "", "", "")},
  /* 127 */ func (INT1 string, TimeSep1 string, INT2 string) civil.Time {return NewTime(INT1, INT2, "", "")},
  /* 128 */ func (INT1 string, TimeSep1 string, INT2 string, TimeSep2 string, INT3 string) civil.Time {return NewTime((mustAtoi(INT1) % 12) + 12, INT2, INT3, "")},
  /* 129 */ func (INT1 string, TimeSep1 string, INT2 string, AM1 string) civil.Time {return NewTime(INT1, INT2, "", "")},
  /* 130 */ func (INT1 string, TimeSep1 string, INT2 string, PM1 string) civil.Time {return NewTime((mustAtoi(INT1) % 12) + 12, INT2, "", "")},
  /* 131 */ func (TimePrefix1 string) string {return TimePrefix1},
  /* 132 */ func (TimePrefixPlus1 string, TimePrefix1 string) string {return TimePrefixPlus1},
  /* 133 */ func (AT1 string) string {return AT1},
  /* 134 */ func (COLON1 string) string {return COLON1},
  /* 135 */ func (FROM1 string) string {return FROM1},
  /* 136 */ func (ON1 string) string {return ON1},
  /* 137 */ func (TIME1 string) string {return TIME1},
  /* 138 */ func (COLON1 string) string {return COLON1},
  /* 139 */ func (PERIOD1 string) string {return PERIOD1},
  /* 140 */ func () string {return ""},
  /* 141 */ func (AND1 string) string {return AND1},
  /* 142 */ func () string {return ""},
  /* 143 */ func (OF1 string) string {return OF1},
}}

var parseStates = &glr.ParseStates{Items:[]glr.ParseState{
  /*   0 */ glr.ParseState{Actions:map[string][]glr.Action{"AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}, "BEGINNING":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:30, Rule:0}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}, "DATE":[]glr.Action{glr.Action{Type:"shift", State:33, Rule:0}}, "FROM":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:16, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "ON":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:32, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "WHEN":[]glr.Action{glr.Action{Type:"shift", State:17, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Date":14, "DatePrefix":27, "DatePrefixPlus":22, "DateTimeTZ":11, "DateTimeTZRange":4, "DateTimeTZRanges":2, "Day":8, "DayPlus":7, "DayPlus1":6, "Month":5, "RangePrefix":18, "RangePrefixPlus":10, "RangesPrefix":9, "RangesPrefixPlus":3, "Time":13, "TimePrefix":26, "TimePrefixPlus":21, "WeekDay":12, "Year":23, "root":1}},
  /*   1 */ glr.ParseState{Actions:map[string][]glr.Action{"$end":[]glr.Action{glr.Action{Type:"accept", State:0, Rule:0}}}, Gotos:map[string]int{}},
  /*   2 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:140}, glr.Action{Type:"reduce", State:0, Rule:1}}, "AND":[]glr.Action{glr.Action{Type:"shift", State:37, Rule:0}}}, Gotos:map[string]int{"$$1":35, "AndOpt":36}},
  /*   3 */ glr.ParseState{Actions:map[string][]glr.Action{"AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}, "BEGINNING":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:30, Rule:0}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}, "DATE":[]glr.Action{glr.Action{Type:"shift", State:33, Rule:0}}, "FROM":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:16, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "ON":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:32, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "WHEN":[]glr.Action{glr.Action{Type:"shift", State:17, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Date":14, "DatePrefix":27, "DatePrefixPlus":22, "DateTimeTZ":11, "DateTimeTZRange":4, "DateTimeTZRanges":38, "Day":8, "DayPlus":7, "DayPlus1":6, "Month":5, "RangePrefix":18, "RangePrefixPlus":10, "RangesPrefix":39, "RangesPrefixPlus":3, "Time":13, "TimePrefix":26, "TimePrefixPlus":21, "WeekDay":12, "Year":23}},
  /*   4 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:9}}}, Gotos:map[string]int{}},
  /*   5 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:79}, glr.Action{Type:"reduce", State:0, Rule:79}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Day":42, "DayPlus":41, "DayPlus1":40, "Year":43}},
  /*   6 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:30}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Day":46, "Month":45}},
  /*   7 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":47}},
  /*   8 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:29}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:53, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:61, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:54, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"DateSep":59, "DateSepPlus":51, "Day":49, "Month":52, "RangeSep":48, "RangeSepPlus":50}},
  /*   9 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:26}}}, Gotos:map[string]int{}},
  /*  10 */ glr.ParseState{Actions:map[string][]glr.Action{"AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}, "BEGINNING":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:30, Rule:0}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}, "DATE":[]glr.Action{glr.Action{Type:"shift", State:33, Rule:0}}, "FROM":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:16, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "ON":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:32, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Date":14, "DatePrefix":27, "DatePrefixPlus":22, "DateTimeTZ":11, "DateTimeTZRange":63, "Day":66, "Month":65, "RangePrefix":64, "RangePrefixPlus":10, "Time":13, "TimePrefix":26, "TimePrefixPlus":21, "WeekDay":12, "Year":23}},
  /*  11 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:34}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"RangeSep":68, "RangeSepPlus":67}},
  /*  12 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:90}, glr.Action{Type:"reduce", State:0, Rule:90}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Day":72, "Month":71}},
  /*  13 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:67}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:82, Rule:0}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}, "DATE":[]glr.Action{glr.Action{Type:"shift", State:33, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:79, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:80, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:83, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Date":74, "DatePrefix":27, "DatePrefixPlus":22, "DateTimeSep":78, "DateTimeSepOpt":73, "DateTimeSepPlus":75, "Day":76, "Month":77, "WeekDay":81, "Year":23}},
  /*  14 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:62}, glr.Action{Type:"reduce", State:0, Rule:62}}, "AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:90, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:79, Rule:0}}, "FROM":[]glr.Action{glr.Action{Type:"shift", State:91, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:87, Rule:0}}, "ON":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:80, Rule:0}}, "T":[]glr.Action{glr.Action{Type:"shift", State:89, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:92, Rule:0}}}, Gotos:map[string]int{"DateSuffix":88, "DateSuffixPlus":86, "DateTimeSep":78, "DateTimeSepPlus":85, "Time":84, "TimePrefix":26, "TimePrefixPlus":21}},
  /*  15 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:112}, glr.Action{Type:"reduce", State:0, Rule:112}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:95, Rule:0}}}, Gotos:map[string]int{"MonthSuffix":94, "MonthSuffixPlus":93}},
  /*  16 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:104}, glr.Action{Type:"reduce", State:0, Rule:104}}, "AM":[]glr.Action{glr.Action{Type:"shift", State:97, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:101, Rule:0}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:103, Rule:0}}, "ORD_IND":[]glr.Action{glr.Action{Type:"shift", State:104, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:102, Rule:0}}, "PM":[]glr.Action{glr.Action{Type:"shift", State:98, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:105, Rule:0}}}, Gotos:map[string]int{"DaySuffix":100, "DaySuffixPlus":96, "TimeSep":99}},
  /*  17 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:28}}}, Gotos:map[string]int{}},
  /*  18 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:50}}}, Gotos:map[string]int{}},
  /*  19 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:122}}}, Gotos:map[string]int{}},
  /*  20 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:123}}}, Gotos:map[string]int{}},
  /*  21 */ glr.ParseState{Actions:map[string][]glr.Action{"AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:90, Rule:0}}, "FROM":[]glr.Action{glr.Action{Type:"shift", State:91, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:87, Rule:0}}, "ON":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:92, Rule:0}}}, Gotos:map[string]int{"Time":106, "TimePrefix":107, "TimePrefixPlus":21}},
  /*  22 */ glr.ParseState{Actions:map[string][]glr.Action{"COLON":[]glr.Action{glr.Action{Type:"shift", State:82, Rule:0}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}, "DATE":[]glr.Action{glr.Action{Type:"shift", State:33, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:83, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Date":108, "DatePrefix":109, "DatePrefixPlus":22, "Day":76, "Month":77, "WeekDay":81, "Year":23}},
  /*  23 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:76}, glr.Action{Type:"reduce", State:0, Rule:76}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:112, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:61, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:113, Rule:0}}}, Gotos:map[string]int{"DateSep":59, "DateSepPlus":110, "Month":111}},
  /*  24 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:52}}}, Gotos:map[string]int{}},
  /*  25 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:135}, glr.Action{Type:"reduce", State:0, Rule:53}}}, Gotos:map[string]int{}},
  /*  26 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:131}}}, Gotos:map[string]int{}},
  /*  27 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:87}}}, Gotos:map[string]int{}},
  /*  28 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:117}, glr.Action{Type:"reduce", State:0, Rule:117}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:116, Rule:0}}}, Gotos:map[string]int{"YearSuffix":115, "YearSuffixPlus":114}},
  /*  29 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:133}}}, Gotos:map[string]int{}},
  /*  30 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:134}, glr.Action{Type:"reduce", State:0, Rule:91}}, "AT":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:134}}, "FROM":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:134}}, "ON":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:134}}}, Gotos:map[string]int{}},
  /*  31 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:136}}}, Gotos:map[string]int{}},
  /*  32 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:137}, glr.Action{Type:"reduce", State:0, Rule:93}}, "AT":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:137}}, "FROM":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:137}}, "ON":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:137}}}, Gotos:map[string]int{}},
  /*  33 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:89}}}, Gotos:map[string]int{}},
  /*  34 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:92}}}, Gotos:map[string]int{}},
  /*  35 */ glr.ParseState{Actions:map[string][]glr.Action{"AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}, "BEGINNING":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:30, Rule:0}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}, "DATE":[]glr.Action{glr.Action{Type:"shift", State:33, Rule:0}}, "FROM":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:16, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "ON":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:32, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "WHEN":[]glr.Action{glr.Action{Type:"shift", State:17, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Date":14, "DatePrefix":27, "DatePrefixPlus":22, "DateTimeTZ":11, "DateTimeTZRange":4, "DateTimeTZRanges":117, "Day":8, "DayPlus":7, "DayPlus1":6, "Month":5, "RangePrefix":18, "RangePrefixPlus":10, "RangesPrefix":9, "RangesPrefixPlus":3, "Time":13, "TimePrefix":26, "TimePrefixPlus":21, "WeekDay":12, "Year":23}},
  /*  36 */ glr.ParseState{Actions:map[string][]glr.Action{"AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}, "BEGINNING":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:30, Rule:0}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}, "DATE":[]glr.Action{glr.Action{Type:"shift", State:33, Rule:0}}, "FROM":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:16, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "ON":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:32, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Date":14, "DatePrefix":27, "DatePrefixPlus":22, "DateTimeTZ":11, "DateTimeTZRange":118, "Day":66, "Month":65, "RangePrefix":18, "RangePrefixPlus":10, "Time":13, "TimePrefix":26, "TimePrefixPlus":21, "WeekDay":12, "Year":23}},
  /*  37 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:141}}}, Gotos:map[string]int{}},
  /*  38 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:8}, glr.Action{Type:"reduce", State:0, Rule:140}, glr.Action{Type:"reduce", State:0, Rule:8}}, "AND":[]glr.Action{glr.Action{Type:"shift", State:37, Rule:0}}}, Gotos:map[string]int{"AndOpt":36}},
  /*  39 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:27}, glr.Action{Type:"reduce", State:0, Rule:26}}}, Gotos:map[string]int{}},
  /*  40 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:11}, glr.Action{Type:"reduce", State:0, Rule:30}, glr.Action{Type:"reduce", State:0, Rule:11}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:119, Rule:0}}}, Gotos:map[string]int{"Day":46}},
  /*  41 */ glr.ParseState{Actions:map[string][]glr.Action{"AND":[]glr.Action{glr.Action{Type:"shift", State:120, Rule:0}}}, Gotos:map[string]int{}},
  /*  42 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:84}, glr.Action{Type:"reduce", State:0, Rule:84}}, "AND":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:29}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Day":49, "Month":122, "RangeSep":121, "RangeSepPlus":123, "Year":124}},
  /*  43 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:80}}}, Gotos:map[string]int{}},
  /*  44 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:104}, glr.Action{Type:"reduce", State:0, Rule:104}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:103, Rule:0}}, "ORD_IND":[]glr.Action{glr.Action{Type:"shift", State:104, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:125, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:105, Rule:0}}}, Gotos:map[string]int{"DaySuffix":100, "DaySuffixPlus":96}},
  /*  45 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:12}, glr.Action{Type:"reduce", State:0, Rule:12}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:126, Rule:0}}}, Gotos:map[string]int{}},
  /*  46 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:32}}}, Gotos:map[string]int{}},
  /*  47 */ glr.ParseState{Actions:map[string][]glr.Action{"AND":[]glr.Action{glr.Action{Type:"shift", State:127, Rule:0}}}, Gotos:map[string]int{}},
  /*  48 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:54}, glr.Action{Type:"reduce", State:0, Rule:54}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":128}},
  /*  49 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:31}}}, Gotos:map[string]int{}},
  /*  50 */ glr.ParseState{Actions:map[string][]glr.Action{"DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"Day":129, "RangeSep":130}},
  /*  51 */ glr.ParseState{Actions:map[string][]glr.Action{"COMMA":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:112, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:61, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:113, Rule:0}}}, Gotos:map[string]int{"DateSep":132, "Day":131}},
  /*  52 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:85}, glr.Action{Type:"reduce", State:0, Rule:85}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Year":133}},
  /*  53 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:97}, glr.Action{Type:"reduce", State:0, Rule:56}}, "COMMA":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:97}}, "PERIOD":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:97}}, "QUO":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:97}}}, Gotos:map[string]int{}},
  /*  54 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:99}, glr.Action{Type:"reduce", State:0, Rule:57}}, "COMMA":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:99}}, "PERIOD":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:99}}, "QUO":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:99}}}, Gotos:map[string]int{}},
  /*  55 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:58}}}, Gotos:map[string]int{}},
  /*  56 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:59}}}, Gotos:map[string]int{}},
  /*  57 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:60}}}, Gotos:map[string]int{}},
  /*  58 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:61}}}, Gotos:map[string]int{}},
  /*  59 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:94}}}, Gotos:map[string]int{}},
  /*  60 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:96}}}, Gotos:map[string]int{}},
  /*  61 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:98}}}, Gotos:map[string]int{}},
  /*  62 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:100}}}, Gotos:map[string]int{}},
  /*  63 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:33}}}, Gotos:map[string]int{}},
  /*  64 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:51}, glr.Action{Type:"reduce", State:0, Rule:50}}}, Gotos:map[string]int{}},
  /*  65 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:79}, glr.Action{Type:"reduce", State:0, Rule:79}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Day":134, "Year":43}},
  /*  66 */ glr.ParseState{Actions:map[string][]glr.Action{"COMMA":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:53, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:61, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:54, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"DateSep":59, "DateSepPlus":51, "Month":52, "RangeSep":68, "RangeSepPlus":50}},
  /*  67 */ glr.ParseState{Actions:map[string][]glr.Action{"AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:30, Rule:0}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}, "DATE":[]glr.Action{glr.Action{Type:"shift", State:33, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "FROM":[]glr.Action{glr.Action{Type:"shift", State:91, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:16, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "ON":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:32, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Date":137, "DatePrefix":27, "DatePrefixPlus":22, "DateTimeTZ":136, "Day":76, "Month":77, "RangeSep":130, "Time":135, "TimePrefix":26, "TimePrefixPlus":21, "WeekDay":81, "Year":23}},
  /*  68 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:54}}}, Gotos:map[string]int{}},
  /*  69 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:56}}}, Gotos:map[string]int{}},
  /*  70 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:57}}}, Gotos:map[string]int{}},
  /*  71 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":138}},
  /*  72 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":139}},
  /*  73 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":140}},
  /*  74 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:65}}, "T":[]glr.Action{glr.Action{Type:"shift", State:89, Rule:0}}}, Gotos:map[string]int{"DateSuffix":88, "DateSuffixPlus":86}},
  /*  75 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:68}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:82, Rule:0}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}, "DATE":[]glr.Action{glr.Action{Type:"shift", State:33, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:79, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:80, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:83, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Date":141, "DatePrefix":27, "DatePrefixPlus":22, "DateTimeSep":142, "Day":76, "Month":77, "WeekDay":81, "Year":23}},
  /*  76 */ glr.ParseState{Actions:map[string][]glr.Action{"COMMA":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:112, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:61, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:113, Rule:0}}}, Gotos:map[string]int{"DateSep":59, "DateSepPlus":51, "Month":52}},
  /*  77 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:79}, glr.Action{Type:"reduce", State:0, Rule:79}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Day":143, "Year":43}},
  /*  78 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:69}}}, Gotos:map[string]int{}},
  /*  79 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:71}}}, Gotos:map[string]int{}},
  /*  80 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:72}}}, Gotos:map[string]int{}},
  /*  81 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:90}}}, Gotos:map[string]int{}},
  /*  82 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:91}}}, Gotos:map[string]int{}},
  /*  83 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:93}}}, Gotos:map[string]int{}},
  /*  84 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:63}, glr.Action{Type:"reduce", State:0, Rule:63}}, "AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:90, Rule:0}}, "FROM":[]glr.Action{glr.Action{Type:"shift", State:91, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:87, Rule:0}}, "ON":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:92, Rule:0}}}, Gotos:map[string]int{"Time":144, "TimePrefix":26, "TimePrefixPlus":21}},
  /*  85 */ glr.ParseState{Actions:map[string][]glr.Action{"AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:90, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:79, Rule:0}}, "FROM":[]glr.Action{glr.Action{Type:"shift", State:91, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:87, Rule:0}}, "ON":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:80, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:92, Rule:0}}}, Gotos:map[string]int{"DateTimeSep":142, "Time":145, "TimePrefix":26, "TimePrefixPlus":21}},
  /*  86 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:74}, glr.Action{Type:"reduce", State:0, Rule:74}}, "T":[]glr.Action{glr.Action{Type:"shift", State:89, Rule:0}}}, Gotos:map[string]int{"DateSuffix":146}},
  /*  87 */ glr.ParseState{Actions:map[string][]glr.Action{"AM":[]glr.Action{glr.Action{Type:"shift", State:97, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:101, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:147, Rule:0}}, "PM":[]glr.Action{glr.Action{Type:"shift", State:98, Rule:0}}}, Gotos:map[string]int{"TimeSep":99}},
  /*  88 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:101}}}, Gotos:map[string]int{}},
  /*  89 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:103}}}, Gotos:map[string]int{}},
  /*  90 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:134}}}, Gotos:map[string]int{}},
  /*  91 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:135}}}, Gotos:map[string]int{}},
  /*  92 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:137}}}, Gotos:map[string]int{}},
  /*  93 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:113}, glr.Action{Type:"reduce", State:0, Rule:113}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:95, Rule:0}}}, Gotos:map[string]int{"MonthSuffix":148}},
  /*  94 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:114}}}, Gotos:map[string]int{}},
  /*  95 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:116}}}, Gotos:map[string]int{}},
  /*  96 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:105}, glr.Action{Type:"reduce", State:0, Rule:105}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:103, Rule:0}}, "ORD_IND":[]glr.Action{glr.Action{Type:"shift", State:104, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:125, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:105, Rule:0}}}, Gotos:map[string]int{"DaySuffix":149}},
  /*  97 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:125}}}, Gotos:map[string]int{}},
  /*  98 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:126}}}, Gotos:map[string]int{}},
  /*  99 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:150, Rule:0}}}, Gotos:map[string]int{}},
  /* 100 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:106}}}, Gotos:map[string]int{}},
  /* 101 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:138}}}, Gotos:map[string]int{}},
  /* 102 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:139}, glr.Action{Type:"reduce", State:0, Rule:110}}}, Gotos:map[string]int{}},
  /* 103 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:108}}}, Gotos:map[string]int{}},
  /* 104 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:109}}}, Gotos:map[string]int{}},
  /* 105 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:111}}}, Gotos:map[string]int{}},
  /* 106 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:124}}}, Gotos:map[string]int{}},
  /* 107 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:132}, glr.Action{Type:"reduce", State:0, Rule:131}}}, Gotos:map[string]int{}},
  /* 108 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:73}, glr.Action{Type:"reduce", State:0, Rule:73}}, "T":[]glr.Action{glr.Action{Type:"shift", State:89, Rule:0}}}, Gotos:map[string]int{"DateSuffix":88, "DateSuffixPlus":86}},
  /* 109 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:88}, glr.Action{Type:"reduce", State:0, Rule:87}}}, Gotos:map[string]int{}},
  /* 110 */ glr.ParseState{Actions:map[string][]glr.Action{"COMMA":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:112, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:61, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:113, Rule:0}}}, Gotos:map[string]int{"DateSep":132, "Day":151}},
  /* 111 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":152}},
  /* 112 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:97}}}, Gotos:map[string]int{}},
  /* 113 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:99}}}, Gotos:map[string]int{}},
  /* 114 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:118}, glr.Action{Type:"reduce", State:0, Rule:118}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:116, Rule:0}}}, Gotos:map[string]int{"YearSuffix":153}},
  /* 115 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:119}}}, Gotos:map[string]int{}},
  /* 116 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:121}}}, Gotos:map[string]int{}},
  /* 117 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:140}}, "AND":[]glr.Action{glr.Action{Type:"shift", State:37, Rule:0}}, "CALENDAR":[]glr.Action{glr.Action{Type:"shift", State:157, Rule:0}}, "GOOGLE":[]glr.Action{glr.Action{Type:"shift", State:156, Rule:0}}, "ICS":[]glr.Action{glr.Action{Type:"shift", State:158, Rule:0}}}, Gotos:map[string]int{"AndOpt":36, "RootSuffix":155, "RootSuffixPlus":154}},
  /* 118 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:10}}}, Gotos:map[string]int{}},
  /* 119 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:13}}}, Gotos:map[string]int{}},
  /* 120 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":159}},
  /* 121 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:54}, glr.Action{Type:"reduce", State:0, Rule:54}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":160}},
  /* 122 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":161}},
  /* 123 */ glr.ParseState{Actions:map[string][]glr.Action{"DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"Day":162, "Month":163, "RangeSep":130}},
  /* 124 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:81}}}, Gotos:map[string]int{}},
  /* 125 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:110}}}, Gotos:map[string]int{}},
  /* 126 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:15}}}, Gotos:map[string]int{}},
  /* 127 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":165, "DayPlus":164, "DayPlus1":166}},
  /* 128 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Day":167, "Month":168}},
  /* 129 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":169}},
  /* 130 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:55}}}, Gotos:map[string]int{}},
  /* 131 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:75}, glr.Action{Type:"reduce", State:0, Rule:75}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:112, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:61, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:113, Rule:0}}}, Gotos:map[string]int{"DateSep":59, "DateSepPlus":170}},
  /* 132 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:95}}}, Gotos:map[string]int{}},
  /* 133 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:82}}}, Gotos:map[string]int{}},
  /* 134 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:84}, glr.Action{Type:"reduce", State:0, Rule:84}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"RangeSep":68, "RangeSepPlus":123, "Year":124}},
  /* 135 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:35}, glr.Action{Type:"reduce", State:0, Rule:35}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:82, Rule:0}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}, "DATE":[]glr.Action{glr.Action{Type:"shift", State:33, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:79, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:80, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:83, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Date":74, "DatePrefix":27, "DatePrefixPlus":22, "DateTimeSep":78, "DateTimeSepPlus":171, "Day":76, "Month":77, "WeekDay":81, "Year":23}},
  /* 136 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:36}}}, Gotos:map[string]int{}},
  /* 137 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:62}, glr.Action{Type:"reduce", State:0, Rule:62}}, "AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:90, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:79, Rule:0}}, "FROM":[]glr.Action{glr.Action{Type:"shift", State:91, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:87, Rule:0}}, "ON":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:80, Rule:0}}, "T":[]glr.Action{glr.Action{Type:"shift", State:89, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:92, Rule:0}}}, Gotos:map[string]int{"DateSuffix":88, "DateSuffixPlus":86, "DateTimeSep":78, "DateTimeSepPlus":85, "Time":172, "TimePrefix":26, "TimePrefixPlus":21}},
  /* 138 */ glr.ParseState{Actions:map[string][]glr.Action{"DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"RangeSep":68, "RangeSepPlus":173}},
  /* 139 */ glr.ParseState{Actions:map[string][]glr.Action{"DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"RangeSep":68, "RangeSepPlus":174}},
  /* 140 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":175}},
  /* 141 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:66}}, "T":[]glr.Action{glr.Action{Type:"shift", State:89, Rule:0}}}, Gotos:map[string]int{"DateSuffix":88, "DateSuffixPlus":86}},
  /* 142 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:70}}}, Gotos:map[string]int{}},
  /* 143 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:84}, glr.Action{Type:"reduce", State:0, Rule:84}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Year":124}},
  /* 144 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:49}}}, Gotos:map[string]int{}},
  /* 145 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:64}}}, Gotos:map[string]int{}},
  /* 146 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:102}}}, Gotos:map[string]int{}},
  /* 147 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:139}}}, Gotos:map[string]int{}},
  /* 148 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:115}}}, Gotos:map[string]int{}},
  /* 149 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:107}}}, Gotos:map[string]int{}},
  /* 150 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:127}, glr.Action{Type:"reduce", State:0, Rule:127}}, "AM":[]glr.Action{glr.Action{Type:"shift", State:177, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:101, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:147, Rule:0}}, "PM":[]glr.Action{glr.Action{Type:"shift", State:178, Rule:0}}}, Gotos:map[string]int{"TimeSep":176}},
  /* 151 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:77}, glr.Action{Type:"reduce", State:0, Rule:77}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:112, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:61, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:113, Rule:0}}}, Gotos:map[string]int{"DateSep":59, "DateSepPlus":179}},
  /* 152 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:86}}}, Gotos:map[string]int{}},
  /* 153 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:120}}}, Gotos:map[string]int{}},
  /* 154 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:2}}, "CALENDAR":[]glr.Action{glr.Action{Type:"shift", State:157, Rule:0}}, "GOOGLE":[]glr.Action{glr.Action{Type:"shift", State:156, Rule:0}}, "ICS":[]glr.Action{glr.Action{Type:"shift", State:158, Rule:0}}}, Gotos:map[string]int{"RootSuffix":180}},
  /* 155 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:3}}}, Gotos:map[string]int{}},
  /* 156 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:5}}}, Gotos:map[string]int{}},
  /* 157 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:6}}}, Gotos:map[string]int{}},
  /* 158 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:7}}}, Gotos:map[string]int{}},
  /* 159 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":165, "DayPlus":181, "DayPlus1":166}},
  /* 160 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Day":182, "Month":183}},
  /* 161 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:184, Rule:0}}}, Gotos:map[string]int{}},
  /* 162 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:37}, glr.Action{Type:"reduce", State:0, Rule:37}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:185, Rule:0}}}, Gotos:map[string]int{}},
  /* 163 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":186}},
  /* 164 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":187}},
  /* 165 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:29}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":49}},
  /* 166 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:30}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":46}},
  /* 167 */ glr.ParseState{Actions:map[string][]glr.Action{"DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"RangeSep":188}},
  /* 168 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":189}},
  /* 169 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:38}, glr.Action{Type:"reduce", State:0, Rule:38}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:190, Rule:0}}}, Gotos:map[string]int{}},
  /* 170 */ glr.ParseState{Actions:map[string][]glr.Action{"COMMA":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:112, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:61, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:113, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"DateSep":132, "Year":191}},
  /* 171 */ glr.ParseState{Actions:map[string][]glr.Action{"COLON":[]glr.Action{glr.Action{Type:"shift", State:82, Rule:0}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}, "DATE":[]glr.Action{glr.Action{Type:"shift", State:33, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:79, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:80, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:83, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}}, Gotos:map[string]int{"Date":141, "DatePrefix":27, "DatePrefixPlus":22, "DateTimeSep":142, "Day":76, "Month":77, "WeekDay":81, "Year":23}},
  /* 172 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:63}}}, Gotos:map[string]int{}},
  /* 173 */ glr.ParseState{Actions:map[string][]glr.Action{"DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}}, Gotos:map[string]int{"Day":193, "RangeSep":130, "WeekDay":192}},
  /* 174 */ glr.ParseState{Actions:map[string][]glr.Action{"DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}}, Gotos:map[string]int{"Day":195, "RangeSep":130, "WeekDay":194}},
  /* 175 */ glr.ParseState{Actions:map[string][]glr.Action{"DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"RangeSep":68, "RangeSepPlus":196}},
  /* 176 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:197, Rule:0}}}, Gotos:map[string]int{}},
  /* 177 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:129}}}, Gotos:map[string]int{}},
  /* 178 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:130}}}, Gotos:map[string]int{}},
  /* 179 */ glr.ParseState{Actions:map[string][]glr.Action{"COMMA":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:112, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:61, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:113, Rule:0}}}, Gotos:map[string]int{"DateSep":132, "Day":198}},
  /* 180 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:4}}}, Gotos:map[string]int{}},
  /* 181 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:199, Rule:0}}}, Gotos:map[string]int{}},
  /* 182 */ glr.ParseState{Actions:map[string][]glr.Action{"DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"RangeSep":200}},
  /* 183 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":201}},
  /* 184 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:25}}}, Gotos:map[string]int{}},
  /* 185 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:39}}}, Gotos:map[string]int{}},
  /* 186 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:202, Rule:0}}}, Gotos:map[string]int{}},
  /* 187 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:203, Rule:0}}}, Gotos:map[string]int{}},
  /* 188 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":204}},
  /* 189 */ glr.ParseState{Actions:map[string][]glr.Action{"DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"RangeSep":205}},
  /* 190 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:40}}}, Gotos:map[string]int{}},
  /* 191 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:83}}}, Gotos:map[string]int{}},
  /* 192 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Day":207, "Month":206}},
  /* 193 */ glr.ParseState{Actions:map[string][]glr.Action{"TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}}, Gotos:map[string]int{"WeekDay":208}},
  /* 194 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Day":210, "Month":209}},
  /* 195 */ glr.ParseState{Actions:map[string][]glr.Action{"TH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}}, Gotos:map[string]int{"WeekDay":211}},
  /* 196 */ glr.ParseState{Actions:map[string][]glr.Action{"DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"Day":212, "RangeSep":130}},
  /* 197 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:128}}}, Gotos:map[string]int{}},
  /* 198 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:78}}}, Gotos:map[string]int{}},
  /* 199 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:14}}}, Gotos:map[string]int{}},
  /* 200 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":213}},
  /* 201 */ glr.ParseState{Actions:map[string][]glr.Action{"DEC":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}, "TILL":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "UNTIL":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{"RangeSep":214}},
  /* 202 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:41}}}, Gotos:map[string]int{}},
  /* 203 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:16}}}, Gotos:map[string]int{}},
  /* 204 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":215}},
  /* 205 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":216}},
  /* 206 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":217}},
  /* 207 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":218}},
  /* 208 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":219}},
  /* 209 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":220}},
  /* 210 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":221}},
  /* 211 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":222}},
  /* 212 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":223}},
  /* 213 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:17}, glr.Action{Type:"reduce", State:0, Rule:17}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:224, Rule:0}}}, Gotos:map[string]int{}},
  /* 214 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"Day":225}},
  /* 215 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:18}, glr.Action{Type:"reduce", State:0, Rule:18}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:226, Rule:0}}}, Gotos:map[string]int{}},
  /* 216 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{"Month":227}},
  /* 217 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:228, Rule:0}}}, Gotos:map[string]int{}},
  /* 218 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:229, Rule:0}}}, Gotos:map[string]int{}},
  /* 219 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:230, Rule:0}}}, Gotos:map[string]int{}},
  /* 220 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:231, Rule:0}}}, Gotos:map[string]int{}},
  /* 221 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:232, Rule:0}}}, Gotos:map[string]int{}},
  /* 222 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:233, Rule:0}}}, Gotos:map[string]int{}},
  /* 223 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:67}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:79, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:80, Rule:0}}}, Gotos:map[string]int{"DateTimeSep":78, "DateTimeSepOpt":234, "DateTimeSepPlus":235}},
  /* 224 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:21}}}, Gotos:map[string]int{}},
  /* 225 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:19}, glr.Action{Type:"reduce", State:0, Rule:19}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:236, Rule:0}}}, Gotos:map[string]int{}},
  /* 226 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:22}}}, Gotos:map[string]int{}},
  /* 227 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:20}, glr.Action{Type:"reduce", State:0, Rule:20}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:237, Rule:0}}}, Gotos:map[string]int{}},
  /* 228 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:42}}}, Gotos:map[string]int{}},
  /* 229 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:43}}}, Gotos:map[string]int{}},
  /* 230 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:44}}}, Gotos:map[string]int{}},
  /* 231 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:45}}}, Gotos:map[string]int{}},
  /* 232 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:46}}}, Gotos:map[string]int{}},
  /* 233 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:47}}}, Gotos:map[string]int{}},
  /* 234 */ glr.ParseState{Actions:map[string][]glr.Action{"AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:90, Rule:0}}, "FROM":[]glr.Action{glr.Action{Type:"shift", State:91, Rule:0}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:87, Rule:0}}, "ON":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "TIME":[]glr.Action{glr.Action{Type:"shift", State:92, Rule:0}}}, Gotos:map[string]int{"Time":238, "TimePrefix":26, "TimePrefixPlus":21}},
  /* 235 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:68}}, "DEC":[]glr.Action{glr.Action{Type:"shift", State:79, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:80, Rule:0}}}, Gotos:map[string]int{"DateTimeSep":142}},
  /* 236 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:23}}}, Gotos:map[string]int{}},
  /* 237 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:24}}}, Gotos:map[string]int{}},
  /* 238 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:239, Rule:0}}}, Gotos:map[string]int{}},
  /* 239 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:48}}}, Gotos:map[string]int{}},
}}

