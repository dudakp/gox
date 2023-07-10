package parsing

import "gox/internal/scanning"

type Visitor interface {
	VisitForLiteral(expr *Literal) any
	VisitForUnary(expr *Unary) any
	VisitForBinary(expr *Binary) any
	VisitForGrouping(expr *Grouping) any
}

type Expr interface {
	accept(visitor Visitor) any
}

type Literal struct {
	value interface{}
}

func (r *Literal) accept(visitor Visitor) any {
	return visitor.VisitForLiteral(r)
}

type Unary struct {
	operator scanning.Token
	right    Expr
}

func (r *Unary) accept(visitor Visitor) any {
	return visitor.VisitForUnary(r)
}

type Binary struct {
	Left    *Expr
	Operand scanning.Token
	Right   *Expr
}

func (r *Binary) accept(visitor Visitor) any {
	return visitor.VisitForBinary(r)
}

type Grouping struct {
	expression Expr
}

func (r *Grouping) accept(visitor Visitor) any {
	return visitor.VisitForGrouping(r)
}
