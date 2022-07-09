package parser

import (
	"fmt"
	"testing"

	"github.com/smiksha1701/buggy/ast"
	"github.com/smiksha1701/buggy/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x 5;
	let 10;
	let = ;
	`
	fmt.Printf("Hello World")
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(p, t)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !CheckLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}
func CheckParserErrors(p *Parser, t *testing.T) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Parser had %d errors", len(errors))
	for _, e := range errors {
		t.Errorf("parser error: %q", e)
	}
	t.FailNow()
}
func CheckLetStatement(t *testing.T, stmt ast.Statement, tt string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("TokenLiteral expected=let got=%s", stmt.TokenLiteral())
		return false
	}
	s, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("*ast.LetStatement expected got=%T", stmt)
		return false
	}
	if s.Name == nil {
		t.Errorf("here")
		return false
	}
	if s.Name.Value != tt {
		t.Errorf("LetStatement.Name.Value expected=%s got=%s", tt, s.Name.Value)
		return false
	}
	if s.Name.TokenLiteral() != tt {
		t.Errorf("LetStatement.Name expected=%s got=%s", tt, s.Name)
		return false
	}

	return true
}
