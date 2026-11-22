module github.com/ritoon/invoice-bench/go/libs/ritoon-zerechnung

go 1.26

// The version of ZeRechnung measured by default. ZERECHNUNG_DIR measures a
// local clone instead (see ../zerechnung.sh).
require (
	github.com/ritoon/invoice-bench/go/harness v0.0.0
	github.com/ritoon/zerechnung v0.1.1-0.20261122170638-724b3887e942
)

replace github.com/ritoon/invoice-bench/go/harness => ../../harness
