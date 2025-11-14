package lexer

import (
	"go-rilla/diag"
	"go-rilla/source"
	"go-rilla/token"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	input             string
	offset            int             // posición actual en input (apunta al carácter actual)
	readOffset        int             // posición de lectura (apunta al siguiente carácter a leer)
	character         rune            // carácter actual bajo revisión
	position          source.Position // posición del carácter actual (línea y columna)
	diagnostics       []diag.Diagnostic
	lastDecodeInvalid bool
}

// Diagnostics devuelve los diagnósticos léxicos acumulados.
func (l *Lexer) Diagnostics() []diag.Diagnostic { return l.diagnostics }

func (l *Lexer) addDiag(level diag.Level, code, msg, hint string, start, end source.Position) {
	l.diagnostics = append(l.diagnostics, diag.Diagnostic{
		Level:   level,
		Code:    code,
		Message: msg,
		Hint:    hint,
		Range:   source.Range{Start: start, End: end},
	})
}

func New(input string) *Lexer {
	l := &Lexer{input: input, position: source.Position{Line: 1, Column: 1}}
	l.readCharacter()
	return l
}

func (l *Lexer) readCharacter() {
	// avanzar a la siguiente runa
	if l.readOffset >= len(l.input) {
		// EOF virtual
		if l.character == '\n' {
			l.position.Line++
			l.position.Column = 1
		} else if l.offset != 0 {
			l.position.Column++
		}
		l.offset = l.readOffset
		l.character = 0
		return
	}

	// actualizar posición según la runa anterior
	if l.offset != 0 || l.position.Line != 1 || l.position.Column != 1 {
		if l.character == '\n' {
			l.position.Line++
			l.position.Column = 1
		} else {
			l.position.Column++
		}
	}

	r, w := utf8.DecodeRuneInString(l.input[l.readOffset:])
	l.character = r

	if r == utf8.RuneError && w == 1 {
		l.lastDecodeInvalid = true
	} else {
		l.lastDecodeInvalid = false
	}

	l.offset = l.readOffset
	l.readOffset += w
}

func (l *Lexer) peekCharacter() rune {
	if l.readOffset >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readOffset:])
	return r
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	start := l.currentStart()

	switch l.character {
	case '&':
		return makeTwoCharacterToken(l, '&', token.AND, token.ILLEGAL, start)
	case '|':
		return makeTwoCharacterToken(l, '|', token.OR, token.ILLEGAL, start)
	case '=':
		return makeTwoCharacterToken(l, '=', token.EQUALS, token.ASSIGN, start)
	case '!':
		return makeTwoCharacterToken(l, '=', token.NOT_EQUAL, token.BANG, start)
	case '+':
		if l.peekCharacter() == '+' {
			return makeTwoCharacterToken(l, '+', token.PLUS_PLUS, token.PLUS, start)
		}
		return makeTwoCharacterToken(l, '=', token.SUM_ASSIGN, token.PLUS, start)
	case '(':
		tok = newToken(token.LEFT_PARENTHESIS, l.character, start, l.afterCurrent())
		l.readCharacter()
		return tok
	case ')':
		tok = newToken(token.RIGHT_PARENTHESIS, l.character, start, l.afterCurrent())
		l.readCharacter()
		return tok
	case '{':
		tok = newToken(token.LEFT_BRACE, l.character, start, l.afterCurrent())
		l.readCharacter()
		return tok
	case '}':
		tok = newToken(token.RIGHT_BRACE, l.character, start, l.afterCurrent())
		l.readCharacter()
		return tok
	case ',':
		tok = newToken(token.COMMA, l.character, start, l.afterCurrent())
		l.readCharacter()
		return tok
	case ';':
		tok = newToken(token.SEMICOLON, l.character, start, l.afterCurrent())
		l.readCharacter()
		return tok
	case '-':
		if l.peekCharacter() == '-' {
			return makeTwoCharacterToken(l, '-', token.MINUS_MINUS, token.MINUS, start)
		}
		return makeTwoCharacterToken(l, '=', token.SUB_ASSIGN, token.MINUS, start)
	case '/':
		tok = newToken(token.SLASH, l.character, start, l.afterCurrent())
		l.readCharacter()
		return tok
	case '*':
		return makeTwoCharacterToken(l, '*', token.STAR_STAR, token.STAR, start)
	case '<':
		return makeTwoCharacterToken(l, '=', token.LESS_EQUAL, token.LESS_THAN, start)
	case '>':
		return makeTwoCharacterToken(l, '=', token.GREATER_EQUAL, token.GREATER_THAN, start)
	case ':':
		tok = newToken(token.COLON, l.character, start, l.afterCurrent())
		l.readCharacter()
		return tok
	case '.':
		tok = newToken(token.DOT, l.character, start, l.afterCurrent())
		l.readCharacter()
		return tok
	case '[':
		tok = newToken(token.LEFT_BRACKET, l.character, start, l.afterCurrent())
		l.readCharacter()
		return tok
	case ']':
		tok = newToken(token.RIGHT_BRACKET, l.character, start, l.afterCurrent())
		l.readCharacter()
		return tok
	case '"':
		s, closed := l.readString()
		end := l.currentStart()
		if !closed {
			l.addDiag(diag.Error, "LEX003", "String without closing quote", "Missing closing quote '\"'", start, end)
			return token.Token{Type: token.ILLEGAL, Literal: s, Range: source.Range{Start: start, End: end}}
		}
		l.readCharacter() // consumir comilla de cierre
		return token.Token{Type: token.STRING, Literal: s, Range: source.Range{Start: start, End: end}}

	case 0:
		return token.Token{Type: token.EOF, Literal: "", Range: source.Range{Start: start, End: start}}
	default:
		if isLetter(l.character) {
			literal := l.read(isLetter)
			return token.Token{Type: token.LookupIdentifier(literal), Literal: literal, Range: source.Range{Start: start, End: l.currentStart()}}
		}
		if isDigit(l.character) {
			intPart := l.read(isDigit)
			if l.character == '.' {
				if isDigit(l.peekCharacter()) {
					l.readCharacter()
					decimalPart := l.read(isDigit)
					literal := intPart + "." + decimalPart
					return token.Token{Type: token.FLOAT, Literal: literal, Range: source.Range{Start: start, End: l.currentStart()}}
				}
				// si no hay dígito tras el punto, no es un float válido
				literal := intPart + string(l.character)
				end := l.afterCurrent()
				l.addDiag(diag.Error, "LEX002", "Malformed float literal", "At least one digit is expected after the decimal point", start, end)
				l.readCharacter()
				return token.Token{Type: token.ILLEGAL, Literal: literal, Range: source.Range{Start: start, End: end}}
			}
			if isLetter(l.character) {
				l.addDiag(diag.Error, "LEX006", "Malformed number literal", "Unexpected character after number", start, l.currentStart())
				return token.Token{Type: token.ILLEGAL, Literal: intPart, Range: source.Range{Start: start, End: l.currentStart()}}
			}
			return token.Token{Type: token.INTEGER, Literal: intPart, Range: source.Range{Start: start, End: l.currentStart()}}
		}
		if l.lastDecodeInvalid {
			end := l.afterCurrent()
			tok := token.Token{Type: token.ILLEGAL, Literal: string(l.character), Range: source.Range{Start: start, End: end}}
			l.addDiag(diag.Error, "LEX005", "Invalid UTF-8 byte", "carácter no decodificable", start, end)
			l.readCharacter()
			return tok
		}

		end := l.afterCurrent()
		tok := newToken(token.ILLEGAL, l.character, start, end)
		l.addDiag(diag.Error, "LEX001", "Illegal Character", "Character not recognized by the language", start, l.afterCurrent())
		l.readCharacter()
		return tok
	}
}

