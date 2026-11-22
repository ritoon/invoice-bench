// FacturXDotNet 0.5.9-alpha (FacturX-NET, MIT): Factur-X documents (PDF + XMP +
// factur-x.xml) with a CII model limited to the MINIMUM profile header (seller,
// buyer, totals), a TurboXml streaming reader, a writer, and business rules
// coded in C#: 19 CII rules of the MINIMUM profile (BR-01..09, BR-13..15, BR-CO-09,
// BR-CO-26, BR-DEC-12/13/14/18, BR-S-01) and 16 hybrid rules (BR-HYBRID-01..15,
// BR-HYBRID-FR-01). PDF on its fork FacturXDotNet.PDFsharp 6.2.0-preview-3.
// The library API is asynchronous; the adapter waits on each task.

using Bench;
using FacturXDotNet;
using FacturXDotNet.Generation.CII;
using FacturXDotNet.Generation.FacturX;
using FacturXDotNet.Models.CII;
using FacturXDotNet.Parsing.CII;
using FacturXDotNet.Validation;

var lib = new Lib("facturxdotnet", "FacturXDotNet", "0.5.9-alpha");

static byte[] Extract(byte[] pdf)
{
    var doc = FacturXDocument.LoadFromBuffer(pdf);
    var attachment = doc.GetCrossIndustryInvoiceAttachmentAsync().GetAwaiter().GetResult()
        ?? throw new InvalidDataException("no factur-x.xml attachment");
    return attachment.ReadAsync().GetAwaiter().GetResult().ToArray();
}

static CrossIndustryInvoice Parse(byte[] xml)
{
    using var ms = new MemoryStream(xml, writable: false);
    return new CrossIndustryInvoiceReader().Read(ms);
}

var ciiValidator = new CrossIndustryInvoiceValidator();

static byte[] Write(CrossIndustryInvoice cii)
{
    using var ms = new MemoryStream();
    new CrossIndustryInvoiceWriter().WriteAsync(ms, cii).GetAwaiter().GetResult();
    return ms.ToArray();
}

static byte[] Embed(byte[] host, byte[] xml)
{
    using var pdf = new MemoryStream(host, writable: false);
    using var cii = new MemoryStream(xml, writable: false);
    var doc = FacturXDocument.Create().WithBasePdf(pdf).WithCrossIndustryInvoice(cii).BuildAsync().GetAwaiter().GetResult();
    return doc.Data.ToArray();
}

// The French hybrid rule is left out, as everywhere in the protocol.
var options = new FacturXValidationOptions();
options.RulesToSkip.Add("BR-HYBRID-FR-01");
var pdfValidator = new FacturXValidator(options);

return Harness.Run(lib, args, h =>
{
    h.Extract(
        "FacturXDocument.LoadFromBuffer + GetCrossIndustryInvoiceAttachmentAsync + ReadAsync: name factur-x.xml only (default parameter), raw bytes",
        Extract);
    h.Parse("CrossIndustryInvoiceReader.Read (TurboXml): model limited to the header of the MINIMUM profile", Parse);
    h.Validate(
        "CrossIndustryInvoiceReader.Read + CrossIndustryInvoiceValidator.GetValidationResult: 19 CII rules coded in C# (MINIMUM profile: BR-01 to BR-09, BR-13 to BR-15, BR-CO-09, BR-CO-26, BR-DEC-12, 13, 14 and 18, BR-S-01); no XSD or schematron",
        x => ciiValidator.GetValidationResult(Parse(x)),
        r => Findings.Of(r.Rules.Where(x => x.HasFailed).Select(x => x.Rule.Name)));
    h.Write("CrossIndustryInvoiceWriter.WriteAsync: rewrites the MINIMUM model only (lines, VAT and most terms lost), model read by CrossIndustryInvoiceReader", Parse, Write);
    h.Embed(
        "FacturXDocument.Create().WithBasePdf(PDF).WithCrossIndustryInvoice(XML).BuildAsync(): factur-x.xml as associated file, Factur-X XMP generated from the XML, PDFsharp fork",
        Embed,
        Extract);
    h.CheckPdf(
        "FacturXValidator.GetValidationResultAsync: 15 hybrid rules (XMP, file name and relationship, profile) + 19 MINIMUM CII rules; BR-HYBRID-FR-01 left out; no PDF/A check",
        pdf => pdfValidator.GetValidationResultAsync(FacturXDocument.LoadFromBuffer(pdf)).GetAwaiter().GetResult());
});
