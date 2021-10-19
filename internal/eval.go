package internal

import "fmt"

func Eval(expr Expr) Expr {
	switch expr.Kind() {
	case exprAtom:
		panic("cannot evaluate atom: " + expr.String())
	case exprNumber:
		return expr
	case exprString:
		return expr
	case exprList:
		lst := expr.(List)
		if len(lst.values) == 0 {
			return nil
		}
		fnName := lst.values[0]
		if fnName.Kind() != exprAtom {
			if len(lst.values) != 1 {
				panic("expected an atom as a function name in list: " + expr.String())
			}
			// cover (1) and ("foo")
			return lst.values[0]
		}
		args := []Expr{}
		for _, subExpr := range lst.values[1:] {
			args = append(args, Eval(subExpr))
		}
		switch fnName.String() {
		case "print":
			return builtinPrint(args)
		case "+":
			return builtinPlus(args)
		case "*":
			return builtinMult(args)
		}
	case exprListLit:
		panic("todo")
	default:
		panic("unreachable")
	}
	return nil
}

func builtinPrint(args []Expr) Expr {
	for _, arg := range args {
		if s, ok := arg.(String); ok {
			fmt.Print(s.value)
			continue
		}
		fmt.Print(arg.String())
	}
	return nil
}
func builtinPlus(args []Expr) Expr {
	total := int64(0)
	for _, arg := range args {
		num, _ := arg.(Number)
		total += num.value
	}
	return Number{value: total}
}
func builtinMult(args []Expr) Expr {
	total := int64(1)
	for _, arg := range args {
		num, _ := arg.(Number)
		total *= num.value
	}
	return Number{value: total}
}
