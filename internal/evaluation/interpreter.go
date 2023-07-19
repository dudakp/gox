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
	left := r.evaluate(expr.Left)
	right := r.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case scanning.GREATER:
		return left.(float64) > right.(float64)
	case scanning.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case scanning.LESS:
		return left.(float64) < right.(float64)
	case scanning.LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case scanning.EQUAL:
		return r.isEqual(left, right)
	case scanning.BANG_EQUAL:
		return !r.isEqual(left, right)
	case scanning.MINUS:
		return left.(float64) - right.(float64)
	case scanning.STRING:
		_, leftIsString := left.(string)
		_, rightIsString := right.(string)
		if leftIsString && rightIsString {
			return left.(string) + right.(string)
		} else {
			return left.(float64) + right.(float64)
		}
	case scanning.SLASH:
		return left.(float64) / right.(float64)
	case scanning.STAR:
		return left.(float64) * right.(float64)
	}
	return nil
}

func (r *Interpreter) VisitForGrouping(expr *parsing.Grouping) any {
	return r.evaluate(expr)
}

func (r *Interpreter) evaluate(expr parsing.Expr) any {
	return expr.Accept(r)
}

func (r *Interpreter) isTruthy(right any) bool {
	if right == nil {
		return false
	}
	if _, ok := right.(bool); ok {
		return right.(bool)
	}
	return false
}

func (r *Interpreter) isEqual(a, b any) bool {
	return a == b
}
