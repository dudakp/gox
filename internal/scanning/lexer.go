package scanning

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	UnexpectedCharacter = errors.New("unexpected character")
	UnterminatedString  = errors.New("unterminated string")
	InvalidNumber       = errors.New("invalid number")

	reserved = map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
)

type SyntaxError struct {
	error
	Line int
}

type Lexer struct {
	Source  string
	tokens  []Token
	start   int // start of lexeme
	current int // current character of lexeme being scanned
	line    int
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		Source: source,
		line:   1,
	}
}

func (r *Lexer) ScanTokens() ([]Token, *SyntaxError) {
	var err error
	for {
		if r.isAtEnd() {
			r.tokens = append(r.tokens, Token{EOF, "", nil, r.line})
			break
		}
		r.start = r.current
		err = r.scanToken()
		if err != nil {
			return nil, &SyntaxError{
				error: err,
				Line:  r.line,
			}
		}
	}
	return r.tokens, nil
}

func (r *Lexer) scanToken() error {
	c := r.advance()
	switch c {
	// simple operators
	case '(':
		r.addSimpleToken(LEFT_PAREN)
		break
	case ')':
		r.addSimpleToken(RIGHT_PAREN)
		break
	case '{':
		r.addSimpleToken(LEFT_BRACE)
		break
	case '}':
		r.addSimpleToken(RIGHT_BRACE)
		break
	case ',':
		r.addSimpleToken(COMMA)
		break
	case '.':
		r.addSimpleToken(DOT)
		break
	case '-':
		r.addSimpleToken(MINUS)
		break
	case '+':
		r.addSimpleToken(PLUS)
		break
	case ';':
		r.addSimpleToken(SEMICOLON)
		break
	case '*':
		r.addSimpleToken(STAR)
		break
		// multi-token operators
	case '!':
		r.addSimpleToken(r.matchReturn("=", BANG_EQUAL, BANG))
		break
	case '=':
		r.addSimpleToken(r.matchReturn("=", EQUAL_EQUAL, EQUAL))
		break
	case '<':
		r.addSimpleToken(r.matchReturn("=", LESS_EQUAL, LESS))
		break
	case '>':
		r.addSimpleToken(r.matchReturn("=", GREATER_EQUAL, GREATER))
		break
	case '/':
		// single line comment
		if r.match("/") {
			for r.peek() != "\n" && !r.isAtEnd() {
				r.advance()
			}
		} else {
			r.addSimpleToken(SLASH)
		}

	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
		break
	case '\n':
		r.line++
		break
	case '"':
		stringLiteral, syntaxError := r.string()
		if syntaxError != nil {
			return syntaxError
		}
		r.addToken(STRING, stringLiteral)
		break
	default:
		if unicode.IsDigit(rune(c)) {
			numberLiteral, err := r.number()
			if err != nil {
				return err
			}
			r.addToken(NUMBER, numberLiteral)
		} else if unicode.IsLetter(rune(c)) {
			r.addSimpleToken(r.identifier())
		} else {
			return UnexpectedCharacter
		}
	}
	return nil
}

func (r *Lexer) identifier() TokenType {
	for unicode.IsLetter(r.peekAsRune()) || unicode.IsDigit(r.peekAsRune()) {
		r.advance()
	}
	literal := r.Source[r.start:r.current]
	keyword, ok := reserved[literal]
	if !ok {
		return IDENTIFIER
	} else {
		return keyword
	}
}

func (r *Lexer) number() (float64, error) {
	for unicode.IsDigit(r.peekAsRune()) {
		r.advance()
	}
	if r.peek() == "." && unicode.IsDigit(rune(r.peekNext()[0])) {
		r.advance()
		for unicode.IsDigit(r.peekAsRune()) {
			r.advance()
		}
	}
	literal := r.Source[r.start:r.current]
	float, err := strconv.ParseFloat(literal, 6)
	if err != nil {
		return 0.0, errors.Join(InvalidNumber, err)
	}
	return float, nil
}

func (r *Lexer) string() (string, error) {
	for r.peek() != "\"" && !r.isAtEnd() {
		if r.peek() == "\n" {
			r.line++
		}
		r.advance()
	}

	if r.isAtEnd() {
		return "", UnterminatedString
	}

	// scan tle closing "
	r.advance()

	return r.clean(r.Source[r.start+1 : r.current-1]), nil
}

func (r *Lexer) advance() byte {
	r.current++
	return r.Source[r.current-1]
}

func (r *Lexer) peekAsRune() rune {
	return rune(r.peek()[0])
}

func (r *Lexer) peek() string {
	if r.isAtEnd() {
		return "\000"
	}
	return string(r.Source[r.current])
}

func (r *Lexer) peekNext() string {
	if r.current+1 >= len(r.Source) {
		return "\000"
	}
	return string(r.Source[r.current+1])
}

func (r *Lexer) matchReturn(expected string, ifTrue, ifFalse TokenType) TokenType {
	if r.match(expected) {
		return ifTrue
	} else {
		return ifFalse
	}
}

func (r *Lexer) match(expected string) bool {
	if r.isAtEnd() {
		return false
	}
	if string(r.Source[r.current]) != expected {
		return false
	}
	r.current++
	return true
}

func (r *Lexer) addSimpleToken(t TokenType) {
	r.addToken(t, nil)
}

func (r *Lexer) addToken(t TokenType, literal interface{}) {
	text := r.Source[r.start:r.current]
	r.tokens = append(r.tokens, Token{
		TokenType: t,
		Lexeme:    r.clean(text),
		Literal:   literal,
		Line:      r.line,
	})
}

func (r *Lexer) isAtEnd() bool {
	return r.current >= len(r.Source)
}

func (r *Lexer) clean(s string) string {
	return strings.Trim(strings.TrimSpace(s), "\"")
}
