FROM golang:1.9.1-alpine3.6 AS builder
WORKDIR src/github.com/marktiedemann/dockerlint
RUN apk add --no-cache git
RUN go get -d -v github.com/docker/docker/builder/dockerfile/instructions
RUN go get -d -v github.com/docker/docker/builder/dockerfile/parser
COPY . .
RUN go test
RUN go build

FROM alpine:3.6
COPY --from=builder go/src/github.com/marktiedemann/dockerlint/dockerlint .
CMD ["./dockerlint"]