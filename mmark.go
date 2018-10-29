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
	"github.com/mmarkdown/mmark/mparser"
	mmarkout "github.com/mmarkdown/mmark/render/markdown"
	"github.com/mmarkdown/mmark/render/mhtml"
	"github.com/mmarkdown/mmark/render/xml"
	"github.com/mmarkdown/mmark/render/xml2"
)

var (
	flagCSS      = flag.String("css", "", "link to a CSS stylesheet (only used with -html)")
	flagHead     = flag.String("head", "", "link to HTML to be included in head (only used with -html)")
	flagAst      = flag.Bool("ast", false, "print abstract syntax tree and exit")
	flagBib      = flag.Bool("bibliography", true, "generate a bibliography section after the back matter")
	flagFragment = flag.Bool("fragment", false, "don't create a full document")
	flagHTML     = flag.Bool("html", false, "create HTML output")
	flagIndex    = flag.Bool("index", true, "generate an index at the end of the document")
	flagTwo      = flag.Bool("2", false, "generate RFC 7749 XML")
	flagMarkdown = flag.Bool("markdown", false, "generate markdown (experimental)")
	flagWrite    = flag.Bool("w", false, "write to source file when generating markdown")
	flagWidth    = flag.Int("width", 100, "text width when generating markdown")
	flagUnsafe   = flag.Bool("unsafe", false, "allow unsafe includes")
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
		if *flagUnsafe {
			init.Flags |= mparser.UnsafeInclude
		}

		if *flagMarkdown {
			mparser.Extensions &^= parser.Includes
		}

		p := parser.NewWithExtensions(mparser.Extensions)
		parserFlags := parser.FlagsNone
		documentTitle := "" // hack to get document title from toml title block and then set it here.
		if !*flagHTML {
			parserFlags |= parser.SkipFootnoteList // both xml formats don't deal with footnotes well.
		}
		p.Opts = parser.ParserOptions{
			ParserHook: func(data []byte) (ast.Node, []byte, int) {
				node, data, consumed := mparser.Hook(data)
				if t, ok := node.(*mast.Title); ok {
					if !t.IsTriggerDash() {
						documentTitle = t.TitleData.Title
					}
				}
				return node, data, consumed
			},
			ReadIncludeFn: init.ReadInclude,
			Flags:         parserFlags,
		}

		doc := markdown.Parse(d, p)
		if *flagBib {
			mparser.AddBibliography(doc)
		}
		if *flagIndex {
			mparser.AddIndex(doc)
		}

		if *flagAst {
			ast.Print(os.Stdout, doc)
			fmt.Print("\n")
			return
		}

		var renderer markdown.Renderer

		switch {
		case *flagHTML:
			opts := html.RendererOptions{
				// TODO(miek): make this an option.
				Comments:       [][]byte{[]byte("//"), []byte("#")},
				RenderNodeHook: mhtml.RenderHook,
				Flags:          html.CommonFlags | html.FootnoteNoHRTag | html.FootnoteReturnLinks,
				Generator:      `  <meta name="GENERATOR" content="github.com/mmarkdown/mmark Mmark Markdown Processor - mmark.nl`,
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
		case *flagTwo:
			opts := xml2.RendererOptions{
				Flags:    xml2.CommonFlags,
				Comments: [][]byte{[]byte("//"), []byte("#")},
			}
			if *flagFragment {
				opts.Flags |= xml2.XMLFragment
			}

			renderer = xml2.NewRenderer(opts)
		case *flagMarkdown:
			opts := mmarkout.RendererOptions{TextWidth: *flagWidth}
			renderer = mmarkout.NewRenderer(opts)
		default:
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
		if *flagMarkdown && *flagWrite && fileName != "os.Stdin" {
			ioutil.WriteFile(fileName, x, 0600)
			continue
		}
		if *flagMarkdown {
			fmt.Print(string(x))
			continue
		}

		fmt.Println(string(x))
	}
}
