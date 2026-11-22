package angelobench

import (
	"testing"

	"github.com/angelodlfrtr/go-invoice-generator/facturx"
	"github.com/ritoon/invoice-bench/go/harness"
)

func extract(p []byte) ([]byte, error)        { return facturx.Extract(p) }
func parse(x []byte) (*facturx.Parsed, error) { return facturx.Parse(x) }
func validate(x []byte) (int, error) { // Parse + Verify (arithmetic checks)
	p, err := parse(x)
	if err != nil {
		return 0, err
	}
	if err := p.Verify(); err != nil {
		return 1, nil
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

func TestSurvey(t *testing.T) {
	harness.SurveyValidate(t, validate)
	harness.SurveyParse(t, parse)
	harness.SurveyExtract(t, extract)
}

func BenchmarkParseValidate_Large(b *testing.B) { harness.Bytes(b, harness.XMLLarge(b), validate) }

// TestCases validates the conformance cases (../../../cases) for the
// confidence scores: run by ./run.sh cases, which sets CASES_OUT.
func TestCases(t *testing.T) {
	harness.Cases(t, "angelodlfrtr-go-invoice-generator", "github.com/angelodlfrtr/go-invoice-generator", "facturx.Parse + Parsed.Verify (arithmetic checks): one error, its message", func(x []byte) (int, []string, error) {
		p, err := parse(x)
		if err != nil {
			return 0, nil, err
		}
		if err := p.Verify(); err != nil {
			return 1, []string{err.Error()}, nil
		}
		return 0, nil, nil
	})
}
