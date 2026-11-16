module github.com/ritoon/factur-x/bench/ecosystem/facturx

go 1.26

require (
	github.com/ritoon/factur-x v0.0.0
	github.com/ritoon/factur-x/bench/ecosystem/common v0.0.0
)

replace (
	github.com/ritoon/factur-x => ../../..
	github.com/ritoon/factur-x/bench/ecosystem/common => ../common
)
