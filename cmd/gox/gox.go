package gox

import (
	"bufio"
	"fmt"
	"gox/internal/parsing"
	"gox/internal/runtime"
	"gox/internal/scanning"
	"os"
	"path/filepath"
)

type Gox struct {
	Interpreter *runtime.Interpreter
}

func (r *Gox) RunFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	file, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Println("An error occurred while reading file. Please try again", err)
		return err
	}
	err = r.run(string(file))
	if err != nil {
		_ = fmt.Errorf("unable to run script: %w", err)
		os.Exit(65)
	}
	return nil
}

func (r *Gox) Repl() {
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		if err.Error() == "EOF" {
			fmt.Println("Bye!")
		}
		fmt.Println("An error occurred while reading input. Please try again", err)
	}
	err = r.run(input)
}

// run makes necessary calls to execute the source code
func (r *Gox) run(source string) error {
	tokens, syntaxErr := scanning.NewLexer(source).ScanTokens()
	if syntaxErr != nil {
		ReportError(syntaxErr.Line, syntaxErr.Error(), "")
		return syntaxErr
	}
	ast, parseErr := parsing.NewParser(tokens).Parse()
	if parseErr != nil {
		ReportParseError(parseErr)
		return parseErr
	}
	interpreterErr := r.Interpreter.Interpret(ast)
	if interpreterErr != nil {
		ReportError(interpreterErr.Token.Line, interpreterErr.Error(), "")
	}
	return nil
}

func ReportError(line int, message, where string) {
	fmt.Printf("[line %d] error %s: %s\n", line, where, message)
}

func ReportParseError(parseError *parsing.ParseError) {
	token := parseError.Token
	if token.TokenType == scanning.EOF {
		ReportError(token.Line, parseError.Error(), "")
	} else {
		ReportError(token.Line, parseError.Error(), token.Lexeme)
	}
}
