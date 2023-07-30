package runtime

import (
	"errors"
	"gox/internal"
	"gox/internal/scanning"
)

var (
	undefinedVariable = errors.New("undefined variable")
)

type environment struct {
	enclosing *environment
	values    map[string]any
}

func newEnvironment(encosing *environment) *environment {
	return &environment{
		values:    make(map[string]any),
		enclosing: encosing,
	}
}

func (r *environment) define(name string, value any) {
	r.values[name] = value
}

func (r *environment) get(token *scanning.Token) (any, *internal.RuntimeError) {
	val, ok := r.values[token.Lexeme]
	if ok {
		return val, nil
	}
	if r.enclosing != nil {
		return r.enclosing.get(token)
	}
	return nil, &internal.RuntimeError{
		Error: undefinedVariable,
		Token: token,
	}
}

func (r *environment) assign(token *scanning.Token, value any) *internal.RuntimeError {
	if _, ok := r.values[token.Lexeme]; ok {
		r.values[token.Lexeme] = value
		return nil
	}
	if r.enclosing != nil {
		return r.enclosing.assign(token, value)
	}
	return &internal.RuntimeError{
		Error: undefinedVariable,
		Token: token,
	}
}
