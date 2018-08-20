package xml

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/mmarkdown/mmark/mast"
)

// TODO, expand this to use regular expressions.
var (
	rfcRe = regexp.MustCompile(`/RFC(/d+)/`)
)

func (r *Renderer) bibliography(w io.Writer, node *mast.Bibliography, entering bool) {
	if entering {
		r.sectionClose(w)
		r.section = nil
		r.outs(w, "<references>\n")
		return
	}

	r.outs(w, "</references>\n")
}

func (r *Renderer) bibliographyItem(w io.Writer, node *mast.BibliographyItem) {
	if node.RawXML != nil {
		r.out(w, node.RawXML)
		r.cr(w)
		return
	}

	tag := ""
	switch {
	case bytes.HasPrefix(node.Anchor, []byte("RFC")): // TODO(miek): use regexp here.
		tag = makeXiInclude(toolsIetfOrg, fmt.Sprintf("reference.RFC.%s.xml", node.Anchor[3:]))
	}
	r.outs(w, tag)
	r.cr(w)
}

func makeXiInclude(url, reference string) string {
	// <xi:include href="https://xml2rfc.tools.ietf.org/public/rfc/bibxml/reference.RFC.2119.xml"/>
	return fmt.Sprintf("<xi:include href=\"%s/%s\"/>", url, reference)
}

var toolsIetfOrg = "https://xml2rfc.tools.ietf.org/public/rfc/bibxml"
