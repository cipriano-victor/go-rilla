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
	ASSIGN       = "="
	PLUS         = "+"
	BANG         = "!"
	MINUS        = "-"
	SLASH        = "/"
	ASTERISK     = "*"
	LESS_THAN    = "<"
	GREATER_THAN = ">"
	EQUALS       = "=="
	NOT_EQUALS   = "!="

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
	RETURN   = "RETURN"

	// Booleanos
	TRUE  = "TRUE"
	FALSE = "FALSE"
	IF    = "IF"
	ELSE  = "ELSE"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
}

func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}
