package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

/**
- use iota to give the following constants incrementing numbers as values
- the _ identifier takes the zero and the following constants get assigned
values 1 to x

note:
- the order of the relations between these constants matter.
- it will allow us to answer questions regarding precedence
ex: "does the * operator have a higher precedence than the == operator?"

TLDR:
- these values will be used to identify the precedence of a token.
- The values are incrementing, from 0 to X (so the order matters)
- + has a lower precedence than *, etc.
**/
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(x)
)

/**
Precedence table.
- associates token types with their precedence.
ex:
- token.PLUS and token.MINUS have the same precedence
- these tokens have a lower precedence than token.ASTERISK and token.SLASH
**/
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

/**
Prefix and infix parsing functions

examples:
prefix expression: --5,

infix: a + b
**/
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	// pointer to an instance of the lexer
	// used for calling NextToken() to get the next token in the input.
	l *lexer.Lexer
	// token values
	curToken  token.Token
	peekToken token.Token
	// slice of error strings
	errors []string

	//parsing functions
	/**
		Note:
		- Since we're using the Pratt Parser implementation it makes sense to use a map here.
		- The token types are associated with a parsing function.
		- Each token type can have up to two parsing functions associated with it, depending on its position (prefix / infix)
		// key: tokenType, res: prefix/infix function
	**/
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	// generate a pointer to this new Parser struct
	p := &Parser{l: l, errors: []string{}}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	//Initialize the prefixParseFn map, register a parsing function for Identifiers.
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	//  if we encounter a token of type token.IDENT the parsing function to call is parseIdentifier
	// ex: x, foobar => call parseIdentifier
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	// If we encounter a token of type token.INT, call parseIntegerLiteral
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	// If we encounter a token of type BANG (!), call this function
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	// parse grouped expressions
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	// if expressions
	p.registerPrefix(token.IF, p.parseIfExpression)

	// Initialize the infix parse function map
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// Helper method to advance token pointers
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	// parser.lexer.nextToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	// pointer to the program
	program := &ast.Program{}
	// slice of statements
	program.Statements = []ast.Statement{}
	// Loop until we reach a null token / no token
	for !p.curTokenIs(token.EOF) {
		// parse the current statement
		stmt := p.parseStatement()

		if stmt != nil {
			// add the current statement to the program statements slice
			program.Statements = append(program.Statements, stmt)
		}
		// move onto the next token
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		// by default we'll parse it as an expression: x, foobar, x + y, etc
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// grabs the 'let' statement
	stmt := &ast.LetStatement{Token: p.curToken}
	// We expect to find an identifier: let x, let a, let etc
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	// Construct an identifier node
	// now we have let <identifier>
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// We then expect to find an equal sign after the identifier
	// ex: let <identifier> <assign>
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	// let <identifier> <assign> <expression> ;
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	// move up to the next token
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Create an error when no prefix parse function has been found
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// See if the current token is registered to a parsing function
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		// If its not, create an error
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	/**
	prefix parsing function exists, call it, grab the value.
	**/
	leftExp := prefix()
	/*
		- find the infix parsing function for the next token (if it exists)
		- If it exists, call it, building up the Infix Expression Node
		- Continue doing this until we encounter a token that has a higher precedence
		than th eone passed or a semicolon
	*/
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// grab the infix parsing function for this specific token (if it exists)
		// ex: curToken => 5, peektoken => +
		infix := p.infixParseFns[p.peekToken.Type]
		// if no function exist (because its not an infix operator), return the leftExp
		if infix == nil {
			return leftExp
		}
		/*
			else, move to the next token (the infix operator).
			this token wil lbe ued in the parseInfixExpression function
		*/
		p.nextToken()
		/*
			Pass the already parsed AST Expression Node to the parseInfixExpression function
			so it can be assigned as its 'left' value.

			As this loop progresses leftExp will change and continue to get passed into the
			parseInfixExpression function, building an expression of multiple, embedded infix expressions:
			ex: 1 + 2 + 3 => ( (1+2) + 3 )
			-> (1+2) being the first infix expression parsed
			-> ( (1+2) + 3) being the final returned expression, where (1+2) is an infix expresssion
			assigned to the 'left' value of the outer infix expression ((inner) + 3)
		*/
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// Create an ExpressionStatement AST Node
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	/*
		Here we pass the lowest possible precedence to parseExpression, since
		we didn't parse anything yet and we can't compare precedences.
	*/
	stmt.Expression = p.parseExpression(LOWEST)
	// If the next token is a semicolon, move onto the next token
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	// Check the type of the next token
	if p.peekTokenIs(t) {
		// If its correct, advance the tokens
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

//Returns any parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// Adds any errors we encountered while peeking in expectPeek()
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

/**
- Returns the precedence associated with the token type of p.peekToken
- Defaults to LOWEST
**/
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

/**
- Returns the precedence associated with the token type of p.curToken
- Defaults to LOWEST
**/
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	// convert string into an int64
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// Parses expressions with prefixes: -5, !true, etc
// anytime this function is called the tokens advance and the current token
// is the one after the prefix operator
func (p *Parser) parsePrefixExpression() ast.Expression {
	// Create the prefix expression
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	// navigate to the next token
	/**
		Note:
		- Unlike most parsing functions this one advances our tokens.
		- The reason for this is because in order to correctly parse
		  a prefix expression like -3 more than one token has to be "consumed".
		- So after we grabe our current token (being the prefix) we advance the tokens
		  and call parseExpression

	**/
	p.nextToken()
	/*
		Now that we have advanced the tokens, the next token will be the one after the prefix operator.
		Take that token and set the current token as the prefix expressions "Right" value (value after prefix)

		ex:
		if we encounter '-5' when parseExpression is called then p.curToken.Type is token.INT.
		parseExpression then checks the registered prefix parsing functions and finds its associated parsing function (parseIntegerLiteral).
		This function builds the an *ast.IntergerLiteral node and returns it.

		parseExpression returns this new node and uses it to fill the Right field of *ast.PrefixExpression
	*/
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

/**
- Takes an ast.Expression argument as the 'left' side of the infix expression
- Grabs the precedence of the current token (operator of the infix expression)
- Advances the tokens, filling the Right field of the node
**/
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// Generate the infix expression struct
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	/*
		Grab the precedence of the current token (the operator)
		before advancing the token pointers.
	*/
	precedence := p.curPrecedence()
	// Point to the next token
	p.nextToken()
	// Parse and grab the next AST Node
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfexpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	// we should expect a left parenthesis as the next token
	// i.e. if ( x )
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	// progress tokens, parse expression
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	/**
		Make sure we encounter a right parenthesis, progress tokens if we do
		Then make sure we encounter a left brace {, progress tokens if we do

		Note:
		- A side effect of the expectPeek function is it progresses the tokens
		if the condition is true.
	**/
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	return expression
}

/**
Dev Notes:

Concepts:

Parser:
- A parser is a software component that takes input data (frequently text) and builds
a data structure – often some kind of parse tree, abstract syntax tree or other
hierarchical structure

– This tree creates a structural representation of the input, checking for
correct syntax in the process.

- The parser is often preceded by a separate lexical analyser, which creates tokens
from the sequence of input characters;

TLDR: This is what parsers do: take source input (in our case tokens) and produce a data structure that represents source code

Pratt Parsing:
- The main idea is to associate parsing functions with token types
- whenever a token type is encountered, the parsing functions are called to
parse the appropriate expression and retur nan AST node that represents it.
- each token type can have up to two parsing functions associated with it (prefx + infix),
depending on its position

Prefix operator:
- an operator "in front" of its operand
ex: --5

here the operator is -- (decrement), the operand is the integer literal 5 and the operator is
in the prefix position

Postfix operator:
- an operator "after" its operand.
ex: 5++
note: we won't have postfix operators in monke-lang (for now)

Infix Operators:
- when the operator sits between its operands
ex: 5 + 5


Other:
------------------------------
How this parser works:
It repeatedly advances the tokens and checks the current token to decide what to do next:
- either call another parsing function (prefix or infix)
- throw an error.

Each function then does its job and possibly constructs an AST node so that
the “main loop” in parseProgram() can advance the tokens and decide what to do again.

Parser approach:
- This parser uses Top Down Operator Precedence Parsing (aka Vaughan Pratt Parsing)
- This is different from Backus-Naur-Form parsing (which foxues on grammer rules)

different how?
- instead of associating parsing functions with grammer rules, Pratt parsing associates
parsing functions with single token types (like we do here).
- Each token type can have two parsing functions associates with it depending on the tokens position: infix or prefix
**/
