package expr

import "golox/tree-walk/object"

type Expr interface {
	Accept(v Visitor) object.Object
}

type Visitor interface {
	//visitAssignExpr(Assign expr) object.Object
	//visitBinaryExpr(Binary expr) object.Object
	//visitCallExpr(Call expr) object.Object
	//visitGetExpr(Get expr) object.Object
	//visitGroupingExpr(Grouping expr) object.Object
	//visitLiteralExpr(Literal expr) object.Object
	//visitLogicalExpr(Logical expr) object.Object
	//visitSetExpr(Set expr) object.Object
	//visitSuperExpr(Super expr) object.Object
	//visitThisExpr(This expr) object.Object
	//visitUnaryExpr(Unary expr) object.Object
	//visitVariableExpr(Variable expr) object.Object
}
