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
  MONTH_NAME INT Sep INT INT Sep INT
DateTimeTZRanges:
  INT Sep INT INT Sep INT MONTH_NAME
DateTimeTZRanges:
  MONTH_NAME INT Sep INT MONTH_NAME INT Sep INT
DateTimeTZRanges:
  INT Sep INT MONTH_NAME INT Sep INT MONTH_NAME
DateTimeTZRanges:
  MONTH_NAME INT Sep INT INT Sep INT YEAR
DateTimeTZRanges:
  INT Sep INT INT Sep INT MONTH_NAME YEAR
DateTimeTZRanges:
  MONTH_NAME INT Sep INT MONTH_NAME INT Sep INT YEAR
DateTimeTZRanges:
  INT Sep INT MONTH_NAME INT Sep INT MONTH_NAME YEAR
DateTimeTZRanges:
  MONTH_NAME INT MONTH_NAME INT YEAR
DateTimeTZRange:
  DateTimeTZ
DateTimeTZRange:
  DateTimeTZ Sep Time
DateTimeTZRange:
  MONTH_NAME INT Sep INT
DateTimeTZRange:
  INT Sep INT MONTH_NAME
DateTimeTZRange:
  MONTH_NAME INT Sep INT YEAR
DateTimeTZRange:
  INT Sep INT MONTH_NAME YEAR
DateTimeTZRange:
  DateTimeTZ Sep DateTimeTZ
DateTimeTZRange:
  MONTH_NAME INT Sep MONTH_NAME INT YEAR
DateTimeTZ:
  Date
DateTimeTZ:
  Date AtOpt Time
Date:
  YEAR
Date:
  YEAR SUB INT
Date:
  YEAR SUB INT SUB INT TOpt
Date:
  MONTH_NAME INT YEAR
Date:
  INT MONTH_NAME YEAR
Date:
  INT QUO INT QUO YEAR
Date:
  WeekDayNameOpt MONTH_NAME INT
Date:
  WeekDayNameOpt INT MONTH_NAME
TOpt:
  <empty>
TOpt:
  T
WeekDayNameOpt:
  <empty>
WeekDayNameOpt:
  WEEKDAY_NAME
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
AndOpt:
  <empty>
AndOpt:
  AND
AtOpt:
  <empty>
AtOpt:
  AT
Sep:
  SUB
Sep:
  THROUGH
Sep:
  TO
PrefixOpt:
  WhenOpt
WhenOpt:
  <empty>
WhenOpt:
  WHEN
SuffixOpt:
  GoogleOpt CalendarOpt ICSOpt
GoogleOpt:
  <empty>
GoogleOpt:
  GOOGLE
CalendarOpt:
  <empty>
CalendarOpt:
  CALENDAR
ICSOpt:
  <empty>
