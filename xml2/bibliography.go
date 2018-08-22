package xml2

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
)

// TODO, expand this to use regular expressions.
var (
	rfcRe = regexp.MustCompile(`/RFC(/d+)/`)
)

func (r *Renderer) bibliography(w io.Writer, node *mast.Bibliography, entering bool) {
	if !entering {
		r.outs(w, "</references>\n")
		return
	}

	r.sectionClose(w)
	r.section = nil

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
	case bytes.HasPrefix(node.Anchor, []byte("RFC")): // TODO(miek): use regexp here.
		tag = makeRFCInclude(toolsIetfOrg, fmt.Sprintf("reference.RFC.%s.xml", node.Anchor[3:]))
	}
	r.outs(w, tag)
	r.cr(w)
}

func makeRFCInclude(url, reference string) string {
	// <xi:include href="https://xml2rfc.tools.ietf.org/public/rfc/bibxml/reference.RFC.2119.xml"/>
	return fmt.Sprintf("<?rfc include href=\"%s/%s\"?>", url, reference)
}

var toolsIetfOrg = "https://xml2rfc.tools.ietf.org/public/rfc/bibxml"
