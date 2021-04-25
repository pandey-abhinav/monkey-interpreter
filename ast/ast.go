package ast

type Node interface {
	TokenLiteral() string
	String() string
}

// Statement interface represents all the statements in the language
type Statement interface {
	Node
	statementNode()
}

// Statement interface represents all expressions in the language
type Expression interface {
	Node
	expressionNode()
}
