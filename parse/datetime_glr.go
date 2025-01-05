package parse

import "github.com/findyourpaths/phil/glr"
import "cloud.google.com/go/civil"

/*
Rules

root:
  PrefixOpt DateTimeTZRanges SuffixOpt
DateTimeTZRanges:
  DateTimeTZRange
DateTimeTZRanges:
  DateTimeTZRanges AndOpt DateTimeTZRange
DateTimeTZRanges:
  MONTH_NAME INT INT
DateTimeTZRanges:
  MONTH_NAME INT INT INT
DateTimeTZRanges:
  MONTH_NAME INT INT INT INT
DateTimeTZRanges:
  MONTH_NAME INT INT INT INT INT
DateTimeTZRanges:
  INT INT MONTH_NAME
DateTimeTZRanges:
  INT INT INT MONTH_NAME
DateTimeTZRanges:
  INT INT INT INT MONTH_NAME
DateTimeTZRanges:
  INT INT INT INT INT MONTH_NAME
DateTimeTZRanges:
  MONTH_NAME INT INT YEAR
DateTimeTZRanges:
  MONTH_NAME INT INT INT YEAR
DateTimeTZRanges:
  MONTH_NAME INT INT INT INT YEAR
DateTimeTZRanges:
  MONTH_NAME INT INT INT INT INT YEAR
DateTimeTZRanges:
  INT INT MONTH_NAME YEAR
DateTimeTZRanges:
  INT INT INT MONTH_NAME YEAR
DateTimeTZRanges:
  INT INT INT INT MONTH_NAME YEAR
DateTimeTZRanges:
  INT INT INT INT INT MONTH_NAME YEAR
DateTimeTZRanges:
  MONTH_NAME INT RangeSep INT INT RangeSep INT
DateTimeTZRanges:
  INT RangeSep INT INT RangeSep INT MONTH_NAME
DateTimeTZRanges:
  MONTH_NAME INT RangeSep INT MONTH_NAME INT RangeSep INT
DateTimeTZRanges:
  INT RangeSep INT MONTH_NAME INT RangeSep INT MONTH_NAME
DateTimeTZRanges:
  MONTH_NAME INT RangeSep INT INT RangeSep INT YEAR
DateTimeTZRanges:
  INT RangeSep INT INT RangeSep INT MONTH_NAME YEAR
DateTimeTZRanges:
  MONTH_NAME INT RangeSep INT MONTH_NAME INT RangeSep INT YEAR
DateTimeTZRanges:
  INT RangeSep INT MONTH_NAME INT RangeSep INT MONTH_NAME YEAR
DateTimeTZRanges:
  MONTH_NAME INT MONTH_NAME INT YEAR
DateTimeTZRange:
  DateTimeTZ
DateTimeTZRange:
  DateTimeTZ RangeSep Time
DateTimeTZRange:
  MONTH_NAME INT RangeSep INT
DateTimeTZRange:
  INT RangeSep INT MONTH_NAME
DateTimeTZRange:
  MONTH_NAME INT RangeSep INT YEAR
DateTimeTZRange:
  INT RangeSep INT MONTH_NAME YEAR
DateTimeTZRange:
  DateTimeTZ RangeSep DateTimeTZ
DateTimeTZRange:
  MONTH_NAME INT RangeSep MONTH_NAME INT YEAR
DateTimeTZ:
  Date
DateTimeTZ:
  Date DateTimeSep Time
Date:
  WEEKDAY_NAME CommaOpt Date
Date:
  Date T
Date:
  INT DateSep INT
Date:
  YEAR
Date:
  YEAR SUB INT
Date:
  YEAR SUB INT SUB INT
Date:
  MONTH_NAME INT YEAR
Date:
  INT MONTH_NAME YEAR
Date:
  INT DateSep INT DateSep YEAR
Date:
  MONTH_NAME INT
Date:
  INT MONTH_NAME
Time:
  INT
Time:
  INT AM
Time:
  INT PM
Time:
  INT COLON INT
Time:
  INT COLON INT COLON INT
Time:
  INT COLON INT AM
Time:
  INT COLON INT PM
DateSep:
  QUO
DateSep:
  PERIOD
DateTimeSep:
  <empty>
DateTimeSep:
  AT
DateTimeSep:
  SUB
RangeSep:
  SUB
RangeSep:
  THROUGH
RangeSep:
  TO
PrefixOpt:
  WhenOpt
SuffixOpt:
  GoogleOpt CalendarOpt ICSOpt
AndOpt:
  <empty>
AndOpt:
  AND
CalendarOpt:
  <empty>
CalendarOpt:
  CALENDAR
CommaOpt:
  <empty>
CommaOpt:
  COMMA
GoogleOpt:
  <empty>
GoogleOpt:
  GOOGLE
ICSOpt:
  <empty>
ICSOpt:
  ICS
WhenOpt:
  <empty>
WhenOpt:
  WHEN
*/

