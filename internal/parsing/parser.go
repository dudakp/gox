package parsing

import (
	"gox/internal/scanning"
)

type Parser struct {
	tokens  []scanning.Token
	current int
}

func NewParser(tokens []scanning.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (r *Parser) expression() Expr {
	return r.equality()
}

func (r *Parser) equality() Expr {
	expr := r.comparison()

	for r.match(scanning.BANG, scanning.BANG_EQUAL) {
		operator := r.previous()
		right := r.comparison()
		return &Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

func (r *Parser) comparison() Expr {
	expr := r.term()

	for r.match(scanning.GREATER, scanning.GREATER_EQUAL, scanning.LESS, scanning.LESS_EQUAL) {
		operator := r.previous()
		right := r.term()
		return &Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

func (r *Parser) term() Expr {
	expr := r.factor()

	for r.match(scanning.MINUS, scanning.PLUS) {
		operator := r.previous()
		right := r.term()
		return &Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

func (r *Parser) factor() Expr {
	expr := r.unary()

	for r.match(scanning.SLASH, scanning.STAR) {
		operator := r.previous()
		right := r.term()
		return &Binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr
}

func (r *Parser) unary() Expr {
	for r.match(scanning.BANG, scanning.MINUS) {
		operator := r.previous()
		right := r.unary()
		return &Binary{
			operator: operator,
			right:    right,
		}
	}
	return r.primary()
}

func (r *Parser) primary() Expr {
	if r.match(scanning.FALSE) {
		return &Literal{value: false}
	}
	if r.match(scanning.TRUE) {
		return &Literal{value: true}
	}
	if r.match(scanning.NIL) {
		return &Literal{value: nil}
	}

	if r.match(scanning.NUMBER, scanning.STRING) {
		return &Literal{value: r.previous().Literal}
	}

	if r.match(scanning.LEFT_PAREN) {
		expr := r.expression()
		r.consume(scanning.RIGHT_PAREN, "expected ) after expression")
		return &Grouping{expression: expr}
	}
	return nil
}

func (r *Parser) consume(t scanning.TokenType, message string) (scanning.Token, error) {
	if r.check(t) {
		return r.advance(), nil
	}
	return scanning.Token{}, errors.
}

func (r *Parser) match(types ...scanning.TokenType) bool {
	for _, t := range types {
		if r.check(t) {
			r.advance()
			return true
		}
	}
	return false
}

func (r *Parser) check(t scanning.TokenType) bool {
	if r.isAtEnd() {
		return false
	} else {
		return r.peek().TokenType == t
	}
}

func (r *Parser) previous() scanning.Token {
	// TODO: add some check for out of bounds error
	return r.tokens[r.current-1]
}

func (r *Parser) isAtEnd() bool {
	return r.peek().TokenType == scanning.EOF
}

func (r *Parser) advance() scanning.Token {
	if !r.isAtEnd() {
		r.current++
	}
	return r.previous()
}

func (r *Parser) peek() scanning.Token {
	return r.tokens[r.current]
}
