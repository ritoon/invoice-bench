package andeebench

import (
	"os"
	"testing"

	xinvoice "github.com/andeedotnet/go-xinvoice"
	"github.com/ritoon/zerechnung/bench/ecosystem/common"
)

func TestDetail(t *testing.T) {
	c := map[string]int{}
	for _, f := range common.CIIFiles(t) {
		data, _ := os.ReadFile(f)
		inv, _ := parse(data)
		for _, x := range xinvoice.Validate(inv).Findings {
			if x.Severity == xinvoice.SeverityError {
				c[x.Rule]++
			}
		}
	}
	t.Logf("RULES %v", c)
}
