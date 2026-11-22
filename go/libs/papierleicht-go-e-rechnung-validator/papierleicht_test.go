package papierleichtbench

import (
	"os"
	"sync"
	"testing"

	validator "github.com/papierleicht/go-e-rechnung-validator"
	"github.com/ritoon/invoice-bench/go/harness"
)

var (
	once sync.Once
	v    *validator.Validator
)

func val(tb testing.TB) *validator.Validator {
	once.Do(func() {
		dir, _ := os.MkdirTemp("", "erv-cache")
		var err error
		v, err = validator.New(validator.Options{CacheDir: dir})
		if err != nil {
			tb.Fatal(err)
		}
	})
	return v
}

func extract(p []byte) ([]byte, error) { return validator.ExtractXML(p) }

func BenchmarkExtract_Small(b *testing.B)    { harness.Extract(b, harness.PDFSmall(b), extract) }
func BenchmarkExtract_Large(b *testing.B)    { harness.Extract(b, harness.PDFLarge(b), extract) }
func BenchmarkExtract_Minimum(b *testing.B)  { harness.Extract(b, harness.PDFMinimum(b), extract) }
func BenchmarkExtract_Extended(b *testing.B) { harness.Extract(b, harness.PDFExtended(b), extract) }

func validate(tb testing.TB) func([]byte) (int, error) {
	vv := val(tb)
	return func(x []byte) (int, error) {
		r := vv.Validate(x)
		return len(r.Errors()), nil
	}
}

// Validation = XSD (official FeRD/KoSIT) + official Schematron compiled to XSLT.
// Rules are compiled once (warm-up excluded from timing).
func BenchmarkParseValidate_Small(b *testing.B) { harness.Bytes(b, harness.XMLSmall(b), validate(b)) }
func BenchmarkParseValidate_Lines(b *testing.B) { harness.Bytes(b, harness.XMLLines(b), validate(b)) }
func BenchmarkParseValidate_Large(b *testing.B) { harness.Bytes(b, harness.XMLLarge(b), validate(b)) }

func TestSurvey(t *testing.T) {
	harness.SurveyValidate(t, validate(t))
	harness.SurveyExtract(t, extract)
}

func BenchmarkCheckPDF_Small(b *testing.B) {
	vv := val(b)
	harness.Bytes(b, harness.PDFSmall(b), func(p []byte) (int, error) { return len(vv.ValidatePDF(p).Errors()), nil })
}
func BenchmarkCheckPDF_Large(b *testing.B) {
	vv := val(b)
	harness.Bytes(b, harness.PDFLarge(b), func(p []byte) (int, error) { return len(vv.ValidatePDF(p).Errors()), nil })
}

// TestCases validates the conformance cases (../../../cases) for the
// confidence scores: run by ./run.sh cases, which sets CASES_OUT.
func TestCases(t *testing.T) {
	vv := val(t)
	harness.Cases(t, "papierleicht-go-e-rechnung-validator", "github.com/papierleicht/go-e-rechnung-validator", "validator.Validate (format detected by the library): errors", func(x []byte) (int, []string, error) {
		r := vv.Validate(x)
		var rules []string
		for _, e := range r.Errors() {
			rules = append(rules, e.RuleID)
		}
		return len(rules), rules, nil
	})
}
