# Agents

## Pipeline

[`.github/workflows/ci.yml`](.github/workflows/ci.yml) is the source of truth for CI. A PR must pass **three** jobs before merge.

| CI job | What it runs | Local equivalent |
|--------|----------------|------------------|
| `check` | `make check` | `make check` |
| `build` | `go build -v ./...` | `go build ./...` |
| `lint` | golangci-lint **v2.11.4** | `golangci-lint run ./...` on that version |

**`make check` is not the whole pipeline.** It covers tracked `gofmt`, `go vet`, and `go test -count=1` only.

Before you say lint is green or the pipeline will pass, run `make check`, `go build ./...`, and `golangci-lint run ./...` (match the workflow version when you can).

Release tags (`v*`) run [`.github/workflows/release.yml`](.github/workflows/release.yml) (GoReleaser + GPG).

The open-unifi **controller** lives in [github.com/lucavb/open-unifi](https://github.com/lucavb/open-unifi). Admin API contract tests run there in `integration/provider/`.
