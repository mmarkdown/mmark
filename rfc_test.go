package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/miekg/markdown/xml2"
	"github.com/mmarkdown/mmark/mparser"
	"github.com/mmarkdown/mmark/xml"
)

// testRFC parses the RFC in the rfc/ directory and runs xml2rfc on them to see if they parse OK.
func testRFC(t *testing.T) { // currently broken because of xml2rfc --debug foo
	dir := "rfc"
	testFiles := []string{
		"2100.md",
		"3514.md",
		// "7511.md",
	}

	for _, f := range testFiles {
		base := f[:len(f)-3]

		opts := xml.RendererOptions{
			Flags: xml.CommonFlags | xml.XMLFragment,
		}
		renderer := xml.NewRenderer(opts)
		doRenderTest(t, dir, base, renderer)

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

	p := parser.NewWithExtensions(Extensions)
	init := mparser.NewInitial(filename)
	p.Opts = parser.ParserOptions{
		ParserHook:    mparser.TitleHook,
		ReadIncludeFn: init.ReadInclude,
	}

	rfcdata := markdown.ToHTML(input, p, renderer)
	switch renderer.(type) {
	case *xml.Renderer:
		out, err := runXml2Rfc([]string{"--v3"}, rfcdata)
		if err != nil {
			t.Errorf("failed to parse XML3 output for %q: %s\n%s", filename, err, out)
		}
	case *xml2.Renderer:
		out, err := runXml2Rfc(nil, rfcdata)
		if err != nil {
			t.Errorf("failed to parse XML2 output for %q: %s\n%s", filename, err, out)
		}
	}
}

func runXml2Rfc(options []string, rfc []byte) ([]byte, error) {
	ioutil.WriteFile("x.xml", rfc, 0600)
	defer os.Remove("x.xml")
	defer os.Remove("x.txt") // if we are lucky.
	xml2rfc := exec.Command("xml2rfc", append([]string{"x.xml"}, options...)...)
	return xml2rfc.CombinedOutput()
}
