package token

// Token represents lexical token
type Token struct {
	Type    Type
	Literal string
}

// Type represents Token Type
type Type string

// Token Types
const (
	ILLEGAL  Type = "-1"
	EOF      Type = "0"
	IDENT    Type = "ident"
	STRING   Type = "string"
	NUM      Type = "num"
	SQUOTE   Type = "'"
	DQUOTE   Type = "\""
	ASTERISK Type = "*"
	DOT      Type = "."
	COMMA    Type = ","
	COLON    Type = ":"
	DCOLON   Type = "::"
	LPAREN   Type = "("
	RPAREN   Type = ")"
	LBRACKET Type = "["
	RBRACKET Type = "]"
	VBAR     Type = "|"
	PLUS     Type = "+"
	MINUS    Type = "-"
	GT       Type = ">"
	EQ       Type = "="
	TILDE    Type = "~"
	S        Type = "w"

	INCLUDES       Type = "~="
	DASHMATCH      Type = "|="
	PREFIXMATCH    Type = "^="
	SUFFIXMATCH    Type = "$="
	SUBSTRINGMATCH Type = "*="
	ATKEYWORD      Type = "@{ident}"
	HASH           Type = "#{name}"
	FUNCTION       Type = "{ident}("
)

var tokens = map[string]Type{
	"'":  SQUOTE,
	"\"": DQUOTE,
	"*":  ASTERISK,
	".":  DOT,
	",":  COMMA,
	":":  COLON,
	"::": DCOLON,
	"(":  LPAREN,
	")":  RPAREN,
	"[":  LBRACKET,
	"]":  RBRACKET,
	"|":  VBAR,
	"+":  PLUS,
	"-":  MINUS,
	">":  GT,
	"=":  EQ,
	"~":  TILDE,
	"w":  S,
	"~=": INCLUDES,
	"|=": DASHMATCH,
	"^=": PREFIXMATCH,
	"$=": SUFFIXMATCH,
	"*=": SUBSTRINGMATCH,
}

// TokenMap ..
func TokenMap(t string) Token {
	if tok, ok := tokens[t]; ok {
		return Token{Type: tok, Literal: t}
	}
	return Token{Type: ILLEGAL, Literal: t}
}
