package interpreter

import "golox/tree-walk/object"

type LoxCallable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []object.Object) object.Object
}
