package lexer

import (
	"go-rilla/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `import "math" as m;

!-/*>=<=¿==

let five = 5;
let ten = 10;
let add = fn(x, y) {
x + y;
};

add(five, ten);

if (five < ten) {
return true;
} else {
return false;
}

if (a && b || c) {}

five += 1;
ten -= 1;

"foo bar"
[1, 2]
{"key": "value"}


m.sqrt(9) != 4
let decimal = 3.1415;
let bigger_int_part = 314.15;
3..1
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IMPORT, "import"},
		{token.STRING, "math"},
		{token.AS, "as"},
		{token.IDENTIFIER, "m"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.GREATER_EQUAL, ">="},
		{token.LESS_EQUAL, "<="},
		{token.ILLEGAL, "¿"},
		{token.EQUALS, "=="},

		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LEFT_PARENTHESIS, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.LEFT_BRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.IDENTIFIER, "add"},
		{token.LEFT_PARENTHESIS, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LEFT_PARENTHESIS, "("},
		{token.IDENTIFIER, "five"},
		{token.LESS_THAN, "<"},
		{token.IDENTIFIER, "ten"},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},

		{token.RIGHT_BRACE, "}"},
		{token.ELSE, "else"},
		{token.LEFT_BRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},

		{token.IF, "if"},
		{token.LEFT_PARENTHESIS, "("},
		{token.IDENTIFIER, "a"},
		{token.AND, "&&"},
		{token.IDENTIFIER, "b"},
		{token.OR, "||"},
		{token.IDENTIFIER, "c"},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RIGHT_BRACE, "}"},

		{token.IDENTIFIER, "five"},
		{token.SUM_ASSIGN, "+="},
		{token.INTEGER, "1"},
		{token.SEMICOLON, ";"},

		{token.IDENTIFIER, "ten"},
		{token.SUB_ASSIGN, "-="},
		{token.INTEGER, "1"},
		{token.SEMICOLON, ";"},

		{token.STRING, "foo bar"},
		{token.LEFT_BRACKET, "["},
		{token.INTEGER, "1"},
		{token.COMMA, ","},
		{token.INTEGER, "2"},
		{token.RIGHT_BRACKET, "]"},

		{token.LEFT_BRACE, "{"},
		{token.STRING, "key"},
		{token.COLON, ":"},
		{token.STRING, "value"},
		{token.RIGHT_BRACE, "}"},

		{token.IDENTIFIER, "m"},
		{token.DOT, "."},
		{token.IDENTIFIER, "sqrt"},
		{token.LEFT_PARENTHESIS, "("},
		{token.INTEGER, "9"},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.NOT_EQUAL, "!="},
		{token.INTEGER, "4"},

		{token.LET, "let"},
		{token.IDENTIFIER, "decimal"},
		{token.ASSIGN, "="},
		{token.FLOAT, "3.1415"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "bigger_int_part"},
		{token.ASSIGN, "="},
		{token.FLOAT, "314.15"},
		{token.SEMICOLON, ";"},

		{token.ILLEGAL, "3."},
		{token.DOT, "."},
		{token.INTEGER, "1"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
