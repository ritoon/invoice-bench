package kanjobench

import (
	"testing"

	"github.com/cyprienbrisset/kanjo/pkg/pdfa"
	"github.com/cyprienbrisset/kanjo/pkg/read"
	"github.com/cyprienbrisset/kanjo/pkg/rules"
	"github.com/cyprienbrisset/kanjo/pkg/write"
	"github.com/ritoon/zerechnung/bench/ecosystem/common"

	_ "github.com/cyprienbrisset/kanjo/pkg/formats"
	_ "github.com/cyprienbrisset/kanjo/pkg/rules/all"
)

func extract(p []byte) ([]byte, error)     { x, _, _, err := pdfa.ExtractInvoiceXML(p); return x, err }
func parse(x []byte) (*read.Result, error) { return read.ReadBytes(x, "invoice.xml") }
func validate(x []byte) (int, error) {
	r, err := parse(x)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, f := range rules.Validate(r.Doc, "en16931").Findings {
		if f.Severity >= rules.SeverityError {
			n++
		}
	}
	return n, nil
}
func write_(r *read.Result) ([]byte, error) {
	return write.WriteBytes("cii", r.Doc, write.Options{Profile: write.ProfileEN16931})
}
func embed(p, x []byte) ([]byte, error) {
	r, err := pdfa.EmbedXML(p, x, "factur-x.xml")
	if err != nil {
		return nil, err
	}
	return r.PDF, nil
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
func BenchmarkWrite_Small(b *testing.B)         { common.Write(b, common.XMLSmall(b), parse, write_) }
func BenchmarkWrite_Lines(b *testing.B)         { common.Write(b, common.XMLLines(b), parse, write_) }
func BenchmarkEmbed_Small(b *testing.B) {
	common.Embed(b, common.HostPDF(b), common.XMLSmall(b), embed, extract)
}

func TestSurvey(t *testing.T) {
	common.SurveyValidate(t, validate)
	common.SurveyParse(t, parse)
	common.SurveyExtract(t, extract)
}
