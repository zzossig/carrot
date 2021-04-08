package eval

import (
	"testing"

	"github.com/zzossig/carrot/lexer"
	"github.com/zzossig/carrot/parser"
	"golang.org/x/net/html"
)

func TestCarrot(t *testing.T) {
	e1 := testEval("h1")
	if len(e1) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e1))
	}
}

func testEval(input string) []*html.Node {
	l := lexer.New(input)
	p := parser.New(l)
	e := p.ParseExpression()
	ctx := NewContext()
	ctx.SetDoc("./testdata/t.html")

	return Eval(e, ctx)
}