ICSOpt:
  ICS
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
  /*  20 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "Sep", "INT", "INT", "Sep", "INT"}, Type:"*DateTimeTZRanges"},
  /*  21 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "Sep", "INT", "INT", "Sep", "INT", "MONTH_NAME"}, Type:"*DateTimeTZRanges"},
  /*  22 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "Sep", "INT", "MONTH_NAME", "INT", "Sep", "INT"}, Type:"*DateTimeTZRanges"},
  /*  23 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "Sep", "INT", "MONTH_NAME", "INT", "Sep", "INT", "MONTH_NAME"}, Type:"*DateTimeTZRanges"},
  /*  24 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "Sep", "INT", "INT", "Sep", "INT", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  25 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "Sep", "INT", "INT", "Sep", "INT", "MONTH_NAME", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  26 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "Sep", "INT", "MONTH_NAME", "INT", "Sep", "INT", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  27 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"INT", "Sep", "INT", "MONTH_NAME", "INT", "Sep", "INT", "MONTH_NAME", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  28 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"MONTH_NAME", "INT", "MONTH_NAME", "INT", "YEAR"}, Type:"*DateTimeTZRanges"},
  /*  29 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"DateTimeTZ"}, Type:"*DateTimeTZRange"},
  /*  30 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"DateTimeTZ", "Sep", "Time"}, Type:"*DateTimeTZRange"},
  /*  31 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"MONTH_NAME", "INT", "Sep", "INT"}, Type:"*DateTimeTZRange"},
  /*  32 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"INT", "Sep", "INT", "MONTH_NAME"}, Type:"*DateTimeTZRange"},
  /*  33 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"MONTH_NAME", "INT", "Sep", "INT", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  34 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"INT", "Sep", "INT", "MONTH_NAME", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  35 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"DateTimeTZ", "Sep", "DateTimeTZ"}, Type:"*DateTimeTZRange"},
  /*  36 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"MONTH_NAME", "INT", "Sep", "MONTH_NAME", "INT", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  37 */ glr.Rule{Nonterminal:"DateTimeTZ", RHS:[]string{"Date"}, Type:"*DateTimeTZ"},
  /*  38 */ glr.Rule{Nonterminal:"DateTimeTZ", RHS:[]string{"Date", "AtOpt", "Time"}, Type:"*DateTimeTZ"},
  /*  39 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"YEAR"}, Type:"civil.Date"},
  /*  40 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"YEAR", "SUB", "INT"}, Type:"civil.Date"},
  /*  41 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"YEAR", "SUB", "INT", "SUB", "INT", "TOpt"}, Type:"civil.Date"},
  /*  42 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"MONTH_NAME", "INT", "YEAR"}, Type:"civil.Date"},
  /*  43 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"INT", "MONTH_NAME", "YEAR"}, Type:"civil.Date"},
  /*  44 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"INT", "QUO", "INT", "QUO", "YEAR"}, Type:"civil.Date"},
  /*  45 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"WeekDayNameOpt", "MONTH_NAME", "INT"}, Type:"civil.Date"},
  /*  46 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"WeekDayNameOpt", "INT", "MONTH_NAME"}, Type:"civil.Date"},
  /*  47 */ glr.Rule{Nonterminal:"TOpt", RHS:[]string(nil), Type:""},
  /*  48 */ glr.Rule{Nonterminal:"TOpt", RHS:[]string{"T"}, Type:""},
  /*  49 */ glr.Rule{Nonterminal:"WeekDayNameOpt", RHS:[]string(nil), Type:""},
  /*  50 */ glr.Rule{Nonterminal:"WeekDayNameOpt", RHS:[]string{"WEEKDAY_NAME"}, Type:""},
  /*  51 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT"}, Type:"civil.Time"},
  /*  52 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "AM"}, Type:"civil.Time"},
  /*  53 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "PM"}, Type:"civil.Time"},
  /*  54 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "COLON", "INT"}, Type:"civil.Time"},
  /*  55 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "COLON", "INT", "COLON", "INT"}, Type:"civil.Time"},
  /*  56 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "COLON", "INT", "AM"}, Type:"civil.Time"},
  /*  57 */ glr.Rule{Nonterminal:"Time", RHS:[]string{"INT", "COLON", "INT", "PM"}, Type:"civil.Time"},
  /*  58 */ glr.Rule{Nonterminal:"AndOpt", RHS:[]string(nil), Type:""},
  /*  59 */ glr.Rule{Nonterminal:"AndOpt", RHS:[]string{"AND"}, Type:""},
  /*  60 */ glr.Rule{Nonterminal:"AtOpt", RHS:[]string(nil), Type:""},
  /*  61 */ glr.Rule{Nonterminal:"AtOpt", RHS:[]string{"AT"}, Type:""},
  /*  62 */ glr.Rule{Nonterminal:"Sep", RHS:[]string{"SUB"}, Type:""},
  /*  63 */ glr.Rule{Nonterminal:"Sep", RHS:[]string{"THROUGH"}, Type:""},
  /*  64 */ glr.Rule{Nonterminal:"Sep", RHS:[]string{"TO"}, Type:""},
  /*  65 */ glr.Rule{Nonterminal:"PrefixOpt", RHS:[]string{"WhenOpt"}, Type:""},
  /*  66 */ glr.Rule{Nonterminal:"WhenOpt", RHS:[]string(nil), Type:""},
  /*  67 */ glr.Rule{Nonterminal:"WhenOpt", RHS:[]string{"WHEN"}, Type:""},
  /*  68 */ glr.Rule{Nonterminal:"SuffixOpt", RHS:[]string{"GoogleOpt", "CalendarOpt", "ICSOpt"}, Type:""},
  /*  69 */ glr.Rule{Nonterminal:"GoogleOpt", RHS:[]string(nil), Type:""},
  /*  70 */ glr.Rule{Nonterminal:"GoogleOpt", RHS:[]string{"GOOGLE"}, Type:""},
  /*  71 */ glr.Rule{Nonterminal:"CalendarOpt", RHS:[]string(nil), Type:""},
  /*  72 */ glr.Rule{Nonterminal:"CalendarOpt", RHS:[]string{"CALENDAR"}, Type:""},
  /*  73 */ glr.Rule{Nonterminal:"ICSOpt", RHS:[]string(nil), Type:""},
  /*  74 */ glr.Rule{Nonterminal:"ICSOpt", RHS:[]string{"ICS"}, Type:""},
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
  /*  20 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, INT2 string, INT3 string, Sep2 string, INT4 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, "")), NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT3, ""), NewMDYDate(MONTH_NAME1, INT4, "")))},
  /*  21 */ func (INT1 string, Sep1 string, INT2 string, INT3 string, Sep2 string, INT4 string, MONTH_NAME1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, "")), NewRangeWithStartEndDates(NewDMYDate(INT3, MONTH_NAME1, ""), NewDMYDate(INT4, MONTH_NAME1, "")))},
  /*  22 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, INT2 string, MONTH_NAME2 string, INT3 string, Sep2 string, INT4 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, "")), NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME2, INT3, ""), NewMDYDate(MONTH_NAME2, INT4, "")))},
  /*  23 */ func (INT1 string, Sep1 string, INT2 string, MONTH_NAME1 string, INT3 string, Sep2 string, INT4 string, MONTH_NAME2 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, "")), NewRangeWithStartEndDates(NewDMYDate(INT3, MONTH_NAME2, ""), NewDMYDate(INT4, MONTH_NAME2, "")))},
  /*  24 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, INT2 string, INT3 string, Sep2 string, INT4 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1)), NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT3, YEAR1), NewMDYDate(MONTH_NAME1, INT4, YEAR1)))},
  /*  25 */ func (INT1 string, Sep1 string, INT2 string, INT3 string, Sep2 string, INT4 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1)), NewRangeWithStartEndDates(NewDMYDate(INT3, MONTH_NAME1, YEAR1), NewDMYDate(INT4, MONTH_NAME1, YEAR1)))},
  /*  26 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, INT2 string, MONTH_NAME2 string, INT3 string, Sep2 string, INT4 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1)), NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME2, INT3, YEAR1), NewMDYDate(MONTH_NAME2, INT4, YEAR1)))},
  /*  27 */ func (INT1 string, Sep1 string, INT2 string, MONTH_NAME1 string, INT3 string, Sep2 string, INT4 string, MONTH_NAME2 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEndDates(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1)), NewRangeWithStartEndDates(NewDMYDate(INT3, MONTH_NAME2, YEAR1), NewDMYDate(INT4, MONTH_NAME2, YEAR1)))},
  /*  28 */ func (MONTH_NAME1 string, INT1 string, MONTH_NAME2 string, INT2 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStart(NewMDYDate(MONTH_NAME1, INT1, YEAR1)), NewRangeWithStart(NewMDYDate(MONTH_NAME2, INT2, YEAR1)))},
  /*  29 */ func (DateTimeTZ1 *DateTimeTZ) *DateTimeTZRange {return &DateTimeTZRange{Start: DateTimeTZ1}},
  /*  30 */ func (DateTimeTZ1 *DateTimeTZ, Sep1 string, Time1 civil.Time) *DateTimeTZRange {return NewRangeWithStartEndDateTimes(DateTimeTZ1, NewDateTime(DateTimeTZ1.Date, Time1, ""))},
  /*  31 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, INT2 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, ""))},
  /*  32 */ func (INT1 string, Sep1 string, INT2 string, MONTH_NAME1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, ""))},
  /*  33 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, INT2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1))},
  /*  34 */ func (INT1 string, Sep1 string, INT2 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1))},
  /*  35 */ func (DateTimeTZ1 *DateTimeTZ, Sep1 string, DateTimeTZ2 *DateTimeTZ) *DateTimeTZRange {return &DateTimeTZRange{Start: DateTimeTZ1, End: DateTimeTZ2}},
  /*  36 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, MONTH_NAME2 string, INT2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEndDates(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME2, INT2, YEAR1))},
  /*  37 */ func (Date1 civil.Date) *DateTimeTZ {return NewDateTimeWithDate(Date1)},
  /*  38 */ func (Date1 civil.Date, AtOpt1 string, Time1 civil.Time) *DateTimeTZ {return NewDateTime(Date1, Time1, "")},
  /*  39 */ func (YEAR1 string) civil.Date {return NewDMYDate("", "", YEAR1)},
  /*  40 */ func (YEAR1 string, SUB1 string, INT1 string) civil.Date {return NewDMYDate("", INT1, YEAR1)},
  /*  41 */ func (YEAR1 string, SUB1 string, INT1 string, SUB2 string, INT2 string, TOpt1 string) civil.Date {return NewDMYDate(INT2, INT1, YEAR1)},
  /*  42 */ func (MONTH_NAME1 string, INT1 string, YEAR1 string) civil.Date {return NewMDYDate(MONTH_NAME1, INT1, YEAR1)},
  /*  43 */ func (INT1 string, MONTH_NAME1 string, YEAR1 string) civil.Date {return NewDMYDate(INT1, MONTH_NAME1, YEAR1)},
  /*  44 */ func (INT1 string, QUO1 string, INT2 string, QUO2 string, YEAR1 string) civil.Date {return NewAmbiguousDate(ambiguousDateMode, INT1, INT2, YEAR1)},
  /*  45 */ func (WeekDayNameOpt1 string, MONTH_NAME1 string, INT1 string) civil.Date {return NewMDYDate(MONTH_NAME1, INT1, "")},
  /*  46 */ func (WeekDayNameOpt1 string, INT1 string, MONTH_NAME1 string) civil.Date {return NewDMYDate(INT1, MONTH_NAME1, "")},
  /*  47 */ func () string {return ""},
  /*  48 */ func (T1 string) string {return T1},
  /*  49 */ func () string {return ""},
  /*  50 */ func (WEEKDAY_NAME1 string) string {return WEEKDAY_NAME1},
  /*  51 */ func (INT1 string) civil.Time {return NewTime(INT1, "", "", "")},
  /*  52 */ func (INT1 string, AM1 string) civil.Time {return NewTime(INT1, "", "", "")},
  /*  53 */ func (INT1 string, PM1 string) civil.Time {return NewTime((mustAtoi(INT1) % 12) + 12, "", "", "")},
  /*  54 */ func (INT1 string, COLON1 string, INT2 string) civil.Time {return NewTime(INT1, INT2, "", "")},
  /*  55 */ func (INT1 string, COLON1 string, INT2 string, COLON2 string, INT3 string) civil.Time {return NewTime((mustAtoi(INT1) % 12) + 12, INT2, INT3, "")},
  /*  56 */ func (INT1 string, COLON1 string, INT2 string, AM1 string) civil.Time {return NewTime(INT1, INT2, "", "")},
  /*  57 */ func (INT1 string, COLON1 string, INT2 string, PM1 string) civil.Time {return NewTime((mustAtoi(INT1) % 12) + 12, INT2, "", "")},
  /*  58 */ func () string {return ""},
  /*  59 */ func (AND1 string) string {return AND1},
  /*  60 */ func () string {return ""},
  /*  61 */ func (AT1 string) string {return AT1},
  /*  62 */ func (SUB1 string) string {return SUB1},
  /*  63 */ func (THROUGH1 string) string {return THROUGH1},
  /*  64 */ func (TO1 string) string {return TO1},
  /*  65 */ func (WhenOpt1 string) string {return WhenOpt1},
  /*  66 */ func () string {return ""},
  /*  67 */ func (WHEN1 string) string {return WHEN1},
  /*  68 */ func (GoogleOpt1 string, CalendarOpt1 string, ICSOpt1 string) string {return GoogleOpt1},
  /*  69 */ func () string {return ""},
  /*  70 */ func (GOOGLE1 string) string {return GOOGLE1},
  /*  71 */ func () string {return ""},
  /*  72 */ func (CALENDAR1 string) string {return CALENDAR1},
  /*  73 */ func () string {return ""},
  /*  74 */ func (ICS1 string) string {return ICS1},
}}

