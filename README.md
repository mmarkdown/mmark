---
title: "About"
date: 2018-07-22T14:05:51+01:00
aliases: [/about/]
---

Mmark is a powerful markdown processor written in Go, geared towards writing IETF documents. It is,
however, *also* suited for writing complete books and other technical documentation, like the
[Learning Go book](https://miek.nl/go) ([mmark source](https://github.com/miekg/learninggo), and
[I-D text output](https://miek.nl/go/learninggo-2.txt)).

Also see [this repository](https://github.com/danyork/writing-internet-drafts-in-markdown) on how to
write RFC using Markdown.

It provides an advanced markdown dialect that processes file(s) to produce internet-drafts in XML
[RFC 7991](https://tools.ietf.org/html/rfc7991) format. Mmark can produce xml2rfc (aforementioned
RFC 7991), HTML5 output, and manual pages.

Example RFCs in Mmark format can be [found in the Github
repository](https://github.com/mmarkdown/mmark/tree/master/rfc).

Mmark uses [gomarkdown](https://github.com/gomarkdown/markdown) which is a fork of
[blackfriday](https://github.com/russross/blackfriday/). See its
[README.md](https://github.com/gomarkdown/markdown/blob/master/README.md) for more documentation.

## Syntax

Mmark's syntax and the extra features compared to plain Markdown are detailed in
[syntax.md](https://mmark.miek.nl/syntax).

Mmark adds the following syntax elements to
[gomarkdown/markdown](https://github.com/gomarkdown/markdown/blob/master/README.md):

* (Extended) [title block](https://mmark.miek.nl/post/syntax/#title-block).
* [Special sections](https://mmark.miek.nl/post/syntax/#special-sections).
* [Including other files](https://mmark.miek.nl/post/syntax/#including-files) with the option to specify line ranges, regular
  expressions and/or prefix each line with a string. By default only files on the same level, or
  below are allowed to be included (see the `-unsafe` flag).
* [Document divisions](https://mmark.miek.nl/post/syntax/#document-divisions).
* [Captions](https://mmark.miek.nl/post/syntax/#captions) for code, tables and quotes
* [Asides](https://mmark.miek.nl/post/syntax/#asides).
* [Figures and Subfigures](https://mmark.miek.nl/post/syntax/#figures-and-subfigures) - bundle (sub)figures
  into a larger figure.
* [Block Level Attributes](https://mmark.miek.nl/post/syntax/#block-level-attributes) that allow to specify attributes, classes and
  IDs for elements.
* [Indices](https://mmark.miek.nl/post/syntax/#indices) to mark an item (and/or a subitem) to be referenced in the document index.
* [Citations](https://mmark.miek.nl/post/syntax/#citations) and adding [XML References](https://mmark.miek.nl/post/syntax/#xml-references)
* [In document cross references](https://mmark.miek.nl/post/syntax/#cross-references), short form of referencing a section in the
  document.
* [Super- and Subscript](https://mmark.miek.nl/post/syntax/#super-and-subscript).
* [Callouts](https://mmark.miek.nl/post/syntax/#callouts) in code and text.
* [BCP14](https://mmark.miek.nl/post/syntax/#bcp14) (RFC 2119) keyword detection.

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

Outputting HTML5 is done with the `-html` switch.

Files edited under Windows/Mac and using Windows style will be converted into Unix style line ending
before parsing. Any output from `mmark` will use Unix line endings.

[1]: https://daringfireball.net/projects/markdown/ "Markdown"
[2]: https://golang.org/ "Go Language"

Note there are no _wrong_ markdown documents, so `mmark` will only warn about things that are not
right. This may result in invalid XML. Any warning from `mmark` are send to standard error, to see
and check for those you can discard standard output to just leave standard error: `./mmark
rfc/3515.md > /dev/null`.

## Example RFC

The rfc/ directory contains a couple of example RFCs that can be build via v3 tool chain.
The build the text files, just run:

~~~ sh
cd rfc
make txt
~~~

Official RFCs are in rfc/orig (so you can compare the text output from mmark).

## Also See

[Kramdown-rfc2629](https://github.com/cabo/kramdown-rfc2629) is another tool to process markdown and
output XML2RFC XML.

See Syntax.md for a primer on how to use the Markdown syntax to create IETF documents.
