package expr

import (
	. "golox/tree-walk/object"
	"golox/tree-walk/token"
)

type Expr interface {
	Accept(v Visitor) Object
}

type Visitor interface {
	visitAssignExpr(expr *Assign) Object
	visitBinaryExpr(expr *Binary) Object
	visitCallExpr(expr *Call) Object
	visitGetExpr(expr *Get) Object
	visitGroupingExpr(expr *Grouping) Object
	visitLiteralExpr(expr *Literal) Object
	visitLogicalExpr(expr *Logical) Object
	visitSetExpr(expr *Set) Object
	visitSuperExpr(expr *Super) Object
	visitThisExpr(expr *This) Object
	visitUnaryExpr(expr *Unary) Object
	visitVariableExpr(expr *Variable) Object
}

type Assign struct {
	Name  token.Token
	Value Expr
}

func (a *Assign) Accept(v Visitor) Object {
	return v.visitAssignExpr(a)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *Binary) Accept(v Visitor) Object {
	return v.visitBinaryExpr(b)
}

type Call struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func (c *Call) Accept(v Visitor) Object {
	return v.visitCallExpr(c)
}

type Get struct {
	Object Expr
	Name   token.Token
}

func (g *Get) Accept(v Visitor) Object {
	return v.visitGetExpr(g)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(v Visitor) Object {
	return v.visitGroupingExpr(g)
}

type Literal struct {
	Value Object
}

func (l *Literal) Accept(v Visitor) Object {
	return v.visitLiteralExpr(l)
}

type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (l *Logical) Accept(v Visitor) Object {
	return v.visitLogicalExpr(l)
}

type Set struct {
	Object Expr
	Name   token.Token
	Value  Expr
}

func (s *Set) Accept(v Visitor) Object {
	return v.visitSetExpr(s)
}

type Super struct {
	Keyword token.Token
	Method  token.Token
}

func (s *Super) Accept(v Visitor) Object {
	return v.visitSuperExpr(s)
}

type This struct {
	Keyword token.Token
}

func (t *This) Accept(v Visitor) Object {
	return v.visitThisExpr(t)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u *Unary) Accept(v Visitor) Object {
	return v.visitUnaryExpr(u)
}

type Variable struct {
	Name token.Token
}

func (va *Variable) Accept(v Visitor) Object {
	return v.visitVariableExpr(va)
}
