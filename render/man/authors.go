package man

import (
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
)

// authors creates a 'Authors' section that can be included. Only the 'fullname' is used.
// If not authors are specified, it will return nil.
func (r *Renderer) authors(w io.Writer, _ *mast.Authors, entering bool) {
	if !entering {
		return
	}
	if r.Title == nil {
		return
	}
	author := r.Title.TitleData.Author
	if len(author) == 0 {
		return
	}

	// create a node and call render on it.
	node := &ast.Heading{Level: 1}
	authors := r.opts.Language.Authors()
	ast.AppendChild(node, &ast.Text{ast.Leaf{Literal: []byte(authors)}})
	la := len(author)

	// Needs to use the translation stuff
	para := &ast.Paragraph{}
	text := r.opts.Language.WrittenBy() + " "
	// combine them with commas, add the last one with 'and'
	switch la {
	case 1:
		text += author[0].Fullname
	case 2:
		text += author[0].Fullname + " " + r.opts.Language.And() + " " + author[1].Fullname
	default:
		names := make([]string, len(author))
		for i, a := range author {
			names[i] = a.Fullname
		}
		text += strings.Join(names[:len(names)-2], ", ")
		text += " " + r.opts.Language.And() + " " + names[len(names)-1]
	}
	text += "."

	ast.AppendChild(para, &ast.Text{ast.Leaf{Literal: []byte(text)}})
	ast.AppendChild(node, para)

	ast.WalkFunc(node, func(node ast.Node, entering bool) ast.WalkStatus {
		return r.RenderNode(w, node, entering)
	})

	return
}
