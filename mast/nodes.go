package mast

import "github.com/gomarkdown/markdown/ast"

// some extra functions for manipulation the AST

// MoveChilderen moves the children from a to b *and* make the parent of each point to b.
// Any children of b are obliterated.
func MoveChildren(a, b ast.Node) {
	a.SetChildren(b.GetChildren())
	b.SetChildren(nil)

	for _, child := range a.GetChildren() {
		child.SetParent(a)
	}
}
