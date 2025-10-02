package ast

import (
	"go-rilla/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
			// &ReturnStatement{
			// 	Token: token.Token{Type: token.RETURN, Literal: "return"},
			// 	ReturnValue: &Identifier{
			// 		Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
			// 		Value: "myVar",
			// 	},
			// },
			// &ExpressionStatement{
			// 	Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
			// 	Expression: &Identifier{
			// 		Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
			// 		Value: "myVar",
			// 	},
			// },
		},
	}

	expected := "let myVar = anotherVar;"
	if program.String() != expected {
		t.Errorf("program.String() wrong. expected=%q, got=%q",
			expected, program.String())
	}
}
