package mhtml

import (
	"fmt"
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
		// outout toml title block in html.
		//title(w, node, entering)
		return ast.GoToNext, true
	case *mast.Indices:
		if !entering {
			return ast.GoToNext, true
		}
		io.WriteString(w, "<h1>Index</h1>\n")
		return ast.GoToNext, true
	case *mast.IndexItem:
		if !entering {
			return ast.GoToNext, true
		}
		indexItem(w, node)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func indexItem(w io.Writer, node *mast.IndexItem) {
	for i := range node.Items {
		w.Write(node.Items[i].Item)
		w.Write([]byte("\n "))
		w.Write(node.Items[i].Subitem)
		w.Write([]byte(" "))
		fmt.Fprintf(w, "%d\n", node.Items[i].ID)
	}
}

func references(w io.Writer, node ast.Node, entering bool) {
	println("references: TODO")
}

func reference(w io.Writer, node ast.Node, entering bool) {
	println("reference: TODO")
}
