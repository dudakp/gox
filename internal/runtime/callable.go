package runtime

import (
	"gox/internal"
	"gox/internal/ast"
)

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, args []any) (any, *internal.RuntimeError)
}

type LoxStdFunction interface {
	Callable
	Name() string
}

type LoxFunction struct {
	declaration ast.Function
}

func (r *LoxFunction) Arity() int {
	//TODO implement me
	panic("implement me")
}

func (r *LoxFunction) Call(interpreter *Interpreter, args []any) (any, *internal.RuntimeError) {
	//TODO implement me
	panic("implement me")
}
