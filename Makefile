VERSION:=$(shell grep 'var Version =' version.go | awk '{ print $$4} ' | tr -d '"' )
GITHUB:=mmarkdown
NAME:=mmark
LINUX_ARCH:=amd64 arm64 arm
DARWIN_ARCH:=amd64 arm64

.PHONY: mmark
mmark:
	@echo $(VERSION)
	go build

define DOCHEADER
%%%
title = "mmark"
date = "2019-04-04T19:23:50+01:00"
area = "User Commands"
workgroup = "Mmark Markdown"
%%%
endef

mmark.1: mmark.1.md
	$(file > mmark.1.docheader,$(DOCHEADER))
	( cat mmark.1.docheader ; tail +8 mmark.1.md ) | ./mmark -man > mmark.1
	rm -f mmark.1.docheader

.PHONY: clean
clean:
	rm -rf build release
	$(MAKE) -C rfc clean

.PHONY: git-commit
git-commit:
	git ci -am"Version $(VERSION)"
	git push --tags
	git push

.PHONY: build
build:
	@echo $(VERSION)
	rm -rf build && mkdir build
	for arch in $(LINUX_ARCH); do GOOS=linux GOARCH=$$arch go build -o build/linux/$$arch/mmark; done
	for arch in $(DARWIN_ARCH); do GOOS=darwin GOARCH=$$arch go build -o build/darwin/$$arch/mmark; done
	GOOS=windows GOARCH=amd64 go build -o build/windows/amd64/mmark.exe

.PHONY: tar
tar:
	rm -rf release && mkdir release
	for arch in $(LINUX_ARCH); do tar -zcf release/mmark_$(VERSION)_linux_$$arch.tgz -C build/linux/$$arch mmark; done
	for arch in $(DARWIN_ARCH); do tar -zcf release/mmark_$(VERSION)_darwin_$$arch.tgz -C build/darwin/$$arch mmark; done
	tar -zcf release/mmark_$(VERSION)_windows_amd64.tgz -C build/windows/amd64 mmark.exe

.PHONY: release
release:
	@echo Releasing: $(VERSION)
	@$(eval RELEASE:=$(shell curl -s -d '{"tag_name": "v$(VERSION)", "name": "v$(VERSION)"}'  -H "Authorization: token ${GITHUB_ACCESS_TOKEN}" "https://api.github.com/repos/$(GITHUB)/$(NAME)/releases" | grep -m 1 '"id"' | tr -cd '[[:digit:]]'))
	@echo ReleaseID: $(RELEASE)
	for asset in `ls -A release`; do \
	    curl -o /dev/null -X POST \
	      -H "Content-Type: application/gzip" \
	      -H "Authorization: token ${GITHUB_ACCESS_TOKEN}" \
	      --data-binary "@release/$$asset" \
	      "https://uploads.github.com/repos/$(GITHUB)/$(NAME)/releases/$(RELEASE)/assets?name=$${asset}" ; \
	done

.PHONY: debian
debian: mmark.1 mmark
	export MY_APP_VERSION=$(VERSION)
	nfpm -f .nfpm.yaml pkg -t mmark.deb
