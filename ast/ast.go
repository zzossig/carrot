package ast

import (
	"fmt"
	"strings"

	"github.com/zzossig/carrot/token"
)

// Expression is a type maker
type Expression interface {
	expression()
	String() string
}

// Group ::= selector [ COMMA S* selector ]*
type Group struct {
	Selectors []Expression
}

func (g *Group) expression() {}
func (g *Group) String() string {
	var sb strings.Builder
	for i, s := range g.Selectors {
		sb.WriteString(s.String())
		if i < len(g.Selectors)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

// Selector ::= simple_selector_sequence [ combinator simple_selector_sequence ]*
type Selector struct {
	Left  Expression
	Right Expression
	Token token.Token
}

func (s *Selector) expression() {}
func (s *Selector) String() string {
	var sb strings.Builder

	if s.Left != nil {
		sb.WriteString(s.Left.String())
	}

	switch s.Token.Type {
	case token.PLUS:
		fallthrough
	case token.GT:
		fallthrough
	case token.TILDE:
		sb.WriteString(" ")
		sb.WriteString(s.Token.Literal)
		sb.WriteString(" ")
	default:
		sb.WriteString(" ")
	}

	if s.Right != nil {
		sb.WriteString(s.Right.String())
	}

	return sb.String()
}

// Sequence ::= [ type_selector | universal ]
//     		 simple_sequence*
//   		 | simple_sequence+
type Sequence struct {
	Expression
	Exprs []Expression
}

func (s *Sequence) expression() {}
func (s *Sequence) String() string {
	var sb strings.Builder

	if s.Expression != nil {
		sb.WriteString(s.Expression.String())
	}
	for _, ss := range s.Exprs {
		sb.WriteString(ss.String())
	}

	return sb.String()
}

// Universal ::= '*'
type Universal struct {
	Token token.Token
}

func (u *Universal) expression() {}
func (u *Universal) String() string {
	return u.Token.Literal
}

// Class ::= '.' Name
type Class struct {
	Name string
}

func (c *Class) expression() {}
func (c *Class) String() string {
	return fmt.Sprintf(".%s", c.Name)
}

// Hash ::= '#' Name
type Hash struct {
	Name string
}

func (h *Hash) expression() {}
func (h *Hash) String() string {
	return fmt.Sprintf("#%s", h.Name)
}

// Attrib ::= '[' AttrExpr ']'
type Attrib struct {
	*AttrExpr
}

func (a *Attrib) expression() {}
func (a *Attrib) String() string {
	return fmt.Sprintf("[%s]", a.AttrExpr.String())
}

// AttrExpr ::= S* [ namespace_prefix ]? IDENT S*
//        [ [ PREFIXMATCH |
//            SUFFIXMATCH |
//            SUBSTRINGMATCH |
//            '=' |
//            INCLUDES |
//            DASHMATCH ] S* [ IDENT | STRING ] S*
//        ]?
type AttrExpr struct {
	Left, Right *Ident
	Token       token.Token
	TypeID      byte
}

func (ae *AttrExpr) expression() {}
func (ae *AttrExpr) String() string {
	switch ae.TypeID {
	case 1:
		return fmt.Sprintf("%s", ae.Left.String())
	case 2:
		return fmt.Sprintf("%s%s%s", ae.Left.String(), ae.Token.Literal, ae.Right.String())
	}
	return ""
}

// Negation ::= NOT S* negation_arg S* ')'
// negation_arg ::= ident | universal | HASH | class | attrib | pseudo
type Negation struct {
	*NArg
}

func (n *Negation) expression() {}
func (n *Negation) String() string {
	var sb strings.Builder
	sb.WriteString(":not(")
	sb.WriteString(n.NArg.String())
	sb.WriteString(")")
	return sb.String()
}

// NArg ::= type_selector | universal | HASH | class | attrib | pseudo
type NArg struct {
	*Ident
	*Universal
	*Hash
	*Class
	*Attrib
	*Pseudo
	TypeID byte
}

func (na *NArg) String() string {
	switch na.TypeID {
	case 1:
		return na.Ident.String()
	case 2:
		return na.Universal.String()
	case 3:
		return na.Hash.String()
	case 4:
		return na.Class.String()
	case 5:
		return na.Attrib.String()
	case 6:
		return na.Pseudo.String()
	}
	return ""
}

// Pseudo ::= ':' ':'? [ IDENT | functional_pseudo ]
type Pseudo struct {
	*Ident
	*FunctionalPseudo
	Token  token.Token
	TypeID byte
}

func (p *Pseudo) expression() {}
func (p *Pseudo) String() string {
	var sb strings.Builder
	sb.WriteString(p.Token.Literal)
	switch p.TypeID {
	case 1:
		sb.WriteString(p.Ident.String())
	case 2:
		sb.WriteString(p.FunctionalPseudo.String())
	}
	return sb.String()
}

// FunctionalPseudo ::= FUNCTION S* arg ')'
type FunctionalPseudo struct {
	Token token.Token
	*Arg
}

func (fp *FunctionalPseudo) expression() {}
func (fp *FunctionalPseudo) String() string {
	return fmt.Sprintf("%s(%s)", fp.Token.Literal, fp.Arg.String())
}

// Arg ::= DIMENSION | NUMBER | STRING | IDENT
type Arg struct {
	*Dimension
	*Number
	*Str
	*Ident
	TypeID byte
}

func (a *Arg) String() string {
	switch a.TypeID {
	case 1:
		return a.Dimension.String()
	case 2:
		return a.Number.String()
	case 3:
		return a.Str.String()
	case 4:
		return a.Ident.String()
	default:
		return ""
	}
}

// Number ::= int
type Number struct {
	Value int
}

func (n *Number) expression() {}
func (n *Number) String() string {
	return fmt.Sprintf("%d", n.Value)
}

// Ident ::= string
type Ident struct {
	Value string
}

func (i *Ident) expression() {}
func (i *Ident) String() string {
	return i.Value
}

// Str ::= string
type Str struct {
	Value string
}

func (s *Str) expression() {}
func (s *Str) String() string {
	return fmt.Sprintf("%q", s.Value)
}

// Dimension ::= an + b
type Dimension struct {
	A, B     int
	Aop, Bop string
}

func (d *Dimension) expression() {}
func (d *Dimension) String() string {
	var sb strings.Builder
	if d.Aop == "-" {
		sb.WriteString("-")
	}
	if d.A != 1 {
		sb.WriteString(fmt.Sprintf("%d", d.A))
	}
	sb.WriteString("n")
	if d.Bop != "" {
		sb.WriteString(d.Bop)
		sb.WriteString(fmt.Sprintf("%d", d.B))
	}
	return sb.String()
}
