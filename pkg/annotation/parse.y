%{
package annotation
%}

%union{
  Expr     Matchers
  Matchers Matchers
  Matcher  Matcher
  str      string
}

%start expr

%type  <Expr>        expr
%type  <Matchers>    matchers
%type  <Matcher>     matcher

%token <str>  IDENTIFIER STRING
%token <val>  EQ NEQ RE NRE OPEN_BRACE CLOSE_BRACE COMMA

%%

expr: OPEN_BRACE matchers CLOSE_BRACE { yylex.(*lexer).output = $2 };

matchers:
      matcher                  { $$ = []Matcher{ $1 } }
    | matchers COMMA matcher   { $$ = append($1, $3) }
    ;

matcher:
      IDENTIFIER EQ STRING     { $$ = Eq($1, $3) }
    | IDENTIFIER NEQ STRING    { $$ = Ne($1, $3) }
    | IDENTIFIER RE STRING     { $$ = Re($1, $3) }
    | IDENTIFIER NRE STRING    { $$ = Nre($1, $3) }
    ;

%%
