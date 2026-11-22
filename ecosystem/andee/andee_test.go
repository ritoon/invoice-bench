package andeebench

import (
	"testing"

	xinvoice "github.com/andeedotnet/go-xinvoice"
	xinvoicepdf "github.com/andeedotnet/go-xinvoice-pdf"
	"github.com/ritoon/invoice-bench/ecosystem/common"
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

func BenchmarkExtract_Small(b *testing.B)       { common.Extract(b, common.PDFSmall(b), extract) }
func BenchmarkExtract_Large(b *testing.B)       { common.Extract(b, common.PDFLarge(b), extract) }
func BenchmarkExtract_Minimum(b *testing.B)     { common.Extract(b, common.PDFMinimum(b), extract) }
func BenchmarkExtract_Extended(b *testing.B)    { common.Extract(b, common.PDFExtended(b), extract) }
func BenchmarkParse_Small(b *testing.B)         { common.Bytes(b, common.XMLSmall(b), parse) }
func BenchmarkParse_Lines(b *testing.B)         { common.Bytes(b, common.XMLLines(b), parse) }
func BenchmarkParse_Large(b *testing.B)         { common.Bytes(b, common.XMLLarge(b), parse) }
func BenchmarkParseValidate_Small(b *testing.B) { common.Bytes(b, common.XMLSmall(b), validate) }
func BenchmarkParseValidate_Lines(b *testing.B) { common.Bytes(b, common.XMLLines(b), validate) }
func BenchmarkParseValidate_Large(b *testing.B) { common.Bytes(b, common.XMLLarge(b), validate) }
func BenchmarkWrite_Small(b *testing.B)         { common.Write(b, common.XMLSmall(b), parse, write) }
func BenchmarkWrite_Lines(b *testing.B)         { common.Write(b, common.XMLLines(b), parse, write) }
func BenchmarkEmbed_Small(b *testing.B) {
	common.Embed(b, common.HostPDF(b), common.XMLSmall(b), embed, extract)
}

func TestSurvey(t *testing.T) {
	common.SurveyValidate(t, validate)
	common.SurveyParse(t, parse)
	common.SurveyExtract(t, extract)
}

// TestCases validates the conformance cases (../../cases) for the
// confidence scores: run by ./run.sh cases, which sets CASES_OUT.
func TestCases(t *testing.T) {
	common.Cases(t, "andee", "github.com/andeedotnet/go-xinvoice", "xinvoice.ParseXML + xinvoice.Validate: findings of severity error", func(x []byte) (int, []string, error) {
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
