package ast

import (
	"boar/token"
	"testing"
)

//Tests the following statement: let myVar = anotherVar;
/**
Notes:
- This test is mostly for demonstration purposes.
- Its to show how we can create an easily readable layer of tests
  for our parser by comparing the parser output with strings.
**/
func TestString(t *testing.T) {
	program := &Program{
		// slice of statements
		// test scenario: let myVar = anotherVar;
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() is wrong. got=%q", program.String())
	}
}
