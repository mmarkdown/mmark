---
title: "FAQ"
date: 2018-07-22T14:05:51+01:00
aliases: [/faq/]
toc: true
---

Mmark Frequently Asked Questions. Also see the XML2RFCv3 FAQ:
<https://www.rfc-editor.org/materials/FAQ-xml2rfcv3.html>, section below will have the same
questions, but then answered in mmark syntax.

# How Do I Create an Independent IETF Document?

Use the following as starting point for your title block, `ipr` and `submissiontype` are the important
settings here.

~~~ toml
title = "Title"
abbrev = "Title"
ipr = "none"
submissiontype = "independent"
keyword = [""]

[seriesInfo]
name = "Internet-Draft"
value = "draft-00"
stream = "independent"
status = "informational"
~~~

# How Do I Create an IRTF Document?

Set `submissiontype` and `stream` in `seriesInfo` to *IRTF*. Items like `workgroup` function as they
do for normal Internet-Draft documents.

# How Do I Create an IAB Document?

Set `submissiontype` and `stream` in `seriesInfo` to *IAB*. Items like `workgroup` are (I believe)
ignored for this stream.

# How Do I Create an FYI Document?

Use this as the `seriesInfo`:

~~~ toml
[seriesInfo]
name = "FYI"
value = "2100"
stream = "IETF"
status = "informational"
~~~

Note this makes xml2rfc still complain, but at least creates valid XML.

# How Do I Make an Author an Editor

Use `role = "editor"` in the author's section in the titleblock.

## How Do Specify a Contact

Use a `[[contact]]` in the toml header:

~~~ toml
[[contact]]
initials="D."
surname="Addison"
fullname="David Addison"
  [contact.address.postal]
  city = "St. Petersburg"
  code = "FL 33709-4819"
~~~

Using the contact is done by referencing it: `[@David Addison]` (using the `fullname` property). If
the reference is the *first* thing after a new paragraph it will be expanded like XML2RFC expands
authors in an Internet-Draft.

# XML2RFCv3 FAQ

## What version of xml2rfc is supported?

The latest version of xml2rfc is the supported version. As it currently stands, the xml2rfc
*implementation* is the *spec*. Older versions may happen to work, with newer features unsupported,
but this is not guaranteed.

Latest version of xml2rfc can be found at [pypi](https://pypi.org/project/xml2rfc/).

## How do I get different kinds of lists?

Use the standard markdown syntax for unordered, ordered and definition lists.

## How do I get a list like (1), (2), (3) or (a), (b), (c)?

Use a block level attribute: `{type="(%d)"}`, `{type="(%c)"}` or `{type="REQ%d"}`.

## How do I get continuous numbering in a list that is split by text (or across sections)?

Set the group attribute with a block level attribute.

~~~
{type="REQ%d" group="reqs"}
1. do a
2. do b

Here is text in between

{type="REQ%d" group="reqs"}
1. do c
2. do d
~~~

## How do I get indentation? or How do I use definition lists?

~~~ markdown
First Term
: This is the definition of the first term.

Second Term
: This is one definition of the second term.
~~~

A non compact definition list can be done like so: (not the block attribute allows for a newline
after the term):

~~~ markdown
{newline="true"}
First Term

: This is the definition of the first term.

Second Term

: This is one definition of the second term.
~~~

## How do I create nested lists?

~~~ markdown
Foo validator
: It performs the following actions:
  * runs
  * jumps
  * walks
~~~

~~~ markdown
{type="Step %d:"}
1. Send it to
   * Alice
   * Bob
   * Carol
~~~

## How do I insert non-ASCII characters?

This is handled for you, mmark will wrap non-ASCII characters in `<u>`. The `asciiFullname` and
friends used in authors and contacts is currently not implemented.

## How do I insert a table?

Use the markdown table syntax.

## How do I get bold, italics, or a fixed-width font?

* bold: `**bold**`
* italics `*italics*`
* fixed-width, wrap in back-ticks

## How do I get subscript and superscript?

* subscript: `_2_`
* superscript: `^10^`

## Do I have to use the bcp14 element each time a keyword (e.g., "MUST") appears in my document?

Just use `**MUST**`, i.e. make the bcp14 element bold and capital, mmark wraps these in `<bcp14>`
tags.
