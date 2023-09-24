package interpreter

import (
	"golox/object"
	"golox/rt"
	"golox/token"
)

type LoxInstance struct {
	class  *LoxClass
	fields map[string]object.Object
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{class: class, fields: make(map[string]object.Object)}
}

func (i *LoxInstance) Get(name token.Token) object.Object {
	if object, ok := i.fields[name.Lexeme]; ok {
		return object
	}

	method := i.class.FindMethod(name.Lexeme)
	if method != nil {
		return method.Bind(i)
	}

	panic(rt.RuntimeError{Token: name, Message: "Undefined property '" + name.Lexeme + "'."})
}

func (i *LoxInstance) Set(name token.Token, value object.Object) {
	i.fields[name.Lexeme] = value
}

func (i *LoxInstance) ToString() string {
	return i.class.Name + " instance"
}