var parseRules = &glr.Rules{Items:[]glr.Rule{
  /*   0 */ glr.Rule{Nonterminal:"", RHS:[]string(nil), Type:""}, // ignored because rule-numbering starts at 1
  /*   1 */ glr.Rule{Nonterminal:"root", RHS:[]string{"PrefixOpt", "DateTimeTZRanges", "SuffixOpt"}, Type:"*DateTimeTZRanges"},
  /*   2 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"DateTimeTZRange"}, Type:"*DateTimeTZRanges"},
  /*   3 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"DateTimeTZRanges", "AndOpt", "DateTimeTZRange"}, Type:"*DateTimeTZRanges"},
  /*   4 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "INT"}, Type:"*DateTimeTZRanges"},
  /*   5 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "INT", "INT"}, Type:"*DateTimeTZRanges"},
  /*   6 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "INT", "INT", "INT"}, Type:"*DateTimeTZRanges"},
  /*   7 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "INT", "INT", "INT", "INT"}, Type:"*DateTimeTZRanges"},
  /*   8 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "INT", "MONTH_NAME"}, Type:"*DateTimeTZRanges"},
  /*   9 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "INT", "INT", "MONTH_NAME"}, Type:"*DateTimeTZRanges"},
  /*  10 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "INT", "INT", "INT", "MONTH_NAME"}, Type:"*DateTimeTZRanges"},
  /*  11 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "INT", "INT", "INT", "INT", "MONTH_NAME"}, Type:"*DateTimeTZRanges"},
  /*  12 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "INT", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  13 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "INT", "INT", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  14 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "INT", "INT", "INT", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  15 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "INT", "INT", "INT", "INT", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  16 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "INT", "MONTH_NAME", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  17 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "INT", "INT", "MONTH_NAME", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  18 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "INT", "INT", "INT", "MONTH_NAME", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  19 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "INT", "INT", "INT", "INT", "MONTH_NAME", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  20 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "RangeSep", "INT", "INT", "RangeSep", "INT"}, Type:"*DateTimeTZRanges"},
  /*  21 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "RangeSep", "INT", "INT", "RangeSep", "INT", "MONTH_NAME"}, Type:"*DateTimeTZRanges"},
  /*  22 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "RangeSep", "INT", "MONTH_NAME", "INT", "RangeSep", "INT"}, Type:"*DateTimeTZRanges"},
  /*  23 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "RangeSep", "INT", "MONTH_NAME", "INT", "RangeSep", "INT", "MONTH_NAME"}, Type:"*DateTimeTZRanges"},
  /*  24 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "RangeSep", "INT", "INT", "RangeSep", "INT", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  25 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "RangeSep", "INT", "INT", "RangeSep", "INT", "MONTH_NAME", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  26 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "RangeSep", "INT", "MONTH_NAME", "INT", "RangeSep", "INT", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  27 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "RangeSep", "INT", "MONTH_NAME", "INT", "RangeSep", "INT", "MONTH_NAME", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  28 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "MONTH_NAME", "INT", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  29 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"DateTimeTZ"}, Type:"*DateTimeTZRange"},
  /*  30 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"DateTimeTZ", "RangeSep", "Time"}, Type:"*DateTimeTZRange"},
  /*  31 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"MONTH_NAME", "INT", "RangeSep", "INT"}, Type:"*DateTimeTZRange"},
  /*  32 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"INT", "RangeSep", "INT", "MONTH_NAME"}, Type:"*DateTimeTZRange"},
  /*  33 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"MONTH_NAME", "INT", "RangeSep", "INT", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  34 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"INT", "RangeSep", "INT", "MONTH_NAME", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  35 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"DateTimeTZ", "RangeSep", "DateTimeTZ"}, Type:"*DateTimeTZRange"},
  /*  36 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"MONTH_NAME", "INT", "RangeSep", "MONTH_NAME", "INT", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  37 */ glr.Rule{Nonterminal:"DateTimeTZ", RHS:[]string{"Date"}, Type:"*DateTimeTZ"},
  /*  38 */ glr.Rule{Nonterminal:"DateTimeTZ", RHS:[]string{"Date", "DateTimeSep", "Time"}, Type:"*DateTimeTZ"},
  /*  39 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"WEEKDAY_NAME", "CommaOpt", "Date"}, Type:"civil.Date"},
  /*  40 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"Date", "T"}, Type:"civil.Date"},
  /*  41 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"INT", "DateSep", "INT"}, Type:"civil.Date"},
  /*  42 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"YEAR"}, Type:"civil.Date"},
  /*  43 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"YEAR", "SUB", "INT"}, Type:"civil.Date"},
  /*  44 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"YEAR", "SUB", "INT", "SUB", "INT"}, Type:"civil.Date"},
  /*  45 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"MONTH_NAME", "INT", "YEAR"}, Type:"civil.Date"},
  /*  46 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"INT", "MONTH_NAME", "YEAR"}, Type:"civil.Date"},
  /*  47 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"INT", "DateSep", "INT", "DateSep", "YEAR"}, Type:"civil.Date"},
  /*  48 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"MONTH_NAME", "INT"}, Type:"civil.Date"},
  /*  49 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"INT", "MONTH_NAME"}, Type:"civil.Date"},
  /*  50 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT"}, Type:"civil.Time"},
  /*  51 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "AM"}, Type:"civil.Time"},
  /*  52 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "PM"}, Type:"civil.Time"},
  /*  53 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "COLON", "INT"}, Type:"civil.Time"},
  /*  54 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "COLON", "INT", "COLON", "INT"}, Type:"civil.Time"},
  /*  55 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "COLON", "INT", "AM"}, Type:"civil.Time"},
  /*  56 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "COLON", "INT", "PM"}, Type:"civil.Time"},
  /*  57 */ glr.Rule{Nonterminal:"DateSep", RHS:[]string{"QUO"}, Type:""},
  /*  58 */ glr.Rule{Nonterminal:"DateSep", RHS:[]string{"PERIOD"}, Type:""},
  /*  59 */ glr.Rule{Nonterminal:"DateTimeSep", RHS:[]string(nil), Type:""},
  /*  60 */ glr.Rule{Nonterminal:"DateTimeSep", RHS:[]string{"AT"}, Type:""},
  /*  61 */ glr.Rule{Nonterminal:"DateTimeSep", RHS:[]string{"SUB"}, Type:""},
  /*  62 */ glr.Rule{Nonterminal:"RangeSep", RHS:[]string{"SUB"}, Type:""},
  /*  63 */ glr.Rule{Nonterminal:"RangeSep", RHS:[]string{"THROUGH"}, Type:""},
  /*  64 */ glr.Rule{Nonterminal:"RangeSep", RHS:[]string{"TO"}, Type:""},
  /*  65 */ glr.Rule{Nonterminal:"PrefixOpt", RHS:[]string{"WhenOpt"}, Type:""},
  /*  66 */ glr.Rule{Nonterminal:"SuffixOpt", RHS:[]string{"GoogleOpt", "CalendarOpt", "ICSOpt"}, Type:""},
  /*  67 */ glr.Rule{Nonterminal:"AndOpt", RHS:[]string(nil), Type:""},
  /*  68 */ glr.Rule{Nonterminal:"AndOpt", RHS:[]string{"AND"}, Type:""},
  /*  69 */ glr.Rule{Nonterminal:"CalendarOpt", RHS:[]string(nil), Type:""},
  /*  70 */ glr.Rule{Nonterminal:"CalendarOpt", RHS:[]string{"CALENDAR"}, Type:""},
  /*  71 */ glr.Rule{Nonterminal:"CommaOpt", RHS:[]string(nil), Type:""},
  /*  72 */ glr.Rule{Nonterminal:"CommaOpt", RHS:[]string{"COMMA"}, Type:""},
  /*  73 */ glr.Rule{Nonterminal:"GoogleOpt", RHS:[]string(nil), Type:""},
  /*  74 */ glr.Rule{Nonterminal:"GoogleOpt", RHS:[]string{"GOOGLE"}, Type:""},
  /*  75 */ glr.Rule{Nonterminal:"ICSOpt", RHS:[]string(nil), Type:""},
  /*  76 */ glr.Rule{Nonterminal:"ICSOpt", RHS:[]string{"ICS"}, Type:""},
  /*  77 */ glr.Rule{Nonterminal:"WhenOpt", RHS:[]string(nil), Type:""},
  /*  78 */ glr.Rule{Nonterminal:"WhenOpt", RHS:[]string{"WHEN"}, Type:""},
}}

