package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	// Identifiers + literals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"
	// Operators
	ASSIGN  = "="
	EQ      = "=="
	NEQ     = "!="
	BANG    = "!"
	PLUS    = "+"
	MINUS   = "-"
	SLASH   = "/"
	ASTERIX = "*"
	LT      = "<"
	GT      = ">"
	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LBRACKET  = "["
	RBRACKET  = "]"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	// Keywords
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func ChecKeywords(tok string) TokenType {
	if token, ok := keywords[tok]; ok {
		return token
	}
	return IDENT
}
