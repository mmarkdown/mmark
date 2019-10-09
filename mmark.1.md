---
title: 'MMARK(1)'
author:
    - Mmark Authors
date: October 2019
---

# NAME

mmark – generate XML, HTML or markdown from mmark markdown documents.

# SYNOPSIS

**mmark** [**OPTIONS**] [*FILE...*]

# DESCRIPTION

**Mmark** is a powerful markdown processor written in Go, geared towards writing IETF documents. It
is, however, *also* suited for writing complete books and other technical documentation.

Mmark provides an advanced markdown dialect that processes file(s) to produce internet-drafts in XML
[RFC 7991](https://tools.ietf.org/html/rfc7991) format. Mmark can produce xml2rfc (aforementioned
RFC 7991), HTML5, manual pages and markdown output.

The syntax is detailed at [https://mmark.miek.nl/syntax](https://mmark.miek.nl/syntax).

Compared to other markdown variants mmark adds:

*  (Extended) Title block to specify authors and IETF specific bits.

*  (Special) sections for abstracts, prefaces or notes.

*  Including other files with the option to specify line ranges, regular expressions and/or
   prefixing each line with a custom string.

*  Document divisions: front matter, main matter and back matter.

*  Captions for code, tables, quotes and subfigures.

*  Asides for small notes.

*  Figures and Subfigures that allow for grouping images into subfigures as well as giving a single
   image metadata (a link, attributes, etc.).

*  Block Level Attributes that allow to specify attributes, classes and IDs for elements.

*  Indices to mark an item (and/or a subitem) to be referenced in the document index.

*  Citations and adding XML References (those used by the IETF).

*  Cross references: short form of referencing a section in the document.

*  Super- and Subscript.

*  Callouts in code and text.

*  BCP14 (RFC 2119) keyword detection.

## RFC 7991

This is the XML format used by the RFC editor for accepting Internet-Drafts.

## HTML5

The HTML5 renderer outputs HTML.

## Markdown

Mmark can also be "translated" into markdown again. This is a useful feature for auto-formatting
markdown files.

## Manual Pages

The man renderer outputs nroff that can be viewed via man(1).

# OPTIONS

**-ast**

:  print abstract syntax tree and exit.

**-css string**

:  link to a CSS stylesheet (only used with -html).

**-fragment**

:  don't create a full document.

**-head string**

:  link to HTML to be included in head (only used with -html).

**-html**

:  create HTML output.

**-markdown**

:  output (normalized) markdown.

**-man**

:  output nroff (manual pages).

**-unsafe**

:  allow includes from anywhere in the filesystem, otherwise they are only allowed *under* the
   current document.

**-textwidth integer**

:  set the text width when generating markdown, defaults to 100 characters.

**-w**

:  write to source file when generating markdown.

**-index**

:  generate an index at the end of the document (default true).

**-bibliography**

:  generate a bibliography section after the back matter (default true), this needs a
   `{{backmatter}}` in the document.

**-version**

:  show mmark's version.

# ALSO SEE

RFC 7991 and (maybe) RFC 7749. The main site for Mmark is
[https://mmark.miek.nl](https://mmark.miek.nl). The syntax used by mmark is explained in the [syntax
page](https://mmark.miek.nl/post/syntax/).