// Semantic action functions

var parseActions = &glr.SemanticActions{Items:[]any{
  /*   0 */ nil, // empty action
  /*   1 */ func (PrefixOpt1 string, DateTimeTZRanges1 *DateTimeTZRanges, SuffixOpt1 string) *DateTimeTZRanges {return DateTimeTZRanges1},
  /*   2 */ func (DateTimeTZRange1 *DateTimeTZRange) *DateTimeTZRanges {return &DateTimeTZRanges{Items: []*DateTimeTZRange{DateTimeTZRange1}}},
  /*   3 */ func (DateTimeTZRanges1 *DateTimeTZRanges, AndOpt1 string, DateTimeTZRange1 *DateTimeTZRange) *DateTimeTZRanges {return AppendDateTimeTZRanges(DateTimeTZRanges1, DateTimeTZRange1)},
  /*   4 */ func (MONTH_NAME1 string, INT1 string, INT2 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, ""))},
  /*   5 */ func (MONTH_NAME1 string, INT1 string, INT2 string, INT3 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, ""), NewMDYDate(MONTH_NAME1, INT3, ""))},
  /*   6 */ func (MONTH_NAME1 string, INT1 string, INT2 string, INT3 string, INT4 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, ""), NewMDYDate(MONTH_NAME1, INT3, ""), NewMDYDate(MONTH_NAME1, INT4, ""))},
  /*   7 */ func (MONTH_NAME1 string, INT1 string, INT2 string, INT3 string, INT4 string, INT5 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, ""), NewMDYDate(MONTH_NAME1, INT3, ""), NewMDYDate(MONTH_NAME1, INT4, ""), NewMDYDate(MONTH_NAME1, INT5, ""))},
  /*   8 */ func (INT1 string, INT2 string, MONTH_NAME1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, ""))},
  /*   9 */ func (INT1 string, INT2 string, INT3 string, MONTH_NAME1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, ""), NewDMYDate(INT3, MONTH_NAME1, ""))},
  /*  10 */ func (INT1 string, INT2 string, INT3 string, INT4 string, MONTH_NAME1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, ""), NewDMYDate(INT3, MONTH_NAME1, ""), NewDMYDate(INT4, MONTH_NAME1, ""))},
  /*  11 */ func (INT1 string, INT2 string, INT3 string, INT4 string, INT5 string, MONTH_NAME1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, ""), NewDMYDate(INT3, MONTH_NAME1, ""), NewDMYDate(INT4, MONTH_NAME1, ""), NewDMYDate(INT5, MONTH_NAME1, ""))},
  /*  12 */ func (MONTH_NAME1 string, INT1 string, INT2 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1))},
  /*  13 */ func (MONTH_NAME1 string, INT1 string, INT2 string, INT3 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1), NewMDYDate(MONTH_NAME1, INT3, YEAR1))},
  /*  14 */ func (MONTH_NAME1 string, INT1 string, INT2 string, INT3 string, INT4 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1), NewMDYDate(MONTH_NAME1, INT3, YEAR1), NewMDYDate(MONTH_NAME1, INT4, YEAR1))},
  /*  15 */ func (MONTH_NAME1 string, INT1 string, INT2 string, INT3 string, INT4 string, INT5 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1), NewMDYDate(MONTH_NAME1, INT3, YEAR1), NewMDYDate(MONTH_NAME1, INT4, YEAR1), NewMDYDate(MONTH_NAME1, INT5, YEAR1))},
  /*  16 */ func (INT1 string, INT2 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1))},
  /*  17 */ func (INT1 string, INT2 string, INT3 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1), NewDMYDate(INT3, MONTH_NAME1, YEAR1))},
  /*  18 */ func (INT1 string, INT2 string, INT3 string, INT4 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1), NewDMYDate(INT3, MONTH_NAME1, YEAR1), NewDMYDate(INT4, MONTH_NAME1, YEAR1))},
  /*  19 */ func (INT1 string, INT2 string, INT3 string, INT4 string, INT5 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStartDates(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1), NewDMYDate(INT3, MONTH_NAME1, YEAR1), NewDMYDate(INT4, MONTH_NAME1, YEAR1), NewDMYDate(INT5, MONTH_NAME1, YEAR1))},
  /*  20 */ func (MONTH_NAME1 string, INT1 string, RangeSep1 string, INT2 string, INT3 string, RangeSep2 string, INT4 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, "")), NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT3, ""), NewMDYDate(MONTH_NAME1, INT4, "")))},
  /*  21 */ func (INT1 string, RangeSep1 string, INT2 string, INT3 string, RangeSep2 string, INT4 string, MONTH_NAME1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, "")), NewRangeWithStartEndDates(NewDMYDate(INT3, MONTH_NAME1, ""), NewDMYDate(INT4, MONTH_NAME1, "")))},
  /*  22 */ func (MONTH_NAME1 string, INT1 string, RangeSep1 string, INT2 string, MONTH_NAME2 string, INT3 string, RangeSep2 string, INT4 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, "")), NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME2, INT3, ""), NewMDYDate(MONTH_NAME2, INT4, "")))},
  /*  23 */ func (INT1 string, RangeSep1 string, INT2 string, MONTH_NAME1 string, INT3 string, RangeSep2 string, INT4 string, MONTH_NAME2 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, "")), NewRangeWithStartEndDates(NewDMYDate(INT3, MONTH_NAME2, ""), NewDMYDate(INT4, MONTH_NAME2, "")))},
  /*  24 */ func (MONTH_NAME1 string, INT1 string, RangeSep1 string, INT2 string, INT3 string, RangeSep2 string, INT4 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1)), NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT3, YEAR1), NewMDYDate(MONTH_NAME1, INT4, YEAR1)))},
  /*  25 */ func (INT1 string, RangeSep1 string, INT2 string, INT3 string, RangeSep2 string, INT4 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1)), NewRangeWithStartEndDates(NewDMYDate(INT3, MONTH_NAME1, YEAR1), NewDMYDate(INT4, MONTH_NAME1, YEAR1)))},
  /*  26 */ func (MONTH_NAME1 string, INT1 string, RangeSep1 string, INT2 string, MONTH_NAME2 string, INT3 string, RangeSep2 string, INT4 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1)), NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME2, INT3, YEAR1), NewMDYDate(MONTH_NAME2, INT4, YEAR1)))},
  /*  27 */ func (INT1 string, RangeSep1 string, INT2 string, MONTH_NAME1 string, INT3 string, RangeSep2 string, INT4 string, MONTH_NAME2 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1)), NewRangeWithStartEndDates(NewDMYDate(INT3, MONTH_NAME2, YEAR1), NewDMYDate(INT4, MONTH_NAME2, YEAR1)))},
  /*  28 */ func (MONTH_NAME1 string, INT1 string, MONTH_NAME2 string, INT2 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStart(NewMDYDate(MONTH_NAME1, INT1, YEAR1)), NewRangeWithStart(NewMDYDate(MONTH_NAME2, INT2, YEAR1)))},
  /*  29 */ func (DateTimeTZ1 *DateTimeTZ) *DateTimeTZRange {return &DateTimeTZRange{Start: DateTimeTZ1}},
  /*  30 */ func (DateTimeTZ1 *DateTimeTZ, RangeSep1 string, Time1 civil.Time) *DateTimeTZRange {return NewRangeWithStartEndDateTimes(DateTimeTZ1, NewDateTime(DateTimeTZ1.Date, Time1, ""))},
  /*  31 */ func (MONTH_NAME1 string, INT1 string, RangeSep1 string, INT2 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, ""))},
  /*  32 */ func (INT1 string, RangeSep1 string, INT2 string, MONTH_NAME1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, ""))},
  /*  33 */ func (MONTH_NAME1 string, INT1 string, RangeSep1 string, INT2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1))},
  /*  34 */ func (INT1 string, RangeSep1 string, INT2 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1))},
  /*  35 */ func (DateTimeTZ1 *DateTimeTZ, RangeSep1 string, DateTimeTZ2 *DateTimeTZ) *DateTimeTZRange {return &DateTimeTZRange{Start: DateTimeTZ1, End: DateTimeTZ2}},
  /*  36 */ func (MONTH_NAME1 string, INT1 string, RangeSep1 string, MONTH_NAME2 string, INT2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME2, INT2, YEAR1))},
  /*  37 */ func (Date1 civil.Date) *DateTimeTZ {return NewDateTimeWithDate(Date1)},
  /*  38 */ func (Date1 civil.Date, DateTimeSep1 string, Time1 civil.Time) *DateTimeTZ {return NewDateTime(Date1, Time1, "")},
  /*  39 */ func (WEEKDAY_NAME1 string, CommaOpt1 string, Date1 civil.Date) civil.Date {return Date1},
  /*  40 */ func (Date1 civil.Date, T1 string) civil.Date {return Date1},
  /*  41 */ func (INT1 string, DateSep1 string, INT2 string) civil.Date {return NewAmbiguousDate(INT1, INT2, "")},
  /*  42 */ func (YEAR1 string) civil.Date {return NewDMYDate("", "", YEAR1)},
  /*  43 */ func (YEAR1 string, SUB1 string, INT1 string) civil.Date {return NewDMYDate("", INT1, YEAR1)},
  /*  44 */ func (YEAR1 string, SUB1 string, INT1 string, SUB2 string, INT2 string) civil.Date {return NewDMYDate(INT2, INT1, YEAR1)},
  /*  45 */ func (MONTH_NAME1 string, INT1 string, YEAR1 string) civil.Date {return NewMDYDate(MONTH_NAME1, INT1, YEAR1)},
  /*  46 */ func (INT1 string, MONTH_NAME1 string, YEAR1 string) civil.Date {return NewDMYDate(INT1, MONTH_NAME1, YEAR1)},
  /*  47 */ func (INT1 string, DateSep1 string, INT2 string, DateSep2 string, YEAR1 string) civil.Date {return NewAmbiguousDate(INT1, INT2, YEAR1)},
  /*  48 */ func (MONTH_NAME1 string, INT1 string) civil.Date {return NewMDYDate(MONTH_NAME1, INT1, "")},
  /*  49 */ func (INT1 string, MONTH_NAME1 string) civil.Date {return NewDMYDate(INT1, MONTH_NAME1, "")},
  /*  50 */ func (INT1 string) civil.Time {return NewTime(INT1, "", "", "")},
  /*  51 */ func (INT1 string, AM1 string) civil.Time {return NewTime(INT1, "", "", "")},
  /*  52 */ func (INT1 string, PM1 string) civil.Time {return NewTime((mustAtoi(INT1) % 12) + 12, "", "", "")},
  /*  53 */ func (INT1 string, COLON1 string, INT2 string) civil.Time {return NewTime(INT1, INT2, "", "")},
  /*  54 */ func (INT1 string, COLON1 string, INT2 string, COLON2 string, INT3 string) civil.Time {return NewTime((mustAtoi(INT1) % 12) + 12, INT2, INT3, "")},
  /*  55 */ func (INT1 string, COLON1 string, INT2 string, AM1 string) civil.Time {return NewTime(INT1, INT2, "", "")},
  /*  56 */ func (INT1 string, COLON1 string, INT2 string, PM1 string) civil.Time {return NewTime((mustAtoi(INT1) % 12) + 12, INT2, "", "")},
  /*  57 */ func (QUO1 string) string {return QUO1},
  /*  58 */ func (PERIOD1 string) string {return PERIOD1},
  /*  59 */ func () string {return ""},
  /*  60 */ func (AT1 string) string {return AT1},
  /*  61 */ func (SUB1 string) string {return SUB1},
  /*  62 */ func (SUB1 string) string {return SUB1},
  /*  63 */ func (THROUGH1 string) string {return THROUGH1},
  /*  64 */ func (TO1 string) string {return TO1},
  /*  65 */ func (WhenOpt1 string) string {return WhenOpt1},
  /*  66 */ func (GoogleOpt1 string, CalendarOpt1 string, ICSOpt1 string) string {return GoogleOpt1},
  /*  67 */ func () string {return ""},
  /*  68 */ func (AND1 string) string {return AND1},
  /*  69 */ func () string {return ""},
  /*  70 */ func (CALENDAR1 string) string {return CALENDAR1},
  /*  71 */ func () string {return ""},
  /*  72 */ func (COMMA1 string) string {return COMMA1},
  /*  73 */ func () string {return ""},
  /*  74 */ func (GOOGLE1 string) string {return GOOGLE1},
  /*  75 */ func () string {return ""},
  /*  76 */ func (ICS1 string) string {return ICS1},
  /*  77 */ func () string {return ""},
  /*  78 */ func (WHEN1 string) string {return WHEN1},
}}

