---
title: "Syntax"
date: 2018-07-22T14:05:51+01:00
aliases: [/syntax/]
---

This is version 2 of [Mmark](https://github.com/mmarkdown/mmark):
based on a [new markdown implementation](https://github.com/mmarkdown/markdown)
and some (small) language changes as well. We think these language changes lead to a more consistent
user experience and lead to less confusion.

See [changes from v1](#changes-from-version-1) if you're comming from version 1.

# Mmark V2 Syntax

This document describes all the *extra* syntax elements that can be used in Mmark. Mmark's syntax is
based on the ["standard" Markdown syntax](https://daringfireball.net/projects/markdown/syntax).

> Read the above document if you haven't already, it helps you understand how markdown looks and feels.

For the rest we build up on <https://github.com/gomarkdown/markdown> and support all syntax
[it supports](https://github.com/gomarkdown/markdown/blob/master/README.md). We enable the following
extensions by default:

* *Strikethrough* allow strikethrough text using `~~test~~`.
* *Footnotes* Pandoc style footnotes.
* *HeadingIDs*, specify heading IDs  with `{#id}`.
* *AutoHeadingIDs*, create the heading ID from the text.
* *DefinitionLists*, parse definition lists.
* *MathJax*, parse MathJax
* *OrderedListStart*, notice start element of ordered list.
* *Attributes* allow block level attributes.

Mmark adds numerous enhancements to make it suitable for writing (IETF) Internet Drafts and even
complete books. It <strike>steals</strike> borrows syntax elements from [pandoc], [kramdown],
[leanpub], [asciidoc], [PHP markdown extra] and [Scholarly markdown].

[kramdown]: https://kramdown.gettalong.org/
[leanpub]: https://leanpub.com/help/manual
[asciidoc]: http://www.methods.co.nz/asciidoc/
[PHP markdown extra]: http://michelf.com/projects/php-markdown/extra/
[pandoc]: http://johnmacfarlane.net/pandoc/
[CommonMark]: http://commonmark.org/
[Scholarly markdown]: http://scholarlymarkdown.com/Scholarly-Markdown-Guide.html

## What does Mmark add?

TODO(miek): this list needs to link to the sections detailing the options.

Mmark adds:

* Extended title block
* Including other files with the option to specify line ranges and/or prefix each line with a string
* Document divisions
* Captions for code, tables and quotes
* Asides and other unnumbered sections (i.e. Abstract)
* Indices
* Citations
* Callouts

### RFC 7991 XML Output

This is the output format used for generating Internet-Drafts and RFCs. The generated XML needs to
be processed by another tool (xml2rfc) to generate to official (final) output. The XML from *mmark*
can be used directly to upload to the IETF tools website.

Title Block:
:   If the document has a [title block](#title-block) the front matter is already open. Closing the
    front matter can only be done by starting the middle matter with `{mainmatter}`. Any open
    "matters" are closed when the document ends.

Abstract:
:   The abstract can be started by using the special header syntax `.# Abstract`

Note:
:   Any special header that is not "abstract" or "preface" will be a
    [note](https://tools.ietf.org/html/rfc7749#section-2.24): a numberless section.

BCP 14/RFC 2119 Keywords:
:   If an RFC 2119 word is found enclosed in `**` it will be rendered
    as an `<bcp14>` element: i.e. `**MUST**` becomes `<bcp14>MUST</bcp14>`.

Artwork:
:   Artwork is added by using a (fenced) code block. If the code block has an caption it will be
    wrapped in a `<figure>`, this is true for source code as well.

Source code:
:   If you want to typeset a source code instead of an artwork you must specify a language to the
    fenced block:

    ~~~
    ``` go
    println(hello)
    ````
    ~~~
    Will be typesets as source code with the language set to `go`.

## Block Elements

### Title Block

A Title Block contains a document's meta data; title, authors, date and other elements. The elements
that can be specified are copied from the [xml2rfc v3
standard](https://tools.ietf.org/html/rfc7791). More on these below. The complete title block is
specified in [TOML](https://github.com/toml-lang/toml). Examples title blocks can be [found in the
repository of mmark](https://github.com/mmarkdown/mmark/tree/master/rfc).

The title block itself needs three or more `%`'s at the start and end of the block. A minimal title
block would look like this:

~~~
%%%
title = "Foo Bar"
%%%
~~~

#### Elements of the Title Block

An I-D needs to have a Title Block with the following items filled out:

* title - the main title of the document.
* abbrev - abbreviation of the title.
* updates/obsoletes - array of integers.
* seriesInfo, containing
   * name - `RFC` or `Internet-Draft` or `DOI`
   * value - draft name or RFC number
   * stream - `IETF` (default), `IAB`, `IRTF` or `independent`.
   * status - `standard`, `informational`, `experimental`, `bcp`, `fyi`, or `full-standard`.
* ipr - usually just set `trust200902`
* area - usually just `Internet`
* workgroup - the workgroup the document is created for
* keyword - array with keywords (optional)
* author(s) - define all the authors.
* date - the date for this I-D/RFC.

An example would be:

~~~
%%%
title = "Using mmark to create I-Ds and RFCs"
abbrev = "mmark2rfc"
updates = [1925, 7511]
ipr= "trust200902"
area = "Internet"
workgroup = ""
keyword = ["markdown", "xml", "mmark"]

[seriesInfo]
status = "informational"
name = "Internet-Draft"
value = "draft-gieben-mmark2rfc-00"
stream = "IETF"

date = 2014-12-10T00:00:00Z

[[author]]
initials="R."
surname="Gieben"
fullname="R. (Miek) Gieben"
organization = "Mmark"
  [author.address]
  email = "miek@miek.nl"
%%%
~~~

An `#` acts as a comment in this block. TOML itself is specified [here](https://github.com/toml-lang/toml).

### Including Files

Including other files can done be with `{{filename}}`, if the path of `filename` is *not* absolute,
the filename is taken relative to *current file being processed*. With `<{{filename}}`
you include a file as a code block. The main difference being it will be returned as a code
block. The file's extension *will be used* as the language. The syntax is:

~~~
{{pathname}}[address]
~~~
And address can be `N,M`, where `N` and `M` are line numbers. Or `/N/,/M/`, where `N` and `M` are
regular expressions that include from where to where to include the file. Each of these can have
an optional `prefix=""` specifier.

~~~
{{filename}}[3,5]
~~~

Only includes the lines 3 to (*not* inclusive) 5 into the current document.

~~~
{{filename}}[3,5;prefix="C: "]
~~~
will include the same lines *and* prefix each include line with `C: `.

Captioning works as well:

~~~
<{{test.go}}[/START/,/END/]
Caption: A sample function.
~~~

### Document Divisions

Mmark support three document divisions, front matter, main matter and the back matter. Mmark
automatically starts the front matter for you *if* the document has a title block. Switching
divisions can be done with `{frontmatter}`, `{mainmatter}` and `{backmatter}`. This must be the only
thing on the line.

## Captions

Mmark supports caption below [tables](#tables), [code blocks](#code-blocks) and [block
quotes](#block-quotes). You can caption each elements with `Caption: `. The caption extends to the
first *empty* line. Some examples:

~~~
Name    | Age
--------|-----:
Bob     | 27
Alice   | 23
Caption: This is the table caption.
~~~

Or for a code block:

     ~~~ go
     func getTrue() bool {
         return true
     }
     ~~~
     Caption: This is a caption for a code block.

And for a quote:

     > Ability is nothing without opportunity.
     Caption: https://example.com, Napoleon Bonaparte

### Asides

Any text prefixed with `A>` will become an
[aside](https://developer.mozilla.org/en/docs/Web/HTML/Element/aside). This is similar to a block
quote, but can be styled differently.

### Figures and Subfigures

> TODO TODO TODO

To group artworks and code blocks into figures, we need an extra syntax element.
[Scholarly markdown] has a neat syntax
for this. It uses a special section syntax and all images in that section become
subfigures of a larger figure. Disadvantage of this syntax is that it can not be
used in lists. Hence we use a quote like solution, just like asides,
but for figures: we prefix the entire paragraph with `F>`. Each of the images and
or code block included will be part of the same over arching figure.

Basic usage:

~~~
F> ~~~ ascii-art
F> +-----+
F> | ART |
F> +-----+
F> ~~~~
F> Caption: This caption is ignored in v3, but used in v2.
F>
F> ~~~ c
F> printf("%s\n", "hello");
F> ~~~
F>
Caption: Caption for both figures in v3 (in v2 this is ignored).
~~~

### Example lists

> TODO TODO TODO

### XML References

Any valid XML reference fragment found anywhere in the document, can be used as a citation reference.
The syntax of the XML reference element is defined in [RFC
7749](https://tools.ietf.org/html/rfc7749#section-2.30). The `anchor` defined can be used in the
[citation](#Citations), which the example below that would be `[@pandoc]`:

~~~
<reference anchor='pandoc' target='http://johnmacfarlane.net/pandoc/'>
    <front>
        <title>Pandoc, a universal document converter</title>
        <author initials='J.' surname='MacFarlane' fullname='John MacFarlane'>
            <organization>University of California, Berkeley</organization>
            <address>
                <email>jgm@berkeley.edu</email>
                <uri>http://johnmacfarlane.net/</uri>
            </address>
        </author>
        <date year='2006'/>
    </front>
</reference>
~~~

Note that for citing I-Ds and RFCs you *don't* need to include any XML, as Mmark will pull these
automatically from their online location: or technically more correct: the xml2rfc post processor
will do this.

## Span Elements

### Indices

Defining indices allows you to create an index. The define an index use the `(!item)`. Sub items can
be added as well, with `(!item; subitem)`. To make `item` primary, use another `!`: `(!!item,
subitem)`. If any index is defined the end of the document contains the list of indices. The
`-index=false` flag suppresses this generation.

### Citations

Mmark uses the citation syntax from Pandoc: `[@RFC2535]`, the citation can either be informative
(default) or normative, this can be indicated by using the `?` or `!` modifier: `[@!RFC2535]` create
a normative reference for RFC 2535. To suppress a citation use `[@-RFC1000]`. It will still add the
citation to the references, but does not show up in the document as a citation.

The first seen modifier determines the type (suppressed, normative or informative).
Multiple citation can separated with a semicolon: `[@RFC1034; @RFC1035]`.

If you reference an RFC or I-D the reference will be added automatically (no need to muck about
with an `<reference>` block.

For I-Ds you may want to add a draft sequence number, which can be done as such: `[@?I-D.blah#06]`.
If you reference an I-D *without* a sequence number it will create a reference to the *last* I-D in
citation index.

A reference section is create by default, but you can suppress it by using the command line flag
`-reference=false`.

### Cross References

Cross references can use the syntax `[](#id)`, but usually the need for the title within the
brackets is not needed, so Mmark has the shorter syntax `(#id)` to cross reference in the document.

Example:

~~~
My header {#header}

Lorem ipsum dolor sit amet, at ultricies ...
See Section (#header).
~~~

### Super- and Subscript

> TODO TODO TODO

For superscript use `^` and for subscripts use `~`. For example:

~~~
H~2~O is a liquid. 2^10^ is 1024.
~~~

Inside a super- or subscript you must escape spaces. Thus, if you want the letter P with 'a cat' in
subscripts, use `P~a\ cat~`, not `P~a cat~`.

## Links and Images

Normal markdown synax.
**SVG TODO and maybe new syntax**

## Block Level Attributes

A "Block Level Attribute" is a list of HTML attributes between braces: `{...}`. It allows you to
set classes, an anchor and other types of *extra* information for the next block level element.

The full syntax is: `{#id .class key="value"}`. Values may be omitted, i,e., just `{.class}` is
valid.

The following example applies the attributes: `type` and `id` to the blockquote:
~~~
{title="The blockquote title" #myid}
> A blockquote with a title
~~~
Gets expanded into:
~~~
<blockquote id="myid" title="The blockquote title">
    <t>A blockquote with a title</t>
</blockquote>
~~~

### Callouts

Callouts are way to reference code from paragraphs following that code. Mmark uses the following
syntax for specifying a callout `<<N>>` where N is integer > 0.

In code blocks you can use the *same* syntax to create a callout:

~~~
    Code  //<<1>>
    More  //<<2>>

As you can see in <<1>> but not in <<2>>. There is no <<3>>.
~~~

Using callouts in source code examples will lead to code examples that do not compile.
To fix this the callout needs to be placed in a comment, but then your source show useless empty comments.
To fix this Mmark will detect (and remove!) the comment from the callout, leaving your
example pristine in the document.

Note that callouts *in code blocks* are only detected if the renderer has been configured to look
for them. The default mmark configuration is to detect them after `//` and `#` comment starters.
See the `-comment` option for to change this.

Lone callouts without them being prefixed with a comment means they are not detected by Mmark.

# Changes from version 1

These are the changes from Mmark version 1:

* Caption under tables, figure, quotes and code block are now *always* done with `Caption: `. No
  more `Table: `, `Quote: `, and `Figure: `.
* Citations:
   * Suppressing a citation is done with `[@-ref]` (it was the reverse `-@` in v1), this is more consistent.
   * Multiple citations are allowed in one go, separated with a semicolons: `[@ref1; @ref2]`.
   * **TODO** Reference text is allowed `[@ref p. 23]`.
* Indices: now just done with `(!item)`, marking one primary will be: `(!!item)`.
* Code block call outs are now a renderer setting, not a [Block Level
  Attribute](#block-level-attributes). Callout in code are *only* detected if they are used after
  a comment.
* Including files with a prefix is now specified in the address specification:
  `{{myfile}}[prefix="C: "]` will use `C: ` as the prefix. No more mucking about with block
  attribute lists that are hard to discover.
* **TODO** Extended table syntax; if this ever comes back it needs to more robust implementation.
* Title Block need to be sandwiched between `%%%`, the prefix `%` does not work anymore.

Syntax that is *not* supported anymore:

* HTML abbreviations.
* The different list syntaxes have been dropped, use a [Block Level
  Attribute](#block-level-attributes) to tweak the output.
* Tasks lists.
* Comment detection, i.e. to support `cref`: dropped. Comments are copied depending on the
  flag `renderer.SkipHTML`.
* Parts
