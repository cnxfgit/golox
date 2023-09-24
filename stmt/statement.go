package stmt

import (
	"golox/expr"
	"golox/object"
	"golox/token"
)

type Stmt interface {
	Accept(v Visitor) object.Object
}

type Visitor interface {
	VisitBlockStmt(stmt *Block) object.Object
	VisitClassStmt(stmt *Class) object.Object
	VisitExpressionStmt(stmt *Expression) object.Object
	VisitFunctionStmt(stmt *Function) object.Object
	VisitIfStmt(stmt *If) object.Object
	VisitPrintStmt(stmt *Print) object.Object
	VisitReturnStmt(stmt *Return) object.Object
	VisitVarStmt(stmt *Var) object.Object
	VisitWhileStmt(stmt *While) object.Object
}

type Block struct {
	Statements []Stmt
}

func (b *Block) Accept(v Visitor) object.Object {
	return v.VisitBlockStmt(b)
}

type Class struct {
	Name       token.Token
	Superclass *expr.Variable
	Methods    []*Function
}

func (c *Class) Accept(v Visitor) object.Object {
	return v.VisitClassStmt(c)
}

type Expression struct {
	Expression expr.Expr
}

func (e *Expression) Accept(v Visitor) object.Object {
	return v.VisitExpressionStmt(e)
}

type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

func (f *Function) Accept(v Visitor) object.Object {
	return v.VisitFunctionStmt(f)
}

type If struct {
	Condition  expr.Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i *If) Accept(v Visitor) object.Object {
	return v.VisitIfStmt(i)
}

type Print struct {
	Expression expr.Expr
}

func (p *Print) Accept(v Visitor) object.Object {
	return v.VisitPrintStmt(p)
}

type Return struct {
	Keyword token.Token
	Value   expr.Expr
}

func (r *Return) Accept(v Visitor) object.Object {
	return v.VisitReturnStmt(r)
}

type Var struct {
	Name        token.Token
	Initializer expr.Expr
}

func (va *Var) Accept(v Visitor) object.Object {
	return v.VisitVarStmt(va)
}

type While struct {
	Condition expr.Expr
	Body      Stmt
}

func (w *While) Accept(v Visitor) object.Object {
	return v.VisitWhileStmt(w)
}
