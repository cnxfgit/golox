package interpreter

import (
	"fmt"
	"golox/tree-walk/expr"
	. "golox/tree-walk/object"
	"golox/tree-walk/rt"
	"golox/tree-walk/stmt"
	"golox/tree-walk/token"
	"strconv"
	"strings"
	"time"
)

type Interpreter struct {
	Globals     *rt.Environment
	Environment *rt.Environment
	locals      map[expr.Expr]int
}

func NewInterpreter() *Interpreter {
	globals := rt.NewEnvironment(nil)
	globals.Define("clock", NewNative(0, func(interpreter *Interpreter, arguments []Object) Object {
		return Number(time.Now().UnixNano() / int64(time.Millisecond))
	}))
	return &Interpreter{
		Globals:     globals,
		Environment: globals,
		locals:      make(map[expr.Expr]int),
	}
}

func (i *Interpreter) Interpret(statements []stmt.Stmt) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(rt.RuntimeError); ok {
				rt.ErrorRuntime(e)
			}
		}
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) execute(statement stmt.Stmt) {
	statement.Accept(i)
}

func (i *Interpreter) evaluate(expression expr.Expr) Object {
	return expression.Accept(i)
}

func (i *Interpreter) Resolve(expression expr.Expr, depth int) {
	i.locals[expression] = depth
}

func (i *Interpreter) ExecuteBlock(statements []stmt.Stmt, environment *rt.Environment) {
	previous := i.Environment

	defer func() {
		i.Environment = previous
	}()

	i.Environment = environment

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) VisitBlockStmt(stmt *stmt.Block) Object {
	i.ExecuteBlock(stmt.Statements, rt.NewEnvironment(i.Environment))
	return nil
}

func (i *Interpreter) VisitClassStmt(stmt *stmt.Class) Object {
	var superclass Object = nil
	if stmt.Superclass != nil {
		superclass = i.evaluate(stmt.Superclass)
		if _, ok := superclass.(*LoxClass); !ok {
			panic(rt.RuntimeError{Token: stmt.Superclass.Name, Message: "Superclass must be a class."})
		}
	}

	i.Environment.Define(stmt.Name.Lexeme, nil)

	if stmt.Superclass != nil {
		i.Environment = rt.NewEnvironment(i.Environment)
		i.Environment.Define("super", superclass)
	}

	methods := make(map[string]*LoxFunction)
	for _, method := range stmt.Methods {
		function := NewLoxFunction(method, i.Environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = function
	}
	sc := superclass.(*LoxClass)
	class := NewLoxClass(stmt.Name.Lexeme, sc, methods)

	if superclass != nil {
		i.Environment = i.Environment.Enclosing
	}

	i.Environment.Assign(stmt.Name, class)

	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *stmt.Expression) Object {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *stmt.Function) Object {
	function := NewLoxFunction(stmt, i.Environment, false)
	i.Environment.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *stmt.If) Object {
	if isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *stmt.Print) Object {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringify(value))

	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *stmt.Return) Object {
	var value Object = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	panic(rt.Return{Value: value})
}

func (i *Interpreter) VisitVarStmt(stmt *stmt.Var) Object {
	var value Object = nil
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.Environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *stmt.While) Object {
	for isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *expr.Assign) Object {
	value := i.evaluate(expr.Value)

	if distance, ok := i.locals[expr]; ok {
		i.Environment.AssignAt(distance, expr.Name, value)
	} else {
		i.Globals.Assign(expr.Name, value)
	}

	return value
}

func (i *Interpreter) VisitBinaryExpr(expr *expr.Binary) Object {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.Greater:
		checkNumberOperands(expr.Operator, left, right)
		l := left.(Number)
		r := right.(Number)
		return Boolean(l > r)
	case token.GreaterEqual:
		checkNumberOperands(expr.Operator, left, right)
		l := left.(Number)
		r := right.(Number)
		return Boolean(l > r)
	case token.Less:
		checkNumberOperands(expr.Operator, left, right)
		l := left.(Number)
		r := right.(Number)
		return Boolean(l < r)
	case token.LessEqual:
		checkNumberOperands(expr.Operator, left, right)
		l := left.(Number)
		r := right.(Number)
		return Boolean(l <= r)
	case token.Minus:
		checkNumberOperands(expr.Operator, left, right)
		l := left.(Number)
		r := right.(Number)
		return l - r
	case token.Plus:
		{
			l1, ok1 := left.(Number)
			r2, ok2 := right.(Number)
			if ok1 && ok2 {
				return l1 + r2
			}
			l3, ok1 := left.(String)
			l4, ok2 := right.(String)
			if ok1 && ok2 {
				return l3 + l4
			}
			panic(rt.RuntimeError{expr.Operator, "Operands must be two numbers or two strings."})
		}
	case token.Slash:
		checkNumberOperands(expr.Operator, left, right)
		l := left.(Number)
		r := right.(Number)
		return l / r
	case token.Star:
		checkNumberOperands(expr.Operator, left, right)
		l := left.(Number)
		r := right.(Number)
		return l * r
	case token.BangEqual:
		return Boolean(!isEqual(left, right))
	case token.EqualEqual:
		return Boolean(isEqual(left, right))
	}

	return nil
}