var parseStates = &glr.ParseStates{Items:[]glr.ParseState{
  /*   0 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:77}}, "WHEN":[]glr.Action{glr.Action{Type:"shift", State:4, Rule:0}}}, Gotos:map[string]int{"PrefixOpt":2, "WhenOpt":3, "root":1}},
  /*   1 */ glr.ParseState{Actions:map[string][]glr.Action{"$end":[]glr.Action{glr.Action{Type:"accept", State:0, Rule:0}}}, Gotos:map[string]int{}},
  /*   2 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:8, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:7, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:11, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:12, Rule:0}}}, Gotos:map[string]int{"Date":10, "DateTimeTZ":9, "DateTimeTZRange":6, "DateTimeTZRanges":5}},
  /*   3 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:65}}}, Gotos:map[string]int{}},
  /*   4 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:78}}}, Gotos:map[string]int{}},
  /*   5 */ glr.ParseState{Actions:map[string][]glr.Action{"$end":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:73}}, ".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:67}}, "AND":[]glr.Action{glr.Action{Type:"shift", State:16, Rule:0}}, "CALENDAR":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:73}}, "GOOGLE":[]glr.Action{glr.Action{Type:"shift", State:17, Rule:0}}, "ICS":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:73}}}, Gotos:map[string]int{"AndOpt":14, "GoogleOpt":15, "SuffixOpt":13}},
  /*   6 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:2}}}, Gotos:map[string]int{}},
  /*   7 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:18, Rule:0}}}, Gotos:map[string]int{}},
  /*   8 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:22, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:27, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}}, Gotos:map[string]int{"DateSep":21, "RangeSep":20}},
  /*   9 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:29}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}}, Gotos:map[string]int{"RangeSep":28}},
  /*  10 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:37}, glr.Action{Type:"reduce", State:0, Rule:59}, glr.Action{Type:"reduce", State:0, Rule:37}}, "AT":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:32, Rule:0}}, "T":[]glr.Action{glr.Action{Type:"shift", State:30, Rule:0}}}, Gotos:map[string]int{"DateTimeSep":29}},
  /*  11 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:71}}, "COMMA":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}}, Gotos:map[string]int{"CommaOpt":33}},
  /*  12 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:42}, glr.Action{Type:"reduce", State:0, Rule:42}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:35, Rule:0}}}, Gotos:map[string]int{}},
  /*  13 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:1}}}, Gotos:map[string]int{}},
  /*  14 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:38, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:37, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:11, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:12, Rule:0}}}, Gotos:map[string]int{"Date":10, "DateTimeTZ":9, "DateTimeTZRange":36}},
  /*  15 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:69}}, "CALENDAR":[]glr.Action{glr.Action{Type:"shift", State:40, Rule:0}}}, Gotos:map[string]int{"CalendarOpt":39}},
  /*  16 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:68}}}, Gotos:map[string]int{}},
  /*  17 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:74}}}, Gotos:map[string]int{}},
  /*  18 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:48}, glr.Action{Type:"reduce", State:0, Rule:48}, glr.Action{Type:"reduce", State:0, Rule:48}, glr.Action{Type:"reduce", State:0, Rule:48}, glr.Action{Type:"reduce", State:0, Rule:48}, glr.Action{Type:"reduce", State:0, Rule:48}, glr.Action{Type:"reduce", State:0, Rule:48}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:41, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:43, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"RangeSep":42}},
  /*  19 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:46, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:45, Rule:0}}}, Gotos:map[string]int{}},
  /*  20 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:47, Rule:0}}}, Gotos:map[string]int{}},
  /*  21 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:48, Rule:0}}}, Gotos:map[string]int{}},
  /*  22 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:49}, glr.Action{Type:"reduce", State:0, Rule:49}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:49, Rule:0}}}, Gotos:map[string]int{}},
  /*  23 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:62}}}, Gotos:map[string]int{}},
  /*  24 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:63}}}, Gotos:map[string]int{}},
  /*  25 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:64}}}, Gotos:map[string]int{}},
  /*  26 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:57}}}, Gotos:map[string]int{}},
  /*  27 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:58}}}, Gotos:map[string]int{}},
  /*  28 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:52, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:53, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:11, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:12, Rule:0}}}, Gotos:map[string]int{"Date":10, "DateTimeTZ":51, "Time":50}},
  /*  29 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}}, Gotos:map[string]int{"Time":54}},
  /*  30 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:40}}}, Gotos:map[string]int{}},
  /*  31 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:60}}}, Gotos:map[string]int{}},
  /*  32 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:61}}}, Gotos:map[string]int{}},
  /*  33 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:53, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:11, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:12, Rule:0}}}, Gotos:map[string]int{"Date":56}},
  /*  34 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:72}}}, Gotos:map[string]int{}},
  /*  35 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}}, Gotos:map[string]int{}},
  /*  36 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:3}}}, Gotos:map[string]int{}},
  /*  37 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:59, Rule:0}}}, Gotos:map[string]int{}},
  /*  38 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:22, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:27, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}}, Gotos:map[string]int{"DateSep":21, "RangeSep":60}},
  /*  39 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:75}}, "ICS":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}}, Gotos:map[string]int{"ICSOpt":61}},
  /*  40 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:70}}}, Gotos:map[string]int{}},
  /*  41 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:4}, glr.Action{Type:"reduce", State:0, Rule:4}, glr.Action{Type:"reduce", State:0, Rule:4}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:63, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:64, Rule:0}}}, Gotos:map[string]int{}},
  /*  42 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:65, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:66, Rule:0}}}, Gotos:map[string]int{}},
  /*  43 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:67, Rule:0}}}, Gotos:map[string]int{}},
  /*  44 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:45}}}, Gotos:map[string]int{}},
  /*  45 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:8}, glr.Action{Type:"reduce", State:0, Rule:8}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:68, Rule:0}}}, Gotos:map[string]int{}},
  /*  46 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}}, Gotos:map[string]int{}},
  /*  47 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:71, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:72, Rule:0}}}, Gotos:map[string]int{}},
  /*  48 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:41}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:27, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}}, Gotos:map[string]int{"DateSep":73}},
  /*  49 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:46}}}, Gotos:map[string]int{}},
  /*  50 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:30}}}, Gotos:map[string]int{}},
  /*  51 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:35}}}, Gotos:map[string]int{}},
  /*  52 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:50}, glr.Action{Type:"reduce", State:0, Rule:50}}, "AM":[]glr.Action{glr.Action{Type:"shift", State:74, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:76, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:22, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:27, Rule:0}}, "PM":[]glr.Action{glr.Action{Type:"shift", State:75, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}}, Gotos:map[string]int{"DateSep":21}},
  /*  53 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:77, Rule:0}}}, Gotos:map[string]int{}},
  /*  54 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:38}}}, Gotos:map[string]int{}},
  /*  55 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:50}}, "AM":[]glr.Action{glr.Action{Type:"shift", State:74, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:76, Rule:0}}, "PM":[]glr.Action{glr.Action{Type:"shift", State:75, Rule:0}}}, Gotos:map[string]int{}},
  /*  56 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:39}, glr.Action{Type:"reduce", State:0, Rule:39}}, "T":[]glr.Action{glr.Action{Type:"shift", State:30, Rule:0}}}, Gotos:map[string]int{}},
  /*  57 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:22, Rule:0}}, "PERIOD":[]glr.Action{glr.Action{Type:"shift", State:27, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}}, Gotos:map[string]int{"DateSep":21}},
  /*  58 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:43}, glr.Action{Type:"reduce", State:0, Rule:43}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:78, Rule:0}}}, Gotos:map[string]int{}},
  /*  59 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:48}, glr.Action{Type:"reduce", State:0, Rule:48}, glr.Action{Type:"reduce", State:0, Rule:48}, glr.Action{Type:"reduce", State:0, Rule:48}, glr.Action{Type:"reduce", State:0, Rule:48}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{"RangeSep":79}},
  /*  60 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:80, Rule:0}}}, Gotos:map[string]int{}},
  /*  61 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:66}}}, Gotos:map[string]int{}},
  /*  62 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:76}}}, Gotos:map[string]int{}},
  /*  63 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:5}, glr.Action{Type:"reduce", State:0, Rule:5}, glr.Action{Type:"reduce", State:0, Rule:5}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:81, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:82, Rule:0}}}, Gotos:map[string]int{}},
  /*  64 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:12}}}, Gotos:map[string]int{}},
  /*  65 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:31}, glr.Action{Type:"reduce", State:0, Rule:31}, glr.Action{Type:"reduce", State:0, Rule:31}, glr.Action{Type:"reduce", State:0, Rule:31}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:83, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:84, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:85, Rule:0}}}, Gotos:map[string]int{}},
  /*  66 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:86, Rule:0}}}, Gotos:map[string]int{}},
  /*  67 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:87, Rule:0}}}, Gotos:map[string]int{}},
  /*  68 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:16}}}, Gotos:map[string]int{}},
  /*  69 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:9}, glr.Action{Type:"reduce", State:0, Rule:9}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:88, Rule:0}}}, Gotos:map[string]int{}},
  /*  70 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:90, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:89, Rule:0}}}, Gotos:map[string]int{}},
  /*  71 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}}, Gotos:map[string]int{"RangeSep":91}},
  /*  72 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:32}, glr.Action{Type:"reduce", State:0, Rule:32}, glr.Action{Type:"reduce", State:0, Rule:32}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:92, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:93, Rule:0}}}, Gotos:map[string]int{}},
  /*  73 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:94, Rule:0}}}, Gotos:map[string]int{}},
  /*  74 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:51}}}, Gotos:map[string]int{}},
  /*  75 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:52}}}, Gotos:map[string]int{}},
  /*  76 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:95, Rule:0}}}, Gotos:map[string]int{}},
  /*  77 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:48}, glr.Action{Type:"reduce", State:0, Rule:48}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{}},
  /*  78 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:96, Rule:0}}}, Gotos:map[string]int{}},
  /*  79 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:97, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:66, Rule:0}}}, Gotos:map[string]int{}},
  /*  80 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:98, Rule:0}}}, Gotos:map[string]int{}},
  /*  81 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:6}, glr.Action{Type:"reduce", State:0, Rule:6}, glr.Action{Type:"reduce", State:0, Rule:6}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:99, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:100, Rule:0}}}, Gotos:map[string]int{}},
  /*  82 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:13}}}, Gotos:map[string]int{}},
  /*  83 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}}, Gotos:map[string]int{"RangeSep":101}},
  /*  84 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:102, Rule:0}}}, Gotos:map[string]int{}},
  /*  85 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:33}}}, Gotos:map[string]int{}},
  /*  86 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:103, Rule:0}}}, Gotos:map[string]int{}},
  /*  87 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:28}}}, Gotos:map[string]int{}},
  /*  88 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:17}}}, Gotos:map[string]int{}},
  /*  89 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:10}, glr.Action{Type:"reduce", State:0, Rule:10}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:104, Rule:0}}}, Gotos:map[string]int{}},
  /*  90 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:105, Rule:0}}}, Gotos:map[string]int{}},
  /*  91 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:106, Rule:0}}}, Gotos:map[string]int{}},
  /*  92 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}}, Gotos:map[string]int{"RangeSep":107}},
  /*  93 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:34}}}, Gotos:map[string]int{}},
  /*  94 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:47}}}, Gotos:map[string]int{}},
  /*  95 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:53}}, "AM":[]glr.Action{glr.Action{Type:"shift", State:109, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:108, Rule:0}}, "PM":[]glr.Action{glr.Action{Type:"shift", State:110, Rule:0}}}, Gotos:map[string]int{}},
  /*  96 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:44}}}, Gotos:map[string]int{}},
  /*  97 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:31}, glr.Action{Type:"reduce", State:0, Rule:31}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:85, Rule:0}}}, Gotos:map[string]int{}},
  /*  98 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:32}, glr.Action{Type:"reduce", State:0, Rule:32}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:93, Rule:0}}}, Gotos:map[string]int{}},
  /*  99 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:7}, glr.Action{Type:"reduce", State:0, Rule:7}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:111, Rule:0}}}, Gotos:map[string]int{}},
  /* 100 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:14}}}, Gotos:map[string]int{}},
  /* 101 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:112, Rule:0}}}, Gotos:map[string]int{}},
  /* 102 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}}, Gotos:map[string]int{"RangeSep":113}},
  /* 103 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:36}}}, Gotos:map[string]int{}},
  /* 104 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:18}}}, Gotos:map[string]int{}},
  /* 105 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:11}, glr.Action{Type:"reduce", State:0, Rule:11}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:114, Rule:0}}}, Gotos:map[string]int{}},
  /* 106 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:115, Rule:0}}}, Gotos:map[string]int{}},
  /* 107 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:116, Rule:0}}}, Gotos:map[string]int{}},
  /* 108 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:117, Rule:0}}}, Gotos:map[string]int{}},
  /* 109 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:55}}}, Gotos:map[string]int{}},
  /* 110 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:56}}}, Gotos:map[string]int{}},
  /* 111 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:15}}}, Gotos:map[string]int{}},
  /* 112 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:20}, glr.Action{Type:"reduce", State:0, Rule:20}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:118, Rule:0}}}, Gotos:map[string]int{}},
  /* 113 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:119, Rule:0}}}, Gotos:map[string]int{}},
  /* 114 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:19}}}, Gotos:map[string]int{}},
  /* 115 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:21}, glr.Action{Type:"reduce", State:0, Rule:21}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:120, Rule:0}}}, Gotos:map[string]int{}},
  /* 116 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:121, Rule:0}}}, Gotos:map[string]int{}},
  /* 117 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:54}}}, Gotos:map[string]int{}},
  /* 118 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:24}}}, Gotos:map[string]int{}},
  /* 119 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:22}, glr.Action{Type:"reduce", State:0, Rule:22}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:122, Rule:0}}}, Gotos:map[string]int{}},
  /* 120 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:25}}}, Gotos:map[string]int{}},
  /* 121 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:23}, glr.Action{Type:"reduce", State:0, Rule:23}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:123, Rule:0}}}, Gotos:map[string]int{}},
  /* 122 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:26}}}, Gotos:map[string]int{}},
  /* 123 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:27}}}, Gotos:map[string]int{}},
}}

