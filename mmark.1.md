% MMARK(1)
% Mmark Authors
% August 2018

# NAME

mmark – generate XML, HTML or markdown from mmark markdown

# SYNOPSIS

**mmark** [**OPTIONS**] [*FILE...*]

# DESCRIPTION

**Mmark** is a powerful markdown processor written in Go, geared towards writing IETF documents. It
is, however, *also* suited for writing complete books and other technical documentation.

Mmark provides an advanced markdown dialect that processes file(s) to produce internet-drafts in XML
[RFC 7991](https://tools.ietf.org/html/rfc7991) format. Mmark can produce xml2rfc (aforementioned
RFC 7991), RFC 7749 (xml2rfc version 2), HTML5 and markdown output.

The syntax is detailed at <https://mmark.nl/syntax>.

# OPTIONS

**-ast**
:    print abstract syntax tree and exit.

**-css string**
:    link to a CSS stylesheet (only used with -html).

**-fragment**
:    don't create a full document.

**-head string**
:    link to HTML to be included in head (only used with -html).

**-html**
:    create HTML output.

**-2**
:   generate RFC 7749 XML.

**-markdown**
:    output (normalized) markdown.

**-unsafe**
:    allow includes from anywhere in the filesystem, otherwise they are only allowed *under* the
     current document.

**-textwidth integer**
:    set the text width when generating markdown, defaults to 100 characters.

**-w**
:    write to source file when generating markdown.

**-index**
:    generate an index at the end of the document (default true).

**-bibliography**
:    generate a bibliography section after the back matter (default true), this needs
     a `{{backmatter}}` in the document.

**-version**
:    show mmark's version.

# ALSO SEE

RFC 7991 and RFC 7749. The main site for Mmark is <https://mmark.nl>
