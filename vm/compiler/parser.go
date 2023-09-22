package compiler

import . "golox/vm/token"

var parser Parser = Parser{}

type Parser struct {
	current   Token
	previous  Token
	hadError  bool
	panicMode bool
}
