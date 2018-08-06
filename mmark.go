package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/miekg/markdown/xml"
	"github.com/mmarkdown/mmark/mparser"
)

// Usage: mmark <markdown-file>

var flagAst = flag.Bool("ast", false, "print abstract syntax tree and exit.")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS] %s\n", os.Args[0], "FILE...")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	for _, fileName := range flag.Args() {
		d, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't open '%s', error: '%s'\n", fileName, err)
			continue
		}

		ext := parser.CommonExtensions | parser.HeadingIDs | parser.AutoHeadingIDs | parser.Footnotes |
			parser.OrderedListStart | parser.Attributes | parser.Mmark

		p := parser.NewWithExtensions(ext)
		p.Opts = parser.ParserOptions{ParserHook: mparser.Hook}

		opts := xml.RendererOptions{
			Flags: xml.CommonFlags,
		}

		doc := markdown.Parse(d, p)
		norm, inform := mparser.CitationToReferences(p, doc)
		if norm != nil {
			ast.AppendChild(doc, norm)
		}
		if inform != nil {
			ast.AppendChild(doc, inform)
		}

		if *flagAst {
			ast.Print(os.Stdout, doc)
			fmt.Print("\n")
			return
		}

		renderer := xml.NewRenderer(opts)
		x := markdown.Render(doc, renderer)
		fmt.Println(string(x))
	}
}
