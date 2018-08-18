package mhtml

import (
	"bytes"
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
	case *mast.IndexLetter:
		if !entering {
			return ast.GoToNext, true
		}
		io.WriteString(w, `<h3 class="idxletter">`)
		io.WriteString(w, string(node.Literal))
		io.WriteString(w, "</h3>")

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
	// First gather all the links for single items, and subitems. Both may have multiple links.
	itemLinks := map[string][]string{}
	subItemLinks := map[string][]string{} // index on item,subitem
	for _, n := range node.Items {
		item := string(n.Item)
		if n.Subitem == nil {
			itemLinks[item] = append(itemLinks[item], n.ID)
			continue
		}
		sub := string(n.Subitem)
		subItemLinks[item+","+sub] = append(subItemLinks[item+","+sub], n.ID)
	}

	// Now range again through the Items and assign each unique on the list of Ids
	links := map[*ast.Index][]string{}
	for _, n := range node.Items {
		item := string(n.Item)
		if refs, ok := itemLinks[item]; ok {
			links[n] = refs
			continue
		}
		sub := string(n.Subitem)
		if refs, ok := subItemLinks[item+","+sub]; ok {
			links[n] = refs
		}
	}

	for i := range node.Items {
		if bytes.Compare(node.Items[i].Item, prevItem) != 0 {
			w.Write(node.Items[i].Item)
			w.Write([]byte("\n "))
		}
		w.Write(node.Items[i].Subitem)
		w.Write([]byte(" "))
		fmt.Fprintf(w, "%s\n", node.Items[i].ID)
		prevItem = node.Items[i].Item
	}
}

func references(w io.Writer, node ast.Node, entering bool) {
	println("references: TODO")
}

func reference(w io.Writer, node ast.Node, entering bool) {
	println("reference: TODO")
}
