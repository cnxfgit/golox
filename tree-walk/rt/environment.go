package rt

import (
	. "golox/tree-walk/object"
	"golox/tree-walk/token"
)

type Environment struct {
	Enclosing *Environment
	values    map[string]Object
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{Enclosing: enclosing, values: make(map[string]Object)}
}

func (e *Environment) Get(name token.Token) Object {
	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	panic(RuntimeError{Token: name, Message: "Undefined variable '" + name.Lexeme + "'."})
}

func (e *Environment) Assign(name token.Token, value Object) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return
	}

	if e.Enclosing != nil {
		e.Enclosing.Assign(name, value)
		return
	}

	panic(RuntimeError{Token: name, Message: "Undefined variable '" + name.Lexeme + "'."})
}

func (e *Environment) Define(name string, value Object) {
	e.values[name] = value
}

func (e *Environment) Ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.Enclosing
	}
	return environment
}

func (e *Environment) GetAt(distance int, name string) Object {
	return e.Ancestor(distance).values[name]
}

func (e *Environment) AssignAt(distance int, name token.Token, value Object) {
	e.Ancestor(distance).values[name.Lexeme] = value
}
