## Semver Compare

### How To Run

#### Docker

In order to run this through a Docker container, you must have [Docker](https://docker.com)
installed.

Once installed, simply run the following command:

```
$ docker-compose up http-server
```

You can choose to run the application in the background by running this
instead:

```
$ docker-compose up -d http-server
```

#### Baremetal

To run in baremental (directly on your Host machine) you must have
[go](https://golang.org) installaed on your machine. Once you have that
installed, run the following command:

```
$ go get ./...
$ go run main.go
```

The application will be accessible at http://localhost:5656

### Code Architecture

### REST Endpoint Documentation

PATH: `/compare-versions`
Method: `POST`

JSON Schema
```
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "description": "Schema to compare semantic versions",
  "type": "object",
  "properties": {
    "data": {
      "type": "object",
      "properties": {
        "compare_from": {
          "description": "Semantic version to be compared from",
          "type": "string",
          "pattern": "^([0-9]+)(\.[0-9]+)?(\.[0-9]+)?(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$",
          "minLength": 1
        },
        "compare_to": {
          "description": "Semantic version to be compared to",
          "type": "string",
          "pattern": "^([0-9]+)(\.[0-9]+)?(\.[0-9]+)?(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$",
          "minLength": 1
        }
      },
      "required": [
        "compare_from",
        "compare_to"
      ],
      "additionalProperties": false
    }
  },
  "required": ["data"],
  "additionalProperties": false
}
```

Sample Request
```
{
  "data": {
    "compare_from": "1.0",
    "compare_to": "1.0.0"
  }
}
```

Sample Resposne
```
{
  "code": "app.compare_success",
  "data": {
      "result": "1.0 is \"equal\" to 1.0.0"
  }
}
```

### Not Handled

This application's endpoint does not handle "pre-releases" in the comparison

### Code Architecture

The pattern followed by the handlers is inspired from the given [article](https://blog.questionable.services/article/http-handler-error-handling-revisited/).
