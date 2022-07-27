package lexer

import (
	"github.com/smiksha1701/buggy/token"
)

type Lexer struct {
	input        string
	position     int
	ReadPosition int
	ch           byte
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespaces()
	switch l.ch {
	case '=':
		if l.PeekChar() == '=' {
			tok.Type = token.EQ
			tok.Literal = "=="
			l.ReadChar()
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERIX, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '!':
		if l.PeekChar() == '=' {
			tok.Type = token.NEQ
			tok.Literal = "!="
			l.ReadChar()
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.ReadString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if IsLetter(l.ch) {
			tok.Literal = l.ReadIdent()
			tok.Type = token.ChecKeywords(tok.Literal)
			return tok
		} else if IsNumber(l.ch) {
			tok.Literal = l.ReadNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.ReadChar()
	return tok
}
func (l *Lexer) skipWhitespaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.ReadChar()
	}
}
func newToken(TokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: TokenType, Literal: string(ch)}
}
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.ReadChar()
	return l
}
func (l *Lexer) PeekChar() byte {
	if l.ReadPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.ReadPosition]
	}
}
func (l *Lexer) ReadChar() {
	if l.ReadPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.ReadPosition]
	}
	l.position = l.ReadPosition
	l.ReadPosition += 1
}
func (l *Lexer) ReadString() string {
	position := l.position + 1
	for {
		l.ReadChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
func (l *Lexer) ReadIdent() string {
	start_pos := l.position
	for IsLetter(l.ch) {
		l.ReadChar()
	}
	return l.input[start_pos:l.position]
}
func (l *Lexer) ReadNumber() string {
	start_pos := l.position
	for IsNumber(l.ch) {
		l.ReadChar()
	}
	return l.input[start_pos:l.position]
}
func IsNumber(ch byte) bool {
	return (ch >= '0' && ch <= '9')
}
func IsLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch == '_')
}
