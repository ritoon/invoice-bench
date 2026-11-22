package pdfcpubench

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/ritoon/invoice-bench/go/harness"
)

func conf() *model.Configuration {
	c := model.NewStatelessConfiguration()
	c.ValidationMode = model.ValidationRelaxed
	c.PreserveInfoDict = true
	return c
}

// ---------- Extraction (generic pdfcpu API) ----------

func extract(pdf []byte) ([]byte, error) {
	aa, err := api.ExtractAttachmentsRaw(context.Background(), bytes.NewReader(pdf), "", nil, conf())
	if err != nil {
		return nil, err
	}
	for _, a := range aa {
		if filepath.Ext(a.FileName) == ".xml" {
			return io.ReadAll(a.Reader)
		}
	}
	return nil, io.EOF
}

func benchExtract(b *testing.B, s harness.Sample) {
	if x, err := extract(s.Data); err != nil || !harness.LooksLikeCII(x) {
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

func BenchmarkExtract_Small(b *testing.B)    { benchExtract(b, harness.PDFSmall(b)) }
func BenchmarkExtract_Large(b *testing.B)    { benchExtract(b, harness.PDFLarge(b)) }
func BenchmarkExtract_Minimum(b *testing.B)  { benchExtract(b, harness.PDFMinimum(b)) }
func BenchmarkExtract_Extended(b *testing.B) { benchExtract(b, harness.PDFExtended(b)) }

// ---------- Host PDF used by every Embed benchmark ----------

// TestMakeHost writes host.pdf: the small Factur-X PDF without its invoice, a
// PDF/A-3 document into which every library injects one.
func TestMakeHost(t *testing.T) {
	s := harness.PDFSmall(t)
	out, err := withoutInvoice(s.Data)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(harness.TestData(), "host.pdf"), out, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("host.pdf: %d bytes", len(out))
}

// withoutInvoice removes from a Factur-X PDF everything that makes it one: the
// EmbeddedFiles tree, the catalog's /AF with the file specification and its
// stream, and the Factur-X part of the XMP. The PDF/A identification stays.
func withoutInvoice(pdf []byte) ([]byte, error) {
	c := context.Background()
	ctx, err := api.ReadContext(c, bytes.NewReader(pdf), conf())
	if err != nil {
		return nil, err
	}
	if _, err := ctx.RemoveAttachments(c, nil); err != nil {
		return nil, err
	}
	xrt := ctx.XRefTable
	cat, err := xrt.Catalog()
	if err != nil {
		return nil, err
	}
	if af := cat.Delete("AF"); af != nil {
		if err := xrt.DeleteObject(c, af); err != nil {
			return nil, err
		}
	}
	ir, ok := cat["Metadata"].(types.IndirectRef)
	if !ok {
		return nil, fmt.Errorf("catalog /Metadata: %T, want an indirect reference", cat["Metadata"])
	}
	e, ok := xrt.FindTableEntryForIndRef(&ir)
	if !ok {
		return nil, fmt.Errorf("catalog /Metadata: %s not found", ir)
	}
	sd, ok := e.Object.(types.StreamDict)
	if !ok {
		return nil, fmt.Errorf("catalog /Metadata: %T, want a stream", e.Object)
	}
	if err := sd.Decode(); err != nil {
		return nil, err
	}
	if sd.Content, err = xmpWithoutFacturX(sd.Content); err != nil {
		return nil, err
	}
	if err := sd.Encode(); err != nil {
		return nil, err
	}
	e.Object = sd
	var out bytes.Buffer
	if err := api.WriteContext(c, ctx, &out); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// Namespaces of the invoice properties in the XMP of Factur-X, ZUGFeRD 2 and
// ZUGFeRD 1 documents.
var invoiceXMPNamespaces = []string{"urn:factur-x:", "urn:zugferd:", "urn:ferd:"}

// xmpWithoutFacturX drops from an XMP packet every rdf:Description that holds
// invoice properties or declares their PDF/A extension schema (in the corpus
// PDF, that extension declares the Factur-X schema only).
func xmpWithoutFacturX(x []byte) ([]byte, error) {
	const open, closing = "<rdf:Description", "</rdf:Description>"
	var out []byte
	for {
		i := bytes.Index(x, []byte(open))
		if i < 0 {
			return append(out, x...), nil
		}
		gt := bytes.IndexByte(x[i:], '>')
		if gt < 0 {
			return nil, errors.New("XMP: unterminated rdf:Description")
		}
		end := i + gt + 1
		if x[end-2] != '/' {
			j := bytes.Index(x[end:], []byte(closing))
			if j < 0 {
				return nil, errors.New("XMP: rdf:Description without its end tag")
			}
			if bytes.Contains(x[end:end+j], []byte(open)) {
				return nil, errors.New("XMP: nested rdf:Description")
			}
			end += j + len(closing)
		}
		out = append(out, x[:i]...)
		if desc := x[i:end]; containsAny(desc, invoiceXMPNamespaces) {
			out = bytes.TrimRight(out, " \t\r\n")
		} else {
			out = append(out, desc...)
		}
		x = x[end:]
	}
}

func containsAny(b []byte, subs []string) bool {
	for _, s := range subs {
		if bytes.Contains(b, []byte(s)) {
			return true
		}
	}
	return false
}

// TestMakeHostClassic rewrites host.pdf with a classic xref table and no object
// streams, to separate libraries that cannot handle xref streams.
func TestMakeHostClassic(t *testing.T) {
	s := harness.HostPDF(t)
	cf := conf()
	cf.WriteObjectStream = false
	cf.WriteXRefStream = false
	ctx, err := api.ReadContext(context.Background(), bytes.NewReader(s.Data), cf)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := api.WriteContext(context.Background(), ctx, &out); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(harness.TestData(), "host_classic.pdf"), out.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("host_classic.pdf: %d bytes", out.Len())
}

// TestHostsStrict checks that both host PDFs pass pdfcpu's strict validation
// and hold no invoice, so that an embedding failure is the library's fault,
// not the host's.
func TestHostsStrict(t *testing.T) {
	for _, h := range []struct {
		s          harness.Sample
		xrefStream bool
	}{{harness.HostPDF(t), true}, {harness.HostPDFClassic(t), false}} {
		t.Run(h.s.Name, func(t *testing.T) {
			cf := model.NewDefaultConfiguration()
			cf.ValidationMode = model.ValidationStrict
			if err := api.Validate(context.Background(), bytes.NewReader(h.s.Data), cf, nil); err != nil {
				t.Errorf("strict validation: %v", err)
			}
			t.Run("NoInvoice", func(t *testing.T) { checkNoInvoice(t, h.s.Data, h.xrefStream) })
		})
	}
}

// checkNoInvoice checks that a host is a PDF/A-3 document without invoice: no
// /AF, no EmbeddedFiles tree, no file specification or embedded file left in
// the file, an XMP that keeps the PDF/A-3 identification but no Factur-X
// property nor extension schema, and the expected kind of cross-reference.
func checkNoInvoice(t *testing.T, pdf []byte, xrefStream bool) {
	ctx, err := api.ReadContext(context.Background(), bytes.NewReader(pdf), conf())
	if err != nil {
		t.Fatal(err)
	}
	if ctx.Read.UsingXRefStreams != xrefStream {
		t.Errorf("xref stream: %v, want %v", ctx.Read.UsingXRefStreams, xrefStream)
	}
	xrt := ctx.XRefTable
	cat, err := xrt.Catalog()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := cat.Find("Names"); ok {
		names, err := xrt.DereferenceDict(cat["Names"])
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := names.Find("EmbeddedFiles"); ok {
			t.Error("catalog /Names has an EmbeddedFiles tree")
		}
	}
	for nr, e := range xrt.Table {
		if e.Free || e.Object == nil {
			continue
		}
		var d types.Dict
		switch o := e.Object.(type) {
		case types.Dict:
			d = o
		case types.StreamDict:
			d = o.Dict
		default:
			continue
		}
		if _, ok := d.Find("AF"); ok {
			t.Errorf("object %d has /AF", nr)
		}
		if typ := d.Type(); typ != nil && (*typ == "Filespec" || *typ == "EmbeddedFile") {
			t.Errorf("object %d: /Type /%s", nr, *typ)
		}
	}
	sd, _, err := xrt.DereferenceStreamDict(cat["Metadata"])
	if err != nil || sd == nil {
		t.Fatalf("catalog /Metadata: %v", err)
	}
	if err := sd.Decode(); err != nil {
		t.Fatal(err)
	}
	x := sd.Content
	if !bytes.Contains(x, []byte(`pdfaid:part="3"`)) && !bytes.Contains(x, []byte("<pdfaid:part>3</pdfaid:part>")) {
		t.Error("XMP: no PDF/A-3 identification (pdfaid:part 3)")
	}
	for _, ns := range append(invoiceXMPNamespaces, "pdfaExtension:schemas") {
		if bytes.Contains(x, []byte(ns)) {
			t.Errorf("XMP: still contains %q", ns)
		}
	}
}

// ---------- Embedding: the low-level pdfcpu recipe (AF + FileSpec + XMP) ----------

const xmpTpl = `<?xpacket begin="" id="W5M0MpCehiHzreSzNTczkc9d"?>
<x:xmpmeta xmlns:x="adobe:ns:meta/"><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
<rdf:Description rdf:about="" xmlns:pdfaid="http://www.aiim.org/pdfa/ns/id/"><pdfaid:part>3</pdfaid:part><pdfaid:conformance>B</pdfaid:conformance></rdf:Description>
<rdf:Description rdf:about="" xmlns:fx="urn:factur-x:pdfa:CrossIndustryDocument:invoice:1p0#"><fx:DocumentType>INVOICE</fx:DocumentType><fx:DocumentFileName>factur-x.xml</fx:DocumentFileName><fx:Version>1.0</fx:Version><fx:ConformanceLevel>EN 16931</fx:ConformanceLevel></rdf:Description>
<rdf:Description rdf:about="" xmlns:pdfaExtension="http://www.aiim.org/pdfa/ns/extension/" xmlns:pdfaSchema="http://www.aiim.org/pdfa/ns/schema#" xmlns:pdfaProperty="http://www.aiim.org/pdfa/ns/property#"><pdfaExtension:schemas><rdf:Bag><rdf:li rdf:parseType="Resource"><pdfaSchema:schema>Factur-X PDFA Extension Schema</pdfaSchema:schema><pdfaSchema:namespaceURI>urn:factur-x:pdfa:CrossIndustryDocument:invoice:1p0#</pdfaSchema:namespaceURI><pdfaSchema:prefix>fx</pdfaSchema:prefix></rdf:li></rdf:Bag></pdfaExtension:schemas></rdf:Description>
</rdf:RDF></x:xmpmeta><?xpacket end="w"?>`

func embed(pdf, xml []byte) ([]byte, error) {
	c := context.Background()
	cf := conf()
	ctx, err := api.ReadContext(c, bytes.NewReader(pdf), cf)
	if err != nil {
		return nil, err
	}
	xrt := ctx.XRefTable
	if err := xrt.LocateNameTree("EmbeddedFiles", true); err != nil {
		return nil, err
	}
	sIR, err := xrt.NewEmbeddedStreamDict(bytes.NewReader(xml), time.Now())
	if err != nil {
		return nil, err
	}
	e, _ := xrt.FindTableEntryForIndRef(sIR)
	sd := e.Object.(types.StreamDict)
	sd.InsertName("Subtype", "text/xml")
	e.Object = sd
	fs, err := xrt.NewFileSpecDict("factur-x.xml", "factur-x.xml", "Factur-X invoice", *sIR)
	if err != nil {
		return nil, err
	}
	fs.InsertName("AFRelationship", "Alternative")
	fsIR, err := xrt.IndRefForNewObject(fs)
	if err != nil {
		return nil, err
	}
	if err := xrt.Names["EmbeddedFiles"].Add(c, xrt, "factur-x.xml", *fsIR, nil, []string{"F", "UF"}); err != nil {
		return nil, err
	}
	cat, err := xrt.Catalog()
	if err != nil {
		return nil, err
	}
	cat.Update("AF", types.Array{*fsIR})
	md := types.StreamDict{Dict: types.NewDict(), Content: []byte(xmpTpl)}
	md.InsertName("Type", "Metadata")
	md.InsertName("Subtype", "XML")
	if err := md.Encode(); err != nil {
		return nil, err
	}
	mdIR, err := xrt.IndRefForNewObject(md)
	if err != nil {
		return nil, err
	}
	cat.Update("Metadata", *mdIR)
	var out bytes.Buffer
	if err := api.Write(c, ctx, &out, cf); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func TestEmbedRoundTrip(t *testing.T) {
	host, x := harness.HostPDF(t), harness.XMLSmall(t)
	out, err := embed(host.Data, x.Data)
	if err != nil {
		t.Fatal(err)
	}
	back, err := extract(out)
	if err != nil || !bytes.Equal(back, x.Data) {
		t.Fatalf("round trip failed: %v", err)
	}
}

func BenchmarkEmbed_Small(b *testing.B) {
	host, x := harness.HostPDF(b), harness.XMLSmall(b)
	b.SetBytes(int64(len(host.Data)))
	b.ReportAllocs()
	for b.Loop() {
		if _, err := embed(host.Data, x.Data); err != nil {
			b.Fatal(err)
		}
	}
}

func TestSurvey(t *testing.T) { harness.SurveyExtract(t, extract) }
