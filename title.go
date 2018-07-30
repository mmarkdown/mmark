package mmark

import (
	"io"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gomarkdown/markdown/ast"
)

// Title represents the TOML encoded title block.
type Title struct {
	ast.Leaf
	*content
}

type content struct {
	Title  string
	Abbrev string

	DocName        string
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
	Author    []author
}

// TitleHook will parse a title and add it to the ast tree.
func TitleHook(data []byte) (ast.Node, []byte, int) {
	// parse text between %%% and %%% and return it as a Title node.
	i := 0
	if len(data) < 3 {
		return nil, nil, 0
	}
	if data[i] != '%' && data[i+1] != '%' && data[i+2] != '%' {
		return nil, nil, 0
	}

	i += 3
	// search for end.
	for i < len(data) {
		if data[i] == '%' && data[i+1] == '%' && data[i+2] == '%' {
			break
		}
		i++
	}
	node := &Title{}

	block := &content{
		PI: pi{
			Header: "__mmark_toml_pi_not_set",
			Footer: "__mmark_toml_pi_not_set",
		},
		Area: "Internet",
		Ipr:  "trust200902",
		Date: time.Now(),
	}

	if _, err := toml.Decode(string(data[4:i]), &block); err != nil {
		log.Printf("Failure to parsing title block: %s", err.Error())
	}
	node.content = block

	return node, nil, i + 3
}

type author struct {
	Initials           string
	Surname            string
	Fullname           string
	Organization       string
	OrganizationAbbrev string `toml:"abbrev"`
	Role               string
	Ascii              string
	Address            address
}

type address struct {
	Phone  string
	Email  string
	Uri    string
	Postal addressPostal
}

type addressPostal struct {
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

func RenderHookHTML(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	t, ok := node.(*Title)
	if !ok {
		return ast.GoToNext, false
	}

	if t.content == nil {
		println("WHAT")
	} else {
		println(t.content.Area)
		println(t.content.Title)
	}

	return ast.GoToNext, true
}
