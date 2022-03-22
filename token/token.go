package token

// Allows us to distinguish between different types of tokens
type TokenType string

const (
	// Signifies a token we don't know about
	ILLEGAL = "ILLEGAL"
	// end of file
	EOF = "EOF"

	// Idenfifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, etc.
	INT   = "INT"   // 123456

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<" // less than
	GT       = ">" // greater than
	EQ       = "=="
	NOT_EQ   = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	//parenthesis + brackets
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

type Token struct {
	Type    TokenType
	Literal string
}

// map these keywords to their token types
// investigate
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

/**
Checks the keywords table to see whether the given identifier in in fact a keyword.

-if it is, return the keywords TokenType constant.
-if not, return token.IDENT, which is a TokenType for all user-defined identifiers

note on go-lang specific syntax used here:

if statements in Go can include both a condition and an initialization statement.
The example above uses both:
- initializes 'tok' with the value at keywords[ident];
- also initializes 'ok', which recieves a bool that will be set to true/false if 'ident' is present in the map

so if 'ident' is present in the map the body of the if statement will be executed and 'tok' will be available in the local scope.
**/
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

/**
Dev Notes:

Concepts:

Tokens:
- small, categorizable data structures that are fed into the parser
- Used to construct the AST (abstract syntax tree)

Other:
------------------------------
- The lexer reads the source code and characterizes what it finds into tokens
**/
