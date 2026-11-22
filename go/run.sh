#!/usr/bin/env bash
# Benchmarks of the Go Factur-X ecosystem. See README.md.
#
#   ./run.sh setup              fetch the corpus and dependencies, build the host PDFs
#                               (corpus and host PDFs go to ../testdata, shared by
#                               every language)
#   ./run.sh survey [module…]   correctness surveys on the whole corpus
#   ./run.sh cases  [module…]   conformance cases of ../cases (validating modules)
#   ./run.sh bench  [module…]   benchmarks (COUNT and BENCHTIME can be overridden)
#   ./run.sh report             tables in results/<date>/tables.md
#   ./run.sh all    [module…]   setup + survey + cases + bench + report
#
# A module is a folder of libs/, named after the Go module it measures
# (ritoon-zerechnung, pdfcpu-pdfcpu…). Results go to results/<date>/ (OUT can
# be overridden): bench_<module>.txt, survey_<module>.txt and
# cases_<module>.jsonl.
# ZERECHNUNG_DIR measures a local clone of ZeRechnung instead of the version
# that libs/ritoon-zerechnung/go.mod requires (see zerechnung.sh).
set -euo pipefail
cd "$(dirname "$0")"
. ./zerechnung.sh

CORPUS_URL=https://github.com/ZUGFeRD/corpus.git
CORPUS_COMMIT=84f8d8f68f6c3698aa3019d701a8b82db9ba0aac
FX_URL=https://github.com/solarpush/fX.git
FX_COMMIT=1314f3d94a8092ab94a2db8909f07f00e0eab875
TESTDATA=../testdata

MODULES=(
	stdlib
	ritoon-zerechnung
	pdfcpu-pdfcpu
	speedata-einvoice
	andeedotnet-go-xinvoice
	dotwavehq-go-einvoice
	invopop-gobl.de.xinvoice
	angelodlfrtr-go-invoice-generator
	papierleicht-go-e-rechnung-validator
	cyprienbrisset-kanjo
	dcnws-invoicekit.go
	mgilbir-pdf0
	gotenberg-litho
	tahayassine-facturx-go
	solarpush-fx
	marlinkuhn-gopdfattach
)
OUT=${OUT:-results/$(date +%F)}
COUNT=${COUNT:-3}
BENCHTIME=${BENCHTIME:-1s}

# One module per library: -mod=mod lets go resolve them without a workspace,
# GOTOOLCHAIN=auto downloads Go 1.27.x for the modules that require it.
export GOFLAGS=-mod=mod GOTOOLCHAIN=auto

setup() {
	if [ ! -d "$TESTDATA/corpus/.git" ]; then
		echo "== corpus ZUGFeRD (sparse, ${CORPUS_COMMIT:0:7})"
		git clone --quiet --filter=blob:none --no-checkout "$CORPUS_URL" "$TESTDATA/corpus"
		git -C "$TESTDATA/corpus" sparse-checkout set XML-Rechnung ZUGFeRDv1/correct ZUGFeRDv2/correct
		git -C "$TESTDATA/corpus" checkout --quiet "$CORPUS_COMMIT"
	fi
	# The module zip of solarpush/fX is rejected by the go command (a file name
	# contains a carriage return), so the fx module builds from a local clone.
	if [ ! -d .deps/solarpush-fx/.git ]; then
		echo "== solarpush/fX (${FX_COMMIT:0:7})"
		git clone --quiet "$FX_URL" .deps/solarpush-fx
		git -C .deps/solarpush-fx checkout --quiet "$FX_COMMIT"
	fi
	echo "== host PDFs ($TESTDATA/host.pdf, $TESTDATA/host_classic.pdf)"
	(cd libs/pdfcpu-pdfcpu && go test -count 1 -run '^(TestMakeHost|TestMakeHostClassic|TestHostsStrict)$' . >/dev/null)
}

modules() { if [ $# -gt 0 ]; then echo "$@"; else echo "${MODULES[@]}"; fi; }

# check_modules stops on a name that is not a folder of libs/.
check_modules() {
	local m
	for m in "$@"; do
		[ -f "libs/$m/go.mod" ] || { echo "unknown module: $m (the folders of libs/)" >&2; exit 2; }
	done
}

# goflags prints the GOFLAGS of module $1: those of zerechnung_goflags for
# ZeRechnung, the exported ones for the others.
goflags() { if [ "$1" = ritoon-zerechnung ]; then zerechnung_goflags libs/ritoon-zerechnung; else echo "$GOFLAGS"; fi; }

survey() {
	check_modules "$@"
	mkdir -p "$OUT"
	for m in $(modules "$@"); do
		[ "$m" = stdlib ] && continue
		echo "== survey $m"
		(export GOFLAGS=$(goflags "$m") && cd "libs/$m" && go test -count 1 -run '^TestSurvey$' -v . 2>&1 || true) |
			grep -E 'SURVEY|FLAG|MISS|PARSEFAIL' | sed -E 's/^ +[a-z_]+\.go:[0-9]+: //' >"$OUT/survey_$m.txt" || true
		grep SURVEY "$OUT/survey_$m.txt" || echo "   (no result)"
	done
}

# cases: the modules with a validation run TestCases, which writes
# $OUT/cases_<module>.jsonl (see ../cases/README.md).
cases() {
	check_modules "$@"
	mkdir -p "$OUT"
	local abs
	abs=$(cd "$OUT" && pwd)
	for m in $(modules "$@"); do
		grep -qs '^func TestCases' "libs/$m"/*_test.go || continue
		echo "== cases $m"
		(export GOFLAGS=$(goflags "$m") CASES_OUT="$abs/cases_$m.jsonl" && cd "libs/$m" &&
			go test -count 1 -run '^TestCases$' -timeout 2h -v . 2>&1 || true) |
			grep -E 'CASES|FAIL|panic:' | sed -E 's/^ +[a-z_]+\.go:[0-9]+: //' || true
	done
}

bench() {
	check_modules "$@"
	mkdir -p "$OUT"
	for m in $(modules "$@"); do
		echo "== bench $m"
		# Some benchmarks fail on purpose (a library returns a PDF that does not
		# read back): keep going and record what was measured. The output of
		# ritoon-zerechnung starts with the version of ZeRechnung it measures.
		{
			[ "$m" = ritoon-zerechnung ] && zerechnung_version libs/ritoon-zerechnung
			(export GOFLAGS=$(goflags "$m") && cd "libs/$m" && go version && go test -run '^$' -bench . -benchmem -count "$COUNT" -benchtime "$BENCHTIME" . 2>&1) || true
		} >"$OUT/bench_$m.txt"
		grep -E -- '--- FAIL|^FAIL' "$OUT/bench_$m.txt" || true
	done
}

report() {
	local abs
	abs=$(cd "$OUT" && pwd)
	(cd report && go run . "$abs") >"$abs/tables.md" && echo "== $OUT/tables.md"
}

cmd=${1:-}
shift || true
case "$cmd" in
setup) setup ;;
survey) survey "$@" ;;
cases) cases "$@" ;;
bench) bench "$@" ;;
report) report ;;
all) setup && survey "$@" && cases "$@" && bench "$@" && report ;;
*)
	sed -n '2,17p' "$0"
	exit 2
	;;
esac
