package ast

import (
	"bytes"
	"monkey/token"
)

// Node Interface for our AST
type Node interface {
	TokenLiteral() string //returns the literal value of the token its associated with
	String() string       //for debugging and comparison
}

// Statements, a type of ndoe in our AST
type Statement interface {
	Node
	statementNode()
}

// Expressions, AST Nodes that generate a value:
// ex: let x = 6 doesn't produce a value, but 6 does (the value being 6)
type Expression interface {
	Node
	expressionNode()
}

// Root node of AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	// If we have any statements return the first one
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	// Create a buffer
	var out bytes.Buffer
	// Return the value of each statements String() method
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	// Return the buffer as a string
	return out.String()
}

// Implements Statement and Node interface
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier //identifier for the binding (ex: x in let x = 5)
	Value Expression  //expression that produces the value (the 5 in let x = 5)
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()

}

/**
Implements expression interface
note:
- Identifiers don't always produce a value.
- We'll use this node to keep the number of different node types small.
- We'll use Identifier here to represent the name in a variable binding
  and also to represent an identifier as part of or as a complete expression.
**/

// ex: the x in let x = 5
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type ReturnStatement struct {
	Token       token.Token //the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.TokenLiteral() + " ")
	}

	out.WriteString(";")

	return out.String()
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token //the prefix token: !, -
	Operator string      // !, -
	Right    Expression  //expression after operator -> ex: !true, -10, etc.
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	/*
		- add the parentheses around the operator and its operand (expression on the right).
		- This allows us to see which operands belong to which operator
	*/
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
