package runtime

import (
	"errors"
	"gox/internal/scanning"
)

var (
	undefinedVariable = errors.New("undefined variable")
)

type environment struct {
	values map[string]any
}

func (r *environment) define(name string, value any) {
	r.values[name] = value
}

func (r *environment) get(token *scanning.Token) (any, error) {
	val, ok := r.values[token.Lexeme]
	if ok {
		return val, nil
	}
	return nil, &RuntimeError{
		error: errors.New("undefined variable: " + token.Lexeme),
		Token: token,
	}
}

func (r *environment) assign(token *scanning.Token, value any) error {
	if _, ok := r.values[token.Lexeme]; ok {
		r.values[token.Lexeme] = value
		return nil
	}
	return undefinedVariable
}
