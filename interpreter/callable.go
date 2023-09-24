package interpreter

import (
	"golox/object"
)

type LoxCallable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []object.Object) object.Object
}
