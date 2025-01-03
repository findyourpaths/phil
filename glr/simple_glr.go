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
  /*   0 */ Rule{Nonterminal:"", RHS:[]string(nil), Type:""}, // ignored because rule-numbering starts at 1
  /*   1 */ Rule{Nonterminal:"root", RHS:[]string{"ABCD"}, Type:"*Alphabet"},
  /*   2 */ Rule{Nonterminal:"root", RHS:[]string{"ABC"}, Type:"*Alphabet"},
  /*   3 */ Rule{Nonterminal:"root", RHS:[]string{"BCD"}, Type:"*Alphabet"},
  /*   4 */ Rule{Nonterminal:"ABCD", RHS:[]string{"A", "B", "C", "D"}, Type:"*ABCD"},
  /*   5 */ Rule{Nonterminal:"ABC", RHS:[]string{"A", "B", "C"}, Type:"*ABC"},
  /*   6 */ Rule{Nonterminal:"BCD", RHS:[]string{"B", "C", "OptD"}, Type:"*BCD"},
  /*   7 */ Rule{Nonterminal:"BCD", RHS:[]string{"B", "C", "WrapD"}, Type:"*BCD"},
  /*   8 */ Rule{Nonterminal:"BCD", RHS:[]string{"B", "WrapC", "D"}, Type:"*BCD"},
  /*   9 */ Rule{Nonterminal:"WrapC", RHS:[]string{"C"}, Type:"string"},
  /*  10 */ Rule{Nonterminal:"OptD", RHS:[]string(nil), Type:"string"},
  /*  11 */ Rule{Nonterminal:"OptD", RHS:[]string{"D"}, Type:"string"},
  /*  12 */ Rule{Nonterminal:"WrapD", RHS:[]string{"D"}, Type:"string"},
}}

// Semantic action functions

var glrActions = &SemanticActions{Items:[]any{
  /*   0 */ nil, // empty type
  /*   1 */ func (ABCD1 *ABCD) *Alphabet {return &Alphabet{ABCD: ABCD1}},
  /*   2 */ func (ABC1 *ABC) *Alphabet {return &Alphabet{ABC: ABC1}},
  /*   3 */ func (BCD1 *BCD) *Alphabet {return &Alphabet{BCD: BCD1}},
  /*   4 */ func (A1 string, B1 string, C1 string, D1 string) *ABCD {return &ABCD{A: A1, B: B1, C: C1, D: D1}},
  /*   5 */ func (A1 string, B1 string, C1 string) *ABC {return &ABC{A: A1, B: B1, C: C1}},
  /*   6 */ func (B1 string, C1 string, OptD1 string) *BCD {return &BCD{B: B1, C: C1, D: OptD1}},
  /*   7 */ func (B1 string, C1 string, WrapD1 string) *BCD {return &BCD{B: B1, C: C1, D: WrapD1}},
  /*   8 */ func (B1 string, WrapC1 string, D1 string) *BCD {return &BCD{B: B1, C: WrapC1, D: D1}},
  /*   9 */ func (C1 string) string {return C1},
  /*  10 */ func () string {return ""},
  /*  11 */ func (D1 string) string {return D1},
  /*  12 */ func (D1 string) string {return D1},
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

