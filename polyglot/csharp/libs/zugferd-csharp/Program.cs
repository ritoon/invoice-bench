// ZUGFeRD-csharp 18.0.1 + ZUGFeRD.PDF-csharp 18.0.1 (Stephan Stapel, s2industries):
// the most downloaded .NET ZUGFeRD / Factur-X / XRechnung library. Model,
// readers and writers for ZUGFeRD 1.0, 2.0, 2.3 (CII) and XRechnung UBL;
// the .PDF package reads and writes hybrid PDFs on PDFsharp-extended.
//
// The PDF package has no public call that returns the XML bytes: LoadFromPdf
// extracts and parses into an InvoiceDescriptor, so Extract_* measure
// extraction + parsing (the check re-serialises the model outside the timing).

using Bench;
using s2industries.ZUGFeRD;
using s2industries.ZUGFeRD.PDF;

var lib = new Lib("zugferd-csharp", "ZUGFeRD-csharp + ZUGFeRD.PDF-csharp", "18.0.1");

static InvoiceDescriptor Parse(byte[] xml)
{
    using var ms = new MemoryStream(xml, writable: false);
    return InvoiceDescriptor.Load(ms);
}

static byte[] Write(InvoiceDescriptor d)
{
    using var ms = new MemoryStream();
    d.Save(ms, ZUGFeRDVersion.Version23, d.Profile, ZUGFeRDFormats.CII);
    return ms.ToArray();
}

static InvoiceDescriptor FromPdf(byte[] pdf)
{
    using var ms = new MemoryStream(pdf, writable: false);
    return InvoicePdfProcessor.LoadFromPdf(ms)
        ?? throw new InvalidDataException("LoadFromPdf returned null (no invoice attachment under a known name)");
}

// InvoiceValidator.Validate recomputes the monetary summation from the model:
// line totals, BR-CO-11..17, BR-DEC-20, totals and due amount. Its messages
// mix the computation trace with the failures; the failures are the lines
// that carry a rule id or a German "computed vs present" mismatch.
static Findings Summarize(ValidationResult r)
{
    if (r.IsValid)
    {
        return Findings.None;
    }
    var errors = r.Messages.Where(m =>
        ((m.StartsWith("BR-", StringComparison.Ordinal) || m.StartsWith("BT-", StringComparison.Ordinal))
            && !m.Contains(" matches ", StringComparison.Ordinal))
        || m.Contains("aber tatsächlich", StringComparison.Ordinal)
        || m.Contains("is not supported", StringComparison.Ordinal)
        || m.Contains(" Kein ", StringComparison.Ordinal)
        || m.StartsWith("Tax type code", StringComparison.Ordinal)
        || m.StartsWith("Invalid invoice", StringComparison.Ordinal))
        .Select(m => m.Length > 80 ? m[..80] : m)
        .ToList();
    return errors.Count == 0 ? Findings.Of(["IsValid=false"]) : Findings.Of(errors);
}

static byte[] Embed(byte[] host, byte[] xml)
{
    var d = Parse(xml);
    using var src = new MemoryStream(host, writable: false);
    using var dst = new MemoryStream();
    InvoicePdfProcessor.SaveToPdf(dst, ZUGFeRDVersion.Version23, Profile.Comfort, ZUGFeRDFormats.CII, src, d);
    return dst.ToArray();
}

return Harness.Run(lib, args, h =>
{
    h.Extract(
        "InvoicePdfProcessor.LoadFromPdf: extraction (catalog /AF, names zugferd.xml, zugferd-invoice.xml, factur-x.xml, xrechnung.xml) + parsing of the model; no public API that returns the bytes",
        FromPdf,
        d => Write(d));
    h.Parse("InvoiceDescriptor.Load (ZUGFeRD 1.0, 2.0, 2.3 CII, UBL readers)", Parse);
    h.Validate(
        "InvoiceDescriptor.Load + InvoiceValidator.Validate: recomputation of the totals only (lines, BR-CO-11 to BR-CO-17, BR-DEC-20); no XSD or schematron",
        x => InvoiceValidator.Validate(Parse(x), ZUGFeRDVersion.Version23),
        Summarize);
    h.Write("InvoiceDescriptor.Save(Version23, profile read, CII), model read by Load", Parse, Write);
    h.Embed(
        "adapter: Load + InvoicePdfProcessor.SaveToPdf(Version23, Comfort, CII); rewrites the XML from the model, copies the pages into a new PDFsharp document, Factur-X XMP, sRGB OutputIntent, system fonts embedded; dates = DateTime.UtcNow (not configurable)",
        Embed,
        pdf => Write(FromPdf(pdf)));
    h.Unsupported(Op.CheckPdf, "no hybrid PDF check");
});
