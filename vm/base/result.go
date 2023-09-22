package base

type InterpretResult int

const (
	Ok           InterpretResult = iota
	CompileError InterpretResult = iota
	RuntimeError InterpretResult = iota
)

type FunctionType int

const (
	Function    FunctionType = iota
	Initializer FunctionType = iota
	Method      FunctionType = iota
	Script      FunctionType = iota
)
