package mparser

import (
	"sort"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/mast"
)

// IndexToIndices crawls the entire doc searching for indices, it will then return
// an ast.DocumentIndex that contains ast.Indices.
func IndexToIndices(p *parser.Parser, doc ast.Node) *mast.Indices {
	all := map[string]*mast.IndexItem{}

	// Gather all indexes.
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		switch i := node.(type) {
		case *ast.Index:
			item := string(i.Item)
			ii, ok := all[item]
			if !ok {
				it := &mast.IndexItem{}
				it.Items = []*ast.Index{i}
				all[item] = it
			} else {
				ii.Items = append(ii.Items, i)
			}
		}
		return ast.GoToNext
	})
	if len(all) == 0 {
		return nil
	}

	keys := []string{}
	for k := range all {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	indices := &mast.Indices{}
	for _, k := range keys {
		ast.AppendChild(indices, all[k])
	}

	return indices
}
