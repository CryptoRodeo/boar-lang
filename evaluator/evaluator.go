package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	//statements
	case *ast.Program:
		return evalStatements(node.Statements)

		// a single statement
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	//expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}

	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

/**
dev notes:

eval()
- should start at the top of the AST, recieving an *ast.Program and
  then traverse every node in it and evaluate every statement (aka tree walking).

self-evaluating expressions:
- what we call literals
- we input an integer into eval() and get that integer back (hence they evaluate themselves.)
- we input an *ast.IntegerLiteral, eval() returns an *object.Literal with a Value of that integer
**/
