package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/mhtml"
	"github.com/mmarkdown/mmark/mparser"
	"github.com/mmarkdown/mmark/xml"
)

// Usage: mmark <markdown-file>

var (
	flagAst       = flag.Bool("ast", false, "print abstract syntax tree and exit")
	flagFragment  = flag.Bool("fragment", false, "don't create a full document")
	flagHTML      = flag.Bool("html", false, "create HTML output")
	flagCss       = flag.String("css", "", "link to a CSS stylesheet")
	flagHead      = flag.String("head", "", "link to HTML to be included in head")
	flagIndex     = flag.Bool("index", true, "generate an index at the end of the document")
	flagReference = flag.Bool("reference", true, "generate a references section at the end of the document")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS] %s\n", os.Args[0], "FILE...")
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"os.Stdin"}
	}

	for _, fileName := range args {
		cwd := mparser.NewCwd()

		var d []byte
		var err error
		if fileName == "os.Stdin" {
			d, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Couldn't read '%s', error: '%s'\n", fileName, err)
				continue
			}
		} else {
			d, err = ioutil.ReadFile(fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Couldn't open '%s', error: '%s'\n", fileName, err)
				continue
			}
			cwd.Update(fileName)
		}

		ext := parser.CommonExtensions | parser.HeadingIDs | parser.AutoHeadingIDs | parser.Footnotes |
			parser.Strikethrough | parser.OrderedListStart | parser.Attributes | parser.Mmark

		p := parser.NewWithExtensions(ext)
		p.Opts = parser.ParserOptions{
			ParserHook:    mparser.Hook,
			ReadIncludeFn: cwd.ReadInclude,
		}

		doc := markdown.Parse(d, p)
		if *flagReference {
			norm, inform := mparser.CitationToReferences(p, doc)
			if norm != nil {
				ast.AppendChild(doc, norm)
			}
			if inform != nil {
				ast.AppendChild(doc, inform)
			}
		}
		if *flagIndex {
			// index stuff
		}

		if *flagAst {
			ast.Print(os.Stdout, doc)
			fmt.Print("\n")
			return
		}

		var renderer markdown.Renderer

		if *flagHTML {
			opts := html.RendererOptions{
				// TODO(miek): makes this an option.
				Comments:       [][]byte{[]byte("//"), []byte("#")},
				RenderNodeHook: mhtml.RenderHook,
			}
			if !*flagFragment {
				opts.Flags |= html.CompletePage
			}
			opts.CSS = *flagCss
			if *flagHead != "" {
				head, err := ioutil.ReadFile(*flagHead)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Couldn't open '%s', error: '%s'\n", *flagHead, err)
					continue
				}
				opts.Head = head
			}

			renderer = html.NewRenderer(opts)
		} else {
			opts := xml.RendererOptions{
				Flags:    xml.CommonFlags,
				Comments: [][]byte{[]byte("//"), []byte("#")},
			}
			if *flagFragment {
				opts.Flags |= xml.XMLFragment
			}

			renderer = xml.NewRenderer(opts)
		}

		x := markdown.Render(doc, renderer)
		fmt.Println(string(x))
	}
}
