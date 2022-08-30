package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/v2/lang"
	"github.com/mmarkdown/mmark/v2/mast"
	"github.com/mmarkdown/mmark/v2/mparser"
	"github.com/mmarkdown/mmark/v2/render/mhtml"
	"github.com/mmarkdown/mmark/v2/render/xml"
)

var (
	testFiles = []string{
		"2100.md",
		"3514.md",
		"5841.md",
		"7511.md",
		"8341.md",
	}
)

// TestRFC3 parses the RFC in the rfc/ directory and runs xml2rfc on them to see if they parse OK.
func TestRFC3(t *testing.T) {
	for _, f := range testFiles {
		base := f[:len(f)-3]

		opts := xml.RendererOptions{
			Flags: xml.CommonFlags | xml.XMLFragment,
		}
		renderer := xml.NewRenderer(opts)
		t.Run("rfc3/"+base, func(t *testing.T) {
			err := doRenderTest(base, renderer)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

// TestHTML parses the RFC in the rfc/ directory to HTMl.
func TestHTML(t *testing.T) {
	for _, f := range testFiles {
		base := f[:len(f)-3]

		mhtmlOpts := mhtml.RendererOptions{Language: lang.New("en")}
		opts := html.RendererOptions{
			RenderNodeHook: mhtmlOpts.RenderHook,
		}
		renderer := html.NewRenderer(opts)
		t.Run("html/"+base, func(t *testing.T) {
			err := doRenderTest(base, renderer)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func doRenderTest(basename string, renderer markdown.Renderer) error {
	filename := filepath.Join("rfc", basename+".md")
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("couldn't open '%s', error: %v\n", filename, err)
	}

	init := mparser.NewInitial(filename)

	p := parser.NewWithExtensions(mparser.Extensions)
	p.Opts = parser.Options{
		ParserHook: func(data []byte) (ast.Node, []byte, int) {
			node, data, consumed := mparser.Hook(data)
			if t, ok := node.(*mast.Title); ok {
				_ = t.TitleData.Title
			}
			return node, data, consumed
		},
		ReadIncludeFn: init.ReadInclude,
	}

	doc := markdown.Parse(input, p)
	mparser.AddBibliography(doc)
	mparser.AddIndex(doc)

	rfcdata := markdown.Render(doc, renderer)

	switch renderer.(type) {
	case *xml.Renderer:
		out, err := runXML2RFC([]string{"--v3"}, rfcdata)
		if err != nil {
			return fmt.Errorf("failed to parse XML3 output for %q: %s\n%s", filename, err, out)
		}
	}
	return nil
}

func runXML2RFC(options []string, rfc []byte) ([]byte, error) {
	if _, err := exec.LookPath("xml2rfc"); err != nil {
		return nil, nil
	}
	ioutil.WriteFile("x.xml", rfc, 0600)
	defer os.Remove("x.xml")
	defer os.Remove("x.txt") // if we are lucky.
	xml2rfc := exec.Command("xml2rfc", append([]string{"x.xml"}, options...)...)
	return xml2rfc.CombinedOutput()
}
