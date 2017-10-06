# dockerlint

**Linting `Dockerfile`s as a service.**

## Quickstart

```sh
# clone repo
git clone https://github.com/marktiedemann/dockerlint && cd dockerlint

# build image
docker build --rm -t dockerlint:0.2.0 .

# run container
docker run --rm -it -p 3000:3000 dockerlint:0.2.0
```

## API

**Base URL:** `http://localhost:3000`

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

**Headers:**

- `Content-Type`: `application/json`

**Body:** *json*, the linting result

- `error`: *boolean*, `true` if a linting error was detected; otherwise `false`
- `message`: *string | undefined*, the error message, if a linting error was detected (*optional*)

*Example:*

- *200*:

```json
{
  "error": false
}
```

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
 - [ ] Support `-addr` and `-path` flags and env vars
 - [ ] Setup CI
 - [ ] Upload image to Docker Hub
 - [ ] Host demo instance

## License

[WTFPL](http://www.wtfpl.net/) – Do What the F*ck You Want to Public License.

Made with :heart: by [@MarkTiedemann](https://twitter.com/MarkTiedemannDE).