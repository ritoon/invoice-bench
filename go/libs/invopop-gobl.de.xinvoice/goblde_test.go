package gobldebench

import (
	"testing"

	xinvoice "github.com/invopop/gobl.de.xinvoice"
	"github.com/ritoon/invoice-bench/go/harness"
)

func parse(x []byte) (*xinvoice.Parsed, error) { return xinvoice.Parse(x) }
func validate(x []byte) (int, error) {
	p, err := parse(x)
	if err != nil {
		return 0, err
	}
	if err := p.Envelope.Validate(); err != nil {
		return 1, nil
	}
	return 0, nil
}
func write(p *xinvoice.Parsed) ([]byte, error) {
	d, err := xinvoice.Convert(p.Envelope, xinvoice.FormatZUGFeRD)
	if err != nil {
		return nil, err
	}
	return d.Data, nil
}

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
}

// TestCases validates the conformance cases (../../../cases) for the
// confidence scores: run by ./run.sh cases, which sets CASES_OUT.
func TestCases(t *testing.T) {
	harness.Cases(t, "invopop-gobl.de.xinvoice", "github.com/invopop/gobl.de.xinvoice", "xinvoice.Parse + GOBL Envelope.Validate: one error, its message", func(x []byte) (int, []string, error) {
		p, err := parse(x)
		if err != nil {
			return 0, nil, err
		}
		if err := p.Envelope.Validate(); err != nil {
			return 1, []string{err.Error()}, nil
		}
		return 0, nil, nil
	})
}
