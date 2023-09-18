package token

import (
	"fmt"
	"golox/tree-walk/object"
)

type TokenType int

const (
	LeftParen  TokenType = iota
	RightParen TokenType = iota
	LeftBrace  TokenType = iota
	RightBrace TokenType = iota
	Comma      TokenType = iota
	Dot        TokenType = iota
	Minus      TokenType = iota
	Plus       TokenType = iota
	Semicolon  TokenType = iota
	Slash      TokenType = iota
	Star       TokenType = iota

	Bang      TokenType = iota
	BangEqual TokenType = iota

	Equal        TokenType = iota
	EqualEqual   TokenType = iota
	Greater      TokenType = iota
	GreaterEqual TokenType = iota
	Less         TokenType = iota
	LessEqual    TokenType = iota

	Identifier TokenType = iota
	String     TokenType = iota
	Number     TokenType = iota

	And    TokenType = iota
	Class  TokenType = iota
	Else   TokenType = iota
	False  TokenType = iota
	Fun    TokenType = iota
	For    TokenType = iota
	If     TokenType = iota
	Nil    TokenType = iota
	Or     TokenType = iota
	Print  TokenType = iota
	Return TokenType = iota
	Super  TokenType = iota
	This   TokenType = iota
	True   TokenType = iota
	Var    TokenType = iota
	While  TokenType = iota

	Eof TokenType = iota
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal object.Object
	Line    uint
}

func NewToken(typ TokenType, lexeme string, literal object.Object, line uint) Token {
	return Token{
		typ,
		lexeme,
		literal,
		line,
	}
}

func (t *Token) toString() string {
	return fmt.Sprintf("%d %s %v", t.Type, t.Lexeme, t.Literal)
}
