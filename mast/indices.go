package mast

import "github.com/gomarkdown/markdown/ast"

// Indices represents markdown document index node.
type Indices struct {
	ast.Container

	HeadingID string // This might hold heading ID, if present
}

// IndexItem contains a single index for the indices section.
type IndexItem struct {
	ast.Leaf

	// map from the main item to index.
	Items []*ast.Index
}
