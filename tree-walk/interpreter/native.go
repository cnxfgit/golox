package interpreter

import "golox/tree-walk/object"

type Native struct {
	arity    int
	Function func(interpreter *Interpreter, arguments []object.Object) object.Object
}

func NewNative(arity int,
	function func(interpreter *Interpreter, arguments []object.Object) object.Object) *Native {
	return &Native{arity: arity, Function: function}
}

func (n *Native) Arity() int {
	return n.arity
}

func (n *Native) Call(interpreter *Interpreter, arguments []object.Object) object.Object {
	return n.Function(interpreter, arguments)
}

func (n *Native) ToString() string {
	return "<fn native>"
}
