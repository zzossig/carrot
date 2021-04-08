package eval

import (
	"bufio"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

type Context struct {
	Doc   *html.Node
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

	InitCNode(c)

	return nil
}

func InitCNode(ctx *Context) {
	if ctx.Doc != nil {
		ctx.CNode = walkDesc(ctx.Doc)
	}
}
