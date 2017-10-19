package annotation

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

func Parse(input string) (Matchers, error) {
	l := lexer{
		Scanner: scanner.Scanner{
			Mode: scanner.SkipComments | scanner.ScanStrings,
		},
	}
	l.Init(strings.NewReader(input))
	parser := yyNewParser()
	e := parser.Parse(&l)
	if e != 0 {
		return nil, fmt.Errorf(l.error)
	}
	return l.output, nil
}

type lexer struct {
	scanner.Scanner
	output Matchers
	error  string
}

func (l *lexer) Lex(lval *yySymType) int {
	r := l.Scan()
	switch r {
	case scanner.EOF:
		return 0
	case scanner.String:
		var err error
		lval.str, err = strconv.Unquote(l.TokenText())
		if err != nil {
			panic(err)
		}
		return STRING
	}

	switch l.TokenText() {
	case "=":
		if l.Peek() == '~' {
			l.Scan()
			return RE
		}
		return EQ
	case "!":
		if l.Peek() == '=' {
			l.Scan()
			return NEQ
		} else if l.Peek() == '~' {
			l.Scan()
			return NRE
		}
	case ",":
		return COMMA
	case "{":
		return OPEN_BRACE
	case "}":
		return CLOSE_BRACE
	}

	lval.str = l.TokenText()
	return IDENTIFIER
}

func (l *lexer) Error(s string) {
	l.error = s
}
