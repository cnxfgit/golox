package scanner

import . "golox/vm/token"

type Scanner struct {
	start   int
	current int
	line    int
	source  string
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
	}
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	b := s.source[s.current]
	s.current++
	return b
}
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}
func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
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

func (s *Scanner) skipWhitespace() {
	for {
		c := s.peek()
		switch c {
		case ' ':
		case '\r':
		case '\t':
			s.advance()
		case '\n':
			s.line++
			s.advance()
		case '/':
			if s.peekNext() == '/' {
				// A comment goes until the end of the line.
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

func (s *Scanner) checkKeyword(rest string, typ TokenType) TokenType {
	if len(rest) == s.current-s.start && rest == s.source[s.start:s.start+len(rest)] {
		return typ
	}
	return Identifier
}

func (s *Scanner) identifierType() TokenType {
	switch s.source[s.start] {
	case 'a':
		return s.checkKeyword("and", And)
	case 'c':
		return s.checkKeyword("class", Class)
	case 'e':
		return s.checkKeyword("else", Else)
	case 'f':
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'a':
				return s.checkKeyword("false", False)
			case 'o':
				return s.checkKeyword("for", For)
			case 'u':
				return s.checkKeyword("fun", Fun)
			}
		}
	case 'i':
		return s.checkKeyword("if", If)
	case 'n':
		return s.checkKeyword("nil", Nil)
	case 'o':
		return s.checkKeyword("or", Or)
	case 'p':
		return s.checkKeyword("print", Print)
	case 'r':
		return s.checkKeyword("return", Return)
	case 's':
		return s.checkKeyword("super", Super)
	case 't':
		if s.current-s.start > 1 {
			switch s.source[s.start+1] {
			case 'h':
				return s.checkKeyword("s", This)
			case 'r':
				return s.checkKeyword("true", True)
			}
		}
	case 'v':
		return s.checkKeyword("var", Var)
	case 'w':
		return s.checkKeyword("while", While)
	}
	return Identifier
}

func (s *Scanner) identifier() Token {
	for isAlpha(s.peek()) || isDigit(s.peek()) {
		s.advance()
	}
	return NewToken(s.identifierType(), s.start, s.current-s.start, s.line, s.source[s.start:s.current])
}

func (s *Scanner) number() Token {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the ".".
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	return NewToken(Number, s.start, s.current, s.line, s.source[s.start:s.current])
}

func (s *Scanner) string() Token {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return ErrorToken("Unterminated string.", s.line)
	}

	// The closing quote.
	s.advance()
	return NewToken(String, s.start, s.current-s.start-2, s.line, s.source[s.start+1:s.current-1])
}

func (s *Scanner) ScanToken() Token {
	s.skipWhitespace()

	s.start = s.current
	// 重新标记扫描仪起点并检查源代码是否结束
	if s.isAtEnd() {
		return NewToken(Eof, s.start, s.current, s.line, s.source[s.start:s.current])
	}

	c := s.advance()
	if isAlpha(c) {
		return s.identifier()
	}
	if isDigit(c) {
		return s.number()
	}

	switch c {
	case '(':
		return NewToken(LeftParen, s.start, s.current, s.line, s.source[s.start:s.current])
	case ')':
		return NewToken(RightParen, s.start, s.current, s.line, s.source[s.start:s.current])
	case '{':
		return NewToken(LeftBrace, s.start, s.current, s.line, s.source[s.start:s.current])
	case '}':
		return NewToken(RightBrace, s.start, s.current, s.line, s.source[s.start:s.current])
	case ';':
		return NewToken(Semicolon, s.start, s.current, s.line, s.source[s.start:s.current])
	case ',':
		return NewToken(Comma, s.start, s.current, s.line, s.source[s.start:s.current])
	case '.':
		return NewToken(Dot, s.start, s.current, s.line, s.source[s.start:s.current])
	case '-':
		return NewToken(Minus, s.start, s.current, s.line, s.source[s.start:s.current])
	case '+':
		return NewToken(Plus, s.start, s.current, s.line, s.source[s.start:s.current])
	case '/':
		return NewToken(Slash, s.start, s.current, s.line, s.source[s.start:s.current])
	case '*':
		return NewToken(Star, s.start, s.current, s.line, s.source[s.start:s.current])
	case '!':
		if s.match('=') {
			return NewToken(BangEqual, s.start, s.current, s.line, s.source[s.start:s.current])
		}
		return NewToken(Bang, s.start, s.current, s.line, s.source[s.start:s.current])
	case '=':
		if s.match('=') {
			return NewToken(EqualEqual, s.start, s.current, s.line, s.source[s.start:s.current])
		}
		return NewToken(Equal, s.start, s.current, s.line, s.source[s.start:s.current])
	case '<':
		if s.match('=') {
			return NewToken(LessEqual, s.start, s.current, s.line, s.source[s.start:s.current])
		}
		return NewToken(Less, s.start, s.current, s.line, s.source[s.start:s.current])
	case '>':
		if s.match('=') {
			return NewToken(GreaterEqual, s.start, s.current, s.line, s.source[s.start:s.current])
		}
		return NewToken(Greater, s.start, s.current, s.line, s.source[s.start:s.current])
	case '"':
		return s.string()
	}
	return ErrorToken("Unexpected character.", s.line)
}
