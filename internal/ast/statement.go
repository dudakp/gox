package ast

import (
	"gox/internal"
	"gox/internal/scanning"
)

type StmtVisitor interface {
	VisitForExpression(stmt *Expression) *internal.RuntimeError
	VisitForPrint(stmt *Print) *internal.RuntimeError
	VisitForVar(stmt *Var) *internal.RuntimeError
	VisitForBlock(block *Block) *internal.RuntimeError
}
type Stmt interface {
	Accept(visitor StmtVisitor) *internal.RuntimeError
}

// Expression
type Expression struct {
	Expression *Expr
}

func (r *Expression) Accept(visitor StmtVisitor) *internal.RuntimeError {
	return visitor.VisitForExpression(r)
}

// Print
type Print struct {
	Expression *Expr
}

func (r *Print) Accept(visitor StmtVisitor) *internal.RuntimeError {
	return visitor.VisitForPrint(r)
}

// Var
type Var struct {
	Name        *scanning.Token
	Initializer *Expr
}

func (r *Var) Accept(visitor StmtVisitor) *internal.RuntimeError {
	return visitor.VisitForVar(r)
}

// Block
type Block struct {
	Statements []*Stmt
}

func (r *Block) Accept(visitor StmtVisitor) *internal.RuntimeError {
	return visitor.VisitForBlock(r)
}
