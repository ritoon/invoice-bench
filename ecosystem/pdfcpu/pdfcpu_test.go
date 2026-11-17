package pdfcpubench

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/ritoon/zerechnung/bench/ecosystem/common"
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

// ---------- Host PDF used by every Embed benchmark ----------

func TestMakeHost(t *testing.T) {
	s := common.PDFSmall(t)
	var out bytes.Buffer
	if err := api.RemoveAttachments(context.Background(), bytes.NewReader(s.Data), &out, nil, conf()); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join("..", "host.pdf"), out.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("host.pdf: %d bytes", out.Len())
}

// TestMakeHostClassic rewrites host.pdf with a classic xref table and no object
// streams, to separate libraries that cannot handle xref streams.
func TestMakeHostClassic(t *testing.T) {
	s := common.HostPDF(t)
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
	if err := os.WriteFile(filepath.Join("..", "host_classic.pdf"), out.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("host_classic.pdf: %d bytes", out.Len())
}

// TestHostsStrict checks that both host PDFs pass pdfcpu's strict validation,
// so that an embedding failure is the library's fault, not the host's.
func TestHostsStrict(t *testing.T) {
	for _, s := range []common.Sample{common.HostPDF(t), common.HostPDFClassic(t)} {
		cf := model.NewDefaultConfiguration()
		cf.ValidationMode = model.ValidationStrict
		if err := api.Validate(context.Background(), bytes.NewReader(s.Data), cf, nil); err != nil {
			t.Errorf("%s: %v", s.Name, err)
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
	host, x := common.HostPDF(t), common.XMLSmall(t)
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
	host, x := common.HostPDF(b), common.XMLSmall(b)
	b.SetBytes(int64(len(host.Data)))
	b.ReportAllocs()
	for b.Loop() {
		if _, err := embed(host.Data, x.Data); err != nil {
			b.Fatal(err)
		}
	}
}

func TestSurvey(t *testing.T) { common.SurveyExtract(t, extract) }
