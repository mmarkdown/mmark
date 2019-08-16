package lang

// New returns a new and initialized Lang.
func New(language string) Lang {
	l := Lang{language: language}

	// Add all lanaguages here, the keys should be named according to BCP47.
	l.m = map[string]Term{
		"en": Term{
			Footnotes:    "Footnotes",
			Bibliography: "Bibliography",
			Index:        "Index",
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
	Footnotes    string
	Bibliography string
	Index        string
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
