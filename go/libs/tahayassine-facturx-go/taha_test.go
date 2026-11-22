package tahabench

import (
	"testing"

	"github.com/ritoon/invoice-bench/go/harness"
	"github.com/tahayassine/facturx-go/pkg/cii"
	"github.com/tahayassine/facturx-go/pkg/en16931"
	"github.com/tahayassine/facturx-go/pkg/facturx"
)

func extract(p []byte) ([]byte, error)         { return facturx.ExtractXML(p) }
func parse(x []byte) (*en16931.Invoice, error) { return cii.Decode(x) }
func validate(x []byte) (int, error) {
	inv, err := parse(x)
	if err != nil {
		return 0, err
	}
	return len(en16931.Validate(inv).Errors), nil
}
func write(inv *en16931.Invoice) ([]byte, error) { return cii.Encode(inv) }

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

func TestSurvey(t *testing.T) {
	harness.SurveyValidate(t, validate)
	harness.SurveyParse(t, parse)
	harness.SurveyExtract(t, extract)
}

// TestCases validates the conformance cases (../../../cases) for the
// confidence scores: run by ./run.sh cases, which sets CASES_OUT.
func TestCases(t *testing.T) {
	harness.Cases(t, "tahayassine-facturx-go", "github.com/tahayassine/facturx-go", "cii.Decode + en16931.Validate: errors", func(x []byte) (int, []string, error) {
		inv, err := parse(x)
		if err != nil {
			return 0, nil, err
		}
		var rules []string
		for _, e := range en16931.Validate(inv).Errors {
			rules = append(rules, e.Rule)
		}
		return len(rules), rules, nil
	})
}
