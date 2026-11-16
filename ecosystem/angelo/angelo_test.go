package angelobench

import (
	"testing"

	"github.com/angelodlfrtr/go-invoice-generator/facturx"
	"github.com/ritoon/factur-x/bench/ecosystem/common"
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

func BenchmarkExtract_Small(b *testing.B)       { common.Extract(b, common.PDFSmall(b), extract) }
func BenchmarkExtract_Large(b *testing.B)       { common.Extract(b, common.PDFLarge(b), extract) }
func BenchmarkExtract_Minimum(b *testing.B)     { common.Extract(b, common.PDFMinimum(b), extract) }
func BenchmarkExtract_Extended(b *testing.B)    { common.Extract(b, common.PDFExtended(b), extract) }
func BenchmarkParse_Small(b *testing.B)         { common.Bytes(b, common.XMLSmall(b), parse) }
func BenchmarkParse_Lines(b *testing.B)         { common.Bytes(b, common.XMLLines(b), parse) }
func BenchmarkParse_Large(b *testing.B)         { common.Bytes(b, common.XMLLarge(b), parse) }
func BenchmarkParseValidate_Small(b *testing.B) { common.Bytes(b, common.XMLSmall(b), validate) }
func BenchmarkParseValidate_Lines(b *testing.B) { common.Bytes(b, common.XMLLines(b), validate) }

func TestSurvey(t *testing.T) {
	common.SurveyValidate(t, validate)
	common.SurveyParse(t, parse)
	common.SurveyExtract(t, extract)
}

func BenchmarkParseValidate_Large(b *testing.B) { common.Bytes(b, common.XMLLarge(b), validate) }
