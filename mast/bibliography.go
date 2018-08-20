package mast

import "github.com/gomarkdown/markdown/ast"

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
	RawXML []byte // If there is a <reference> in the doc
}
