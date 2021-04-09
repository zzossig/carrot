package eval

import (
	"github.com/zzossig/carrot/ast"
	"golang.org/x/net/html"
)

func evalPIdent(ident *ast.Ident, ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node
	switch ident.Value {
	case "first-child":
		return fnFirstChild(ctx, isNeg)
	case "last-child":
		return fnLastChild(ctx, isNeg)
	case "first-of-type":
		return fnFirstOfType(ctx, isNeg)
	case "last-of-type":
		return fnLastOfType(ctx, isNeg)
	case "only-child":
		return fnOnlyChild(ctx, isNeg)
	case "only-of-type":
		return fnOnlyOfType(ctx, isNeg)
	case "empty":
		return fnEmpty(ctx, isNeg)
	case "root":
		return fnRoot(ctx, isNeg)
	}
	return nodes
}

func evalPFP(fp *ast.FunctionalPseudo, ctx *Context, isNeg bool) []*html.Node {
	switch fp.Token.Literal {
	case "nth-child":
		return nthChild(fp.Arg, ctx, isNeg)
	case "nth-last-child":
		return nthLastChild(fp.Arg, ctx, isNeg)
	case "nth-of-type":
		return nthOfType(fp.Arg, ctx, isNeg)
	case "nth-last-of-type":
		return nthLastOfType(fp.Arg, ctx, isNeg)
	}
	return nil
}

func fnFirstChild(ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node

	if isNeg {
		for _, n := range ctx.CNode {
			if !isFirstChild(n) {
				nodes = appendNode(nodes, n)
			}
		}
	} else {
		for _, n := range ctx.CNode {
			if isFirstChild(n) {
				nodes = appendNode(nodes, n)
			}
		}
	}

	return nodes
}

func fnLastChild(ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node

	if isNeg {
		for _, n := range ctx.CNode {
			if !isLastChild(n) {
				nodes = appendNode(nodes, n)
			}
		}
	} else {
		for _, n := range ctx.CNode {
			if isLastChild(n) {
				nodes = appendNode(nodes, n)
			}
		}
	}

	return nodes
}

func fnFirstOfType(ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node

	if isNeg {
		for _, n := range ctx.CNode {
			if !isFirstOfType(n, ctx.CType) {
				nodes = appendNode(nodes, n)
			}
		}
	} else {
		for _, n := range ctx.CNode {
			if isFirstOfType(n, ctx.CType) {
				nodes = appendNode(nodes, n)
			}
		}
	}

	return nodes
}

func fnLastOfType(ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node

	if isNeg {
		for _, n := range ctx.CNode {
			if !isLastOfType(n, ctx.CType) {
				nodes = appendNode(nodes, n)
			}
		}
	} else {
		for _, n := range ctx.CNode {
			if isLastOfType(n, ctx.CType) {
				nodes = appendNode(nodes, n)
			}
		}
	}

	return nodes
}

func fnOnlyChild(ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node

	if isNeg {
		for _, n := range ctx.CNode {
			if !isOnlyChlid(n) {
				nodes = appendNode(nodes, n)
			}
		}
	} else {
		for _, n := range ctx.CNode {
			if isOnlyChlid(n) {
				nodes = appendNode(nodes, n)
			}
		}
	}

	return nodes
}

func fnOnlyOfType(ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node

	if isNeg {
		for _, n := range ctx.CNode {
			if !isOnlyOfType(n, ctx.CType) {
				nodes = appendNode(nodes, n)
			}
		}
	} else {
		for _, n := range ctx.CNode {
			if isOnlyOfType(n, ctx.CType) {
				nodes = appendNode(nodes, n)
			}
		}
	}

	return nodes
}

func fnEmpty(ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node

	if isNeg {
		for _, n := range ctx.CNode {
			if n.FirstChild != nil {
				nodes = appendNode(nodes, n)
			}
		}
	} else {
		for _, n := range ctx.CNode {
			if n.FirstChild == nil {
				nodes = appendNode(nodes, n)
			}
		}
	}

	return nodes
}

func fnRoot(ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node
	if ctx.Doc.FirstChild != nil && ctx.CType == ctx.Doc.FirstChild.Data {
		if isNeg {
			ctx.CNode = []*html.Node{ctx.Doc.FirstChild}
			nodes = collectDesc(ctx)
		} else {
			nodes = appendNode(nodes, ctx.Doc.FirstChild)
		}
	}
	return nodes
}

func nthChild(arg *ast.Arg, ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node

	switch arg.TypeID {
	case 1:
		if isNeg {
			for _, n := range ctx.CNode {
				if !isNthChild(n, arg.Dimension) {
					nodes = appendNode(nodes, n)
				}
			}
		} else {
			for _, n := range ctx.CNode {
				if isNthChild(n, arg.Dimension) {
					nodes = appendNode(nodes, n)
				}
			}
		}
	case 2:
		if isNeg {
			for _, n := range ctx.CNode {
				if !isNChild(n, arg.Number.Value) {
					nodes = appendNode(nodes, n)
				}
			}
		} else {
			for _, n := range ctx.CNode {
				if isNChild(n, arg.Number.Value) {
					nodes = appendNode(nodes, n)
				}
			}
		}
	case 4:
		switch arg.Ident.Value {
		case "even":
			if isNeg {
				for _, n := range ctx.CNode {
					if !isEvenChild(n) {
						nodes = appendNode(nodes, n)
					}
				}
			} else {
				for _, n := range ctx.CNode {
					if isEvenChild(n) {
						nodes = appendNode(nodes, n)
					}
				}
			}
		case "odd":
			if isNeg {
				for _, n := range ctx.CNode {
					if !isOddChild(n) {
						nodes = appendNode(nodes, n)
					}
				}
			} else {
				for _, n := range ctx.CNode {
					if isOddChild(n) {
						nodes = appendNode(nodes, n)
					}
				}
			}
		}
	}

	return nodes
}

