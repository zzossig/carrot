package eval

import (
	"github.com/zzossig/carrot/ast"
	"github.com/zzossig/carrot/object"
	"golang.org/x/net/html"
)

func appendNode(src []*html.Node, target *html.Node) []*html.Node {
	if !isContain(src, target) {
		src = append(src, target)
	}
	return src
}

func isContain(src []*html.Node, target *html.Node) bool {
	for _, item := range src {
		if item == target {
			return true
		}
	}
	return false
}

func collectSubSibling(ctx *object.Context) []*html.Node {
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		for s := n.NextSibling; s != nil; s = s.NextSibling {
			if s.Type == html.ElementNode {
				nodes = appendNode(nodes, n)
			}
		}
	}

	return nodes
}

func collectNextSibling(ctx *object.Context) []*html.Node {
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		for s := n.NextSibling; s != nil; s = s.NextSibling {
			if s.Type == html.ElementNode {
				nodes = appendNode(nodes, n)
				break
			}
		}
	}

	return nodes
}

func collectChild(ctx *object.Context) []*html.Node {
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				nodes = appendNode(nodes, n)
			}
		}
	}

	return nodes
}

func collectDesc(ctx *object.Context) []*html.Node {
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		for _, d := range walkDesc(n) {
			nodes = appendNode(nodes, d)
		}
	}

	return nodes
}

func collectDescOrSelf(ctx *object.Context) []*html.Node {
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		nodes = appendNode(nodes, n)
		for _, d := range walkDesc(n) {
			nodes = appendNode(nodes, d)
		}
	}

	return nodes
}

func walkDesc(n *html.Node) []*html.Node {
	var nodes []*html.Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			nodes = appendNode(nodes, n)

			if c.FirstChild != nil {
				for _, d := range walkDesc(c) {
					nodes = appendNode(nodes, d)
				}
			}
		}
	}

	return nodes
}

func isFirstChild(n *html.Node) bool {
	if n.Parent == nil {
		return false
	}

	for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if c == n {
				return true
			}
			break
		}
	}

	return false
}

func isLastChild(n *html.Node) bool {
	if n.Parent == nil {
		return false
	}

	for c := n.Parent.LastChild; c != nil; c = c.PrevSibling {
		if c.Type == html.ElementNode {
			if c == n {
				return true
			}
			break
		}
	}

	return false
}

func isOnlyChlid(n *html.Node) bool {
	if n.Parent == nil {
		return false
	}

	cnt := 0
	for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			cnt++
		}
	}

	return cnt == 1
}

func isNChild(n *html.Node, num int) bool {
	if n.Parent == nil {
		return false
	}

	i := 0
	for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			i++
			if i == num {
				return true
			}
		}
	}

	return false
}

func isNLastChild(n *html.Node, num int) bool {
	if n.Parent == nil {
		return false
	}

	i := 0
	for c := n.Parent.LastChild; c != nil; c = c.PrevSibling {
		if c.Type == html.ElementNode {
			i++
			if i == num {
				return true
			}
		}
	}

	return false
}

func isNOfType(n *html.Node, num int, t string) bool {
	if n.Parent == nil {
		return false
	}

	i := 0
	for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			i++
			if i == num && c.Data == t {
				return true
			}
		}
	}

	return false
}

func isNLastOfType(n *html.Node, num int, t string) bool {
	if n.Parent == nil {
		return false
	}

	i := 0
	for c := n.Parent.LastChild; c != nil; c = c.PrevSibling {
		if c.Type == html.ElementNode {
			i++
			if i == num && c.Data == t {
				return true
			}
		}
	}

	return false
}

func isNthChild(n *html.Node, d *ast.Dimension) bool {
	if n.Parent == nil {
		return false
	}
	if d.A == 0 {
		return isNChild(n, d.B)
	}

	var remaining = func(dd *ast.Dimension, i int) int {
		a := dd.A
		b := dd.B

		if dd.Aop == "-" {
			a = -a
		}
		if dd.Bop == "-" {
			b = -b
		}

		return (i - b) % a
	}

	i := 0
	for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			i++
			if remaining(d, i) == 0 && i >= d.B {
				return true
			}
		}
	}

	return false
}

func isNthLastChild(n *html.Node, d *ast.Dimension) bool {
	if n.Parent == nil {
		return false
	}
	if d.A == 0 {
		return isNLastChild(n, d.B)
	}

	var remaining = func(dd *ast.Dimension, i int) int {
		a := dd.A
		b := dd.B

		if dd.Aop == "-" {
			a = -a
		}
		if dd.Bop == "-" {
			b = -b
		}

		return (i - b) % a
	}

	i := 0
	for c := n.Parent.LastChild; c != nil; c = c.PrevSibling {
		if c.Type == html.ElementNode {
			i++
			if remaining(d, i) == 0 && i >= d.B {
				return true
			}
		}
	}

	return false
}

func isNthOfType(n *html.Node, d *ast.Dimension, t string) bool {
	if n.Parent == nil {
		return false
	}
	if d.A == 0 {
		return isNOfType(n, d.B, t)
	}

	var remaining = func(dd *ast.Dimension, i int) int {
		a := dd.A
		b := dd.B

		if dd.Aop == "-" {
			a = -a
		}
		if dd.Bop == "-" {
			b = -b
		}

		return (i - b) % a
	}

	i := 0
	for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			i++
			if remaining(d, i) == 0 && i >= d.B && c.Data == t {
				return true
			}
		}
	}

	return false
}

func isNthLastOfType(n *html.Node, d *ast.Dimension, t string) bool {
	if n.Parent == nil {
		return false
	}
	if d.A == 0 {
		return isNLastOfType(n, d.B, t)
	}

	var remaining = func(dd *ast.Dimension, i int) int {
		a := dd.A
		b := dd.B

		if dd.Aop == "-" {
			a = -a
		}
		if dd.Bop == "-" {
			b = -b
		}

		return (i - b) % a
	}

	i := 0
	for c := n.Parent.LastChild; c != nil; c = c.PrevSibling {
		if c.Type == html.ElementNode {
			i++
			if remaining(d, i) == 0 && i >= d.B && c.Data == t {
				return true
			}
		}
	}

	return false
}

func isEvenChild(n *html.Node) bool {
	d := &ast.Dimension{A: 2, Aop: "+"}
	return isNthChild(n, d)
}

func isEvenLastChild(n *html.Node) bool {
	d := &ast.Dimension{A: 2, Aop: "+"}
	return isNthLastChild(n, d)
}

func isOddChild(n *html.Node) bool {
	d := &ast.Dimension{A: 2, Aop: "+", B: 1, Bop: "+"}
	return isNthChild(n, d)
}

func isOddLastChild(n *html.Node) bool {
	d := &ast.Dimension{A: 2, Aop: "+", B: 1, Bop: "+"}
	return isNthLastChild(n, d)
}

func isEvenNthOfType(n *html.Node, t string) bool {
	d := &ast.Dimension{A: 2, Aop: "+"}
	return isNthOfType(n, d, t)
}

func isOddNthOfType(n *html.Node, t string) bool {
	d := &ast.Dimension{A: 2, Aop: "+", B: 1, Bop: "+"}
	return isNthOfType(n, d, t)
}

func isEvenNthLastOfType(n *html.Node, t string) bool {
	d := &ast.Dimension{A: 2, Aop: "+"}
	return isNthLastOfType(n, d, t)
}

func isOddNthLastOfType(n *html.Node, t string) bool {
	d := &ast.Dimension{A: 2, Aop: "+", B: 1, Bop: "+"}
	return isNthLastOfType(n, d, t)
}
