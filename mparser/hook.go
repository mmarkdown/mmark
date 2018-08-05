package mparser

import "github.com/gomarkdown/markdown/ast"

// Hook will call both TitleHook and ReferenceHook.
func Hook(data []byte) (ast.Node, []byte, int) {
	n, b, i := TitleHook(data)
	if n != nil {
		return n, b, i
	}

	return ReferenceHook(data)
}
