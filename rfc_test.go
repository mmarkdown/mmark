// +build xml2rfc

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/mast"
	"github.com/mmarkdown/mmark/mparser"
	"github.com/mmarkdown/mmark/xml"
	"github.com/mmarkdown/mmark/xml2"
)

// TestRFC3 parses the RFC in the rfc/ directory and runs xml2rfc on them to see if they parse OK.
func TestRFC3(t *testing.T) { // currently broken because of xml2rfc --debug foo
	dir := "rfc"
	testFiles := []string{
		"2100.md",
		"3514.md",
		"7511.md",
	}

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
func TestRFC2(t *testing.T) { // currently broken because of xml2rfc --debug foo
	dir := "rfc"
	testFiles := []string{
		"2100.md",
		"3514.md",
		"7511.md",
	}

	for _, f := range testFiles {
		base := f[:len(f)-3]

		opts2 := xml2.RendererOptions{
			Flags: xml2.CommonFlags | xml2.XMLFragment,
		}
		renderer2 := xml2.NewRenderer(opts2)
		doRenderTest(t, dir, base, renderer2)
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

	p := parser.NewWithExtensions(Extensions)
	p.Opts = parser.ParserOptions{
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
	addBibliography(doc)
	addIndex(doc)

	rfcdata := markdown.Render(doc, renderer)

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
