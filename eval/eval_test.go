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
	if e1[0].Data != "h1" {
		t.Errorf("selected node should be h1")
	}

	e2 := testEval("h4")
	if len(e2) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e2))
	}
	if e2[0].Data != "h4" {
		t.Errorf("selected node should be h4")
	}

	e3 := testEval("h1,h2,h3")
	if len(e3) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(e3))
	}

	e4 := testEval("*.depth")
	if len(e4) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e4))
	}

	e5 := testEval(".depth")
	if len(e5) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e5))
	}

	e6 := testEval("*#1")
	if len(e6) != 0 {
		t.Errorf("wrong number of items. got=%d, expected=0", len(e6))
	}

	e7 := testEval("*#a")
	if len(e7) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e7))
	}
	if e7[0].Attr[0].Key != "id" || e7[0].Attr[0].Val != "a" {
		t.Errorf("wrong node selected")
	}

	e8 := testEval("p:first-child")
	if len(e8) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e8))
	}

	e9 := testEval("h1[title]")
	if len(e9) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e9))
	}
	if e9[0].Attr[0].Key != "title" || e9[0].Attr[0].Val != "carrot" {
		t.Errorf("wrong node selected")
	}

	e10 := testEval("span[class='example']")
	if len(e10) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e10))
	}
	if e10[0].Attr[0].Key != "class" || e10[0].Attr[0].Val != "example" {
		t.Errorf("wrong node selected")
	}

	e11 := testEval("span[hello='Cleveland']")
	if len(e11) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(e11))
	}

	e12 := testEval("span[hello='Cleveland'][goodbye='Columbus']")
	if len(e12) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e12))
	}
	if e12[0].Attr[1].Key != "goodbye" || e12[0].Attr[1].Val != "Columbus" {
		t.Errorf("wrong node selected")
	}

	e13 := testEval("div#a")
	if len(e13) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e13))
	}
	if e13[0].Attr[0].Key != "id" || e13[0].Attr[0].Val != "a" {
		t.Errorf("wrong node selected")
	}

	e14 := testEval("p[class='foo']")
	if len(e14) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(e14))
	}

	e15 := testEval("p[class~='foo']")
	if len(e15) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(e15))
	}

	e16 := testEval("p[class|='foo']")
	if len(e16) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e16))
	}

	e17 := testEval("p[class^=foo]")
	if len(e17) != 6 {
		t.Errorf("wrong number of items. got=%d, expected=6", len(e17))
	}

	e18 := testEval("p[class$=foo]")
	if len(e18) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(e18))
	}

	e19 := testEval("p[class*=foo]")
	if len(e19) != 7 {
		t.Errorf("wrong number of items. got=%d, expected=7", len(e19))
	}

	e20 := testEval("*.p.q.r")
	if len(e20) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e20))
	}

	e21 := testEval(".foobar")
	if len(e21) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(e21))
	}

	e22 := testEval(".depth .foobar")
	if len(e22) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e22))
	}

	e23 := testEval("div[id='e'] ~ p")
	if len(e23) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e23))
	}

	e24 := testEval("div[id='e'] > p")
	if len(e24) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e24))
	}

	e25 := testEval("p > span")
	if len(e25) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e25))
	}

	e26 := testEval("div + p")
	if len(e26) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(e26))
	}

	e27 := testEval("p:nth-child(2n+1)")
	if len(e27) != 8 {
		t.Errorf("wrong number of items. got=%d, expected=8", len(e27))
	}

	e28 := testEval("p:nth-child(odd)")
	if len(e28) != 8 {
		t.Errorf("wrong number of items. got=%d, expected=8", len(e28))
	}

	e29 := testEval("p:nth-child(2n)")
	if len(e29) != 9 {
		t.Errorf("wrong number of items. got=%d, expected=9", len(e29))
	}

	e30 := testEval("p:nth-child(even)")
	if len(e30) != 9 {
		t.Errorf("wrong number of items. got=%d, expected=9", len(e30))
	}

	e31 := testEval("p:nth-child(4n + 1)")
	if len(e31) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e31))
	}

	e32 := testEval("p:nth-child(4n + 2)")
	if len(e32) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e32))
	}

	e33 := testEval("p:nth-child(4n + 3)")
	if len(e33) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(e33))
	}

	e34 := testEval("p:nth-child(4n + 4)")
	if len(e34) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e34))
	}

	e35 := testEval("p:nth-child(n)")
	if len(e35) != 17 {
		t.Errorf("wrong number of items. got=%d, expected=17", len(e35))
	}

	e36 := testEval(":nth-child(10n-1)")
	if len(e36) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(e36))
	}
	if e36[0].Data != "h4" {
		t.Errorf("wrong node selected")
	}

	e37 := testEval(":nth-child(10n+9)")
	if len(e37) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(e37))
	}
	if e37[0].Data != "h4" {
		t.Errorf("wrong node selected")
	}

	e38 := testEval(":nth-child(0n+5)")
	if len(e38) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e38))
	}
	if e38[0].Data != "h2" {
		t.Errorf("wrong node selected")
	}

	e39 := testEval("p:nth-child(0n+2)")
	if len(e39) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(e39))
	}

	e40 := testEval(":nth-child(5)")
	if len(e40) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e40))
	}
	if e40[0].Data != "h2" {
		t.Errorf("wrong node selected")
	}

	e41 := testEval("p:nth-child(1n+0)")
	if len(e41) != 17 {
		t.Errorf("wrong number of items. got=%d, expected=17", len(e41))
	}

	e42 := testEval("p:nth-child( +3n - 2 )")
	if len(e42) != 7 {
		t.Errorf("wrong number of items. got=%d, expected=7", len(e42))
	}

	e43 := testEval("p:nth-child( -n+ 6)")
	if len(e43) != 8 {
		t.Errorf("wrong number of items. got=%d, expected=8", len(e43))
	}

	e44 := testEval("p:nth-last-child(-n+2)")
	if len(e44) != 7 {
		t.Errorf("wrong number of items. got=%d, expected=7", len(e44))
	}

	e45 := testEval("p:nth-last-child(odd)")
	if len(e45) != 6 {
		t.Errorf("wrong number of items. got=%d, expected=6", len(e45))
	}

	e46 := testEval("p:nth-last-child(2n+1)")
	if len(e46) != 6 {
		t.Errorf("wrong number of items. got=%d, expected=6", len(e46))
	}

	e47 := testEval("p:nth-last-child(2n-1)")
	if len(e47) != 6 {
		t.Errorf("wrong number of items. got=%d, expected=6", len(e47))
	}

	e48 := testEval("p:nth-last-child(2n)")
	if len(e48) != 11 {
		t.Errorf("wrong number of items. got=%d, expected=11", len(e48))
	}

	e49 := testEval("p:nth-last-child(even)")
	if len(e49) != 11 {
		t.Errorf("wrong number of items. got=%d, expected=11", len(e49))
	}

	e50 := testEval("div:nth-of-type(2n+1)")
	if len(e50) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e50))
	}

	e51 := testEval("div:nth-of-type(-n+3)")
	if len(e51) != 6 {
		t.Errorf("wrong number of items. got=%d, expected=6", len(e51))
	}

	e52 := testEval("div:nth-of-type(-n+2)")
	if len(e52) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e52))
	}

	e53 := testEval("div:nth-of-type(-n+1)")
	if len(e53) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e53))
	}

	e54 := testEval("div:nth-of-type(-n)")
	if len(e54) != 0 {
		t.Errorf("wrong number of items. got=%d, expected=0", len(e54))
	}

	e56 := testEval("p:nth-last-of-type(2n+3)")
	if len(e56) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e56))
	}

	e57 := testEval("p:nth-last-of-type(-n+1)")
	if len(e57) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e56))
	}

	e58 := testEval("p:nth-last-of-type(-n+2)")
	if len(e58) != 8 {
		t.Errorf("wrong number of items. got=%d, expected=8", len(e58))
	}

	e59 := testEval("p:nth-last-of-type(-3n+1)")
	if len(e59) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e59))
	}

	e60 := testEval("p:nth-of-type(-2n+2)")
	if len(e60) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(e60))
	}

	e61 := testEval("p:nth-of-type(-2n-1)")
	if len(e61) != 0 {
		t.Errorf("wrong number of items. got=%d, expected=0", len(e61))
	}

	e62 := testEval("p:nth-of-type(-2n+1)")
	if len(e62) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e62))
	}

	e63 := testEval("p:nth-last-child(-2n+2)")
	if len(e63) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e63))
	}

	e64 := testEval("p:nth-last-child(2)")
	if len(e64) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e64))
	}

	e65 := testEval("p:nth-last-child(-4n+2)")
	if len(e65) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e65))
	}

	e66 := testEval("body > p:nth-of-type(n+2):nth-last-of-type(n+2)")
	if len(e66) != 9 {
		t.Errorf("wrong number of items. got=%d, expected=9", len(e66))
	}

	e67 := testEval("body > p:not(:first-of-type):not(:last-of-type)")
	if len(e67) != 9 {
		t.Errorf("wrong number of items. got=%d, expected=9", len(e67))
	}

	e68 := testEval("div > p:first-child")
	if len(e68) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e68))
	}

	e69 := testEval("p:first-child")
	if len(e69) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e69))
	}

	e70 := testEval("p:last-child")
	if len(e70) != 3 {
		t.Errorf("wrong number of items. got=%d, expected=3", len(e70))
	}

	e71 := testEval("p:first-of-type")
	if len(e71) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e71))
	}

	e72 := testEval("p:last-of-type")
	if len(e72) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e72))
	}

	e73 := testEval("*:not(p)")
	if len(e73) != 22 {
		t.Errorf("wrong number of items. got=%d, expected=22", len(e73))
	}

	e74 := testEval("body p *")
	if len(e74) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e74))
	}
	if e74[0].Data != "span" {
		t.Errorf("wrong selected")
	}

	e75 := testEval("h6 ~ p")
	if len(e75) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e75))
	}

	e76 := testEval("p.foo-bar + h6")
	if len(e76) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e76))
	}

	e77 := testEval("body > p")
	if len(e77) != 11 {
		t.Errorf("wrong number of items. got=%d, expected=11", len(e77))
	}

	e78 := testEval("p:only-of-type")
	if len(e78) != 2 {
		t.Errorf("wrong number of items. got=%d, expected=2", len(e78))
	}

	e79 := testEval("span:only-child")
	if len(e79) != 5 {
		t.Errorf("wrong number of items. got=%d, expected=5", len(e79))
	}

	e80 := testEval(":empty")
	if len(e80) != 1 {
		t.Errorf("wrong number of items. got=%d, expected=1", len(e80))
	}

	e81 := testEval("p:not(even)")
	if len(e81) != 17 {
		t.Errorf("wrong number of items. got=%d, expected=17", len(e81))
	}

	e82 := testEval("body *")
	if len(e82) != 35 {
		t.Errorf("wrong number of items. got=%d, expected=35", len(e82))
	}

	e83 := testEval("body *:not(h1,h2,h3,h4,h5,h6)")
	if len(e83) != 29 {
		t.Errorf("wrong number of items. got=%d, expected=29", len(e83))
	}

	e84 := testEval("body p:not(:nth-child(4n+1),:nth-child(4n+2),:nth-child(4n+3))")
	if len(e84) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e84))
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
