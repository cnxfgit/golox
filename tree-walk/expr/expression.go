package expr

import (
	. "golox/tree-walk/object"
	"golox/tree-walk/token"
)

type Expr interface {
	Accept(v Visitor) Object
}

type Visitor interface {
	VisitAssignExpr(expr *Assign) Object
	VisitBinaryExpr(expr *Binary) Object
	VisitCallExpr(expr *Call) Object
	VisitGetExpr(expr *Get) Object
	VisitGroupingExpr(expr *Grouping) Object
	VisitLiteralExpr(expr *Literal) Object
	VisitLogicalExpr(expr *Logical) Object
	VisitSetExpr(expr *Set) Object
	VisitSuperExpr(expr *Super) Object
	VisitThisExpr(expr *This) Object
	VisitUnaryExpr(expr *Unary) Object
	VisitVariableExpr(expr *Variable) Object
}

type Assign struct {
	Name  token.Token
	Value Expr
}

func (a *Assign) Accept(v Visitor) Object {
	return v.VisitAssignExpr(a)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (b *Binary) Accept(v Visitor) Object {
	return v.VisitBinaryExpr(b)
}

type Call struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func (c *Call) Accept(v Visitor) Object {
	return v.VisitCallExpr(c)
}

type Get struct {
	Object Expr
	Name   token.Token
}

func (g *Get) Accept(v Visitor) Object {
	return v.VisitGetExpr(g)
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(v Visitor) Object {
	return v.VisitGroupingExpr(g)
}

type Literal struct {
	Value Object
}

func (l *Literal) Accept(v Visitor) Object {
	return v.VisitLiteralExpr(l)
}

type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (l *Logical) Accept(v Visitor) Object {
	return v.VisitLogicalExpr(l)
}

type Set struct {
	Object Expr
	Name   token.Token
	Value  Expr
}

func (s *Set) Accept(v Visitor) Object {
	return v.VisitSetExpr(s)
}

type Super struct {
	Keyword token.Token
	Method  token.Token
}

func (s *Super) Accept(v Visitor) Object {
	return v.VisitSuperExpr(s)
}

type This struct {
	Keyword token.Token
}

func (t *This) Accept(v Visitor) Object {
	return v.VisitThisExpr(t)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u *Unary) Accept(v Visitor) Object {
	return v.VisitUnaryExpr(u)
}

type Variable struct {
	Name token.Token
}

func (va *Variable) Accept(v Visitor) Object {
	return v.VisitVariableExpr(va)
}
