package gobldebench

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ritoon/factur-x/bench/ecosystem/common"
)

func TestDetail(t *testing.T) {
	for _, f := range common.CIIFiles(t)[:4] {
		data, _ := os.ReadFile(f)
		p, _ := parse(data)
		t.Logf("%s: %.300v", filepath.Base(f), p.Envelope.Validate())
	}
}
