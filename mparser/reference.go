package mparser

import (
	"bytes"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/mast"
)

func CitationToReferences(p *parser.Parser, doc ast.Node) (normative, informative ast.Node) {
	seen := map[string]*mast.Reference{}

	// Gather all citations.
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if c, ok := node.(*ast.Citation); ok {
			for i, d := range c.Destination {
				if _, ok := seen[string(bytes.ToLower(d))]; ok {
					continue
				}
				ref := &mast.Reference{}
				ref.Anchor = d
				ref.Type = c.Type[i]

				seen[string(d)] = ref
			}
		}
		return ast.GoToNext
	})

	for _, r := range seen {
		switch r.Type {
		case ast.CitationTypeNone:
			fallthrough
		case ast.CitationTypeInformative:
			if informative == nil {
				informative = &mast.References{}
				p.Inline(informative, []byte("Normative References"))
			}
			ast.AppendChild(informative, r)
		case ast.CitationTypeNormative:
			if normative == nil {
				normative = &mast.References{}
				p.Inline(normative, []byte("Normative References"))
			}
			ast.AppendChild(normative, r)
		case ast.CitationTypeSuppressed:
			// Don't add it.
		}
	}
	return normative, informative
}
