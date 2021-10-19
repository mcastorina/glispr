package internal

import "testing"

func TestParseLiterals(t *testing.T) {
	type testInput struct {
		input    string
		expected Expr
	}
	parse := func(input string) Expr {
		parser := NewStringParser(input)
		return parser.Expression()
	}
	tests := []testInput{
		{"1  ", Number{value: 1}},
		{"-1", Number{value: -1}},
		{"+", Atom{value: "+"}},
		{"-", Atom{value: "-"}},
		{"#foo", Atom{value: "#foo"}},
		{"bar", Atom{value: "bar"}},
		{"  \"string\"", String{value: "string"}},
	}

	for _, test := range tests {
		if have := parse(test.input); have != test.expected {
			t.Errorf("Found %T, want %T", have, test.expected)
		}
	}
}

func TestParseList(t *testing.T) {
	type testInput struct {
		input    string
		expected string
	}
	parse := func(input string) Expr {
		parser := NewStringParser(input)
		return parser.Expression()
	}
	tests := []testInput{
		{
			"   (   foo    bar    baz ) ",
			"(foo bar baz)",
		},
		{
			"'(foo bar-baz)",
			"'(foo bar-baz)",
		},
		{
			"(print(*(+ 1 2)3))",
			"(print (* (+ 1 2) 3))",
		},
		{
			"(print '(foo bar))",
			"(print '(foo bar))",
		},
	}

	for _, test := range tests {
		if have := parse(test.input).String(); have != test.expected {
			t.Errorf("Found %s, want %s", have, test.expected)
		}
	}
}

func TestEval(t *testing.T) {
	type testInput struct {
		input        string
		expected     string
		expectedKind ExprType
	}
	tests := []testInput{
		{
			"(* (+ 1 2) 3)",
			"9",
			exprNumber,
		},
	}
	for _, test := range tests {
		parser := NewStringParser(test.input)
		expr := parser.Expression()
		result := Eval(expr)
		if result.Kind() != test.expectedKind || result.String() != test.expected {
			t.Errorf("Expression %s evaluated to %s, want %d", expr, result, 9)
		}
	}
}
