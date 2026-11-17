package pdf0bench

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mgilbir/formalis"
	"github.com/ritoon/zerechnung/bench/ecosystem/common"
)

func TestDetail(t *testing.T) {
	for _, f := range common.CIIFiles(t) {
		data, _ := os.ReadFile(f)
		r, _ := formalis.Validate(context.Background(), data, formalis.ProfileEN16931)
		for _, v := range r.Fatal() {
			t.Logf("%s | %s | %.100s", filepath.Base(f), v.Rule, v.Message)
		}
	}
}
