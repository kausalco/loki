%{
package annotation
%}

%union{
  Expr       Matchers
  Matchers   Matchers
  Matcher    Matcher
  str        string
  int        int64
  Identifier string
}

%start expr

%type  <Expr>        expr
%type  <Matchers>    matchers
%type  <Matcher>     matcher
%type  <Identifier>  identifier

%token <str>  IDENTIFIER STRING
%token <int>  INT
%token <val>  EQ NEQ RE NRE OPEN_BRACE CLOSE_BRACE COMMA DOT

%%

expr: OPEN_BRACE matchers CLOSE_BRACE { yylex.(*lexer).output = $2 };

matchers:
      matcher                  { $$ = []Matcher{ $1 } }
    | matchers COMMA matcher   { $$ = append($1, $3) }
    ;

matcher:
      identifier EQ STRING     { $$ = EqStr($1, $3) }
    | identifier NEQ STRING    { $$ = NeStr($1, $3) }
    | identifier RE STRING     { $$ = Re($1, $3) }
    | identifier NRE STRING    { $$ = Nre($1, $3) }
    | identifier EQ INT        { $$ = EqInt($1, $3) }
    | identifier NEQ INT       { $$ = NeInt($1, $3) }
    ;

identifier:
      IDENTIFIER                { $$ = $1 }
    | identifier DOT IDENTIFIER { $$ = $1 + "." + $3 }
    ;

%%
