package vm

import (
	. "golox/vm/base"
	. "golox/vm/compiler"
	. "golox/vm/scanner"
)

type Vm struct {
}

func NewVm() *Vm {
	return &Vm{}
}

func (v *Vm) Interpret(source string) InterpretResult {
	compiler := NewCompiler(NewScanner(source), Script)
	_ = compiler.Compile()

	return Ok
}
