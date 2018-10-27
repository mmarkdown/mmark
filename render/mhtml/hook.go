package mhtml

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
)

var (
	// IndexReturnLinkContents is the string to use for index item return links.
	IndexReturnLinkContents = "<sup>[go]</sup>"
)

// RenderHook is used to render mmark specific AST nodes.
func RenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node := node.(type) {
	case *ast.Footnotes:
		if !entering {
			io.WriteString(w, "</h1>\n")
			return ast.GoToNext, true
		}
		io.WriteString(w, `<h1 id="footnote-section">Footnotes`)
	case *mast.Bibliography:
		if !entering {
			io.WriteString(w, "</dl>\n")
			return ast.GoToNext, true
		}
		io.WriteString(w, "<h1 id=\"bibliography-section\">Bibliography</h1>\n<div class=\"bibliography\">\n")
		io.WriteString(w, "<dl>\n")
		return ast.GoToNext, true
	case *mast.BibliographyItem:
		if !entering {
			return ast.GoToNext, true
		}
		bibliographyItem(w, node, entering)
		return ast.GoToNext, true
	case *mast.Title:
		// we out if in mmark.go with a hack to capture it.
		return ast.GoToNext, true
	case *mast.DocumentIndex:
		if !entering {
			io.WriteString(w, "\n</div>\n")
			return ast.GoToNext, true
		}
		io.WriteString(w, "<h1 id=\"index-section\">Index</h1>\n<div class=\"index\">\n")
		return ast.GoToNext, true
	case *mast.IndexLetter:
		if !entering {
			io.WriteString(w, "</ul>\n")
			io.WriteString(w, "</dd>\n")
			io.WriteString(w, "</dl>\n")
			return ast.GoToNext, true
		}
		// use id= idxref idxitm.
		io.WriteString(w, "<dl>\n")
		io.WriteString(w, `<dt>`)
		io.WriteString(w, string(node.Literal))
		io.WriteString(w, "</dt>\n")
		io.WriteString(w, "<dd>\n")
		io.WriteString(w, "<ul>\n")
		return ast.GoToNext, true
	case *mast.IndexItem:
		if !entering {
			io.WriteString(w, "</li>\n")
			return ast.GoToNext, true
		}
		io.WriteString(w, "<li>\n")
		w.Write(node.Item)
		return ast.GoToNext, true
	case *mast.IndexSubItem:
		if !entering {
			if lastSubItem(node) {
				io.WriteString(w, "</ul>\n")
			}
			io.WriteString(w, "</li>\n")
			return ast.GoToNext, true
		}
		if firstSubItem(node) {
			io.WriteString(w, "<ul>\n")
		}
		io.WriteString(w, "<li>\n")
		w.Write(node.Subitem)
		return ast.GoToNext, true
	case *mast.IndexLink:
		if !entering {
			io.WriteString(w, "</a>")
			return ast.GoToNext, true
		}
		io.WriteString(w, ` <a class="index-return" href="#`+string(node.Destination)+`">`)
		io.WriteString(w, IndexReturnLinkContents)
		return ast.GoToNext, true

	}
	return ast.GoToNext, false
}

func bibliographyItem(w io.Writer, bib *mast.BibliographyItem, entering bool) {
	io.WriteString(w, `<dt class="bibliography-cite" id="`+string(bib.Anchor)+`">`+fmt.Sprintf("[%s]", bib.Anchor)+"</dt>\n")
	io.WriteString(w, `<dd>`)
	io.WriteString(w, `<span class="bibliography-author">`+bib.Reference.Front.Author.Fullname+"</span>\n")
	io.WriteString(w, `<span class="bibliography-title">`+bib.Reference.Front.Title+"</span>\n")
	if bib.Reference.Format != nil && bib.Reference.Format.Target != "" {
		io.WriteString(w, `<a class="bliography-target" href="`+bib.Reference.Format.Target+"\">"+bib.Reference.Format.Target+"</a>\n")
	}
	if bib.Reference.Front.Date.Year != "" {
		io.WriteString(w, `<date class="bibliography-date">`+bib.Reference.Front.Date.Year+"</date>\n")
	}
	io.WriteString(w, "</dd>\n")
}

func firstSubItem(node ast.Node) bool {
	prev := ast.GetPrevNode(node)
	if prev == nil {
		return true
	}
	for prev != nil {
		_, ok := prev.(*mast.IndexSubItem)
		if ok {
			return false
		}
		prev = ast.GetPrevNode(prev)
	}
	return true
}

func lastSubItem(node ast.Node) bool {
	next := ast.GetNextNode(node)
	if next == nil {
		return true
	}
	for next != nil {
		_, ok := next.(*mast.IndexSubItem)
		if ok {
			return false
		}
		next = ast.GetNextNode(next)
	}
	return true
}
