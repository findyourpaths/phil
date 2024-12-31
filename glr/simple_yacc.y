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
  A B C D {$$ = nil}
;

ABC:
  A B C {$$ = nil}
;

BCD:
  B C OptD {$$ = nil}
| B C WrapD {$$ = nil}
| B WrapC D {$$ = nil}
;

WrapC:
  C
;

OptD:

| D
;

WrapD:
  D
;

%%
