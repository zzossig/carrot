package eval

import (
	"github.com/zzossig/carrot/ast"
	"golang.org/x/net/html"
)

// Eval function evaluate CSS selector
func Eval(expr ast.Expression, ctx *Context) []*html.Node {
	switch expr := expr.(type) {
	case *ast.Group:
		return evalGroup(expr, ctx)
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
	case *ast.Ident:
		return evalIdent(expr, ctx)
	case *ast.Negation:
		return evalNegation(expr, ctx)
	case *ast.Has:
		return evalHas(expr, ctx)
	case *ast.Attrib:
		return evalAttrib(expr, ctx, false)
	case *ast.Pseudo:
		return evalPseudo(expr, ctx, false)
	}
	return nil
}
