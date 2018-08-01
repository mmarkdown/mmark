package xml3

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
)

// RenderHook renders nodes that are defined outside of the main markdown.ast. Currently
// we render:
//
// * the TOML title block.
func RenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	t, ok := node.(*mast.Title)
	if !ok {
		return ast.GoToNext, false
	}

	if t.TitleData == nil {
		println("nothing defined")
	} else {
		println(t.TitleData.Area)
		println(t.TitleData.Title)
	}

	return ast.GoToNext, true
}
