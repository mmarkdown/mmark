// Package reference defines the elements of a <reference> block.
package reference

import "encoding/xml"

// Author is the reference author.
type Author struct {
	Fullname string `xml:"fullname,attr"`
	Initials string `xml:"initials,attr"`
	Surname  string `xml:"surname,attr"`
}

// Date is the reference date.
type Date struct {
	Year  string `xml:"year,attr,omitempty"`
	Month string `xml:"month,attr,omitempty"`
	Day   string `xml:"day,attr,omitempty"`
}

// Front the reference <front>.
type Front struct {
	Title  string `xml:"title"`
	Author Author `xml:"author"`
	Date   Date   `xml:"date"`
}

// Format is the reference <format>.
type Format struct {
	Type   string `xml:"type,attr,omitempty"`
	Target string `xml:"target,attr"`
}

// Reference is the entire <reference> structure.
type Reference struct {
	XMLName xml.Name `xml:"reference"`
	Anchor  string   `xml:"anchor,attr"`
	Front   Front    `xml:"front"`
	Format  *Format  `xml:"format,omitempty"`
}
