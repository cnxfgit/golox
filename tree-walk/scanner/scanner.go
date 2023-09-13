package scanner

import (
	"golox/tree-walk/object"
	"golox/tree-walk/runtime"
	"strconv"
)

var keywords map[string]TokenType

func init() {
	keywords = make(map[string]TokenType)
	keywords["and"] = And
	keywords["class"] = Class
	keywords["else"] = Else
	keywords["false"] = False
	keywords["for"] = For
	keywords["fun"] = Fun
	keywords["if"] = If
	keywords["nil"] = Nil
	keywords["or"] = Or
	keywords["print"] = Print
	keywords["return"] = Return
	keywords["super"] = Super
	keywords["this"] = This
	keywords["true"] = True
	keywords["var"] = Var
	keywords["while"] = While
}

type Scanner struct {
	source  string
	tokens  []Token
	start   uint
	current uint
	line    uint
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, NewToken(Eof, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addTokenTyp(LeftParen)
	case ')':
		s.addTokenTyp(RightParen)
	case '{':
		s.addTokenTyp(LeftBrace)
	case '}':
		s.addTokenTyp(RightBrace)
	case ',':
		s.addTokenTyp(Comma)
	case '.':
		s.addTokenTyp(Dot)
	case '-':
		s.addTokenTyp(Minus)
	case '+':
		s.addTokenTyp(Plus)
	case ';':
		s.addTokenTyp(Semicolon)
	case '*':
		s.addTokenTyp(Star)
	case '!':
		{
			if s.match('=') {
				s.addTokenTyp(BangEqual)
			} else {
				s.addTokenTyp(Bang)
			}
		}
	case '=':
		if s.match('=') {
			s.addTokenTyp(EqualEqual)
		} else {
			s.addTokenTyp(Equal)
		}
	case '<':
		if s.match('=') {
			s.addTokenTyp(LessEqual)
		} else {
			s.addTokenTyp(Less)
		}
	case '>':
		if s.match('=') {
			s.addTokenTyp(GreaterEqual)
		} else {
			s.addTokenTyp(Greater)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addTokenTyp(Slash)
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
			runtime.ErrorLine(s.line, "Unexpected character.")
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
		typ = Identifier
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
	s.addToken(Number, num)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		runtime.ErrorLine(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addToken(String, value)
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

func (s *Scanner) addTokenTyp(typ TokenType) {
	s.addToken(typ, nil)
}

func (s *Scanner) addToken(typ TokenType, literal object.Object) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(typ, text, literal, s.line))
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
