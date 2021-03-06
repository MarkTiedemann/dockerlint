# dockerlint [![Build Status](https://travis-ci.org/MarkTiedemann/dockerlint.svg?branch=master)](https://travis-ci.org/MarkTiedemann/dockerlint)

**Linting `Dockerfile`s as a service.**

- *Allows you to lint `Dockerfile`s with a simple HTTP API, without having `docker` installed.*
- *Uses `docker`'s [`dockerfile/parser`](https://github.com/moby/moby/blob/master/builder/dockerfile/parser/parser.go) and [`dockerfile/instructions/parser`](https://github.com/moby/moby/blob/master/builder/dockerfile/instructions/parse.go) internally.*
- *100% test coverage.*

## Installation

```sh
go get github.com/marktiedemann/dockerlint
```

**From source:**

```sh
git clone https://github.com/marktiedemann/dockerlint
cd dockerlint
make install
```

**From Docker Hub ([marktiedemann/dockerlint](https://hub.docker.com/r/marktiedemann/dockerlint/)):**

```sh
docker pull marktiedemann/dockerlint
```

## Flags

- `-addr`: the address of the server, *default:* `:3000`
- `-path`: the path of the handler, *default:* `/`

## API

### Request `POST /`

**Body:** *plain text / binary*, the content of the `Dockerfile`

*Example:*

```
FROM golang
```

#### Response

**Status:**

- `200 OK`: linting succeeded
- `400 Bad Request`: linting failed
- `404 Not Found`
- `405 Method Not Allowed`

**Headers:**

- `Content-Type`: `application/json`

**Body:** *json*, the linting result

- `error`: *boolean*, `true` if a linting error was detected; otherwise `false`
- `message`: *string | undefined*, the error message, if a linting error was detected (*optional*)

*Example:*

`curl --data-binary "FROM golang" localhost:3000`

- *200*:

```json
{
  "error": false
}
```

`curl --data-binary "FROM" localhost:3000`

- *400*:

```json
{
  "error": true,
  "message": "Dockerfile parse error line 1: FROM requires either one or three arguments"
}
```

## Todos

 - [ ] Implement proper content-type negotiation
 - [ ] Support `plain/text` content-type
 - [ ] Host demo instance

## License

[WTFPL](http://www.wtfpl.net/) – Do What the F*ck You Want to Public License.

Made with :heart: by [@MarkTiedemann](https://twitter.com/MarkTiedemannDE).