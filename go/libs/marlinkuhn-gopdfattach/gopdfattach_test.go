package gopdfattachbench

import (
	"bytes"
	"testing"

	"github.com/MarlinKuhn/gopdfattach"
	"github.com/ritoon/invoice-bench/go/harness"
)

func extract(p []byte) ([]byte, error) {
	x, _, err := gopdfattach.Extract(bytes.NewReader(p))
	return x, err
}
func embed(p, x []byte) ([]byte, error) {
	return gopdfattach.AttachFacturX(bytes.NewReader(x), bytes.NewReader(p), &gopdfattach.AttachConfig{
		ConformanceLevel: "EN 16931", AFRelationship: gopdfattach.AFAlternative})
}

func BenchmarkExtract_Small(b *testing.B)    { harness.Extract(b, harness.PDFSmall(b), extract) }
func BenchmarkExtract_Large(b *testing.B)    { harness.Extract(b, harness.PDFLarge(b), extract) }
func BenchmarkExtract_Minimum(b *testing.B)  { harness.Extract(b, harness.PDFMinimum(b), extract) }
func BenchmarkExtract_Extended(b *testing.B) { harness.Extract(b, harness.PDFExtended(b), extract) }
func BenchmarkEmbed_Small(b *testing.B) {
	harness.Embed(b, harness.HostPDF(b), harness.XMLSmall(b), embed, extract)
}

func TestSurvey(t *testing.T) { harness.SurveyExtract(t, extract) }
