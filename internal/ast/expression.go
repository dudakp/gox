package ast

import (
	"gox/internal"
	"gox/internal/scanning"
)

type ExprVisitor interface {
	VisitForLiteral(expr *Literal) (any, *internal.RuntimeError)
	VisitForUnary(expr *Unary) (any, *internal.RuntimeError)
	VisitForBinary(expr *Binary) (any, *internal.RuntimeError)
	VisitForGrouping(expr *Grouping) (any, *internal.RuntimeError)
	VisitForVariableExpression(expr *VarExpr) (any, *internal.RuntimeError)
	VisitForAssignExpression(expr *Assign) (any, *internal.RuntimeError)
	VisitForLogical(expr *Logical) (any, *internal.RuntimeError)
	VisitForFunctionCall(expr *Call) (any, *internal.RuntimeError)
}

type Expr interface {
	Accept(visitor ExprVisitor) (any, *internal.RuntimeError)
}

// Literal
type Literal struct {
	Value interface{}
}

func (r *Literal) Accept(visitor ExprVisitor) (any, *internal.RuntimeError) {
	return visitor.VisitForLiteral(r)
}

// Unary
type Unary struct {
	Operator *scanning.Token
	Right    *Expr
}

func (r *Unary) Accept(visitor ExprVisitor) (any, *internal.RuntimeError) {
	return visitor.VisitForUnary(r)
}

// Binary
type Binary struct {
	Left     *Expr
	Operator *scanning.Token
	Right    *Expr
}

func (r *Binary) Accept(visitor ExprVisitor) (any, *internal.RuntimeError) {
	return visitor.VisitForBinary(r)
}

// Grouping
type Grouping struct {
	Expression *Expr
}

func (r *Grouping) Accept(visitor ExprVisitor) (any, *internal.RuntimeError) {
	return visitor.VisitForGrouping(r)
}

// VarExpr
type VarExpr struct {
	Name *scanning.Token
}

func (r *VarExpr) Accept(visitor ExprVisitor) (any, *internal.RuntimeError) {
	return visitor.VisitForVariableExpression(r)
}

// Assign
type Assign struct {
	Name  *scanning.Token
	Value Expr
}

func (r *Assign) Accept(visitor ExprVisitor) (any, *internal.RuntimeError) {
	return visitor.VisitForAssignExpression(r)
}

// Logical
type Logical struct {
	Left     Expr
	Operator *scanning.Token
	Right    Expr
}

func (r *Logical) Accept(visitor ExprVisitor) (any, *internal.RuntimeError) {
	return visitor.VisitForLogical(r)
}

// Call
type Call struct {
	Callee Expr
	Paren  *scanning.Token
	Params []Expr
}

func (r *Call) Accept(visitor ExprVisitor) (any, *internal.RuntimeError) {
	return visitor.VisitForFunctionCall(r)
}
