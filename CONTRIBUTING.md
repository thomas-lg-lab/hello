# Contributing

Checks the CI runs, locally with Go 1.27 (or `docker run --rm -v "$PWD:/src" -w /src golang:1.27 sh -c '...'`):

```bash
gofmt -l .        # prints nothing
go vet ./...
go test -race ./...
docker build -t hello .
```
