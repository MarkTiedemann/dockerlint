FROM golang:1.9.1-alpine3.6 AS builder
WORKDIR /go/src/github.com/marktiedemann/dockerlint
RUN apk add --no-cache git make
COPY . .
RUN make install

FROM alpine:3.6
WORKDIR /usr/local/bin
COPY --from=builder /go/bin/dockerlint .
EXPOSE 3000
ENTRYPOINT ["./dockerlint"]
CMD ["-addr=:3000", "-path=/"]