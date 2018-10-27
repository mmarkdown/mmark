package xml2

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
	"github.com/mmarkdown/mmark/render/xml"
)

func (r *Renderer) bibliography(w io.Writer, node *mast.Bibliography, entering bool) {
	if len(node.GetChildren()) == 0 {
		return
	}

	if !entering {
		r.outs(w, "</references>\n")
		return
	}

	r.sectionClose(w, nil)

	switch node.Type {
	case ast.CitationTypeInformative:
		r.outs(w, `<references title="Informative References">`)
	case ast.CitationTypeNormative:
		r.outs(w, `<references title="Normative References">`)
	}
	r.cr(w)
}

func (r *Renderer) bibliographyItem(w io.Writer, node *mast.BibliographyItem) {
	if node.Raw != nil {
		r.out(w, node.Raw)
		r.cr(w)
		return
	}

	tag := ""
	switch {
	case bytes.HasPrefix(node.Anchor, []byte("RFC")):
		tag = makeRFCInclude(xml.ToolsRFC, fmt.Sprintf("reference.RFC.%s.xml", node.Anchor[3:]))

	case bytes.HasPrefix(node.Anchor, []byte("W3C.")):
		tag = makeRFCInclude(xml.ToolsW3C, fmt.Sprintf("reference.W3C.%s.xml", node.Anchor[4:]))

	case bytes.HasPrefix(node.Anchor, []byte("I-D.")):
		hash := bytes.Index(node.Anchor, []byte("#"))
		if hash > 0 {
			// rewrite # to - and we have our link
			node.Anchor[hash] = '-'
			defer func() { node.Anchor[hash] = '#' }() // never know if this will be used again
		}
		tag = makeRFCInclude(xml.ToolsID, fmt.Sprintf("reference.I-D.%s.xml", node.Anchor[4:]))
	}
	r.outs(w, tag)
	r.cr(w)
}

func makeRFCInclude(url, reference string) string {
	return fmt.Sprintf("<?rfc include=\"%s/%s\"?>", url, reference)
}
