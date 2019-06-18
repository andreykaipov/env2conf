module := $(shell head -n1 go.mod | cut -d' ' -f2)
version := $(shell grep -Eo '\d+\.\d+\.\d+' version.go)
gitsha := $(shell git rev-parse --short HEAD)

default: build

build:
	CGO_ENABLED=0 go build -ldflags='-w -s -X main.version=${version}+dev -X main.gitsha=${gitsha}' -o bin/env2star ${module}

lint:
	@if ! command -v golangci-lint >/dev/null; then\
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b /usr/local/bin latest;\
	fi;\
	golangci-lint run --enable-all ./...
	golint ./...
	@echo All good!

test: build
	go test -v -count=1 ./...
	PATH=$$PWD/bin:$$PATH ./test.sh

release:
	@for os in linux darwin windows; do\
		if [ $$os = "windows" ]; then ext=".exe"; fi;\
		out="bin/release/env2star-${version}-$$os-amd64$$ext";\
		echo "Building $$out";\
		GOOS=$$os GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-w -s -X main.gitsha=${gitsha}' -o $$out ${module};\
	done
	upx --best bin/release/*

clean:
	rm -rf bin
