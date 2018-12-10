package mast

import "github.com/gomarkdown/markdown/ast"

// LatexSpan represents markdown LaTeX span node, i.e. any string that matches:
// \\[a-zA-Z]{.+}.
type LatexSpan struct {
	ast.Leaf
}
