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
			Mode: scanner.SkipComments | scanner.ScanStrings | scanner.ScanInts,
		},
	}
	l.Init(strings.NewReader(input))
	parser := yyNewParser()
	e := parser.Parse(&l)
	if e != 0 {
		return nil, l.err
	}
	return l.output, nil
}

type lexer struct {
	scanner.Scanner
	output Matchers
	err    error
}

var tokens = map[string]int{
	",": COMMA,
	".": DOT,
	"{": OPEN_BRACE,
	"}": CLOSE_BRACE,
}

func (l *lexer) Lex(lval *yySymType) int {
	r := l.Scan()
	var err error
	switch r {
	case scanner.EOF:
		return 0

	case scanner.Int:
		lval.int, err = strconv.ParseInt(l.TokenText(), 10, 64)
		if err != nil {
			l.err = err
			return 0
		}
		return INT

	case scanner.String:
		lval.str, err = strconv.Unquote(l.TokenText())
		if err != nil {
			l.err = err
			return 0
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
	}

	if token, ok := tokens[l.TokenText()]; ok {
		return token
	}

	lval.str = l.TokenText()
	return IDENTIFIER
}

func (l *lexer) Error(s string) {
	l.err = fmt.Errorf(s)
}
