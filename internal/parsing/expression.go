package parsing

import "gox/internal/scanning"

type Visitor interface {
	VisitForLiteral(expr *Literal) any
	VisitForUnary(expr *Unary) any
	VisitForBinary(expr *Binary) any
	VisitForGrouping(expr *Grouping) any
}

// TODO: add some methods to this because like this it can be renamed to "Visitable"
type Expr interface {
	Accept(visitor Visitor) any
}

type Literal struct {
	Value interface{}
}

func (r *Literal) Accept(visitor Visitor) any {
	return visitor.VisitForLiteral(r)
}

type Unary struct {
	Operator scanning.Token
	Right    Expr
}

func (r *Unary) Accept(visitor Visitor) any {
	return visitor.VisitForUnary(r)
}

type Binary struct {
	Left     Expr
	Operator scanning.Token
	Right    Expr
}

func (r *Binary) Accept(visitor Visitor) any {
	return visitor.VisitForBinary(r)
}

type Grouping struct {
	expression Expr
}

func (r *Grouping) Accept(visitor Visitor) any {
	return visitor.VisitForGrouping(r)
}
