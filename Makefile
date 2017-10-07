.PHONY: all deps test build build-cross install clean docker-build docker-push

DOCKER_TAG := $(DOCKER_USER)/dockerlint:$(shell cat version.txt)

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

install: test
	go install

clean:
	rm -f dockerlint
	rm -rf build

docker-build:
	docker build --rm -t $(DOCKER_TAG) .

docker-push: docker-build
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASSWORD)
	docker push $(DOCKER_TAG)