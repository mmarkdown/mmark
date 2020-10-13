---
title: "Syntax"
date: 2018-07-22T14:05:51+01:00
aliases: [/syntax/]
toc: true
---

This is version 2 of [Mmark](https://github.com/mmarkdown/mmark): based on a [new markdown
implementation](https://github.com/gomarkdown/markdown) and some (small) language changes as well. We
think these language changes lead to a more consistent user experience and lead to less confusion.

See [changes from v1](#changes-from-version-1) if you're coming from version 1.

Biggest changes:

* Including files is now done relative to the file being parsed (i.e. the *sane* way).
* Block attributes apply to block elements *only*.
* Callouts
    * *always* rendered and require double greater/less-than signs, `<<1>>`.
    * *always* require a comment in the code, i.e. `//<<1>>` will be rendered as a callout, a plain
    `<<1>>` will not.
* Block Tables have been dropped.
* Example lists (originally copied from Pandoc) have been dropped.
* Plain citations, i.e. `@RFC5412`, when the reference was previously seen don't work anymore,
  always use the full syntax `[@RFC5412]`.

# Why this new version?

It fixes a bunch of long standing bugs and the parser generates an abstract syntax tree (AST). It
will be easier to add new renderers with this setup. It is also closer to Common Mark. So we took
this opportunity to support RFC 7991 XML (xml2rfc version 3), HTML5, RFC 7749 XML (xml2rfc version 2
- now deprecated) and Markdown output (use mmark is a markdown formatter). Also with code upstreamed
(to [gomarkdown](https://github.com/gomarkdown)), we have less code to maintain.

Because of the abstract syntax tree it will also be easier to write helper tools, like, for instance
a tool that checks if all referenced labels in the document are actually defined. Another idea could
be to write a "check-the-code" tool that syntax checks all code in code blocks. Eventually these
could be build into the `mmark` binary itself. See some [fun
ideas](https://github.com/mmarkdown/filter) here.

# Mmark V2 Syntax

This document describes all the *extra* syntax elements that can be used in Mmark. Mmark's syntax is
based on the ["standard" Markdown syntax](https://daringfireball.net/projects/markdown/syntax).
A good primer on what blackfriday implements [is this
article](https://notes.peter-baumgartner.net/archive/content-organisation/blackfriday-markdown/).

> Read the above documents if you haven't already, it helps you understand how markdown looks and feels.

For the rest we build up on <https://github.com/gomarkdown/markdown> and support all syntax
[it supports](https://github.com/gomarkdown/markdown/blob/master/README.md). We enable the following
[extensions](https://github.com/gomarkdown/markdown/blob/master/README.md#extensions) by default:

* *Strikethrough*, allow strike through text using `~~test~~`.
* *Autolink*, detect embedded URLs that are not explicitly marked.
* *Footnotes* Pandoc style footnotes.
* *HeadingIDs*, specify heading IDs  with `{#id}`.
* *AutoHeadingIDs*, create the heading ID from the text.
* *DefinitionLists*, parse definition lists.
* *MathJax*, parse MathJax
* *OrderedListStart*, notice start element of ordered list.
* *Attributes*, allow block level attributes.
* *Smartypants*, expand `--` and `---` into ndash and mdashes.
* *SuperSubscript*, parse super- and subscript: H~2~O is water and 2^10^ is 1024.
* *Tables*, parse tables.
* *NonBlockingSpace*, convert "backslash space" into a non blocking space.

Mmark adds numerous enhancements to make it suitable for writing ([IETF](https://ietf.org)) Internet
Drafts and even complete books. It <strike>steals</strike> borrows syntax elements from [pandoc],
[kramdown], [leanpub], [asciidoc], [PHP markdown extra] and [Scholarly markdown].

[kramdown]: https://kramdown.gettalong.org/
[leanpub]: https://leanpub.com/help/manual
[asciidoc]: http://www.methods.co.nz/asciidoc/
[PHP markdown extra]: http://michelf.com/projects/php-markdown/extra/
[pandoc]: http://johnmacfarlane.net/pandoc/
[CommonMark]: http://commonmark.org/
[Scholarly markdown]: http://scholarlymarkdown.com/Scholarly-Markdown-Guide.html

## What does Mmark add?

Mmark adds:

* (Extended) [title block](#title-block) to specify authors and IETF specific bits in
  [TOML](https://en.wikipedia.org/wiki/TOML) format.
* [Special sections](#special-sections), for abstracts, prefaces or notes.
* [Including other files](#including-files) with the option to specify line ranges, regular
  expressions and/or prefixing each line with a custom string.
* [Document divisions](#document-divisions).
* [Captions](#captions) for code, tables, quotes and subfigures.
* [Asides](#asides).
* [Figures and Subfigures](#figures-and-subfigures) - allows grouping images into subfigures as
  well as giving a single image metadata (a link, attributes, etc.). See [Images in
  Mmark](/syntax/images) for more details.
* [Block Level Attributes](#block-level-attributes) that allow to specify attributes, classes and
  IDs for elements.
* [Indices](#indices) to mark an item (and/or a subitem) to be referenced in the document index.
* [Citations](#citations) and adding [XML References](#xml-references).
* [In document cross references](#cross-references), short form of referencing a section in the
  document.
* [Super- and Subscript](#super-and-subscript)
* [Callouts](#callouts) in code and text.
* [BCP14](#bcp14) (RFC 2119) keyword detection.

### Syntax Gotchas

Because markdown is not perfect, there are some gotchas you have to be aware of:

* Adding a caption under a quote block (`Quote: `) needs a newline before it, otherwise the caption text
  will be detected as being part of the quote.
* Including files (and code includes) requires an empty line before them, as they are block level
  elements and we need to trigger *that* scan from the parser.
* Including files in lists requires a empty line to be present in the list item; otherwise Mmark
  will only assume inline elements and not parse the includes (which are block level elements).
* A bibliography is *only added* if a `{backmatter}` has been specified, because we need to add just
  before that point.
* Intra-work emphasis is enabled so a string like `SSH_MSG_KEXECDH_REPLY` is interpreted as
  `SSH<em>MSG</em>...`. You need to escape the underscores: `SSH\_MSG...`.

### RFC 7991 XML Output

This is the output format used for generating Internet-Drafts and RFCs. The generated XML needs to
be processed by another tool (xml2rfc) to generate to official (final) output. The XML from *mmark*
can be used directly to upload to the IETF tools website.

Title Block:
:   If the document has a [title block](#title-block) the front matter is already open. Closing the
    front matter can only be done by starting the middle matter with `{mainmatter}`. Any open
    "matters" are closed when the document ends. *Area* defaults to "Internet" and *Ipr* defaults to
    `trust200902`.

    Not giving a date will output `<date/>` which mean the current date will be applied *when
    xml2rfc is run*.

Abstract:
:   The abstract can be started by using the special header syntax `.# Abstract`

Note:
:   Any special header that is not "abstract" or "preface" will be a
    [note](https://tools.ietf.org/html/rfc7749#section-2.24): a numberless section.
    These notes are only allowed in the `<front>` section of the document.
    Note [sic] that notes can only contain `<t>` and not other block level elements,
    Mmark will filter these out for: `blockquote` currently (2020 September).

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

Block Level Attributes:
:   We use the attributes as specified in RFC 7991, e.g. to specify an empty list style use:
    `{empty="true"}` before the list. The renderer for this output format filters unknown attributes
    away. The current list is to allow IDs (translated into 'anchor'), remove any `class=` and `style=`
    attributes, so `{style="empty" empty="true"}`, will make a document both RFC 7991 and RFC 7749
    compliant.

Footnotes:
:   Are discarded from the final output, don't use them.

Images:
:   Images are supported. We convert this to an `<artwork>` with `src` set to the image URL of path.
    I.e. `![alt text](img.svg "title")` becomes `<artwork src="img.svg" type="svg" name="title"/>`.
    Note the first `svg` (the alt text) is used as the `type=` attribute. Also note that an image
    like this will be wrapped in `<t>` which is not allowed in RFC 7991 syntax. So to make this
    fully work you need to the image in a subfigure: `!---`. See [Images in Mmark](/syntax/images)
    for more details.

Horizontal Line:
:   Outputs a paragraph with 60 dashes `-`.

Comments:
:   HTML Comments are detected and discarded. These can be useful to make the parser parse certain
    constructs as a block element without meddling with the output.

HTML:
:   The `<br>` tag is detected and converted into a hard break.

### HTML5 Output

Title Block:
:   From the title block only the title is used, in the `<title>` tag.

### Markdown Output

This outputs markdown again, but pretty printed.

### Manual Page Output

Title Block:
:   The title block needs a few elements to correctly generate a manual page

    * `title`, title needs to *end in a digit* to signal the *section*, defaults to "1" if nothing is
      found.
    * `area`, what is it, e.g. "User Commands".
    * `workgroup`, who wrote this e.g.  "Mmark Markdown".
    * `date`, date of the man page, optional, defaults to "today".
    * `author`, to add an Authors section at the end.

Images:
:   See [Images in Mmark](/syntax/images) for details, `ascii-art` images from a sub-figure are
    included.

References and citations:
:   Supported, a "Bibliography" section is added. Note that unlike XML2RFC, references for IDs and
    RFCs *are not* automatically added.

Code Block:
:   Tabs are converted into four spaces.


## Block Elements

### Title Block

A Title Block contains a document's meta data; title, authors, date and other elements. The elements
that can be specified are copied from the [xml2rfc v3
standard](https://tools.ietf.org/html/rfc7991). More on these below. The complete title block is
specified in [TOML](https://github.com/toml-lang/toml). Examples title blocks can be [found in the
repository of Mmark](https://github.com/mmarkdown/mmark/tree/master/rfc).

The title block itself needs three or more `%`'s (or `-`'s) at the start and end of the block. A
minimal title block would look like this:

~~~
%%%
title = "Foo Bar"
%%%
~~~
or

~~~
---
title = "Foo Bar"
---
~~~

The difference between the two is:

* `%%%`: block is assumed to be encoded in TOML and *parsed*.
* `---`: block is not parsed just outputted as-is again (for markdown output), all other formats
  ignore the contents.

#### Elements of the Title Block

An I-D needs to have a Title Block with the following items filled out:

* `title` - the main title of the document.
* `abbrev` - abbreviation of the title.
* `updates/obsoletes` - array of integers.
* `seriesInfo`, containing (this might change with the new *new* XMLv3 output)
   * `name` - `RFC` or `Internet-Draft` or `DOI`
   * `value` - draft name or RFC number
   * `stream` - `IETF` (default), `IAB`, `IRTF` or `independent`.
   * `status` - `standard`, `informational`, `experimental`, `bcp`, `fyi`, or `full-standard`.
* `ipr` - usually just set `trust200902`.
* `area` - usually just `Internet`.
* `workgroup` - the workgroup the document is created for.
* `keyword` - array with keywords (optional).
* `author(s)` - define all the authors.
* `date` - the date for this I-D/RFC.
* `language` - the language for this document, this uses localized names for `Index`, `Footnotes`
  and `References`, etc. Valid values are from [BCP47](https://tools.ietf.org/html/bcp47). This
  defaults to `en` (English). See the [current
  list](https://github.com/mmarkdown/mmark/blob/master/lang/lang.go).

For a manual page the `title`, `area` and `workgroup` are mandatory, if `date` is not specified,
"today" is assumed.

An example would be:

~~~ toml
%%%
title = "Using Mmark to create I-Ds and RFCs"
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

### Special Sections

Any section that needs special handling, like an abstract or preface can be started with `.#
Heading`. This creates a special section that is usually unnumbered.

### Including Files

Including other files can done be with `{{filename}}`, if the path of `filename` is *not* absolute,
the filename is taken relative to *current file being processed*. With `<{{filename}}`
you include a file as a code block. The main difference being it will be returned as a code
block. The file's extension *will be used* as the language. The syntax is:

~~~
{{pathname}}[address]
~~~
And address can be `N,M`, where `N` and `M` are line numbers. If `M` is not specified, i.e. `N,` it
is taken that we should include the entire file starting from `N`.

Or you can use regular expression with: `/N/,/M/`, where `N` and `M` are regular expressions that
specify from where to where to include lines from file.

Each of these can have an optional `prefix=""` specifier.

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
Figure: A sample function.
~~~

Note that because the extension of the file above is "go", this include will lead to the following
block being parsed:

    ~~~ go
    // test.go data
    ~~~
    Figure: A sample function.

### Document Divisions

Mmark support three document divisions, front matter, main matter and the back matter. Mmark
automatically starts the front matter for you *if* the document has a title block. Switching
divisions can be done with `{frontmatter}`, `{mainmatter}` and `{backmatter}`. This must be the only
thing on the line.

### Captions

Mmark supports caption below [tables](#tables), [code blocks](#code-blocks) and [block
quotes](#block-quotes). You can caption each elements with `Table: `, `Figure: ` and `Quote: `
respectively. The caption extends to the first *empty* line. Some examples:

~~~
Name    | Age
--------|-----:
Bob     | 27
Alice   | 23
Table: This is the table caption.
~~~

Or for a code block:

     ~~~ go
     func getTrue() bool {
         return true
     }
     ~~~
     Figure: This is a caption for a code block.

And for a quote:

     > Ability is nothing without opportunity.

     Quote: https://example.com, Napoleon Bonaparte

A caption can potentially contain a "heading ID": `{#id}` as the *last* text in the caption. If this
is found that ID is used as the ID for the entire figure:

~~~
Name    | Age
--------|-----:
Bob     | 27
Alice   | 23
Table: This is the table caption. {#ages}
~~~

### Asides

Any text prefixed with `A>` will become an
[aside](https://developer.mozilla.org/en/docs/Web/HTML/Element/aside). This is similar to a block
quote, but can be styled differently.

### Figures and Subfigures

To *group* artworks and code blocks into figures, we need an extra syntax element. [Scholarly
markdown] has a neat syntax for this. It uses a special section syntax and all images in that
section become subfigures of a larger figure. Disadvantage of this syntax is that it can not be used
in lists. We use a fenced code block like syntax: `!---` as the opening and closing "tag".
Note: only inline elements are parsed inside a figure block.

Basic usage:

~~~
!---
![Alt text](/path/to/img.jpg "Optional title")
!---
~~~

if the figure block has a caption that will be used as well:

~~~
!---
![Alt text](/path/to/img.jpg "Optional title")
![Alt2 text](/path/to/img2.jpg "Optional title2")
!---
Figure: this is a figure containing subfigures.
~~~

Or when just using fenced code blocks:

```
!---
~~~ ascii-art
+-----+
| ART |
+-----+
~~~
Figure: Caption for this ascii-art

~~~ c
printf("%s\n", "hello");
~~~
!---
Figure: Caption for both figures.
```

### Block Level Attributes

A "Block Level Attribute" is a list of HTML attributes between braces: `{...}`. It allows you to
set classes, an anchor and other types of *extra* information for the next block level element.

The full syntax is: `{#id .class key="value"}`. Values may be omitted, i.e., just `{.class}` is
valid.

The following example applies the attributes: `title` and `anchor` to the blockquote:
~~~
{title="The blockquote" #myid}
> A blockquote with a title
~~~
Gets expanded into:
~~~
<blockquote anchor="myid" title="The blockquote">
    <t>A blockquote with a title</t>
</blockquote>
~~~

### Paragraphs

Text that is separated from the rest of the content with empty lines.

### Tables

Tables can be entered by using a simple syntax:

~~~
Name    | Age
--------|------
Bob     | 27
Alice   | 23
~~~

Table footers are supported as well and can be added with equal signs (=):

~~~
Name    | Age
--------|------
Bob     | 27
Alice   | 23
========|======
Total   | 50
~~~

The pipe symbol (`|`) to mark columns does not need to be aligned. Each row must be on a single
line.

## Inline Elements

### Indices

Defining indices allows you to create an index. The define an index use the `(!item)`. Sub items can
be added as well, with `(!item, subitem)`. To make `item` primary, use another `!`: `(!!item,
subitem)`. If any index is defined the end of the document contains the list of indices. The
`-index=false` flag suppresses this generation.

### Citations

Mmark uses the citation syntax from Pandoc: `[@RFC2535]`, the citation can either be informative
(default) or normative, this can be indicated by using the `?` or `!` modifier: `[@!RFC2535]` create
a normative reference for RFC 2535. To suppress a citation use `[@-RFC1000]`. It will still add the
citation to the references, but does not show up in the document as a citation.

The first seen modifier determines the type (suppressed, normative or informative). Multiple
citation can separated with a semicolon: `[@RFC1034;@RFC1035]`.

If you reference an RFC, I-D or W3C document the reference will be added automatically (no need to
muck about with an `<reference>` block). This is to say:

Any reference starting with *RFC*, *I-D.* or *W3C.* will be automatically added to the correct
reference section.

For I-Ds you may want to add a draft sequence number, which can be done as such: `[@?I-D.blah#06]`.
If you reference an I-D *without* a sequence number it will create a reference to the *last* I-D in
citation index.

A bibliography section is created by default if a `{backmatter}` is given, but you can suppress it
by using the command line flag `-bibliography=false`.

### XML References

Any valid XML reference fragment found anywhere in the document, can be used as a citation reference.
The syntax of the XML reference element is defined in [RFC
7749](https://tools.ietf.org/html/rfc7749#section-2.30). The `anchor` defined can be used in the
[citation](#Citations), which the example below that would be `[@pandoc]`:

~~~ html
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

For superscript use `^` and for subscripts use `~`. For example:

~~~
H~2~O is a liquid. 2^10^ is 1024.
~~~

Inside a super- or subscript you must escape spaces. Thus, if you want the letter P with 'a cat' in
subscripts, use `P~a\ cat~`, not `P~a cat~`.

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

Lone callouts (in code blocks) without them being prefixed with a comment means they are not
detected by Mmark.

### BCP14

Phrases that are defined in RFC 2119 (i.e. MUST, SHOULD, etc) are detected when being type set as
strong elements: `**MUST**`, in the RFC 7991 output these will typeset as `<bcp14>MUST</bcp14>`.

# Changes from version 1

These are the changes from Mmark version 1:

* Citations:
   * Suppressing a citation is done with `[@-ref]` (it was the reverse `-@` in v1), this is more consistent.
   * Multiple citations are allowed in one go, separated with a semicolons: `[@ref1; @ref2]`.
   * A reference text suffix is allowed `[@ref, p. 23]`, the separation character is a comma; this
     mirrors the pandoc syntax.
* Indices: now just done with `(!item)`, marking one primary will be: `(!!item)`.
* Code block callouts are now a renderer setting, not a [Block Level
  Attribute](#block-level-attributes). Callout in code are *only* detected if they are used after
  a comment.
* Including files with a prefix is now specified in the address specification:
  `{{myfile}}[prefix="C: "]` will use `C: ` as the prefix. No more mucking about with block
  attribute lists that are hard to discover.
* There no extended table syntax; if this ever comes back it needs to more robust implementation.
* Title Block need to be sandwiched between `%%%`, the prefix `%` does not work anymore.

Syntax that is *not* supported anymore:

* HTML abbreviations.
* The different list syntaxes have been dropped, use a [Block Level
  Attribute](#block-level-attributes) to tweak the output.
* Tasks lists and example lists.
* Comment detection, i.e. to support `cref`: dropped. Comments are copied depending on the output
  renderer.
* Parts
* Extended table syntax.
