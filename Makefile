module := $(shell head -n1 go.mod | cut -d' ' -f2)
version := $(shell grep -Eo '[0-9]+\.[0-9]+\.[0-9]+' version.go)
gitsha := $(shell git rev-parse --short HEAD)

default: build

build:
	CGO_ENABLED=0 go build -ldflags='-w -s -X main.version=${version}+dev -X main.gitsha=${gitsha}' -o bin/env2star ${module}

lint:
	golangci-lint run --enable-all ./...
	if command -v golint >/dev/null; then golint ./...; else echo "golint isn't installed; skipping it"; fi
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

compress:
	upx --best bin/release/*

clean:
	rm -rf bin
