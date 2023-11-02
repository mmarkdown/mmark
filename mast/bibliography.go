package mast

import (
	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/v2/mast/reference"
)

// BibliographyWrapper represents markdown bibliography section wrapper.
type BibliographyWrapper struct {
	ast.Container
}

// Bibliography represents markdown bibliography node.
type Bibliography struct {
	ast.Container

	Type ast.CitationTypes
}

// BibliographyItem contains a single bibliography item.
type BibliographyItem struct {
	ast.Leaf

	Anchor []byte
	Type   ast.CitationTypes

	Reference      *reference.Reference // parsed reference XML
	ReferenceGroup []byte               // raw, unparsed reference group  XML
}
