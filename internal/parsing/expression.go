package parsing

import "gox/internal/scanning"

type Visitor interface {
	visitForLiteral(expr *Literal) any
	visitForUnary(expr *Unary) any
	visitForBinary(expr *Binary) any
	visitForGrouping(expr *Grouping) any
}

// TODO: add some methods to this because like this it can be renamed to "Visitable"
type Expr interface {
	accept(visitor Visitor) any
}

type Literal struct {
	value interface{}
}

func (r *Literal) accept(visitor Visitor) any {
	return visitor.visitForLiteral(r)
}

type Unary struct {
	operator scanning.Token
	right    Expr
}

func (r *Unary) accept(visitor Visitor) any {
	return visitor.visitForUnary(r)
}

type Binary struct {
	left     Expr
	operator scanning.Token
	right    Expr
}

func (r *Binary) accept(visitor Visitor) any {
	return visitor.visitForBinary(r)
}

type Grouping struct {
	expression Expr
}

func (r *Grouping) accept(visitor Visitor) any {
	return visitor.visitForGrouping(r)
}
