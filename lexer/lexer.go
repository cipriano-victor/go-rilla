package lexer

import "go-rilla/token"

type Lexer struct {
	input        string
	position     int  // posición actual en input (apunta al carácter actual)
	readPosition int  // posición de lectura (apunta al siguiente carácter a leer)
	character    byte // carácter actual bajo revisión
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readCharacter()
	return l
}

func (l *Lexer) readCharacter() {
	if l.readPosition >= len(l.input) {
		l.character = 0
	} else {
		l.character = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.character {
	case '=':
		tok = newToken(token.ASSIGN, l.character)
	case '+':
		tok = newToken(token.PLUS, l.character)
	case '(':
		tok = newToken(token.LEFT_PARENTHESIS, l.character)
	case ')':
		tok = newToken(token.RIGHT_PARENTHESIS, l.character)
	case '{':
		tok = newToken(token.LEFT_BRACE, l.character)
	case '}':
		tok = newToken(token.RIGHT_BRACE, l.character)
	case ',':
		tok = newToken(token.COMMA, l.character)
	case ';':
		tok = newToken(token.SEMICOLON, l.character)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.character) {
			tok.Literal = l.read(isLetter)
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.character) {
			tok.Literal = l.read(isDigit)
			tok.Type = token.INTEGER
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.character)
		}
	}

	l.readCharacter()
	return tok
}

func newToken(tokenType token.TokenType, character byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(character)}
}

func (l *Lexer) skipWhitespace() {
	for l.character == ' ' || l.character == '\t' || l.character == '\n' || l.character == '\r' {
		l.readCharacter()
	}
}

func (l *Lexer) read(isValid func(byte) bool) string {
	position := l.position
	for isValid(l.character) {
		l.readCharacter()
	}
	return l.input[position:l.position]
}

func isLetter(character byte) bool {
	return 'a' <= character && character <= 'z' || 'A' <= character && character <= 'Z' || character == '_'
}

func isDigit(character byte) bool {
	return '0' <= character && character <= '9'
}
