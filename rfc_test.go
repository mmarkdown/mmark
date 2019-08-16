package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/mmarkdown/mmark/lang"
	"github.com/mmarkdown/mmark/mast"
	"github.com/mmarkdown/mmark/mparser"
	mmarkdown "github.com/mmarkdown/mmark/render/markdown"
	"github.com/mmarkdown/mmark/render/mhtml"
	"github.com/mmarkdown/mmark/render/xml"
	"github.com/mmarkdown/mmark/render/xml2"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var (
	testFiles = []string{
		"2100.md",
		"3514.md",
		"5841.md",
		"7511.md",
		"8341.md",
	}
	dir      = "rfc"
	doXMLRFC = false // xml2rfc is broken as we need a very new version of it (in travis)
)

// TestRFC3 parses the RFC in the rfc/ directory and runs xml2rfc on them to see if they parse OK.
func TestRFC3(t *testing.T) {
	for _, f := range testFiles {
		base := f[:len(f)-3]

		opts := xml.RendererOptions{
			Flags: xml.CommonFlags | xml.XMLFragment,
		}
		renderer := xml.NewRenderer(opts)
		doRenderTest(t, dir, base, renderer)
	}
}

// TestRFC2 parses the RFC in the rfc/ directory and runs xml2rfc on them to see if they parse OK.
func TestRFC2(t *testing.T) {
	for _, f := range testFiles {
		base := f[:len(f)-3]

		opts := xml2.RendererOptions{
			Flags: xml2.CommonFlags | xml2.XMLFragment,
		}
		renderer := xml2.NewRenderer(opts)
		doRenderTest(t, dir, base, renderer)
	}
}

// TestHTML parses the RFC in the rfc/ directory to HTMl.
func TestHTML(t *testing.T) {
	for _, f := range testFiles {
		base := f[:len(f)-3]

		mhtmlOpts := mhtml.RenderOptions{Language: lang.New("en")}
		opts := html.RendererOptions{
			RenderNodeHook: mhtmlOpts.RenderHook,
		}
		renderer := html.NewRenderer(opts)
		doRenderTest(t, dir, base, renderer)
	}
}

// TestMarkdown parses the RFC in the rfc/ directory to markdown.
func TestMarkdown(t *testing.T) {
	for _, f := range testFiles {
		base := f[:len(f)-3]

		opts := mmarkdown.RendererOptions{}
		renderer := mmarkdown.NewRenderer(opts)
		doRenderTest(t, dir, base, renderer)
	}
}

func doRenderTest(t *testing.T, dir, basename string, renderer markdown.Renderer) {
	filename := filepath.Join(dir, basename+".md")
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("couldn't open '%s', error: %v\n", filename, err)
		return
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
	if !doXMLRFC {
		return
	}

	switch renderer.(type) {
	case *xml.Renderer:
		out, err := runXML2RFC([]string{"--v3"}, rfcdata)
		if err != nil {
			t.Errorf("failed to parse XML3 output for %q: %s\n%s", filename, err, out)
		} else {
			t.Logf("successfully parsed %s for XML3 output:\n%s", filename, out)
		}
	case *xml2.Renderer:
		out, err := runXML2RFC(nil, rfcdata)
		if err != nil {
			t.Errorf("failed to parse XML2 output for %q: %s\n%s", filename, err, out)
		} else {
			t.Logf("successfully parsed %s for XML2 output:\n%s", filename, out)
		}
	}
}

func runXML2RFC(options []string, rfc []byte) ([]byte, error) {
	ioutil.WriteFile("x.xml", rfc, 0600)
	defer os.Remove("x.xml")
	defer os.Remove("x.txt") // if we are lucky.
	xml2rfc := exec.Command("xml2rfc", append([]string{"x.xml"}, options...)...)
	return xml2rfc.CombinedOutput()
}
