---
title: "About"
date: 2018-07-22T14:05:51+01:00
aliases: [/about/]
---

Mmark is a powerful markdown processor written in Go, geared towards writing IETF documents. It is,
however, *also* suited for writing complete books and other technical documentation, like the
[Learning Go book](https://miek.nl/go) ([mmark source](https://github.com/miekg/learninggo)).

It provides an advanced markdown dialect that processes file(s) to produce internet-drafts in XML
[RFC 7991](https://tools.ietf.org/html/rfc7991) format. Mmark can produce XML2RFC (aforementioned
RFC 7991) and HTML5 output.

Example RFCs can be [found in the Github repository](https://github.com/mmarkdown/mmark/tree/master/rfc).

Mmark uses [gomarkdown](https://github.com/gomarkdown/markdown) which is a fork of
[blackfriday](https://github.com/russross/blackfriday/).

If you like Go and parsing text, drop me (<mailto:miek@miek.nl>) a line if you want to be part of
the *Mmarkdown* Github org, and help develop Mmark!

## Syntax

Mmark's syntax and the extra feature compaired to plain Markdown are detailed in [syntax.md](/syntax).

Mmark adds the following syntax elements to [gomarkdown/black
friday](https://github.com/russross/blackfriday/blob/master/README.md):

* TOML titleblock
* Including other files
* Table and codeblock captions
* Quote attribution (quote "captions")
* Table footers
* Subfigures
* Inline Attribute Lists
* Indices
* Citations
* Abstract/Preface/Notes sections
* Asides
* Main-, middle- and backmatter divisions
* Example lists
* BCP14 (RFC2119) keyword detection
* Include raw XML references
* Callouts in code blocks.

TODO(miek): reference these in the syntax doc.

## Usage

In the mmark subdirectory you can build the mmark tool:

    % cd mmark
    % go build
    % ./mmark -version
    2.0.0

To output v3 xml just give it a markdown file and:

    % ./mmark/mmark mmark2rfc.md

Making a draft in text form:

    % ./mmark/mmark mmark2rfc.md > x.xml
    % xml2rfc --v3 --text x.xml

Outputting HTML is done with the `-html` switch.

[1]: https://daringfireball.net/projects/markdown/ "Markdown"
[2]: https://golang.org/ "Go Language"
