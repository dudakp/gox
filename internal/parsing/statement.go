package parsing

import (
	"gox/internal/scanning"
)

type StmtVisitor interface {
	VisitForExpression(stmt *Expression) error
	VisitForPrint(stmt *Print) error
	VisitForVar(stmt *Var) error
}

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

// Expression
type Expression struct {
	Expression *Expr
}

func (r *Expression) Accept(visitor StmtVisitor) error {
	return visitor.VisitForExpression(r)
}

// Print
type Print struct {
	Expression *Expr
}

func (r *Print) Accept(visitor StmtVisitor) error {
	return visitor.VisitForPrint(r)
}

// Var
type Var struct {
	Name        *scanning.Token
	Initializer *Expr
}

func (r *Var) Accept(visitor StmtVisitor) error {
	return visitor.VisitForVar(r)
}
