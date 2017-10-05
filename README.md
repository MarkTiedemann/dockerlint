# dockerlint

**Linting `Dockerfile`s as a service.**

## Quickstart

```sh
# clone repo
git clone https://github.com/marktiedemann/dockerlint && cd dockerlint

# build image
docker build --rm -t dockerlint .

# run container
docker run --rm -it -p 3000:3000 dockerlint
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

- `error`: *boolean*, whether linting errors were detected
- `messages`: *string[]*, list of linting error messages

*Example:*

- *200*:

```json
{
  "error": false,
  "messages": []
}
```

- *400*:

```json
{
  "error": true,
  "messages": [
    "Dockerfile parse error line 1: FROM requires either one or three arguments"
  ]
}
```

## License

[WTFPL](http://www.wtfpl.net/) – Do What the F*ck You Want to Public License.

Made with :heart: by [@MarkTiedemann](https://twitter.com/MarkTiedemannDE).