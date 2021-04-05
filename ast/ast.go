package ast

import "github.com/zzossig/carrot/token"

// Group ::= selector [ COMMA S* selector ]*
type Group struct {
	Selectors []Selector
}

// Selector ::= simple_selector_sequence [ combinator simple_selector_sequence ]*
type Selector struct {
	Left  SSS
	Right SSS
	Token token.Token
}

// SSS ::= [ type_selector | universal ]
//     		 [ HASH | class | attrib | pseudo | negation ]*
//   		 | [ HASH | class | attrib | pseudo | negation ]+
type SSS struct {
	TypeSelector
	Universal
	Hash
	Class
	Attrib
	Pseudo
	Negation
	TypeID byte
}

// TypeSelector ::= element_name
type TypeSelector struct {
	ElementName
}

// element_name ::= IDENT
type ElementName struct {
	Ident
}

// Universal ::= '*'
type Universal struct {
	Token token.Token
}

// Class ::= '.' IDENT
type Class struct {
	Ident
}

// Hash ::= '#' Ident
type Hash struct {
	Ident
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

// Pseudo ::= ':' ':'? [ IDENT | functional_pseudo ]
type Pseudo struct {
	Ident
	FunctionalPseudo
	Token  token.Token
	TypeID byte
}

// FunctionalPseudo
type FunctionalPseudo struct {
	Expression
}

// Expression ::= DIMENSION | NUMBER | STRING | IDENT
type Expression struct {
	Ident
	String
	Number
	Dimension
	TypeID byte
}

// Number ::= int
type Number struct {
	Value int
}

// Ident ::= string
type Ident struct {
	Value string
}

// String ::= string
type String struct {
	Value string
}

// Dimension ::= an + b
type Dimension struct {
	A, B int
}
