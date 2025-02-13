# kbbi

[![godoc](https://godoc.org/github.com/raf555/kbbi-api/pkg/kbbi?status.svg)][godoc]

kbbi package contains structs that are directly used by the [kbbi-api][] API server response.
Clients can expect the latest version of this package to be always the same as the response used by the API respective to the version.

E.g. `/api/v1/*` API uses `v1.x.y-z` package.  

[godoc]: https://godoc.org/github.com/raf555/kbbi-api/pkg/kbbi
[kbbi-api]: https://github.com/raf555/kbbi-api

## Installing

You can import directly to your application.

```go
import "github.com/raf555/kbbi-api/pkg/kbbi"
```

Or use `go get`

```sh
go get github.com/raf555/kbbi-api/pkg/kbbi
```
