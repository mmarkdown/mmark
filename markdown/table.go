package markdown

import (
	"bytes"

	"github.com/gomarkdown/markdown/ast"
)

func (r *Renderer) tableColWidth(tab *ast.Table) []int {
	cells := 0
	ast.WalkFunc(tab, func(node ast.Node, entering bool) ast.WalkStatus {
		switch node := node.(type) {
		case *ast.TableRow:
			cells = len(node.GetChildren())
			break
		}
		return ast.GoToNext
	})

	width := make([]int, cells)

	ast.WalkFunc(tab, func(node ast.Node, entering bool) ast.WalkStatus {
		switch node := node.(type) {
		case *ast.TableRow:
			for col, cell := range node.GetChildren() {
				buf := &bytes.Buffer{}
				ast.WalkFunc(cell, func(node1 ast.Node, entering bool) ast.WalkStatus {
					r.RenderNode(buf, node1, entering)
					return ast.GoToNext
				})
				if l := buf.Len(); l > width[col] {
					width[col] = l + 1 // space in beginning or end
				}
			}
		}
		return ast.GoToNext
	})
	return width
}
