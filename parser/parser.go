package parser

import (
	"fmt"

	"github.com/smiksha1701/buggy/ast"
	"github.com/smiksha1701/buggy/lexer"
	"github.com/smiksha1701/buggy/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()

	return p
}
func (p *Parser) Errors() []string {
	return p.errors
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
	default:
		return nil
	}

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
	if !p.CurTokenIs(token.SEMICOLON) {
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
