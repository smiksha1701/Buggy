package parser

import (
	"fmt"
	"testing"

	"github.com/smiksha1701/buggy/ast"
	"github.com/smiksha1701/buggy/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParserErrors(p, t)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}
		stmt := program.Statements[0]
		if !CheckLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}
func TestReturnStatements(t *testing.T) {
	input := `
	return b;
	return a+b;
	return a(1,2);
	`

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
	for _, stmt := range program.Statements {
		returnstmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Expected ast.ReturnStatement got=%f", stmt)
			continue
		}
		if returnstmt.TokenLiteral() != "return" {
			t.Errorf("Expected return got=%s", returnstmt.TokenLiteral())
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
func TestIdentifiersExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(p, t)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ast.ExpressionStatement got=%f", program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Fatalf("ident.Value expected=foobar got=%T", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident.TokenLiteral() expected=foobar got=%T", ident.TokenLiteral())
	}

}
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Fatalf("ident.Value expected=%T got=%T", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Fatalf("ident.TokenLiteral() expected=%T got=%T", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got = %T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value is wrong expected=%d, got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral() not %d, got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("il not *ast.. got = %T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("integ.Value is wrong expected=%t, got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("integ.TokenLiteral() not %t, got=%s", value, bo.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}
func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(p, t)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ast.ExpressionStatement got=%f", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Fatalf("ident.Value expected=5 got=%T", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Fatalf("ident.TokenLiteral() expected=5 got=%T", literal.TokenLiteral())
	}

}
func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(p, t)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ast.ExpressionStatement got=%f", program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Fatalf("ident.Value expected=foobar got=%s", literal.Value)
	}
	if literal.TokenLiteral() != "hello world" {
		t.Fatalf("ident.TokenLiteral() expected=foobar got=%s", literal.TokenLiteral())
	}

}
func TestParsingPrefixExpressions(t *testing.T) {
	prefix_tests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"-15;", "-", 15},
		{"!5;", "!", 5},
	}
	for _, tt := range prefix_tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParserErrors(p, t)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected ast.ExpressionStatement got=%f", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not *ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator expected=%s got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}
func TestParsingInfixExpressions(t *testing.T) {
	infix_tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5-5", 5, "-", 5},
		{"5+5", 5, "+", 5},
		{"5*5", 5, "*", 5},
		{"5/5", 5, "/", 5},
		{"5<5", 5, "<", 5},
		{"5>5", 5, ">", 5},
		{"5==5", 5, "==", 5},
		{"5!=5", 5, "!=", 5},
	}
	for _, tt := range infix_tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParserErrors(p, t)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected ast.ExpressionStatement got=%f", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp not *ast.InfixExpression. got=%T", stmt.Expression)
		}
		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator expected=%s got=%s", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5-5",
			"(5 - 5)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a+(b+c)",
			"(a + (b + c))",
		},
		{
			"a+b+c==a*b/3",
			"(((a + b) + c) == ((a * b) / 3))",
		},
		{
			"3<5*3>5",
			"((3 < (5 * 3)) > 5)",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParserErrors(p, t)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
func TestBooleanExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(p, t)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ast.ExpressionStatement got=%f", program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
	}
	if ident.Value != true {
		t.Fatalf("ident.Value expected=foobar got=%T", ident.Value)
	}
	if ident.TokenLiteral() != "true" {
		t.Fatalf("ident.TokenLiteral() expected=foobar got=%T", ident.TokenLiteral())
	}

}

func TestParsingIndexExpression(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(p, t)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Errorf("Expected not *ast.IndexExpression got=%f", stmt.Expression)
	}
	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}

}

func TestIfStatements(t *testing.T) {
	input := `
	if (x<y){x}
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(p, t)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression got=%T", stmt.Expression)
	}
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}
	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}
	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Fatalf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `
	fn(x, y){ x + y; }
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(p, t)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement got=%T", program.Statements[0])
	}
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral got=%T", stmt.Expression)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong amount. Wanted=2 got=%T", len(function.Parameters))
	}
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(function.Body.Statements))
	}
	bodystmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}
	testInfixExpression(t, bodystmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(a) {};", expectedParams: []string{"a"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		CheckParserErrors(p, t)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n", len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestArrayParsing(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(p, t)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ArrayLiteral got=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, x + y)"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	CheckParserErrors(p, t)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression got=%T", stmt.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], "x", "+", "y")

}
