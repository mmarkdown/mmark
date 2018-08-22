package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/mast"
	"github.com/mmarkdown/mmark/mhtml"
	"github.com/mmarkdown/mmark/mparser"
	"github.com/mmarkdown/mmark/xml"
	"github.com/mmarkdown/mmark/xml2"
)

var (
	flagAst      = flag.Bool("ast", false, "print abstract syntax tree and exit")
	flagFragment = flag.Bool("fragment", false, "don't create a full document")
	flagHTML     = flag.Bool("html", false, "create HTML output")
	flagCSS      = flag.String("css", "", "link to a CSS stylesheet (only used with -html)")
	flagHead     = flag.String("head", "", "link to HTML to be included in head (only used with -html)")
	flagIndex    = flag.Bool("index", true, "generate an index at the end of the document")
	flagBib      = flag.Bool("bibliography", true, "generate a bibliography section after the back matter")
	flagTwo      = flag.Bool("2", false, "generate RFC 7749 XML")
	flagVersion  = flag.Bool("version", false, "show mmark version")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "SYNOPSIS: %s [OPTIONS] %s\n", os.Args[0], "[FILE...]")
		fmt.Println("\nOPTIONS:")
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"os.Stdin"}
	}
	if *flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	for _, fileName := range args {
		var (
			d    []byte
			err  error
			init mparser.Initial
		)
		if fileName == "os.Stdin" {
			init = mparser.NewInitial("")
			d, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Printf("Couldn't read %q: %q", fileName, err)
				continue
			}
		} else {
			init = mparser.NewInitial(fileName)
			d, err = ioutil.ReadFile(fileName)
			if err != nil {
				log.Printf("Couldn't open %q: %q", fileName, err)
				continue
			}
		}

		documentTitle := "" // hack to get document title from toml title block and then set it here.

		p := parser.NewWithExtensions(Extensions)
		p.Opts = parser.ParserOptions{
			ParserHook: func(data []byte) (ast.Node, []byte, int) {
				node, data, consumed := mparser.Hook(data)
				if t, ok := node.(*mast.Title); ok {
					documentTitle = t.TitleData.Title
				}
				return node, data, consumed
			},
			ReadIncludeFn: init.ReadInclude,
		}

		doc := markdown.Parse(d, p)
		if *flagBib {
			where := mparser.NodeBackMatter(doc)
			if where != nil {
				norm, inform := mparser.CitationToBibliography(p, doc)
				if norm != nil {
					ast.AppendChild(where, norm)
				}
				if inform != nil {
					ast.AppendChild(where, inform)
				}
			}
		}
		if *flagIndex {
			if idx := mparser.IndexToDocumentIndex(p, doc); idx != nil {
				ast.AppendChild(doc, idx)
			}
		}

		if *flagAst {
			ast.Print(os.Stdout, doc)
			fmt.Print("\n")
			return
		}

		var renderer markdown.Renderer

		if *flagHTML {
			opts := html.RendererOptions{
				// TODO(miek): make this an option.
				Comments:       [][]byte{[]byte("//"), []byte("#")},
				RenderNodeHook: mhtml.RenderHook,
				Flags:          html.CommonFlags,
			}
			if !*flagFragment {
				opts.Flags |= html.CompletePage
			}
			opts.CSS = *flagCSS
			if *flagHead != "" {
				head, err := ioutil.ReadFile(*flagHead)
				if err != nil {
					log.Printf("Couldn't open %q, error: %q", *flagHead, err)
					continue
				}
				opts.Head = head
			}
			if documentTitle != "" {
				opts.Title = documentTitle
			}

			renderer = html.NewRenderer(opts)
		} else if *flagTwo {
			opts := xml2.RendererOptions{
				Flags:    xml2.CommonFlags,
				Comments: [][]byte{[]byte("//"), []byte("#")},
			}
			if *flagFragment {
				opts.Flags |= xml2.XMLFragment
			}

			renderer = xml2.NewRenderer(opts)
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

// Extensions is exported to we can use it in tests.
var Extensions = parser.CommonExtensions | parser.HeadingIDs | parser.AutoHeadingIDs | parser.Footnotes |
	parser.Strikethrough | parser.OrderedListStart | parser.Attributes | parser.Mmark | parser.Autolink