func nthLastChild(arg *ast.Arg, ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node

	switch arg.TypeID {
	case 1:
		if isNeg {
			for _, n := range ctx.CNode {
				if !isNthLastChild(n, arg.Dimension) {
					nodes = appendNode(nodes, n)
				}
			}
		} else {
			for _, n := range ctx.CNode {
				if isNthLastChild(n, arg.Dimension) {
					nodes = appendNode(nodes, n)
				}
			}
		}
	case 2:
		if isNeg {
			for _, n := range ctx.CNode {
				if !isNLastChild(n, arg.Number.Value) {
					nodes = appendNode(nodes, n)
				}
			}
		} else {
			for _, n := range ctx.CNode {
				if isNLastChild(n, arg.Number.Value) {
					nodes = appendNode(nodes, n)
				}
			}
		}
	case 4:
		switch arg.Ident.Value {
		case "even":
			if isNeg {
				for _, n := range ctx.CNode {
					if !isEvenLastChild(n) {
						nodes = appendNode(nodes, n)
					}
				}
			} else {
				for _, n := range ctx.CNode {
					if isEvenLastChild(n) {
						nodes = appendNode(nodes, n)
					}
				}
			}
		case "odd":
			if isNeg {
				for _, n := range ctx.CNode {
					if !isOddLastChild(n) {
						nodes = appendNode(nodes, n)
					}
				}
			} else {
				for _, n := range ctx.CNode {
					if isOddLastChild(n) {
						nodes = appendNode(nodes, n)
					}
				}
			}
		}
	}

	return nodes
}

func nthOfType(arg *ast.Arg, ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node

	switch arg.TypeID {
	case 1:
		if isNeg {
			for _, n := range ctx.CNode {
				if !isNthOfType(n, arg.Dimension, ctx.CType) {
					nodes = appendNode(nodes, n)
				}
			}
		} else {
			for _, n := range ctx.CNode {
				if isNthOfType(n, arg.Dimension, ctx.CType) {
					nodes = appendNode(nodes, n)
				}
			}
		}
	case 2:
		if isNeg {
			for _, n := range ctx.CNode {
				if !isNOfType(n, arg.Number.Value, ctx.CType) {
					nodes = appendNode(nodes, n)
				}
			}
		} else {
			for _, n := range ctx.CNode {
				if isNOfType(n, arg.Number.Value, ctx.CType) {
					nodes = appendNode(nodes, n)
				}
			}
		}
	case 4:
		switch arg.Ident.Value {
		case "even":
			if isNeg {
				for _, n := range ctx.CNode {
					if !isEvenNthOfType(n, ctx.CType) {
						nodes = appendNode(nodes, n)
					}
				}
			} else {
				for _, n := range ctx.CNode {
					if isEvenNthOfType(n, ctx.CType) {
						nodes = appendNode(nodes, n)
					}
				}
			}
		case "odd":
			if isNeg {
				for _, n := range ctx.CNode {
					if !isOddNthOfType(n, ctx.CType) {
						nodes = appendNode(nodes, n)
					}
				}
			} else {
				for _, n := range ctx.CNode {
					if isOddNthOfType(n, ctx.CType) {
						nodes = appendNode(nodes, n)
					}
				}
			}
		}
	}

	return nodes
}

func nthLastOfType(arg *ast.Arg, ctx *Context, isNeg bool) []*html.Node {
	var nodes []*html.Node

	switch arg.TypeID {
	case 1:
		if isNeg {
			for _, n := range ctx.CNode {
				if !isNthLastOfType(n, arg.Dimension, ctx.CType) {
					nodes = appendNode(nodes, n)
				}
			}
		} else {
			for _, n := range ctx.CNode {
				if isNthLastOfType(n, arg.Dimension, ctx.CType) {
					nodes = appendNode(nodes, n)
				}
			}
		}
	case 2:
		if isNeg {
			for _, n := range ctx.CNode {
				if !isNLastOfType(n, arg.Number.Value, ctx.CType) {
					nodes = appendNode(nodes, n)
				}
			}
		} else {
			for _, n := range ctx.CNode {
				if isNLastOfType(n, arg.Number.Value, ctx.CType) {
					nodes = appendNode(nodes, n)
				}
			}
		}
	case 4:
		switch arg.Ident.Value {
		case "even":
			if isNeg {
				for _, n := range ctx.CNode {
					if !isEvenNthLastOfType(n, ctx.CType) {
						nodes = appendNode(nodes, n)
					}
				}
			} else {
				for _, n := range ctx.CNode {
					if isEvenNthLastOfType(n, ctx.CType) {
						nodes = appendNode(nodes, n)
					}
				}
			}
		case "odd":
			if isNeg {
				for _, n := range ctx.CNode {
					if !isOddNthLastOfType(n, ctx.CType) {
						nodes = appendNode(nodes, n)
					}
				}
			} else {
				for _, n := range ctx.CNode {
					if isOddNthLastOfType(n, ctx.CType) {
						nodes = appendNode(nodes, n)
					}
				}
			}
		}
	}

	return nodes
}
