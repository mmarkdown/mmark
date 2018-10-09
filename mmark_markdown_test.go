package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	mmarkdown "github.com/mmarkdown/mmark/markdown"
	"github.com/mmarkdown/mmark/mparser"
)

func TestMmarkMarkdown(t *testing.T) {
	dir := "testdata/markdown"
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
		opts := mmarkdown.RendererOptions{}

		renderer := mmarkdown.NewRenderer(opts)

		doTestMarkdown(t, dir, base, renderer)
	}
}

func doTestMarkdown(t *testing.T, dir, basename string, renderer markdown.Renderer) {
	filename := filepath.Join(dir, basename+".md")
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("couldn't open '%s', error: %v\n", filename, err)
		return
	}

	filename = filepath.Join(dir, basename+".fmt")
	expected, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("couldn't open '%s', error: %v\n", filename, err)
	}
	expected = bytes.TrimSpace(expected)

	p := parser.NewWithExtensions(Extensions) // no includes

	p.Opts = parser.ParserOptions{
		ParserHook: mparser.TitleHook,
	}

	doc := markdown.Parse(input, p)
	actual := markdown.Render(doc, renderer)
	actual = bytes.TrimSpace(actual)
	if bytes.Compare(actual, expected) != 0 {
		t.Errorf("\n    [%#v]\nExpected[%s]\nActual  [%s]",
			basename+".md", expected, actual)
	}
}
