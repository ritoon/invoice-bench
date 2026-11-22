// International.EInvoicing 1.1.0 (Epic Dice Studio, MIT): a managed model of
// EN 16931 with UBL, CII, Factur-X / ZUGFeRD and XRechnung readers and
// writers, and its own Schematron engine (an XPath subset evaluated over the
// preprocessed official .sch files, no XSLT). Packages used here:
//
// - International.EInvoicing: facade, CII/UBL readers and writers, EN 16931
//   rules (CEN artefacts 1.3.16, embedded) through AddDefaults();
// - .Validation.XRechnung: KoSIT XRechnung CII/UBL rules (embedded);
// - .Validation.Xsd: UN/CEFACT CII D16B schemas (embedded);
// - .Countries.Germany: the XRechnung profiles;
// - .FacturX.PdfSharp: reading and writing the PDF part on PDFsharp 6.2.4.
//
// The Factur-X profile rules (FeRD schematron) and the French rules are
// fetched, not shipped (AddFacturXRulesFrom, AddFrance): they are not used, so
// the validation is the most complete offline set the packages carry.

using System.Text;
using Bench;
using International.EInvoicing;
using International.EInvoicing.Countries.Germany;
using International.EInvoicing.FacturX;
using International.EInvoicing.FacturX.Pdf;
using International.EInvoicing.FacturX.PdfSharp;
using International.EInvoicing.Model;
using International.EInvoicing.Validation;
using International.EInvoicing.Validation.XRechnung;
using International.EInvoicing.Validation.Xsd;

var lib = new Lib("international-einvoicing", "International.EInvoicing", "1.1.0");

// One instance, built once (as the README advises): rule sets are loaded and
// compiled lazily on the first Validate, which lands in first_call_ns.
var pdfReader = new PdfSharpAttachmentReader();
var pdfWriter = new PdfSharpAttachmentWriter();
var einvoicing = EInvoicing.Create(
    e => e.AddDefaults().AddGermany().AddXRechnungRules().AddCiiSchema(),
    pdfReader);

const long MaxAttachmentBytes = 64L * 1024 * 1024; // the library's default limit

byte[] Extract(byte[] pdf)
{
    using var ms = new MemoryStream(pdf, writable: false);
    return pdfReader.FindAttachment(ms, FacturXAttachment.KnownFileNames, MaxAttachmentBytes)?.Xml
        ?? throw new InvalidDataException("FindAttachment returned null (no factur-x.xml, zugferd-invoice.xml or xrechnung.xml)");
}

EInvoice Parse(byte[] xml)
{
    var read = einvoicing.Read(xml);
    if (read.TryGetInvoice(out var invoice))
    {
        return invoice;
    }
    var first = read.Errors.FirstOrDefault();
    throw new InvalidDataException($"Read: no invoice ({read.Kind}){(first is null ? "" : ": " + first)}");
}

string Text(byte[] xml)
{
    using var reader = new StreamReader(new MemoryStream(xml, writable: false), Encoding.UTF8, detectEncodingFromByteOrderMarks: true);
    return reader.ReadToEnd();
}

static Findings Summarize(ValidationReport r) =>
    Findings.Of(r.Errors.Select(m => m.RuleIdentifier));

byte[] Embed(byte[] host, byte[] xml)
{
    using var src = new MemoryStream(host, writable: false);
    using var dst = new MemoryStream();
    pdfWriter.Attach(src, new FacturXAttachment(FacturXAttachment.FacturXFileName, xml), FacturXProfiles.En16931, dst);
    return dst.ToArray();
}

return Harness.Run(lib, args, h =>
{
    h.Extract(
        "PdfSharpAttachmentReader.FindAttachment (PDFsharp 6.2.4): names factur-x.xml, zugferd-invoice.xml, ZUGFeRD-invoice.xml, xrechnung.xml, raw bytes",
        Extract);
    h.Parse("EInvoicing.Read(byte[]) → EInvoice (syntax and profile detection, CII reader)", Parse);
    h.Validate(
        "EInvoicing.Validate(string): CII D16B XSD + EN 16931 CII 1.3.16 rules + XRechnung CII (KoSIT, embedded) if BT-24 is XRechnung, in-house Schematron engine (XPath subset); Factur-X and French rules not embedded (to be downloaded), hence absent",
        x => einvoicing.Validate(Text(x)),
        Summarize);
    h.Write("EInvoicing.Write(EInvoice, DocumentFormat.Cii), model read by Read", Parse, m => einvoicing.Write(m, DocumentFormat.Cii));
    h.Embed(
        "PdfSharpAttachmentWriter.Attach(EN 16931 profile): factur-x.xml associated file (/AF, name tree), Factur-X XMP added or completed, XML embedded as is; ModDate = DateTimeOffset.UtcNow",
        Embed,
        Extract);
    h.Unsupported(Op.CheckPdf, "no hybrid PDF check: parsing a PDF only compares the XMP with the XML (warnings), no PDF/A");
});
