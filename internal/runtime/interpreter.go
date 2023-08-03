package runtime

import (
	"errors"
	"fmt"
	"gox/internal"
	ast2 "gox/internal/ast"
	"gox/internal/scanning"
	"strings"
)

// TODO: write tests
type Interpreter struct {
	env *environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: newEnvironment(nil),
	}
}

func (r *Interpreter) Interpret(statements []*ast2.Stmt) *internal.RuntimeError {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("unable to interpret given code: %s", err)
		}
	}()
	for _, stmt := range statements {
		if stmt == nil {
			continue
		}
		err := r.execute(*stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// expressions
func (r *Interpreter) VisitForLiteral(expr *ast2.Literal) (any, *internal.RuntimeError) {
	return expr.Value, nil
}

func (r *Interpreter) VisitForUnary(expr *ast2.Unary) (any, *internal.RuntimeError) {
	right, err := r.evaluate(*expr.Right)
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

func (r *Interpreter) VisitForBinary(expr *ast2.Binary) (any, *internal.RuntimeError) {
	left, err := r.evaluate(*expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := r.evaluate(*expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case scanning.EQUAL_EQUAL:
		bothOperandsAreSting := checkBothOperandsAreString(left, right)
		if bothOperandsAreSting {
			return strings.Compare(left.(string), right.(string)) == 0, nil
		} else {
			if err = r.checkNumberOperands(*expr.Operator, left, right); err != nil {
				return nil, err
			}
			return r.isEqual(left, right), nil
		}
	case scanning.BANG_EQUAL:
		bothOperandsAreSting := checkBothOperandsAreString(left, right)
		if bothOperandsAreSting {
			return !(strings.Compare(left.(string), right.(string)) == 0), nil
		} else {
			if err = r.checkNumberOperands(*expr.Operator, left, right); err != nil {
				return nil, err
			}
			return !r.isEqual(left, right), nil
		}
	case scanning.GREATER:
		if err = r.checkNumberOperands(*expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case scanning.GREATER_EQUAL:
		if err = r.checkNumberOperands(*expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case scanning.LESS:
		if err = r.checkNumberOperands(*expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case scanning.LESS_EQUAL:
		if err = r.checkNumberOperands(*expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case scanning.MINUS:
		if err = r.checkNumberOperands(*expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case scanning.PLUS:
		bothOperandsAreSting := checkBothOperandsAreString(left, right)
		if bothOperandsAreSting {
			return fmt.Sprint(left.(string) + right.(string)), nil
		} else {
			if err = r.checkNumberOperands(*expr.Operator, left, right); err != nil {
				return nil, err
			}
			return left.(float64) + right.(float64), nil
		}
	case scanning.SLASH:
		if err = r.checkNumberOperands(*expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case scanning.STAR:
		if err = r.checkNumberOperands(*expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	}
	return nil, nil
}

func (r *Interpreter) VisitForGrouping(expr *ast2.Grouping) (any, *internal.RuntimeError) {
	return r.evaluate(*expr.Expression)
}

func (r *Interpreter) VisitForVariableExpression(expr *ast2.VarExpr) (any, *internal.RuntimeError) {
	res, err := r.env.get(expr.Name)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Interpreter) VisitForAssignExpression(expr *ast2.Assign) (any, *internal.RuntimeError) {
	val, err := r.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	err = r.env.assign(expr.Name, val)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *Interpreter) VisitForLogical(expr *ast2.Logical) (any, *internal.RuntimeError) {
	left, err := r.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.TokenType == scanning.OR {
		if r.isTruthy(left) {
			return left, nil
		}
	} else {
		if !r.isTruthy(left) {
			return left, nil
		}
	}

	return r.evaluate(expr.Right)
}

// statements
func (r *Interpreter) VisitForExpression(stmt *ast2.Expression) *internal.RuntimeError {
	_, err := r.evaluate(*stmt.Expression)
	return err
}

func (r *Interpreter) VisitForPrint(stmt *ast2.Print) *internal.RuntimeError {
	value, err := r.evaluate(*stmt.Expression)
	if err == nil {
		fmt.Println(toString(value))
	}
	return err
}

func (r *Interpreter) VisitForVar(stmt *ast2.Var) *internal.RuntimeError {
	if stmt.Initializer != nil {
		value, err := r.evaluate(*stmt.Initializer)
		if err != nil {
			return err
		}
		r.env.define(stmt.Name.Lexeme, value)
	} else {
		r.env.define(stmt.Name.Lexeme, nil)
	}
	return nil
}

func (r *Interpreter) VisitForBlock(block *ast2.Block) *internal.RuntimeError {
	return r.executeBlock(block.Statements, newEnvironment(r.env))
}

func (r *Interpreter) VisitForIf(ifStmt *ast2.If) *internal.RuntimeError {
	conditionRes, err := r.evaluate(ifStmt.Condition)
	if err != nil {
		return err
	}
	if r.isTruthy(conditionRes) {
		r.execute(ifStmt.Then)
	} else if ifStmt.Else != nil {
		r.execute(ifStmt.Else)
	}
	return nil
}

func (r *Interpreter) evaluate(expr ast2.Expr) (any, *internal.RuntimeError) {
	return expr.Accept(r)
}

func (r *Interpreter) isTruthy(right any) bool {
	if right == nil {
		return false
	}
	if _, ok := right.(bool); ok {
		return right.(bool)
	} else {
		if _, ok := right.(string); ok {
			return len(right.(string)) > 0
		}
	}
	return false
}

func (r *Interpreter) isEqual(a, b any) bool {
	return (a == nil && b == nil) && (a == b)
}

func (r *Interpreter) checkNumberOperand(operator *scanning.Token, operand any) *internal.RuntimeError {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return &internal.RuntimeError{
		Error: errors.New("operand must be number"),
		Token: operator,
	}
}

func checkBothOperandsAreString(left any, right any) bool {
	_, leftIsString := left.(string)
	_, rightIsString := right.(string)
	bothOperandsAreSting := leftIsString && rightIsString
	return bothOperandsAreSting
}

func (r *Interpreter) checkNumberOperands(operator scanning.Token, operand1, operand2 any) *internal.RuntimeError {
	if _, ok := operand1.(float64); ok {
		if _, ok := operand2.(float64); ok {
			return nil
		}
	}
	return &internal.RuntimeError{
		Error: errors.New("both operands must be numbers"),
		Token: &operator,
	}
}

func (r *Interpreter) execute(stmt ast2.Stmt) *internal.RuntimeError {
	return stmt.Accept(r)
}

func (r *Interpreter) executeBlock(statements []*ast2.Stmt, env *environment) *internal.RuntimeError {
	prevEnv := r.env

	// after execution
	defer func() {
		r.env = prevEnv
	}()

	r.env = env

	for _, statement := range statements {
		err := r.execute(*statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func toString(a any) string {
	return fmt.Sprint(a)
}
