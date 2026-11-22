// OfficeIMO.Invoicing 3.6.2 (+ .Validation and .Pdf 3.6.2, Evotec, MIT): one
// typed invoice model with bounded CII/UBL parsers and writers (EN 16931
// 1.3.16, XRechnung 3.0.2, Factur-X 1.09.2 / ZUGFeRD 2.5.2, Peppol BIS 3.0.21),
// model checks (semantics, arithmetic), and OfficeIMO.Pdf, a managed PDF
// engine that reads attachments, attaches Factur-X XML to an existing PDF and
// assesses PDF/A-3 + Factur-X readiness.
//
// Validation: OfficeIMO.Invoicing.Validation runs the XSD stage in .NET from
// the KoSIT archive (validator-configuration-xrechnung v2026-08-31, the one it
// pins; run.sh setup downloads it) but needs Java and SaxonJ-HE for the
// business-rule stage. That stage is left out (no JVM in this bench), so the
// validation measured here is XSD + the core model checks.

using Bench;
using OfficeIMO.Invoicing;
using OfficeIMO.Invoicing.Validation;
using OfficeIMO.Pdf;

var lib = new Lib("officeimo", "OfficeIMO.Invoicing + .Validation + .Pdf", "3.6.2");

string[] names = ["factur-x.xml", "zugferd-invoice.xml", "ZUGFeRD-invoice.xml", "xrechnung.xml"];

bool IsInvoice(PdfExtractedAttachment a) =>
    names.Any(n => string.Equals(n, a.FileName, StringComparison.OrdinalIgnoreCase)
        || string.Equals(n, a.UnicodeFileName, StringComparison.OrdinalIgnoreCase)
        || string.Equals(n, a.Name, StringComparison.OrdinalIgnoreCase));

byte[] Extract(byte[] pdf)
{
    using var ms = new MemoryStream(pdf, writable: false);
    var doc = PdfDocument.Load(ms);
    return doc.ExtractAttachments().FirstOrDefault(IsInvoice)?.Bytes
        ?? throw new InvalidDataException("no attachment under a standard invoice name");
}

static InvoiceReadResult Parse(byte[] xml) => InvoiceParser.Read(xml);

// The XSD stage only: without a SaxonInvoiceRulesRunner the business rules are NotRun.
InvoiceValidator? validator = null;
InvoiceValidator Validator()
{
    if (validator != null)
    {
        return validator;
    }
    var zip = Path.Combine(Corpus.Artefacts, "dl", "xrechnung-3.0.2-validator-configuration-2026-08-31.zip");
    if (!File.Exists(zip))
    {
        throw new FileNotFoundException($"{zip} missing: run ./run.sh setup");
    }
    return validator = new InvoiceValidator(InvoiceRuleBundle.Load(zip, null, null), null!);
}

(InvoiceModelValidationResult Model, InvoiceValidationReport Xsd) Validate(byte[] xml)
{
    var read = InvoiceParser.Read(xml);
    var release = read.Declaration?.Profile == InvoiceProfile.XRechnung
        ? InvoiceSpecificationRelease.XRechnung_3_0_2_2026_08_31
        : InvoiceSpecificationRelease.En16931_1_3_16;
    var model = InvoiceModelValidator.ValidateSource(read);
    var xsd = Validator().ValidateAsync(xml, release, CancellationToken.None).GetAwaiter().GetResult();
    return (model, xsd);
}

static Findings Summarize((InvoiceModelValidationResult Model, InvoiceValidationReport Xsd) r)
{
    if (r.Xsd.SchemaStatus is InvoiceValidationStatus.Failed or InvoiceValidationStatus.NotRun)
    {
        throw new InvalidOperationException($"XSD stage {r.Xsd.SchemaStatus}: {r.Xsd.Diagnostics.FirstOrDefault()?.Message}");
    }
    return Findings.Of(r.Model.Diagnostics.Where(d => d.Severity == InvoiceDiagnosticSeverity.Error).Select(d => d.Code)
        .Concat(r.Xsd.Diagnostics.Where(d => d.Severity == InvoiceDiagnosticSeverity.Error).Select(d => "schema:" + d.Code)));
}

static byte[] Write(InvoiceReadResult read) =>
    read.Write(new InvoiceXmlOptions(
        InvoiceSpecificationRelease.En16931_1_3_16,
        InvoiceSyntax.Cii,
        read.Declaration?.Profile ?? InvoiceProfile.En16931), false);

// The Factur-X helpers (AddFacturXInvoiceXml, XMP) belong to PdfCreateOptions,
// i.e. to PDFs OfficeIMO generates; an existing PDF only takes a generic
// associated file through the attachment editor.
static byte[] Embed(byte[] host, byte[] xml)
{
    using var ms = new MemoryStream(host, writable: false);
    var doc = PdfDocument.Load(ms);
    var file = new PdfEmbeddedFile("factur-x.xml", xml, "text/xml", PdfAssociatedFileRelationship.Alternative, "Factur-X invoice");
    // Replace when the host already carries a factur-x.xml (Add refuses duplicates).
    var result = doc.Attachments.Edit(session =>
    {
        if (session.Attachments.Any(a => string.Equals(a.FileName, file.FileName, StringComparison.OrdinalIgnoreCase)))
        {
            session.Replace(file.FileName, file);
        }
        else
        {
            session.Add(file);
        }
    });
    return result.ToBytes();
}

static PdfComplianceReadinessReport Check(byte[] pdf)
{
    using var ms = new MemoryStream(pdf, writable: false);
    return PdfDocument.Load(ms).AssessCompliance(PdfComplianceProfile.FacturX);
}

return Harness.Run(lib, args, h =>
{
    h.Extract(
        "OfficeIMO.Pdf PdfDocument.Load + ExtractAttachments: full parse of the PDF, /Names/EmbeddedFiles and /AF, decoded streams; first file named factur-x.xml, zugferd-invoice.xml, ZUGFeRD-invoice.xml or xrechnung.xml",
        Extract);
    h.Parse("InvoiceParser.Read (bounded CII/UBL parser, data not carried over is reported)", Parse);
    h.Validate(
        "InvoiceParser.Read + InvoiceModelValidator.ValidateSource (semantics, arithmetic) + InvoiceValidator without a runner: XSD stage only (CII D16B from the KoSIT archive v2026-08-31); EN 16931 / XRechnung business rules not run (Java + SaxonJ-HE required)",
        Validate,
        Summarize);
    h.Write("InvoiceReadResult.Write(EN 16931 1.3.16, CII, profile read), model read by InvoiceParser.Read; model checks before writing", Parse, Write);
    h.Embed(
        "adapter: OfficeIMO.Pdf PdfDocument.Load + Attachments.Edit (Add, or Replace if the host already carries a factur-x.xml) + ToBytes: generic associated file, without Factur-X XMP or OutputIntent (the Factur-X helpers only serve PDFs generated by OfficeIMO)",
        Embed,
        Extract);
    h.CheckPdf(
        "OfficeIMO.Pdf PdfDocument.Load + AssessCompliance(FacturX): check of the file against the PDF/A-3 + Factur-X requirements (XMP, associated file, OutputIntent, fonts…); the XML is not validated",
        Check);
});