func (i *Interpreter) VisitCallExpr(expr *expr.Call) Object {
	callee := i.evaluate(expr.Callee)

	arguments := make([]Object, 0)
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		panic(rt.RuntimeError{expr.Paren, "Can only call functions and classes."})
	}

	if len(arguments) != function.Arity() {
		panic(rt.RuntimeError{Token: expr.Paren,
			Message: "Expected " + strconv.FormatInt(int64(function.Arity()), 64) +
				" arguments but got " + strconv.FormatInt(int64(len(arguments)), 64) + "."})
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitGetExpr(expr *expr.Get) Object {
	object := i.evaluate(expr.Object)
	if o, ok := object.(*LoxInstance); ok {
		return o.Get(expr.Name)
	}

	panic(rt.RuntimeError{Token: expr.Name, Message: "Only instances have properties."})
}

func (i *Interpreter) VisitGroupingExpr(expr *expr.Grouping) Object {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *expr.Literal) Object {
	return expr.Value
}

func (i *Interpreter) VisitLogicalExpr(expr *expr.Logical) Object {
	left := i.evaluate(expr.Left)

	if expr.Operator.Type == token.Or {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitSetExpr(expr *expr.Set) Object {
	object := i.evaluate(expr.Object)

	instance, ok := object.(*LoxInstance)
	if !ok {
		panic(rt.RuntimeError{expr.Name, "Only instances have fields."})
	}

	value := i.evaluate(expr.Value)
	instance.Set(expr.Name, value)
	return value
}

func (i *Interpreter) VisitSuperExpr(expr *expr.Super) Object {
	distance := i.locals[expr]
	superclass := i.Environment.GetAt(distance, "super").(*LoxClass)

	object := i.Environment.GetAt(distance-1, "this").(*LoxInstance)

	method := superclass.FindMethod(expr.Method.Lexeme)

	if method == nil {
		panic(rt.RuntimeError{expr.Method,
			"Undefined property '" + expr.Method.Lexeme + "'."})
	}

	return method.Bind(object)
}

func (i *Interpreter) VisitThisExpr(expr *expr.This) Object {
	return i.lookUpVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitUnaryExpr(expr *expr.Unary) Object {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.Bang:
		return Boolean(!isTruthy(right))
	case token.Minus:
		checkNumberOperand(expr.Operator, right)
		r := right.(Number)
		return -r
	}

	return nil
}

func (i *Interpreter) VisitVariableExpr(expr *expr.Variable) Object {
	return i.lookUpVariable(expr.Name, expr)
}

func checkNumberOperand(operator token.Token, operand Object) {
	if _, ok := operand.(Number); ok {
		return
	}

	panic(rt.RuntimeError{operator, "Operand must be a number."})
}

func (i *Interpreter) lookUpVariable(name token.Token, expr expr.Expr) Object {
	distance, ok := i.locals[expr]
	if ok {
		return i.Environment.GetAt(distance, name.Lexeme)
	} else {
		return i.Globals.Get(name)
	}
}

func isEqual(a Object, b Object) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func checkNumberOperands(operator token.Token, left Object, right Object) {
	_, ok1 := left.(Number)
	_, ok2 := right.(Number)
	if ok1 && ok2 {
		return
	}
	panic(rt.RuntimeError{Token: operator, Message: "Operands must be numbers."})
}

func stringify(object Object) string {
	if object == nil {
		return "nil"
	}

	if _, ok := object.(Number); ok {
		text := object.ToString()
		if strings.HasSuffix(text, ".0") {
			text = text[0 : len(text)-2]
		}
		return text
	}

	return object.ToString()
}

func isTruthy(object Object) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(Boolean); ok {
		return bool(b)
	}
	return true
}
