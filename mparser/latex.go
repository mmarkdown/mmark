package mparser

import (
	"unicode"

	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
)

func LatexHook(data []byte) (ast.Node, []byte, int) {
	if len(data) < 4 {
		return nil, nil, 0
	}
	i := 0
	if data[i] != '\\' {
		return nil, nil, 0
	}
	i++
	for i < len(data) && data[i] != '{' {
		c := data[i]
		// chars in between need to be a-z or A-Z or 0-9
		if !unicode.IsLetter(rune(c)) && !unicode.IsNumber(rune(c)) {
			return nil, nil, 0
		}
		i++
	}
	if i == len(data) || i == 1 {
		return nil, nil, 0
	}

	// find first } - this isn't perfect but works for now.
	for i < len(data) && data[i] != '}' {
		i++
	}
	if i == len(data) {
		return nil, nil, 0
	}
	node := &mast.LatexSpan{}
	node.Content = data[:i+1]
	return node, nil, i + 1
}
