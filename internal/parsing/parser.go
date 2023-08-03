package parsing

import (
	"errors"
	ast2 "gox/internal/ast"
	"gox/internal/scanning"
)

const (
	expectedSemicolonMsg       = "expected ; after statement"
	missingRightParenMsg       = "expected ) after expression"
	varNameExpectedMsg         = "expected variable name"
	expectedRightBraceMsg      = "expected } after block"
	missingLeftParenAfterIfMsg = "expected ( after 'if' "
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

func (r *Parser) Parse() ([]*ast2.Stmt, *ParseError) {
	var res = make([]*ast2.Stmt, 0)
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

func (r *Parser) declaration() (ast2.Stmt, *TokenError) {
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

func (r *Parser) varDeclaration() (ast2.Stmt, *TokenError) {
	identifier, tokenError := r.consume(scanning.IDENTIFIER, varNameExpectedMsg)
	if tokenError != nil {
		return nil, tokenError
	}

	if r.match(scanning.EQUAL) {
		initializer, err := r.expression()
		if err != nil {
			return nil, err
		}
		_, tokenError := r.consume(scanning.SEMICOLON, expectedSemicolonMsg)

		return &ast2.Var{
			Name:        identifier,
			Initializer: &initializer,
		}, tokenError
	}
	_, tokenError = r.consume(scanning.SEMICOLON, expectedSemicolonMsg)
	return &ast2.Var{
		Name: identifier,
	}, tokenError
}

func (r *Parser) statement() (ast2.Stmt, *TokenError) {
	if r.match(scanning.IF) {
		return r.ifStatement()
	}
	if r.match(scanning.PRINT) {
		return r.printStatement()
	}
	if r.match(scanning.LEFT_BRACE) {
		return r.block()
	}
	return r.expressionStatement()

}

func (r *Parser) expression() (ast2.Expr, *TokenError) {
	return r.assignment()
}

func (r *Parser) equality() (ast2.Expr, *TokenError) {
	expr, err := r.comparison()

	for r.match(scanning.EQUAL_EQUAL, scanning.BANG_EQUAL) {
		operator := r.previous()
		right, err := r.comparison()
		return &ast2.Binary{
			Left:     &expr,
			Operator: operator,
			Right:    &right,
		}, err
	}
	return expr, err
}

func (r *Parser) comparison() (ast2.Expr, *TokenError) {
	expr, err := r.term()

	for r.match(scanning.GREATER, scanning.GREATER_EQUAL, scanning.LESS, scanning.LESS_EQUAL) {
		operator := r.previous()
		right, err := r.term()
		return &ast2.Binary{
			Left:     &expr,
			Operator: operator,
			Right:    &right,
		}, err
	}
	return expr, err
}

func (r *Parser) term() (ast2.Expr, *TokenError) {
	expr, err := r.factor()

	for r.match(scanning.MINUS, scanning.PLUS) {
		operator := r.previous()
		right, err := r.term()
		return &ast2.Binary{
			Left:     &expr,
			Operator: operator,
			Right:    &right,
		}, err
	}
	return expr, err
}

func (r *Parser) factor() (ast2.Expr, *TokenError) {
	expr, err := r.unary()

	for r.match(scanning.SLASH, scanning.STAR) {
		operator := r.previous()
		right, err := r.term()
		return &ast2.Binary{
			Left:     &expr,
			Operator: operator,
			Right:    &right,
		}, err
	}
	return expr, err
}

func (r *Parser) unary() (ast2.Expr, *TokenError) {
	for r.match(scanning.BANG, scanning.MINUS) {
		operator := r.previous()
		right, err := r.unary()
		return &ast2.Unary{
			Operator: operator,
			Right:    &right,
		}, err
	}
	return r.primary()
}

func (r *Parser) primary() (ast2.Expr, *TokenError) {
	if r.match(scanning.FALSE) {
		return &ast2.Literal{Value: false}, nil
	}
	if r.match(scanning.TRUE) {
		return &ast2.Literal{Value: true}, nil
	}
	if r.match(scanning.NIL) {
		return &ast2.Literal{Value: nil}, nil
	}

	if r.match(scanning.NUMBER, scanning.STRING) {
		return &ast2.Literal{Value: r.previous().Literal}, nil
	}

	if r.match(scanning.IDENTIFIER) {
		prev := r.previous()
		return &ast2.VarExpr{
			Name: prev,
		}, nil
	}

	if r.match(scanning.LEFT_PAREN) {
		expr, tokenErr := r.expression()
		_, tokenErr = r.consume(scanning.RIGHT_PAREN, missingRightParenMsg)
		if tokenErr != nil {
			return nil, tokenErr
		}
		return &ast2.Grouping{Expression: &expr}, nil
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

func (r *Parser) consumeExpression() (ast2.Expr, *TokenError) {
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

func (r *Parser) printStatement() (ast2.Stmt, *TokenError) {
	expr, err := r.consumeExpression()
	return &ast2.Print{
		Expression: &expr,
	}, err
}

func (r *Parser) expressionStatement() (ast2.Stmt, *TokenError) {
	expr, err := r.consumeExpression()
	return &ast2.Expression{
		Expression: &expr,
	}, err
}

func (r *Parser) assignment() (ast2.Expr, *TokenError) {
	expr, tokenError := r.or()
	if tokenError != nil {
		return nil, tokenError
	}
	if r.match(scanning.EQUAL) {
		eq := r.previous()
		value, tokenError := r.assignment()
		if tokenError != nil {
			return nil, tokenError
		}

		if _, ok := expr.(*ast2.VarExpr); ok {
			name := expr.(*ast2.VarExpr).Name
			return &ast2.Assign{
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

func (r *Parser) or() (ast2.Expr, *TokenError) {
	expr, tokenError := r.and()
	if tokenError != nil {
		return nil, tokenError
	}
	for r.match(scanning.OR) {
		operator := r.previous()
		right, tokenErr := r.and()
		if tokenErr != nil {
			return nil, tokenErr
		}
		return &ast2.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}, nil
	}
	return expr, nil
}

func (r *Parser) and() (ast2.Expr, *TokenError) {
	expr, tokenError := r.equality()
	if tokenError != nil {
		return nil, tokenError
	}
	for r.match(scanning.AND) {
		operator := r.previous()
		right, tokenErr := r.equality()
		if tokenErr != nil {
			return nil, tokenErr
		}
		return &ast2.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}, nil
	}
	return expr, nil
}

func (r *Parser) block() (*ast2.Block, *TokenError) {
	res := make([]*ast2.Stmt, 0)
	for !r.check(scanning.RIGHT_BRACE) && !r.isAtEnd() {
		declaration, tokenError := r.declaration()
		if tokenError != nil {
			return nil, tokenError
		}
		res = append(res, &declaration)
	}

	_, err := r.consume(scanning.RIGHT_BRACE, expectedRightBraceMsg)
	if err != nil {
		return nil, err
	}
	return &ast2.Block{Statements: res}, nil
}

func (r *Parser) ifStatement() (ast2.Stmt, *TokenError) {
	_, err := r.consume(scanning.LEFT_PAREN, missingLeftParenAfterIfMsg)
	if err != nil {
		return nil, err
	}
	condition, err := r.expression()
	if err != nil {
		return nil, err
	}
	_, err = r.consume(scanning.RIGHT_PAREN, missingRightParenMsg)
	if err != nil {
		return nil, err
	}

	var elseStmt ast2.Stmt
	then, err := r.statement()
	if err != nil {
		return nil, err
	}
	if r.match(scanning.ELSE) {
		elseStmt, err = r.statement()
	}

	return &ast2.If{
		Condition: condition,
		Then:      then,
		Else:      elseStmt,
	}, nil
}
