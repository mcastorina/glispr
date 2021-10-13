package internal

import "fmt"

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

func (a *Atom) Kind() ExprType { return exprAtom }
func (a *Atom) String() string { return a.value }

type Number struct{ value int64 }

func (n *Number) Kind() ExprType { return exprNumber }
func (n *Number) String() string { return fmt.Sprintf("%d", n.value) }

type String struct{ value string }

func (s *String) Kind() ExprType { return exprString }
func (s *String) String() string { return fmt.Sprintf(`"%s"`, s.value) }

type List struct {
	values []Expr
	isLit  bool
}

func (l *List) Kind() ExprType {
	if l.isLit {
		return exprListLit
	}
	return exprList
}
func (l *List) String() string {
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

func Run() {
}
