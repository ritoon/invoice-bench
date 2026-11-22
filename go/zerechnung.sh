# Sourced by run.sh and ../run.sh.
#
# The libs/ritoon-zerechnung module measures the version of ZeRechnung that
# its go.mod requires. ZERECHNUNG_DIR, if set, is a local clone of
# github.com/ritoon/zerechnung to measure instead (a branch before it is
# merged, say): zerechnung_goflags then writes local.mod in the module, a copy
# of its go.mod with a replace directive to that clone, and prints GOFLAGS
# that make the go command read it.
#
#   GOFLAGS=$(zerechnung_goflags path/to/libs/ritoon-zerechnung) go test …

if [ -n "${ZERECHNUNG_DIR:-}" ] && ! grep -qs '^module github.com/ritoon/zerechnung$' "$ZERECHNUNG_DIR/go.mod"; then
	echo "ZERECHNUNG_DIR=$ZERECHNUNG_DIR is not a clone of github.com/ritoon/zerechnung" >&2
	exit 2
fi

zerechnung_goflags() {
	local dir zdir flags=${GOFLAGS:-}
	dir=$(cd "$1" && pwd)
	if [ -n "${ZERECHNUNG_DIR:-}" ]; then
		zdir=$(cd "$ZERECHNUNG_DIR" && pwd)
		cp "$dir/go.mod" "$dir/local.mod"
		cp "$dir/go.sum" "$dir/local.sum"
		go mod edit -replace "github.com/ritoon/zerechnung=$zdir" "$dir/local.mod"
		flags="${flags:+$flags }-modfile=$dir/local.mod"
	fi
	echo "$flags"
}

# zerechnung_version prints the version of ZeRechnung the module in $1 builds
# with, and the local clone it is replaced by, if any.
zerechnung_version() {
	local flags
	flags=$(zerechnung_goflags "$1")
	(cd "$1" && GOFLAGS=$flags go list -m github.com/ritoon/zerechnung)
}
