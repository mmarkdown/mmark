VERSION:=$(shell grep 'var Version =' version.go | awk '{ print $$4} ' | tr -d '"' )
GITHUB:=mmarkdown
NAME:=mmark
LINUX_ARCH:=amd64 arm64 arm
DARWIN_ARCH:=amd64 arm64

.PHONY: mmark
mmark:
	@echo $(VERSION)
	CGO_ENABLED=0 go build

define DOCHEADER
%%%
title = "mmark 1"
date = "2019-04-04T19:23:50+01:00"
area = "User Commands"
workgroup = "Mmark Markdown"
%%%
endef

define SYNHEADER
%%%
title = "mmark-syntax 7"
date = "2019-04-04T19:23:50+01:00"
area = "User Commands"
workgroup = "Mmark Markdown syntax"
%%%
endef

mmark.1: mmark.1.md
	$(file > mmark.1.docheader,$(DOCHEADER))
	( cat mmark.1.docheader ; tail +8 mmark.1.md ) | ./mmark -man > mmark.1 && rm -f mmark.1.docheader

mmark-syntax.7: Syntax.md
	$(file > mmark-syntax.7.synheader,$(SYNHEADER))
	( cat mmark-syntax.7.synheader ; tail +8 Syntax.md ) | ./mmark -man > mmark-syntax.7 && rm -f mmark-syntax.7.synheader

mmark-syntax-images.7: Syntax-images.md
	$(file > mmark-syntax.7.synheader,$(SYNHEADER))
	( cat mmark-syntax.7.synheader ; tail +8 Syntax-images.md ) | ./mmark -man > mmark-syntax-images.7 && rm -f mmark-syntax.7.synheader

.PHONY: clean
clean:
	rm -rf build release
	$(MAKE) -C rfc clean

.PHONY: man
man: mmark.1 mmark-syntax.7 mmark-syntax-images.7

# up the version in version.go and 'make release'
.PHONY: release
release:
	git ci -am"Version $(VERSION)"
	git tag v$(VERSION)
	git push --tags
	git push
