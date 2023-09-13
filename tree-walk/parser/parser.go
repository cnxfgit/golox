package parser

import "golox/tree-walk/scanner"

type Parser struct {
	tokens  []scanner.Token
	current uint
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}
