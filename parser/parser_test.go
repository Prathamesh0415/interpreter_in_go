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