package parse

import "github.com/findyourpaths/phil/glr"

/*
Rules

root:
  DateTimeTZRanges
DateTimeTZRanges:
  DateTimeTZRange
DateTimeTZRange:
  DateTimeTZ
DateTimeTZ:
  Date
Date:
  MonthName Day Year
MonthName:
  MONTH
Day:
  INT
Year:
  INT
*/

var parseRules = &glr.Rules{Items:[]glr.Rule{
  /*   0 */ glr.Rule{Nonterminal:"", RHS:[]string(nil)}, // ignored because rule-numbering starts at 1
  /*   1 */ glr.Rule{Nonterminal:"root", RHS:[]string{"DateTimeTZRanges"}},
  /*   2 */ glr.Rule{Nonterminal:"DateTimeTZRanges", RHS:[]string{"DateTimeTZRange"}},
  /*   3 */ glr.Rule{Nonterminal:"DateTimeTZRange", RHS:[]string{"DateTimeTZ"}},
  /*   4 */ glr.Rule{Nonterminal:"DateTimeTZ", RHS:[]string{"Date"}},
  /*   5 */ glr.Rule{Nonterminal:"Date", RHS:[]string{"MonthName", "Day", "Year"}},
  /*   6 */ glr.Rule{Nonterminal:"MonthName", RHS:[]string{"MONTH"}},
  /*   7 */ glr.Rule{Nonterminal:"Day", RHS:[]string{"INT"}},
  /*   8 */ glr.Rule{Nonterminal:"Year", RHS:[]string{"INT"}},
}}

var parseStates = &glr.ParseStates{Items:[]glr.ParseState{
  /*   0 */ glr.ParseState{Actions:map[string][]glr.Action{"MONTH":[]glr.Action{glr.Action{Type:"shift", State:7, Rule:0}}}, Gotos:map[string]int{"Date":5, "DateTimeTZ":4, "DateTimeTZRange":3, "DateTimeTZRanges":2, "MonthName":6, "root":1}},
  /*   1 */ glr.ParseState{Actions:map[string][]glr.Action{"$end":[]glr.Action{glr.Action{Type:"accept", State:0, Rule:0}}}, Gotos:map[string]int{}},
  /*   2 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:1}}}, Gotos:map[string]int{}},
  /*   3 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:2}}}, Gotos:map[string]int{}},
  /*   4 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:3}}}, Gotos:map[string]int{}},
  /*   5 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:4}}}, Gotos:map[string]int{}},
  /*   6 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:9, Rule:0}}}, Gotos:map[string]int{"Day":8}},
  /*   7 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:6}}}, Gotos:map[string]int{}},
  /*   8 */ glr.ParseState{Actions:map[string][]glr.Action{"INT":[]glr.Action{glr.Action{Type:"shift", State:11, Rule:0}}}, Gotos:map[string]int{"Year":10}},
  /*   9 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:7}}}, Gotos:map[string]int{}},
  /*  10 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:5}}}, Gotos:map[string]int{}},
  /*  11 */ glr.ParseState{Actions:map[string][]glr.Action{".":[]glr.Action{glr.Action{Type:"reduce", State:0, Rule:8}}}, Gotos:map[string]int{}},
}}

