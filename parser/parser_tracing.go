package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 0

const traceIdentPlaceholder string = "\t"

func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}
func incIdent() { traceLevel = traceLevel + 1 }
func decIdent() { traceLevel = traceLevel - 1 }

func trace(msg string) string {
	incIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	tracePrint("END " + msg)
	decIdent()
}

/**
This file is used to trace the parser as it goes along creating AT nodes.

It can be used like this:

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
defer untrace(trace("parseExpressionStatement"))
...
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
defer untrace(trace("parseExpression"))
...
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
defer untrace(trace("parseIntegerLiteral"))
...
}

func (p *Parser) parsePrefixExpression() ast.Expression {
defer untrace(trace("parsePrefixExpression"))
...
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
defer untrace(trace("parseInfixExpression"))
...
}

and when running the parser tests it will generate output like this:
(example test case: -1 * 2 + 3)

BEGIN parseExpressionStatement
	BEGIN parseExpression
		BEGIN parsePrefixExpression
			BEGIN parseExpression
				BEGIN parseIntegerLiteral
				END parseIntegerLiteral
			END parseExpression
		END parsePrefixExpression
		BEGIN parseInfixExpression
			BEGIN parseExpression
				BEGIN parseIntegerLiteral
				END parseIntegerLiteral
			END parseExpression
		END parseInfixExpression
		BEGIN parseInfixExpression
			BEGIN parseExpression
				BEGIN parseIntegerLiteral
				END parseIntegerLiteral
			END parseExpression
		END parseInfixExpression
	END parseExpression
END parseExpressionStatement
**/
