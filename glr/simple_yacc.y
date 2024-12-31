%{
package glr

type node struct {
    val string
    children []*node
}

%}

%union {
    token string
    node  *node
}

%token A B C D X Y ILLEGAL

%type <node> ABCD ABC BCD

%%

root:
  ABCD
| ABC
| BCD
;

ABCD:
  A B C D {$$ = &node{val: "ABCD", children: []*node{{val: "A"}, {val: "B"}, {val: "C"}, {val: "D"}}}}
;

ABC:
  A B C {$$ = &node{val: "ABC", children: []*node{{val: "A"}, {val: "B"}, {val: "C"}}}}
;

BCD:
  B C OptD {$$ = &node{val: "BCD", children: []*node{{val: "B"}, {val: "C"}, {val: "D"}}}}
;

OptD:

| D
;

%%
