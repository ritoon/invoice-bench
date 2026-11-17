package gobldebench

import (
	"testing"

	xinvoice "github.com/invopop/gobl.de.xinvoice"
	"github.com/ritoon/zerechnung/bench/ecosystem/common"
)

func parse(x []byte) (*xinvoice.Parsed, error) { return xinvoice.Parse(x) }
func validate(x []byte) (int, error) {
	p, err := parse(x)
	if err != nil {
		return 0, err
	}
	if err := p.Envelope.Validate(); err != nil {
		return 1, nil
	}
	return 0, nil
}
func write(p *xinvoice.Parsed) ([]byte, error) {
	d, err := xinvoice.Convert(p.Envelope, xinvoice.FormatZUGFeRD)
	if err != nil {
		return nil, err
	}
	return d.Data, nil
}

func BenchmarkParse_Small(b *testing.B)         { common.Bytes(b, common.XMLSmall(b), parse) }
func BenchmarkParse_Lines(b *testing.B)         { common.Bytes(b, common.XMLLines(b), parse) }
func BenchmarkParse_Large(b *testing.B)         { common.Bytes(b, common.XMLLarge(b), parse) }
func BenchmarkParseValidate_Small(b *testing.B) { common.Bytes(b, common.XMLSmall(b), validate) }
func BenchmarkParseValidate_Lines(b *testing.B) { common.Bytes(b, common.XMLLines(b), validate) }
func BenchmarkParseValidate_Large(b *testing.B) { common.Bytes(b, common.XMLLarge(b), validate) }
func BenchmarkWrite_Small(b *testing.B)         { common.Write(b, common.XMLSmall(b), parse, write) }
func BenchmarkWrite_Lines(b *testing.B)         { common.Write(b, common.XMLLines(b), parse, write) }

func TestSurvey(t *testing.T) {
	common.SurveyValidate(t, validate)
	common.SurveyParse(t, parse)
}
