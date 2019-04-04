VERSION:=$(shell grep 'var Version =' version.go | awk '{ print $$4} ' | tr -d '"' )
GITHUB:=mmarkdown
NAME:=mmark

.PHONY: mmark
mmark:
	@echo $(VERSION)
	go build

mmark.1: mmark.1.md
	pandoc mmark.1.md -s -t man > mmark.1

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
	GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/mmark
	GOOS=linux GOARCH=arm64 go build -o build/linux/arm64/mmark
	GOOS=linux GOARCH=arm go build -o build/linux/arm/mmark
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/amd64/mmark
	GOOS=windows GOARCH=amd64 go build -o build/windows/amd64/mmark.exe

.PHONY: tar
tar:
	rm -rf release && mkdir release
	tar -zcf release/mmark_$(VERSION)_linux_amd64.tgz -C build/linux/amd64 mmark
	tar -zcf release/mmark_$(VERSION)_linux_arm64.tgz -C build/linux/arm64 mmark
	tar -zcf release/mmark_$(VERSION)_linux_arm.tgz -C build/linux/arm mmark
	tar -zcf release/mmark_$(VERSION)_darwin_amd64.tgz -C build/darwin/amd64 mmark
	tar -zcf release/mmark_$(VERSION)_windows_amd64.tgz -C build/windows/amd64 mmark.exe


.PHONY: release
release:
	@echo Releasing: $(VERSION)
	@$(eval RELEASE:=$(shell curl -s -d '{"tag_name": "v$(VERSION)", "name": "v$(VERSION)"}' "https://api.github.com/repos/$(GITHUB)/$(NAME)/releases?access_token=${GITHUB_ACCESS_TOKEN}" | grep -m 1 '"id"' | tr -cd '[[:digit:]]'))
	@echo ReleaseID: $(RELEASE)
	for asset in `ls -A release`; do \
	    curl -o /dev/null -X POST \
	      -H "Content-Type: application/gzip" \
	      --data-binary "@release/$$asset" \
	      "https://uploads.github.com/repos/$(GITHUB)/$(NAME)/releases/$(RELEASE)/assets?name=$${asset}&access_token=${GITHUB_ACCESS_TOKEN}" ; \
	done
