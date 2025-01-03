package parse

import "github.com/findyourpaths/phil/glr"
import "cloud.google.com/go/civil"

/*
Rules

root:
  DateTimeTZRanges
DateTimeTZRanges:
  DateTimeTZRange
DateTimeTZRanges:
  DateTimeTZRanges DateTimeTZRange
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
WeekDayNameOpt:
  <empty>
WeekDayNameOpt:
  WEEKDAY_NAME
Sep:
  SUB
Sep:
  THROUGH
Sep:
  TO
*/

var parseRules = &glr.Rules{Items:[]glr.Rule{
  /*   0 */ glr.Rule{Nonterminal:"", RHS:[]string(nil), Type:""}, // ignored because rule-numbering starts at 1
  /*   1 */ glr.Rule{Nonterminal:"root", RHS:[]string{"DateTimeTZRanges"}, Type:"*DateTimeTZRanges"},
  /*   2 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"DateTimeTZRange"}, Type:"*DateTimeTZRanges"},
  /*   3 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"DateTimeTZRanges", "DateTimeTZRange"}, Type:"*DateTimeTZRanges"},
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
  /*  30 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"MONTH_NAME", "INT", "Sep", "INT"}, Type:"*DateTimeTZRange"},
  /*  31 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"INT", "Sep", "INT", "MONTH_NAME"}, Type:"*DateTimeTZRange"},
  /*  32 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"MONTH_NAME", "INT", "Sep", "INT", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  33 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"INT", "Sep", "INT", "MONTH_NAME", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  34 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"DateTimeTZ", "Sep", "DateTimeTZ"}, Type:"*DateTimeTZRange"},
  /*  35 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"MONTH_NAME", "INT", "Sep", "MONTH_NAME", "INT", "YEAR"}, Type:"*DateTimeTZRange"},
  /*  36 */ glr.Rule{Nonterminal:"DateTimeTZ", RHS:[]string{"Date"}, Type:"*DateTimeTZ"},
  /*  37 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"MONTH_NAME", "INT", "YEAR"}, Type:"civil.Date"},
  /*  38 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"INT", "MONTH_NAME", "YEAR"}, Type:"civil.Date"},
  /*  39 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"INT", "QUO", "INT", "QUO", "YEAR"}, Type:"civil.Date"},
  /*  40 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"WeekDayNameOpt", "MONTH_NAME", "INT"}, Type:"civil.Date"},
  /*  41 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"WeekDayNameOpt", "INT", "MONTH_NAME"}, Type:"civil.Date"},
  /*  42 */ glr.Rule{Nonterminal:"WeekDayNameOpt", RHS:[]string(nil), Type:""},
  /*  43 */ glr.Rule{Nonterminal:"WeekDayNameOpt", RHS:[]string{"WEEKDAY_NAME"}, Type:""},
  /*  44 */ glr.Rule{Nonterminal:"Sep", RHS:[]string{"SUB"}, Type:""},
  /*  45 */ glr.Rule{Nonterminal:"Sep", RHS:[]string{"THROUGH"}, Type:""},
  /*  46 */ glr.Rule{Nonterminal:"Sep", RHS:[]string{"TO"}, Type:""},
}}

// Semantic action functions

