package eval

import (
	"testing"

	"github.com/zzossig/carrot/lexer"
	"github.com/zzossig/carrot/parser"
	"golang.org/x/net/html"
)

func TestCarrot(t *testing.T) {
	// e1 := testEval("h1")
	// if len(e1) != 1 {
	// 	t.Errorf("wrong number of items. got=%d, expected=1", len(e1))
	// }
	// if e1[0].Data != "h1" {
	// 	t.Errorf("selected node should be h1")
	// }

	// e2 := testEval("h4")
	// if len(e2) != 1 {
	// 	t.Errorf("wrong number of items. got=%d, expected=1", len(e2))
	// }
	// if e2[0].Data != "h4" {
	// 	t.Errorf("selected node should be h4")
	// }

	// e3 := testEval("h1,h2,h3")
	// if len(e3) != 3 {
	// 	t.Errorf("wrong number of items. got=%d, expected=3", len(e3))
	// }

	// e4 := testEval("*.depth")
	// if len(e4) != 4 {
	// 	t.Errorf("wrong number of items. got=%d, expected=4", len(e4))
	// }

	// e5 := testEval(".depth")
	// if len(e5) != 4 {
	// 	t.Errorf("wrong number of items. got=%d, expected=4", len(e5))
	// }

	// e6 := testEval("*#1")
	// if len(e6) != 0 {
	// 	t.Errorf("wrong number of items. got=%d, expected=0", len(e6))
	// }

	// e7 := testEval("*#a")
	// if len(e7) != 1 {
	// 	t.Errorf("wrong number of items. got=%d, expected=1", len(e7))
	// }
	// if e7[0].Attr[0].Key != "id" || e7[0].Attr[0].Val != "a" {
	// 	t.Errorf("wrong node selected")
	// }

	// e8 := testEval("p:first-child")
	// if len(e8) != 2 {
	// 	t.Errorf("wrong number of items. got=%d, expected=2", len(e8))
	// }

	// e9 := testEval("h1[title]")
	// if len(e9) != 1 {
	// 	t.Errorf("wrong number of items. got=%d, expected=1", len(e9))
	// }
	// if e9[0].Attr[0].Key != "title" || e9[0].Attr[0].Val != "carrot" {
	// 	t.Errorf("wrong node selected")
	// }

	// e10 := testEval("span[class='example']")
	// if len(e10) != 1 {
	// 	t.Errorf("wrong number of items. got=%d, expected=1", len(e10))
	// }
	// if e10[0].Attr[0].Key != "class" || e10[0].Attr[0].Val != "example" {
	// 	t.Errorf("wrong node selected")
	// }

	// e11 := testEval("span[hello='Cleveland']")
	// if len(e11) != 2 {
	// 	t.Errorf("wrong number of items. got=%d, expected=2", len(e11))
	// }

	// e12 := testEval("span[hello='Cleveland'][goodbye='Columbus']")
	// if len(e12) != 1 {
	// 	t.Errorf("wrong number of items. got=%d, expected=1", len(e12))
	// }
	// if e12[0].Attr[1].Key != "goodbye" || e12[0].Attr[1].Val != "Columbus" {
	// 	t.Errorf("wrong node selected")
	// }

	// e13 := testEval("div[id='e'] ~ p")
	// if len(e13) != 1 {
	// 	t.Errorf("wrong number of items. got=%d, expected=1", len(e13))
	// }

	// e14 := testEval("p[class='foo']")
	// if len(e14) != 2 {
	// 	t.Errorf("wrong number of items. got=%d, expected=2", len(e14))
	// }

	// e15 := testEval("p[class~='foo']")
	// if len(e15) != 2 {
	// 	t.Errorf("wrong number of items. got=%d, expected=2", len(e15))
	// }

	// e16 := testEval("p[class|='foo']")
	// if len(e16) != 4 {
	// 	t.Errorf("wrong number of items. got=%d, expected=4", len(e16))
	// }

	// e17 := testEval("p[class^=foo]")
	// if len(e17) != 6 {
	// 	t.Errorf("wrong number of items. got=%d, expected=6", len(e17))
	// }

	// e18 := testEval("p[class$=foo]")
	// if len(e18) != 3 {
	// 	t.Errorf("wrong number of items. got=%d, expected=3", len(e18))
	// }

	e19 := testEval("p[class*=foo]")
	if len(e19) != 7 {
		t.Errorf("wrong number of items. got=%d, expected=7", len(e19))
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
