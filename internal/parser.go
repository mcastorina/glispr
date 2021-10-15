package internal

import (
	"fmt"
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

type Expr interface {
	Kind() ExprType
	String() string
}

type Atom struct{ value string }

func (a Atom) Kind() ExprType { return exprAtom }
func (a Atom) String() string { return a.value }

type Number struct{ value int64 }

func (n Number) Kind() ExprType { return exprNumber }
func (n Number) String() string { return fmt.Sprintf("%d", n.value) }

type String struct{ value string }

func (s String) Kind() ExprType { return exprString }
func (s String) String() string { return fmt.Sprintf(`"%s"`, s.value) }

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

type Parser struct {
	lexer   *Lexer
	peekTok *Token
	peekStr string
}

func NewParser(input string) *Parser {
	return &Parser{
		lexer: NewLexer(strings.NewReader(input)),
	}
}

func (p *Parser) next() (Token, string) {
	if p.peekTok != nil {
		defer func() {
			p.peekTok = nil
		}()
		return *p.peekTok, p.peekStr
	}
	return p.lexer.Next()
}

func (p *Parser) peek() Token {
	if p.peekTok != nil {
		return *p.peekTok
	}
	tok, val := p.next()
	p.peekTok = &tok
	p.peekStr = val
	return tok
}

func (p *Parser) consume(expected Token) {
	tok, _ := p.next()
	if tok != expected {
		panic(fmt.Sprintf("Expected %s, found %s", expected, tok))
	}
}

func (p *Parser) Expression() Expr {
	tok, val := p.next()
	switch tok {
	case tokAtom:
		return Atom{value: val}
	case tokMinus:
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
		return String{value: val[1 : len(val)-1]}
	case tokLParen:
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
