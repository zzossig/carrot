package eval

import (
	"strings"

	"github.com/zzossig/carrot/ast"
	"github.com/zzossig/carrot/token"
	"golang.org/x/net/html"
)

func evalGroup(expr ast.Expression, ctx *Context) []*html.Node {
	g := expr.(*ast.Group)
	var nodes []*html.Node

	for _, selector := range g.Selectors {
		e := Eval(selector, ctx)
		for _, ee := range e {
			nodes = appendNode(nodes, ee)
		}
		ctx.GetBackCtx()
	}

	return nodes
}

func evalSelector(expr ast.Expression, ctx *Context) []*html.Node {
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

func evalSequence(expr ast.Expression, ctx *Context) []*html.Node {
	s := expr.(*ast.Sequence)

	if ident, ok := s.Expression.(*ast.Ident); ok {
		ctx.CType = ident.Value
	}

	if s.Expression != nil {
		h := Eval(s.Expression, ctx)
		ctx.CNode = h
	}

	for _, e := range s.Exprs {
		ss := Eval(e, ctx)
		ctx.CNode = ss
	}

	ctx.CType = ""
	return ctx.CNode
}

func evalUniversal(expr ast.Expression, ctx *Context) []*html.Node {
	return collectDescOrSelf(ctx)
}

func evalClass(expr ast.Expression, ctx *Context) []*html.Node {
	c := expr.(*ast.Class)
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		for _, a := range n.Attr {
			if a.Key == "class" {
				f := strings.Fields(a.Val)
				for _, s := range f {
					if s == c.Name {
						nodes = appendNode(nodes, n)
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	return nodes
}

func evalHash(expr ast.Expression, ctx *Context) []*html.Node {
	h := expr.(*ast.Hash)
	var nodes []*html.Node

	if h.Name == "" {
		return nodes
	} else if strings.ContainsRune("0123456789", rune(h.Name[0])) {
		return nodes
	}

	for _, n := range ctx.CNode {
		for _, a := range n.Attr {
			if a.Key == "id" && a.Val == h.Name {
				nodes = appendNode(nodes, n)
			}
		}
	}

	ctx.CNode = nodes
	return nodes
}

func evalAttrib(expr ast.Expression, ctx *Context, isNeg bool) []*html.Node {
	ae := expr.(*ast.Attrib).AttrExpr
	var nodes []*html.Node

	switch ae.TypeID {
	case 1:
		if isNeg {
			has := false
			for _, n := range ctx.CNode {
				for _, a := range n.Attr {
					if a.Key == ae.Left.Value {
						has = true
					}
				}

				if !has {
					nodes = appendNode(nodes, n)
				}
				has = false
			}
		} else {
			for _, n := range ctx.CNode {
				for _, a := range n.Attr {
					if a.Key == ae.Left.Value {
						nodes = appendNode(nodes, n)
					}
				}
			}
		}
	case 2:
		switch ae.Token.Type {
		case token.EQ:
			if isNeg {
				has := false
				for _, n := range ctx.CNode {
					for _, a := range n.Attr {
						if a.Key == ae.Left.Value && a.Val == ae.Right.Value {
							has = true
						}
					}

					if !has {
						nodes = appendNode(nodes, n)
					}
					has = false
				}
			} else {
				for _, n := range ctx.CNode {
					for _, a := range n.Attr {
						if a.Key == ae.Left.Value && a.Val == ae.Right.Value {
							nodes = appendNode(nodes, n)
						}
					}
				}
			}
		case token.PREFIXMATCH:
			if isNeg {
				has := false
				for _, n := range ctx.CNode {
					for _, a := range n.Attr {
						if a.Key == ae.Left.Value {
							if strings.HasPrefix(a.Val, ae.Right.Value) {
								has = true
							}
						}
					}

					if !has {
						nodes = appendNode(nodes, n)
					}
					has = false
				}
			} else {
				for _, n := range ctx.CNode {
					for _, a := range n.Attr {
						if a.Key == ae.Left.Value {
							if strings.HasPrefix(a.Val, ae.Right.Value) {
								nodes = appendNode(nodes, n)
							}
						}
					}
				}
			}
		case token.SUFFIXMATCH:
			if isNeg {
				has := false
				for _, n := range ctx.CNode {
					for _, a := range n.Attr {
						if a.Key == ae.Left.Value {
							if strings.HasSuffix(a.Val, ae.Right.Value) {
								has = true
							}
						}
					}

					if !has {
						nodes = appendNode(nodes, n)
					}
					has = false
				}
			} else {
				for _, n := range ctx.CNode {
					for _, a := range n.Attr {
						if a.Key == ae.Left.Value {
							if strings.HasSuffix(a.Val, ae.Right.Value) {
								nodes = appendNode(nodes, n)
							}
						}
					}
				}
			}
		case token.SUBSTRINGMATCH:
			if isNeg {
				has := false
				for _, n := range ctx.CNode {
					for _, a := range n.Attr {
						if a.Key == ae.Left.Value {
							if strings.Contains(a.Val, ae.Right.Value) {
								has = true
							}
						}
					}

					if !has {
						nodes = appendNode(nodes, n)
					}
					has = false
				}
			} else {
				for _, n := range ctx.CNode {
					for _, a := range n.Attr {
						if a.Key == ae.Left.Value {
							if strings.Contains(a.Val, ae.Right.Value) {
								nodes = appendNode(nodes, n)
							}
						}
					}
				}
			}
		case token.INCLUDES:
			if isNeg {
				has := false
				for _, n := range ctx.CNode {
					for _, a := range n.Attr {
						if a.Key == ae.Left.Value {
							f := strings.Fields(a.Val)
							for _, s := range f {
								if s == ae.Right.Value {
									has = true
								}
							}
						}
					}

					if !has {
						nodes = appendNode(nodes, n)
					}
					has = false
				}
			} else {
				for _, n := range ctx.CNode {
					for _, a := range n.Attr {
						if a.Key == ae.Left.Value {
							f := strings.Fields(a.Val)
							for _, s := range f {
								if s == ae.Right.Value {
									nodes = appendNode(nodes, n)
								}
							}
						}
					}
				}
			}
		case token.DASHMATCH:
			if isNeg {
				has := false
				for _, n := range ctx.CNode {
					for _, a := range n.Attr {
						if a.Key == ae.Left.Value {
							if isDashMatched(a.Val, ae.Right.Value) {
								has = true
							}
						}
					}

					if !has {
						nodes = appendNode(nodes, n)
					}
					has = false
				}
			} else {
				for _, n := range ctx.CNode {
					for _, a := range n.Attr {
						if a.Key == ae.Left.Value {
							if isDashMatched(a.Val, ae.Right.Value) {
								nodes = appendNode(nodes, n)
							}
						}
					}
				}
			}
		}
	}

	ctx.CNode = nodes
	return nodes
}

func evalIdent(expr ast.Expression, ctx *Context) []*html.Node {
	i := expr.(*ast.Ident)
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		if n.Data == i.Value {
			nodes = appendNode(nodes, n)
		}
	}

	ctx.CNode = nodes
	return nodes
}

func evalNegation(expr ast.Expression, ctx *Context) []*html.Node {
	na := expr.(*ast.Negation).NArg
	var nodes []*html.Node

	switch na.TypeID {
	case 1:
		for _, n := range ctx.CNode {
			if n.Data != na.Ident.Value {
				nodes = appendNode(nodes, n)
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
				nodes = appendNode(nodes, n)
			} else if idVal != "" && idVal != na.Hash.Name {
				nodes = appendNode(nodes, n)
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
				nodes = appendNode(nodes, n)
			} else if classVal != "" && classVal != na.Class.Name {
				nodes = appendNode(nodes, n)
			}

			hasClass = false
			classVal = ""
		}
	case 5:
		nodes = evalAttrib(na.Attrib, ctx, true)
	case 6:
		nodes = evalPseudo(na.Pseudo, ctx, true)
	case 7:
		for _, selector := range na.Group.Selectors {
			negation := makeNegation(selector)
			if negation == nil {
				continue
			}
			nodes = evalNegation(negation, ctx)
			ctx.CNode = nodes
		}
	}

	ctx.CNode = nodes
	return nodes
}

func evalHas(expr ast.Expression, ctx *Context) []*html.Node {
	has := expr.(*ast.Has).HArg
	var nodes []*html.Node

	switch has.TypeID {
	case 1:
	case 2:
	case 3:
	case 4:
	case 5:
	case 6:
	case 7:
	case 8:
	}

	ctx.CNode = nodes
	return nodes
}

func evalPseudo(expr ast.Expression, ctx *Context, isNeg bool) []*html.Node {
	p := expr.(*ast.Pseudo)
	var nodes []*html.Node

	switch p.TypeID {
	case 1:
		nodes = evalPIdent(p.Ident, ctx, isNeg)
	case 2:
		nodes = evalPFP(p.FunctionalPseudo, ctx, isNeg)
	}

	ctx.CNode = nodes
	return nodes
}
