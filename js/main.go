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

var sampleData = `
%%%
title = "Network Configuration Access Control Model"
abbrev = "NACM"
obsoletes = [6536]
ipr= "trust200902"
area = "Internet"
workgroup = "Network Working Group"
submissiontype = "IETF"
keyword = [""]
date = 2018-03-01T00:00:00Z

[seriesInfo]
name = "RFC"
value = "8341"
stream = "IETF"
status = "standard"

[[author]]
initials="A."
surname="Bierman"
fullname="Andy Bierman"
abbrev = "YumaWorks"
organization = "YumaWorks"
  [author.address]
  email = "andy@yumaworks.com"
  [author.address.postal]
  city = "Simi Valley"
  street = "685 Cochran St."
  code = "CA 93065"
  postalline= ["Suite #160"]
  country = "United States of America"
[[author]]
initials="M."
surname="Bjorklund"
fullname="Martin Bjorklund"
abbrev = "Tail-f Systems"
organization = "Tail-f Systems"
  [author.address]
  email = "mbj@tail-f.com"
%%%

.# Abstract

The standardization of network configuration interfaces for use with the Network Configuration
Protocol (NETCONF) or the RESTCONF protocol requires a structured and secure operating environment
that promotes human usability and multi-vendor interoperability.  There is a need for standard
mechanisms to restrict NETCONF or RESTCONF protocol access for particular users to a preconfigured
subset of all available NETCONF or RESTCONF protocol operations and content.  This document defines
such an access control model.

This document obsoletes RFC 6536.

{mainmatter}

# Introduction

The Network Configuration Protocol (NETCONF) and the RESTCONF protocol do not provide any standard
mechanisms to restrict the protocol operations and content that each user is authorized to access.

There is a need for interoperable management of the controlled access to administrator-selected
portions of the available NETCONF or RESTCONF content within a particular server.

This document addresses access control mechanisms for the Operations and Content layers of NETCONF,
as defined in [@!RFC6241]; and RESTCONF, as defined in [@!RFC8040].  It contains three main
sections:
`

func main() {
	js.Global.Set("mmark", map[string]interface{}{
		"NewDocument": NewDocument,
		"SampleData":  sampleData,
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

func (doc *Document) HTMLFragment() string {
	opts := html.RendererOptions{
		Comments:       [][]byte{[]byte("//"), []byte("#")},
		RenderNodeHook: mhtml.RenderHook,
		Flags:          html.CommonFlags | html.FootnoteNoHRTag | html.FootnoteReturnLinks,
		Title:          doc.title,
	}

	renderer := html.NewRenderer(opts)
	x := markdown.Render(doc.root, renderer)
	return string(x)
}
