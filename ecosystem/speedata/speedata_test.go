package speedatabench

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/ritoon/zerechnung/bench/ecosystem/common"
	"github.com/speedata/einvoice"
	pdf "github.com/speedata/pdfdisassembler"
)

// speedata/einvoice has no exported PDF extraction: we reproduce its CLI
// (cmd/einvoice/pdf.go) on top of github.com/speedata/pdfdisassembler.
func extract(data []byte) ([]byte, error) {
	r, err := pdf.Open(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var fallback *pdf.EmbeddedFile
	for i, f := range r.EmbeddedFiles() {
		switch strings.ToLower(f.Name) {
		case "factur-x.xml", "zugferd-invoice.xml", "xrechnung.xml":
		default:
			if strings.HasSuffix(strings.ToLower(f.Name), ".xml") && fallback == nil {
				fallback = &r.EmbeddedFiles()[i]
			}
			continue
		}
		return content(f)
	}
	if fallback != nil {
		return content(*fallback)
	}
	return nil, errors.New("no invoice xml")
}

func content(f pdf.EmbeddedFile) ([]byte, error) {
	ef, ok := f.Spec.Dict("EF")
	if !ok {
		return nil, errors.New("no /EF")
	}
	s, ok := ef.Stream("F")
	if !ok {
		return nil, errors.New("no /F stream")
	}
	return s.Content()
}

func benchExtract(b *testing.B, s common.Sample) {
	if x, err := extract(s.Data); err != nil || !common.LooksLikeCII(x) {
		b.Skipf("extraction failed on %s: %v", s.Name, err)
	}
	b.SetBytes(int64(len(s.Data)))
	b.ReportAllocs()
	for b.Loop() {
		if _, err := extract(s.Data); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExtract_Small(b *testing.B)    { benchExtract(b, common.PDFSmall(b)) }
func BenchmarkExtract_Large(b *testing.B)    { benchExtract(b, common.PDFLarge(b)) }
func BenchmarkExtract_Minimum(b *testing.B)  { benchExtract(b, common.PDFMinimum(b)) }
func BenchmarkExtract_Extended(b *testing.B) { benchExtract(b, common.PDFExtended(b)) }

func parse(x []byte) (*einvoice.Invoice, error) { return einvoice.ParseReader(bytes.NewReader(x)) }

func benchParse(b *testing.B, s common.Sample) {
	if _, err := parse(s.Data); err != nil {
		b.Skipf("parse failed on %s: %v", s.Name, err)
	}
	b.SetBytes(int64(len(s.Data)))
	b.ReportAllocs()
	for b.Loop() {
		if _, err := parse(s.Data); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParse_Small(b *testing.B) { benchParse(b, common.XMLSmall(b)) }
func BenchmarkParse_Lines(b *testing.B) { benchParse(b, common.XMLLines(b)) }
func BenchmarkParse_Large(b *testing.B) { benchParse(b, common.XMLLarge(b)) }

func validate(x []byte) (int, error) {
	inv, err := parse(x)
	if err != nil {
		return 0, err
	}
	if err := inv.Validate(); err != nil {
		var ve *einvoice.ValidationError
		if errors.As(err, &ve) {
			return len(ve.Violations()), nil
		}
		return 0, err
	}
	return 0, nil
}

func benchValidate(b *testing.B, s common.Sample) {
	if _, err := validate(s.Data); err != nil {
		b.Skipf("validate failed on %s: %v", s.Name, err)
	}
	b.SetBytes(int64(len(s.Data)))
	b.ReportAllocs()
	for b.Loop() {
		validate(s.Data)
	}
}

func BenchmarkParseValidate_Small(b *testing.B) { benchValidate(b, common.XMLSmall(b)) }
func BenchmarkParseValidate_Lines(b *testing.B) { benchValidate(b, common.XMLLines(b)) }
func BenchmarkParseValidate_Large(b *testing.B) { benchValidate(b, common.XMLLarge(b)) }

func benchWrite(b *testing.B, s common.Sample) {
	inv, err := parse(s.Data)
	if err != nil {
		b.Skip(err)
	}
	var buf bytes.Buffer
	if err := inv.Write(&buf); err != nil {
		b.Skip(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		buf.Reset()
		inv.Write(&buf)
	}
}

func BenchmarkWrite_Small(b *testing.B) { benchWrite(b, common.XMLSmall(b)) }
func BenchmarkWrite_Lines(b *testing.B) { benchWrite(b, common.XMLLines(b)) }

// Survey: false positives on the valid official CII samples + extraction coverage.
func TestSurvey(t *testing.T) {
	flagged, total := 0, 0
	for _, f := range common.CIIFiles(t) {
		data, _ := os.ReadFile(f)
		n, err := validate(data)
		total++
		if err != nil || n > 0 {
			flagged++
			t.Logf("FLAG %s: %d violations, err=%v", f[strings.LastIndex(f, "/")+1:], n, err)
		}
	}
	t.Logf("SURVEY validate: %d/%d valid samples flagged", flagged, total)
	ok, all := 0, 0
	for _, f := range common.PDFFiles(t) {
		data, _ := os.ReadFile(f)
		all++
		if x, err := extract(data); err == nil && common.LooksLikeCII(x) {
			ok++
		}
	}
	t.Logf("SURVEY extract: %d/%d PDFs", ok, all)
}
