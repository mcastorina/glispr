package internal

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ExprType int

const (
	// Expr types
	exprAtom = iota
	exprNumber
	exprString
	exprList
	exprListLit
)

// Expr is used to allow many concrete types to be an expression
type Expr interface {
	Kind() ExprType
	String() string
}

// Atom represents an atom expression
type Atom struct{ value string }

func (a Atom) Kind() ExprType { return exprAtom }
func (a Atom) String() string { return a.value }

// Number represents a number expression
type Number struct{ value int64 }

func (n Number) Kind() ExprType { return exprNumber }
func (n Number) String() string { return fmt.Sprintf("%d", n.value) }

// String represents a string expression
type String struct{ value string }

func (s String) Kind() ExprType { return exprString }
func (s String) String() string { return fmt.Sprintf(`"%s"`, s.value) }

// List represents a literal or non-literal list of expressions
type List struct {
	values []Expr
	isLit  bool
}

func (l List) Kind() ExprType {
	if l.isLit {
		return exprListLit
	}
	return exprList
}
func (l List) String() string {
	if len(l.values) == 0 {
		if l.isLit {
			return "'()"
		}
		return "()"
	}
	out := "("
	if l.isLit {
		out = "'("
	}
	for _, expr := range l.values[:len(l.values)-1] {
		out += expr.String() + " "
	}
	out += l.values[len(l.values)-1].String() + ")"
	return out
}

// Parser parses an input into an Expr AST
type Parser struct {
	lexer   *Lexer
	peekTok *Token
	peekStr string
}

func NewParser(input io.Reader) *Parser {
	return &Parser{
		lexer: NewLexer(input),
	}
}

func NewStringParser(input string) *Parser {
	return &Parser{
		lexer: NewLexer(strings.NewReader(input)),
	}
}

// next returns the next Token and literal string from the input
func (p *Parser) next() (Token, string) {
	if p.peekTok != nil {
		defer func() {
			p.peekTok = nil
		}()
		return *p.peekTok, p.peekStr
	}
	return p.lexer.Next()
}

// peek returns the next Token from the input without consuming it
func (p *Parser) peek() Token {
	if p.peekTok != nil {
		return *p.peekTok
	}
	tok, val := p.next()
	p.peekTok = &tok
	p.peekStr = val
	return tok
}

// consume takes the next token and panics if it was not the expected Token
func (p *Parser) consume(expected Token) {
	tok, _ := p.next()
	if tok != expected {
		panic(fmt.Sprintf("Expected %s, found %s", expected, tok))
	}
}

// Expression parses the input in its current state and produces an Expr
func (p *Parser) Expression() Expr {
	tok, val := p.next()
	switch tok {
	case tokAtom:
		return Atom{value: val}
	case tokMinus:
		// could be a number or a `-` atom
		if p.peek() != tokNumber {
			return Atom{value: val}
		}
		tok, val = p.next()
		num, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			panic("Unexpected error: " + err.Error())
		}
		return Number{value: -num}
	case tokNumber:
		num, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			panic("Unexpected error: " + err.Error())
		}
		return Number{value: num}
	case tokString:
		// remove surrounding double quotes
		return String{value: val[1 : len(val)-1]}
	case tokLParen:
		// build a slice of Expr until we reach a closing parenthesis
		exprs := []Expr{}
		for p.peek() != tokRParen {
			exprs = append(exprs, p.Expression())
		}
		p.consume(tokRParen)
		return List{values: exprs}
	case tokQuote:
		list := p.Expression().(List)
		list.isLit = true
		return list
	default:
		panic("Unknown start of expression: " + tok.String())
	}
}
