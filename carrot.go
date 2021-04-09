package carrot

import (
	"net/http"

	"github.com/zzossig/carrot/eval"
	"github.com/zzossig/carrot/lexer"
	"github.com/zzossig/carrot/parser"
	"golang.org/x/net/html"
)

// CSS is a base object to evaluate css selectors.
type CSS struct {
	selector string
	context  *eval.Context
	errors   []error
}

// New creates new CSS object.
func New() *CSS {
	return &CSS{context: eval.NewContext()}
}

// SetDoc set document to a context.
// input param can be url or local filepath.
func (c *CSS) SetDoc(input string) *CSS {
	c.context.GetBackCtx()
	c.selector = ""

	err := c.context.SetDoc(input)
	if err != nil {
		c.errors = append(c.errors, err)
	}

	return c
}

// SetDocR is another version of SetDoc.
func (c *CSS) SetDocR(r *http.Response) *CSS {
	c.context.GetBackCtx()
	c.selector = ""

	err := c.context.SetDocR(r)
	if err != nil {
		c.errors = append(c.errors, err)
	}

	return c
}

// SetDocN is another version of SetDoc.
func (c *CSS) SetDocN(n *html.Node) *CSS {
	c.context.GetBackCtx()
	c.selector = ""
	c.context.SetDocN(n)
	return c
}

// SetDocS is another version of SetDoc.
func (c *CSS) SetDocS(s string) *CSS {
	c.context.GetBackCtx()
	c.selector = ""

	err := c.context.SetDocS(s)
	if err != nil {
		c.errors = append(c.errors, err)
	}

	return c
}

// Eval evaluates a css selector
func (c *CSS) Eval(input string) []*html.Node {
	if len(c.errors) > 0 {
		return nil
	}

	c.selector = input

	l := lexer.New(input)
	p := parser.New(l)
	pe := p.ParseExpression()

	if len(p.Errors()) != 0 {
		c.errors = append(c.errors, p.Errors()...)
		return nil
	}

	e := eval.Eval(pe, c.context)
	c.context.GetBackCtx()
	return e
}

// Errors returns errors field
func (c *CSS) Errors() []error {
	return c.errors
}

// String returns input field
func (c *CSS) String() string {
	return c.selector
}
