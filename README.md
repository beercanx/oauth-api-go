# A Go based OAuth API

An exercise into how to create an HTTP service using GO, following guidance from:
* https://go.dev/doc/modules/layout
* https://github.com/golang-standards/project-layout/blob/master/README.md
* https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years

Then gave up doing vanilla Go HTTP server; because I'm a wimp and used a framework instead:
* https://github.com/gin-gonic/gin

Went back a day later and continued to re-read the "after 13 years" blog post to create the `Hello, World!` http sample. 

Regardless, I have decided to continue with `gin` itself, but following the coding styles suggested; such as not storing dependencies in a `struct` but passing them through the functions. 

## Requirements
* Go 1.24

## Testing

The standard Go approach to unit testing
```bash
go test
```

## Building

### Hello, World!
A hello world HTTP server, using just the built-in GO [http](https://pkg.go.dev/net/http) library.
```bash
go build ./cmd/http
```
```bash
./http
```

### The OAuth server
An OAuth server, using mostly the [Gin web framework](https://github.com/gin-gonic/gin) that wraps around the Go [http](https://pkg.go.dev/net/http) library.
```bash
go build ./cmd/server
```
```bash
./server
```
