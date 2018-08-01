package mast

import (
	"time"

	"github.com/gomarkdown/markdown/ast"
)

// Title represents the TOML encoded title block.
type Title struct {
	ast.Leaf
	*TitleData
}

// NewTitle returns a pointer to TitleData with some defaults set.
func NewTitle() *Title {
	t := &Title{
		TitleData: &TitleData{
			PI: pi{
				Header: "__mmark_toml_pi_not_set",
				Footer: "__mmark_toml_pi_not_set",
			},
			Area: "Internet",
			Ipr:  "trust200902",
			Date: time.Now(),
		},
	}
	return t
}

// TitleData holds all the elements of the title.
type TitleData struct {
	Title  string
	Abbrev string

	DocName        string
	Consensus      bool
	Ipr            string
	Category       string
	Number         int // RFC number
	Obsoletes      []int
	Updates        []int
	PI             pi // Processing Instructions
	SubmissionType string

	Date      time.Time
	Area      string
	Workgroup string
	Keyword   []string
	Author    []Author
}

type Author struct {
	Initials           string
	Surname            string
	Fullname           string
	Organization       string
	OrganizationAbbrev string `toml:"abbrev"`
	Role               string
	Ascii              string
	Address            Address
}

type Address struct {
	Phone  string
	Email  string
	Uri    string
	Postal AddressPostal
}

type AddressPostal struct {
	Street     string
	City       string
	Code       string
	Country    string
	Region     string
	PostalLine []string

	// Plurals when these need to be specified multiple times.
	Streets   []string
	Cities    []string
	Codes     []string
	Countries []string
	Regions   []string
}

// PIs the processing instructions.
var PIs = []string{"toc", "symrefs", "sortrefs", "compact", "subcompact", "private", "topblock", "header", "footer", "comments"}

type pi struct {
	Toc        string
	Symrefs    string
	Sortrefs   string
	Compact    string
	Subcompact string
	Private    string
	Topblock   string
	Comments   string // Typeset cref's in the text.
	Header     string // Top-Left header, usually Internet-Draft.
	Footer     string // Bottom-Center footer, usually Expires ...
}
