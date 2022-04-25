package lexer

import (
	"monkey/token"
)

//Struct to read "tokens"
type Lexer struct {
	input string // the entire string of characters that we've captured / 'source code'
	// current position in input (points to current character) points to the char in the input that corresponds to ch byte.
	position int
	// current position in reading (after current character), points to the "next" character in the input
	readPosition int
	//current char under examination
	ch byte
}

//Return a reference to a lexer struct value
func New(input string) *Lexer {
	// point to the new Lexer struct we're creating
	// initialize that struct with the source code we want to tokenize / lex
	l := &Lexer{input: input}
	// Lets make sure that our *Lexer is in a fully working state before anyone calls NextToken()
	// with l.ch, l.position and l.readPosition already initialized.
	l.readChar()
	return l //return the address of the new Lexer
}

/**
	- give us the next char
	- advances our position pointers used on the input string
**/
func (l *Lexer) readChar() {
	// If we've reached the end of the input
	if l.readPosition >= len(l.input) {
		// Set ch to 0 (ASCII for "NUL" char. Signifies nothing read or EOF)
		l.ch = 0
	} else {
		// Else, set l.ch to the next character
		l.ch = l.input[l.readPosition] //Access the specific char in the string using the current read position
	}
	// Move the current position in input to the next character
	l.position = l.readPosition
	// Increment so we point to the next char
	l.readPosition += 1
}

/**
	returns: Token struct for the current char we're lexing

	purpose:
	- Look at the current character under examination by the lexer (l.ch) and return a token of a specific type,
	depending on which character it is.
**/
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	// Ignore any whitespace found in the current char, (Monke-Lang doesn't add meaning to white spaces)
	l.skipWhitespace()

	// Read the char the lexer is currently on
	// tokenize it (figure out what it is)
	switch l.ch {
	case '=':
		// lets see if the next character is an equal sign
		if l.peekChar() == '=' {
			// save the current char so we don't lose it calling l.readChar()
			ch := l.ch
			// move to the next character (the other equal sign)
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			// save the ucrrent char so we don't lose it
			ch := l.ch
			// progress the position pointers
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		// reached EOF
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// This branch checks for identifiers whenever l.ch is not a recognized character.
		// ex: this could be the 'x' in 'let x = 5;' or also the 5 in that statement
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			/**
				The early exit here is necessary because when calling readIdentifier() we call readChar()
				repeatedly and advance our readPosition and position fields past the last character of the current
				identifier.

				So we don't need to call readChar after the switchStatement again.
			**/
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			// If we cant identify the char, consider it illegal.
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	// Read next character so l.ch is already updated when we call this method again.
	l.readChar()

	return tok
}

/**
  Reads the identifier and advances the lexer's position until it encounters a non-letter character.
**/
func (l *Lexer) readIdentifier() string {
	// position where we first encountered the potential identifier
	position := l.position
	// while the current character is a letter lets read each character and advance our lexers position
	for isLetter(l.ch) {
		l.readChar()
	}

	/*
		return the subset of the string at these positions
		position being the index of when we first found our identifier
		l.position being the index right before we're no longer reading a character
	*/
	return l.input[position:l.position]
}

// Checks for whether a char is a letter
/**
note:
- because we'll consider _ as a letter we can allow it in identifiers and keywords.
- this means we can use variables with names like foo_bar
- we can also sneak in other identifiers like ! and ? here too.
**/
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

/**
params:
- tokenType
- character byte

returns: Token
**/
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Skips any whitespace so our lexer can ignore it.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		// Skip to the next character
		l.readChar()
	}
}

/**
note:

- We only read ints here, not floats, hex notation, octal, etc.
This is to keep things simple...for now :)
**/
func (l *Lexer) readNumber() string {
	position := l.position
	// if the character is a digit
	for isDigit(l.ch) {
		// update the position of the lexer
		l.readChar()
	}

	// return the subset of the string at these positions
	/*
		position being the index of when we first found our number
		l.position being the index of when its no longer a number
	*/
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Allows us to look ahead in the input but not move around it.
func (l *Lexer) peekChar() byte {
	// if we've reached EOF, return NULL
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readString() string {
	// skip over the " char
	position := l.position + 1

	// read characters until we reach the end of the string
	for {
		l.readChar()

		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	// return chars that are in between the quotes
	return l.input[position:l.position]
}

/**
Dev Notes:

Concepts:

lexical analysis / lexing:
------------------------------
- Transforms source code to tokens via lexing / lexical analysis
- This is what is being done in NextToken()

Lexer:
- what does the lexical anaylsis / lexing
- Its job is not to tell us whether the code make sense, worse or contains error. It should only turn input into tokens

Example of source code being lexed into tokens:

input = "let x = 5 + 5;"

result would look something like:
[
  LET,
  IDENTIFIER("x"),
  EQUAL_SIGN,
  INTEGER(5),
  PLUS_SIGN,
  INTEGER(5),
  SEMICOLON
]


Other:
------------------------------
- The lexer only supports ASCII characters instead of the full Unicode range.
- This lets us keep things simple.
- In order to fully support Unicode and UTF-8 we would need to:
  - change l.ch from a byte to a rune
	- change the way we read the next characters, since they would be multiple bytes not.
	- Using l.input[l.readPosition] wouldn't work anymore..

	**/
