package lang

import (
	"strings"
)

// New returns a new and initialized Lang.
func New(language string) Lang {
	l := Lang{language: strings.ToLower(language)} // case insensitivity
	// Add all languages here, the keys should be named according to BCP47.
	// The keys must be in all lower case for normalized lookup.
	l.m = map[string]Term{
		"en": {
			And:          "and",
			Of:           "of",
			Authors:      "Authors",
			Bibliography: "Bibliography",
			Footnotes:    "Footnotes",
			Index:        "Index",
			WrittenBy:    "Written by",
			See:          "see",
			Section:      "section",
			UseCounter:   "use counter",
			UseTitle:     "use title",
		},
		"nl": {
			And:          "en",
			Of:           "of",
			Bibliography: "Bibliografie",
			Footnotes:    "Voetnoten",
			Index:        "Index",
			WrittenBy:    "Geschreven door",
			See:          "zie",
			Section:      "sectie",
			UseCounter:   "gebruik nummer",
			UseTitle:     "gebruik titel",
		},
		"de": {
			And:          "und",
			Of:           "von",
			Bibliography: "Literaturverzeichnis",
			Footnotes:    "Fußnoten",
			Index:        "Index",
			WrittenBy:    "Geschrieben von",
			See:          "siehe",
			Section:      "Abschnitt",
			UseCounter:   "Zähler benutzen",
			UseTitle:     "Titel benutzen",
		},
		"ja": {
			And:          "(no translation!)",
			Of:           "(no translation!)",
			Bibliography: "参考文献",
			Footnotes:    "脚注",
			Index:        "索引",
			WrittenBy:    "(no translation!)",
			See:          "(no translation!)",
			Section:      "(no translation!)",
			UseCounter:   "(no translation!)",
			UseTitle:     "(no translation!)",
		},
		"zh-cn": {
			And:          "(no translation!)",
			Of:           "(no translation!)",
			Bibliography: "参考文献",
			Footnotes:    "注释",
			Index:        "索引",
			WrittenBy:    "(no translation!)",
			See:          "(no translation!)",
			Section:      "(no translation!)",
			UseCounter:   "(no translation!)",
			UseTitle:     "(no translation!)",
		},
		"zh-tw": {
			And:          "(no translation!)",
			Of:           "(no translation!)",
			Bibliography: "參考文獻",
			Footnotes:    "註釋",
			Index:        "索引",
			WrittenBy:    "(no translation!)",
			See:          "(no translation!)",
			Section:      "(no translation!)",
			UseCounter:   "(no translation!)",
			UseTitle:     "(no translation!)",
		},
	}

	return l
}

// Lang maps a language to the terms we use in the document. We use an 'int' as to use the parser.Flags
// to indicate which language we'are using.
type Lang struct {
	language string
	m        map[string]Term
}

// Term contains the specific terms for translation.
type Term struct {
	And          string
	Of           string
	Authors      string
	Bibliography string
	Footnotes    string
	Index        string
	WrittenBy    string

	// for cross references
	See        string
	Section    string
	UseCounter string
	UseTitle   string
}

func (l Lang) Field(f string) string {
	m, ok := l.m[l.language]
	if !ok {
		m = l.m["en"]
	}
	switch strings.ToLower(f) {
	case "and":
		return m.And
	case "of":
		return m.Of
	case "Authors":
		return m.Authors
	case "bibliography":
		return m.Bibliography
	case "footnotes":
		return m.Footnotes
	case "index":
		return m.Index
	case "writtenby":
		return m.WrittenBy
	case "see":
		return m.See
	case "section":
		return m.Section
	case "usecounter":
		return m.UseCounter
	case "usetitle":
		return m.UseTitle
	}
	return ""
}

func (l Lang) Footnotes() string    { return l.Field("footnotes") }
func (l Lang) Bibliography() string { return l.Field("bibliography") }
func (l Lang) Index() string        { return l.Field("index") }
func (l Lang) Authors() string      { return l.Field("authors") }
func (l Lang) And() string          { return l.Field("and") }
func (l Lang) Of() string           { return l.Field("of") }
func (l Lang) WrittenBy() string    { return l.Field("writtenby") }
func (l Lang) See() string          { return l.Field("see") }
func (l Lang) Section() string      { return l.Field("section") }
func (l Lang) UseCounter() string   { return l.Field("usecounter") }
func (l Lang) UseTitle() string     { return l.Field("usetitle") }
