package mast

import "github.com/gomarkdown/markdown/ast"

// Indices represents markdown document index node.
type Indices struct {
	ast.Container

	HeadingID string // This might hold heading ID, if present
}

// IndexItem contains an index for the indices section. It has all the ID
// for the main item of an index.
type IndexItem struct {
	ast.Container

	*ast.Index
	Primary int      // index into IDs to signal the primary item
	IDs     []string // all the of the item's ID that have the item in common
}

// IndexSubItem contains an sub item index for the indices section. It has all the ID
// for the sub item's of an index.
type IndexSubItem struct {
	ast.Leaf

	*ast.Index
	IDs []string // all the of the sub item's ID that have the item in common
}

// IndexLetter has the Letter of this index item.
type IndexLetter struct {
	ast.Leaf
}
