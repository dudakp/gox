package parsing

import "gox/internal/scanning"

type Visitor interface {
	VisitForLiteral(expr *Literal) (any, error)
	VisitForUnary(expr *Unary) (any, error)
	VisitForBinary(expr *Binary) (any, error)
	VisitForGrouping(expr *Grouping) (any, error)
}

// TODO: add some methods to this because like this it can be renamed to "Visitable"
type Expr interface {
	Accept(visitor Visitor) (any, error)
}

type Literal struct {
	Value interface{}
}

func (r *Literal) Accept(visitor Visitor) (any, error) {
	return visitor.VisitForLiteral(r)
}

type Unary struct {
	Operator scanning.Token
	Right    Expr
}

func (r *Unary) Accept(visitor Visitor) (any, error) {
	return visitor.VisitForUnary(r)
}

type Binary struct {
	Left     Expr
	Operator scanning.Token
	Right    Expr
}

func (r *Binary) Accept(visitor Visitor) (any, error) {
	return visitor.VisitForBinary(r)
}

type Grouping struct {
	Expression Expr
}

func (r *Grouping) Accept(visitor Visitor) (any, error) {
	return visitor.VisitForGrouping(r)
}
