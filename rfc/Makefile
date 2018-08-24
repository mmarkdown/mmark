MMARK :=../mmark
TXT := $(patsubst %.md,%.txt,$(wildcard *.md))
XML := $(patsubst %.md,%.xml,$(wildcard *.md))

txt: $(TXT)

%.txt: %.md
	$(MMARK) $< > $(basename $<).xml
	xml2rfc --text --v3 $(basename $<).xml && rm $(basename $<).xml

xml: $(XML)

%.xml: %.md
	$(MMARK) $< > $(basename $<).xml

.PHONY: clean
clean:
	rm -f *.txt *.xml
