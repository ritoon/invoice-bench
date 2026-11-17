package dcnwsbench

import (
	"testing"

	einvoice "github.com/dcnws/invoicekit.go"
	"github.com/dcnws/invoicekit.go/facturx"
	"github.com/dcnws/invoicekit.go/hybrid"
	"github.com/ritoon/zerechnung/bench/ecosystem/common"
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
func BenchmarkCheckPDF_Small(b *testing.B) { common.Bytes(b, common.PDFSmall(b), check) }

func TestSurvey(t *testing.T) {
	common.SurveyValidate(t, validate)
	common.SurveyParse(t, parse)
	common.SurveyExtract(t, extract)
}

func BenchmarkCheckPDF_Large(b *testing.B) { common.Bytes(b, common.PDFLarge(b), check) }
