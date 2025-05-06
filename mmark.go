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
	"github.com/mmarkdown/mmark/v2/lang"
	"github.com/mmarkdown/mmark/v2/mast"
	"github.com/mmarkdown/mmark/v2/mparser"
	"github.com/mmarkdown/mmark/v2/render/man"
	"github.com/mmarkdown/mmark/v2/render/mhtml"
	"github.com/mmarkdown/mmark/v2/render/text"
	"github.com/mmarkdown/mmark/v2/render/xml"
)

var (
	flagCSS       = flag.String("css", "", "link to a CSS stylesheet (only used with -html)")
	flagHead      = flag.String("head", "", "link to HTML to be included in head (only used with -html)")
	flagAst       = flag.Bool("ast", false, "print abstract syntax tree and exit")
	flagBib       = flag.Bool("bibliography", true, "generate a bibliography section after the back matter")
	flagFragment  = flag.Bool("fragment", false, "don't create a full document")
	flagHTML      = flag.Bool("html", false, "create HTML output")
	flagIndex     = flag.Bool("index", true, "generate an index at the end of the document")
	flagMan       = flag.Bool("man", false, "generate manual pages (nroff)")
	flagText      = flag.Bool("text", false, "generate text for ANSI escapes")
	flagUnsafe    = flag.Bool("unsafe", false, "allow unsafe includes")
	flagIntraEmph = flag.Bool("intra-emphasis", false, "interpret camel_case_value as emphasizing \"case\" (legacy behavior)")
	flagVersion   = flag.Bool("version", false, "show mmark version")
	flagUnicode   = flag.Bool("unicode", true, "from xml2rfc 3.16 onwards unicode is allowed in <t>")
	flagWidth     = flag.Int("width", 100, "text width when generating markdown")
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

		d = markdown.NormalizeNewlines(d)

		if *flagUnsafe {
			init.Flags |= mparser.UnsafeInclude
		}

		if !*flagIntraEmph {
			mparser.Extensions |= parser.NoIntraEmphasis
		}

		p := parser.NewWithExtensions(mparser.Extensions)
		parserFlags := parser.FlagsNone
		documentTitle := ""      // hack to get document title from toml title block and then set it here.
		documentLanguage := "en" // get document language from title block if it is set.
		if !*flagHTML && !*flagMan {
			parserFlags |= parser.SkipFootnoteList // both xml formats don't deal with footnotes well.
		}
		p.Opts = parser.Options{
			ParserHook: func(data []byte) (ast.Node, []byte, int) {
				node, data, consumed := mparser.Hook(data)
				if t, ok := node.(*mast.Title); ok {
					documentTitle = t.TitleData.Title
					documentLanguage = t.TitleData.Language
				}
				return node, data, consumed
			},
			ReadIncludeFn: init.ReadInclude,
			Flags:         parserFlags,
		}

		doc := markdown.Parse(d, p)
		if *flagMan {
			title := false
			// If there isn't a title block the resulting manual page does not start
			// with .TH, this messes up the entire rendering. Walk to AST to check for
			// a title block, and if none is found inject an empty one.
			ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
				if _, ok := node.(*mast.Title); ok {
					title = true
					return ast.Terminate
				}
				return ast.GoToNext
			})
			if !title {
				t := &mast.Title{TitleData: &mast.TitleData{Title: "User Commands 1"}}
				c := doc.GetChildren()
				newc := append([]ast.Node{t}, c...)
				doc.SetChildren(newc) // t must be the first element.
			} else {
				ast.AppendChild(doc, &mast.Authors{})
			}

		}
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
			mhtmlOpts := mhtml.RendererOptions{
				Language: lang.New(documentLanguage),
			}
			opts := html.RendererOptions{
				Comments:       [][]byte{[]byte("//"), []byte("#")}, // TODO(miek): make this an option.
				RenderNodeHook: mhtmlOpts.RenderHook,
				Flags:          html.CommonFlags | html.FootnoteNoHRTag | html.FootnoteReturnLinks,
				Generator:      `  <meta name="GENERATOR" content="github.com/mmarkdown/mmark Mmark Markdown Processor - mmark.miek.nl`,
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
		case *flagMan:
			opts := man.RendererOptions{
				Comments: [][]byte{[]byte("//"), []byte("#")},
				Language: lang.New(documentLanguage),
			}
			if *flagFragment {
				opts.Flags |= man.ManFragment
			}
			renderer = man.NewRenderer(opts)
		case *flagText:
			opts := text.RendererOptions{TextWidth: *flagWidth}
			renderer = text.NewRenderer(opts)
		default:
			opts := xml.RendererOptions{
				Flags:    xml.CommonFlags,
				Comments: [][]byte{[]byte("//"), []byte("#")},
				Language: lang.New(documentLanguage),
			}
			if *flagFragment {
				opts.Flags |= xml.XMLFragment
			}
			if *flagUnicode {
				opts.Flags |= xml.AllowUnicode
			}

			renderer = xml.NewRenderer(opts)
		}

		x := markdown.Render(doc, renderer)

		fmt.Println(string(x))
	}
}
