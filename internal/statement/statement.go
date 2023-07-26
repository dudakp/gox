package statement

import "gox/internal/expression"

type Visitor interface {
	VisitForExpression(stmt *Expression) error
	VisitForPrint(stmt *Print) error
}

type Stmt interface {
	Accept(visitor Visitor) error
}

type Expression struct {
	Expression *expression.Expr
}

func (r *Expression) Accept(visitor Visitor) error {
	return visitor.VisitForExpression(r)
}

type Print struct {
	Expression *expression.Expr
}

func (r *Print) Accept(visitor Visitor) error {
	return visitor.VisitForPrint(r)
}
