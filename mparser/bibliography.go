package mparser

import (
	"bytes"
	"encoding/xml"
	"log"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mmarkdown/mmark/mast"
	"github.com/mmarkdown/mmark/mast/reference"
)

func CitationToBibliography(p *parser.Parser, doc ast.Node) (normative ast.Node, informative ast.Node) {
	seen := map[string]*mast.BibliographyItem{}
	raw := map[string][]byte{}

	// Gather all citations.
	// Gather all reference HTML Blocks to see if we have XML we can output.
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		switch c := node.(type) {
		case *ast.Citation:
			for i, d := range c.Destination {
				if _, ok := seen[string(bytes.ToLower(d))]; ok {
					continue
				}
				ref := &mast.BibliographyItem{}
				ref.Anchor = d
				ref.Type = c.Type[i]

				seen[string(d)] = ref
			}
		case *ast.HTMLBlock:
			anchor := anchorFromReference(c.Content)
			raw[string(bytes.ToLower(anchor))] = c.Content
		}
		return ast.GoToNext
	})

	for _, r := range seen {
		// If we have a reference anchor and the raw XML add that here.
		if raw, ok := raw[string(bytes.ToLower(r.Anchor))]; ok {
			var x reference.Reference
			if e := xml.Unmarshal(raw, &x); e != nil {
				log.Printf("Failed to unmarshal reference: %q: %s", r.Anchor, e)
				continue
			}
			r.Raw = raw
			r.Reference = x
		}

		switch r.Type {
		case ast.CitationTypeInformative:
			if informative == nil {
				informative = &mast.Bibliography{}
			}
			ast.AppendChild(informative, r)
		case ast.CitationTypeSuppressed:
			fallthrough
		case ast.CitationTypeNormative:
			if normative == nil {
				normative = &mast.Bibliography{}
			}
			ast.AppendChild(normative, r)
		}
	}
	return normative, informative
}

// Parse '<reference anchor='CBR03' target=''>' and return the string after anchor= is the ID for the reference.
func anchorFromReference(data []byte) []byte {
	if !bytes.HasPrefix(data, []byte("<reference ")) {
		return nil
	}

	anchor := bytes.Index(data, []byte("anchor="))
	if anchor < 0 {
		return nil
	}

	beg := anchor + 7
	if beg >= len(data) {
		return nil
	}

	quote := data[beg]

	i := beg + 1
	// scan for an end-of-reference marker
	for i < len(data) && data[i] != quote {
		i++
	}
	// no end-of-reference marker
	if i >= len(data) {
		return nil
	}
	return data[beg+1 : i]
}

func ReferenceHook(data []byte) (ast.Node, []byte, int) {
	if !bytes.HasPrefix(data, []byte("<reference ")) {
		return nil, nil, 0
	}

	i := 12
	// scan for an end-of-reference marker, across lines if necessary
	for i < len(data) &&
		!(data[i-12] == '<' && data[i-11] == '/' && data[i-10] == 'r' && data[i-9] == 'e' && data[i-8] == 'f' &&
			data[i-7] == 'e' && data[i-6] == 'r' && data[i-5] == 'e' &&
			data[i-4] == 'n' && data[i-3] == 'c' && data[i-2] == 'e' &&
			data[i-1] == '>') {
		i++
	}

	// no end-of-reference marker
	if i > len(data) {
		return nil, nil, 0
	}

	node := &ast.HTMLBlock{}
	node.Content = data[:i]
	return node, nil, i
}
