package interpreter

import (
	"golox/tree-walk/object"
	"golox/tree-walk/rt"
	"golox/tree-walk/stmt"
)

type LoxFunction struct {
	declaration   *stmt.Function
	closure       *rt.Environment
	isInitializer bool
}

func NewLoxFunction(declaration *stmt.Function, closure *rt.Environment,
	isInitializer bool) *LoxFunction {
	return &LoxFunction{
		declaration,
		closure,
		isInitializer,
	}
}

func (f *LoxFunction) Bind(instance *LoxInstance) *LoxFunction {
	environment := rt.NewEnvironment(f.closure)
	environment.Define("this", instance)
	return NewLoxFunction(f.declaration, environment, f.isInitializer)
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) Call(inter *Interpreter, arguments []object.Object) (ret object.Object) {
	environment := rt.NewEnvironment(f.closure)
	for i := 0; i < len(f.declaration.Params); i++ {
		environment.Define(f.declaration.Params[i].Lexeme, arguments[i])
	}

	defer func() {
		if err := recover(); err != nil {
			if rv, ok := err.(rt.Return); ok {
				ret = rv.Value
			}
		}
	}()

	inter.ExecuteBlock(f.declaration.Body, environment)

	if f.isInitializer {
		return f.closure.GetAt(0, "this")
	}
	return nil
}

func (f *LoxFunction) ToString() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}
