.PHONY: all deps dev-deps update-deps lint test cover build build-cross install clean docker-build docker-push

DOCKER_TAG := $(DOCKER_USER)/dockerlint:$(shell cat version.txt)
DOCKER_TAG_LATEST := $(DOCKER_USER)/dockerlint:latest

all: build

deps:
	go get $(GO_GET) github.com/docker/docker/builder/dockerfile/parser
	go get $(GO_GET) github.com/docker/docker/builder/dockerfile/instructions

dev-deps: deps
	go get $(GO_GET) github.com/golang/lint/golint
	go get $(GO_GET) github.com/mitchellh/gox

update-deps:
	$(MAKE) dev-deps GO_GET=-u

lint: dev-deps
	go fmt
	go vet
	golint

test: lint
	go test

cover: lint
	go test -cover -coverprofile=.coverprofile
	go tool cover -html=.coverprofile

build: test
	go build

build-cross: test clean
	mkdir build
	gox -os="windows darwin linux" -arch="386 amd64" -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}"

install: test
	go install

clean:
	rm -f dockerlint
	rm -f .coverprofile
	rm -rf build

docker-build:
	docker build -t $(DOCKER_TAG) -t $(DOCKER_TAG_LATEST) .

docker-push: docker-build
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASSWORD)
	docker push $(DOCKER_TAG)
	docker push $(DOCKER_TAG_LATEST)