package mast

import "github.com/gomarkdown/markdown/ast"

// References represents markdown references node.
type References struct {
	ast.Container

	HeadingID string // This might hold heading ID, if present
}

// Reference contains a single citation.
type Reference struct {
	ast.Leaf

	Anchor []byte
	Type   ast.CitationTypes
	RawXML []byte // If there is a <reference> in the doc
}
