package parser

import (
	"interpreter_go/lexer"
	"interpreter_go/ast"
	"testing"
	"fmt"
)

func TestLetStatements(t *testing.T) {
	input := ` let x = 5;
			let y = 10;
			let foobar = 838383;`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
		len(program.Statements))
	}

	tests := []struct{
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatements(t, statement, tt.expectedIdentifier){
			return
		} 
	}

}

func TestReturnStatements(t *testing.T) {
	input := ` return 5;
			return 10;
			return add(10);`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
		len(program.Statements))
	}

	for _, s := range program.Statements {
		stmt, ok := s.(*ast.ReturnStatement) 
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if stmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
			stmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
		ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
		len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
		program.Statements[0])
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
		literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input string
		operator string
		integerValue interface{}
	}{
		{"!true;", "!", true},
		{"!false;", "!", false},
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if _, ok := tt.integerValue.(int64); ok {
			if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
				return
			}
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
infixTests := []struct {
		input string
		leftValue interface{}
		operator string
		rightValue interface{}
	}{
	{"true == true", true, "==", true},
	{"true != false", true, "!=", false},
	{"false == false", false, "==", false},
	{"5 + 5;", 5, "+", 5},
	{"5- 5;", 5, "-", 5},
	{"5 * 5;", 5, "*", 5},
	{"5 / 5;", 5, "/", 5},
	{"5 > 5;", 5, ">", 5},
	{"5 < 5;", 5, "<", 5},
	{"5 == 5;", 5, "==", 5},
	{"5 != 5;", 5, "!=", 5},}
for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}
		if _, ok := tt.leftValue.(int64); ok {
			if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
				return
			}
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}	
		if _, ok := tt.rightValue.(int64); ok { 
			if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
				return
			}
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
tests := []struct {
		input string
		expected string
	}{
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
		"-a * b",
		"((-a) * b)",
	},
	{
		"!-a",
		"(!(-a))",
	},
	{
		"a + b + c",
		"((a + b) + c)",
	},
	{
		"a + b- c",
		"((a + b) - c)",
	},
	{
		"a * b * c",
		"((a * b) * c)",
	},
	{
		"a * b / c",
		"((a * b) / c)",
	},
	{
		"a + b / c",
		"(a + (b / c))",
	},
	{
		"a + b * c + d / e- f",
		"(((a + (b * c)) + (d / e)) - f)",
	},
	{
		"3 + 4;-5 * 5",
		"(3 + 4)((-5) * 5)",
	},
	{
		"5 > 4 == 3 < 4",
		"((5 > 4) == (3 < 4))",
	},
	{
		"5 < 4 != 3 > 4",
		"((5 < 4) != (3 > 4))",
	},
	{
		"3 + 4 * 5 == 3 * 1 + 4 * 5",
		"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
	},
	{
		"3 + 4 * 5 == 3 * 1 + 4 * 5",
		"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
	},
}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
	t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
		case int:
			return testIntegerLiteral(t, exp, int64(v))
		case int64:
			return testIntegerLiteral(t, exp, v)
		case string:
			return testIdentifier(t, exp, v)
		case bool:
			return testBooleanLiteral(t, exp , v)
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

func testIntegerLiteral(t *testing.T, il ast.Expression, value interface{}) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func testLetStatements(t *testing.T, statement ast.Statement, expected string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", statement.TokenLiteral())
		return false
	}
	letStmt, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", statement)
		return false
	}
	if letStmt.Name.Value != expected {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", expected, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != expected {
		t.Errorf("s.Name not '%s'. got=%s", expected, letStmt.Name)
		return false
	}
	return true

}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
		value, bo.TokenLiteral())
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.errors
	if len(errors) == 0 {
		return
	}
	fmt.Printf("Parser has %d errors \n", len(errors))
	for _, error := range errors {
		t.Errorf("parse Error: %q", error)
	}
	t.FailNow()
}