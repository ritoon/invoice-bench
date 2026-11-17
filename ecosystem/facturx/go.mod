module github.com/ritoon/zerechnung/bench/ecosystem/facturx

go 1.26

require (
	github.com/ritoon/zerechnung v0.0.0
	github.com/ritoon/zerechnung/bench/ecosystem/common v0.0.0
)

replace (
	github.com/ritoon/zerechnung => ../../..
	github.com/ritoon/zerechnung/bench/ecosystem/common => ../common
)
