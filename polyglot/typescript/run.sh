#!/usr/bin/env bash
# Banc TypeScript / JavaScript des librairies Factur-X / ZUGFeRD / EN 16931.
# Voir ../PROTOCOL.md et COVERAGE.md.
#
#   ./run.sh setup             dépendances figées (npm ci --ignore-scripts)
#   ./run.sh survey [lib…]     relevés d'exactitude  -> $OUT/survey_typescript.jsonl
#   ./run.sh bench  [lib…]     benchmarks            -> $OUT/bench_typescript.jsonl
#
# Variables : OUT (../results/<date>), COUNT (3), BENCHTIME (1, secondes).
# Librairies (toutes par défaut) : node adapters/index.mjs
set -euo pipefail
cd "$(dirname "$0")"
here=$(pwd)

OUT=${OUT:-../results/$(date +%F)}
mkdir -p "$OUT"
OUT=$(cd "$OUT" && pwd)
export OUT COUNT=${COUNT:-3} BENCHTIME=${BENCHTIME:-1}
PARTS=$OUT/typescript-parts
# One-line progress display of a bench_typescript.jsonl line.
SHOW='let s="";process.stdin.on("data",d=>s+=d).on("end",()=>{const r=JSON.parse(s);const ns=r.ns_per_op.length?Math.min(...r.ns_per_op):null;console.log(`   ${r.bench.padEnd(20)} ${r.status.padEnd(11)} ${ns==null?"":(ns/1e3).toFixed(1)+" µs/op"}`)})'

setup() {
	# No install script runs (also pinned in .npmrc): none of the benchmarked
	# packages needs one, and @fin.cx/einvoice's postinstall downloads
	# unpinned Schematron files from GitHub.
	npm ci --ignore-scripts --no-audit --no-fund
	[ -d ../../ecosystem/corpus/XML-Rechnung ] || echo "corpus absent : lancer ../../ecosystem/run.sh setup" >&2
	[ -f ../../ecosystem/host.pdf ] || echo "host.pdf absent : lancer ../../ecosystem/run.sh setup" >&2
	echo "== node $(node --version), $(node -p "Object.keys(require('./package.json').dependencies).length") paquets figés"
}

need() {
	[ -d node_modules/@attestwire/en16931 ] || setup
}

libs() { if [ $# -gt 0 ]; then echo "$@"; else node adapters/index.mjs; fi; }

# merge <kind>: survey_typescript.jsonl / bench_typescript.jsonl are the
# per-library parts concatenated in campaign order, so that rerunning one
# library replaces its lines and keeps the others.
merge() {
	local kind=$1 out=$OUT/${1}_typescript.jsonl
	: >"$out.tmp"
	for id in $(node adapters/index.mjs); do
		[ -f "$PARTS/${kind}_$id.jsonl" ] && cat "$PARTS/${kind}_$id.jsonl" >>"$out.tmp"
	done
	mv "$out.tmp" "$out"
	echo "== $out ($(wc -l <"$out" | tr -d ' ') lignes)"
}

survey() {
	mkdir -p "$PARTS"
	for id in $(libs "$@"); do
		echo "== survey $id"
		node survey.mjs "$id" >"$PARTS/survey_$id.jsonl.tmp"
		mv "$PARTS/survey_$id.jsonl.tmp" "$PARTS/survey_$id.jsonl"
	done
	merge survey
}

bench() {
	mkdir -p "$PARTS"
	local benches
	benches=$(node bench.mjs --list)
	for id in $(libs "$@"); do
		echo "== bench $id"
		: >"$PARTS/bench_$id.jsonl.tmp"
		for b in $benches; do
			# One process per benchmark: first_call_ns is a cold first call. A
			# process that dies without its line still leaves a "fail" line.
			line=$(node bench.mjs "$id" "$b") || true
			[ -n "$line" ] || line=$(node bench.mjs --failed "$id" "$b" "processus node terminé sans résultat")
			echo "$line" >>"$PARTS/bench_$id.jsonl.tmp"
			echo "$line" | node -e "$SHOW"
		done
		mv "$PARTS/bench_$id.jsonl.tmp" "$PARTS/bench_$id.jsonl"
	done
	merge bench
}

cmd=${1:-}
shift || true
case "$cmd" in
setup) setup ;;
survey)
	need
	survey "$@"
	;;
bench)
	need
	bench "$@"
	;;
*)
	sed -n '2,10p' "$0"
	exit 2
	;;
esac
cd "$here"
