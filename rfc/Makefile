MMARK :=../mmark
TXT := $(patsubst %.md,%.txt,$(wildcard *.md))
XML := $(patsubst %.md,%.xml,$(wildcard *.md))
HTML := $(patsubst %.md,%.html,$(wildcard *.md))
PDF := $(patsubst %.md,%.pdf,$(wildcard *.md))

txt: $(TXT)

%.txt: %.xml
	xml2rfc --text --v3 $<
	rm -f $(basename $<).xml

html: $(HTML)

%.html: %.xml
	xml2rfc --html --v3 $<
	rm -f $(basename $<).xml

xml: $(XML)

%.xml: %.md
	$(MMARK) $< > $(basename $<).xml

pdf: $(PDF)

%.pdf: %.xml
	xml2rfc --pdf --v3 $<
	rm $(basename $<).xml

.PHONY: clean
clean:
	rm -f *.txt *.xml
