package main

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gopherjs/gopherjs/js"
	"github.com/mmarkdown/mmark/mast"
	"github.com/mmarkdown/mmark/mparser"
	"github.com/mmarkdown/mmark/render/mhtml"
)

var standardComments = [][]byte{[]byte("//"), []byte("#")}

func main() {
	js.Global.Set("mmark", map[string]interface{}{
		"NewDocument": NewDocument,
	})
}

type Document struct {
	title string
	root  ast.Node
}

func NewDocument(data string) *js.Object {
	init := mparser.NewInitial("")

	doc := &Document{}

	p := parser.NewWithExtensions(mparser.Extensions)
	parserFlags := parser.FlagsNone
	p.Opts = parser.Options{
		ParserHook: func(data []byte) (ast.Node, []byte, int) {
			node, data, consumed := mparser.Hook(data)
			if t, ok := node.(*mast.Title); ok {
				if !t.IsTriggerDash() {
					doc.title = t.TitleData.Title
				}
			}
			return node, data, consumed
		},
		ReadIncludeFn: init.ReadInclude,
		Flags:         parserFlags,
	}

	doc.root = markdown.Parse([]byte(data), p)
	return js.MakeWrapper(doc)
}

func (doc *Document) ToHTML() string {
	opts := html.RendererOptions{
		Comments:       [][]byte{[]byte("//"), []byte("#")},
		RenderNodeHook: mhtml.RenderHook,
		Flags:          html.CommonFlags | html.FootnoteNoHRTag | html.FootnoteReturnLinks | html.CompletePage,
		Title:          doc.title,
	}

	renderer := html.NewRenderer(opts)
	x := markdown.Render(doc.root, renderer)
	return string(x)
}
