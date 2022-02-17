package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
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
}

func New(l *lexer.Lexer) *Parser {
	// generate a pointer to this new Parser struct
	p := &Parser{l: l, errors: []string{}}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
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
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// investigate
	stmt := &ast.LetStatement{Token: p.curToken}
	// We expect to find an identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	// Construct an identifier node
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// We then expect to find an equal sign after the identifier
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// Finally, lets skip the expression and stop when encountering a semicolon.
	// TODO: we're skipping the expressions until we encounter a semicolon.
	for !p.curTokenIs(token.SEMICOLON) {
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
