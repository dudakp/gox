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
	Env *environment
}

func NewInterpreter() *Interpreter {
	glob := newEnvironment(nil)
	for _, fn := range StdFunctions {
		glob.define(fn.Name(), fn)
	}

	return &Interpreter{
		Env: glob,
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
	res, err := r.Env.get(expr.Name)
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
	err = r.Env.assign(expr.Name, val)
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

func (r *Interpreter) VisitForFunctionCall(call *ast2.Call) (any, *internal.RuntimeError) {
	callee, runtimeError := r.evaluate(call.Callee)
	if runtimeError != nil {
		return nil, runtimeError
	}

	args := make([]any, len(call.Params))
	for i, arg := range call.Params {
		val, err := r.evaluate(arg)
		if err != nil {
			return nil, err
		}
		args[i] = val
	}

	_, isStdFunc := callee.(Callable)
	_, isLoxFunc := callee.(LoxFunction)
	if !isStdFunc && !isLoxFunc {
		return nil, &internal.RuntimeError{
			Error: errors.New("non-callable element"),
			Token: call.Paren,
		}
	}

	// arity check
	if len(args) != len(call.Params) {
		return nil, &internal.RuntimeError{
			Error: errors.New("invalid number of arguments"),
			Token: call.Paren,
		}
	}

	// TODO: this is shit, find a better, idiomatic way
	if isStdFunc {
		return callee.(LoxStdFunction).Call(r, args)
	}
	if isLoxFunc {
		// TODO: I dont want to assign to variable, I just want to call "Call" method
		function := callee.(LoxFunction)
		return function.Call(r, args)
	}
	return nil, nil
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
		r.Env.define(stmt.Name.Lexeme, value)
	} else {
		r.Env.define(stmt.Name.Lexeme, nil)
	}
	return nil
}

func (r *Interpreter) VisitForBlock(block *ast2.Block) *internal.RuntimeError {
	return r.executeBlock(block.Statements, newEnvironment(r.Env))
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

func (r *Interpreter) VisitForWhile(while *ast2.While) *internal.RuntimeError {
	conditionRes, err := r.evaluate(while.Condition)
	if err != nil {
		return err
	}
	for r.isTruthy(conditionRes) {
		err := r.execute(while.Statement)
		if err != nil {
			return err
		}
		conditionRes, err = r.evaluate(while.Condition)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Interpreter) VisitForFunction(function *ast2.Function) *internal.RuntimeError {
	fun := LoxFunction{declaration: *function}
	r.Env.define(function.Name.Lexeme, fun)
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

func (r *Interpreter) executeBlock(statements []ast2.Stmt, env *environment) *internal.RuntimeError {
	prevEnv := r.Env

	// after execution
	defer func() {
		r.Env = prevEnv
	}()

	r.Env = env

	for _, statement := range statements {
		err := r.execute(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func toString(a any) string {
	return fmt.Sprint(a)
}
