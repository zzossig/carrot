package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

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
	errors []error

	curToken  token.Token
	peekToken token.Token
	peekSpace bool

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []error{},
	}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.prefixParseFns[token.STRING] = p.parseString
	p.prefixParseFns[token.NUM] = p.parseNumber
	p.prefixParseFns[token.IDENT] = p.parseSequence
	p.prefixParseFns[token.HASH] = p.parseSequence
	p.prefixParseFns[token.DOT] = p.parseSequence
	p.prefixParseFns[token.COLON] = p.parseSequence
	p.prefixParseFns[token.DCOLON] = p.parseSequence
	p.prefixParseFns[token.ASTERISK] = p.parseSequence
	p.prefixParseFns[token.LBRACKET] = p.parseSequence

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.infixParseFns[token.COMMA] = p.parseGroup
	p.infixParseFns[token.PLUS] = p.parseSelector
	p.infixParseFns[token.GT] = p.parseSelector
	p.infixParseFns[token.TILDE] = p.parseSelector

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseExpression() ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.newError("no prefix parse function for %s found", p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.EOF) && !p.peekTokenIs(token.ILLEGAL) {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	if p.peekTokenIs(token.ILLEGAL) {
		p.newError("parsing error: illegal token - %s", p.peekToken.Literal)
		return nil
	}

	return leftExp
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) newError(format string, a ...interface{}) {
	p.errors = append(p.errors, fmt.Errorf(format, a...))
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Errorf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	if p.peekTokenIs(token.ILLEGAL) {
		p.newError("parsing error: illegal token")
	}
	p.peekSpace = p.l.PeekSpace()
	p.peekToken = p.l.NextToken()
}

