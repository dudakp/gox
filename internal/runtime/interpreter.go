package runtime

import (
	"errors"
	"fmt"
	"gox/internal/parsing"
	"gox/internal/scanning"
)

type RuntimeError struct {
	error
	Token scanning.Token
}

// TODO: write tests
type Interpreter struct {
}

func (r *Interpreter) Interpret(expr parsing.Expr) (string, *RuntimeError) {
	res, err := r.evaluate(expr)
	if err != nil {
		return "", err.(*RuntimeError)
	}
	return toString(res), nil
}

func (r *Interpreter) VisitForLiteral(expr *parsing.Literal) (any, error) {
	return expr.Value, nil
}

func (r *Interpreter) VisitForUnary(expr *parsing.Unary) (any, error) {
	right, err := r.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case scanning.BANG:
		return !r.isTruthy(right), nil
	case scanning.MINUS:
		return -right.(float64), r.checkNumberOperand(expr.Operator, right)
	}

	return nil, nil
}

func (r *Interpreter) VisitForBinary(expr *parsing.Binary) (any, error) {
	left, err := r.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := r.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case scanning.GREATER:
		if err = r.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case scanning.GREATER_EQUAL:
		if err = r.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case scanning.LESS:
		if err = r.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case scanning.LESS_EQUAL:
		if err = r.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case scanning.EQUAL:
		if err = r.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return r.isEqual(left, right), nil
	case scanning.BANG_EQUAL:
		if err = r.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return !r.isEqual(left, right), nil
	case scanning.MINUS:
		if err = r.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case scanning.PLUS:
		_, leftIsString := left.(string)
		_, rightIsString := right.(string)
		if leftIsString && rightIsString {
			return left.(string) + right.(string), nil
		} else {
			if err = r.checkNumberOperands(expr.Operator, left, right); err != nil {
				return nil, err
			}
			return left.(float64) + right.(float64), nil
		}
	case scanning.SLASH:
		if err = r.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case scanning.STAR:
		if err = r.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	}
	return nil, nil
}

func (r *Interpreter) VisitForGrouping(expr *parsing.Grouping) (any, error) {
	return r.evaluate(expr.Expression)
}

func (r *Interpreter) evaluate(expr parsing.Expr) (any, error) {
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
	return (a == nil && b == nil) && (a == b)
}

func (r *Interpreter) checkNumberOperand(operator scanning.Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return &RuntimeError{
		error: errors.New("operand must be number"),
		Token: operator,
	}
}

func (r *Interpreter) checkNumberOperands(operator scanning.Token, operand1, operand2 any) error {
	if _, ok := operand1.(float64); ok {
		if _, ok := operand2.(float64); ok {
			return nil
		}
	}
	return &RuntimeError{
		error: errors.New("both operands must be numbers"),
		Token: operator,
	}
}

func toString(a any) string {
	return fmt.Sprint(a)
}
