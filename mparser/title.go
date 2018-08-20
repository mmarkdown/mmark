package mparser

import (
	"log"

	"github.com/BurntSushi/toml"
	"github.com/gomarkdown/markdown/ast"
	"github.com/mmarkdown/mmark/mast"
)

// TitleHook will parse a title and returns it.
func TitleHook(data []byte) (ast.Node, []byte, int) {
	i := 0
	if len(data) < 3 {
		return nil, nil, 0
	}
	if data[i] != '%' || data[i+1] != '%' || data[i+2] != '%' {
		return nil, nil, 0
	}

	i += 3
	beg := i
	// search for end.
	for i < len(data) {
		if data[i] == '%' || data[i+1] == '%' || data[i+2] == '%' {
			break
		}
		i++
	}

	node := mast.NewTitle()

	if _, err := toml.Decode(string(data[beg:i+1]), node.TitleData); err != nil {
		log.Printf("Failure parsing title block: %s", err)
	}

	return node, nil, i + 5
}
