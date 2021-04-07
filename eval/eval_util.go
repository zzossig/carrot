package eval

import (
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
	return nodes
}

func collectNextSibling(ctx *object.Context) []*html.Node {
	var nodes []*html.Node
	return nodes
}

func collectChild(ctx *object.Context) []*html.Node {
	var nodes []*html.Node
	return nodes
}

func collectDesc(ctx *object.Context) []*html.Node {
	var nodes []*html.Node
	return nodes
}
