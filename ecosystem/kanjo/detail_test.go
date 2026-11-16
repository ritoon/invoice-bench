package kanjobench

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cyprienbrisset/kanjo/pkg/rules"
	"github.com/ritoon/factur-x/bench/ecosystem/common"
)

func TestDetail(t *testing.T) {
	for _, f := range common.CIIFiles(t)[:3] {
		data, _ := os.ReadFile(f)
		r, _ := parse(data)
		rep := rules.Validate(r.Doc)
		t.Logf("%s sets=%v run=%d", filepath.Base(f), rep.RuleSets, rep.RulesRun)
		for _, x := range rep.Findings {
			if x.Severity >= rules.SeverityError {
				t.Logf("  %s | %.140s", x.RuleID, x.Message)
			}
		}
	}
}
