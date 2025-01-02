package glr

/*
Rules

root:
  ABCD
root:
  ABC
root:
  BCD
ABCD:
  A B C D
ABC:
  A B C
BCD:
  B C OptD
BCD:
  B C WrapD
BCD:
  B WrapC D
WrapC:
  C
OptD:
  <empty>
OptD:
  D
WrapD:
  D
*/

var glrRules = &Rules{Items:[]Rule{
  /*   0 */ Rule{Nonterminal:"", RHS:[]string(nil)}, // ignored because rule-numbering starts at 1
  /*   1 */ Rule{Nonterminal:"root", RHS:[]string{"ABCD"}},
  /*   2 */ Rule{Nonterminal:"root", RHS:[]string{"ABC"}},
  /*   3 */ Rule{Nonterminal:"root", RHS:[]string{"BCD"}},
  /*   4 */ Rule{Nonterminal:"ABCD", RHS:[]string{"A", "B", "C", "D"}},
  /*   5 */ Rule{Nonterminal:"ABC", RHS:[]string{"A", "B", "C"}},
  /*   6 */ Rule{Nonterminal:"BCD", RHS:[]string{"B", "C", "OptD"}},
  /*   7 */ Rule{Nonterminal:"BCD", RHS:[]string{"B", "C", "WrapD"}},
  /*   8 */ Rule{Nonterminal:"BCD", RHS:[]string{"B", "WrapC", "D"}},
  /*   9 */ Rule{Nonterminal:"WrapC", RHS:[]string{"C"}},
  /*  10 */ Rule{Nonterminal:"OptD", RHS:[]string(nil)},
  /*  11 */ Rule{Nonterminal:"OptD", RHS:[]string{"D"}},
  /*  12 */ Rule{Nonterminal:"WrapD", RHS:[]string{"D"}},
}}

var glrStates = &ParseStates{Items:[]ParseState{
  /*   0 */ ParseState{Actions:map[string][]Action{"A":[]Action{Action{Type:"shift", State:5, Rule:0}}, "B":[]Action{Action{Type:"shift", State:6, Rule:0}}}, Gotos:map[string]int{"ABC":3, "ABCD":2, "BCD":4, "root":1}},
  /*   1 */ ParseState{Actions:map[string][]Action{"$end":[]Action{Action{Type:"accept", State:0, Rule:0}}}, Gotos:map[string]int{}},
  /*   2 */ ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:1}}}, Gotos:map[string]int{}},
  /*   3 */ ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:2}}}, Gotos:map[string]int{}},
  /*   4 */ ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:3}}}, Gotos:map[string]int{}},
  /*   5 */ ParseState{Actions:map[string][]Action{"B":[]Action{Action{Type:"shift", State:7, Rule:0}}}, Gotos:map[string]int{}},
  /*   6 */ ParseState{Actions:map[string][]Action{"C":[]Action{Action{Type:"shift", State:8, Rule:0}}}, Gotos:map[string]int{"WrapC":9}},
  /*   7 */ ParseState{Actions:map[string][]Action{"C":[]Action{Action{Type:"shift", State:10, Rule:0}}}, Gotos:map[string]int{}},
  /*   8 */ ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:9}, Action{Type:"reduce", State:0, Rule:10}}, "D":[]Action{Action{Type:"shift", State:13, Rule:0}}}, Gotos:map[string]int{"OptD":11, "WrapD":12}},
  /*   9 */ ParseState{Actions:map[string][]Action{"D":[]Action{Action{Type:"shift", State:14, Rule:0}}}, Gotos:map[string]int{}},
  /*  10 */ ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:5}}, "D":[]Action{Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{}},
  /*  11 */ ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:6}}}, Gotos:map[string]int{}},
  /*  12 */ ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:7}}}, Gotos:map[string]int{}},
  /*  13 */ ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:12}, Action{Type:"reduce", State:0, Rule:11}}}, Gotos:map[string]int{}},
  /*  14 */ ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:8}}}, Gotos:map[string]int{}},
  /*  15 */ ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:4}}}, Gotos:map[string]int{}},
}}

