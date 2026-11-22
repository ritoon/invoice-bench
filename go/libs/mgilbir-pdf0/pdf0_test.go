package pdf0bench

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/mgilbir/formalis"
	"github.com/mgilbir/pdf0"
	"github.com/ritoon/invoice-bench/go/common"
)

// pdf0 exposes the invoice XML only through ValidateFacturX, which also runs
// the container, PDF/A-3b and invoice rule checks: it is a full check.
func check(p []byte) (int, error) {
	doc, err := pdf0.Read(bytes.NewReader(p), int64(len(p)))
	if err != nil {
		return 0, err
	}
	res := pdf0.ValidateFacturX(doc)
	if len(res.XML) == 0 {
		return 0, errors.New("no xml")
	}
	return len(res.Violations), nil
}
func extractViaCheck(p []byte) ([]byte, error) {
	doc, err := pdf0.Read(bytes.NewReader(p), int64(len(p)))
	if err != nil {
		return nil, err
	}
	return pdf0.ValidateFacturX(doc).XML, nil
}
func validate(x []byte) (int, error) {
	r, err := formalis.Validate(context.Background(), x, formalis.ProfileEN16931)
	if err != nil {
		return 0, err
	}
	return len(r.Fatal()), nil
}
func embed(p, x []byte) ([]byte, error) {
	doc, err := pdf0.Read(bytes.NewReader(p), int64(len(p)))
	if err != nil {
		return nil, err
	}
	if err := pdf0.EmbedFacturX(doc, x, formalis.ProfileEN16931, "Facture"); err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if err := doc.Write(&out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func BenchmarkCheckPDF_Small(b *testing.B)      { common.Bytes(b, common.PDFSmall(b), check) }
func BenchmarkCheckPDF_Large(b *testing.B)      { common.Bytes(b, common.PDFLarge(b), check) }
func BenchmarkParseValidate_Small(b *testing.B) { common.Bytes(b, common.XMLSmall(b), validate) }
func BenchmarkParseValidate_Lines(b *testing.B) { common.Bytes(b, common.XMLLines(b), validate) }
func BenchmarkParseValidate_Large(b *testing.B) { common.Bytes(b, common.XMLLarge(b), validate) }
func BenchmarkEmbed_Small(b *testing.B) {
	common.Embed(b, common.HostPDF(b), common.XMLSmall(b), embed, extractViaCheck)
}

func TestSurvey(t *testing.T) {
	common.SurveyValidate(t, validate)
	common.SurveyExtract(t, extractViaCheck)
}

// TestCases validates the conformance cases (../../../cases) for the
// confidence scores: run by ./run.sh cases, which sets CASES_OUT.
func TestCases(t *testing.T) {
	// The bench validates at the Factur-X EN16931 profile; the cases, CII and
	// UBL, EN 16931 and XRechnung, go through the library's own router.
	common.Cases(t, "mgilbir-pdf0", "github.com/mgilbir/formalis", "formalis.Detect, then the validator it returns: fatal violations", func(x []byte) (int, []string, error) {
		det, err := formalis.Detect(x)
		if err != nil {
			return 0, nil, err
		}
		v := det.Validator()
		if v == nil {
			return 0, nil, errors.New("formalis.Detect: format not recognised")
		}
		r, err := v(context.Background(), x)
		if err != nil {
			return 0, nil, err
		}
		var rules []string
		for _, v := range r.Fatal() {
			rules = append(rules, v.Rule)
		}
		return len(rules), rules, nil
	})
}
