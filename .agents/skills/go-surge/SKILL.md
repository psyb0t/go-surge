---
name: go-surge
description: Unofficial Go client for the Surge SMS API. Package surge (github.com/psyb0t/go-surge) provides typed REST methods for accounts, audiences, blasts, campaigns, contacts, messages, phone numbers, recordings, users, tokens and verifications. Use when Go code needs Surge SMS operations or already imports go-surge.
homepage: https://github.com/psyb0t/go-surge
user-invocable: true
metadata:
  openclaw:
    emoji: "✉️"
    requires:
      bins:
        - go
      env:
        - SURGE_API_KEY
permissions:
  filesystem: "Guides edits to the current project's Go source. The library reads and writes no files."
  shell: "The Go toolchain to add and compile the dependency."
  network: "Module downloads use the configured Go proxy. Runtime requests go to https://api.surge.app."
---

# go-surge

Unofficial, typed Go client for the Surge SMS API, in package `surge` at the module root.

For install, operations, pagination and responses, read [references/setup.md](references/setup.md).

## Security & safety

- Read the bearer token from `SURGE_API_KEY`. Set `Authorization: Bearer <token>` with `WithRequestEditorFn`. Never hardcode, commit or log the token or authorization header.
- Use `https://api.surge.app`. Most resource creation and list methods also need an account ID; the token and account ID are separate inputs.
- Sending messages or blasts, creating verification codes, buying or releasing phone numbers, and changing accounts or contacts affect real resources. Perform them only when requested.
- The models include webhook payloads. The client does not verify webhook signatures or run a webhook server.

## When to use

- Go code needs to send or read Surge messages, manage contacts, campaigns or phone numbers, or check verification codes.
- Go code needs account, user, audience, blast or recording operations.
- The project already imports `github.com/psyb0t/go-surge`.

## When not to use

- The code is not Go.
- You need an official provider SDK. This module is unofficial.
- You want to modify the generated client. Read `AGENTS.md` in the repository.

## Usage

Use the compiled example in [example_test.go](../../../example_test.go), also shown in the README. It creates a client with a bearer request editor and calls `ListAccountsWithResponse(ctx, &surge.ListAccountsParams{})`.

- `NewClientWithResponses` returns parsed response methods. `WithHTTPClient` accepts a custom HTTP client with a timeout.
- Go errors cover transport and decoding failures. Check `StatusCode()` and `JSON200` or `JSON201` before accessing the response. Error bodies usually use `JSONDefault`.
- Required fields can be values; optional or nullable fields can be pointers. Check pointers before dereferencing.
- To send directly to a recipient, fill `MessageParams1` with `To` and `Body` or `Attachments`, then call `MessageParams.FromMessageParams1`. The conversation form is `MessageParams0` with `FromMessageParams0`. Pass the union to `CreateMessageWithResponse(ctx, accountID, body)`.
- Message creation queues delivery. Check later status or delivery webhooks to establish delivery.
- Cursor pagination is manual. Pass `Pagination.NextCursor` to the next request's `After` field until nil. Do not decode or alter opaque cursors.
