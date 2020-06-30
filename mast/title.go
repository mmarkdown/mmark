package mast

import (
	"time"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast/reference"
)

// Title represents the TOML encoded title block.
type Title struct {
	ast.Leaf
	*TitleData
	Trigger string // either triggered by %%% or ---
}

// NewTitle returns a pointer to TitleData with some defaults set.
func NewTitle(trigger byte) *Title {
	t := &Title{
		TitleData: &TitleData{
			Area:      "Internet",
			Ipr:       "trust200902",
			Consensus: true,
		},
	}
	t.Trigger = string([]byte{trigger, trigger, trigger})
	return t
}

const triggerDash = "---"

func (t *Title) IsTriggerDash() bool { return t.Trigger == triggerDash }

// TitleData holds all the elements of the title.
type TitleData struct {
	Title  string
	Abbrev string

	SeriesInfo     reference.SeriesInfo
	Consensus      bool
	Ipr            string // See https://tools.ietf.org/html/rfc7991#appendix-A.1
	Obsoletes      []int
	Updates        []int
	SubmissionType string // IETF, IAB, IRTF or independent, defaults to IETF.

	Date      time.Time
	Area      string
	Workgroup string
	Keyword   []string
	Author    []Author

	Language string
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
