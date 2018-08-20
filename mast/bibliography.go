package mast

import (
	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast/reference"
)

// Bibliography represents markdown bibliography node.
type Bibliography struct {
	ast.Container

	HeadingID string // This might hold heading ID, if present
}

// BibliographyItem contains a single bibliography item.
type BibliographyItem struct {
	ast.Leaf

	Anchor []byte
	Type   ast.CitationTypes

	Raw       []byte              // raw reference XML
	Reference reference.Reference // parsed reference XML
}
