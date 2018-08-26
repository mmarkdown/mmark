---
title: "About"
date: 2018-07-22T14:05:51+01:00
aliases: [/about/]
---

[![Build Status](https://img.shields.io/travis/mmarkdown/mmark/master.svg?label=build)](https://travis-ci.org/mmarkdown/mmark)

Mmark is a powerful markdown processor written in Go, geared towards writing IETF documents. It is,
however, *also* suited for writing complete books and other technical documentation, like the
[Learning Go book](https://miek.nl/go) ([mmark source](https://github.com/miekg/learninggo)).

It provides an advanced markdown dialect that processes file(s) to produce internet-drafts in XML
[RFC 7991](https://tools.ietf.org/html/rfc7991) format. Mmark can produce xml2rfc (aforementioned
RFC 7991), RFC 7749 (xml2rfc version 2) and HTML5 output.

Example RFCs can be [found in the Github repository](https://github.com/mmarkdown/mmark/tree/master/rfc).

Mmark uses [gomarkdown](https://github.com/gomarkdown/markdown) which is a fork of
[blackfriday](https://github.com/russross/blackfriday/).

If you like Go and parsing text, drop me (<mailto:miek@miek.nl>) a line if you want to be part of
the *Mmarkdown* Github org, and help develop Mmark!

## Syntax

Mmark's syntax and the extra features compared to plain Markdown are detailed in
[syntax.md](https://mmark.nl/syntax).

Mmark adds the following syntax elements to
[gomarkdown/markdown](https://github.com/gomarkdown/markdown/blob/master/README.md):

* (Extended) [title block](https://mmark.nl/syntax#title-block).
* [Special sections](https://mmark.nl/syntax#special-sections).
* [Including other files](https://mmark.nl/syntax#including-files) with the option to specify line ranges, regular
  expressions and/or prefix each line with a string.
* [Document divisions](https://mmark.nl/syntax#document-divisions).
* [Captions](https://mmark.nl/syntax#captions) for code, tables and quotes
* [Asides](https://mmark.nl/syntax#asides).
* [Figures and Subfigures](https://mmark.nl/syntax#figures-and-subfigures) - this syntax is still under consideration as is
  "do we really need this?"
* [Block Level Attributes](https://mmark.nl/syntax#block-level-attributes) that allow to specify attributes, classes and
  IDs for elements.
* [Indices](https://mmark.nl/syntax#indices) to mark an item (and/or a subitem) to be referenced in the document index.
* [Citations](https://mmark.nl/syntax#citations) and adding [XML References](https://mmark.nl/syntax#xml-references)
* [In document cross references](https://mmark.nl/syntax#cross-references), short form of referencing a section in the
  document.
* [Super- and Subscript](https://mmark.nl/syntax#super-and-subscript).
* [Callouts](https://mmark.nl/syntax#callouts) in code and text.
* [BCP14](https://mmark.nl/syntax#bcp14) (RFC 2119) keyword detection.

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

Making a draft in text form (v2 output)

    % ./mmark -2 rfc/3514.md > x.xml
    % xml2rfc --text x.xml

Outputting HTML5 is done with the `-html` switch. Outputting RFC 7749 is done with `-2`.

[1]: https://daringfireball.net/projects/markdown/ "Markdown"
[2]: https://golang.org/ "Go Language"
