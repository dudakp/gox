package parsing

import (
	"errors"
	"gox/internal/expression"
	"gox/internal/scanning"
	"gox/internal/statement"
)

type ParseError struct {
	TokenError
}

type TokenError struct {
	error
	Token *scanning.Token
}

// TODO: write tests
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

func (r *Parser) Parse() ([]*statement.Stmt, *ParseError) {
	var res = make([]*statement.Stmt, len(r.tokens))
	for !r.isAtEnd() {
		stmt, err := r.statement()
		if err != nil {
			return nil, &ParseError{
				TokenError{
					error: err,
					Token: err.Token,
				}}
		}
		res = append(res, &stmt)
	}
	return res, nil
}

func (r *Parser) statement() (statement.Stmt, *TokenError) {
	if r.match(scanning.PRINT) {
		return r.printStatement()
	} else {
		return r.expressionStatement()
	}
}

func (r *Parser) expression() (expression.Expr, *TokenError) {
	return r.equality()
}

func (r *Parser) equality() (expression.Expr, *TokenError) {
	expr, err := r.comparison()

	for r.match(scanning.BANG, scanning.BANG_EQUAL) {
		operator := r.previous()
		right, err := r.comparison()
		return &expression.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}, err
	}
	return expr, err
}

func (r *Parser) comparison() (expression.Expr, *TokenError) {
	expr, err := r.term()

	for r.match(scanning.GREATER, scanning.GREATER_EQUAL, scanning.LESS, scanning.LESS_EQUAL) {
		operator := r.previous()
		right, err := r.term()
		return &expression.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}, err
	}
	return expr, err
}

func (r *Parser) term() (expression.Expr, *TokenError) {
	expr, err := r.factor()

	for r.match(scanning.MINUS, scanning.PLUS) {
		operator := r.previous()
		right, err := r.term()
		return &expression.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}, err
	}
	return expr, err
}

func (r *Parser) factor() (expression.Expr, *TokenError) {
	expr, err := r.unary()

	for r.match(scanning.SLASH, scanning.STAR) {
		operator := r.previous()
		right, err := r.term()
		return &expression.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}, err
	}
	return expr, err
}

func (r *Parser) unary() (expression.Expr, *TokenError) {
	for r.match(scanning.BANG, scanning.MINUS) {
		operator := r.previous()
		right, err := r.unary()
		return &expression.Unary{
			Operator: operator,
			Right:    right,
		}, err
	}
	return r.primary()
}

func (r *Parser) primary() (expression.Expr, *TokenError) {
	if r.match(scanning.FALSE) {
		return &expression.Literal{Value: false}, nil
	}
	if r.match(scanning.TRUE) {
		return &expression.Literal{Value: true}, nil
	}
	if r.match(scanning.NIL) {
		return &expression.Literal{Value: nil}, nil
	}

	if r.match(scanning.NUMBER, scanning.STRING) {
		return &expression.Literal{Value: r.previous().Literal}, nil
	}

	if r.match(scanning.LEFT_PAREN) {
		expr, tokenErr := r.expression()
		_, tokenErr = r.consume(scanning.RIGHT_PAREN, "expected ) after Expression")
		if tokenErr != nil {
			return nil, tokenErr
		}
		return &expression.Grouping{Expression: expr}, nil
	}
	return nil, nil
}

func (r *Parser) consume(t scanning.TokenType, message string) (scanning.Token, *TokenError) {
	if r.check(t) {
		return r.advance(), nil
	}
	peek := r.peek()
	return scanning.Token{}, &TokenError{
		error: errors.New(message),
		Token: &peek,
	}
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

func (r *Parser) synchronize() {
	r.advance()

	for !r.isAtEnd() {
		if r.previous().TokenType == scanning.SEMICOLON {
			return
		}
		switch r.peek().TokenType {
		case scanning.CLASS:
		case scanning.FUN:
		case scanning.VAR:
		case scanning.FOR:
		case scanning.IF:
		case scanning.WHILE:
		case scanning.PRINT:
		case scanning.RETURN:
			return
		}
	}
	r.advance()
}

func (r *Parser) printStatement() (statement.Stmt, *TokenError) {
	expr, err := r.consumeExpression()
	return &statement.Print{
		Expression: &expr,
	}, err
}

func (r *Parser) expressionStatement() (statement.Stmt, *TokenError) {
	expr, err := r.consumeExpression()
	return &statement.Expression{
		Expression: &expr,
	}, err
}

func (r *Parser) consumeExpression() (expression.Expr, *TokenError) {
	value, tokenError := r.expression()
	if tokenError != nil {
		return nil, tokenError
	}
	_, tokenError = r.consume(scanning.SEMICOLON, "expected ; after statement")
	if tokenError != nil {
		return nil, tokenError
	}
	return value, nil
}
