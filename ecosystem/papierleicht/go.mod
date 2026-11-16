module github.com/ritoon/factur-x/bench/ecosystem/papierleicht

go 1.26

require (
	github.com/ritoon/factur-x/bench/ecosystem/common v0.0.0
	github.com/papierleicht/go-e-rechnung-validator v0.1.0
)

require (
	github.com/andybalholm/brotli v1.2.7 // indirect
	github.com/knroy/go-xml v1.8.1 // indirect
	github.com/speedata/pdfdisassembler v0.0.9 // indirect
	golang.org/x/text v0.41.0 // indirect
)

replace github.com/ritoon/factur-x/bench/ecosystem/common => ../common
