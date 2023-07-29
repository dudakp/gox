package parsing

import "gox/internal/scanning"

type ExprVisitor interface {
	VisitForLiteral(expr *Literal) (any, error)
	VisitForUnary(expr *Unary) (any, error)
	VisitForBinary(expr *Binary) (any, error)
	VisitForGrouping(expr *Grouping) (any, error)
	VisitForVariableExpression(expr *VarExpr) (any, error)
	VisitForAssignExpression(expr *Assign) (any, error)
}

type Expr interface {
	Accept(visitor ExprVisitor) (any, error)
}

// Literal
type Literal struct {
	Value interface{}
}

func (r *Literal) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitForLiteral(r)
}

// Unary
type Unary struct {
	Operator *scanning.Token
	Right    *Expr
}

func (r *Unary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitForUnary(r)
}

// Binary
type Binary struct {
	Left     *Expr
	Operator *scanning.Token
	Right    *Expr
}

func (r *Binary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitForBinary(r)
}

// Grouping
type Grouping struct {
	Expression *Expr
}

func (r *Grouping) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitForGrouping(r)
}

// VarExpr
type VarExpr struct {
	Name *scanning.Token
}

func (r *VarExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitForVariableExpression(r)
}

// Assign
type Assign struct {
	Name  *scanning.Token
	Value Expr
}

func (r *Assign) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitForAssignExpression(r)
}
