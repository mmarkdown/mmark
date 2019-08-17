---
title: "FAQ"
date: 2018-07-22T14:05:51+01:00
aliases: [/faq/]
toc: true
---

Mmark Frequently Asked Questions.

# How Do I Create an Independent IETF submission?

Use the following as starting point for your title block, `ipr` and `submissiontype` or the important
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
