#!/usr/bin/env bash
# Runs the reference validator of each suite on the documents of cases.tsv
# and records its verdicts and the rules it raises as errors. See README.md.
#
#   ./confirm.sh
#
# The reference of the en16931 and xrechnung suites is the KoSIT validator
# with the official XRechnung 3.0.2 configuration of 2026-08-31, as the Java
# bench runs it (java/, library kosit-validator): ../java/run.sh
# setup must have been run.
set -euo pipefail
cd "$(dirname "$0")"

./fetch.sh
(cd tool && go run . index)

out=$(mktemp -d)
trap 'rm -rf "$out"' EXIT
OUT=$out ../java/run.sh cases kosit-validator
(cd tool && go run . confirm "$out")
