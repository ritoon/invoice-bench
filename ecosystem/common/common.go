// Package common holds the shared corpus used by every library benchmark, so
// that all libraries are measured on exactly the same bytes.
package common

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Root of the ZUGFeRD/corpus checkout (Apache-2.0).
func Root() string {
	if v := os.Getenv("FX_CORPUS"); v != "" {
		return v
	}
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "..", "corpus")
}

type Sample struct {
	Name string
	Path string
	Data []byte
}

func Load(tb testing.TB, rel ...string) Sample {
	tb.Helper()
	p := filepath.Join(append([]string{Root()}, rel...)...)
	b, err := os.ReadFile(p)
	if err != nil {
		tb.Fatalf("load %s: %v", p, err)
	}
	return Sample{Name: filepath.Base(p), Path: p, Data: b}
}

// PDF inputs (Factur-X / ZUGFeRD 2.x hybrids).
func PDFSmall(tb testing.TB) Sample { return Load(tb, "XML-Rechnung", "FX", "EN16931_Einfach.pdf") }  // 149 KB, EN 16931
func PDFLarge(tb testing.TB) Sample { return Load(tb, "XML-Rechnung", "FX", "EN16931_Elektron.pdf") } // 1.1 MB, EN 16931
func PDFMinimum(tb testing.TB) Sample {
	return Load(tb, "ZUGFeRDv2", "correct", "FNFE-factur-x-examples", "Facture_FR_MINIMUM.pdf")
}
func PDFExtended(tb testing.TB) Sample {
	return Load(tb, "ZUGFeRDv2", "correct", "intarsys", "EXTENDED", "zugferd_2p0_EXTENDED_Warenrechnung.pdf")
}

// XML inputs (CII D16B).
func XMLSmall(tb testing.TB) Sample {
	return Load(tb, "XML-Rechnung", "CII", "EN16931_Einfach.cii.xml")
} // 13 KB
func XMLLines(tb testing.TB) Sample {
	return Load(tb, "XML-Rechnung", "CII", "EN16931_Rabatte.cii.xml")
} // 19 KB, allowances/charges
func XMLLarge(tb testing.TB) Sample {
	return Load(tb, "XML-Rechnung", "CII", "XRECHNUNG_Elektron.cii.xml")
} // 1.3 MB, base64 attachment

// HostPDF is a plain PDF (attachments removed) used as input for embedding.
func HostPDF(tb testing.TB) Sample {
	_, f, _, _ := runtime.Caller(0)
	p := filepath.Join(filepath.Dir(f), "..", "host.pdf")
	b, err := os.ReadFile(p)
	if err != nil {
		tb.Fatalf("host.pdf missing (run the pdfcpu module TestMakeHost first): %v", err)
	}
	return Sample{Name: "host.pdf", Path: p, Data: b}
}

// LooksLikeCII is the sanity check applied before timing an extraction.
func LooksLikeCII(b []byte) bool {
	return bytes.Contains(b, []byte("CrossIndustryInvoice")) || bytes.Contains(b, []byte("CrossIndustryDocument"))
}

// Corpus lists every valid CII XML used for the false-positive survey.
func CIIFiles(tb testing.TB) []string {
	m, _ := filepath.Glob(filepath.Join(Root(), "XML-Rechnung", "CII", "*.cii.xml"))
	var out []string
	for _, f := range m {
		if filepath.Base(f) != "not_validating_full_invoice_based_onTest_EeISI_300_CENfullmodel.cii.xml" {
			out = append(out, f)
		}
	}
	return out
}

// PDFFiles lists every valid hybrid PDF of ZUGFeRD v2 / Factur-X in the corpus.
func PDFFiles(tb testing.TB) []string {
	var out []string
	filepath.WalkDir(filepath.Join(Root(), "ZUGFeRDv2", "correct"), func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && filepath.Ext(p) == ".pdf" {
			out = append(out, p)
		}
		return nil
	})
	m, _ := filepath.Glob(filepath.Join(Root(), "XML-Rechnung", "FX", "*.pdf"))
	return append(out, m...)
}

// HostPDFClassic is host.pdf rewritten with a classic xref table and no object streams.
func HostPDFClassic(tb testing.TB) Sample {
	_, f, _, _ := runtime.Caller(0)
	b, err := os.ReadFile(filepath.Join(filepath.Dir(f), "..", "host_classic.pdf"))
	if err != nil {
		tb.Fatalf("host_classic.pdf missing: %v", err)
	}
	return Sample{Name: "host_classic.pdf", Data: b}
}
