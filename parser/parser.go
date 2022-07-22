package parser

import (
	"fmt"
	"strconv"

	"github.com/smiksha1701/buggy/ast"
	"github.com/smiksha1701/buggy/lexer"
	"github.com/smiksha1701/buggy/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:      EQUALS,
	token.NEQ:     EQUALS,
	token.LT:      LESSGREATER,
	token.GT:      LESSGREATER,
	token.PLUS:    SUM,
	token.MINUS:   SUM,
	token.LPAREN:  CALL,
	token.SLASH:   PRODUCT,
	token.ASTERIX: PRODUCT,
}

func (p *Parser) PeekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) CurPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

type Parser struct {
	l                *lexer.Lexer
	curToken         token.Token
	peekToken        token.Token
	errors           []string
	prefixParsingFns map[token.TokenType]prefixParsingFn
	infixParsingFns  map[token.TokenType]infixParsingFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()
	p.prefixParsingFns = make(map[token.TokenType]prefixParsingFn)
	p.RegisterPrefix(token.STRING, p.parseString)
	p.RegisterPrefix(token.IDENT, p.parseIdentifier)
	p.RegisterPrefix(token.INT, p.parseIntegerLiteral)
	p.RegisterPrefix(token.BANG, p.parsePrefixExpression)
	p.RegisterPrefix(token.MINUS, p.parsePrefixExpression)
	p.RegisterPrefix(token.TRUE, p.parseBoolean)
	p.RegisterPrefix(token.FALSE, p.parseBoolean)
	p.RegisterPrefix(token.LPAREN, p.parseGroupExpression)
	p.RegisterPrefix(token.IF, p.parseIfExpression)
	p.RegisterPrefix(token.FUNCTION, p.parseFunctionExpression)
	p.infixParsingFns = make(map[token.TokenType]infixParsingFn)
	p.RegisterInfix(token.PLUS, p.parseInfixExpression)
	p.RegisterInfix(token.MINUS, p.parseInfixExpression)
	p.RegisterInfix(token.LPAREN, p.parseCallFunction)
	p.RegisterInfix(token.ASTERIX, p.parseInfixExpression)
	p.RegisterInfix(token.SLASH, p.parseInfixExpression)
	p.RegisterInfix(token.EQ, p.parseInfixExpression)
	p.RegisterInfix(token.NEQ, p.parseInfixExpression)
	p.RegisterInfix(token.LT, p.parseInfixExpression)
	p.RegisterInfix(token.GT, p.parseInfixExpression)
	return p
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}
	if !p.ExpectedPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.ExpectedPeek(token.RPAREN) {
		return nil
	}

	if !p.ExpectedPeek(token.LBRACE) {
		return nil
	}
	exp.Consequence = p.parseBlockStatement()
	if p.PeekTypeIs(token.ELSE) {
		p.nextToken()

		if !p.ExpectedPeek(token.LBRACE) {
			return nil
		}
		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}
	if !p.ExpectedPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.ExpectedPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()

	return lit
}
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.PeekTypeIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}
	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)
	for p.PeekTypeIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !p.ExpectedPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}
func (p *Parser) parseCallFunction(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}

	exp.Arguments = p.parseFunctionCallParameters()

	return exp
}
func (p *Parser) parseFunctionCallParameters() []ast.Expression {
	args := []ast.Expression{}

	if p.PeekTypeIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.PeekTypeIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.ExpectedPeek(token.RPAREN) {
		return nil
	}
	return args
}
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.CurTokenIs(token.RBRACE) && !p.CurTokenIs(token.EOF) {
		stmt := p.ParseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.ExpectedPeek(token.RPAREN) {
		return nil
	}
	return exp
}
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseString() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) noPrefixFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found ", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ErrorExpectedPeek(t token.TokenType) {
	msg := fmt.Sprintf("expected peek type was = %s got = %s instead ", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) ParseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.PeekTypeIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParsingFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.PeekTypeIs(token.SEMICOLON) && precedence < p.PeekPrecedence() {
		infix := p.infixParsingFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.Return = p.parseExpression(LOWEST)

	if p.PeekTypeIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.ExpectedPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Value: p.curToken.Literal, Token: p.curToken}

	if !p.ExpectedPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	if p.PeekTypeIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) CurTokenIs(tt token.TokenType) bool {
	return p.curToken.Type == tt
}

func (p *Parser) PeekTypeIs(tt token.TokenType) bool {
	return p.peekToken.Type == tt
}

func (p *Parser) ExpectedPeek(tt token.TokenType) bool {
	if p.PeekTypeIs(tt) {
		p.nextToken()
		return true
	} else {
		p.ErrorExpectedPeek(tt)
		return false
	}

}

func (p *Parser) RegisterPrefix(toktype token.TokenType, fn prefixParsingFn) {
	p.prefixParsingFns[toktype] = fn
}

func (p *Parser) RegisterInfix(toktype token.TokenType, fn infixParsingFn) {
	p.infixParsingFns[toktype] = fn
}

type (
	prefixParsingFn func() ast.Expression
	infixParsingFn  func(ast.Expression) ast.Expression
)

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.CurPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.CurTokenIs(token.TRUE)}
}
