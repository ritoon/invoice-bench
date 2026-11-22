module github.com/ritoon/invoice-bench/go/libs/tahayassine-facturx-go

go 1.26

require github.com/ritoon/invoice-bench/go/harness v0.0.0

require (
	github.com/tahayassine/facturx-go v0.1.0
	golang.org/x/text v0.29.0 // indirect
)

replace github.com/ritoon/invoice-bench/go/harness => ../../harness
