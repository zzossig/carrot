package eval

import (
	"strings"

	"github.com/zzossig/carrot/ast"
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

func walkDesc(n *html.Node) []*html.Node {
	var nodes []*html.Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			nodes = appendNode(nodes, c)
			if c.FirstChild != nil {
				for _, d := range walkDesc(c) {
					nodes = appendNode(nodes, d)
				}
			}
		}
	}

	return nodes
}

func collectSubSibling(ctx *Context) []*html.Node {
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		for s := n.NextSibling; s != nil; s = s.NextSibling {
			if s.Type == html.ElementNode {
				nodes = append(nodes, s)
			}
		}
	}

	return nodes
}

func collectNextSibling(ctx *Context) []*html.Node {
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		for s := n.NextSibling; s != nil; s = s.NextSibling {
			if s.Type == html.ElementNode {
				nodes = append(nodes, s)
				break
			}
		}
	}

	return nodes
}

func collectChild(ctx *Context) []*html.Node {
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				nodes = append(nodes, c)
			}
		}
	}

	return nodes
}

func collectDesc(ctx *Context) []*html.Node {
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		nodes = append(nodes, walkDesc(n)...)
	}

	return nodes
}

func collectDescOrSelf(ctx *Context) []*html.Node {
	var nodes []*html.Node

	for _, n := range ctx.CNode {
		nodes = appendNode(nodes, n)
		for _, d := range walkDesc(n) {
			nodes = appendNode(nodes, d)
		}
	}

	return nodes
}

func isDashMatched(s, substr string) bool {
	if s == substr {
		return true
	}
	if strings.HasPrefix(s, substr+"-") {
		return true
	}
	return false
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

func isFirstOfType(n *html.Node, t string) bool {
	if n.Parent == nil {
		return false
	}

	for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == t {
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

func isLastOfType(n *html.Node, t string) bool {
	if n.Parent == nil {
		return false
	}

	for c := n.Parent.LastChild; c != nil; c = c.PrevSibling {
		if c.Type == html.ElementNode && c.Data == t {
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

func isOnlyOfType(n *html.Node, t string) bool {
	if n.Parent == nil {
		return false
	}

	cnt := 0
	for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == t {
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
			if i == num && n == c {
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
			if i == num && n == c {
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
		if c.Type == html.ElementNode && c.Data == t {
			i++
			if i == num && n == c {
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
		if c.Type == html.ElementNode && c.Data == t {
			i++
			if i == num && n == c {
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

	a := d.A
	b := d.B

	if d.Aop == "-" {
		a = -a
	}
	if d.Bop == "-" {
		b = -b
	}

	if d.A == 0 {
		return isNChild(n, b)
	}

	i := 0
	for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			i++
			if (i-b)%a == 0 && i >= b && n == c {
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

	a := d.A
	b := d.B

	if d.Aop == "-" {
		a = -a
	}
	if d.Bop == "-" {
		b = -b
	}

	if d.A == 0 {
		return isNLastChild(n, b)
	}

	i := 0
	for c := n.Parent.LastChild; c != nil; c = c.PrevSibling {
		if c.Type == html.ElementNode {
			i++
			if a < 0 {
				if (i-b)%a == 0 && i <= b && n == c {
					return true
				}
			} else {
				if (i-b)%a == 0 && i >= b && n == c {
					return true
				}
			}
		}
	}

	return false
}

func isNthOfType(n *html.Node, d *ast.Dimension, t string) bool {
	if n.Parent == nil {
		return false
	}

	a := d.A
	b := d.B

	if d.Aop == "-" {
		a = -a
	}
	if d.Bop == "-" {
		b = -b
	}

	if d.A == 0 || (a < 0 && d.A > 1) {
		return isNOfType(n, b, t)
	}

	i := 0
	for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == t {
			i++
			if a < 0 {
				if a*(i-1)+b > 0 && n == c {
					return true
				}
			} else {
				if (i-b)%a == 0 && i >= b && n == c {
					return true
				}
			}
		}
	}

	return false
}

func isNthLastOfType(n *html.Node, d *ast.Dimension, t string) bool {
	if n.Parent == nil {
		return false
	}

	a := d.A
	b := d.B

	if d.Aop == "-" {
		a = -a
	}
	if d.Bop == "-" {
		b = -b
	}

	if d.A == 0 || (a < 0 && d.A > 1) {
		return isNLastOfType(n, b, t)
	}

	i := 0
	for c := n.Parent.LastChild; c != nil; c = c.PrevSibling {
		if c.Type == html.ElementNode && c.Data == t {
			i++
			if a < 0 {
				if a*(i-1)+b > 0 && n == c {
					return true
				}
			} else {
				if (i-b)%a == 0 && i >= b && n == c {
					return true
				}
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
