package xml2

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
)

func (r *Renderer) bibliography(w io.Writer, node *mast.Bibliography, entering bool) {
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
		tag = makeRFCInclude(toolsIetfOrg, fmt.Sprintf("reference.RFC.%s.xml", node.Anchor[3:]))
	case bytes.HasPrefix(node.Anchor, []byte("I-D.")):
		hash := bytes.Index(node.Anchor, []byte("#"))
		if hash > 0 {
			// rewrite # to - and we have our link
			node.Anchor[hash] = '-'
		}
		tag = makeRFCInclude(toolsIetfOrg, fmt.Sprintf("reference.I-D.draft-%s.xml", node.Anchor[4:]))
	}
	r.outs(w, tag)
	r.cr(w)
}

func makeRFCInclude(url, reference string) string {
	return fmt.Sprintf("<?rfc include=\"%s/%s\"?>", url, reference)
}

var toolsIetfOrg = "https://xml2rfc.tools.ietf.org/public/rfc/bibxml"
