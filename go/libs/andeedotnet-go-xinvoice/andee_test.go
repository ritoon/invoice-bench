package andeebench

import (
	"testing"

	xinvoice "github.com/andeedotnet/go-xinvoice"
	xinvoicepdf "github.com/andeedotnet/go-xinvoice-pdf"
	"github.com/ritoon/invoice-bench/go/harness"
)

func extract(p []byte) ([]byte, error)          { x, _, err := xinvoicepdf.ExtractXML(p); return x, err }
func parse(x []byte) (*xinvoice.Invoice, error) { return xinvoice.ParseXML(x) }
func validate(x []byte) (int, error) {
	inv, err := parse(x)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, f := range xinvoice.Validate(inv).Findings {
		if f.Severity == xinvoice.SeverityError {
			n++
		}
	}
	return n, nil
}
func write(inv *xinvoice.Invoice) ([]byte, error) { return inv.ToXML(xinvoice.CII) }
func embed(p, x []byte) ([]byte, error) {
	out, _, err := xinvoicepdf.Embed(p, x, &xinvoicepdf.Options{Profile: xinvoicepdf.ProfileEN16931})
	return out, err
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

func TestSurvey(t *testing.T) {
	harness.SurveyValidate(t, validate)
	harness.SurveyParse(t, parse)
	harness.SurveyExtract(t, extract)
}

// TestCases validates the conformance cases (../../../cases) for the
// confidence scores: run by ./run.sh cases, which sets CASES_OUT.
func TestCases(t *testing.T) {
	harness.Cases(t, "andeedotnet-go-xinvoice", "github.com/andeedotnet/go-xinvoice", "xinvoice.ParseXML + xinvoice.Validate: findings of severity error", func(x []byte) (int, []string, error) {
		inv, err := parse(x)
		if err != nil {
			return 0, nil, err
		}
		var rules []string
		for _, f := range xinvoice.Validate(inv).Findings {
			if f.Severity == xinvoice.SeverityError {
				rules = append(rules, f.Rule)
			}
		}
		return len(rules), rules, nil
	})
}
