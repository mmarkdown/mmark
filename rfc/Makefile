MMARK :=../mmark
TWO := ""
TXT := $(patsubst %.md,%.txt,$(wildcard *.md))
XML := $(patsubst %.md,%.xml,$(wildcard *.md))

txt: $(TXT)

%.txt: %.md
	if [ -z $(TWO) ]; then \
	    $(MMARK) $< > $(basename $<).xml; \
	    xml2rfc --text --v3 $(basename $<).xml && rm $(basename $<).xml; \
	else \
	    $(MMARK) -2 $< > $(basename $<).xml; \
	    xml2rfc --text $(basename $<).xml && rm $(basename $<).xml; \
	fi

xml: $(XML)

%.xml: %.md
	if [ -z $(TWO) ]; then \
	    $(MMARK) $< > $(basename $<).xml; \
	else \
	    $(MMARK) -2 $< > $(basename $<).xml; \
	fi

.PHONY: clean
clean:
	rm -f *.txt *.xml
