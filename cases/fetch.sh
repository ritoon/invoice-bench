#!/usr/bin/env bash
# Fetches the official test files that cases.tsv lists under external/, at a
# pinned commit, without versioning them. See README.md.
#
#   ./fetch.sh
set -euo pipefail
cd "$(dirname "$0")"

# KoSIT XRechnung test suite, release v2026-08-31 for XRechnung 3.0.2
# (Apache-2.0): the business and technical cases of src/test, all valid.
KOSIT_URL=https://github.com/itplr-kosit/xrechnung-testsuite.git
KOSIT_COMMIT=44c3905c69ca236aa3b94c04be1b7e9bd066ba03
dir=external/xrechnung-testsuite

if [ "$(git -C "$dir" rev-parse HEAD 2>/dev/null)" != "$KOSIT_COMMIT" ] || [ -n "$(git -C "$dir" status --porcelain)" ]; then
	echo "== xrechnung-testsuite @ ${KOSIT_COMMIT:0:7}"
	rm -rf "$dir"
	git init -q "$dir"
	git -C "$dir" sparse-checkout set --cone src/test
	git -C "$dir" fetch -q --depth 1 --filter=blob:none "$KOSIT_URL" "$KOSIT_COMMIT"
	git -C "$dir" -c advice.detachedHead=false checkout -q FETCH_HEAD
	[ "$(git -C "$dir" rev-parse HEAD)" = "$KOSIT_COMMIT" ] || { echo "$KOSIT_URL: HEAD is not $KOSIT_COMMIT" >&2; exit 1; }
fi
