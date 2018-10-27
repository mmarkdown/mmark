package xml2

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
	"github.com/mmarkdown/mmark/render/xml"
)

func (r *Renderer) out(w io.Writer, d []byte)  { w.Write(d) }
func (r *Renderer) outs(w io.Writer, s string) { io.WriteString(w, s) }
func (r *Renderer) cr(w io.Writer)             { r.outs(w, "\n") }

func (r *Renderer) outAttr(w io.Writer, attrs []string) {
	if len(attrs) > 0 {
		io.WriteString(w, " ")
		io.WriteString(w, strings.Join(attrs, " "))
	}
}

func (r *Renderer) outTag(w io.Writer, name string, attrs []string) {
	s := name
	if len(attrs) > 0 {
		s += " " + strings.Join(attrs, " ")
	}
	io.WriteString(w, s+">")
}

func (r *Renderer) outOneOf(w io.Writer, outFirst bool, first string, second string) {
	if outFirst {
		r.outs(w, first)
	} else {
		r.outs(w, second)
	}
}

func (r *Renderer) outOneOfCr(w io.Writer, outFirst bool, first string, second string) {
	if outFirst {
		r.cr(w)
		r.outs(w, first)
	} else {
		r.outs(w, second)
		r.cr(w)
	}
}

func (r *Renderer) sectionClose(w io.Writer, new *ast.Heading) {
	defer func() {
		r.section = new
	}()

	if r.section == nil {
		return
	}

	if r.section.IsSpecial {
		tag := "</note>"
		if xml.IsAbstract(r.section.Literal) {
			tag = "</abstract>"
		}
		r.outs(w, tag)
		r.cr(w)
		return
	}

	tag := "</section>"
	curLevel := r.section.Level
	newLevel := 1 // close them all
	if new != nil {
		newLevel = new.Level
	}

	// subheading in a heading
	if newLevel > curLevel {
		return
	}

	for i := newLevel; i <= curLevel; i++ {
		r.outs(w, tag)
		r.cr(w)
	}
}

func (r *Renderer) ensureUniqueHeadingID(id string) string {
	for count, found := r.headingIDs[id]; found; count, found = r.headingIDs[id] {
		tmp := fmt.Sprintf("%s-%d", id, count+1)

		if _, tmpFound := r.headingIDs[tmp]; !tmpFound {
			r.headingIDs[id] = count + 1
			id = tmp
		} else {
			id = id + "-1"
		}
	}

	if _, found := r.headingIDs[id]; !found {
		r.headingIDs[id] = 0
	}

	return id
}

func appendLanguageAttr(node ast.Node, info []byte) {
	if len(info) == 0 {
		return
	}
	endOfLang := bytes.IndexAny(info, "\t ")
	if endOfLang < 0 {
		endOfLang = len(info)
	}
	mast.SetAttribute(node, "type", info[:endOfLang])
}

// isHangText returns true if the grandparent is a definition list term.
func isHangText(node ast.Node) bool {
	grandparent := node.GetParent().GetParent()
	if grandparent == nil {
		return false
	}
	if grandparent != nil {
		if li, ok := grandparent.(*ast.ListItem); ok {
			return li.ListFlags&ast.ListTypeTerm != 0
		}
	}
	return false
}
