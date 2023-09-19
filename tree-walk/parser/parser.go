package parser

import (
	"golox/tree-walk/expr"
	"golox/tree-walk/object"
	"golox/tree-walk/rt"
	. "golox/tree-walk/stmt"
	"golox/tree-walk/token"
)

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

func (p *Parser) Parse() []Stmt {
	statements := make([]Stmt, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			p.synchronize()
		}
	}()

	if p.match(token.Class) {
		return p.classDeclaration()
	}
	if p.match(token.Fun) {
		return p.function("function")
	}
	if p.match(token.Var) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) forStatement() Stmt {
	p.consume(token.LeftParen, "Expect '(' after 'for'.")

	var initializer Stmt
	if p.match(token.Semicolon) {
		initializer = nil
	} else if p.match(token.Var) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition expr.Expr = nil
	if !p.check(token.Semicolon) {
		condition = p.expression()
	}
	p.consume(token.Semicolon, "Expect ';' after loop condition.")

	var increment expr.Expr = nil
	if !p.check(token.RightParen) {
		increment = p.expression()
	}
	p.consume(token.RightParen, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = &Block{Statements: []Stmt{body, &Expression{Expression: increment}}}
	}

	if condition == nil {
		condition = &expr.Literal{Value: object.Boolean(true)}
	}
	body = &While{Condition: condition, Body: body}

	if initializer != nil {
		body = &Block{Statements: []Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(token.LeftParen, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(token.RightParen, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt = nil
	if !p.match(token.Else) {
		elseBranch = p.statement()
	}

	return &If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()

	p.consume(token.Semicolon, "Expect ';' after value.")
	return &Print{Expression: value}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()

	var value expr.Expr = nil
	if !p.check(token.Semicolon) {
		value = p.expression()
	}

	p.consume(token.Semicolon, "Expect ';' after return value.")
	return &Return{Keyword: keyword, Value: value}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(token.LeftParen, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(token.RightParen, "Expect ')' after condition.")
	body := p.statement()

	return &While{condition, body}
}

func (p *Parser) expressionStatement() Stmt {
	expression := p.expression()

	p.consume(token.Semicolon, "Expect ';' after expression.")
	return &Expression{Expression: expression}
}

func (p *Parser) statement() Stmt {
	if p.match(token.For) {
		return p.forStatement()
	}
	if p.match(token.If) {
		return p.ifStatement()
	}
	if p.match(token.Print) {
		return p.printStatement()
	}
	if p.match(token.Return) {
		return p.returnStatement()
	}
	if p.match(token.While) {
		return p.whileStatement()
	}
	if p.match(token.LeftBrace) {
		return &Block{Statements: p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) classDeclaration() Stmt {
	name := p.consume(token.Identifier, "Expect class name.")
	var superclass *expr.Variable = nil

	if p.match(token.Less) {
		p.consume(token.Identifier, "Expect superclass name.")
		superclass = &expr.Variable{Name: p.previous()}
	}

	p.consume(token.LeftBrace, "Expect '{' before class body.")

	methods := make([]*Function, 0)
	for !p.check(token.RightBrace) && !p.isAtEnd() {
		methods = append(methods, p.function("method"))
	}

	p.consume(token.RightBrace, "Expect '}' after class body.")

	return &Class{Name: name, Superclass: superclass, Methods: methods}
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(token.Identifier, "Expect variable name.")

	var initializer expr.Expr = nil
	if p.match(token.Equal) {
		initializer = p.expression()
	}
	p.consume(token.Semicolon, "Expect ';' after variable declaration.")
	return &Var{Name: name, Initializer: initializer}
}

func (p *Parser) assignment() expr.Expr {
	expression := p.or()

	if p.match(token.Equal) {
		equals := p.previous()
		value := p.assignment()

		if variable, ok := expression.(*expr.Variable); ok {
			name := variable.Name
			return &expr.Assign{Name: name, Value: value}
		} else if get, ok := expression.(*expr.Get); ok {
			return &expr.Set{Object: get.Object, Name: get.Name, Value: value}
		}
		error(equals, "Invalid assignment target.")
	}

	return expression
}

func (p *Parser) or() expr.Expr {
	expression := p.and()

	for p.match(token.Or) {
		operator := p.previous()
		right := p.and()
		expression = &expr.Logical{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (p *Parser) and() expr.Expr {
	expression := p.equality()

	for p.match(token.And) {
		operator := p.previous()
		right := p.equality()
		expression = &expr.Logical{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (p *Parser) equality() expr.Expr {
	expression := p.comparison()

	for p.match(token.BangEqual, token.EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expression = &expr.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (p *Parser) comparison() expr.Expr {
	expression := p.term()

	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		operator := p.previous()
		right := p.term()
		expression = &expr.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (p *Parser) term() expr.Expr {
	expression := p.factor()

	for p.match(token.Minus, token.Plus) {
		operator := p.previous()
		right := p.factor()
		expression = &expr.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (p *Parser) factor() expr.Expr {
	expression := p.unary()

	for p.match(token.Slash, token.Star) {
		operator := p.previous()
		right := p.unary()
		expression = &expr.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (p *Parser) unary() expr.Expr {
	if p.match(token.Bang, token.Minus) {
		operator := p.previous()
		right := p.unary()
		return &expr.Unary{Operator: operator, Right: right}
	}

	return p.call()
}

func (p *Parser) call() expr.Expr {
	expression := p.primary()

	for {
		if p.match(token.LeftParen) {
			expression = p.finishCall(expression)
		} else if p.match(token.Dot) {
			name := p.consume(token.Identifier, "Expect property name after '.'.")
			expression = &expr.Get{Object: expression, Name: name}
		} else {
			break
		}
	}

	return expression
}

func (p *Parser) finishCall(callee expr.Expr) expr.Expr {
	arguments := make([]expr.Expr, 0)

	if !p.check(token.RightParen) {
		for {
			if len(arguments) >= 255 {
				error(p.peek(), "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, p.expression())
			if !p.match(token.Comma) {
				break
			}
		}
	}

	paren := p.consume(token.RightParen, "Expect ')' after arguments.")

	return &expr.Call{Callee: callee, Paren: paren, Arguments: arguments}
}

func (p *Parser) primary() expr.Expr {
	if p.match(token.False) {
		return &expr.Literal{Value: object.Boolean(false)}
	}
	if p.match(token.True) {
		return &expr.Literal{Value: object.Boolean(true)}
	}
	if p.match(token.Nil) {
		return &expr.Literal{Value: nil}
	}

	if p.match(token.Number, token.String) {
		return &expr.Literal{Value: p.previous().Literal}
	}

	if p.match(token.Super) {
		keyword := p.previous()
		p.consume(token.Dot, "Expect '.' after 'super'.")
		method := p.consume(token.Identifier, "Expect superclass method name.")
		return &expr.Super{Keyword: keyword, Method: method}
	}

	if p.match(token.This) {
		return &expr.This{Keyword: p.previous()}
	}

	if p.match(token.Identifier) {
		return &expr.Variable{Name: p.previous()}
	}

	if p.match(token.LeftParen) {
		expression := p.expression()
		p.consume(token.RightParen, "Expect ')' after expression.")
		return &expr.Grouping{Expression: expression}
	}

	panic(error(p.peek(), "Expect expression."))
}

func (p *Parser) expression() expr.Expr {
	return p.assignment()
}

func (p *Parser) function(kind string) *Function {
	name := p.consume(token.Identifier, "Expect "+kind+" name.")
	p.consume(token.LeftParen, "Expect '(' after "+kind+" name.")

	parameters := make([]token.Token, 0)

	if !p.check(token.RightParen) {
		for {
			if len(parameters) >= 255 {
				error(p.peek(), "Can't have more than 255 parameters.")
			}

			parameters = append(parameters, p.consume(token.Identifier, "Expect parameter name."))
			if !p.match(token.Comma) {
				break
			}
		}
	}

	p.consume(token.RightParen, "Expect ')' after parameters.")

	p.consume(token.LeftBrace, "Expect '{' before "+kind+" body.")
	body := p.block()
	return &Function{Name: name, Params: parameters, Body: body}
}

func (p *Parser) block() []Stmt {
	statements := make([]Stmt, 0)

	for !p.check(token.RightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(token.RightBrace, "Expect '}' after block.")
	return statements
}

func error(token token.Token, message string) ParseError {
	rt.ErrorToken(token, message)
	return ParseError{message: message}
}

func (p *Parser) consume(typ token.TokenType, message string) token.Token {
	if p.check(typ) {
		return p.advance()
	}
	panic(ParseError{message: message})
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(typ token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == typ
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.Semicolon {
			return
		}

		switch p.peek().Type {
		case token.Class:
		case token.Fun:
		case token.Var:
		case token.For:
		case token.If:
		case token.While:
		case token.Print:
		case token.Return:
			return
		}

		p.advance()
	}

}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.Eof
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}
