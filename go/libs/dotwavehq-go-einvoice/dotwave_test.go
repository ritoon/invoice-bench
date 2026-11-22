package dotwavebench

import (
	"os"
	"path/filepath"
	"testing"

	einvoice "github.com/dotwavehq/go-einvoice"
	"github.com/dotwavehq/go-einvoice/cii"
	"github.com/dotwavehq/go-einvoice/zugferd"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/ritoon/invoice-bench/go/harness"
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

func BenchmarkExtract_Small(b *testing.B)    { harness.Extract(b, harness.PDFSmall(b), extract) }
func BenchmarkExtract_Large(b *testing.B)    { harness.Extract(b, harness.PDFLarge(b), extract) }
func BenchmarkExtract_Minimum(b *testing.B)  { harness.Extract(b, harness.PDFMinimum(b), extract) }
func BenchmarkExtract_Extended(b *testing.B) { harness.Extract(b, harness.PDFExtended(b), extract) }
func BenchmarkParse_Small(b *testing.B)      { harness.Bytes(b, harness.XMLSmall(b), parse) }
func BenchmarkParse_Lines(b *testing.B)      { harness.Bytes(b, harness.XMLLines(b), parse) }
func BenchmarkParse_Large(b *testing.B)      { harness.Bytes(b, harness.XMLLarge(b), parse) }
func BenchmarkWrite_Small(b *testing.B)      { harness.Write(b, harness.XMLSmall(b), parse, write) }
func BenchmarkWrite_Lines(b *testing.B)      { harness.Write(b, harness.XMLLines(b), parse, write) }
func BenchmarkEmbed_Small(b *testing.B) {
	harness.Embed(b, harness.HostPDF(b), harness.XMLSmall(b), embed, extract)
}

func TestSurvey(t *testing.T) {
	harness.SurveyParse(t, parse)
	harness.SurveyExtract(t, extract)
}
