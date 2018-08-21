VERSION:=$(shell grep 'var Version =' mmark.go | awk '{ print $$4} ' | tr -d '"' )
GITHUB:=mmarkdown
NAME:=mmark

.PHONY: mmark
mmark:
	@echo $(VERSION)
	go build


.PHONY: build
build:
	@echo $(VERSION)
	rm -rf build
	mkdir build
	GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/mmark
	GOOS=linux GOARCH=arm64 go build -o build/linux/arm64/mmark
	GOOS=linux GOARCH=arm go build -o build/linux/arm/mmark
	GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/mmark
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/amd64/mmark
	GOOS=windows GOARCH=amd64 go build -o build/windows/amd64/mmark.exe


.PHONY: release
release:
	@echo Releasing: $(VERSION)
	@$(eval RELEASE:=$(shell curl -s -d '{"tag_name": "v$(VERSION)", "name": "v$(VERSION)"}' "https://api.github.com/repos/$(GITHUB)/$(NAME)/releases?access_token=${GITHUB_ACCESS_TOKEN}" | grep -m 1 '"id"' | tr -cd '[[:digit:]]'))
	@echo ReleaseID: $(RELEASE)
	@for asset in `find build -type f`; do \
	    curl -o /dev/null -X POST \
	      -H "Content-Type: application/binary" \
	      --data-binary "$$asset" \
	      "https://uploads.github.com/repos/$(GITHUB)/$(NAME)/releases/$(RELEASE)/assets?name=$${asset}&access_token=${GITHUB_ACCESS_TOKEN}" ; \
	done

