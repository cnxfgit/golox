package stmt

import (
	"golox/tree-walk/expr"
	"golox/tree-walk/object"
	"golox/tree-walk/token"
)

type Stmt interface {
	Accept(v Visitor) object.Object
}

type Visitor interface {
	visitBlockStmt(stmt *Block) object.Object
	visitClassStmt(stmt *Class) object.Object
	visitExpressionStmt(stmt *Expression) object.Object
	visitFunctionStmt(stmt *Function) object.Object
	visitIfStmt(stmt *If) object.Object
	visitPrintStmt(stmt *Print) object.Object
	visitReturnStmt(stmt *Return) object.Object
	visitVarStmt(stmt *Var) object.Object
	visitWhileStmt(stmt *While) object.Object
}

type Block struct {
	statements []Stmt
}

func (b *Block) Accept(v Visitor) object.Object {
	return v.visitBlockStmt(b)
}

type Class struct {
	name       token.Token
	superclass expr.Variable
	methods    []Function
}

func (c *Class) Accept(v Visitor) object.Object {
	return v.visitClassStmt(c)
}

type Expression struct {
	expression expr.Expr
}

func (e *Expression) Accept(v Visitor) object.Object {
	return v.visitExpressionStmt(e)
}

type Function struct {
	name   token.Token
	params []token.Token
	body   []Stmt
}

func (f *Function) Accept(v Visitor) object.Object {
	return v.visitFunctionStmt(f)
}

type If struct {
	condition  expr.Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (i *If) Accept(v Visitor) object.Object {
	return v.visitIfStmt(i)
}

type Print struct {
	expression expr.Expr
}

func (p *Print) Accept(v Visitor) object.Object {
	return v.visitPrintStmt(p)
}

type Return struct {
	keyword token.Token
	value   expr.Expr
}

func (r *Return) Accept(v Visitor) object.Object {
	return v.visitReturnStmt(r)
}

type Var struct {
	name        token.Token
	initializer expr.Expr
}

func (va *Var) Accept(v Visitor) object.Object {
	return v.visitVarStmt(va)
}

type While struct {
	condition expr.Expr
	body      Stmt
}
