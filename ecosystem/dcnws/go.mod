module github.com/ritoon/factur-x/bench/ecosystem/dcnws

go 1.27

require (
	github.com/ritoon/factur-x/bench/ecosystem/common v0.0.0
	github.com/dcnws/invoicekit.go v0.2.13
)

replace github.com/ritoon/factur-x/bench/ecosystem/common => ../common
