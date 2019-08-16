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
		"nl": Term{
			Footnotes:    "Voetnoten",
			Bibliography: "Bibliografie",
			Index:        "Index",
		},
		"de": Term{
			Footnotes:    "Fußnoten",
			Bibliography: "Literaturverzeichnis",
			Index:        "Index",
		},
		"ja": Term{
			Footnotes:    "脚注",
			Bibliography: "参考文献",
			Index:        "索引",
		},
		"zh-cn": Term{
			Footnotes:    "注释",
			Bibliography: "参考文献",
			Index:        "索引",
		},
		"zh-tw": Term{
			Footnotes:    "註釋",
			Bibliography: "參考文獻",
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
