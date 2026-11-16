package dotwavebench

import (
	"os"
	"path/filepath"
	"testing"

	einvoice "github.com/dotwavehq/go-einvoice"
	"github.com/dotwavehq/go-einvoice/cii"
	"github.com/dotwavehq/go-einvoice/zugferd"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/ritoon/factur-x/bench/ecosystem/common"
)

func init() { model.ConfigPath = "disable" } // pdfcpu would otherwise write in $HOME

func extract(p []byte) ([]byte, error)                { return zugferd.ExtractXML(p) }
func parse(x []byte) (*einvoice.ParsedInvoice, error) { return einvoice.Parse(x) }
func write(p *einvoice.ParsedInvoice) ([]byte, error) {
	return cii.NewSerializer().Serialize(&p.Invoice)
}

// EmbedXML only accepts file paths: the benchmark includes writing the host
// to a temp file and reading the result back (as a caller would have to).
func embed(p, x []byte) ([]byte, error) {
	dir := os.TempDir()
	in, out := filepath.Join(dir, "dw-in.pdf"), filepath.Join(dir, "dw-out.pdf")
	if err := os.WriteFile(in, p, 0o644); err != nil {
		return nil, err
	}
	if err := zugferd.EmbedXML(in, x, out); err != nil {
		return nil, err
	}
	return os.ReadFile(out)
}

func BenchmarkExtract_Small(b *testing.B)    { common.Extract(b, common.PDFSmall(b), extract) }
func BenchmarkExtract_Large(b *testing.B)    { common.Extract(b, common.PDFLarge(b), extract) }
func BenchmarkExtract_Minimum(b *testing.B)  { common.Extract(b, common.PDFMinimum(b), extract) }
func BenchmarkExtract_Extended(b *testing.B) { common.Extract(b, common.PDFExtended(b), extract) }
func BenchmarkParse_Small(b *testing.B)      { common.Bytes(b, common.XMLSmall(b), parse) }
func BenchmarkParse_Lines(b *testing.B)      { common.Bytes(b, common.XMLLines(b), parse) }
func BenchmarkParse_Large(b *testing.B)      { common.Bytes(b, common.XMLLarge(b), parse) }
func BenchmarkWrite_Small(b *testing.B)      { common.Write(b, common.XMLSmall(b), parse, write) }
func BenchmarkWrite_Lines(b *testing.B)      { common.Write(b, common.XMLLines(b), parse, write) }
func BenchmarkEmbed_Small(b *testing.B) {
	common.Embed(b, common.HostPDF(b), common.XMLSmall(b), embed, extract)
}

func TestSurvey(t *testing.T) {
	common.SurveyParse(t, parse)
	common.SurveyExtract(t, extract)
}
