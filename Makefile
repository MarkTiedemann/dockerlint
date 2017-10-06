.PHONY: all deps test build build-cross build-docker install clean

all: build

deps:
	go get github.com/docker/docker/builder/dockerfile/parser
	go get github.com/docker/docker/builder/dockerfile/instructions

test: deps
	go fmt
	go test

build: test
	go build

build-cross: test clean
	go get github.com/mitchellh/gox
	mkdir build
	gox -os="windows darwin linux" -arch="386 amd64" -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}"

build-docker: test
	docker build --rm -t dockerlint .

install: test
	go install

clean:
	rm -f dockerlint
	rm -rf build