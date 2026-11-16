package kanjobench

import (
	"bytes"
	"testing"

	"github.com/ritoon/factur-x/bench/ecosystem/common"
)

// TestEmbedOutputSize documents why BenchmarkEmbed_Small fails for this
// library: on a host that stores objects in object streams, the rewritten PDF
// loses those objects and comes out much smaller than its input.
func TestEmbedOutputSize(t *testing.T) {
	for _, h := range []common.Sample{common.HostPDF(t), common.HostPDFClassic(t)} {
		out, err := embed(h.Data, common.XMLSmall(t).Data)
		if err != nil {
			t.Fatalf("%s: %v", h.Name, err)
		}
		t.Logf("%s: in=%d out=%d objstm_in=%d", h.Name, len(h.Data), len(out), bytes.Count(h.Data, []byte("/ObjStm")))
	}
}
