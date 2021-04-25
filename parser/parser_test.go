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
	fmt.Printf("==> parser has %d errors\n", len(errors))
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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program.Statements has %d statments", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("identifier expression type error, got type = %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("identifier expression type error, got type = %T", program.Statements[0])
	}
	if ident.Value != "foobar" {
		t.Errorf("parser value error")
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("parser token literal error")
	}

}

func TestIntegerLiteralExpression(t *testing.T) {

	fmt.Println("TestIntegerLiteralExpression")

	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program.Statements has %d statments", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("identifier expression type error, got type = %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("identifier expression type error, got type = %T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("parser value error")
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("parser token literal error")
	}

}

func TestParsingPrefixExpresison(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("parser error statement counts %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expression type error, got type = %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expression type error, got type = %T", program.Statements[0])
		}
		if exp.Operator != tt.operator {
			t.Fatalf("expression operator error, got type = %T", program.Statements[0])
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("integer literal type error, got type = %T", il)
	}

	if integ.Value != value {
		t.Errorf("integer literal value error, got type = %d", integ.Value)
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer literal token literal error, got type = %s", integ.TokenLiteral())
	}
	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("parser error statement counts %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expression type error, got type = %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expression type error, got type = %T", program.Statements[0])
		}
		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("expression operator error, got type = %T", program.Statements[0])
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}

}
