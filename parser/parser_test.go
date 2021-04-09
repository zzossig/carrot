package parser

import (
	"testing"

	"github.com/zzossig/carrot/lexer"
)

func TestExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1", "1"},
		{`"str"`, `"str"`},
		{`'str'`, `"str"`},
		{"h1", "h1"},
		{"h1, h2, h3", "h1, h2, h3"},
		{"*[hreflang|=en]", "*[hreflang|=en]"},
		{"[hreflang|=en]", "[hreflang|=en]"},
		{"*.warning", "*.warning"},
		{".warning", ".warning"},
		{"*#myid", "*#myid"},
		{"#myid", "#myid"},
		{"*", "*"},
		{"[att]", "[att]"},
		{"[att=val]", "[att=val]"},
		{"[att~=val]", "[att~=val]"},
		{"[att|=val]", "[att|=val]"},
		{"[att^=val]", "[att^=val]"},
		{"[att$=val]", "[att$=val]"},
		{"[att*=val]", "[att*=val]"},
		{"h1[title]", "h1[title]"},
		{`span[class='example']`, `span[class=example]`},
		{`span[hello="Cleveland"][goodbye="Columbus"]`, `span[hello=Cleveland][goodbye=Columbus]`},
		{`a[rel~="copyright"]`, `a[rel~=copyright]`},
		{`p.pastoral.marine`, `p.pastoral.marine`},
		{`.pastoral .marine`, `.pastoral .marine`},
		{`.pastoral > .marine`, `.pastoral > .marine`},
		{`.pastoral ~ .marine`, `.pastoral ~ .marine`},
		{`.pastoral + .marine`, `.pastoral + .marine`},
		{`.pastoral[a] .marine[b=c]`, `.pastoral[a] .marine[b=c]`},
		{`.pastoral[a][x][z] .marine[b=c][y][z]`, `.pastoral[a][x][z] .marine[b=c][y][z]`},
		{`h1#chapter1`, `h1#chapter1`},
		{`#chapter1`, `#chapter1`},
		{`*#z98y`, `*#z98y`},
		{`a.external:visited`, `a.external:visited`},
		{`a:focus:hover`, `a:focus:hover`},
		{`p.note:target`, `p.note:target`},
		{`*:target::before`, `*:target::before`},
		{`html:lang(fr-be)`, `html:lang(fr-be)`},
		{`:lang`, `:lang`},
		{`:lang(fr-be) > q`, `:lang(fr-be) > q`},
		{`tr:nth-child(2n+1)`, `tr:nth-child(2n+1)`},
		{`tr:nth-child(2n-1)`, `tr:nth-child(2n-1)`},
		{`tr:nth-child(odd)`, `tr:nth-child(odd)`},
		{`tr:nth-child(2n+0)`, `tr:nth-child(2n+0)`},
		{`:nth-child(10n-1)`, `:nth-child(10n-1)`},
		{`:nth-child(10n+9)`, `:nth-child(10n+9)`},
		{`foo:nth-child(0n+5)`, `foo:nth-child(0n+5)`},
		{`foo:nth-child(5)`, `foo:nth-child(5)`},
		{`bar:nth-child(1n+0)`, `bar:nth-child(n+0)`},
		{`bar:nth-child(n+0)`, `bar:nth-child(n+0)`},
		{`bar:nth-child(n)`, `bar:nth-child(n)`},
		{`tr:nth-child(2n+0)`, `tr:nth-child(2n+0)`},
		{`tr:nth-child(2n)`, `tr:nth-child(2n)`},
		{`:nth-child( 3n + 1 )`, `:nth-child(3n+1)`},
		{`:nth-child( +3n - 2 )`, `:nth-child(3n-2)`},
		{`:nth-child( -n+ 6)`, `:nth-child(-n+6)`},
		{`:nth-child( +6 )`, `:nth-child(6)`},
		{`:nth-child( -6 )`, `:nth-child(-6)`},
		{`body > h2:nth-of-type(n+2):nth-last-of-type(n+2)`, `body > h2:nth-of-type(n+2):nth-last-of-type(n+2)`},
		{`body > h2:not(:first-of-type):not(:last-of-type)`, `body > h2:not(:first-of-type):not(:last-of-type)`},
		{`div > p:first-child`, `div > p:first-child`},
		{`* > a:first-child`, `* > a:first-child`},
		{`dl dt:first-of-type`, `dl dt:first-of-type`},
		{`tr > td:last-of-type`, `tr > td:last-of-type`},
		{`button:not([DISABLED])`, `button:not([DISABLED])`},
		{`*:not(FOO)`, `*:not(FOO)`},
		{`h1 em`, `h1 em`},
		{`div * p`, `div * p`},
		{`div p *[href]`, `div p *[href]`},
		{`div ol>li p`, `div ol > li p`},
		{`div ol > li p`, `div ol > li p`},
		{`math + p`, `math + p`},
		{`h1.opener + h2`, `h1.opener + h2`},
		{`h1 ~ pre`, `h1 ~ pre`},
		{`H1 + *[REL=up]`, `H1 + *[REL=up]`},
		{`body *:not(h1,h2,h3,h4,h5,h6)`, `body *:not(h1, h2, h3, h4, h5, h6)`},
		{`a:has(> img)`, `a:has(> img)`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		e := p.ParseExpression()

		actual := e.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
