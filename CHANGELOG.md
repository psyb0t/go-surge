# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.1, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.1.0, 2026-10-17

- Added `surge`, a REST client generated with oapi-codegen from Surge's public OpenAPI document (spec version 1.0, OpenAPI 3.1.1). It covers all 41 operations across 25 paths, including accounts, audiences, blasts, campaigns, contacts, messages, phone numbers, recordings, users, tokens and verifications.
- Number fields are `float64` instead of oapi-codegen's default `float32`. All component schemas are retained, including webhook payload models.
- `make generate` reads the current Stainless spec URL from `https://raw.githubusercontent.com/surgeapi/python-sdk/main/.stats.yml`, downloads it, downconverts nullable type lists, constants and schema examples for oapi-codegen, and regenerates the client. The initial source is `https://storage.googleapis.com/stainless-sdk-openapi-specs/surge/surge-ffd563fe3318efc4a305510ec4ae66fea4a32d96ed34abd985b44692edb3535e.yml`. The spec is not committed.
