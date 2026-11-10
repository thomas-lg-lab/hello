# hello

[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/thomas-lg-lab/hello/badge)](https://scorecard.dev/viewer/?uri=github.com/thomas-lg-lab/hello)

The smallest lab app: it greets the user the proxy signed in (`Remote-User` header) and counts the visits in
`/data/visits`. It exists to prove the lab's whole path, from a commit to a running deployment behind the proxy.

| | |
| --- | --- |
| HTTP | port 8080; `GET /` answers `Bonjour <user>`, the visit count and the image version, `GET /healthz` answers `ok` |
| Data | `/data` (override with `DATA_DIR`); empty on first start is fine |
| Secrets | none |
| Image | `ghcr.io/thomas-lg-lab/hello`, non-root (uid 10001), static binary on distroless |

The CI runs `gofmt`, `go vet` and the tests, then builds the image and runs the lab's runtime contract test on it.
