package runtime

import (
	"gox/internal"
	"gox/internal/ast"
)

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, args []any) (any, *internal.RuntimeError)
}

// LoxStdFunction represents native function implemented by interpreter
type LoxStdFunction interface {
	Callable
	Name() string
}

// LoxFunction represents user-defined function/method
type LoxFunction struct {
	declaration ast.Function
}

func (r *LoxFunction) Arity() int {
	return len(r.declaration.Params)
}

func (r *LoxFunction) Call(interpreter *Interpreter, args []any) (any, *internal.RuntimeError) {
	env := newEnvironment(interpreter.Env)
	for i, param := range r.declaration.Params {
		env.define(param.Lexeme, args[i])
	}

	return nil, interpreter.executeBlock(r.declaration.Body, env)
}
