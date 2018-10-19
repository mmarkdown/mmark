package xml

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/mmarkdown/mmark/mast"
)

// TODO(miek): double check if this is how it works.

// StatusToCategory translate the status to a category.
var StatusToCategory = map[string]string{
	"standard":      "std",
	"informational": "info",
	"experimental":  "exp",
	"bcp":           "bcp",
	"fyi":           "fyi",
	"full-standard": "std",
	// historic??
}

func (r *Renderer) titleBlock(w io.Writer, t *mast.Title) {
	// Order is fixed in RFC 7991.

	if t.IsTriggerDash() {
		// it was not parsed, leave it alone.
		return
	}

	d := t.TitleData
	if d == nil {
		return
	}

	// rfc tag
	attrs := Attributes(
		[]string{"version", "ipr", "submissionType", "category", "xml:lang", "consensus", "xmlns:xi"},
		[]string{"3", d.Ipr, "IETF", StatusToCategory[d.SeriesInfo.Status], "en", fmt.Sprintf("%t", d.Consensus), "http://www.w3.org/2001/XInclude"},
	)
	attrs = append(attrs, Attributes(
		[]string{"updates", "obsoletes"},
		[]string{IntSliceToString(d.Updates), IntSliceToString(d.Obsoletes)},
	)...)
	// number is deprecated, but xml2rfc want's it here to generate an actual RFC.
	// But only if number is a integer (what a mess).
	if _, err := strconv.Atoi(t.SeriesInfo.Value); err == nil {
		attrs = append(attrs, Attributes(
			[]string{"number"},
			[]string{t.SeriesInfo.Value},
		)...)
	}
	r.outTag(w, "<rfc", attrs)
	r.cr(w)

	r.matter(w, &ast.DocumentMatter{Matter: ast.DocumentMatterFront})

	attrs = Attributes([]string{"abbrev"}, []string{d.Abbrev})
	r.outTag(w, "<title", attrs)
	r.outs(w, d.Title)
	r.outs(w, "</title>")

	r.titleSeriesInfo(w, d.SeriesInfo)

	for _, author := range d.Author {
		r.TitleAuthor(w, author)
	}

	r.TitleDate(w, d.Date)

	r.outTagContent(w, "<area", d.Area)

	r.outTagContent(w, "<workgroup", d.Workgroup)

	r.TitleKeyword(w, d.Keyword)

	// abstract - handled by paragraph
	// note - handled by paragraph
	// boilerplate - not supported.

	return
}

// TitleAuthor outputs the author.
func (r *Renderer) TitleAuthor(w io.Writer, a mast.Author) {

	attrs := Attributes(
		[]string{"role", "initials", "surname", "fullname"},
		[]string{a.Role, a.Initials, a.Surname, a.Fullname},
	)

	r.outTag(w, "<author", attrs)

	r.outTag(w, "<organization", Attributes([]string{"abbrev"}, []string{a.OrganizationAbbrev}))
	html.EscapeHTML(w, []byte(a.Organization))
	r.outs(w, "</organization>")

	r.outs(w, "<address>")
	r.outs(w, "<postal>")

	r.outTagContent(w, "<street", a.Address.Postal.Street)
	for _, street := range a.Address.Postal.Streets {
		r.outTagContent(w, "<street", street)
	}

	r.outTagMaybe(w, "<city", a.Address.Postal.City)
	for _, city := range a.Address.Postal.Cities {
		r.outTagContent(w, "<city", city)
	}

	r.outTagMaybe(w, "<code", a.Address.Postal.Code)
	for _, code := range a.Address.Postal.Codes {
		r.outTagContent(w, "<code", code)
	}

	r.outTagMaybe(w, "<country", a.Address.Postal.Country)
	for _, country := range a.Address.Postal.Countries {
		r.outTagContent(w, "<country", country)
	}

	r.outTagMaybe(w, "<region", a.Address.Postal.Region)
	for _, region := range a.Address.Postal.Regions {
		r.outTagContent(w, "<region", region)
	}

	r.outs(w, "</postal>")

	r.outTagMaybe(w, "<phone", a.Address.Phone)
	r.outTagMaybe(w, "<email", a.Address.Email)
	r.outTagMaybe(w, "<uri", a.Address.URI)

	r.outs(w, "</address>")
	r.outs(w, "</author>")
	r.cr(w)
}

// TitleDate outputs the date from the TOML title block.
func (r *Renderer) TitleDate(w io.Writer, d time.Time) {
	if d.IsZero() { // not specified
		r.outs(w, "<date/>\n")
		return
	}

	var attr = []string{}
	if x := d.Year(); x > 0 {
		attr = append(attr, fmt.Sprintf(`year="%d"`, x))
	}
	if d.Month() > 0 {
		attr = append(attr, d.Format("month=\"January\""))
	}
	if x := d.Day(); x > 0 {
		attr = append(attr, fmt.Sprintf(`day="%d"`, x))
	}
	r.outTag(w, "<date", attr)
	r.outs(w, "</date>\n")
}

// TitleKeyword outputs the keywords from the TOML title block.
func (r *Renderer) TitleKeyword(w io.Writer, keyword []string) {
	for _, k := range keyword {
		if k == "" {
			continue
		}
		r.outTagContent(w, "<keyword", k)
	}
}

// titleSeriesInfo outputs the seriesInfo from the TOML title block.
func (r *Renderer) titleSeriesInfo(w io.Writer, s mast.SeriesInfo) {
	attr := Attributes(
		[]string{"value", "stream", "status", "name"},
		[]string{s.Value, s.Stream, s.Status, s.Name},
	)

	r.outTag(w, "<seriesInfo", attr)
	r.outs(w, "</seriesInfo>\n")
}

// IntSliceToString converts and int slice to a string.
func IntSliceToString(is []int) string {
	if len(is) == 0 {
		return ""
	}
	s := []string{}
	for i := range is {
		s = append(s, strconv.Itoa(is[i]))
	}
	return strings.Join(s, ", ")
}
