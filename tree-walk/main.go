package main

import (
	"fmt"
	"golox/tree-walk/interpreter"
	"golox/tree-walk/parser"
	"golox/tree-walk/resolver"
	"golox/tree-walk/rt"
	"golox/tree-walk/scan"
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

	if rt.HadError {
		os.Exit(65)
	}
	if rt.HadRuntimeError {
		os.Exit(70)
	}
}

func run(source string) {
	i := interpreter.NewInterpreter()
	s := scan.NewScanner(source)
	tokens := s.ScanTokens()
	p := parser.NewParser(tokens)
	statements := p.Parse()

	if rt.HadError {
		return
	}

	r := resolver.NewResolver(i)
	r.Resolve(statements)

	if rt.HadError {
		return
	}

	i.Interpret(statements)
}
