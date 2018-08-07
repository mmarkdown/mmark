package mparser

import (
	"io/ioutil"
	"log"

	"github.com/gomarkdown/markdown/ast"
)

// Hook will call both TitleHook and ReferenceHook.
func Hook(data []byte) (ast.Node, []byte, int) {
	n, b, i := TitleHook(data)
	if n != nil {
		return n, b, i
	}

	return ReferenceHook(data)
}

// ReadInclude is the hook to read includes.
// Its supports:
//
// 4,5 - line numbers separated by commas
// /start/,/end/ - regexp separated by commas
//
// for address.
func ReadInclude(path string, address []byte) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Failure to read %s: %s", path, err)
		return nil
	}

	data, err = parseAddress(address, data)

	return data
}
