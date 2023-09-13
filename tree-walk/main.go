package main

import (
	"fmt"
	"golox/tree-walk/parser"
	"golox/tree-walk/runtime"
	"golox/tree-walk/scanner"
	"os"
)

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

	if runtime.HadError {
		os.Exit(65)
	}
	if runtime.HadRuntimeError {
		os.Exit(70)
	}
}

func run(source string) {
	s := scanner.NewScanner(source)
	tokens := s.ScanTokens()
	_ = parser.NewParser(tokens)
}