func newToken(tokenType token.TokenType, character rune, start, end source.Position) token.Token {
	return token.Token{Type: tokenType, Literal: string(character), Range: source.Range{Start: start, End: end}}
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.character)) {
		l.readCharacter()
	}
}

func (l *Lexer) read(isValid func(rune) bool) string {
	offset := l.offset
	for isValid(l.character) {
		l.readCharacter()
	}
	return l.input[offset:l.offset]
}

func (l *Lexer) readString() (string, bool) {
	l.readCharacter() // consumir la comilla de apertura
	startContent := l.offset
	for l.character != '"' && l.character != 0 {
		if l.character == '\\' {
			escStart := l.currentStart()
			next := l.peekCharacter()
			switch next {
			case '"', '\\', 'n', 't', 'r':
				l.readCharacter()
			default:
				escEnd := l.afterCurrent()
				l.addDiag(diag.Error, "LEX004", "Invalid escape sequence",
					"Use \\\" \\\\ \\n \\t or \\r", escStart, escEnd)
				if next != 0 {
					l.readCharacter()
				}
			}
		}
		l.readCharacter()
	}
	closed := (l.character == '"')
	return l.input[startContent:l.offset], closed
}

func (l *Lexer) currentStart() source.Position {
	return source.Position{Offset: l.offset, Line: l.position.Line, Column: l.position.Column}
}

func (l *Lexer) afterCurrent() source.Position {
	// Posición inmediatamente después del rune actual
	// (End exclusivo) — columna aproximada +1
	col := l.position.Column
	if l.character == '\n' {
		// tras un \n, la próxima runa estará en la línea siguiente, col=1
		return source.Position{Offset: l.readOffset, Line: l.position.Line + 1, Column: 1}
	}
	return source.Position{Offset: l.readOffset, Line: l.position.Line, Column: col + 1}
}

func isLetter(character rune) bool {
	return ('a' <= character && character <= 'z') || ('A' <= character && character <= 'Z') || character == '_'
}

func isDigit(character rune) bool {
	return '0' <= character && character <= '9'
}

func makeTwoCharacterToken(l *Lexer, expected rune, twoCharType, oneCharType token.TokenType, start source.Position) token.Token {
	if l.peekCharacter() == expected {
		first := l.character
		l.readCharacter()
		literal := string(first) + string(l.character)
		end := l.afterCurrent()
		l.readCharacter()
		return token.Token{Type: twoCharType, Literal: literal, Range: source.Range{Start: start, End: end}}
	}
	literal := string(l.character)
	end := l.afterCurrent()
	l.readCharacter()
	return token.Token{Type: oneCharType, Literal: literal, Range: source.Range{Start: start, End: end}}
}