var parseStates = &glr.ParseStates{Items:[]glr.ParseState{
  /*   0 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:66}}, "WHEN":[]glr.Action{glr.Action{Type:"shift", State:4, Rule:0}}}, Gotos:map[string]int{"PrefixOpt":2, "WhenOpt":3, "root":1}},
  /*   1 */ glr.ParseState{Actions:map[string][]glr.Action{"$end":[]glr.Action{glr.Action{Type:"accept", State:0, Rule:0}}}, Gotos:map[string]int{}},
  /*   2 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:49}, glr.Action{Type:"reduce", State:0, Rule:49}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:8, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:7, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:13, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:11, Rule:0}}}, Gotos:map[string]int{"Date":10, "DateTimeTZ":9, "DateTimeTZRange":6, "DateTimeTZRanges":5, "WeekDayNameOpt":12}},
  /*   3 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:65}}}, Gotos:map[string]int{}},
  /*   4 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:67}}}, Gotos:map[string]int{}},
  /*   5 */ glr.ParseState{Actions:map[string][]glr.Action{"$end":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:69}}, ".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:58}}, "AND":[]glr.Action{glr.Action{Type:"shift", State:17, Rule:0}}, "CALENDAR":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:69}}, "GOOGLE":[]glr.Action{glr.Action{Type:"shift", State:18, Rule:0}}, "ICS":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:69}}}, Gotos:map[string]int{"AndOpt":15, "GoogleOpt":16, "SuffixOpt":14}},
  /*   6 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:2}}}, Gotos:map[string]int{}},
  /*   7 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}}, Gotos:map[string]int{}},
  /*   8 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:22, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}}, Gotos:map[string]int{"Sep":21}},
  /*   9 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:29}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}}, Gotos:map[string]int{"Sep":27}},
  /*  10 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:60}, glr.Action{Type:"reduce", State:0, Rule:37}}, "AT":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}}, Gotos:map[string]int{"AtOpt":28}},
  /*  11 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:39}, glr.Action{Type:"reduce", State:0, Rule:39}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:30, Rule:0}}}, Gotos:map[string]int{}},
  /*  12 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:32, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}}, Gotos:map[string]int{}},
  /*  13 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:50}}}, Gotos:map[string]int{}},
  /*  14 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:1}}}, Gotos:map[string]int{}},
  /*  15 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:49}, glr.Action{Type:"reduce", State:0, Rule:49}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:35, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:13, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:11, Rule:0}}}, Gotos:map[string]int{"Date":10, "DateTimeTZ":9, "DateTimeTZRange":33, "WeekDayNameOpt":12}},
  /*  16 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:71}}, "CALENDAR":[]glr.Action{glr.Action{Type:"shift", State:37, Rule:0}}}, Gotos:map[string]int{"CalendarOpt":36}},
  /*  17 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:59}}}, Gotos:map[string]int{}},
  /*  18 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:70}}}, Gotos:map[string]int{}},
  /*  19 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:38, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:40, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:41, Rule:0}}}, Gotos:map[string]int{"Sep":39}},
  /*  20 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:43, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:42, Rule:0}}}, Gotos:map[string]int{}},
  /*  21 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}}, Gotos:map[string]int{}},
  /*  22 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:45, Rule:0}}}, Gotos:map[string]int{}},
  /*  23 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:46, Rule:0}}}, Gotos:map[string]int{}},
  /*  24 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:62}}}, Gotos:map[string]int{}},
  /*  25 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:63}}}, Gotos:map[string]int{}},
  /*  26 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:64}}}, Gotos:map[string]int{}},
  /*  27 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:49}, glr.Action{Type:"reduce", State:0, Rule:49}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:49, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:50, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:13, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:11, Rule:0}}}, Gotos:map[string]int{"Date":10, "DateTimeTZ":48, "Time":47, "WeekDayNameOpt":12}},
  /*  28 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:52, Rule:0}}}, Gotos:map[string]int{"Time":51}},
  /*  29 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:61}}}, Gotos:map[string]int{}},
  /*  30 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:53, Rule:0}}}, Gotos:map[string]int{}},
  /*  31 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:54, Rule:0}}}, Gotos:map[string]int{}},
  /*  32 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}}, Gotos:map[string]int{}},
  /*  33 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:3}}}, Gotos:map[string]int{}},
  /*  34 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}}, Gotos:map[string]int{}},
  /*  35 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:22, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}}, Gotos:map[string]int{"Sep":57}},
  /*  36 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:73}}, "ICS":[]glr.Action{glr.Action{Type:"shift", State:59, Rule:0}}}, Gotos:map[string]int{"ICSOpt":58}},
  /*  37 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:72}}}, Gotos:map[string]int{}},
  /*  38 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:4}, glr.Action{Type:"reduce", State:0, Rule:4}, glr.Action{Type:"reduce", State:0, Rule:4}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:61, Rule:0}}}, Gotos:map[string]int{}},
  /*  39 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:63, Rule:0}}}, Gotos:map[string]int{}},
  /*  40 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:64, Rule:0}}}, Gotos:map[string]int{}},
  /*  41 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:42}}}, Gotos:map[string]int{}},
  /*  42 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:8}, glr.Action{Type:"reduce", State:0, Rule:8}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:65, Rule:0}}}, Gotos:map[string]int{}},
  /*  43 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:67, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:66, Rule:0}}}, Gotos:map[string]int{}},
  /*  44 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:68, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}}, Gotos:map[string]int{}},
  /*  45 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:43}}}, Gotos:map[string]int{}},
  /*  46 */ glr.ParseState{Actions:map[string][]glr.Action{"QUO":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}}, Gotos:map[string]int{}},
  /*  47 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:30}}}, Gotos:map[string]int{}},
  /*  48 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:35}}}, Gotos:map[string]int{}},
  /*  49 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:51}, glr.Action{Type:"reduce", State:0, Rule:51}}, "AM":[]glr.Action{glr.Action{Type:"shift", State:71, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:73, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:22, Rule:0}}, "PM":[]glr.Action{glr.Action{Type:"shift", State:72, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}}, Gotos:map[string]int{}},
  /*  50 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:74, Rule:0}}}, Gotos:map[string]int{}},
  /*  51 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:38}}}, Gotos:map[string]int{}},
  /*  52 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:51}}, "AM":[]glr.Action{glr.Action{Type:"shift", State:71, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:73, Rule:0}}, "PM":[]glr.Action{glr.Action{Type:"shift", State:72, Rule:0}}}, Gotos:map[string]int{}},
  /*  53 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:40}, glr.Action{Type:"reduce", State:0, Rule:40}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:75, Rule:0}}}, Gotos:map[string]int{}},
  /*  54 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:45}}}, Gotos:map[string]int{}},
  /*  55 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:46}}}, Gotos:map[string]int{}},
  /*  56 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:41, Rule:0}}}, Gotos:map[string]int{"Sep":76}},
  /*  57 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:77, Rule:0}}}, Gotos:map[string]int{}},
  /*  58 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:68}}}, Gotos:map[string]int{}},
  /*  59 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:74}}}, Gotos:map[string]int{}},
  /*  60 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:5}, glr.Action{Type:"reduce", State:0, Rule:5}, glr.Action{Type:"reduce", State:0, Rule:5}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:78, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:79, Rule:0}}}, Gotos:map[string]int{}},
  /*  61 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:12}}}, Gotos:map[string]int{}},
  /*  62 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:31}, glr.Action{Type:"reduce", State:0, Rule:31}, glr.Action{Type:"reduce", State:0, Rule:31}, glr.Action{Type:"reduce", State:0, Rule:31}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:80, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:81, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:82, Rule:0}}}, Gotos:map[string]int{}},
  /*  63 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:83, Rule:0}}}, Gotos:map[string]int{}},
  /*  64 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:84, Rule:0}}}, Gotos:map[string]int{}},
  /*  65 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:16}}}, Gotos:map[string]int{}},
  /*  66 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:9}, glr.Action{Type:"reduce", State:0, Rule:9}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:85, Rule:0}}}, Gotos:map[string]int{}},
  /*  67 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:87, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:86, Rule:0}}}, Gotos:map[string]int{}},
  /*  68 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}}, Gotos:map[string]int{"Sep":88}},
  /*  69 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:32}, glr.Action{Type:"reduce", State:0, Rule:32}, glr.Action{Type:"reduce", State:0, Rule:32}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:89, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:90, Rule:0}}}, Gotos:map[string]int{}},
  /*  70 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:91, Rule:0}}}, Gotos:map[string]int{}},
  /*  71 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:52}}}, Gotos:map[string]int{}},
  /*  72 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:53}}}, Gotos:map[string]int{}},
  /*  73 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:92, Rule:0}}}, Gotos:map[string]int{}},
  /*  74 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:41, Rule:0}}}, Gotos:map[string]int{}},
  /*  75 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:93, Rule:0}}}, Gotos:map[string]int{}},
  /*  76 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:94, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:63, Rule:0}}}, Gotos:map[string]int{}},
  /*  77 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:95, Rule:0}}}, Gotos:map[string]int{}},
  /*  78 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:6}, glr.Action{Type:"reduce", State:0, Rule:6}, glr.Action{Type:"reduce", State:0, Rule:6}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:96, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:97, Rule:0}}}, Gotos:map[string]int{}},
  /*  79 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:13}}}, Gotos:map[string]int{}},
  /*  80 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}}, Gotos:map[string]int{"Sep":98}},
  /*  81 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:99, Rule:0}}}, Gotos:map[string]int{}},
  /*  82 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:33}}}, Gotos:map[string]int{}},
  /*  83 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:100, Rule:0}}}, Gotos:map[string]int{}},
  /*  84 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:28}}}, Gotos:map[string]int{}},
  /*  85 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:17}}}, Gotos:map[string]int{}},
  /*  86 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:10}, glr.Action{Type:"reduce", State:0, Rule:10}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:101, Rule:0}}}, Gotos:map[string]int{}},
  /*  87 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:102, Rule:0}}}, Gotos:map[string]int{}},
  /*  88 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:103, Rule:0}}}, Gotos:map[string]int{}},
  /*  89 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}}, Gotos:map[string]int{"Sep":104}},
  /*  90 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:34}}}, Gotos:map[string]int{}},
  /*  91 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:44}}}, Gotos:map[string]int{}},
  /*  92 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:54}}, "AM":[]glr.Action{glr.Action{Type:"shift", State:106, Rule:0}}, "COLON":[]glr.Action{glr.Action{Type:"shift", State:105, Rule:0}}, "PM":[]glr.Action{glr.Action{Type:"shift", State:107, Rule:0}}}, Gotos:map[string]int{}},
  /*  93 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:47}}, "T":[]glr.Action{glr.Action{Type:"shift", State:109, Rule:0}}}, Gotos:map[string]int{"TOpt":108}},
  /*  94 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:31}, glr.Action{Type:"reduce", State:0, Rule:31}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:82, Rule:0}}}, Gotos:map[string]int{}},
  /*  95 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:32}, glr.Action{Type:"reduce", State:0, Rule:32}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:90, Rule:0}}}, Gotos:map[string]int{}},
  /*  96 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:7}, glr.Action{Type:"reduce", State:0, Rule:7}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:110, Rule:0}}}, Gotos:map[string]int{}},
  /*  97 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:14}}}, Gotos:map[string]int{}},
  /*  98 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:111, Rule:0}}}, Gotos:map[string]int{}},
  /*  99 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:25, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}}, Gotos:map[string]int{"Sep":112}},
  /* 100 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:36}}}, Gotos:map[string]int{}},
  /* 101 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:18}}}, Gotos:map[string]int{}},
  /* 102 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:11}, glr.Action{Type:"reduce", State:0, Rule:11}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:113, Rule:0}}}, Gotos:map[string]int{}},
  /* 103 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:114, Rule:0}}}, Gotos:map[string]int{}},
  /* 104 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:115, Rule:0}}}, Gotos:map[string]int{}},
  /* 105 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:116, Rule:0}}}, Gotos:map[string]int{}},
  /* 106 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:56}}}, Gotos:map[string]int{}},
  /* 107 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:57}}}, Gotos:map[string]int{}},
  /* 108 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:41}}}, Gotos:map[string]int{}},
  /* 109 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:48}}}, Gotos:map[string]int{}},
  /* 110 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:15}}}, Gotos:map[string]int{}},
  /* 111 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:20}, glr.Action{Type:"reduce", State:0, Rule:20}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:117, Rule:0}}}, Gotos:map[string]int{}},
  /* 112 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:118, Rule:0}}}, Gotos:map[string]int{}},
  /* 113 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:19}}}, Gotos:map[string]int{}},
  /* 114 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:21}, glr.Action{Type:"reduce", State:0, Rule:21}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:119, Rule:0}}}, Gotos:map[string]int{}},
  /* 115 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:120, Rule:0}}}, Gotos:map[string]int{}},
  /* 116 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:55}}}, Gotos:map[string]int{}},
  /* 117 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:24}}}, Gotos:map[string]int{}},
  /* 118 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:22}, glr.Action{Type:"reduce", State:0, Rule:22}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:121, Rule:0}}}, Gotos:map[string]int{}},
  /* 119 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:25}}}, Gotos:map[string]int{}},
  /* 120 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:23}, glr.Action{Type:"reduce", State:0, Rule:23}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:122, Rule:0}}}, Gotos:map[string]int{}},
  /* 121 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:26}}}, Gotos:map[string]int{}},
  /* 122 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:27}}}, Gotos:map[string]int{}},
}}

