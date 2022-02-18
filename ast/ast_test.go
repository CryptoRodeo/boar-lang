package ast

import (
	"monkey/token"
	"testing"
)

//Tests the following statement: let myVar = anotherVar;
func TestString(t *testing.T) {
	program := &Program{
		// slice of statements
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
