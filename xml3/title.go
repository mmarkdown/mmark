package xml3

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
)

func (r *Renderer) titleBlock(w io.Writer, t *mast.Title) {
	// Order is fixed in RFC 7991.

	d := t.TitleData
	if d == nil {
		return
	}

	// rfc tag
	attrs := attributes(
		[]string{"version", "ipr", "submissionType", "xml:lang", "consensus", "xmlns:xi"},
		[]string{"3", d.Ipr, "IETF", "en", fmt.Sprintf("%t", d.Consensus), "http://www.w3.org/2001/XInclude"},
	)
	attrs = append(attrs, attributes(
		[]string{"updates", "obsoletes"},
		[]string{intSliceToString(d.Updates), intSliceToString(d.Obsoletes)},
	)...)
	r.outTag(w, "<rfc", attrs)
	r.cr(w)

	r.matter(w, &ast.DocumentMatter{Matter: ast.DocumentMatterFront})

	attrs = attributes([]string{"abbrev"}, []string{d.Abbrev})
	r.outTag(w, "<title", attrs)
	r.outs(w, d.Title)
	r.outs(w, "</title>")

	r.titleSeriesInfo(w, d.SeriesInfo)

	for _, author := range d.Author {
		r.titleAuthor(w, author)
	}

	r.titleDate(w, d.Date)

	r.outTagContent(w, "<area", nil, d.Area)

	r.outTagContent(w, "<workgroup", nil, d.Workgroup)

	for _, k := range d.Keyword {
		if k == "" {
			continue
		}
		r.outTagContent(w, "<keyword", nil, k)
	}
	// abstract - handled by paragraph
	// note - handled by paragraph
	// boilerplate - not supported.

	return
}

// titleAuthor outputs the author.
func (r *Renderer) titleAuthor(w io.Writer, a mast.Author) {

	attrs := attributes(
		[]string{"role", "initials", "surname", "fullname"},
		[]string{a.Role, a.Initials, a.Surname, a.Fullname},
	)

	r.outTag(w, "<author", attrs)

	r.outTagContent(w, "<organization", attributes([]string{"abbrev"}, []string{a.OrganizationAbbrev}), a.Organization)

	r.outs(w, "<address>")
	r.outs(w, "<postal>")

	r.outTagContent(w, "<street", nil, a.Address.Postal.Street)
	for _, street := range a.Address.Postal.Streets {
		r.outTagContent(w, "<street", nil, street)
	}

	r.outTagContent(w, "<city", nil, a.Address.Postal.City)
	for _, city := range a.Address.Postal.Cities {
		r.outTagContent(w, "<city", nil, city)
	}

	r.outTagContent(w, "<code", nil, a.Address.Postal.Code)
	for _, code := range a.Address.Postal.Codes {
		r.outTagContent(w, "<code", nil, code)
	}

	r.outTagContent(w, "<country", nil, a.Address.Postal.Country)
	for _, country := range a.Address.Postal.Countries {
		r.outTagContent(w, "<country", nil, country)
	}

	r.outTagContent(w, "<region", nil, a.Address.Postal.Region)
	for _, region := range a.Address.Postal.Regions {
		r.outTagContent(w, "<region", nil, region)
	}

	r.outs(w, "</postal>")

	r.outTagContent(w, "<phone", nil, a.Address.Phone)
	r.outTagContent(w, "<email", nil, a.Address.Email)
	r.outTagContent(w, "<uri", nil, a.Address.URI)

	r.outs(w, "</address>")
	r.outs(w, "</author>")
	r.cr(w)
}

// titleDate outputs the date from the TOML title block.
func (r *Renderer) titleDate(w io.Writer, d time.Time) {
	var attr = []string{}

	if x := d.Year(); x > 0 {
		attr = append(attr, fmt.Sprintf(`year="%d"`, x))
	}
	if x := d.Month(); x > 0 {
		attr = append(attr, fmt.Sprintf(`month="%d"`, x))
	}
	if x := d.Day(); x > 0 {
		attr = append(attr, fmt.Sprintf(`day="%d"`, x))
	}
	r.outTag(w, "<date", attr)
	r.outs(w, "</date>\n")
}

// titleSeriesInfo outputs the seriesInfo from the TOML title block.
func (r *Renderer) titleSeriesInfo(w io.Writer, s mast.SeriesInfo) {
	attr := attributes(
		[]string{"value", "stream", "status", "name"},
		[]string{s.Value, s.Stream, s.Status, s.Name},
	)

	r.outTag(w, "<seriesInfo", attr)
	r.outs(w, "</seriesInfo>\n")
}

func intSliceToString(is []int) string {
	if len(is) == 0 {
		return ""
	}
	s := []string{}
	for i := range is {
		s = append(s, strconv.Itoa(is[i]))
	}
	return strings.Join(s, ", ")
}
