package internal

import (
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	type nextOut struct {
		tok Token
		val string
	}
	runTest := func(input string, expected []nextOut) {
		lexer := NewLexer(strings.NewReader(input))
		for _, exp := range expected {
			tok, val := lexer.Next()
			if tok != exp.tok || val != exp.val {
				t.Errorf("Found (%s, %s), want (%s, %s)", tok, val, exp.tok, exp.val)
			}
		}
	}

	runTest("foo", []nextOut{
		{tokAtom, "foo"},
		{tokEOF, ""},
	})
	runTest("foo bar baz", []nextOut{
		{tokAtom, "foo"},
		{tokAtom, "bar"},
		{tokAtom, "baz"},
		{tokEOF, ""},
	})
	runTest("(foo '(bar))", []nextOut{
		{tokLParen, "("},
		{tokAtom, "foo"},
		{tokQuote, "'"},
		{tokLParen, "("},
		{tokAtom, "bar"},
		{tokRParen, ")"},
		{tokRParen, ")"},
		{tokEOF, ""},
	})
	runTest("-1234 1_000", []nextOut{
		{tokMinus, "-"},
		{tokNumber, "1234"},
		{tokNumber, "1_000"},
		{tokEOF, ""},
	})
	runTest(`#foo      "bar()"`, []nextOut{
		{tokAtom, "#foo"},
		{tokString, `"bar()"`},
		{tokEOF, ""},
	})
	runTest(`"back\"sla;sh"`, []nextOut{
		{tokString, `"back"sla;sh"`},
		{tokEOF, ""},
	})
}
