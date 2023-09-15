package parser

import "golox/tree-walk/token"

type Parser struct {
	tokens  []token.Token
	current uint
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}
