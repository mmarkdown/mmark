MMARK :=../mmark
TXT := $(patsubst %.md,%.md.txt,$(wildcard *.md))
XML := $(patsubst %.md,%.md.xml,$(wildcard *.md))

txt: $(TXT)

%.md.txt: %.md
	$(MMARK) $< > $<.xml
	xml2rfc --text --v3 $<.xml && rm $<.xml

xml: $(XML)

%.md.xml: %.md
	$(MMARK) $< > $<.xml

.PHONY: clean
clean:
	rm -f *.md.txt *md.xml
