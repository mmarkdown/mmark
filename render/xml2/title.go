package xml2

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/mmarkdown/mmark/mast"
	"github.com/mmarkdown/mmark/render/xml"
)

func (r *Renderer) titleBlock(w io.Writer, t *mast.Title) {
	// Order is fixed in RFC 7749.

	if t.IsTriggerDash() {
		return
	}

	d := t.TitleData
	if d == nil {
		return
	}
	consensusToTerm := map[bool]string{
		false: "no",
		true:  "yes",
	}

	attrs := xml.Attributes(
		[]string{"ipr", "submissionType", "category", "xml:lang", "consensus"},
		[]string{d.Ipr, d.SeriesInfo.Stream, xml.StatusToCategory[d.SeriesInfo.Status], "en", consensusToTerm[d.Consensus]},
	)
	attrs = append(attrs, xml.Attributes(
		[]string{"updates", "obsoletes"},
		[]string{xml.IntSliceToString(d.Updates), xml.IntSliceToString(d.Obsoletes)},
	)...)

	// Depending on the SeriesInfo.Name we're dealing with an RFC or Internet-Draft.
	switch d.SeriesInfo.Name {
	case "RFC":
		attrs = append(attrs, `number="`+d.SeriesInfo.Value+"\"")
	case "Internet-Draft": // case sensitive? Or throw error in toml checker?
		attrs = append(attrs, `docName="`+d.SeriesInfo.Value+"\"")
	case "DOI":
		// ?
	}

	r.outTag(w, "<rfc", attrs)
	r.cr(w)

	r.outs(w, `<?rfc toc="yes"?>`)
	r.outs(w, `<?rfc symrefs="yes"?>`)
	r.outs(w, `<?rfc sortrefs="yes"?>`)
	r.outs(w, `<?rfc compact="yes"?>`)
	r.outs(w, `<?rfc subcompact="no"?>`)
	r.outs(w, `<?rfc comments="no"?>`)

	r.matter(w, &ast.DocumentMatter{Matter: ast.DocumentMatterFront})

	attrs = xml.Attributes([]string{"abbrev"}, []string{d.Abbrev})
	r.outTag(w, "<title", attrs)
	html.EscapeHTML(w, []byte(d.Title))
	r.outs(w, "</title>")

	// use a fake xml rendering to hook into the generation of these title elements defined there.
	faker := xml.NewRenderer(xml.RendererOptions{})

	for _, author := range d.Author {
		faker.TitleAuthor(w, author)
	}

	faker.TitleDate(w, d.Date)

	r.outs(w, "<area>")
	html.EscapeHTML(w, []byte(d.Area))
	r.outs(w, "</area>")

	r.outs(w, "<workgroup>")
	html.EscapeHTML(w, []byte(d.Workgroup))
	r.outs(w, "</workgroup>")

	faker.TitleKeyword(w, d.Keyword)

	// abstract - handled by paragraph
	// note - handled by paragraph
	// boilerplate - not supported.

	return
}
