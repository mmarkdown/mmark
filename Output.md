---
title: "How to Input Mmark"
date: 2018-08-31T07:05:51+01:00
aliases: [/ouput/]
toc: true
---

# Block Level Elements

## Headings

To format a line as a heading, use `#`, the number of hashes determines the level.

### You Type

~~~
# Introduction
~~~

### You Get

html

~~~
<h1 id="introduction">Introduction</h1>
~~~

xml2

~~~
<section anchor="introduction" title="Introduction">
</section>
~~~

xml3

~~~
<section anchor="introduction"><name>Introduction</name>
</section>
~~~
