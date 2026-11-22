package pdf0bench

import (
	"testing"

	"github.com/ritoon/invoice-bench/ecosystem/common"
)

func BenchmarkEmbedClassic_Small(b *testing.B) {
	common.Embed(b, common.HostPDFClassic(b), common.XMLSmall(b), embed, extractViaCheck)
}
