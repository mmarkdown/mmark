package mparser

import (
	"bytes"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
)

func CitationToAST(doc ast.Node) ast.Node {
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

	refs := &mast.References{}
	for _, r := range seen {
		ast.AppendChild(refs, r)
	}
	return refs
}
