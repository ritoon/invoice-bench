package dcnwsbench

import (
	"testing"

	einvoice "github.com/dcnws/invoicekit.go"
	"github.com/dcnws/invoicekit.go/facturx"
	"github.com/dcnws/invoicekit.go/hybrid"
	"github.com/ritoon/invoice-bench/go/harness"
)

func extract(p []byte) ([]byte, error)              { x, _, err := facturx.Extract(p); return x, err }
func parse(x []byte) (*einvoice.ParseResult, error) { return einvoice.Parse(x) }
func validate(x []byte) (int, error) { // Schematron CEN 1.3.17 interpreted in pure Go
	r, err := einvoice.Validate(x, einvoice.Options{})
	if err != nil {
		return 0, err
	}
	return r.Fatal, nil
}
func write(p *einvoice.ParseResult) ([]byte, error) { return einvoice.GenerateCII(p.Invoice) }
func embed(p, x []byte) ([]byte, error)             { return facturx.Embed(p, x, facturx.Options{}) }
func check(p []byte) (int, error) {
	r, err := hybrid.Check(p, einvoice.Options{})
	if err != nil {
		return 0, err
	}
	if r.Rules != nil {
		return r.Rules.Fatal, nil
	}
	return 0, nil
}

func BenchmarkExtract_Small(b *testing.B)       { harness.Extract(b, harness.PDFSmall(b), extract) }
func BenchmarkExtract_Large(b *testing.B)       { harness.Extract(b, harness.PDFLarge(b), extract) }
func BenchmarkExtract_Minimum(b *testing.B)     { harness.Extract(b, harness.PDFMinimum(b), extract) }
func BenchmarkExtract_Extended(b *testing.B)    { harness.Extract(b, harness.PDFExtended(b), extract) }
func BenchmarkParse_Small(b *testing.B)         { harness.Bytes(b, harness.XMLSmall(b), parse) }
func BenchmarkParse_Lines(b *testing.B)         { harness.Bytes(b, harness.XMLLines(b), parse) }
func BenchmarkParse_Large(b *testing.B)         { harness.Bytes(b, harness.XMLLarge(b), parse) }
func BenchmarkParseValidate_Small(b *testing.B) { harness.Bytes(b, harness.XMLSmall(b), validate) }
func BenchmarkParseValidate_Lines(b *testing.B) { harness.Bytes(b, harness.XMLLines(b), validate) }
func BenchmarkParseValidate_Large(b *testing.B) { harness.Bytes(b, harness.XMLLarge(b), validate) }
func BenchmarkWrite_Small(b *testing.B)         { harness.Write(b, harness.XMLSmall(b), parse, write) }
func BenchmarkWrite_Lines(b *testing.B)         { harness.Write(b, harness.XMLLines(b), parse, write) }
func BenchmarkEmbed_Small(b *testing.B) {
	harness.Embed(b, harness.HostPDF(b), harness.XMLSmall(b), embed, extract)
}
func BenchmarkCheckPDF_Small(b *testing.B) { harness.Bytes(b, harness.PDFSmall(b), check) }

func TestSurvey(t *testing.T) {
	harness.SurveyValidate(t, validate)
	harness.SurveyParse(t, parse)
	harness.SurveyExtract(t, extract)
}

func BenchmarkCheckPDF_Large(b *testing.B) { harness.Bytes(b, harness.PDFLarge(b), check) }

// TestCases validates the conformance cases (../../../cases) for the
// confidence scores: run by ./run.sh cases, which sets CASES_OUT.
func TestCases(t *testing.T) {
	harness.Cases(t, "dcnws-invoicekit.go", "github.com/dcnws/invoicekit.go", "einvoice.Validate (CEN 1.3.17 schematron interpreted in Go): fatal findings", func(x []byte) (int, []string, error) {
		r, err := einvoice.Validate(x, einvoice.Options{})
		if err != nil {
			return 0, nil, err
		}
		var rules []string
		for _, f := range r.Findings {
			if f.Flag == "fatal" {
				rules = append(rules, f.ID)
			}
		}
		return r.Fatal, rules, nil
	})
}
