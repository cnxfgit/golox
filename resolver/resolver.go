package resolver

import (
	"golox/expr"
	"golox/interpreter"
	"golox/object"
	"golox/rt"
	"golox/stmt"
	"golox/token"
)

type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          *stack
	currentFunction functionType
	currentClass    classType
}

type functionType int

type classType int

const (
	None        = iota
	Function    = iota
	Initializer = iota
	Method      = iota

	Class    = iota
	Subclass = iota
)

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          newStack(),
		currentFunction: None,
		currentClass:    None,
	}
}

func (r *Resolver) Resolve(statements []stmt.Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) resolveStmt(stmt stmt.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr expr.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveFunction(function *stmt.Function, funcType functionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = funcType

	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.Resolve(function.Body)
	r.endScope()

	r.currentFunction = enclosingFunction
}

func (r *Resolver) beginScope() {
	r.scopes.push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes.pop()
}

func (r *Resolver) declare(name token.Token) {
	if r.scopes.isEmpty() {
		return
	}

	scope, _ := r.scopes.peek()
	if _, ok := scope[name.Lexeme]; ok {
		rt.ErrorToken(name, "Already a variable with this name in this scope.")
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name token.Token) {
	if r.scopes.isEmpty() {
		return
	}
	scope, _ := r.scopes.peek()
	scope[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr expr.Expr, name token.Token) {
	for i := r.scopes.len() - 1; i >= 0; i-- {
		if _, ok := r.scopes.get(i)[name.Lexeme]; ok {
			r.interpreter.Resolve(expr, r.scopes.len()-1-i)
			return
		}
	}
}

func (r *Resolver) VisitBlockStmt(stmt *stmt.Block) object.Object {
	r.beginScope()
	r.Resolve(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitClassStmt(stmt *stmt.Class) object.Object {
	enclosingClass := r.currentClass
	r.currentClass = Class

	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.Superclass != nil && stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
		rt.ErrorToken(stmt.Superclass.Name, "A class can't inherit from itself.")
	}

	if stmt.Superclass != nil {
		r.currentClass = Subclass
		r.resolveExpr(stmt.Superclass)
	}

	if stmt.Superclass != nil {
		r.beginScope()
		scope, _ := r.scopes.peek()
		scope["super"] = true
	}

	r.beginScope()
	scope, _ := r.scopes.peek()
	scope["this"] = true

	for _, method := range stmt.Methods {
		declaration := Method
		if method.Name.Lexeme == "init" {
			declaration = Initializer
		}
		r.resolveFunction(method, functionType(declaration))
	}

	r.endScope()

	if stmt.Superclass != nil {
		r.endScope()
	}

	r.currentClass = enclosingClass
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *stmt.Expression) object.Object {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *stmt.Function) object.Object {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, functionType(Function))
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *stmt.If) object.Object {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *stmt.Print) object.Object {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *stmt.Return) object.Object {
	if r.currentFunction == None {
		rt.ErrorToken(stmt.Keyword, "Can't return from top-level code.")
	}

	if stmt.Value != nil {
		if r.currentFunction == Initializer {
			rt.ErrorToken(stmt.Keyword, "Can't return a value from an initializer.")
		}
		r.resolveExpr(stmt.Value)
	}

	return nil
}

func (r *Resolver) VisitVarStmt(stmt *stmt.Var) object.Object {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *stmt.While) object.Object {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *expr.Assign) object.Object {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *expr.Binary) object.Object {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *expr.Call) object.Object {
	r.resolveExpr(expr.Callee)

	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}
	return nil
}

func (r *Resolver) VisitGetExpr(expr *expr.Get) object.Object {
	r.resolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *expr.Grouping) object.Object {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *expr.Literal) object.Object {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *expr.Logical) object.Object {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitSetExpr(expr *expr.Set) object.Object {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitSuperExpr(expr *expr.Super) object.Object {
	if r.currentClass == None {
		rt.ErrorToken(expr.Keyword, "Can't use 'super' outside of a class.")
	} else if r.currentClass != Subclass {
		rt.ErrorToken(expr.Keyword, "Can't use 'super' in a class with no superclass.")
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitThisExpr(expr *expr.This) object.Object {
	if r.currentClass == None {
		rt.ErrorToken(expr.Keyword, "Can't use 'this' outside of a class.")
		return nil
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *expr.Unary) object.Object {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *expr.Variable) object.Object {
	scope, _ := r.scopes.peek()
	v, ok := scope[expr.Name.Lexeme]
	if !r.scopes.isEmpty() && ok && v == false {
		rt.ErrorToken(expr.Name, "Can't read local variable in its own initializer.")
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}
