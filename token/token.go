package token


const (
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	IDENT = "IDENT" // variables
	INT = "INT" // integer data type
	STRING = "STRING"

	// Operators
	ASSIGN = "="
	PLUS = "+"
	BANG = "!"
	ASTERIK = "*"
	MINUS = "-"
	SLASH = "/"
	
	LT = "<"
	GT = ">"
	EQ = "=="
	NOT_EQ = "!="

	//Delimiters
	COLON = ":"
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACKET = "["
	RBRACKET = "]"

	//Keywords
	FUNCTION = "FUNCTION"
	LET = "LET"
	TRUE = "TRUE"
	FALSE = "FALSE"
	IF = "IF"
	ELSE = "ELSE"
	RETURN = "RETURN"
)

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

var keyword = map[string]TokenType { 
	"fn": FUNCTION,
	"let": LET,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keyword[ident]; ok {
		return tok
	}
	return IDENT
}
