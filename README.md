Furlib: Gopher Client/Server for Go
===================================

Furlib is a work-in-progress Gopher client and server library for Go. Furlib
tries to keep things as close as reasonable to the shape of the http stdlib
package.

## Expectation Management

This is a work in progress. The Gopher protocol appears simple, but it's missing
a lot of mod cons and it is taking time, care and effort to work out how to
handle things correctly.

This project is a bit of fun for me, not a product. Issues may be responded to
whenever I happen to get around to them, but PRs are unlikely to be accepted.

## Client quickstart

```go
var client gopher.Client
response, err := client.Fetch(
    context.Background(), gopher.MustParseURL("gopher://gopher.floodgap.com/"))

switch response.(type) {
case *gopher.DirResponse:
    // handle dir
case *gopher.TextResponse:
    // handle text
case *gopher.BinaryResponse:
    // handle bin
}
```

## Server quickstart:

```go
var mux = gopher.NewMux()
mux.Handle("/foo", func(ctx context.Context, w ResponseWriter, r *Request) {
    tw := gopher.NewTextWriter(w, r)
    defer tw.MustFlush()
    tw.WriteString("hello!")
}, nil)

mux.Handle("", func(ctx context.Context, w ResponseWriter, r *Request) {
    dw := gopher.NewDirWriter(w, r)
    defer dw.MustFlush()
    dw.Info("yep")
    dw.Text("Test!", "/foo")
}, nil)

var server = gopher.Server{
    Handler: mux,
}
err := server.ListenAndServe(":7070", "")
```

