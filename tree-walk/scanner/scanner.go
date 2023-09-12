package scanner

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
	return make([]Token, 0)
}
