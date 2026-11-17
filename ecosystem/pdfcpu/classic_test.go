package pdfcpubench

import (
	"testing"

	"github.com/ritoon/zerechnung/bench/ecosystem/common"
)

func BenchmarkEmbedClassic_Small(b *testing.B) {
	common.Embed(b, common.HostPDFClassic(b), common.XMLSmall(b), embed, extract)
}
