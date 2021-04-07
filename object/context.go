package object

import (
	"golang.org/x/net/html"
)

type Context struct {
	Doc   *html.Node
	CNode []*html.Node
}

func NewContext() *Context {
	return &Context{}
}
