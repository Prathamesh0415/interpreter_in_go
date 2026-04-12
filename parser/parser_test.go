package parser

import (
	"interpreter_go/lexer"
	"interpreter_go/ast"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := ` let x = 5;
			let y = 10;
			let foobar = 838383;`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

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