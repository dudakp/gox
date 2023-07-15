package gox

import (
	"bufio"
	"fmt"
	"gox/internal/scanning"
	"os"
)

type Gox struct {
}

func (r *Gox) RunFile(path string) error {
	file, err := os.ReadFile(path)
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
	lexer := scanning.NewLexer(source)
	tokens, err := lexer.ScanTokens()
	if err != nil {
		Error(err.Line, err.Error(), "")
		return err
	}
	for _, token := range tokens {
		fmt.Println(token.String())
	}
	return nil
}

func Error(line int, message, where string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
}
