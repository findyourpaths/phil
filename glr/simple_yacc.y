%{
package glr
%}

 /* Type of each nonterminal. */
%type <Alphabet> root
%type <ABCD> ABCD
%type <ABC> ABC
%type <BCD> BCD
%type <BCDEF> BCDEF
%type <string> WrapC
%type <string> OptD
%type <string> WrapD

%start root

%union {
    string string
    Alphabet *Alphabet
    ABCD *ABCD
    ABC *ABC
    BCD *BCD
    BCDEF *BCDEF
}

%token <string> A B C D E F X Y ILLEGAL

%%

root:
  ABCD {$$ = &Alphabet{ABCD: $1}}
| ABC {$$ = &Alphabet{ABC: $1}}
| BCD {$$ = &Alphabet{BCD: $1}}
| BCDEF {$$ = &Alphabet{BCDEF: $1}}
;

ABCD:
  A B C D {$$ = &ABCD{A: $1, B: $2, C: $3, D: $4}}
;

ABC:
  A B C {$$ = &ABC{A: $1, B: $2, C: $3}}
;

BCD:
  B C OptD {$$ = &BCD{B: $1, C: $2, D: $3}}
| B C WrapD {$$ = &BCD{B: $1, C: $2, D: $3}}
| B WrapC D {$$ = &BCD{B: $1, C: $2, D: $3}}
;

BCDEF:
  B C D E F {$$ = &BCDEF{B: $1, C: $2, D: $3, E: $4, F: $5}}
;

WrapC:
  C
;

OptD:
     {$$ = ""}
| D
;

WrapD:
  D
;

%%
