package xml3

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"regexp"

	"github.com/mmarkdown/mmark/mast"
)

// TODO, expand this to use regular expressions.
var (
	rfcRe = regexp.MustCompile(`/RFC(/d+)/`)
)

func (r *Renderer) references(w io.Writer, node *mast.References, entering bool) {
	if entering {
		r.sectionClose(w)
		r.section = nil
		r.outs(w, "<references>\n")
		return
	}

	r.outs(w, "</references>\n")
}

func (r *Renderer) reference(w io.Writer, node *mast.Reference) {
	tag := ""
	switch {
	case bytes.HasPrefix(node.Anchor, []byte("RFC")):
		tag = makeXiInclude(toolsIetfOrg, fmt.Sprintf("reference.RFC.%s.xml", node.Anchor[3:]))
	}
	r.outs(w, tag)
	r.cr(w)
}

func makeXiInclude(url, reference string) string {
	// <xi:include href="http://xml2rfc.tools.ietf.org/public/rfc/bibxml/reference.RFC.2119.xml"/>
	return fmt.Sprintf("<xi:include href=\"%s\"/>", path.Join(url, reference))
}

var toolsIetfOrg = "https://xml2rfc.tools.ietf.org/public/rfc/bibxml"
