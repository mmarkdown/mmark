package mhtml

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
)

type RenderNodeFunc func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool)

func RenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node := node.(type) {
	case *mast.References:
		references(w, node, entering)
		return ast.GoToNext, true
	case *mast.Reference:
		reference(w, node, entering)
		return ast.GoToNext, true
	case *mast.Title:
		// nothing yet, here
		//title(w, node, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func references(w io.Writer, node ast.Node, entering bool) {
	println("references: TODO")
}

func reference(w io.Writer, node ast.Node, entering bool) {
	println("reference: TODO")
}
