package expr

import (
	"golox/tree-walk/object"
	"golox/tree-walk/token"
)

type Expr interface {
	Accept(v Visitor) object.Object
}

type Visitor interface {
	visitAssignExpr(expr *Assign) object.Object
	visitBinaryExpr(expr *Binary) object.Object
	visitCallExpr(expr *Call) object.Object
	visitGetExpr(expr *Get) object.Object
	visitGroupingExpr(expr *Grouping) object.Object
	visitLiteralExpr(expr *Literal) object.Object
	visitLogicalExpr(expr *Logical) object.Object
	visitSetExpr(expr *Set) object.Object
	visitSuperExpr(expr *Super) object.Object
	visitThisExpr(expr *This) object.Object
	visitUnaryExpr(expr *Unary) object.Object
	visitVariableExpr(expr *Variable) object.Object
}

type Assign struct {
	name  token.Token
	value Expr
}

func (a *Assign) Accept(v Visitor) object.Object {
	return v.visitAssignExpr(a)
}

type Binary struct {
	left     Expr
	operator token.Token
	right    Expr
}

func (b *Binary) Accept(v Visitor) object.Object {
	return v.visitBinaryExpr(b)
}

type Call struct {
	callee    Expr
	paren     token.Token
	arguments []Expr
}

func (c *Call) Accept(v Visitor) object.Object {
	return v.visitCallExpr(c)
}

type Get struct {
	object Expr
	name   token.Token
}

func (g *Get) Accept(v Visitor) object.Object {
	return v.visitGetExpr(g)
}

type Grouping struct {
	expression Expr
}

func (g *Grouping) Accept(v Visitor) object.Object {
	return v.visitGroupingExpr(g)
}

type Literal struct {
	value object.Object
}

func (l *Literal) Accept(v Visitor) object.Object {
	return v.visitLiteralExpr(l)
}

type Logical struct {
	left     Expr
	operator token.Token
	right    Expr
}

func (l *Logical) Accept(v Visitor) object.Object {
	return v.visitLogicalExpr(l)
}

type Set struct {
	object Expr
	name   token.Token
	value  Expr
}

func (s *Set) Accept(v Visitor) object.Object {
	return v.visitSetExpr(s)
}

type Super struct {
	keyword token.Token
	method  token.Token
}

func (s *Super) Accept(v Visitor) object.Object {
	return v.visitSuperExpr(s)
}

type This struct {
	keyword token.Token
}

func (t *This) Accept(v Visitor) object.Object {
	return v.visitThisExpr(t)
}

type Unary struct {
	operator token.Token
	right    Expr
}

func (u *Unary) Accept(v Visitor) object.Object {
	return v.visitUnaryExpr(u)
}

type Variable struct {
	name token.Token
}

func (va *Variable) Accept(v Visitor) object.Object {
	return v.visitVariableExpr(va)
}
