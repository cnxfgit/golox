package main

import (
	"bufio"
	"fmt"
	. "golox/interpreter"
	"golox/parser"
	"golox/resolver"
	"golox/rt"
	"golox/scan"
	"os"
)

var interpreter = NewInterpreter()

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
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		ok := scanner.Scan()
		line := scanner.Text()
		if !ok {
			break
		}
		run(line)
		rt.HadError = false
	}
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

	s := scan.NewScanner(source)
	tokens := s.ScanTokens()
	p := parser.NewParser(tokens)
	statements := p.Parse()

	if rt.HadError {
		return
	}

	r := resolver.NewResolver(interpreter)
	r.Resolve(statements)

	if rt.HadError {
		return
	}

	interpreter.Interpret(statements)
}
