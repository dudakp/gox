package main

import (
	"fmt"
	"gox/cmd/gox"
	"os"
)

func main() {
	interpreter := gox.Gox{}

	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Too many arguments")
		os.Exit(64)
	}

	var err error
	if len(args) == 0 {
		for {
			interpreter.Repl()
		}
	} else {
		err = interpreter.RunFile(args[0])
	}
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		os.Exit(65)
	}
}
