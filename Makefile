version := $(shell grep -Eo '\d+\.\d+\.\d+' ./src/version.go)
gitsha := $(shell git rev-parse --short HEAD)

default: build

build:
	CGO_ENABLED=0 go build -ldflags='-w -s -X main.gitsha=${gitsha}' -o bin/env2star ./src

install:
	install bin/env2star /usr/local/bin/env2star

lint:
	@if ! command -v golangci-lint >/dev/null; then\
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b /usr/local/bin latest;\
	fi;\
	golangci-lint run --enable-all ./src
	golint ./src
	@echo All good!

test: build install
	go test -v -count=1 ./src
	PATH=$$PWD/bin:$$PATH ./test.sh

release:
	@for os in linux darwin windows; do\
		if [ $$os = "windows" ]; then ext=".exe"; fi;\
		out="bin/release/env2star-${version}-$$os-amd64$$ext";\
		echo "Building $$out";\
		GOOS=$$os GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-w -s -X main.gitsha=${gitsha}' -o $$out ./src;\
	done
	upx --best bin/release/*

clean:
	rm -rf bin
