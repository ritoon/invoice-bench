package fxbench

import (
	"testing"

	"github.com/ritoon/invoice-bench/go/common"
	"github.com/solarpush/fx/pkg/cii"
	"github.com/solarpush/fx/pkg/invoice"
	"github.com/solarpush/fx/pkg/pdf"
)

func extract(p []byte) ([]byte, error)         { return pdf.ExtractXML(p) }
func parse(x []byte) (*invoice.Invoice, error) { return cii.Parse(x) }
func validate(x []byte) (int, error) {
	inv, err := parse(x)
	if err != nil {
		return 0, err
	}
	return len(invoice.ValidateReport(inv).Errors()), nil
}
func write(inv *invoice.Invoice) ([]byte, error) { cp := *inv; return cii.Generate(&cp) } // Generate mutates its input

func BenchmarkExtract_Small(b *testing.B)       { common.Extract(b, common.PDFSmall(b), extract) }
func BenchmarkExtract_Large(b *testing.B)       { common.Extract(b, common.PDFLarge(b), extract) }
func BenchmarkExtract_Minimum(b *testing.B)     { common.Extract(b, common.PDFMinimum(b), extract) }
func BenchmarkExtract_Extended(b *testing.B)    { common.Extract(b, common.PDFExtended(b), extract) }
func BenchmarkParse_Small(b *testing.B)         { common.Bytes(b, common.XMLSmall(b), parse) }
func BenchmarkParse_Lines(b *testing.B)         { common.Bytes(b, common.XMLLines(b), parse) }
func BenchmarkParse_Large(b *testing.B)         { common.Bytes(b, common.XMLLarge(b), parse) }
func BenchmarkParseValidate_Small(b *testing.B) { common.Bytes(b, common.XMLSmall(b), validate) }
func BenchmarkParseValidate_Lines(b *testing.B) { common.Bytes(b, common.XMLLines(b), validate) }
func BenchmarkWrite_Small(b *testing.B)         { common.Write(b, common.XMLSmall(b), parse, write) }
func BenchmarkWrite_Lines(b *testing.B)         { common.Write(b, common.XMLLines(b), parse, write) }

func TestSurvey(t *testing.T) {
	common.SurveyValidate(t, validate)
	common.SurveyParse(t, parse)
	common.SurveyExtract(t, extract)
}

func BenchmarkParseValidate_Large(b *testing.B) { common.Bytes(b, common.XMLLarge(b), validate) }

// TestCases validates the conformance cases (../../../cases) for the
// confidence scores: run by ./run.sh cases, which sets CASES_OUT.
func TestCases(t *testing.T) {
	common.Cases(t, "solarpush-fx", "github.com/solarpush/fx", "cii.Parse + invoice.ValidateReport: errors", func(x []byte) (int, []string, error) {
		inv, err := parse(x)
		if err != nil {
			return 0, nil, err
		}
		var rules []string
		for _, e := range invoice.ValidateReport(inv).Errors() {
			rules = append(rules, e.Code)
		}
		return len(rules), rules, nil
	})
}
