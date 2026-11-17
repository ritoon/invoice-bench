package floor

import (
	"bytes"
	"encoding/xml"
	"io"
	"testing"

	"github.com/ritoon/zerechnung/bench/ecosystem/common"
)

// Floor: cost of merely tokenizing the XML with encoding/xml. Any parser built
// on encoding/xml cannot go below this.
func benchTokenScan(b *testing.B, s common.Sample) {
	b.SetBytes(int64(len(s.Data)))
	b.ReportAllocs()
	for b.Loop() {
		d := xml.NewDecoder(bytes.NewReader(s.Data))
		for {
			_, err := d.RawToken()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkFloor_XMLRawToken_Small(b *testing.B) { benchTokenScan(b, common.XMLSmall(b)) }
func BenchmarkFloor_XMLRawToken_Lines(b *testing.B) { benchTokenScan(b, common.XMLLines(b)) }
func BenchmarkFloor_XMLRawToken_Large(b *testing.B) { benchTokenScan(b, common.XMLLarge(b)) }

// Floor for PDF: copying the bytes once (memmove) — the absolute lower bound.
func BenchmarkFloor_Memcpy_PDFSmall(b *testing.B) {
	s := common.PDFSmall(b)
	dst := make([]byte, len(s.Data))
	b.SetBytes(int64(len(s.Data)))
	for b.Loop() {
		copy(dst, s.Data)
	}
}
