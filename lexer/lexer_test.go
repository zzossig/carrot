package lexer

import (
	"testing"

	"github.com/zzossig/carrot/token"
)

func TestNextToken(t *testing.T) {
	input := `
		.class
		"string"
		#id
		*
		div
		div.class
		div, p
		div p
		div+p
		div ~ div
		[attr]
		[attr="value"]
		[attr~=html]
		[attr^=head]
		[attr$=body]
		[attr|=main]
		[attr*=href]
		p:nth-child(2)
		. class
		`

	tokens := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.DOT, "."},
		{token.IDENT, "class"},
		{token.STRING, "string"},
		{token.HASH, "id"},
		{token.ASTERISK, "*"},
		{token.IDENT, "div"},
		{token.IDENT, "div"},
		{token.DOT, "."},
		{token.IDENT, "class"},
		{token.IDENT, "div"},
		{token.COMMA, ","},
		{token.IDENT, "p"},
		{token.IDENT, "div"},
		{token.IDENT, "p"},
		{token.IDENT, "div"},
		{token.PLUS, "+"},
		{token.IDENT, "p"},
		{token.IDENT, "div"},
		{token.TILDE, "~"},
		{token.IDENT, "div"},
		{token.LBRACKET, "["},
		{token.IDENT, "attr"},
		{token.RBRACKET, "]"},
		{token.LBRACKET, "["},
		{token.IDENT, "attr"},
		{token.EQ, "="},
		{token.STRING, "value"},
		{token.RBRACKET, "]"},
		{token.LBRACKET, "["},
		{token.IDENT, "attr"},
		{token.INCLUDES, "~="},
		{token.IDENT, "html"},
		{token.RBRACKET, "]"},
		{token.LBRACKET, "["},
		{token.IDENT, "attr"},
		{token.PREFIXMATCH, "^="},
		{token.IDENT, "head"},
		{token.RBRACKET, "]"},
		{token.LBRACKET, "["},
		{token.IDENT, "attr"},
		{token.SUFFIXMATCH, "$="},
		{token.IDENT, "body"},
		{token.RBRACKET, "]"},
		{token.LBRACKET, "["},
		{token.IDENT, "attr"},
		{token.DASHMATCH, "|="},
		{token.IDENT, "main"},
		{token.RBRACKET, "]"},
		{token.LBRACKET, "["},
		{token.IDENT, "attr"},
		{token.SUBSTRINGMATCH, "*="},
		{token.IDENT, "href"},
		{token.RBRACKET, "]"},
		{token.IDENT, "p"},
		{token.COLON, ":"},
		{token.FUNCTION, "nth-child"},
		{token.NUM, "2"},
		{token.RPAREN, ")"},
		{token.ILLEGAL, "."},
		{token.IDENT, "class"},
	}

	lexer := New(input)

	for i, tt := range tokens {
		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("TestNextToken:type[%d] - expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("TestNextToken:literal[%d] - expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
