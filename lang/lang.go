package lang

import (
	"strings"
)

// New returns a new and initialized Lang.
func New(language string) Lang {
	l := Lang{language: strings.ToLower(language)} // case insensitivity

	// Add all lanaguages here, the keys should be named according to BCP47.
	// The keys must be in all lower case for normalized lookup.
	l.m = map[string]Term{
		"en": {
			And:          "and",
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
			Bibliography: "Bibliografie",
			Footnotes:    "Voetnoten",
			Index:        "Index",
			See:          "zie",
			Section:      "sectie",
			UseCounter:   "gebruik nummer",
			UseTitle:     "gebruik titel",
		},
		"de": {
			And:          "und",
			Bibliography: "Literaturverzeichnis",
			Footnotes:    "Fußnoten",
			Index:        "Index",
			See:          "siehe",
			Section:      "abschnit",
		},
		"ja": {
			Bibliography: "参考文献",
			Footnotes:    "脚注",
			Index:        "索引",
		},
		"zh-cn": {
			Bibliography: "参考文献",
			Footnotes:    "注释",
			Index:        "索引",
		},
		"zh-tw": {
			Bibliography: "參考文獻",
			Footnotes:    "註釋",
			Index:        "索引",
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

func (l Lang) Footnotes() string {
	t, ok := l.m[l.language]
	if !ok {
		return l.m["en"].Footnotes
	}
	return t.Footnotes
}

func (l Lang) Bibliography() string {
	t, ok := l.m[l.language]
	if !ok {
		return l.m["en"].Bibliography
	}
	return t.Bibliography
}

func (l Lang) Index() string {
	t, ok := l.m[l.language]
	if !ok {
		return l.m["en"].Index
	}
	return t.Index
}

func (l Lang) Authors() string {
	t, ok := l.m[l.language]
	if !ok {
		return l.m["en"].Authors
	}
	return t.Authors
}

func (l Lang) And() string {
	t, ok := l.m[l.language]
	if !ok {
		return l.m["en"].And
	}
	return t.And
}

func (l Lang) WrittenBy() string {
	t, ok := l.m[l.language]
	if !ok {
		return l.m["en"].WrittenBy
	}
	return t.WrittenBy
}

func (l Lang) See() string {
	t, ok := l.m[l.language]
	if !ok {
		return l.m["en"].See
	}
	return t.See
}

func (l Lang) Section() string {
	t, ok := l.m[l.language]
	if !ok {
		return l.m["en"].Section
	}
	return t.Section
}

func (l Lang) UseCounter() string {
	t, ok := l.m[l.language]
	if !ok {
		return l.m["en"].UseCounter
	}
	return t.UseCounter
}

func (l Lang) UseTitle() string {
	t, ok := l.m[l.language]
	if !ok {
		return l.m["en"].UseTitle
	}
	return t.UseTitle
}
