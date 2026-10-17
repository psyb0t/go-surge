#!/bin/bash
# Downconverts Surge's JSON OpenAPI 3.1 document for oapi-codegen:
# 1. Nullable type lists become a single type plus nullable: true.
# 2. const becomes a single-value enum so event type constants survive.
# 3. Schema examples arrays become example; response examples maps stay intact.
# The version marker becomes 3.0.3 after these conversions.
# The original download stays untouched in api/openapi.yaml.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR

on_error() {
	local exit_code="$1"
	local line="$2"

	printf '{"level":"ERROR","file":"%s","line":%d,"msg":"command failed","exit_code":%d}\n' \
		"${BASH_SOURCE[1]}" "$line" "$exit_code" >&2
}

log() {
	local level="$1"
	shift

	jq -nc \
		--arg time "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--arg level "$level" \
		--arg file "${BASH_SOURCE[1]##*/}" \
		--arg line "${BASH_LINENO[0]}" \
		--arg func "${FUNCNAME[1]:-main}" \
		--arg msg "$*" \
		'{time:$time,level:$level,file:$file,line:($line|tonumber),func:$func,msg:$msg}' >&2
}

main() {
	local tmp
	tmp="$(mktemp "${SPEC_DIR}/.openapi.fixed.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp while the local exists
	trap "rm -f -- '${tmp}'" EXIT

	jq '
	  walk(if type == "object" then
	    (if (.type | type) == "array" then
	      if (.type | length) == 2 and (.type | index("null")) != null then
	        .type = (.type | map(select(. != "null"))[0]) | .nullable = true
	      else error("unsupported type union") end
	    else . end)
	    | (if has("const") then .enum = [.const] | del(.const) else . end)
	    | (if (.examples | type) == "array" then
	        .example = (.example // .examples[0].value // .examples[0])
	        | del(.examples)
	      else . end)
	  else . end)
	  | .openapi = "3.0.3"
	' "${SPEC_DIR}/openapi.yaml" >"$tmp"
	chmod 0644 "$tmp"
	mv -- "$tmp" "${SPEC_DIR}/openapi.fixed.yaml"
	log INFO "downconverted spec to OpenAPI 3.0.3"
}

main "$@"
