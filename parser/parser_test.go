package parser

import (
	"fmt"
	"pandey-abhinav/monkey-interpreter/ast"
	"pandey-abhinav/monkey-interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() return nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program.Statements has %d statments", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		fmt.Printf("i = %v\n", i)
		stmt := program.Statements[i]
		fmt.Printf("statement = %+v\n", stmt)
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	fmt.Printf("==> parser has %d errors", len(errors))
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error : %q", err)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("token literal is not let, got = %q", s.TokenLiteral())
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("statement not an ast, got=%T", s)
	}

	if letStmt.Name.Value != name {
		t.Errorf("let statement value is not equal to name , value = %s, name = %s",
			letStmt.Name.Value, name)
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("token literal is not equal to name, token literal = %s, name = %s",
			letStmt.Name.TokenLiteral(), name)
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("Program.Statements has %d statments", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt type not correct, got = %T\n", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("return liternal not correct, got = %q", returnStmt.TokenLiteral())
		}
	}

}
