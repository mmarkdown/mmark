% MMARK(1)
% Mmark Authors
% August 2018

# NAME

mmark – generate XML or HTML from mmark markdown

# SYNOPSIS

**mmark** [**OPTIONS**] [*FILE...*]

# DESCRIPTION

**Mmark** is a powerful markdown processor written in Go, geared towards writing IETF documents. It
 is, however, *also* suited for writing complete books and other technical documentation.

It provides an advanced markdown dialect that processes file(s) to produce internet-drafts in XML
[RFC 7991](https://tools.ietf.org/html/rfc7991) format. Mmark can produce xml2rfc (aforementioned
RFC 7991), RFC 7749 (xml2rfc version 2) and HTML5 output.

The syntax is detailed at <https://mmark.nl/syntax>.

# OPTIONS

**-2**
:   generate RFC 7749 XML

**-ast**
:    print abstract syntax tree and exit

**-css string**
:    link to a CSS stylesheet (only used with -html)

**-fragment**
:    don't create a full document

**-head string**
:    link to HTML to be included in head (only used with -html)

**-html**
:    create HTML output

**-index**
:    generate an index at the end of the document (default true)

**-bibliography**
:    generate a bibliographtysection after the back matter (default true)

**-version**
:    show mmark version

# ALSO SEE

RFC 7791 and RFC 7749. The main site for Mmark is <https://mmark.nl>
