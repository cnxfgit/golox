package scan

import (
	"golox/tree-walk/object"
	"golox/tree-walk/rt"
	"golox/tree-walk/token"
	"strconv"
)

var keywords map[string]token.TokenType

func init() {
	keywords = make(map[string]token.TokenType)
	keywords["and"] = token.And
	keywords["class"] = token.Class
	keywords["else"] = token.Else
	keywords["false"] = token.False
	keywords["for"] = token.For
	keywords["fun"] = token.Fun
	keywords["if"] = token.If
	keywords["nil"] = token.Nil
	keywords["or"] = token.Or
	keywords["print"] = token.Print
	keywords["return"] = token.Return
	keywords["super"] = token.Super
	keywords["this"] = token.This
	keywords["true"] = token.True
	keywords["var"] = token.Var
	keywords["while"] = token.While
}

type Scanner struct {
	source  string
	tokens  []token.Token
	start   uint
	current uint
	line    uint
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  make([]token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, token.NewToken(token.Eof, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addTokenTyp(token.LeftParen)
	case ')':
		s.addTokenTyp(token.RightParen)
	case '{':
		s.addTokenTyp(token.LeftBrace)
	case '}':
		s.addTokenTyp(token.RightBrace)
	case ',':
		s.addTokenTyp(token.Comma)
	case '.':
		s.addTokenTyp(token.Dot)
	case '-':
		s.addTokenTyp(token.Minus)
	case '+':
		s.addTokenTyp(token.Plus)
	case ';':
		s.addTokenTyp(token.Semicolon)
	case '*':
		s.addTokenTyp(token.Star)
	case '!':
		{
			if s.match('=') {
				s.addTokenTyp(token.BangEqual)
			} else {
				s.addTokenTyp(token.Bang)
			}
		}
	case '=':
		if s.match('=') {
			s.addTokenTyp(token.EqualEqual)
		} else {
			s.addTokenTyp(token.Equal)
		}
	case '<':
		if s.match('=') {
			s.addTokenTyp(token.LessEqual)
		} else {
			s.addTokenTyp(token.Less)
		}
	case '>':
		if s.match('=') {
			s.addTokenTyp(token.GreaterEqual)
		} else {
			s.addTokenTyp(token.Greater)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addTokenTyp(token.Slash)
		}
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			rt.ErrorLine(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	typ, ok := keywords[text]
	if !ok {
		typ = token.Identifier
	}

	s.addTokenTyp(typ)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}
	num, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addToken(token.Number, num)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		rt.ErrorLine(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addToken(token.String, value)
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= uint(len(s.source)) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) addTokenTyp(typ token.TokenType) {
	s.addToken(typ, nil)
}

func (s *Scanner) addToken(typ token.TokenType, literal object.Object) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(typ, text, literal, s.line))
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= uint(len(s.source))
}

func (s *Scanner) advance() byte {
	return s.source[s.current]
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
