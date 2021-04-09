package carrot

import "testing"

func TestCSS(t *testing.T) {
	carrot := New().SetDoc("./eval/testdata/t.html")

	e1 := carrot.Eval("h1")
	if len(e1) != 1 {
		t.Errorf("length should be 1. got=%d", len(e1))
	}
	if e1[0].Data != "h1" {
		t.Errorf("wrong selected")
	}

	e2 := carrot.Eval("h6 ~ p")
	if len(e2) != 4 {
		t.Errorf("wrong number of items. got=%d, expected=4", len(e2))
	}

	e3 := New().SetDocS("<div>a</div><div>b</div>").Eval("div")
	if len(e3) != 2 {
		t.Errorf("length should be 2. got=%d", len(e3))
	}
}
