package compiler

import (
	"fmt"
	. "golox/vm/base"
	. "golox/vm/object"
	. "golox/vm/scanner"
	. "golox/vm/token"
	"math"
	"os"
)

type Compiler struct {
	current    *Compiler
	parser     *Parser
	enclosing  *Compiler
	function   *ObjFunction
	typ        FunctionType
	locals     [math.MaxUint8]local
	localCount int
	upvalues   [math.MaxUint8]upvalue
	scopeDepth int
	scanner    *Scanner
}

func NewCompiler(scanner *Scanner, typ FunctionType) *Compiler {
	c := &Compiler{}
	c.current = nil
	c.enclosing = c.current
	c.typ = typ
	c.scanner = scanner

	c.localCount = 0
	c.scopeDepth = 0

	c.function = NewFunction()
	c.current = c
	if typ != Script {
		c.current.function.Name = parser.previous.Message
	}

	lo := c.current.locals[c.localCount]
	lo.depth = 0
	lo.isCaptured = false
	c.current.localCount++
	if c.typ != Function {
		lo.name = NewToken(Identifier, 0, 4, 0, "this")
	} else {
		lo.name = NewToken(Identifier, 0, 0, 0, "")
	}
	return c
}

func (c *Compiler) declaration() {

}

func (c *Compiler) endCompiler() *ObjFunction {
	return nil
}

func (c *Compiler) Compile() *ObjFunction {
	parser.hadError = false
	parser.panicMode = false

	c.advance()

	for !c.match(Eof) {
		c.declaration()
	}

	objFunction := c.endCompiler()
	if parser.hadError {
		return nil
	}
	return objFunction
}

func (c *Compiler) match(typ TokenType) bool {
	if !c.check(typ) {
		return false
	}
	c.advance()
	return true
}

func (c *Compiler) check(typ TokenType) bool {
	return parser.current.Type == typ
}

func (c *Compiler) advance() {
	parser.previous = parser.current
	for {
		parser.current = c.scanner.ScanToken()
		if parser.current.Type != Error {
			break
		}
		errorAtCurrent(parser.current.Message)
	}
}

func errorAtCurrent(message string) {
	errorAt(parser.current, message)
}

func errorAt(token Token, message string) {
	if parser.panicMode {
		return
	}
	parser.panicMode = true
	_, _ = fmt.Fprintf(os.Stderr, "[line %d] Error\n", token.Line)

	if token.Type == Eof {
		_, _ = fmt.Fprintf(os.Stderr, " at end\n")
	} else if token.Type == Error {
		// Nothing.
	} else {
		_, _ = fmt.Fprintf(os.Stderr, " at \"%s\"\n", message)
	}

	_, _ = fmt.Fprintf(os.Stderr, ": %s\n", token.Message)
	parser.hadError = true
}
