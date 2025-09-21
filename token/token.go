package token

import "go-rilla/source"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Range   source.Range
}

const (
	ILLEGAL = "ILLEGAL" // Caracter desconocido
	EOF     = "EOF"     // Delimitador de fin de archivo

	// Identificadores + Literales
	IDENTIFIER = "IDENTIFIER"
	INTEGER    = "INTEGER"
	STRING     = "STRING"

	// Operadores
	ASSIGN       = "="
	PLUS         = "+"
	BANG         = "!"
	MINUS        = "-"
	SLASH        = "/"
	ASTERISK     = "*"
	LESS_THAN    = "<"
	GREATER_THAN = ">"
	EQUALS       = "=="
	NOT_EQUAL    = "!="

	// Delimitadores
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	DOT       = "." // para acceso a miembros

	LEFT_PARENTHESIS  = "("
	RIGHT_PARENTHESIS = ")"
	LEFT_BRACE        = "{"
	RIGHT_BRACE       = "}"
	LEFT_BRACKET      = "["
	RIGHT_BRACKET     = "]"

	// Palabras Clave
	FUNCTION = "FUNCTION"
	LET      = "LET"
	RETURN   = "RETURN"

	// Booleanos
	TRUE  = "TRUE"
	FALSE = "FALSE"
	IF    = "IF"
	ELSE  = "ELSE"

	// Keywords para módulos y paquetes
	IMPORT = "IMPORT"
	AS     = "AS"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"import": IMPORT,
	"as":     AS,
}

func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}
