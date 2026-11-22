package angelobench

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ritoon/invoice-bench/go/harness"
)

func TestDetail(t *testing.T) {
	for _, f := range harness.CIIFiles(t) {
		data, _ := os.ReadFile(f)
		p, _ := parse(data)
		if err := p.Verify(); err != nil {
			t.Logf("%s: %.200v", filepath.Base(f), err)
		}
	}
}
