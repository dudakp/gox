package parsing

import (
	"errors"
	"gox/internal/scanning"
)

const (
	expectedSemicolonMsg = "expected ; after statement"
)

var (
	invalidAssignmentTarget = errors.New("invalid assignment target")
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

func (r *Parser) Parse() ([]*Stmt, *ParseError) {
	var res = make([]*Stmt, 0)
	for !r.isAtEnd() {
		stmt, err := r.declaration()
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

func (r *Parser) declaration() (Stmt, *TokenError) {
	if r.match(scanning.VAR) {
		declaration, tokenError := r.varDeclaration()
		if tokenError != nil {
			r.synchronize()
			// TODO: really dont return tokenErr?
			return nil, nil
		}
		return declaration, tokenError
	} else {
		return r.statement()
	}
}

func (r *Parser) varDeclaration() (Stmt, *TokenError) {
	identifier, tokenError := r.consume(scanning.IDENTIFIER, "expected variable name")
	if tokenError != nil {
		return nil, tokenError
	}

	if r.match(scanning.EQUAL) {
		initializer, err := r.expression()
		if err != nil {
			return nil, err
		}
		_, tokenError := r.consume(scanning.SEMICOLON, expectedSemicolonMsg)

		return &Var{
			Name:        identifier,
			Initializer: &initializer,
		}, tokenError
	}
	_, tokenError = r.consume(scanning.SEMICOLON, expectedSemicolonMsg)
	return &Var{
		Name: identifier,
	}, tokenError
}

func (r *Parser) statement() (Stmt, *TokenError) {
	if r.match(scanning.PRINT) {
		return r.printStatement()
	} else {
		return r.expressionStatement()
	}
}

func (r *Parser) expression() (Expr, *TokenError) {
	return r.assignment()
}

func (r *Parser) equality() (Expr, *TokenError) {
	expr, err := r.comparison()

	for r.match(scanning.BANG, scanning.BANG_EQUAL) {
		operator := r.previous()
		right, err := r.comparison()
		return &Binary{
			Left:     &expr,
			Operator: operator,
			Right:    &right,
		}, err
	}
	return expr, err
}

func (r *Parser) comparison() (Expr, *TokenError) {
	expr, err := r.term()

	for r.match(scanning.GREATER, scanning.GREATER_EQUAL, scanning.LESS, scanning.LESS_EQUAL) {
		operator := r.previous()
		right, err := r.term()
		return &Binary{
			Left:     &expr,
			Operator: operator,
			Right:    &right,
		}, err
	}
	return expr, err
}

func (r *Parser) term() (Expr, *TokenError) {
	expr, err := r.factor()

	for r.match(scanning.MINUS, scanning.PLUS) {
		operator := r.previous()
		right, err := r.term()
		return &Binary{
			Left:     &expr,
			Operator: operator,
			Right:    &right,
		}, err
	}
	return expr, err
}

func (r *Parser) factor() (Expr, *TokenError) {
	expr, err := r.unary()

	for r.match(scanning.SLASH, scanning.STAR) {
		operator := r.previous()
		right, err := r.term()
		return &Binary{
			Left:     &expr,
			Operator: operator,
			Right:    &right,
		}, err
	}
	return expr, err
}

func (r *Parser) unary() (Expr, *TokenError) {
	for r.match(scanning.BANG, scanning.MINUS) {
		operator := r.previous()
		right, err := r.unary()
		return &Unary{
			Operator: operator,
			Right:    &right,
		}, err
	}
	return r.primary()
}

func (r *Parser) primary() (Expr, *TokenError) {
	if r.match(scanning.FALSE) {
		return &Literal{Value: false}, nil
	}
	if r.match(scanning.TRUE) {
		return &Literal{Value: true}, nil
	}
	if r.match(scanning.NIL) {
		return &Literal{Value: nil}, nil
	}

	if r.match(scanning.NUMBER, scanning.STRING) {
		return &Literal{Value: r.previous().Literal}, nil
	}

	if r.match(scanning.IDENTIFIER) {
		prev := r.previous()
		return &VarExpr{
			Name: prev,
		}, nil
	}

	if r.match(scanning.LEFT_PAREN) {
		expr, tokenErr := r.expression()
		_, tokenErr = r.consume(scanning.RIGHT_PAREN, "expected ) after Expression")
		if tokenErr != nil {
			return nil, tokenErr
		}
		return &Grouping{Expression: &expr}, nil
	}
	return nil, nil
}

func (r *Parser) consume(t scanning.TokenType, message string) (*scanning.Token, *TokenError) {
	if r.check(t) {
		return r.advance(), nil
	}
	peek := r.peek()
	return nil, &TokenError{
		error: errors.New(message),
		Token: peek,
	}
}

func (r *Parser) consumeExpression() (Expr, *TokenError) {
	value, tokenError := r.expression()
	if tokenError != nil {
		return nil, tokenError
	}
	_, tokenError = r.consume(scanning.SEMICOLON, expectedSemicolonMsg)
	if tokenError != nil {
		return nil, tokenError
	}
	return value, nil
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

func (r *Parser) previous() *scanning.Token {
	// TODO: add some check for out of bounds error
	return &r.tokens[r.current-1]
}

func (r *Parser) isAtEnd() bool {
	return r.peek().TokenType == scanning.EOF
}

func (r *Parser) advance() *scanning.Token {
	if !r.isAtEnd() {
		r.current++
	}
	return r.previous()
}

func (r *Parser) peek() *scanning.Token {
	return &r.tokens[r.current]
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

func (r *Parser) printStatement() (Stmt, *TokenError) {
	expr, err := r.consumeExpression()
	return &Print{
		Expression: &expr,
	}, err
}

func (r *Parser) expressionStatement() (Stmt, *TokenError) {
	expr, err := r.consumeExpression()
	return &Expression{
		Expression: &expr,
	}, err
}

func (r *Parser) assignment() (Expr, *TokenError) {
	expr, tokenError := r.equality()
	if tokenError != nil {
		return nil, tokenError
	}
	if r.match(scanning.EQUAL) {
		eq := r.previous()
		value, tokenError := r.assignment()
		if tokenError != nil {
			return nil, tokenError
		}

		if _, ok := expr.(*VarExpr); ok {
			name := expr.(*VarExpr).Name
			return &Assign{
				Name:  name,
				Value: value,
			}, nil
		} else {
			return nil, &TokenError{
				error: invalidAssignmentTarget,
				Token: eq,
			}
		}
	}
	return expr, nil
}
