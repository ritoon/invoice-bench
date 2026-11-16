module github.com/ritoon/factur-x/bench/ecosystem/kanjo

go 1.26

require github.com/ritoon/factur-x/bench/ecosystem/common v0.0.0

require (
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/cyprienbrisset/kanjo v0.3.2
	github.com/hhrutter/tiff v1.0.6 // indirect
	github.com/mattn/go-runewidth v0.0.27 // indirect
	github.com/pdfcpu/pdfcpu v0.15.0
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/image v0.44.0 // indirect
	golang.org/x/text v0.40.0 // indirect
)

replace github.com/ritoon/factur-x/bench/ecosystem/common => ../common
