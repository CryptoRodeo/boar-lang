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

// Prefix and infix parsing functions
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
	// array of statements
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		// parse the current statement
		stmt := p.parseStatement()

		if stmt != nil {
			// add the current statement to the program statements array
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
	// Finally, lets skip the expression and stop when encountering a semicolon.
	// TODO: we're skipping the expressions until we encounter a semicolon.
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	// let <identifier> <assign> <expression> ;
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	// move up to the next token
	p.nextToken()

	//TODO: we're skipping expressions until we encounter a semicolon.
	for !p.curTokenIs(token.SEMICOLON) {
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

	// If it is, return it.
	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	/*
		Here we pass the lowest possible precedence to parseExpression, since
		we didn't parse anything yet and we can't compare precedences.
	*/
	stmt.Expression = p.parseExpression(LOWEST)

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
	// CHeck the type of the next oken
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
