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
	Statements []Stmt
}

func (b *Block) Accept(v Visitor) object.Object {
	return v.visitBlockStmt(b)
}

type Class struct {
	Name       token.Token
	Superclass *expr.Variable
	Methods    []*Function
}

func (c *Class) Accept(v Visitor) object.Object {
	return v.visitClassStmt(c)
}

type Expression struct {
	Expression expr.Expr
}

func (e *Expression) Accept(v Visitor) object.Object {
	return v.visitExpressionStmt(e)
}

type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

func (f *Function) Accept(v Visitor) object.Object {
	return v.visitFunctionStmt(f)
}

type If struct {
	Condition  expr.Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i *If) Accept(v Visitor) object.Object {
	return v.visitIfStmt(i)
}

type Print struct {
	Expression expr.Expr
}

func (p *Print) Accept(v Visitor) object.Object {
	return v.visitPrintStmt(p)
}

type Return struct {
	Keyword token.Token
	Value   expr.Expr
}

func (r *Return) Accept(v Visitor) object.Object {
	return v.visitReturnStmt(r)
}

type Var struct {
	Name        token.Token
	Initializer expr.Expr
}

func (va *Var) Accept(v Visitor) object.Object {
	return v.visitVarStmt(va)
}

type While struct {
	Condition expr.Expr
	Body      Stmt
}

func (w *While) Accept(v Visitor) object.Object {
	return v.visitWhileStmt(w)
}