func (p *Parser) peekTokenIs(t token.Type, ts ...token.Type) bool {
	if len(ts) > 0 {
		for _, tt := range ts {
			if p.peekToken.Type == tt {
				return true
			}
		}
	}
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
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

func (p *Parser) parseString() ast.Expression {
	return &ast.Str{Value: p.curToken.Literal}
}

func (p *Parser) parseNumber() ast.Expression {
	num := &ast.Number{}

	value, _ := strconv.ParseInt(p.curToken.Literal, 0, 64)
	num.Value = int(value)

	return num
}

func (p *Parser) parseSequence() ast.Expression {
	seq := &ast.Sequence{}

	switch p.curToken.Type {
	case token.IDENT:
		seq.Expression = p.parseIdent()
	case token.ASTERISK:
		seq.Expression = p.parseUniversal()
	case token.HASH:
		seq.Exprs = append(seq.Exprs, p.parseHash())
	case token.DOT:
		seq.Exprs = append(seq.Exprs, p.parseClass())
	case token.LBRACKET:
		seq.Exprs = append(seq.Exprs, p.parseAttr())
	case token.COLON:
		fallthrough
	case token.DCOLON:
		seq.Exprs = append(seq.Exprs, p.parsePseudo())
	}

	if p.peekSpace {
		if !p.peekTokenIs(token.EOF, token.PLUS, token.GT, token.TILDE) {
			p.nextToken()
			selector := &ast.Selector{Left: seq, Token: token.TokenMap("w")}
			selector.Right = p.ParseExpression()
			return selector
		} else {
			p.nextToken()
			selector := &ast.Selector{Left: seq, Token: p.curToken}
			p.nextToken()
			selector.Right = p.ParseExpression()
			return selector
		}
	}

	return p.parseSimpleSequence(seq)
}

func (p *Parser) parseSimpleSequence(seq *ast.Sequence) ast.Expression {
	for p.peekTokenIs(token.HASH, token.DOT, token.LBRACKET, token.COLON, token.DCOLON) {
		if p.peekSpace {
			if !p.peekTokenIs(token.EOF, token.PLUS, token.GT, token.TILDE) {
				p.nextToken()
				selector := &ast.Selector{Left: seq, Token: token.TokenMap("w")}
				selector.Right = p.ParseExpression()
				return selector
			} else {
				p.nextToken()
				selector := &ast.Selector{Left: seq, Token: p.curToken}
				p.nextToken()
				selector.Right = p.ParseExpression()
				return selector
			}
		}

		switch p.peekToken.Type {
		case token.HASH:
			p.nextToken()
			seq.Exprs = append(seq.Exprs, p.parseHash())
		case token.DOT:
			p.nextToken()
			seq.Exprs = append(seq.Exprs, p.parseClass())
		case token.LBRACKET:
			p.nextToken()
			seq.Exprs = append(seq.Exprs, p.parseAttr())
		case token.COLON:
			fallthrough
		case token.DCOLON:
			p.nextToken()
			seq.Exprs = append(seq.Exprs, p.parsePseudo())
		}
	}

	return seq
}

func (p *Parser) parseIdent() ast.Expression {
	return &ast.Ident{Value: p.curToken.Literal}
}

func (p *Parser) parseHash() ast.Expression {
	return &ast.Hash{Name: p.curToken.Literal}
}

func (p *Parser) parseClass() ast.Expression {
	if !p.expectPeek(token.IDENT) {
		p.newError("parsing error: expectPeek=token.IDENT, got=%s", p.curToken.Type)
		return nil
	}

	return &ast.Class{Name: p.curToken.Literal}
}

func (p *Parser) parsePseudo() ast.Expression {
	psd := &ast.Pseudo{Token: p.curToken}

	switch p.peekToken.Type {
	case token.IDENT:
		p.nextToken()
		psd.Ident = &ast.Ident{Value: p.curToken.Literal}
		psd.TypeID = 1
	case token.FUNCTION:
		p.nextToken()
		if p.curToken.Literal == "not" {
			neg := &ast.Negation{}
			p.nextToken()
			neg.NArg = p.parseNArg()
			if !p.expectPeek(token.RPAREN) {
				return nil
			}
			return neg
		} else {
			fp := &ast.FunctionalPseudo{Token: p.curToken}
			p.nextToken()
			fp.Arg = p.parseArg()
			psd.FunctionalPseudo = fp
			psd.TypeID = 2
		}

		if !p.expectPeek(token.RPAREN) {
			return nil
		}
	default:
		p.errors = append(p.errors, fmt.Errorf("expected next token to be token.IDENT or token.FUNCTION, got %s instead", p.peekToken.Literal))
	}

	return psd
}

func (p *Parser) parseUniversal() ast.Expression {
	return &ast.Universal{Token: p.curToken}
}

func (p *Parser) parseGroup(left ast.Expression) ast.Expression {
	switch left := left.(type) {
	case *ast.Group:
		g := &ast.Group{}
		g.Selectors = append(g.Selectors, left.Selectors...)

		p.nextToken()
		right := p.ParseExpression()
		g.Selectors = append(g.Selectors, right)
		return g
	default:
		g := &ast.Group{}
		g.Selectors = append(g.Selectors, left)

		p.nextToken()
		right := p.ParseExpression()

		if r, ok := right.(*ast.Group); ok {
			g.Selectors = append(g.Selectors, r.Selectors...)
		} else {
			g.Selectors = append(g.Selectors, right)
		}

		return g
	}
}

func (p *Parser) parseSelector(left ast.Expression) ast.Expression {
	s := &ast.Selector{Left: left, Token: p.curToken}

	p.nextToken()
	s.Right = p.ParseExpression()

	return s
}

func (p *Parser) parseAttr() ast.Expression {
	attr := &ast.Attrib{}

	if !p.expectPeek(token.IDENT) {
		p.newError("parsing error: expectPeek=token.IDENT, got=%s", p.curToken)
		return nil
	}

	if p.peekTokenIs(token.EQ) ||
		p.peekTokenIs(token.INCLUDES) ||
		p.peekTokenIs(token.DASHMATCH) ||
		p.peekTokenIs(token.PREFIXMATCH) ||
		p.peekTokenIs(token.SUFFIXMATCH) ||
		p.peekTokenIs(token.SUBSTRINGMATCH) {
		ident := p.parseIdent().(*ast.Ident)
		attr.AttrExpr = p.parseAttrExpr(ident)
		if !p.expectPeek(token.RBRACKET) {
			return nil
		}

		return attr
	}

	ident := p.parseIdent().(*ast.Ident)
	attr.AttrExpr = &ast.AttrExpr{Left: ident, TypeID: 1}
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return attr
}

func (p *Parser) parseAttrExpr(left *ast.Ident) *ast.AttrExpr {
	ae := &ast.AttrExpr{Left: left}

	p.nextToken()

	switch p.curToken.Type {
	case token.PREFIXMATCH, token.SUFFIXMATCH, token.SUBSTRINGMATCH,
		token.EQ, token.INCLUDES, token.DASHMATCH:
		ae.TypeID = 2
		ae.Token = p.curToken

		p.nextToken()
		if p.curTokenIs(token.IDENT) {
			ae.Right = p.parseIdent().(*ast.Ident)
		} else if p.curTokenIs(token.STRING) {
			str := p.parseString().(*ast.Str)
			ae.Right = &ast.Ident{Value: str.Value}
		}
	default:
		ae.TypeID = 1
	}

	return ae
}

func (p *Parser) parseArg() *ast.Arg {
	arg := &ast.Arg{}

	switch p.curToken.Type {
	case token.STRING:
		arg.TypeID = 3
		arg.Str = p.parseString().(*ast.Str)
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
		return arg
	default:
		var sb strings.Builder
		for {
			sb.WriteString(p.curToken.Literal)
			if p.peekTokenIs(token.RPAREN) || p.peekTokenIs(token.EOF) {
				break
			}
			p.nextToken()
		}

		str := sb.String()

		numRe := regexp.MustCompile("^[-+]?[0-9]+$")
		if numRe.MatchString(str) {
			numMatch := regexp.MustCompile(`[0-9]+`)
			opMatch := regexp.MustCompile(`[+-]+`)

			nStr := numMatch.FindString(str)
			opStr := opMatch.FindString(str)

			arg.TypeID = 2
			num, _ := strconv.ParseInt(nStr, 0, 64)
			if opStr == "-" {
				arg.Number = &ast.Number{Value: -1 * int(num)}
			} else {
				arg.Number = &ast.Number{Value: int(num)}
			}

			return arg
		}

		dimRe := regexp.MustCompile(`^[-+]?[0-9]*[n]+[-+]?[0-9]*$`)
		if dimRe.MatchString(str) {
			arg.TypeID = 1
			arg.Dimension = p.parseDimension(str)
			return arg
		}

		identRe := regexp.MustCompile("^[A-Za-z]?[A-Za-z-_]*$")
		if identRe.MatchString(str) {
			arg.TypeID = 4
			arg.Ident = &ast.Ident{Value: str}
			return arg
		}

		return nil
	}
}

func (p *Parser) parseNArg() *ast.NArg {
	narg := &ast.NArg{}

	switch p.curToken.Type {
	case token.IDENT:
		narg.TypeID = 1
		narg.Ident = p.parseIdent().(*ast.Ident)
	case token.ASTERISK:
		narg.TypeID = 2
		narg.Universal = p.parseUniversal().(*ast.Universal)
	case token.HASH:
		narg.TypeID = 3
		narg.Hash = p.parseHash().(*ast.Hash)
	case token.DOT:
		if c := p.parseClass(); c != nil {
			narg.TypeID = 4
			narg.Class = c.(*ast.Class)
		}
	case token.LBRACKET:
		if a := p.parseAttr(); a != nil {
			narg.TypeID = 5
			narg.Attrib = a.(*ast.Attrib)
		}
	case token.COLON:
		if pd := p.parsePseudo(); pd != nil {
			pdTyped, ok := pd.(*ast.Pseudo)
			if ok {
				narg.TypeID = 6
				narg.Pseudo = pdTyped
			}
		}
	}

	return narg
}

func (p *Parser) parseDimension(str string) *ast.Dimension {
	d := &ast.Dimension{}

	numRe := regexp.MustCompile(`[0-9]+`)
	opRe := regexp.MustCompile(`[+-]+`)

	case1Re := regexp.MustCompile(`^[-+]?[0-9]*[A-Za-z]+$`)
	case2Re := regexp.MustCompile(`^[-+]?[A-Za-z]+[-+]?[0-9]+$`)
	if case1Re.MatchString(str) {
		n1 := numRe.FindString(str)
		op1 := opRe.FindString(str)
		if n1 != "" {
			n, _ := strconv.ParseInt(n1, 0, 64)
			d.A = int(n)
		} else {
			d.A = 1
		}

		if op1 != "" {
			d.Aop = op1
		} else {
			d.Aop = "+"
		}
	} else if case2Re.MatchString(str) {
		n1 := numRe.FindString(str)
		if n1 != "" {
			n, _ := strconv.ParseInt(n1, 0, 64)
			d.B = int(n)
			d.A = 1
		}

		ops := opRe.FindAllString(str, -1)
		if len(ops) == 2 {
			d.Aop = ops[0]
			d.Bop = ops[1]
		} else if len(ops) == 1 {
			d.Bop = ops[0]
		}
	} else {
		nums := numRe.FindAllString(str, -1)
		ops := opRe.FindAllString(str, -1)

		if len(nums) == 2 {
			n1, _ := strconv.ParseInt(nums[0], 0, 64)
			n2, _ := strconv.ParseInt(nums[1], 0, 64)
			d.A = int(n1)
			d.B = int(n2)
		} else if len(nums) == 1 {
			n1, _ := strconv.ParseInt(nums[0], 0, 64)
			d.A = int(n1)
		}

		if len(ops) == 2 {
			d.Aop = ops[0]
			d.Bop = ops[1]
		} else if len(ops) == 1 {
			d.Bop = ops[0]
		}
	}

	return d
}
