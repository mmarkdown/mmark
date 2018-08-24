package mast

import (
	"github.com/gomarkdown/markdown/ast"
)

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

// Some attribute helper functions.

func attributeFromNode(node ast.Node) *ast.Attribute {
	if l := node.AsLeaf(); l != nil && l.Attribute != nil {
		return l.Attribute
	}
	if c := node.AsContainer(); c != nil && c.Attribute != nil {
		return c.Attribute
	}
	return nil
}

// DeleteAttribute delete the attribute under key from a.
func DeleteAttribute(node ast.Node, key string) {
	a := attributeFromNode(node)
	if a == nil {
		return
	}

	switch key {
	case "id":
		a.ID = nil
	case "class":
		// TODO
	default:
		delete(a.Attrs, key)
	}
}

// SetAttribute sets the attribute under key to value.
func SetAttribute(node ast.Node, key string, value []byte) {
	a := attributeFromNode(node)
	if a == nil {
		return
	}
	switch key {
	case "id":
		a.ID = value
	case "class":
		// TODO
	default:
		a.Attrs[key] = value
	}
}

// Attribute return the attribute value under key.
func Attribute(node ast.Node, key string) []byte {
	a := attributeFromNode(node)
	if a == nil {
		return nil
	}
	switch key {
	case "id":
		return a.ID
	case "class":
		// TODO
	}

	return a.Attrs[key]
}
