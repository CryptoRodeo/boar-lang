package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

// Node Interface for our AST
type Node interface {
	TokenLiteral() string //returns the literal value of the token its associated with
	String() string       //for debugging and comparison
}

// Statements, a type of ndoe in our AST
// Note: Statements do not generate values, expressions do.
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

// A statement that consists of only one expression
// ex: let x = 5;
// the expression here being 5 (which generates a value)
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

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

type InfixExpression struct {
	Token    token.Token // the operator token: -, +, etc
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// if (condition) <consequence> else <alternative>
type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

/**
Function literals are Expressions.
They can be used anywhere expressions are valid.

ex:
a function literal as the expression in a let statement:
let someFunction = fn(x, y) { return x + y; }

as the expression in a return statement inside another function literal:
fn() {
	return fn(y,z) { return y > x; }
}

as an argument when calling another function:
outerFunc(a, b, fn(c, d) { return a > b; });

**/
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier   // (x,y,z)
	Body       *BlockStatement // { x + y; }, { foo > bar; }
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	// fn(x,y,z)
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	//fn(params)
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()

}

/**
Call expression
<expression>(<comma separated expressions>)

call expressions consist of an expression that results in:
- a function when evaluated
- a list of expressions that are the arguments to this function call.

example: add(1,2) => fn(a,b) { a + b; }(1,2)
**/

type CallExpression struct {
	Token     token.Token // the '(' token
	Function  Expression  // idenfifier or function literal
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token // the [ token
	Left  Expression  // the object being access (array, some identifier, function call, etc)
	Index Expression  // any expression as long as it produces an integer
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

type IndexAssignment struct {
	Token token.Token // the = token
	Left  Expression
	Index Expression
	Value Expression
}

func (ia *IndexAssignment) expressionNode()      {}
func (ia *IndexAssignment) TokenLiteral() string { return ia.Token.Literal }
func (ia *IndexAssignment) String() string {
	var out bytes.Buffer
	out.WriteString(ia.Left.String())
	out.WriteString("[")
	out.WriteString(ia.Index.String())
	out.WriteString("]")
	out.WriteString("=")
	out.WriteString(ia.Value.String())
	return out.String()
}

/**
The basic syntactic structure of a hash literal is:
{<expression> : <expression>, ... }
**/
type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}

	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	// { <expression> : <expression>, ...}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type InternalFunctionCall struct {
	Token              token.Token  // the '.' token
	CallerIdentifier   *Identifier  //someArray, someHash, etc
	FunctionIdentifier *Identifier  // pop, delete, etc.
	Arguments          []Expression //(1,2,3), (), etc.
}

func (ifc *InternalFunctionCall) expressionNode()      {}
func (ifc *InternalFunctionCall) TokenLiteral() string { return ifc.Token.Literal }
func (ifc *InternalFunctionCall) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, a := range ifc.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ifc.CallerIdentifier.String())   //someArray, someHash, etc
	out.WriteString(ifc.Token.Literal)               // .
	out.WriteString(ifc.FunctionIdentifier.String()) // delete, pop
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", ")) // (), (1,2,3), ("a", "b", "c"), etc
	out.WriteString(")")

	return out.String()
}

type AssignmentExpression struct {
	Token token.Token // the = token
	Name  *Identifier //identifier for the binding (ex: x in x = 5)
	Value Expression  //expression that produces the value (the 5 in let x = 5)
}

func (as *AssignmentExpression) expressionNode()      {}
func (as *AssignmentExpression) TokenLiteral() string { return as.Token.Literal }
func (as *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(as.Name.String())
	out.WriteString(as.TokenLiteral())

	if as.Value != nil {
		out.WriteString(as.Value.String())
	}

	out.WriteString(";")

	return out.String()

}

// for (<counter variable init>;<loop conditional>;<counterVar increment>) { <statements> };
type ForLoopStatement struct {
	Token         token.Token   // the 'for' token
	CounterVar    *LetStatement //identifier for the binding (ex: x in x = 5)
	LoopCondition Expression
	CounterUpdate *AssignmentExpression //expression that produces the value (the 5 in let x = 5)
	LoopBlock     *BlockStatement
}

func (fl *ForLoopStatement) statementNode()       {}
func (fl *ForLoopStatement) TokenLiteral() string { return fl.Token.Literal }
func (fl *ForLoopStatement) String() string {
	var out bytes.Buffer
	out.WriteString(fl.Token.Literal)
	out.WriteString("(")
	out.WriteString(fl.CounterVar.String())
	out.WriteString(";")
	out.WriteString(fl.LoopCondition.String())
	out.WriteString(";")
	out.WriteString(fl.CounterUpdate.String())
	out.WriteString("{")
	out.WriteString(fl.LoopBlock.String())
	out.WriteString("};")

	return out.String()

}
