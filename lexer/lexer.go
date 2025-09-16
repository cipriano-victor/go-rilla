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
		// default:
		// 	tok = newToken(token.ILLEGAL, l.character)
	}

	l.readCharacter()
	return tok
}

func newToken(tokenType token.TokenType, character byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(character)}
}
