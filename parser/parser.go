package parser

import (
	"fmt"

	"github.com/zzossig/carrot/ast"
	"github.com/zzossig/carrot/lexer"
	"github.com/zzossig/carrot/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.prefixParseFns[token.STRING] = p.parseString
	p.prefixParseFns[token.IDENT] = p.parseIdent
	p.prefixParseFns[token.NUM] = p.parseNumber
	p.prefixParseFns[token.HASH] = p.parseHash
	p.prefixParseFns[token.DOT] = p.parseClass
	p.prefixParseFns[token.COLON] = p.parsePseudo
	p.prefixParseFns[token.DCOLON] = p.parsePseudo
	p.prefixParseFns[token.ASTERISK] = p.parseUniversal

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.infixParseFns[token.PLUS] = p.parsePlus
	p.infixParseFns[token.GT] = p.parseGreater
	p.infixParseFns[token.EQ] = p.parseEqual
	p.infixParseFns[token.TILDE] = p.parseTilde
	p.infixParseFns[token.INCLUDES] = p.parseIncludes
	p.infixParseFns[token.DASHMATCH] = p.parseDashMatch
	p.infixParseFns[token.PREFIXMATCH] = p.parsePrefixMatch
	p.infixParseFns[token.SUFFIXMATCH] = p.parseSuffixMatch
	p.infixParseFns[token.SUBSTRINGMATCH] = p.parseSubstringMatch

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression() ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseString() ast.Expression {
	return nil
}

func (p *Parser) parseIdent() ast.Expression {
	return nil
}

func (p *Parser) parseNumber() ast.Expression {
	return nil
}

func (p *Parser) parseHash() ast.Expression {
	return nil
}

func (p *Parser) parseClass() ast.Expression {
	return nil
}

func (p *Parser) parsePseudo() ast.Expression {
	return nil
}

func (p *Parser) parseUniversal() ast.Expression {
	return nil
}

func (p *Parser) parsePlus(left ast.Expression) ast.Expression {
	return nil
}

func (p *Parser) parseGreater(left ast.Expression) ast.Expression {
	return nil
}

func (p *Parser) parseEqual(left ast.Expression) ast.Expression {
	return nil
}

func (p *Parser) parseTilde(left ast.Expression) ast.Expression {
	return nil
}

func (p *Parser) parseIncludes(left ast.Expression) ast.Expression {
	return nil
}

func (p *Parser) parseDashMatch(left ast.Expression) ast.Expression {
	return nil
}

func (p *Parser) parsePrefixMatch(left ast.Expression) ast.Expression {
	return nil
}

func (p *Parser) parseSuffixMatch(left ast.Expression) ast.Expression {
	return nil
}

func (p *Parser) parseSubstringMatch(left ast.Expression) ast.Expression {
	return nil
}
