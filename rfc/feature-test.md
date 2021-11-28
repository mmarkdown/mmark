%%%
# This "RFC" uses a lot of xml2rfc features, so that we can see a) mmark correctly parses it and b)
# xml2rfc makes a document out of it.
title = "The Naming of Hosts"
abbrev = "The Naming of Hosts"
ipr = "trust200902"
area = "Internet"
workgroup = "Network Working Group"
submissiontype = "IETF"
keyword = [""]
tocdepth = 5
#date = 1997-04-01T00:00:00Z

[seriesInfo]
name = "RFC"
value = "2100"
stream = "IETF"
status = "informational"

[[author]]
initials="J. R."
surname="Ashworth"
fullname="Jay R. Ashworth"
abbrev = "Ashworth & Associates"
organization = "Advanced Technology Consulting"
  [author.address]
  email = "jra@scfn.thpl.lib.fl.us"
  phone = "+1 813 790 7592"
  [author.address.postal]
  city = "St. Petersburg"
  code = "FL 33709-4819"
  pobox = "aaaa"
  cityarea = "bbbb"

[[contact]]
initials="D."
surname="Addison"
fullname="David Addison"
  [contact.address.postal]
  city = "St. Petersburg"
  code = "FL 33709-4819"
  pobox = "aaaa"
  cityarea = "bbbb"
%%%

{mainmatter}

# Introduction

This RFC is a commentary on the difficulty of deciding upon an acceptably
distinctive hostname for one's computer, a problem which grows in direct
proportion to the logarithmically increasing size of the Internet.

Distribution of this memo is unlimited.

Except to TS Eliot.

And, for that matter, to [@David Addison], who hates iambic pentameter.

# Features

## Definition List

First Term
: This is the definition of the first term.

Second Term
: This is one definition of the second term.

## Ordered List

1) First item.
2) Second item.

## Lists in Lists

Foo validator:
: It performs the following actions:
  * runs
  * jumps
  * walks

Another example:

{type="Step %d:"}
1. Send it to
   * Alice
   * Bob
   * Carol

# Credits

[@David Addison]

[@-libes; @-lottor; @-wong; @-ts]

# Security Considerations

Security issues are not discussed in this memo.

Particularly the cardiac security of certain famous poets.

{backmatter}

<reference anchor='libes' target=''>
 <front>
 <title>Choosing a Name for Your Computer</title>
  <author initials='D.' surname='Libes' fullname='D. Libes'></author>
  <date year='1989' month='November'/>
 </front>
 <seriesInfo name="Communications of the ACM" value='Vol. 32, No. 11, Pg. 1289' />
 </reference>

<reference anchor='lottor' target='namedroppers@internic.net'>
 <front>
 <title>Domain Name Survey</title>
  <author initials='M.' surname='Lottor' fullname='M. Lottor'></author>
  <date year='1997' month='January'/>
 </front>
 </reference>

<reference anchor='wong' target='http://www.seas.upenn.edu/~mengwong/coolhosts.html'>
 <front>
 <title>Cool Hostnames</title>
  <author initials='M.' surname='Wong' fullname='M. Wong'></author>
  <date/>
 </front>
 </reference>

<reference anchor='ts' target=''>
 <front>
 <title>Old Possum's Book of Practical Cats</title>
  <author initials='TS' surname='Stearns' fullname='TS. Stearns'></author>
  <date/>
 </front>
 </reference>
