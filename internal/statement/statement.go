package statement

import (
	"gox/internal/expression"
	"gox/internal/scanning"
)

type Visitor interface {
	VisitForExpression(stmt *Expression) error
	VisitForPrint(stmt *Print) error
	VisitForVar(stmt *Var) error
}

type Stmt interface {
	Accept(visitor Visitor) error
}

// Expression
type Expression struct {
	Expression *expression.Expr
}

func (r *Expression) Accept(visitor Visitor) error {
	return visitor.VisitForExpression(r)
}

// Print
type Print struct {
	Expression *expression.Expr
}

func (r *Print) Accept(visitor Visitor) error {
	return visitor.VisitForPrint(r)
}

// Var
type Var struct {
	Name        *scanning.Token
	Initializer *expression.Expr
}

func (r *Var) Accept(visitor Visitor) error {
	return visitor.VisitForVar(r)
}
