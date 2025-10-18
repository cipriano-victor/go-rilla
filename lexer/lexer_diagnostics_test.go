package lexer

import (
	"go-rilla/token"
	"testing"
)

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
