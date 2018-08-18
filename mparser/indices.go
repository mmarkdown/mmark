package mparser

import (
	"bytes"
	"sort"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/mast"
)

// IndexToIndices crawls the entire doc searching for indices, it will then return
// an mast.Indices that contains mast.IndexItems that group all indices with the same
// item.
func IndexToIndices(p *parser.Parser, doc ast.Node) *mast.Indices {
	main := map[string]*mast.IndexItem{} // main item -> Index Items

	// Gather all indexes.
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		switch i := node.(type) {
		case *ast.Index:
			item := string(i.Item)

			if _, ok := main[item]; !ok {
				main[item] = &mast.IndexItem{Index: i}
			}
			// only the main item
			if i.Subitem == nil {
				main[item].IDs = append(main[item].IDs, i.ID)
				return ast.GoToNext
			}
			// check if we already have a child with the subitem and then just add the ID
			for _, c := range main[item].GetChildren() {
				si := c.(*mast.IndexSubItem)
				if bytes.Compare(si.Subitem, i.Subitem) == 0 {
					si.IDs = append(si.IDs, i.ID)
					return ast.GoToNext
				}
			}

			sub := &mast.IndexSubItem{Index: i, IDs: []string{i.ID}}
			ast.AppendChild(main[item], sub)
		}
		return ast.GoToNext
	})
	if len(main) == 0 {
		return nil
	}

	keys := []string{}
	for k := range main {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	indices := &mast.Indices{}
	prevLetter := ""
	for _, k := range keys {
		letter := string(k[0])
		if letter != prevLetter {
			il := &mast.IndexLetter{}
			il.Literal = []byte(letter)
			ast.AppendChild(indices, il)
		}
		ast.AppendChild(indices, main[k])
		prevLetter = letter
	}

	return indices
}
