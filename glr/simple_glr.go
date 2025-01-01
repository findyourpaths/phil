package glr

var glrRules = &Rules{Items:[]Rule{
  Rule{Nonterminal:"", RHS:[]string(nil)},
  Rule{Nonterminal:"root", RHS:[]string{"ABCD"}},
  Rule{Nonterminal:"root", RHS:[]string{"ABC"}},
  Rule{Nonterminal:"root", RHS:[]string{"BCD"}},
  Rule{Nonterminal:"ABCD", RHS:[]string{"A", "B", "C", "D"}},
  Rule{Nonterminal:"ABC", RHS:[]string{"A", "B", "C"}},
  Rule{Nonterminal:"BCD", RHS:[]string{"B", "C", "OptD"}},
  Rule{Nonterminal:"BCD", RHS:[]string{"B", "C", "WrapD"}},
  Rule{Nonterminal:"BCD", RHS:[]string{"B", "WrapC", "D"}},
  Rule{Nonterminal:"WrapC", RHS:[]string{"C"}},
  Rule{Nonterminal:"OptD", RHS:[]string(nil)},
  Rule{Nonterminal:"OptD", RHS:[]string{"D"}},
  Rule{Nonterminal:"WrapD", RHS:[]string{"D"}},
}}

var glrStates = &ParseStates{Items:[]ParseState{
  ParseState{Actions:map[string][]Action{"A":[]Action{Action{Type:"shift", State:5, Rule:0}}, "B":[]Action{Action{Type:"shift", State:6, Rule:0}}}, Gotos:map[string]int{"ABC":3, "ABCD":2, "BCD":4, "root":1}},
  ParseState{Actions:map[string][]Action{"$end":[]Action{Action{Type:"accept", State:0, Rule:0}}}, Gotos:map[string]int{}},
  ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:1}}}, Gotos:map[string]int{}},
  ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:2}}}, Gotos:map[string]int{}},
  ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:3}}}, Gotos:map[string]int{}},
  ParseState{Actions:map[string][]Action{"B":[]Action{Action{Type:"shift", State:7, Rule:0}}}, Gotos:map[string]int{}},
  ParseState{Actions:map[string][]Action{"C":[]Action{Action{Type:"shift", State:8, Rule:0}}}, Gotos:map[string]int{"WrapC":9}},
  ParseState{Actions:map[string][]Action{"C":[]Action{Action{Type:"shift", State:10, Rule:0}}}, Gotos:map[string]int{}},
  ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:9}, Action{Type:"reduce", State:0, Rule:10}}, "D":[]Action{Action{Type:"shift", State:13, Rule:0}}}, Gotos:map[string]int{"OptD":11, "WrapD":12}},
  ParseState{Actions:map[string][]Action{"D":[]Action{Action{Type:"shift", State:14, Rule:0}}}, Gotos:map[string]int{}},
  ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:5}}, "D":[]Action{Action{Type:"shift", State:15, Rule:0}}}, Gotos:map[string]int{}},
  ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:6}}}, Gotos:map[string]int{}},
  ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:7}}}, Gotos:map[string]int{}},
  ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:12}, Action{Type:"reduce", State:0, Rule:11}}}, Gotos:map[string]int{}},
  ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:8}}}, Gotos:map[string]int{}},
  ParseState{Actions:map[string][]Action{".":[]Action{Action{Type:"reduce", State:0, Rule:4}}}, Gotos:map[string]int{}},
}}

