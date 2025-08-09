# go-webcrawler

This repository contains the code of a web crawler which analyzes given URL.

## Installation

Clone and build:

```bash
git clone https://github.com/dorukco/go-webcrawler.git
cd go-webcrawler
go build
```

## Development

Run tests:
```bash
go test -v ./...
```

Run locally:
```bash
go run main.go
```

When running the application locally, it is located at `http://localhost:8080/`.

## Building the Docker image

* Run `docker build -t webcrawler .` to build a docker image.
* Run `docker run -d -p 8080:8080 webcrawler` to run the docker image.

## License

MIT