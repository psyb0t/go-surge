# go-surge

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-surge.svg)](https://pkg.go.dev/github.com/psyb0t/go-surge)
[![CI](https://github.com/psyb0t/go-surge/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-surge/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-surge/badges/version.svg)](https://github.com/psyb0t/go-surge/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-surge/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-surge/badges/importers.svg)](https://github.com/psyb0t/go-surge/blob/badges/importers.md)

Unofficial Go client for the [Surge SMS API](https://surge.app), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from Surge's public OpenAPI document. I am not affiliated with Surge.

The package is `surge`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-surge
```

## Usage

Every REST endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200`, `JSON201` and the other status fields.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-surge"
)

func main() {
	auth := surge.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set(
				"Authorization", "Bearer "+os.Getenv("SURGE_API_KEY"),
			)

			return nil
		},
	)

	client, err := surge.NewClientWithResponses("https://api.surge.app", auth)
	if err != nil {
		panic(err)
	}

	res, err := client.ListAccountsWithResponse(
		context.Background(),
		&surge.ListAccountsParams{},
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles. Check `StatusCode()` and the expected JSON field before using a response. HTTP errors are decoded responses, not Go errors; the Go error covers transport and decoding failures.

The full API, with the spec's descriptions as doc comments, is in the [Go reference](https://pkg.go.dev/github.com/psyb0t/go-surge) and in `go doc github.com/psyb0t/go-surge` from inside your module.

The client covers all 41 operations across 25 paths: accounts, audiences, blasts, campaigns, contacts, messages, phone numbers, recordings, users, tokens and verification codes. The component models include webhook payloads, but the package does not run a webhook server or verify webhook signatures.

## Auth

Surge uses HTTP bearer authentication. Set `SURGE_API_KEY` to your API token and attach `Authorization: Bearer <token>` with the `WithRequestEditorFn` shown above. The base URL is `https://api.surge.app`. Account-scoped endpoints take an account ID separately from the token.

List operations use cursor pagination. Pass `Pagination.NextCursor` as `After` on the next request until the cursor is nil. The client does not paginate or retry automatically. Set a timeout through `WithHTTPClient` when needed.

Sending messages queues delivery asynchronously. A successful creation response does not mean the carrier delivered the message. Use message status or delivery webhooks to track delivery. `MessageParams` is a union: use its generated `FromMessageParams0` or `FromMessageParams1` helper to select the conversation or direct recipient form. Supply body text or attachments, and do not mix `conversation` with `to` or `from`.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads the spec the Surge API serves at `https://api.surge.app/openapi.json` to `api/openapi.yaml`. The content is JSON despite the `.yaml` file name.
- `scripts/patch-spec.sh` writes `api/openapi.fixed.yaml`: nullable type lists become `nullable: true`, `const` becomes a single-value enum, schema examples arrays become `example`, and the OpenAPI version becomes 3.0.3 for oapi-codegen. Response examples maps remain intact.
- oapi-codegen runs with `config.yaml`. The config maps `number` to `float64` and retains all component schemas, including webhook models.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Agent integrations

The [skill](.agents/skills/go-surge) teaches a coding agent how to use go-surge when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-surge@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-surge@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The release workflow publishes the skill to ClawHub:

```bash
openclaw skills install @psyb0t/go-surge
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

`make sec` scans Git history when Git metadata is present. Source archives still receive the full directory scan.

The spec's `TokenResponse` example contains a sample JWT. Gitleaks allows that exact value only in `client.gen.go`, where oapi-codegen copies the example into a comment.

## License

[MIT](LICENSE). The generated code comes from Surge's public API specification. The specification itself is not redistributed here. The spec supplies no license or terms URL; the official Python SDK's Apache-2.0 license does not establish a license for the separately hosted specification.
