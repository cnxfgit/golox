package token

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

	Bang         TokenType = iota
	BangEqual    TokenType = iota
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
	For    TokenType = iota
	Fun    TokenType = iota
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
	Error  TokenType = iota
	Eof    TokenType = iota
)

type Token struct {
	Type    TokenType
	Start   int
	Length  int
	Line    int
	Message string
}

func NewToken(typ TokenType, start int, length int, line int, message string) Token {
	return Token{
		Type:    typ,
		Start:   start,
		Length:  length,
		Line:    line,
		Message: message,
	}
}

func ErrorToken(message string, line int) Token {
	return Token{
		Line:    line,
		Message: message,
		Length:  len(message),
		Type:    Error,
	}
}
