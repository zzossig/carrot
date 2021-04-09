package eval

import (
	"bufio"
	"net/http"
	"os"
	"strings"

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

	if c.Doc != nil {
		c.CNode = walkDesc(c.Doc)
		c.Nodes = make([]*html.Node, len(c.CNode))
		copy(c.Nodes, c.CNode)
	}

	return nil
}

func (c *Context) SetDocR(r *http.Response) error {
	defer r.Body.Close()

	nr := bufio.NewReader(r.Body)
	parsedHTML, err := html.Parse(nr)
	if err != nil {
		return err
	}

	c.Doc = parsedHTML
	if c.Doc != nil {
		c.CNode = walkDesc(c.Doc)
		c.Nodes = make([]*html.Node, len(c.CNode))
		copy(c.Nodes, c.CNode)
	}

	return nil
}

func (c *Context) SetDocN(n *html.Node) {
	c.Doc = n
	if c.Doc != nil {
		c.CNode = walkDesc(c.Doc)
		c.Nodes = make([]*html.Node, len(c.CNode))
		copy(c.Nodes, c.CNode)
	}
}

func (c *Context) SetDocS(s string) error {
	nr := strings.NewReader(s)
	parsedHTML, err := html.Parse(nr)
	if err != nil {
		return err
	}

	c.Doc = parsedHTML
	if c.Doc != nil {
		c.CNode = walkDesc(c.Doc)
		c.Nodes = make([]*html.Node, len(c.CNode))
		copy(c.Nodes, c.CNode)
	}

	return nil
}

func (c *Context) GetBackCtx() {
	c.CType = ""
	c.CNode = make([]*html.Node, len(c.Nodes))
	copy(c.CNode, c.Nodes)
}
