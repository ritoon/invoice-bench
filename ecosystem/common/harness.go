package common

import (
	"os"
	"path/filepath"
	"testing"
)

// Extract benchmarks fn on s after a correctness check.
func Extract(b *testing.B, s Sample, fn func([]byte) ([]byte, error)) {
	if x, err := fn(s.Data); err != nil || !LooksLikeCII(x) {
		b.Skipf("UNSUPPORTED extraction on %s: %v", s.Name, err)
	}
	b.SetBytes(int64(len(s.Data)))
	b.ReportAllocs()
	for b.Loop() {
		if _, err := fn(s.Data); err != nil {
			b.Fatal(err)
		}
	}
}

// Bytes benchmarks any []byte -> (T, error) operation (parse, validate...).
func Bytes[T any](b *testing.B, s Sample, fn func([]byte) (T, error)) {
	if _, err := fn(s.Data); err != nil {
		b.Skipf("UNSUPPORTED on %s: %v", s.Name, err)
	}
	b.SetBytes(int64(len(s.Data)))
	b.ReportAllocs()
	for b.Loop() {
		fn(s.Data)
	}
}

// Write benchmarks serialisation of a model parsed once outside the loop.
func Write[M any](b *testing.B, s Sample, parse func([]byte) (M, error), write func(M) ([]byte, error)) {
	m, err := parse(s.Data)
	if err != nil {
		b.Skipf("UNSUPPORTED parse on %s: %v", s.Name, err)
	}
	if out, err := write(m); err != nil || len(out) == 0 {
		b.Skipf("UNSUPPORTED write on %s: %v", s.Name, err)
	}
	b.ReportAllocs()
	for b.Loop() {
		write(m)
	}
}

// Embed benchmarks host PDF + XML -> Factur-X PDF.
func Embed(b *testing.B, host, xml Sample, fn func(pdf, xml []byte) ([]byte, error), extract func([]byte) ([]byte, error)) {
	out, err := fn(host.Data, xml.Data)
	if err != nil {
		b.Skipf("UNSUPPORTED embed: %v", err)
	}
	if extract != nil {
		if back, err := extract(out); err != nil || !LooksLikeCII(back) {
			b.Fatalf("embedded PDF does not round-trip: %v", err)
		}
	}
	b.SetBytes(int64(len(host.Data)))
	b.ReportAllocs()
	for b.Loop() {
		if _, err := fn(host.Data, xml.Data); err != nil {
			b.Fatal(err)
		}
	}
}

// SurveyValidate counts valid official samples that a validator flags (false positives).
// count returns the number of errors (0 = accepted).
func SurveyValidate(t *testing.T, count func([]byte) (int, error)) {
	flagged, total := 0, 0
	for _, f := range CIIFiles(t) {
		data, _ := os.ReadFile(f)
		n, err := count(data)
		total++
		if err != nil || n > 0 {
			flagged++
			t.Logf("FLAG %s: %d errors, err=%v", filepath.Base(f), n, err)
		}
	}
	t.Logf("SURVEY validate: %d/%d valid samples flagged", flagged, total)
}

// SurveyExtract counts corpus PDFs from which an invoice XML is extracted.
func SurveyExtract(t *testing.T, fn func([]byte) ([]byte, error)) {
	ok, all := 0, 0
	for _, f := range PDFFiles(t) {
		data, _ := os.ReadFile(f)
		all++
		if x, err := fn(data); err == nil && LooksLikeCII(x) {
			ok++
		} else {
			t.Logf("MISS %s: %v", filepath.Base(f), err)
		}
	}
	t.Logf("SURVEY extract: %d/%d PDFs", ok, all)
}

// SurveyParse counts valid CII samples a parser accepts.
func SurveyParse[T any](t *testing.T, fn func([]byte) (T, error)) {
	ok, all := 0, 0
	for _, f := range CIIFiles(t) {
		data, _ := os.ReadFile(f)
		all++
		if _, err := fn(data); err == nil {
			ok++
		} else {
			t.Logf("PARSEFAIL %s: %v", filepath.Base(f), err)
		}
	}
	t.Logf("SURVEY parse: %d/%d CII", ok, all)
}
