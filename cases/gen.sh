#!/usr/bin/env bash
# Regenerates the mutation cases of the suites from the mutation tests of
# ZeRechnung, at the commit below (or from the clone ZERECHNUNG_DIR), then
# rebuilds cases.tsv. See README.md.
#
#   ./gen.sh
#
# Each mutation test writes one document per case and expected.tsv, the rules
# the case was written to break (internal/variants of ZeRechnung). These rules
# become the "intent" column of cases.tsv; the expected verdict and rules come
# from the reference validator afterwards (./confirm.sh), not from ZeRechnung.
set -euo pipefail
cd "$(dirname "$0")"

ZERECHNUNG_URL=https://github.com/ritoon/zerechnung.git
ZERECHNUNG_COMMIT=9b9f489814bb03ca48b4c4a869b185441e5bfbc2

if [ -n "${ZERECHNUNG_DIR:-}" ]; then
	src=$(cd "$ZERECHNUNG_DIR" && pwd)
	grep -qs '^module github.com/ritoon/zerechnung$' "$src/go.mod" ||
		{ echo "ZERECHNUNG_DIR=$ZERECHNUNG_DIR is not a clone of github.com/ritoon/zerechnung" >&2; exit 2; }
else
	src=.deps/zerechnung
	if [ "$(git -C "$src" rev-parse HEAD 2>/dev/null)" != "$ZERECHNUNG_COMMIT" ]; then
		echo "== zerechnung @ ${ZERECHNUNG_COMMIT:0:7}"
		rm -rf "$src"
		git init -q "$src"
		git -C "$src" fetch -q --depth 1 "$ZERECHNUNG_URL" "$ZERECHNUNG_COMMIT"
		git -C "$src" -c advice.detachedHead=false checkout -q FETCH_HEAD
	fi
fi
echo "== zerechnung $(git -C "$src" rev-parse --short HEAD)"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

# variants pkg test dir: the documents of one mutation test.
variants() {
	(cd "$src" && GOFLAGS=-mod=mod go test -count 1 "./internal/rules/$1" -run "^$2\$" -variants "$tmp/$3" >/dev/null)
}

# suite name dir...: replaces the documents of a suite by those of the
# mutation tests written to dir..., and merges their expected.tsv.
suite() {
	local name=$1
	shift
	rm -rf "$name"
	mkdir -p "$name"
	: >"$name/intent.tsv"
	for d in "$@"; do
		cp "$tmp/$d"/*.xml "$name/"
		cat "$tmp/$d/expected.tsv" >>"$name/intent.tsv"
	done
	sort -o "$name/intent.tsv" "$name/intent.tsv"
	echo "== $name: $(ls "$name"/*.xml | wc -l | tr -d ' ') documents"
}

variants en16931 TestVariants en16931
variants ubl TestVariants ubl
variants de TestVariants de
suite en16931 en16931 ubl
suite xrechnung de

(cd tool && go run . index)
