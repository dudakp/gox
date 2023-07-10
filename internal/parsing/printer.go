package parsing

import "fmt"

type AstPrinter struct {
}

func (r *AstPrinter) Print(expr Expr) (string, error) {
	return expr.accept(r).(string), nil
}

func (r *AstPrinter) visitForLiteral(expr *Literal) any {
	if expr.value == nil {
		return "nil"
	}
	return fmt.Sprint(expr.value)
}

func (r *AstPrinter) visitForUnary(expr *Unary) any {
	return r.parenthesize(expr.operator.Lexeme, expr.right)
}

func (r *AstPrinter) visitForBinary(expr *Binary) any {
	return r.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (r *AstPrinter) visitForGrouping(expr *Grouping) any {
	return r.parenthesize("group", expr.expression)
}

func (r *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	res := "(" + name
	for _, expr := range exprs {
		accept := expr.accept(r)
		res += " " + accept.(string)
	}
	res += ")"
	return res
}
