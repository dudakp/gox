package parsing

import "fmt"

type AstPrinter struct {
}

func (r *AstPrinter) Print(expr Expr) (string, error) {
	return expr.Accept(r).(string), nil
}

func (r *AstPrinter) VisitForLiteral(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprint(expr.Value)
}

func (r *AstPrinter) VisitForUnary(expr *Unary) any {
	return r.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (r *AstPrinter) VisitForBinary(expr *Binary) any {
	return r.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (r *AstPrinter) VisitForGrouping(expr *Grouping) any {
	return r.parenthesize("group", expr.expression)
}

func (r *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	res := "(" + name
	for _, expr := range exprs {
		accept := expr.Accept(r)
		res += " " + accept.(string)
	}
	res += ")"
	return res
}
