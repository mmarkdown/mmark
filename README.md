---
title: "About"
date: 2018-07-22T14:05:51+01:00
aliases: [/about/]
---

[![Build Status](https://img.shields.io/travis/mmarkdown/mmark/master.svg?label=build)](https://travis-ci.org/mmarkdown/mmark)

Mmark is a powerful markdown processor written in Go, geared towards writing IETF documents. It is,
however, *also* suited for writing complete books and other technical documentation, like the
[Learning Go book](https://miek.nl/go) ([mmark source](https://github.com/miekg/learninggo), and
[I-D text output](https://miek.nl/go/learninggo-2.txt)).

Also see [this repository](https://github.com/danyork/writing-internet-drafts-in-markdown) on how to
write RFC using Markdown.

It provides an advanced markdown dialect that processes file(s) to produce internet-drafts in XML
[RFC 7991](https://tools.ietf.org/html/rfc7991) format. Mmark can produce xml2rfc (aforementioned
RFC 7991), RFC 7749 (xml2rfc version 2 - now deprecated), HTML5 output, markdown and manual pages.

Example RFCs in Mmark format can be [found in the Github
repository](https://github.com/mmarkdown/mmark/tree/master/rfc).

Mmark uses [gomarkdown](https://github.com/gomarkdown/markdown) which is a fork of
[blackfriday](https://github.com/russross/blackfriday/). See its
[README.md](https://github.com/gomarkdown/markdown/blob/master/README.md) for more documentation.

If you like Go and parsing text, drop me (<mailto:miek@miek.nl>) a line if you want to be part of
the *Mmarkdown* Github org, and help develop Mmark!

## Syntax

Mmark's syntax and the extra features compared to plain Markdown are detailed in
[syntax.md](https://mmark.miek.nl/syntax).

Mmark adds the following syntax elements to
[gomarkdown/markdown](https://github.com/gomarkdown/markdown/blob/master/README.md):

* (Extended) [title block](https://mmark.miek.nl/syntax#title-block).
* [Special sections](https://mmark.miek.nl/syntax#special-sections).
* [Including other files](https://mmark.miek.nl/syntax#including-files) with the option to specify line ranges, regular
  expressions and/or prefix each line with a string. By default only files on the same level, or
  below are allowed to be included (see the `-unsafe` flag).
* [Document divisions](https://mmark.miek.nl/syntax#document-divisions).
* [Captions](https://mmark.miek.nl/syntax#captions) for code, tables and quotes
* [Asides](https://mmark.miek.nl/syntax#asides).
* [Figures and Subfigures](https://mmark.miek.nl/syntax#figures-and-subfigures) - bundle (sub)figures
  into a larger figure.
* [Block Level Attributes](https://mmark.miek.nl/syntax#block-level-attributes) that allow to specify attributes, classes and
  IDs for elements.
* [Indices](https://mmark.miek.nl/syntax#indices) to mark an item (and/or a subitem) to be referenced in the document index.
* [Citations](https://mmark.miek.nl/syntax#citations) and adding [XML References](https://mmark.miek.nl/syntax#xml-references)
* [In document cross references](https://mmark.miek.nl/syntax#cross-references), short form of referencing a section in the
  document.
* [Super- and Subscript](https://mmark.miek.nl/syntax#super-and-subscript).
* [Callouts](https://mmark.miek.nl/syntax#callouts) in code and text.
* [BCP14](https://mmark.miek.nl/syntax#bcp14) (RFC 2119) keyword detection.

## Usage

You can [download a binary](https://github.com/mmarkdown/mmark/releases) or optionally build mmark
your self. You'll need a working [Go environment](https://golang.org), then check out the code and:

    % go get && go build
    % ./mmark -version
    2.0.0

To output XML2RFC v3 xml just give it a markdown file and:

    % ./mmark rfc/3514.md

Making a draft in text form (v3 output)

    % ./mmark rfc/3514.md > x.xml
    % xml2rfc --v3 --text x.xml

Outputting HTML5 is done with the `-html` switch. Outputting markdown is done with the `-markdown`
switch (optionally you can use `-width` to set the text width).

[1]: https://daringfireball.net/projects/markdown/ "Markdown"
[2]: https://golang.org/ "Go Language"

## Example RFC

The rfc/ directory contains a couple of example RFCs that can be build via the v2 or v3 tool chain.
The build the text files, just run:

~~~ sh
cd rfc
make txt
~~~

Official RFCs are in rfc/orig (so you can compare the text output from mmark).

## Using Mmark as a library

By default Mmark gives you a binary you can run, if you want to include the parser and renderers in
your own code you'll have to lift some of it out of `mmark.go`.

Create a parser with the correct options and flags. The that `init` is used to track file includes.
In this snippet we set if to `fileName` which is the file we're currently reading. If reading from
standard input, this can be set to `""`.

~~~ go
p := parser.NewWithExtensions(mparser.Extensions)
init := mparser.NewInitial(fileName)
documentTitle := "" // hack to get document title from TOML title block and then set it here.
p.Opts = parser.Options{
    ParserHook: func(data []byte) (ast.Node, []byte, int) {
        node, data, consumed := mparser.Hook(data)
        if t, ok := node.(*mast.Title); ok {
            if !t.IsTriggerDash() {
                documentTitle = t.TitleData.Title
            }
        }
        return node, data, consumed
    },
    ReadIncludeFn: init.ReadInclude,
    Flags:         parserFlags,
}
~~~

Then parser the document (`d` is a `[]byte` containing the document text):

~~~ go
doc := markdown.Parse(d, p)
mparser.AddBibliography(doc)
mparser.AddIndex(doc)
~~~

After this `doc` is ready to be rendered. Create a renderer, with a bunch of options.

~~~ go
opts := html.RendererOptions{
    Comments:       [][]byte{[]byte("//"), []byte("#")}, // used for callouts.
	RenderNodeHook: mhtml.RenderHook,
	Flags:          html.CommonFlags | html.FootnoteNoHRTag | html.FootnoteReturnLinks| html.CompletePage,
	Generator:      `  <meta name="GENERATOR" content="github.com/mmarkdown/mmark Mmark Markdown Processor - mmark.miek.nl`,
}
opts.Title = documentTitle // hack to add-in discovered title

renderer := html.NewRenderer(opts)
~~~

Next we we only need to generate the HTML: `x := markdown.Render(doc, renderer)`. Now `x` contains
a `[]byte` with the HTML.

## Also See

[Kramdown-rfc2629](https://github.com/cabo/kramdown-rfc2629) is another tool to process markdown and
output XML2RFC XML.
