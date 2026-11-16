package gopdfattachbench

import (
	"bytes"
	"testing"

	"github.com/MarlinKuhn/gopdfattach"
	"github.com/ritoon/factur-x/bench/ecosystem/common"
)

func extract(p []byte) ([]byte, error) {
	x, _, err := gopdfattach.Extract(bytes.NewReader(p))
	return x, err
}
func embed(p, x []byte) ([]byte, error) {
	return gopdfattach.AttachFacturX(bytes.NewReader(x), bytes.NewReader(p), &gopdfattach.AttachConfig{
		ConformanceLevel: "EN 16931", AFRelationship: gopdfattach.AFAlternative})
}

func BenchmarkExtract_Small(b *testing.B)    { common.Extract(b, common.PDFSmall(b), extract) }
func BenchmarkExtract_Large(b *testing.B)    { common.Extract(b, common.PDFLarge(b), extract) }
func BenchmarkExtract_Minimum(b *testing.B)  { common.Extract(b, common.PDFMinimum(b), extract) }
func BenchmarkExtract_Extended(b *testing.B) { common.Extract(b, common.PDFExtended(b), extract) }
func BenchmarkEmbed_Small(b *testing.B) {
	common.Embed(b, common.HostPDF(b), common.XMLSmall(b), embed, extract)
}

func TestSurvey(t *testing.T) { common.SurveyExtract(t, extract) }
