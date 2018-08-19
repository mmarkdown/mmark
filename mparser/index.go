package mparser

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/mast"
)

// IndexToDocIndices crawls the entire doc searching for indices, it will then return
// an mast.DocumentIndex that contains a tree:
//
// IndexLetter
// - IndexItem
//   - IndexLink
//   - IndexSubItem
//     - IndexLink
//     - IndexLink
//
// Which can then be rendered by the renderer.
func IndexToDocumentIndex(p *parser.Parser, doc ast.Node) *mast.DocumentIndex {
	main := map[string]*mast.IndexItem{}
	subitem := map[string][]*mast.IndexSubItem{} // gather these so we can add them in one swoop at the end

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
				ast.AppendChild(main[item], newLink(i.ID, len(main[item].GetChildren()), i.Primary))
				return ast.GoToNext
			}
			// check if we already have a child with the subitem and then just add the link
			for _, sub := range subitem[item] {
				if bytes.Compare(sub.Subitem, i.Subitem) == 0 {
					ast.AppendChild(sub, newLink(i.ID, len(sub.GetChildren()), i.Primary))
					return ast.GoToNext
				}
			}

			sub := &mast.IndexSubItem{Index: i}
			ast.AppendChild(sub, newLink(i.ID, len(subitem[item]), i.Primary))
			subitem[item] = append(subitem[item], sub)
		}
		return ast.GoToNext
	})
	if len(main) == 0 {
		return nil
	}

	// Now add a subitem children to the correct main item.
	for k, sub := range subitem {
		// sort sub here ideally
		for j := range sub {
			ast.AppendChild(main[k], sub[j])
		}
	}

	keys := []string{}
	for k := range main {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	docIndex := &mast.DocumentIndex{}
	prevLetter := ""
	for _, k := range keys {
		letter := string(k[0])
		if letter != prevLetter {
			il := &mast.IndexLetter{}
			il.Literal = []byte(letter)
			ast.AppendChild(docIndex, il)
		}
		ast.AppendChild(docIndex, main[k])
		prevLetter = letter
	}

	return docIndex
}

func newLink(id string, number int, primary bool) *mast.IndexLink {
	link := &ast.Link{Destination: []byte(id)}
	il := &mast.IndexLink{Link: link, Primary: primary}
	il.Literal = []byte(fmt.Sprintf("%d", number))
	return il
}
