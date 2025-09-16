package lexer

import (
	"go-rilla/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LEFT_PARENTHESIS, "("},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RIGHT_BRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		token := l.NextToken()

		if token.Type != tt.expectedType {
			t.Fatalf("tests[%d] - TokenType Wrong. Expected=%q, Got=%q", i, tt.expectedType, token.Type)
		}

		if token.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] -  Literal Wrong. Expected=%q, Got=%q", i, tt.expectedLiteral, token.Literal)
		}
	}
}
