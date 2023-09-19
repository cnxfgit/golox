package interpreter

import "golox/tree-walk/object"

type LoxClass struct {
	Name       string
	Superclass *LoxClass
	Methods    map[string]*LoxFunction
}

func NewLoxClass(Name string, Superclass *LoxClass, Methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		Name:       Name,
		Superclass: Superclass,
		Methods:    Methods,
	}
}

func (c *LoxClass) FindMethod(name string) *LoxFunction {
	if method, ok := c.Methods[name]; ok {
		return method
	}

	if c.Superclass != nil {
		return c.Superclass.FindMethod(name)
	}

	return nil
}

func (c *LoxClass) ToString() string {
	return c.Name
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []object.Object) object.Object {
	instance := NewLoxInstance(c)
	initializer := c.FindMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}
	return instance
}

func (c *LoxClass) Arity() int {
	initializer := c.FindMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}
