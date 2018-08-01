package xml3

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mmarkdown/mmark/mast"
)

func (r *Renderer) titleBlock(w io.Writer, t *mast.Title) {
	d := t.TitleData
	if d == nil {
		return
	}

	// rfc tag
	attrs := attributes(
		[]string{"ipr", "submissionType", "xml:lang", "xmlns", "consensus"},
		[]string{d.Ipr, "IETF", "en", "http://www.w3.org/2001/XInclude", fmt.Sprintf("%t", d.Consensus)},
	)
	attrs = append(attrs, attributes(
		[]string{"updates", "obsoletes"},
		[]string{intSliceToString(d.Updates), intSliceToString(d.Obsoletes)},
	)...)
	r.outTag(w, "<rfc", attrs)
	r.cr(w)

	// front tag - todo

	attrs = attributes([]string{"abbrev"}, []string{d.Abbrev})
	r.outTag(w, "<title", attrs)
	r.outs(w, d.Title)
	r.outs(w, "</title>")
	r.cr(w)

	for _, author := range d.Author {
		r.titleAuthor(w, author)
	}

	return
}

// titleAuthor outputs the author.
func (r *Renderer) titleAuthor(w io.Writer, a mast.Author) {

	attrs := attributes(
		[]string{"role", "initials", "surname", "fullname"},
		[]string{a.Role, a.Initials, a.Surname, a.Fullname},
	)

	r.outTag(w, "<author", attrs)
	r.cr(w)

	r.outTagContent(w, "<organization", attributes([]string{"abbrev"}, []string{a.OrganizationAbbrev}), a.Organization)

	r.outs(w, "<address>")
	r.cr(w)
	r.outs(w, "<postal>")
	r.cr(w)

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

	r.cr(w)
	r.outs(w, "</postal>")

	r.outTagContent(w, "<phone", nil, a.Address.Phone)
	r.outTagContent(w, "<email", nil, a.Address.Email)
	r.outTagContent(w, "<uri", nil, a.Address.Uri)

	r.outs(w, "</address>")
	r.cr(w)
	r.outs(w, "</author>")
	r.cr(w)
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
