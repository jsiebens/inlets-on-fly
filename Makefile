SHELL := bash
Version := $(shell git describe --tags --dirty)
# Version := "dev"
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X github.com/jsiebens/inlets-on-fly/pkg/cmd.Version=$(Version) -X github.com/jsiebens/inlets-on-fly/pkg/cmd.GitCommit=$(GitCommit)"
.PHONY: all

.PHONY: build
build:
	go build -ldflags $(LDFLAGS)

.PHONY: dist
dist:
	mkdir -p dist
	GOOS=linux go build -ldflags $(LDFLAGS) -o dist/inlets-on-fly
	GOOS=darwin go build -ldflags $(LDFLAGS) -o dist/inlets-on-fly-darwin
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -o dist/inlets-on-fly-armhf
	GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -o dist/inlets-on-fly-arm64
	GOOS=windows go build -ldflags $(LDFLAGS) -o dist/inlets-on-fly.exe

.PHONY: hash
hash:
	for f in dist/inlets-on-fly*; do shasum -a 256 $$f > $$f.sha256; done