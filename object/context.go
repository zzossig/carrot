package object

import (
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
