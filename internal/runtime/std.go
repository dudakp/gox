package runtime

import (
	"gox/internal"
	"time"
)

var StdFunctions []LoxStdFunction

func init() {
	StdFunctions = make([]LoxStdFunction, 1)
	StdFunctions[0] = &clock{}
}

type clock struct {
}

func (c *clock) Name() string {
	return "clock"
}

func (c *clock) Arity() int {
	return 0
}

func (c *clock) Call(interpreter *Interpreter, args []any) (any, *internal.RuntimeError) {
	return time.Now().Unix(), nil
}
