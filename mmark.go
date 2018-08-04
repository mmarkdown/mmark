package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/miekg/markdown/xml3"
	"github.com/mmarkdown/mmark/mparser"
)

// Usage: mmark <markdown-file>

func usageAndExit() {
	fmt.Printf("Usage: mmark <markdown-file>\n")
	os.Exit(1)
}

func main() {
	nFiles := len(os.Args) - 1
	if nFiles < 1 {
		usageAndExit()
	}
	for i := 0; i < nFiles; i++ {
		fileName := os.Args[i+1]
		d, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't open '%s', error: '%s'\n", fileName, err)
			continue
		}

		ext := parser.CommonExtensions | parser.HeadingIDs | parser.AutoHeadingIDs | parser.Footnotes |
			parser.OrderedListStart | parser.Attributes | parser.Mmark

		p := parser.NewWithExtensions(ext)
		p.Opts = parser.ParserOptions{ParserHook: mparser.TitleHook}

		doc := markdown.Parse(d, p)
		fmt.Printf("Ast of file '%s':\n", fileName)
		ast.Print(os.Stdout, doc)
		fmt.Print("\n")

		p = parser.NewWithExtensions(ext)
		p.Opts = parser.ParserOptions{
			ParserHook: mparser.TitleHook,
		}

		opts := xml3.RendererOptions{
			Flags: xml3.CommonFlags,
		}
		renderer := xml3.NewRenderer(opts)
		xml := markdown.ToHTML(d, p, renderer)
		fmt.Println(string(xml))
	}
}
