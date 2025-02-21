# A Go based OAuth API

An exercise into how to create an HTTP service using GO, following guidance from:
* https://go.dev/doc/modules/layout
* https://github.com/golang-standards/project-layout/blob/master/README.md
* https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years

Then gave up doing vanilla Go HTTP server; because I'm a wimp and used a framework instead:
* https://github.com/gin-gonic/gin

## Requirements
* Go 1.24

## Building

### Hello, World!
```bash
go build ./cmd/hello-world
```
```bash
./hello-world
```

### The OAuth server
```bash
go build ./cmd/server
```
```bash
./server
```
