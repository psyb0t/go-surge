// Package surge is a Go client for the Surge SMS API,
// generated with oapi-codegen from the OpenAPI document Surge publishes.
package surge

//go:generate bash scripts/fetch-spec.sh
//go:generate bash scripts/patch-spec.sh
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.fixed.yaml
