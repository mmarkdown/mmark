---
title: "Images in Mmark"
date: 2020-10-12T10:05:51+01:00
aliases: [/syntax/images]
toc: true
---

Images in Mmark are somewhat complicated, not in the least, because XML2RFC needs to output both
HTML and text. To make that work you can specify multiple images in an `artset` and the renderer will
pick the correct one, depending on the output.

To include artwork/source/images, you can:

* Use a code block. If this has a language specified, it will become a `sourcecode` otherwise a will
  be an `artwork`. The contents of both must be in plain text.
* Use subfigures in a figure block (`!--`) to group figure and potentially make them have a caption.
  We also use this syntax to support an `artset`, but only under special conditions (see below).

## Code Blocks

A code block *with* a language will be turned into a `sourcecode`:

    ``` go
    println("hello!")
    ```

If no language (the `go` above) is given it will be an `artwork`.

## Figures in an Artset

To support `artset` we do the following. If multiple images are present as subfigures, we check
if the name *without* the extension of the image destination (the file to be shown) is equal for all
subfigures. If so, we assume an `artset` needs to be outputted and do so.

For example the following will result in a artset where an `svg` and an `ascii-art` version of the
(hopefully) same image exists. Note the extension **must** be `ascii-art` because we use that to set
the type and XML2RFC checks for that string.

~~~
!---
![Array vs Slice](array-vs-slice.svg "Title of the svg image")
![Array vs Slice](array-vs-slice.ascii-art "Title of the ascii-art image")
!---
~~~

Note this syntax is also supported for the *manual page output* and it does the same thing by only
using the `ascii-art` version. This is true for all included imagery; only `ascii-art` ones are
included in the output.

By some happy co-incidence a browser will not show the `ascii-art` version of the image when
generating HTML. It remains to be seen if we need some code to actually filter these out.
