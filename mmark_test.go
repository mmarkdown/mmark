package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/mparser"
	"github.com/mmarkdown/mmark/render/xml"
	"github.com/mmarkdown/mmark/render/xml2"
)

func TestMmarkXML(t *testing.T) {
	dir := "testdata"
	testFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf("could not read %s: %q", dir, err)
	}
	for _, f := range testFiles {
		if f.IsDir() {
			continue
		}

		if filepath.Ext(f.Name()) != ".md" {
			continue
		}
		base := f.Name()[:len(f.Name())-3]
		opts := xml.RendererOptions{
			Flags: xml.CommonFlags | xml.XMLFragment,
		}
		renderer := xml.NewRenderer(opts)

		doTest(t, dir, base, renderer)
	}
}

func TestMmarkXML2(t *testing.T) {
	dir := "testdata/2"
	testFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf("could not read %s: %q", dir, err)
	}
	for _, f := range testFiles {
		if f.IsDir() {
			continue
		}

		if filepath.Ext(f.Name()) != ".md" {
			continue
		}
		base := f.Name()[:len(f.Name())-3]
		opts := xml2.RendererOptions{
			Flags: xml2.CommonFlags | xml2.XMLFragment,
		}
		renderer := xml2.NewRenderer(opts)

		doTest(t, dir, base, renderer)
	}
}

func doTest(t *testing.T, dir, basename string, renderer markdown.Renderer) {
	filename := filepath.Join(dir, basename+".md")
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("couldn't open '%s', error: %v\n", filename, err)
		return
	}

	filename = filepath.Join(dir, basename+".xml")
	expected, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("couldn't open '%s', error: %v\n", filename, err)
	}
	expected = bytes.TrimSpace(expected)

	p := parser.NewWithExtensions(mparser.Extensions)

	init := mparser.NewInitial(filename)
	p.Opts = parser.ParserOptions{
		ParserHook:    mparser.TitleHook,
		ReadIncludeFn: init.ReadInclude,
	}

	actual := markdown.ToHTML(input, p, renderer)
	actual = bytes.TrimSpace(actual)
	if bytes.Compare(actual, expected) != 0 {
		t.Errorf("\n    [%#v]\nExpected[%s]\nActual  [%s]",
			basename+".md", expected, actual)
	}
}
