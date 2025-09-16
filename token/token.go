package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" // Caracter desconocido
	EOF     = "EOF"     // Delimitador de fin de archivo

	// Identificadores + Literales
	IDENTIFIER = "IDENTIFIER"
	INTEGER    = "INTEGER"

	// Operadores
	ASSIGN = "="
	PLUS   = "+"

	// Delimitadores
	COMMA     = ","
	SEMICOLON = ";"

	LEFT_PARENTHESIS  = "("
	RIGHT_PARENTHESIS = ")"

	LEFT_BRACE  = "{"
	RIGHT_BRACE = "}"

	// Palabras Clave
	FUNCTION = "FUNCTION"
	LET      = "LET"
)
