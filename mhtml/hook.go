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
		// outout toml title block in html.
		//title(w, node, entering)
		return ast.GoToNext, true
	case *mast.Indices:
		if !entering {
			return ast.GoToNext, true
		}
		io.WriteString(w, "<h1>Index</h1>\n")
		return ast.GoToNext, true
	case *mast.IndexLetter:
		if !entering {
			return ast.GoToNext, true
		}
		io.WriteString(w, `<h3 class="index letter">`)
		io.WriteString(w, string(node.Literal))
		io.WriteString(w, "</h3>\n")

		return ast.GoToNext, true
	case *mast.IndexItem:
		if !entering {
			return ast.GoToNext, true
		}
		span := wrapInSpan(node.Item, "index item")
		io.WriteString(w, span)
		for i := range node.IDs {
			io.WriteString(w, "#"+node.IDs[i])
		}

		return ast.GoToNext, true
	case *mast.IndexSubItem:
		if !entering {
			return ast.GoToNext, true
		}
		span := wrapInSpan(node.Subitem, "index subitem")
		io.WriteString(w, span)
		for i := range node.IDs {
			io.WriteString(w, "#"+node.IDs[i])
		}
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

func wrapInSpan(content []byte, class string) string {
	s := "<span "
	s += `class="` + class + `">` + string(content) + "</span>\n"
	return s
}
