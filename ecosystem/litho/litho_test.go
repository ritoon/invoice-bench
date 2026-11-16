package lithobench

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/gotenberg/litho/pkg/models"
	"github.com/gotenberg/litho/pkg/pdf"
	"github.com/ritoon/factur-x/bench/ecosystem/common"
)

var ctx = context.Background()

func open(p []byte) (*pdf.Reader, error) { return pdf.Open(ctx, bytes.NewReader(p), int64(len(p))) }

func extract(p []byte) ([]byte, error) {
	r, err := open(p)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	files, err := r.EmbeddedFiles(ctx)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		n := strings.ToLower(f.Name)
		if n == "factur-x.xml" || n == "zugferd-invoice.xml" || n == "xrechnung.xml" {
			rc, err := f.Open(ctx)
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, errors.New("no invoice xml")
}

func check(p []byte) (int, error) {
	r, err := open(p)
	if err != nil {
		return 0, err
	}
	defer r.Close()
	rep, err := r.ValidateFacturXReport(ctx)
	if err != nil {
		return 0, err
	}
	return len(rep.Fatals()), nil
}

func embed(p, x []byte) ([]byte, error) {
	r, err := open(p)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	w := pdf.NewWriter(r)
	defer w.Close()
	if _, err := w.RemediateFacturX(ctx, pdf.FacturXConfig{InvoiceXML: x, Level: models.FacturXEN16931}); err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if err := w.WriteTo(ctx, &out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func BenchmarkExtract_Small(b *testing.B)    { common.Extract(b, common.PDFSmall(b), extract) }
func BenchmarkExtract_Large(b *testing.B)    { common.Extract(b, common.PDFLarge(b), extract) }
func BenchmarkExtract_Minimum(b *testing.B)  { common.Extract(b, common.PDFMinimum(b), extract) }
func BenchmarkExtract_Extended(b *testing.B) { common.Extract(b, common.PDFExtended(b), extract) }
func BenchmarkCheckPDF_Small(b *testing.B)   { common.Bytes(b, common.PDFSmall(b), check) }
func BenchmarkCheckPDF_Large(b *testing.B)   { common.Bytes(b, common.PDFLarge(b), check) }
func BenchmarkEmbed_Small(b *testing.B) {
	common.Embed(b, common.HostPDF(b), common.XMLSmall(b), embed, extract)
}

func TestSurvey(t *testing.T) { common.SurveyExtract(t, extract) }
