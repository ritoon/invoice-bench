package dcnwsbench

import (
	"os"
	"path/filepath"
	"testing"

	einvoice "github.com/dcnws/invoicekit.go"
	"github.com/ritoon/factur-x/bench/ecosystem/common"
)

func TestDetail(t *testing.T) {
	for _, f := range common.CIIFiles(t) {
		data, _ := os.ReadFile(f)
		r, _ := einvoice.Validate(data, einvoice.Options{})
		for _, x := range r.Findings {
			if x.Flag == "fatal" {
				t.Logf("%s | %s | %s", filepath.Base(f), x.ID, x.Source)
			}
		}
	}
}
