package eval

import (
	"strings"

	"github.com/zzossig/carrot/ast"
	"github.com/zzossig/carrot/object"
	"github.com/zzossig/carrot/token"
	"golang.org/x/net/html"
)

func Eval(expr ast.Expression, ctx *object.Context) []*html.Node {
	switch expr := expr.(type) {
	case *ast.Group:
		return evalGroup(expr)
	case *ast.Selector:
		return evalSelector(expr, ctx)
	case *ast.Sequence:
		return evalSequence(expr, ctx)
	case *ast.Universal:
		return evalUniversal(expr, ctx)
	case *ast.Class:
		return evalClass(expr, ctx)
	case *ast.Hash:
		return evalHash(expr, ctx)
	case *ast.Attrib:
		return evalAttrib(expr, ctx)
	case *ast.Ident:
		return evalIdent(expr, ctx)
	case *ast.Pseudo:
		return evalPseudo(expr, ctx)
	case *ast.Negation:
		return evalNegation(expr, ctx)
	}
	return nil
}

func evalGroup(expr ast.Expression) []*html.Node {
	g := expr.(*ast.Group)
	var nodes []*html.Node

	for _, selector := range g.Selectors {
		c := object.NewContext()
		e := Eval(selector, c)
		for _, ee := range e {
			nodes = appendNode(nodes, ee)
		}
	}

	return nodes
}

func evalSelector(expr ast.Expression, ctx *object.Context) []*html.Node {
	s := expr.(*ast.Selector)

	leftNodes := Eval(s.Left, ctx)
	ctx.CNode = leftNodes

	switch s.Token.Type {
	case token.TILDE:
		ctx.CNode = collectSubSibling(ctx)
	case token.PLUS:
		ctx.CNode = collectNextSibling(ctx)
	case token.GT:
		ctx.CNode = collectChild(ctx)
	case token.S:
		ctx.CNode = collectDesc(ctx)
	}

	rightNodes := Eval(s.Right, ctx)
	ctx.CNode = rightNodes

	return ctx.CNode
}

func evalSequence(expr ast.Expression, ctx *object.Context) []*html.Node {
	s := expr.(*ast.Sequence)

	if s.Expression != nil {
		h := Eval(s.Expression, ctx)
		ctx.CNode = h
	}

	for _, e := range s.Exprs {
		ss := Eval(e, ctx)
		ctx.CNode = ss
	}

	return ctx.CNode
}

func evalUniversal(expr ast.Expression, ctx *object.Context) []*html.Node {
	return ctx.CNode
}

func evalClass(expr ast.Expression, ctx *object.Context) []*html.Node {
	c := expr.(*ast.Class)
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		for _, a := range n.Attr {
			if a.Key == "class" {
				f := strings.Fields(a.Val)
				for _, s := range f {
					if s == c.Name {
						nodes = append(nodes, n)
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	return nodes
}

func evalHash(expr ast.Expression, ctx *object.Context) []*html.Node {
	h := expr.(*ast.Hash)
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		for _, a := range n.Attr {
			if a.Key == "id" && a.Val == h.Name {
				nodes = append(nodes, n)
			}
		}
	}

	ctx.CNode = nodes
	return nodes
}

func evalAttrib(expr ast.Expression, ctx *object.Context) []*html.Node {
	ae := expr.(*ast.Attrib).AttrExpr
	var nodes []*html.Node

	switch ae.TypeID {
	case 1:
		for _, n := range ctx.CNode {
			for _, a := range n.Attr {
				if a.Key == ae.Left.Value {
					nodes = append(nodes, n)
				}
			}
		}
	case 2:
		switch ae.Token.Type {
		case token.EQ:
			for _, n := range ctx.CNode {
				for _, a := range n.Attr {
					if a.Key == ae.Left.Value && a.Val == ae.Right.Value {
						nodes = append(nodes, n)
					}
				}
			}
		case token.PREFIXMATCH:
			for _, n := range ctx.CNode {
				for _, a := range n.Attr {
					if a.Key == ae.Left.Value {
						if strings.HasPrefix(a.Val, ae.Right.Value) {
							nodes = append(nodes, n)
						}
					}
				}
			}
		case token.SUFFIXMATCH:
			for _, n := range ctx.CNode {
				for _, a := range n.Attr {
					if a.Key == ae.Left.Value {
						if strings.HasSuffix(a.Val, ae.Right.Value) {
							nodes = append(nodes, n)
						}
					}
				}
			}
		case token.SUBSTRINGMATCH:
			for _, n := range ctx.CNode {
				for _, a := range n.Attr {
					if a.Key == ae.Left.Value {
						if strings.Contains(a.Val, ae.Right.Value) {
							nodes = append(nodes, n)
						}
					}
				}
			}
		case token.INCLUDES:
			for _, n := range ctx.CNode {
				for _, a := range n.Attr {
					if a.Key == ae.Left.Value {
						f := strings.Fields(a.Val)
						for _, s := range f {
							if s == ae.Right.Value {
								nodes = append(nodes, n)
							}
						}
					}
				}
			}
		case token.DASHMATCH:
			for _, n := range ctx.CNode {
				for _, a := range n.Attr {
					if a.Key == ae.Left.Value {
						if strings.HasPrefix(a.Val, "en") {
							nodes = append(nodes, n)
						}
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	return nodes
}

func evalIdent(expr ast.Expression, ctx *object.Context) []*html.Node {
	i := expr.(*ast.Ident)
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		if n.Data == i.Value {
			nodes = append(nodes, n)
		}
	}

	ctx.CNode = nodes
	return nodes
}

func evalNegation(expr ast.Expression, ctx *object.Context) []*html.Node {
	na := expr.(*ast.Negation).NArg
	var nodes []*html.Node

	switch na.TypeID {
	case 1:
		for _, n := range ctx.CNode {
			if n.Data != na.Ident.Value {
				nodes = append(nodes, n)
			}
		}
	case 2:
		// do nothing
	case 3:
		hasID := false
		idVal := ""

		for _, n := range ctx.CNode {
			for _, a := range n.Attr {
				if a.Key == "id" {
					hasID = true
					idVal = a.Val
				}
			}

			if !hasID {
				nodes = append(nodes, n)
			} else if idVal != "" && idVal != na.Hash.Name {
				nodes = append(nodes, n)
			}

			hasID = false
			idVal = ""
		}
	case 4:
		hasClass := false
		classVal := ""

		for _, n := range ctx.CNode {
			for _, a := range n.Attr {
				if a.Key == "class" {
					hasClass = true
					f := strings.Fields(a.Val)

					for _, s := range f {
						if s == na.Class.Name {
							classVal = s
						}
					}
				}
			}

			if !hasClass {
				nodes = append(nodes, n)
			} else if classVal != "" && classVal != na.Class.Name {
				nodes = append(nodes, n)
			}

			hasClass = false
			classVal = ""
		}
	case 5:
	case 6:
	}

	ctx.CNode = nodes
	return nodes
}

func evalPseudo(expr ast.Expression, ctx *object.Context) []*html.Node {
	return nil
}