var parseActions = &glr.SemanticActions{Items:[]any{
  /*   0 */ nil, // empty action
  /*   1 */ func (DateTimeTZRanges1 *DateTimeTZRanges) *DateTimeTZRanges {return DateTimeTZRanges1},
  /*   2 */ func (DateTimeTZRange1 *DateTimeTZRange) *DateTimeTZRanges {return &DateTimeTZRanges{Items: []*DateTimeTZRange{DateTimeTZRange1}}},
  /*   3 */ func (DateTimeTZRanges1 *DateTimeTZRanges, DateTimeTZRange1 *DateTimeTZRange) *DateTimeTZRanges {return AppendDateTimeTZRanges(DateTimeTZRanges1, DateTimeTZRange1)},
  /*   4 */ func (MONTH_NAME1 string, INT1 string, INT2 string) *DateTimeTZRanges {return NewRangesWithStarts(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, ""))},
  /*   5 */ func (MONTH_NAME1 string, INT1 string, INT2 string, INT3 string) *DateTimeTZRanges {return NewRangesWithStarts(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, ""), NewMDYDate(MONTH_NAME1, INT3, ""))},
  /*   6 */ func (MONTH_NAME1 string, INT1 string, INT2 string, INT3 string, INT4 string) *DateTimeTZRanges {return NewRangesWithStarts(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, ""), NewMDYDate(MONTH_NAME1, INT3, ""), NewMDYDate(MONTH_NAME1, INT4, ""))},
  /*   7 */ func (MONTH_NAME1 string, INT1 string, INT2 string, INT3 string, INT4 string, INT5 string) *DateTimeTZRanges {return NewRangesWithStarts(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, ""), NewMDYDate(MONTH_NAME1, INT3, ""), NewMDYDate(MONTH_NAME1, INT4, ""), NewMDYDate(MONTH_NAME1, INT5, ""))},
  /*   8 */ func (INT1 string, INT2 string, MONTH_NAME1 string) *DateTimeTZRanges {return NewRangesWithStarts(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, ""))},
  /*   9 */ func (INT1 string, INT2 string, INT3 string, MONTH_NAME1 string) *DateTimeTZRanges {return NewRangesWithStarts(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, ""), NewDMYDate(INT3, MONTH_NAME1, ""))},
  /*  10 */ func (INT1 string, INT2 string, INT3 string, INT4 string, MONTH_NAME1 string) *DateTimeTZRanges {return NewRangesWithStarts(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, ""), NewDMYDate(INT3, MONTH_NAME1, ""), NewDMYDate(INT4, MONTH_NAME1, ""))},
  /*  11 */ func (INT1 string, INT2 string, INT3 string, INT4 string, INT5 string, MONTH_NAME1 string) *DateTimeTZRanges {return NewRangesWithStarts(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, ""), NewDMYDate(INT3, MONTH_NAME1, ""), NewDMYDate(INT4, MONTH_NAME1, ""), NewDMYDate(INT5, MONTH_NAME1, ""))},
  /*  12 */ func (MONTH_NAME1 string, INT1 string, INT2 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStarts(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1))},
  /*  13 */ func (MONTH_NAME1 string, INT1 string, INT2 string, INT3 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStarts(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1), NewMDYDate(MONTH_NAME1, INT3, YEAR1))},
  /*  14 */ func (MONTH_NAME1 string, INT1 string, INT2 string, INT3 string, INT4 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStarts(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1), NewMDYDate(MONTH_NAME1, INT3, YEAR1), NewMDYDate(MONTH_NAME1, INT4, YEAR1))},
  /*  15 */ func (MONTH_NAME1 string, INT1 string, INT2 string, INT3 string, INT4 string, INT5 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStarts(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1), NewMDYDate(MONTH_NAME1, INT3, YEAR1), NewMDYDate(MONTH_NAME1, INT4, YEAR1), NewMDYDate(MONTH_NAME1, INT5, YEAR1))},
  /*  16 */ func (INT1 string, INT2 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStarts(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1))},
  /*  17 */ func (INT1 string, INT2 string, INT3 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStarts(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1), NewDMYDate(INT3, MONTH_NAME1, YEAR1))},
  /*  18 */ func (INT1 string, INT2 string, INT3 string, INT4 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStarts(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1), NewDMYDate(INT3, MONTH_NAME1, YEAR1), NewDMYDate(INT4, MONTH_NAME1, YEAR1))},
  /*  19 */ func (INT1 string, INT2 string, INT3 string, INT4 string, INT5 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRanges {return NewRangesWithStarts(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1), NewDMYDate(INT3, MONTH_NAME1, YEAR1), NewDMYDate(INT4, MONTH_NAME1, YEAR1), NewDMYDate(INT5, MONTH_NAME1, YEAR1))},
  /*  20 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, INT2 string, INT3 string, Sep2 string, INT4 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEnd(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, "")), NewRangeWithStartEnd(NewMDYDate(MONTH_NAME1, INT3, ""), NewMDYDate(MONTH_NAME1, INT4, "")))},
  /*  21 */ func (INT1 string, Sep1 string, INT2 string, INT3 string, Sep2 string, INT4 string, MONTH_NAME1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEnd(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, "")), NewRangeWithStartEnd(NewDMYDate(INT3, MONTH_NAME1, ""), NewDMYDate(INT4, MONTH_NAME1, "")))},
  /*  22 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, INT2 string, MONTH_NAME2 string, INT3 string, Sep2 string, INT4 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEnd(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, "")), NewRangeWithStartEnd(NewMDYDate(MONTH_NAME2, INT3, ""), NewMDYDate(MONTH_NAME2, INT4, "")))},
  /*  23 */ func (INT1 string, Sep1 string, INT2 string, MONTH_NAME1 string, INT3 string, Sep2 string, INT4 string, MONTH_NAME2 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEnd(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, "")), NewRangeWithStartEnd(NewDMYDate(INT3, MONTH_NAME2, ""), NewDMYDate(INT4, MONTH_NAME2, "")))},
  /*  24 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, INT2 string, INT3 string, Sep2 string, INT4 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEnd(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1)), NewRangeWithStartEnd(NewMDYDate(MONTH_NAME1, INT3, YEAR1), NewMDYDate(MONTH_NAME1, INT4, YEAR1)))},
  /*  25 */ func (INT1 string, Sep1 string, INT2 string, INT3 string, Sep2 string, INT4 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEnd(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1)), NewRangeWithStartEnd(NewDMYDate(INT3, MONTH_NAME1, YEAR1), NewDMYDate(INT4, MONTH_NAME1, YEAR1)))},
  /*  26 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, INT2 string, MONTH_NAME2 string, INT3 string, Sep2 string, INT4 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEnd(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1)), NewRangeWithStartEnd(NewMDYDate(MONTH_NAME2, INT3, YEAR1), NewMDYDate(MONTH_NAME2, INT4, YEAR1)))},
  /*  27 */ func (INT1 string, Sep1 string, INT2 string, MONTH_NAME1 string, INT3 string, Sep2 string, INT4 string, MONTH_NAME2 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStartEnd(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1)), NewRangeWithStartEnd(NewDMYDate(INT3, MONTH_NAME2, YEAR1), NewDMYDate(INT4, MONTH_NAME2, YEAR1)))},
  /*  28 */ func (MONTH_NAME1 string, INT1 string, MONTH_NAME2 string, INT2 string, YEAR1 string) *DateTimeTZRanges {return NewRanges(NewRangeWithStart(NewMDYDate(MONTH_NAME1, INT1, YEAR1)), NewRangeWithStart(NewMDYDate(MONTH_NAME2, INT2, YEAR1)))},
  /*  29 */ func (DateTimeTZ1 *DateTimeTZ) *DateTimeTZRange {return &DateTimeTZRange{Start: DateTimeTZ1}},
  /*  30 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, INT2 string) *DateTimeTZRange {return NewRangeWithStartEnd(NewMDYDate(MONTH_NAME1, INT1, ""), NewMDYDate(MONTH_NAME1, INT2, ""))},
  /*  31 */ func (INT1 string, Sep1 string, INT2 string, MONTH_NAME1 string) *DateTimeTZRange {return NewRangeWithStartEnd(NewDMYDate(INT1, MONTH_NAME1, ""), NewDMYDate(INT2, MONTH_NAME1, ""))},
  /*  32 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, INT2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEnd(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME1, INT2, YEAR1))},
  /*  33 */ func (INT1 string, Sep1 string, INT2 string, MONTH_NAME1 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEnd(NewDMYDate(INT1, MONTH_NAME1, YEAR1), NewDMYDate(INT2, MONTH_NAME1, YEAR1))},
  /*  34 */ func (DateTimeTZ1 *DateTimeTZ, Sep1 string, DateTimeTZ2 *DateTimeTZ) *DateTimeTZRange {return &DateTimeTZRange{Start: DateTimeTZ1, End: DateTimeTZ2}},
  /*  35 */ func (MONTH_NAME1 string, INT1 string, Sep1 string, MONTH_NAME2 string, INT2 string, YEAR1 string) *DateTimeTZRange {return NewRangeWithStartEnd(NewMDYDate(MONTH_NAME1, INT1, YEAR1), NewMDYDate(MONTH_NAME2, INT2, YEAR1))},
  /*  36 */ func (Date1 civil.Date) *DateTimeTZ {return &DateTimeTZ{DateTime: civil.DateTime{Date: Date1}}},
  /*  37 */ func (MONTH_NAME1 string, INT1 string, YEAR1 string) civil.Date {return NewMDYDate(MONTH_NAME1, INT1, YEAR1)},
  /*  38 */ func (INT1 string, MONTH_NAME1 string, YEAR1 string) civil.Date {return NewDMYDate(INT1, MONTH_NAME1, YEAR1)},
  /*  39 */ func (INT1 string, QUO1 string, INT2 string, QUO2 string, YEAR1 string) civil.Date {return NewAmbiguousDate(ambiguousDateMode, INT1, INT2, YEAR1)},
  /*  40 */ func (WeekDayNameOpt1 string, MONTH_NAME1 string, INT1 string) civil.Date {return NewMDYDate(MONTH_NAME1, INT1, "")},
  /*  41 */ func (WeekDayNameOpt1 string, INT1 string, MONTH_NAME1 string) civil.Date {return NewDMYDate(INT1, MONTH_NAME1, "")},
  /*  42 */ func () string {return ""},
  /*  43 */ func (WEEKDAY_NAME1 string) string {return WEEKDAY_NAME1},
  /*  44 */ func (SUB1 string) string {return SUB1},
  /*  45 */ func (THROUGH1 string) string {return THROUGH1},
  /*  46 */ func (TO1 string) string {return TO1},
}}

