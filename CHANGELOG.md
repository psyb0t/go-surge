# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.0, 2026-10-17

- Added `surge`, a REST client generated with oapi-codegen from Surge's public OpenAPI document (spec version 1.0, OpenAPI 3.1.1). It covers all 41 operations across 25 paths, including accounts, audiences, blasts, campaigns, contacts, messages, phone numbers, recordings, users, tokens and verifications.
- Number fields are `float64` instead of oapi-codegen's default `float32`. All component schemas are retained, including webhook payload models.
- `make generate` reads the current Stainless spec URL from `https://raw.githubusercontent.com/surgeapi/python-sdk/main/.stats.yml`, downloads it, downconverts nullable type lists, constants and schema examples for oapi-codegen, and regenerates the client. The initial source is `https://storage.googleapis.com/stainless-sdk-openapi-specs/surge/surge-ffd563fe3318efc4a305510ec4ae66fea4a32d96ed34abd985b44692edb3535e.yml`. The spec is not committed.
