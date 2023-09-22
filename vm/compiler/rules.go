package compiler

var rules []parseRule

func init() {
	rules = []parseRule{}
}

type Precedence int

const (
	NONE       Precedence = iota
	ASSIGNMENT Precedence = iota
	OR         Precedence = iota
	AND        Precedence = iota
	EQUALITY   Precedence = iota
	COMPARISON Precedence = iota
	TERM       Precedence = iota
	FACTOR     Precedence = iota
	UNARY      Precedence = iota
	CALL       Precedence = iota
	PRIMARY    Precedence = iota
)

type parseRule struct {
	prefix     func(canAssign bool) // 前缀
	infix      func(canAssign bool) // 中缀
	precedence Precedence           // 优先级
}
