package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark"
)

// This prints AST of parsed markdown document.
// Usage: printast <markdown-file>

func usageAndExit() {
	fmt.Printf("Usage: printast <markdown-file>\n")
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

		p := parser.New()
		p.Opts = parser.ParserOptions{ParserHook: mmark.TitleHook}

		doc := markdown.Parse(d, p)
		fmt.Printf("Ast of file '%s':\n", fileName)
		ast.Print(os.Stdout, doc)
		fmt.Print("\n")

		p = parser.New()
		p.Opts = parser.ParserOptions{ParserHook: mmark.TitleHook}
		opts := html.RendererOptions{
			Flags:          html.CommonFlags,
			RenderNodeHook: mmark.RenderHookHTML,
		}
		renderer := html.NewRenderer(opts)
		html := markdown.ToHTML(d, p, renderer)
		fmt.Println(string(html))
	}
}
