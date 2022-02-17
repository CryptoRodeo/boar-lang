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
**/
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok { //investigate
		return tok
	}
	return IDENT
}
