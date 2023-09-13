package stmt

import (
	"golox/tree-walk/object"
	"golox/tree-walk/scanner"
)

type Stmt interface {
	Accept(v Visitor) object.Object
}

type Visitor interface {
	visitBlockStmt(stmt *Block) object.Object
	visitClassStmt(stmt *Class) object.Object
	//visitExpressionStmt(stmt *Expression) object.Object
	//visitFunctionStmt(stmt *Function) object.Object
	//visitIfStmt(stmt *If) object.Object
	//visitPrintStmt(stmt *Print) object.Object
	//visitReturnStmt(stmt *Return) object.Object
	//visitVarStmt(stmt *Var) object.Object
	//visitWhileStmt(stmt *While) object.Object
}

type Block struct {
	statements []Stmt
}

func (b *Block) Accept(v Visitor) object.Object {
	return v.visitBlockStmt(b)
}

type Class struct {
	name scanner.Token
}
