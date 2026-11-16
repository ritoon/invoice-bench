// Benchmarks of this repository's own package, github.com/ritoon/factur-x,
// measured with the same harness as the other libraries.
package facturxbench

import (
	"errors"
	"testing"
	"time"

	facturx "github.com/ritoon/factur-x"
	"github.com/ritoon/factur-x/bench/ecosystem/common"
	"github.com/ritoon/factur-x/cii"
	"github.com/ritoon/factur-x/invoice"
	"github.com/ritoon/factur-x/validate"
)

func extract(p []byte) ([]byte, error) {
	att, err := facturx.Extract(p)
	if err != nil {
		return nil, err
	}
	return att.XML, nil
}

func parse(x []byte) (*invoice.Invoice, error) { return cii.Parse(x) }

func write(inv *invoice.Invoice) ([]byte, error) { return cii.Marshal(inv, nil) }

// validateXML parses and validates with the default rule sets (schema,
// EN 16931, profile rules; BR-DE for XRECHNUNG invoices) and returns the
// number of errors.
func validateXML(x []byte) (int, error) {
	r := validate.XML(x, nil)
	return len(r.Errors()), nil
}

var modDate = time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)

func embed(p, x []byte) ([]byte, error) {
	return facturx.Embed(p, x, &facturx.EmbedOptions{ModDate: modDate})
}

func check(p []byte) (int, error) {
	r, att, err := facturx.Check(p, nil)
	if err != nil {
		return 0, err
	}
	if att == nil {
		return 0, errors.New("no invoice")
	}
	return len(r.Errors()), nil
}

func BenchmarkExtract_Small(b *testing.B)    { common.Extract(b, common.PDFSmall(b), extract) }
func BenchmarkExtract_Large(b *testing.B)    { common.Extract(b, common.PDFLarge(b), extract) }
func BenchmarkExtract_Minimum(b *testing.B)  { common.Extract(b, common.PDFMinimum(b), extract) }
func BenchmarkExtract_Extended(b *testing.B) { common.Extract(b, common.PDFExtended(b), extract) }

func BenchmarkParse_Small(b *testing.B) { common.Bytes(b, common.XMLSmall(b), parse) }
func BenchmarkParse_Lines(b *testing.B) { common.Bytes(b, common.XMLLines(b), parse) }
func BenchmarkParse_Large(b *testing.B) { common.Bytes(b, common.XMLLarge(b), parse) }

func BenchmarkParseValidate_Small(b *testing.B) { common.Bytes(b, common.XMLSmall(b), validateXML) }
func BenchmarkParseValidate_Lines(b *testing.B) { common.Bytes(b, common.XMLLines(b), validateXML) }
func BenchmarkParseValidate_Large(b *testing.B) { common.Bytes(b, common.XMLLarge(b), validateXML) }

func BenchmarkWrite_Small(b *testing.B) { common.Write(b, common.XMLSmall(b), parse, write) }
func BenchmarkWrite_Lines(b *testing.B) { common.Write(b, common.XMLLines(b), parse, write) }

func BenchmarkEmbed_Small(b *testing.B) {
	common.Embed(b, common.HostPDF(b), common.XMLSmall(b), embed, extract)
}

func BenchmarkEmbedClassic_Small(b *testing.B) {
	common.Embed(b, common.HostPDFClassic(b), common.XMLSmall(b), embed, extract)
}

func BenchmarkCheckPDF_Small(b *testing.B) { common.Bytes(b, common.PDFSmall(b), check) }
func BenchmarkCheckPDF_Large(b *testing.B) { common.Bytes(b, common.PDFLarge(b), check) }

func TestSurvey(t *testing.T) {
	common.SurveyValidate(t, validateXML)
	common.SurveyParse(t, parse)
	common.SurveyExtract(t, extract)
}
