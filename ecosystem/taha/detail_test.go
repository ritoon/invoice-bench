package tahabench

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ritoon/factur-x/bench/ecosystem/common"
	"github.com/tahayassine/facturx-go/pkg/cii"
	"github.com/tahayassine/facturx-go/pkg/en16931"
)

func TestDetail(t *testing.T) {
	for _, f := range common.CIIFiles(t) {
		data, _ := os.ReadFile(f)
		inv, _ := cii.Decode(data)
		for _, e := range en16931.Validate(inv).Errors {
			t.Logf("%s | %s | %.100s", filepath.Base(f), e.Rule, e.Message)
		}
	}
}
