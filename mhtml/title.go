package mhtml

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
)

// RenderHook renders nodes that are defined outside of the main markdown.ast. Currently
// we render:
//
// * the TOML title block.
func RenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	t, ok := node.(*Title)
	if !ok {
		return ast.GoToNext, false
	}

	if t.content == nil {
		println("nothing defined")
	} else {
		println(t.content.Area)
		println(t.content.Title)
	}

	return ast.GoToNext, true
}
