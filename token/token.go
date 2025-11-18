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
	FLOAT      = "FLOAT"

	// Operadores
	ASSIGN       = "="
	PLUS         = "+"
	BANG         = "!"
	MINUS        = "-"
	SLASH        = "/"
	STAR         = "*"
	LESS_THAN    = "<"
	GREATER_THAN = ">"
	PERCENT      = "%"

	// Operadores de dos caracteres
	LESS_EQUAL    = "<="
	GREATER_EQUAL = ">="
	EQUALS        = "=="
	NOT_EQUAL     = "!="
	AND           = "&&"
	OR            = "||"
	SUM_ASSIGN    = "+="
	SUB_ASSIGN    = "-="
	PLUS_PLUS     = "++"
	MINUS_MINUS   = "--"
	STAR_STAR     = "**"

	// Delimitadores
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	DOT       = "."

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

	// Bucles
	FOR      = "FOR"
	WHILE    = "WHILE"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
)

var keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"let":      LET,
	"return":   RETURN,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"while":    WHILE,
	"break":    BREAK,
	"continue": CONTINUE,
}

func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}
