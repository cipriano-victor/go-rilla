package lexer

import (
	"testing"
)

func TestIllegalChar(t *testing.T) {
	l := New("¿")
	_ = l.NextToken()
	diags := l.Diagnostics()
	if len(diags) == 0 {
		t.Fatalf("expected diagnostics, got none")
	}
	if diags[0].Code != "LEX001" {
		t.Fatalf("expected LEX001, got %s", diags[0].Code)
	}
}

func TestMalformedFloat(t *testing.T) {
	l := New("3.")
	_ = l.NextToken()
	diags := l.Diagnostics()
	if len(diags) == 0 {
		t.Fatalf("expected diagnostics, got none")
	}
	if diags[0].Code != "LEX002" {
		t.Fatalf("expected LEX002, got %s", diags[0].Code)
	}
}

func TestUnterminatedString(t *testing.T) {
	l := New("\"foo")
	_ = l.NextToken()
	diags := l.Diagnostics()
	if len(diags) == 0 {
		t.Fatalf("expected diagnostics, got none")
	}
	if diags[0].Code != "LEX003" {
		t.Fatalf("expected LEX003, got %s", diags[0].Code)
	}
}
