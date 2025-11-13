package lexer

import (
	"go-rilla/token"
	"testing"
	"unicode/utf8"
)

func TestNextToken(t *testing.T) {
	input := `!-/*>=<=¿==

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

{"key": "value"}

m.sqrt(9) != 4
let decimal = 3.1415;
let bigger_int_part = 314.15;
3..1

"foobar"
"foo bar"

[1, 2];
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
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

		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},

		{token.LEFT_BRACKET, "["},
		{token.INTEGER, "1"},
		{token.COMMA, ","},
		{token.INTEGER, "2"},
		{token.RIGHT_BRACKET, "]"},
		{token.SEMICOLON, ";"},

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

func TestStringEmptyNoDiag(t *testing.T) {
	l := New("\"\"")
	tok := l.NextToken()
	if tok.Type != token.STRING {
		t.Fatalf("expected STRING, got %s", tok.Type)
	}
	if ds := l.Diagnostics(); len(ds) != 0 {
		t.Fatalf("expected 0 diagnostics, got %d: %#v", len(ds), ds)
	}
}

func TestStringWithEscapedQuoteNoDiag(t *testing.T) {
	l := New("\"foo\\\"bar\"")
	tok := l.NextToken()
	if tok.Type != token.STRING {
		t.Fatalf("expected STRING, got %s", tok.Type)
	}
	if tok.Literal != "foo\\\"bar" {
		t.Fatalf("unexpected literal: %q", tok.Literal)
	}
	if ds := l.Diagnostics(); len(ds) != 0 {
		t.Fatalf("expected 0 diagnostics, got %d: %#v", len(ds), ds)
	}
}

func TestMalformedFloat(t *testing.T) {
	l := New("3.")
	tok := l.NextToken()
	if tok.Type != token.ILLEGAL {
		t.Fatalf("expected ILLEGAL, got %s", tok.Type)
	}
	ds := l.Diagnostics()
	if len(ds) == 0 {
		t.Fatalf("expected diagnostics, got none")
	}
	if ds[0].Code != "LEX002" {
		t.Fatalf("expected LEX002, got %s", ds[0].Code)
	}
}

func TestUnterminatedString(t *testing.T) {
	l := New("\"foo")
	tok := l.NextToken()
	if tok.Type != token.ILLEGAL {
		t.Fatalf("expected ILLEGAL, got %s", tok.Type)
	}
	ds := l.Diagnostics()
	if len(ds) == 0 {
		t.Fatalf("expected diagnostics, got none")
	}
	if ds[0].Code != "LEX003" {
		t.Fatalf("expected LEX003, got %s", ds[0].Code)
	}
}

func TestIllegalChar(t *testing.T) {
	l := New("¿")
	tok := l.NextToken()
	if tok.Type != token.ILLEGAL {
		t.Fatalf("expected ILLEGAL, got %s", tok.Type)
	}
	ds := l.Diagnostics()
	if len(ds) == 0 {
		t.Fatalf("expected diagnostics, got none")
	}
	if ds[0].Code != "LEX001" {
		t.Fatalf("expected LEX001, got %s", ds[0].Code)
	}
}

func TestInvalidEscape(t *testing.T) {
	l := New("\"foo\\z\"")
	tok := l.NextToken()
	if tok.Type != token.STRING {
		t.Fatalf("expected STRING, got %s", tok.Type)
	}
	ds := l.Diagnostics()
	if len(ds) == 0 {
		t.Fatalf("expected diagnostics, got none")
	}
	found := false
	for _, d := range ds {
		if d.Code == "LEX004" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected LEX004, got: %#v", ds)
	}
}

func TestInvalidUTF8(t *testing.T) {
	input := string([]byte{0xff})
	l := New(input)
	tok := l.NextToken()
	if tok.Type != token.ILLEGAL {
		t.Fatalf("expected ILLEGAL, got %s", tok.Type)
	}
	ds := l.Diagnostics()
	if len(ds) == 0 {
		t.Fatalf("expected diagnostics, got none")
	}
	if ds[0].Code != "LEX005" {
		t.Fatalf("expected LEX005, got %s", ds[0].Code)
	}
}

func FuzzLexerNeverPanicsAndRangesAreMonotonic(f *testing.F) {
	seeds := []string{
		``, `a`, `3.`, `"x`, `"foo\"bar"`, "¿",
		"let x = 1;", `import "math" m;`, "!;",
		string([]byte{0xff}),
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		_ = utf8.ValidString(s)

		l := New(s)

		const maxSteps = 10000
		prevOff := -1
		for i := 0; i < maxSteps; i++ {
			tok := l.NextToken()
			if tok.Range.Start.Offset < prevOff {
				t.Fatalf("token range offset went backwards: prev=%d, now=%d (tok=%+v)", prevOff, tok.Range.Start.Offset, tok)
			}
			prevOff = tok.Range.Start.Offset

			if tok.Type == token.EOF {
				break
			}
		}

		for _, d := range l.Diagnostics() {
			if d.Range.End.Offset <= d.Range.Start.Offset {
				t.Fatalf("diag with invalid offsets: %+v (input=%q)", d.Range, s)
			}
			if d.Range.Start.Line < 1 || d.Range.Start.Column < 1 ||
				d.Range.End.Line < 1 || d.Range.End.Column < 1 {
				t.Fatalf("diag with invalid line/col: %+v (input=%q)", d.Range, s)
			}
		}
	})
}
