package papierleichtbench

import (
	"os"
	"sync"
	"testing"

	validator "github.com/papierleicht/go-e-rechnung-validator"
	"github.com/ritoon/factur-x/bench/ecosystem/common"
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

func BenchmarkExtract_Small(b *testing.B)    { common.Extract(b, common.PDFSmall(b), extract) }
func BenchmarkExtract_Large(b *testing.B)    { common.Extract(b, common.PDFLarge(b), extract) }
func BenchmarkExtract_Minimum(b *testing.B)  { common.Extract(b, common.PDFMinimum(b), extract) }
func BenchmarkExtract_Extended(b *testing.B) { common.Extract(b, common.PDFExtended(b), extract) }

func validate(tb testing.TB) func([]byte) (int, error) {
	vv := val(tb)
	return func(x []byte) (int, error) {
		r := vv.Validate(x)
		return len(r.Errors()), nil
	}
}

// Validation = XSD (official FeRD/KoSIT) + official Schematron compiled to XSLT.
// Rules are compiled once (warm-up excluded from timing).
func BenchmarkParseValidate_Small(b *testing.B) { common.Bytes(b, common.XMLSmall(b), validate(b)) }
func BenchmarkParseValidate_Lines(b *testing.B) { common.Bytes(b, common.XMLLines(b), validate(b)) }
func BenchmarkParseValidate_Large(b *testing.B) { common.Bytes(b, common.XMLLarge(b), validate(b)) }

func TestSurvey(t *testing.T) {
	common.SurveyValidate(t, validate(t))
	common.SurveyExtract(t, extract)
}

func BenchmarkCheckPDF_Small(b *testing.B) {
	vv := val(b)
	common.Bytes(b, common.PDFSmall(b), func(p []byte) (int, error) { return len(vv.ValidatePDF(p).Errors()), nil })
}
func BenchmarkCheckPDF_Large(b *testing.B) {
	vv := val(b)
	common.Bytes(b, common.PDFLarge(b), func(p []byte) (int, error) { return len(vv.ValidatePDF(p).Errors()), nil })
}
