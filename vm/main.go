package main

import (
	"bufio"
	"fmt"
	"golox/vm/base"
	"golox/vm/rt"
	"golox/vm/vm"
	"os"
)

func main() {
	rt.Vm = vm.NewVm()

	args := os.Args
	if len(args) == 1 {
		repl()
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: golox [path]")
		os.Exit(64)
	}

}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		ok := scanner.Scan()
		line := scanner.Text()
		if !ok {
			break
		}
		rt.Vm.Interpret(line)
	}
}

func runFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Could not read file %s.\n", path)
		os.Exit(74)
	}
	source := string(file)

	result := rt.Vm.Interpret(source)
	switch result {
	case base.CompileError:
		os.Exit(65)
	case base.RuntimeError:
		os.Exit(70)
	}
}