var parseStates = &glr.ParseStates{Items:[]glr.ParseState{
  /*   0 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:42}, glr.Action{Type:"reduce", State:0, Rule:42}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:5, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:4, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:9, Rule:0}}}, Gotos:map[string]int{"Date":7, "DateTimeTZ":6, "DateTimeTZRange":3, "DateTimeTZRanges":2, "WeekDayNameOpt":8, "root":1}},
  /*   1 */ glr.ParseState{Actions:map[string][]glr.Action{"$end":[]glr.Action{glr.Action{Type:"accept", State:0, Rule:0}}}, Gotos:map[string]int{}},
  /*   2 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:42}, glr.Action{Type:"reduce", State:0, Rule:42}, glr.Action{Type:"reduce", State:0, Rule:1}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:12, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:11, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:9, Rule:0}}}, Gotos:map[string]int{"Date":7, "DateTimeTZ":6, "DateTimeTZRange":10, "WeekDayNameOpt":8}},
  /*   3 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:2}}}, Gotos:map[string]int{}},
  /*   4 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:13, Rule:0}}}, Gotos:map[string]int{}},
  /*   5 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:14, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:16, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:17, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:18, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}}, Gotos:map[string]int{"Sep":15}},
  /*   6 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:29}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:18, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}}, Gotos:map[string]int{"Sep":21}},
  /*   7 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:36}}}, Gotos:map[string]int{}},
  /*   8 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:23, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:22, Rule:0}}}, Gotos:map[string]int{}},
  /*   9 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:43}}}, Gotos:map[string]int{}},
  /*  10 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:3}}}, Gotos:map[string]int{}},
  /*  11 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:24, Rule:0}}}, Gotos:map[string]int{}},
  /*  12 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:16, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:17, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:18, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}}, Gotos:map[string]int{"Sep":25}},
  /*  13 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:26, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:28, Rule:0}}, "SUB":[]glr.Action{glr.Action{Type:"shift", State:18, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}}, Gotos:map[string]int{"Sep":27}},
  /*  14 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:31, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:30, Rule:0}}}, Gotos:map[string]int{}},
  /*  15 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:32, Rule:0}}}, Gotos:map[string]int{}},
  /*  16 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:33, Rule:0}}}, Gotos:map[string]int{}},
  /*  17 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:34, Rule:0}}}, Gotos:map[string]int{}},
  /*  18 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:44}}}, Gotos:map[string]int{}},
  /*  19 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:45}}}, Gotos:map[string]int{}},
  /*  20 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:46}}}, Gotos:map[string]int{}},
  /*  21 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:42}, glr.Action{Type:"reduce", State:0, Rule:42}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:37, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:36, Rule:0}}, "WEEKDAY_NAME":[]glr.Action{glr.Action{Type:"shift", State:9, Rule:0}}}, Gotos:map[string]int{"Date":7, "DateTimeTZ":35, "WeekDayNameOpt":8}},
  /*  22 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:38, Rule:0}}}, Gotos:map[string]int{}},
  /*  23 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:39, Rule:0}}}, Gotos:map[string]int{}},
  /*  24 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:18, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}}, Gotos:map[string]int{"Sep":40}},
  /*  25 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:41, Rule:0}}}, Gotos:map[string]int{}},
  /*  26 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:4}, glr.Action{Type:"reduce", State:0, Rule:4}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:42, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:43, Rule:0}}}, Gotos:map[string]int{}},
  /*  27 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:44, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:45, Rule:0}}}, Gotos:map[string]int{}},
  /*  28 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:46, Rule:0}}}, Gotos:map[string]int{}},
  /*  29 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:37}}}, Gotos:map[string]int{}},
  /*  30 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:8}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:47, Rule:0}}}, Gotos:map[string]int{}},
  /*  31 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:49, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:48, Rule:0}}}, Gotos:map[string]int{}},
  /*  32 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:50, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:51, Rule:0}}}, Gotos:map[string]int{}},
  /*  33 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:38}}}, Gotos:map[string]int{}},
  /*  34 */ glr.ParseState{Actions:map[string][]glr.Action{"QUO":[]glr.Action{glr.Action{Type:"shift", State:52, Rule:0}}}, Gotos:map[string]int{}},
  /*  35 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:34}}}, Gotos:map[string]int{}},
  /*  36 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:53, Rule:0}}}, Gotos:map[string]int{}},
  /*  37 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:16, Rule:0}}, "QUO":[]glr.Action{glr.Action{Type:"shift", State:17, Rule:0}}}, Gotos:map[string]int{}},
  /*  38 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:40}}}, Gotos:map[string]int{}},
  /*  39 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:41}}}, Gotos:map[string]int{}},
  /*  40 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:54, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:45, Rule:0}}}, Gotos:map[string]int{}},
  /*  41 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:55, Rule:0}}}, Gotos:map[string]int{}},
  /*  42 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:5}, glr.Action{Type:"reduce", State:0, Rule:5}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:56, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:57, Rule:0}}}, Gotos:map[string]int{}},
  /*  43 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:12}}}, Gotos:map[string]int{}},
  /*  44 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:30}, glr.Action{Type:"reduce", State:0, Rule:30}, glr.Action{Type:"reduce", State:0, Rule:30}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:58, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:59, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}}, Gotos:map[string]int{}},
  /*  45 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:61, Rule:0}}}, Gotos:map[string]int{}},
  /*  46 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:62, Rule:0}}}, Gotos:map[string]int{}},
  /*  47 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:16}}}, Gotos:map[string]int{}},
  /*  48 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:9}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:63, Rule:0}}}, Gotos:map[string]int{}},
  /*  49 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:65, Rule:0}}, "MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:64, Rule:0}}}, Gotos:map[string]int{}},
  /*  50 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:18, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}}, Gotos:map[string]int{"Sep":66}},
  /*  51 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:31}, glr.Action{Type:"reduce", State:0, Rule:31}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:67, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:68, Rule:0}}}, Gotos:map[string]int{}},
  /*  52 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:69, Rule:0}}}, Gotos:map[string]int{}},
  /*  53 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:29, Rule:0}}}, Gotos:map[string]int{}},
  /*  54 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:30}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:60, Rule:0}}}, Gotos:map[string]int{}},
  /*  55 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:31}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:68, Rule:0}}}, Gotos:map[string]int{}},
  /*  56 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:6}, glr.Action{Type:"reduce", State:0, Rule:6}}, "INT":[]glr.Action{glr.Action{Type:"shift", State:70, Rule:0}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:71, Rule:0}}}, Gotos:map[string]int{}},
  /*  57 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:13}}}, Gotos:map[string]int{}},
  /*  58 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:18, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}}, Gotos:map[string]int{"Sep":72}},
  /*  59 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:73, Rule:0}}}, Gotos:map[string]int{}},
  /*  60 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:32}}}, Gotos:map[string]int{}},
  /*  61 */ glr.ParseState{Actions:map[string][]glr.Action{"YEAR":[]glr.Action{glr.Action{Type:"shift", State:74, Rule:0}}}, Gotos:map[string]int{}},
  /*  62 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:28}}}, Gotos:map[string]int{}},
  /*  63 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:17}}}, Gotos:map[string]int{}},
  /*  64 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:10}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:75, Rule:0}}}, Gotos:map[string]int{}},
  /*  65 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:76, Rule:0}}}, Gotos:map[string]int{}},
  /*  66 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:77, Rule:0}}}, Gotos:map[string]int{}},
  /*  67 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:18, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}}, Gotos:map[string]int{"Sep":78}},
  /*  68 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:33}}}, Gotos:map[string]int{}},
  /*  69 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:39}}}, Gotos:map[string]int{}},
  /*  70 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:7}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:79, Rule:0}}}, Gotos:map[string]int{}},
  /*  71 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:14}}}, Gotos:map[string]int{}},
  /*  72 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:80, Rule:0}}}, Gotos:map[string]int{}},
  /*  73 */ glr.ParseState{Actions:map[string][]glr.Action{"SUB":[]glr.Action{glr.Action{Type:"shift", State:18, Rule:0}}, "THROUGH":[]glr.Action{glr.Action{Type:"shift", State:19, Rule:0}}, "TO":[]glr.Action{glr.Action{Type:"shift", State:20, Rule:0}}}, Gotos:map[string]int{"Sep":81}},
  /*  74 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:35}}}, Gotos:map[string]int{}},
  /*  75 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:18}}}, Gotos:map[string]int{}},
  /*  76 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:11}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:82, Rule:0}}}, Gotos:map[string]int{}},
  /*  77 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:83, Rule:0}}}, Gotos:map[string]int{}},
  /*  78 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:84, Rule:0}}}, Gotos:map[string]int{}},
  /*  79 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:15}}}, Gotos:map[string]int{}},
  /*  80 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:20}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:85, Rule:0}}}, Gotos:map[string]int{}},
  /*  81 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:86, Rule:0}}}, Gotos:map[string]int{}},
  /*  82 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:19}}}, Gotos:map[string]int{}},
  /*  83 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:21}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:87, Rule:0}}}, Gotos:map[string]int{}},
  /*  84 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH_NAME":[]glr.Action{glr.Action{Type:"shift", State:88, Rule:0}}}, Gotos:map[string]int{}},
  /*  85 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:24}}}, Gotos:map[string]int{}},
  /*  86 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:22}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:89, Rule:0}}}, Gotos:map[string]int{}},
  /*  87 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:25}}}, Gotos:map[string]int{}},
  /*  88 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:23}}, "YEAR":[]glr.Action{glr.Action{Type:"shift", State:90, Rule:0}}}, Gotos:map[string]int{}},
  /*  89 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:26}}}, Gotos:map[string]int{}},
  /*  90 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:27}}}, Gotos:map[string]int{}},
}}

