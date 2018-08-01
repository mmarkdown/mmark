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
	SubmissionType string

	Date      time.Time
	Area      string
	Workgroup string
	Keyword   []string
	Author    []Author
}

// Author denotes an RFC author.
type Author struct {
	Initials           string
	Surname            string
	Fullname           string
	Organization       string
	OrganizationAbbrev string `toml:"abbrev"`
	Role               string
	ASCII              string
	Address            Address
}

// Address denotes the address of an RFC author.
type Address struct {
	Phone  string
	Email  string
	URI    string
	Postal AddressPostal
}

// AddressPostal denotes the postal address of an RFC author.
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
