package lexer

import (
	"unicode"

	"github.com/zzossig/carrot/token"
)

// Lexer reads input string one by one
type Lexer struct {
	input string // user input
	pos   int    // current position within input
	fPos  int    // following position
	ch    byte   // current char under examination
}

// New returns Lexer pointer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// PeekSpace checks if next char is space or not
func (l *Lexer) PeekSpace() bool {
	return unicode.IsSpace(rune(l.ch))
}

// NextToken returns next token by reading the input characters
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipSpace()

	switch l.ch {
	case '"', '\'':
		tok = token.Token{Type: token.STRING, Literal: l.readString()}
	case '+':
		tok = token.Token{Type: token.PLUS, Literal: "+"}
	case '*':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.SUBSTRINGMATCH, Literal: "*="}
		} else {
			tok = token.Token{Type: token.ASTERISK, Literal: "*"}
		}
	case '.':
		if !isLetter(l.peekChar()) {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
		} else {
			tok = token.Token{Type: token.DOT, Literal: "."}
		}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: ","}
	case ':':
		if l.peekChar() == ':' {
			l.readChar()
			tok = token.Token{Type: token.DCOLON, Literal: "::"}
		} else {
			tok = token.Token{Type: token.COLON, Literal: ":"}
		}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: "("}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: ")"}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: "["}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: "]"}
	case '|':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.DASHMATCH, Literal: "|="}
		} else {
			tok = token.Token{Type: token.VBAR, Literal: "|"}
		}
	case '>':
		tok = token.Token{Type: token.GT, Literal: ">"}
	case '=':
		tok = token.Token{Type: token.EQ, Literal: "="}
	case '~':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.INCLUDES, Literal: "~="}
		} else {
			tok = token.Token{Type: token.TILDE, Literal: "~"}
		}
	case '^':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.PREFIXMATCH, Literal: "^="}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
		}
	case '$':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.SUFFIXMATCH, Literal: "$="}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
		}
	case '@':
		if !isLetter(l.peekChar()) {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
		} else {
			l.readChar()
			tok = token.Token{Type: token.ATKEYWORD, Literal: l.readIdent()}
		}
	case '#':
		if !isLetter(l.peekChar()) {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
		} else {
			l.readChar()
			tok = token.Token{Type: token.HASH, Literal: l.readIdent()}
		}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdent()
			if l.ch == '(' {
				l.readChar()
				tok.Type = token.FUNCTION
			} else {
				tok.Type = token.IDENT
			}
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.NUM
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.fPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.fPos]
	}
	l.pos = l.fPos
	l.fPos++
}

func (l *Lexer) peekChar() byte {
	if l.fPos >= len(l.input) {
		return 0
	}
	return l.input[l.fPos]
}

func (l *Lexer) skipSpace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	pos := l.pos + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == '\'' || l.ch == 0 {
			break
		}
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readIdent() string {
	pos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '-'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
