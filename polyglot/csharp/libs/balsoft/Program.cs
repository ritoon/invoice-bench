// Balsoft.Hive.EInvoice 1.1.0 + Balsoft.Hive.EInvoice.Pdf 1.1.0 (Balsoft,
// Apache-2.0): a hand-written EN 16931 model with CII and UBL readers and
// writers, a validator coded in C# (business rules on the model, no XSD, no
// schematron) and hybrid PDF creation on PDFsharp 6.2.4.

using Balsoft.Hive.EInvoice;
using Balsoft.Hive.EInvoice.Cii;
using Balsoft.Hive.EInvoice.Pdf;
using Balsoft.Hive.EInvoice.Reading;
using Balsoft.Hive.EInvoice.Validation;
using Bench;

var lib = new Lib("balsoft", "Balsoft.Hive.EInvoice + Balsoft.Hive.EInvoice.Pdf", "1.1.0");

static byte[] Extract(byte[] pdf) =>
    HybridPdf.ExtractXml(pdf) ?? throw new InvalidDataException("ExtractXml returned null (no embedded XML)");

static ReadResult Parse(byte[] xml) => InvoiceReader.Read(xml);

// The validator works on the model; the stated totals that the model does not
// keep are compared with the computed ones by ReadResult.TotalsDiscrepancies.
static (ValidationResult Rules, IReadOnlyList<string> Totals) Validate(byte[] xml)
{
    var read = InvoiceReader.Read(xml);
    return (InvoiceValidator.Validate(read.Invoice, read.Profile ?? InvoiceProfile.EN16931), read.TotalsDiscrepancies());
}

static Findings Summarize((ValidationResult Rules, IReadOnlyList<string> Totals) r) =>
    Findings.Of(r.Rules.Errors.Select(i => i.RuleId).Concat(r.Totals.Select(t => t.Split(' ')[0])));

var noValidation = new WriterOptions { ValidateBeforeWriting = false };

return Harness.Run(lib, args, h =>
{
    h.Extract("HybridPdf.ExtractXml: /EmbeddedFiles name tree (factur-x.xml, xrechnung.xml, zugferd-invoice.xml, otherwise the first .xml), raw bytes", Extract);
    h.Parse("InvoiceReader.Read (XDocument then CII or UBL reader depending on the root namespace)", Parse);
    h.Validate(
        "InvoiceReader.Read + InvoiceValidator.Validate(profile read, else EN16931) + ReadResult.TotalsDiscrepancies: rules coded in C# on the model (BR-*, codes, XRechnung/Factur-X profiles), no XSD or schematron",
        Validate,
        Summarize);
    h.Write(
        "CiiWriter.Write(profile read) with ValidateBeforeWriting=false (enabled by default: the validation is measured by ParseValidate); model read by InvoiceReader.Read",
        Parse,
        r => CiiWriter.Write(r.Invoice, r.Profile ?? InvoiceProfile.EN16931, noValidation));
    h.Embed(
        "HybridPdf.Create(XML bytes, EN16931): copies the pages into a new PDFsharp document, factur-x.xml as associated file, sRGB OutputIntent, missing fonts taken from the system folders, XMP replaced by incremental update; dates = DateTimeOffset.Now (configurable, left at the default)",
        (host, xml) => HybridPdf.Create(host, xml, InvoiceProfile.EN16931),
        Extract);
    h.Unsupported(Op.CheckPdf, "no hybrid PDF check");
});
