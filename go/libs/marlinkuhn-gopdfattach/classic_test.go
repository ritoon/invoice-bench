package gopdfattachbench

import (
	"testing"

	"github.com/ritoon/invoice-bench/go/harness"
)

func BenchmarkEmbedClassic_Small(b *testing.B) {
	harness.Embed(b, harness.HostPDFClassic(b), harness.XMLSmall(b), embed, extract)
}
