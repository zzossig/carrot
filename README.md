# 🥕carrot

[![Go Report Card](https://goreportcard.com/badge/github.com/zzossig/carrot)](https://goreportcard.com/report/github.com/zzossig/carrot)
> CSS Selectors Level 3 implementation

Carrot is built to select HTML node using CSS3 selectors.

## Basic Usage

```go
carrot := New().SetDoc("./eval/testdata/t.html")
e1 := carrot.Eval("h1") // return []*html.Node
err := carrot.Errors() // return []error
```