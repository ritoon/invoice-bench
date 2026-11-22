package cases

import (
	"math"
	"testing"
)

func TestNormalize(t *testing.T) {
	got := Normalize([]string{
		"[BR-CO-10] Sum of Invoice line net amount",
		"br-de-15",
		"FX-SCH-A-000571",
		"BR-AF-05",
		"cvc-complex-type.2.4.a: Invalid content was found",
		"XSD: Element 'x': This element is not expected.",
		"a message without identifier",
	})
	want := []string{"BR-CO-10", "BR-DE-15", "BR-CL-25", "BR-IP-05", "XSD"}
	if len(got) != len(want) {
		t.Fatalf("Normalize = %v, want %v", got, want)
	}
	for _, w := range want {
		if !got[w] {
			t.Errorf("Normalize: %s missing from %v", w, got)
		}
	}
}

func testIndex() *Index {
	return &Index{Suites: []Suite{{Name: "s", Syntaxes: []string{"cii"}, Countries: []string{"DE"}, Unchecked: []string{"BR-CO-25"}}}}
}

func TestJudge(t *testing.T) {
	idx := testIndex()
	cases := []Case{
		{Suite: "s", File: "v1", Verdict: Valid},
		{Suite: "s", File: "v2", Verdict: Valid},
		{Suite: "s", File: "i1", Verdict: Invalid, Rules: []string{"BR-01"}},
		{Suite: "s", File: "i2", Verdict: Invalid, Rules: []string{"BR-01"}},
		{Suite: "s", File: "i3", Verdict: Invalid, Rules: []string{"BR-02"}},
		{Suite: "s", File: "i4", Verdict: Invalid, Rules: []string{"XSD"}, Intent: []string{"BR-03"}},
		{Suite: "s", File: "x", Verdict: Excluded},
	}
	outs := map[string]Outcome{
		"v1": {},
		"v2": {Errors: 1, Rules: []string{"BR-05"}},
		"i1": {Errors: 1, Rules: []string{"BR-01"}},
		"i2": {Errors: 1, Rules: []string{"BR-09"}}, // rejected under another rule: 0.5
		// i3 accepted: 0
		"i3": {},
		"i4": {Errors: 1, Rules: []string{"[BR-03] intended rule"}}, // XSD case, intent raised: 1
		"x":  {Errors: 1, Error: "boom"},
	}
	s := idx.Judge(cases, outs)
	if s.Valid != 2 || s.Accepted != 1 || s.Invalid != 4 || s.Rejected != 3 || s.Exact != 2 || s.Missing != 0 {
		t.Fatalf("counts: %+v", s)
	}
	want := map[string]float64{"BR-01": 0.75, "BR-02": 0, "XSD": 1}
	for r, c := range want {
		if s.Credit[r] != c {
			t.Errorf("credit %s = %v, want %v", r, s.Credit[r], c)
		}
	}
	d := (0.75 + 0 + 1) / 3
	if math.Abs(s.Note()-100*math.Sqrt(0.5*d)) > 1e-9 {
		t.Errorf("Note = %v, want %v", s.Note(), 100*math.Sqrt(0.5*d))
	}
}

func TestJudgeUnchecked(t *testing.T) {
	idx := testIndex()
	cases := []Case{
		{Suite: "s", File: "a", Verdict: Valid},
		{Suite: "s", File: "b", Verdict: Valid},
		{Suite: "s", File: "c", Verdict: Valid},
		{Suite: "s", File: "d", Verdict: Valid},
	}
	outs := map[string]Outcome{
		"a": {Errors: 2, Rules: []string{"BR-CO-25", "[BR-CO-25] payment terms"}}, // tolerated
		"b": {Errors: 2, Rules: []string{"BR-CO-25", "BR-DE-15"}},                 // another rule
		"c": {Errors: 1, Rules: []string{"payment terms missing"}},                // no identifier
		"d": {Errors: 1, Rules: []string{"BR-CO-25"}, Error: "exception"},         // exception
	}
	s := idx.Judge(cases, outs)
	if s.Accepted != 1 || s.Tolerated != 1 {
		t.Fatalf("accepted %d, tolerated %d, want 1 and 1", s.Accepted, s.Tolerated)
	}
	// A missing outcome counts as an exception.
	s = idx.Judge([]Case{{Suite: "s", File: "zz", Verdict: Valid}}, outs)
	if s.Accepted != 0 || s.Missing != 1 {
		t.Fatalf("missing outcome: %+v", s)
	}
}
