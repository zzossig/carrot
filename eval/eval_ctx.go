package eval

import (
	"bufio"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

type Context struct {
	Doc   *html.Node
	Nodes []*html.Node
	CNode []*html.Node
	CType string
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) SetDoc(input string) error {
	if file, err := os.Open(input); err == nil {
		defer file.Close()

		buf := bufio.NewReader(file)
		parsedHTML, err := html.Parse(buf)
		if err != nil {
			return err
		}

		c.Doc = parsedHTML
	}

	if resp, err := http.Get(input); err == nil {
		defer resp.Body.Close()

		buf := bufio.NewReader(resp.Body)
		parsedHTML, err := html.Parse(buf)
		if err != nil {
			return err
		}

		c.Doc = parsedHTML
	}

	setCNode(c)

	return nil
}

func (c *Context) InitContext() {
	c.CType = ""
	c.CNode = make([]*html.Node, len(c.Nodes))
	copy(c.CNode, c.Nodes)
}

func setCNode(ctx *Context) {
	if ctx.Doc != nil {
		ctx.CNode = walkDesc(ctx.Doc)
		ctx.Nodes = make([]*html.Node, len(ctx.CNode))
		copy(ctx.Nodes, ctx.CNode)
	}
}
