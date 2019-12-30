---
title: "FAQ"
date: 2018-07-22T14:05:51+01:00
aliases: [/faq/]
toc: true
---

Mmark Frequently Asked Questions.

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

# How Do I Make an Author an Editor

Use `role = "editor"` in the author's section in the titleblock.
