package fxbench

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ritoon/zerechnung/bench/ecosystem/common"
	"github.com/solarpush/fx/pkg/cii"
	"github.com/solarpush/fx/pkg/invoice"
)

func TestDetail(t *testing.T) {
	for _, f := range common.CIIFiles(t) {
		data, _ := os.ReadFile(f)
		inv, _ := cii.Parse(data)
		for _, e := range invoice.ValidateReport(inv).Errors() {
			t.Logf("%s | %s | %s | %.80s", filepath.Base(f), e.Code, e.Field, e.Message)
		}
	}
}
