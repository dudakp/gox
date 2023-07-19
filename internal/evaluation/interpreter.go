package evaluation

import (
	"gox/internal/parsing"
	"gox/internal/scanning"
)

type Interpreter struct {
}

func (r *Interpreter) VisitForLiteral(expr *parsing.Literal) any {
	return expr.Value
}

func (r *Interpreter) VisitForUnary(expr *parsing.Unary) any {
	right := r.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case scanning.BANG:
		return !r.isTruthy(right)
	case scanning.MINUS:
		return -right.(float64)
	}

	return nil
}

func (r *Interpreter) VisitForBinary(expr *parsing.Binary) any {
	//TODO implement me
	panic("implement me")
}

func (r *Interpreter) VisitForGrouping(expr *parsing.Grouping) any {
	return r.evaluate(expr)
}

func (r *Interpreter) evaluate(expr parsing.Expr) any {
	return expr.Accept(r)
}

func (r *Interpreter) isTruthy(right any) bool {

}
