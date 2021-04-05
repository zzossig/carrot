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
	Selectors []Selector
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
	Left  TypeSelector
	Right TypeSelector
	Token token.Token
}

func (s *Selector) expression() {}
func (s *Selector) String() string {
	var sb strings.Builder

	sb.WriteString(s.Left.String())
	switch s.Token.Type {
	case token.PLUS:
		fallthrough
	case token.GT:
		fallthrough
	case token.TILDE:
		sb.WriteString(s.Token.Literal)
	default:
		sb.WriteString(" ")
	}
	sb.WriteString(s.Right.String())

	return sb.String()
}

// Sequence ::= [ type_selector | universal ]
//     		 simple_sequence*
//   		 | simple_sequence+
type Sequence struct {
	TypeSelector
	Universal
	Seq    []SimpleSelector
	TypeID byte
}

func (s *Sequence) expression() {}
func (s *Sequence) String() string {
	var sb strings.Builder

	switch s.TypeID {
	case 1:
		sb.WriteString(s.TypeSelector.String())
	case 2:
		sb.WriteString(s.Universal.String())
	}

	for _, ss := range s.Seq {
		sb.WriteString(ss.String())
	}

	return sb.String()
}

// SimpleSelector ::= [ HASH | class | attrib | pseudo | negation ]
type SimpleSelector struct {
	Hash
	Class
	Attrib
	Pseudo
	Negation
	TypeID byte
}

func (ss *SimpleSelector) expression() {}
func (ss *SimpleSelector) String() string {
	switch ss.TypeID {
	case 1:
		return ss.Hash.String()
	case 2:
		return ss.Class.String()
	case 3:
		return ss.Attrib.String()
	case 4:
		return ss.Pseudo.String()
	case 5:
		return ss.Negation.String()
	}
	return ""
}

// TypeSelector ::= element_name
type TypeSelector struct {
	ElementName
}

func (ts *TypeSelector) expression() {}
func (ts *TypeSelector) String() string {
	return ts.ElementName.String()
}

// element_name ::= IDENT
type ElementName struct {
	Ident
}

func (en *ElementName) expression() {}
func (en *ElementName) String() string {
	return en.Ident.String()
}

// Universal ::= '*'
type Universal struct {
	Token token.Token
}

func (u *Universal) expression() {}
func (u *Universal) String() string {
	return u.Token.Literal
}

// Class ::= '.' IDENT
type Class struct {
	Ident
}

func (c *Class) expression() {}
func (c *Class) String() string {
	return fmt.Sprintf(".%s", c.Ident.String())
}

// Hash ::= '#' Ident
type Hash struct {
	Ident
}

func (h *Hash) expression() {}
func (h *Hash) String() string {
	return fmt.Sprintf("#%s", h.Ident.String())
}

// Attrib ::= '[' S* [ namespace_prefix ]? IDENT S*
//        [ [ PREFIXMATCH |
//            SUFFIXMATCH |
//            SUBSTRINGMATCH |
//            '=' |
//            INCLUDES |
//            DASHMATCH ] S* [ IDENT | STRING ] S*
//        ]? ']'
type Attrib struct {
	Left, Right Ident
	Token       token.Token
	TypeID      byte
}

func (a *Attrib) expression() {}
func (a *Attrib) String() string {
	switch a.TypeID {
	case 1:
		return fmt.Sprintf("[%s]", a.Left.String())
	case 2:
		return fmt.Sprintf("[%s%s%s]", a.Left.Value, a.Token.Literal, a.Right.Value)
	}
	return ""
}

// Negation ::= NOT S* negation_arg S* ')'
// negation_arg ::= type_selector | universal | HASH | class | attrib | pseudo
type Negation struct {
	TypeSelector
	Universal
	Hash
	Class
	Attrib
	Pseudo
	TypeID byte
}

func (n *Negation) expression() {}
func (n *Negation) String() string {
	var sb strings.Builder
	sb.WriteString(":not(")
	switch n.TypeID {
	case 1:
		sb.WriteString(n.TypeSelector.String())
	case 2:
		sb.WriteString(n.Universal.String())
	case 3:
		sb.WriteString(n.Hash.String())
	case 4:
		sb.WriteString(n.Class.String())
	case 5:
		sb.WriteString(n.Attrib.String())
	case 6:
		sb.WriteString(n.Pseudo.String())
	}
	sb.WriteString(")")
	return sb.String()
}

// Pseudo ::= ':' ':'? [ IDENT | functional_pseudo ]
type Pseudo struct {
	Ident
	FunctionalPseudo
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

// FunctionalPseudo ::= FUNCTION S* expression ')'
type FunctionalPseudo struct {
	Token token.Token
	Arg
}

func (fp *FunctionalPseudo) expression() {}
func (fp *FunctionalPseudo) String() string {
	return fmt.Sprintf("%s(%s)", fp.Token.Literal, fp.Arg.String())
}

// Arg ::= DIMENSION | NUMBER | STRING | IDENT
type Arg struct {
	Ident
	Str
	Number
	Dimension
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

func (n *Number) String() string {
	return fmt.Sprintf("%d", n.Value)
}

// Ident ::= string
type Ident struct {
	Value string
}

func (i *Ident) String() string {
	return i.Value
}

// Str ::= string
type Str struct {
	Value string
}

func (s *Str) String() string {
	return s.Value
}

// Dimension ::= an + b
type Dimension struct {
	A, B int
}

func (d *Dimension) String() string {
	return fmt.Sprintf("%dn + %d", d.A, d.B)
}
