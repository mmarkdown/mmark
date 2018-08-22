package xml2

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
	"github.com/mmarkdown/mmark/xml"
)

func (r *Renderer) titleBlock(w io.Writer, t *mast.Title) {
	// Order is fixed in RFC 7741.

	d := t.TitleData
	if d == nil {
		return
	}

	attrs := attributes(
		[]string{"ipr", "submissionType", "category", "xml:lang", "consensus"},
		[]string{d.Ipr, d.SeriesInfo.Stream, d.SeriesInfo.Status, "en", fmt.Sprintf("%t", d.Consensus)},
	)
	attrs = append(attrs, attributes(
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

	r.matter(w, &ast.DocumentMatter{Matter: ast.DocumentMatterFront})

	attrs = attributes([]string{"abbrev"}, []string{d.Abbrev})
	r.outTag(w, "<title", attrs)
	r.outs(w, d.Title)
	r.outs(w, "</title>")

	// use a fake xml rendering to hook into the generation of these title elements
	// defined there.
	faker := xml.NewRenderer(xml.RendererOptions{})

	for _, author := range d.Author {
		faker.TitleAuthor(w, author)
	}

	faker.TitleDate(w, d.Date)

	r.outTagContent(w, "<area", nil, d.Area)

	r.outTagContent(w, "<workgroup", nil, d.Workgroup)

	faker.TitleKeyword(w, d.Keyword)

	// abstract - handled by paragraph
	// note - handled by paragraph
	// boilerplate - not supported.

	return
}
