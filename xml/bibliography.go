package xml

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
		r.outs(w, `<references><name>Informative References</name>`)
	case ast.CitationTypeNormative:
		r.outs(w, `<references><name>Normative References</name>`)
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
		tag = makeXiInclude(toolsIetfOrg, fmt.Sprintf("reference.RFC.%s.xml", node.Anchor[3:]))

	case bytes.HasPrefix(node.Anchor, []byte("I-D.")):
		hash := bytes.Index(node.Anchor, []byte("#"))
		if hash > 0 {
			// rewrite # to - and we have our link
			node.Anchor[hash] = '-'
			defer func() { node.Anchor[hash] = '#' }() // never know if this will be used again
		}
		tag = makeXiInclude(toolsIetfOrg, fmt.Sprintf("reference.I-D.draft-%s.xml", node.Anchor[4:]))
	}
	r.outs(w, tag)
	r.cr(w)
}

func makeXiInclude(url, reference string) string {
	// <xi:include href="https://xml2rfc.tools.ietf.org/public/rfc/bibxml/reference.RFC.2119.xml"/>
	return fmt.Sprintf("<xi:include href=\"%s/%s\"/>", url, reference)
}

var toolsIetfOrg = "https://xml2rfc.tools.ietf.org/public/rfc/bibxml"
