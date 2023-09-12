package main

import (
	"fmt"
	"golox/tree-walk/scanner"
	"os"
)

var hadError bool = false
var hadRuntimeError bool = false

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runPrompt() {

}

func runFile(path string) {
	bytes, _ := os.ReadFile(path)
	source := string(bytes)

	run(source)

	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}
}

func run(source string) {
	_ = scanner.NewScanner(source)

}
