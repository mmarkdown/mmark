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
// Its supports the following options for address.
//
// 4,5 - line numbers separated by commas
// /start/,/end/ - regexp separated by commas
// optional a prefix="" string.
func (c *Cwd) ReadInclude(p string, address []byte) []byte {
	path := c.Path(p)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Failure to read %s: %s", path, err)
		return nil
	}

	data, err = parseAddress(address, data)
	if err != nil {
		log.Printf("Failure to read %s: %s", path, err)
		return nil
	}
	if data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}

	c.Update(p)
	return data
}
