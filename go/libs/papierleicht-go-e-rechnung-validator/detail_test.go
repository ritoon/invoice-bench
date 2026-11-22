package papierleichtbench

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ritoon/invoice-bench/go/harness"
)

func TestDetail(t *testing.T) {
	vv := val(t)
	for _, f := range harness.CIIFiles(t) {
		data, _ := os.ReadFile(f)
		r := vv.Validate(data)
		for _, e := range r.Errors() {
			t.Logf("%s | %s | %v | %.120s", filepath.Base(f), e.RuleID, r.Format, e.Message)
		}
	}
}
