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
	EQUALS        // ==
	LESSGREATER   // < or >
	SUM           // +
	PRODUCT       // *
	PREFIX        // -X or !X
	ASSIGN        // =
	CALL          // myFunction(x)
	INDEX         //array[index]
	INTERNAL_CALL // arr.pop, hash.delete, etc
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
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
	token.DOT:      INTERNAL_CALL,
	token.ASSIGN:   ASSIGN,
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
	// function expressions
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

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
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parseInternalCallExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignmentExpression)

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
		than the one currently passed or we encounter a semicolon
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
			this token will be used in the parseInfixExpression function
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

func (p *Parser) parseIfExpression() ast.Expression {
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

	// The tokens get advanced enough so we are now sitting on the LBRACE

	expression.Consequence = p.parseBlockStatement()

	// At this point we should be sitting on the right brace }
	// Check if there is an 'else', move up tokens if there is
	if p.peekTokenIs(token.ELSE) {
		// we're currently sisting on the 'else' token, move up the tokens
		p.nextToken()
		// If for some reason theres not a LBRACE token immediately after the else the expression is invalid
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	// Jump over the LBRACE token
	p.nextToken()
	// Continue parsing the statement until we reach the end of the block or token.EOF
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	func_lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	func_lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	func_lit.Body = p.parseBlockStatement()

	return func_lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	//empty parameter list
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	// Move past the parenthesis we're currently on, point to the first identifier
	p.nextToken()

	// Grab first identifier
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	// Every time we encounter a comma, move token pointers up
	// so we point to the identifiers
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	// no closing parenthesis
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// recieves the already parsed function as argument, uses it to construct call expression node.
// "leftExp" in parseExpressions gets passed to this infix parsing function
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	/*
		p.curToken => left parenthesis
		function => identifier (i.e.: add, subtract, doTheThing)
	*/
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	// (x,y,z)
	exp.Arguments = p.parseExpressionList(token.RPAREN)

	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	// No arguments
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	// most past lparen
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	// At the end of the argument list we should see a right parenthesis / closing parenthesis
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parses a list of expressions until we reach the end of the list (via the end token type)
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	// empty list
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	// first expression
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	// For each comma separated expression
	for p.peekTokenIs(token.COMMA) {
		// move up 2 tokens (so we land on the expressoin)
		p.nextToken()
		p.nextToken()
		// parse the expression
		list = append(list, p.parseExpression(LOWEST))
	}

	// If for some reason we haven't reached the end token we passed, return nil
	if !p.expectPeek(end) {
		return nil
	}

	return list

}

func (p *Parser) parseArrayLiteral() ast.Expression {
	// the [ token
	array := &ast.ArrayLiteral{Token: p.curToken}
	// Grab all the elements before we reach the right bracket (end of array)
	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	/**
		In this function the [  in someArray[0] is treated as the infix operator.
		someArray being the left operand and 0 being the right operand
	**/
	// curToken => [
	// left => some identifier, array literal, etc
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	// move pointer to where the index expression is
	p.nextToken()

	// parse the index expression (1, 1+1, a*b, etc)
	index := p.parseExpression(LOWEST)
	exp.Index = index

	// we should reach a ] after the index expression
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	// Check if we're trying to assign a value at some index / key.
	// If we are, the next token should be '='
	// hash[a] = 2, arr[0] = 1
	if p.peekTokenIs(token.ASSIGN) {
		return p.parseIndexAssignment(exp, index)
	}

	return exp
}

func (p *Parser) parseIndexAssignment(node, index ast.Expression) ast.Expression {

	indexExp, ok := node.(*ast.IndexExpression)

	if !ok {
		return nil
	}

	// Grab the identifier from the original IndexExpression
	// a[index], hash[key] => a, hash
	identifier := indexExp.Left

	// we should currently be at the ']' token, traverse to the assign token
	p.nextToken()
	// We should now be at the assign token '='
	token := p.curToken
	// move onto what should be a value
	p.nextToken()
	// we should now have some value to parse
	value := p.parseExpression(LOWEST)

	return &ast.IndexAssignment{Left: identifier, Index: index, Token: token, Value: value}
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken} // the { symbol
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	// Until we reach an } (the end of the hash)
	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()

		key := p.parseExpression(LOWEST)

		// A key should be followed by a : (ex: "key":"value")
		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		// grab the value
		value := p.parseExpression(LOWEST)
		// create the pair
		hash.Pairs[key] = value

		// If we haven't reached a right brace (end of hash) or comma (next pair)
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	// If we haven't reach the end of the hash by now, somethings wrong.
	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseInternalCallExpression(left ast.Expression) ast.Expression {
	// the current token should be '.'
	if !p.curTokenIs(token.DOT) {
		return nil
	}

	dot := p.curToken

	// Grab the identifier: arr, hash, etc.
	ident, ok := left.(*ast.Identifier)

	if !ok {
		return nil
	}

	// We should now be at the function name: pop, delete, etc
	p.nextToken()
	internal_function_ident := p.parseIdentifier()
	func_ident, ok := internal_function_ident.(*ast.Identifier)

	if !ok {
		return nil
	}

	//After the function name we should expect a '('
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// After the '(' we should have either 0 -> expressions
	args := p.parseExpressionList(token.RPAREN)

	if !ok {
		return nil
	}

	ifc := &ast.InternalFunctionCall{
		CallerIdentifier:   ident,
		Token:              dot,
		FunctionIdentifier: func_ident,
		Arguments:          args,
	}
	return ifc
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	// the current token should be '.'
	if !p.curTokenIs(token.ASSIGN) {
		return nil
	}

	assignment := &ast.AssignmentExpression{}

	assign := p.curToken

	assignment.Token = assign

	// Grab the identifier: arr, hash, etc.
	ident, ok := left.(*ast.Identifier)

	if !ok {
		return nil
	}

	assignment.Name = ident

	// Now lets grab the expression after '='
	p.nextToken()
	assignment.Value = p.parseExpression(LOWEST)

	// If we reach a semicolon (terminating the expression), move onto the next token.
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return assignment

}

/**
Dev Notes:

Concepts:
TLDR:
- Parser generates AST Nodes from the tokens generated by the lexer

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
