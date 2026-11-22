// Antiphon.Invoicing 1.0.0 (Clockwork Otter Foundry, commercial, closed source;
// the evaluation licence is free, fully functional, needs no key and has no
// time limit). The suite is split into packages: a typed EN 16931 model
// (Antiphon.Invoicing), a CII parser and generator, Factur-X PDF extraction and
// embedding on PDFsharp 6.2.4, and a layered validator (well-formedness, XSD,
// Schematron) whose Schematron rule sets are shipped pre-compiled to an
// intermediate representation run by its own XPath 2 evaluator (no XSLT, no
// Java). Rule sets used here: CEN EN 16931 CII and KoSIT XRechnung CII. There
// is no Factur-X profile rule set (FeRD) and no French one; the Peppol CII set
// is a separate package that no example of the corpus implies.
//
// Every call goes through the plain (non-licence) entry points: no licence key
// is read or checked.

using System.Text;
using System.Xml.Linq;
using Antiphon.Generation.Cii;
using Antiphon.Parsing.Cii;
using Antiphon.Pdf;
using Antiphon.Rules.Ir;
using Antiphon.Validation.Model;
using Antiphon.Validation.Pipeline;
using Bench;
using CenCii = Antiphon.Rules.Cen.Cii.RuleSetInfo;
using XrCii = Antiphon.Rules.XRechnung.Cii.RuleSetInfo;

var lib = new Lib("antiphon", "Antiphon.Invoicing (+ Parsing.Cii, Generation.Cii, Validation.Pipeline, Rules.Cen.Cii, Rules.XRechnung.Cii, Xsd.Cii, Pdf)", "1.0.0");

// The compiled rule sets are deserialised once, on the first validation (part
// of first_call_ns), then reused, as a caller would keep them.
Dictionary<string, RuleSetIr>? ruleSets = null;
Dictionary<string, RuleSetIr> RuleSets() => ruleSets ??= new()
{
    [CenCii.SourceId] = IrSerializer.Deserialize(CenCii.GetCompiledIrJson()),
    [XrCii.SourceId] = IrSerializer.Deserialize(XrCii.GetCompiledIrJson()),
};

static byte[] Extract(byte[] pdf)
{
    using var ms = new MemoryStream(pdf, writable: false);
    return FacturXExtractor.ExtractBytes(ms)
        ?? throw new InvalidDataException("ExtractBytes returned null (no attachment under a known name)");
}

static ParseResult Parse(byte[] xml)
{
    using var ms = new MemoryStream(xml, writable: false);
    return CiiInvoiceParser.Parse(XDocument.Load(ms));
}

static string Text(byte[] xml)
{
    using var reader = new StreamReader(new MemoryStream(xml, writable: false), Encoding.UTF8, detectEncodingFromByteOrderMarks: true);
    return reader.ReadToEnd();
}

static Findings Summarize(ValidationResult r)
{
    if (r.ResolvedFormat is null || r.ResolvedFormat.Format is null)
    {
        throw new InvalidOperationException("format not detected: no rule set run");
    }
    // OfficialFatalViolations: the fatal violations of the XSD layer (no
    // SourceId) and of the official rule sets.
    return Findings.Of(r.OfficialFatalViolations.Select(v => string.IsNullOrEmpty(v.SourceId) ? "schema:" + v.RuleId : v.RuleId));
}

return Harness.Run(lib, args, h =>
{
    h.Extract(
        "FacturXExtractor.ExtractBytes: names factur-x.xml, zugferd-invoice.xml, ZUGFeRD-invoice.xml, xrechnung.xml, order-x.xml (then case-insensitive), PDFsharp 6.2.4, raw bytes",
        Extract);
    h.Parse("CiiInvoiceParser.Parse(XDocument) → Invoice + elements not carried over (COMFORT content)", Parse);
    h.Validate(
        $"LayeredValidator.Validate: CII D16B XSD (Antiphon.Xsd.Cii) then rules chosen by FormatDetector, CEN EN 16931 CII Schematron ({CenCii.Version}, {CenCii.AssertCount} assertions) and KoSIT XRechnung CII ({XrCii.Version}, {XrCii.AssertCount} assertions), precompiled, in-house XPath evaluator; errors = OfficialFatalViolations; no Factur-X or French rules",
        x => LayeredValidator.Validate(Text(x), RuleSets()),
        Summarize);
    h.Write(
        "CiiInvoiceGenerator.GenerateXml(Invoice, Comfort), model read by CiiInvoiceParser.Parse",
        Parse,
        r => Encoding.UTF8.GetBytes(CiiInvoiceGenerator.GenerateXml(r.Invoice, CiiProfile.Comfort)));
    h.Embed(
        "FacturXEmbedder.Embed(PDF, XML, EN 16931): factur-x.xml associated file (/AF, /Names/EmbeddedFiles), PDF/A-3 + Factur-X XMP; the host must already be PDF/A (documented contract, not checked by the library)",
        (host, xml) => FacturXEmbedder.Embed(host, xml),
        Extract);
    h.Unsupported(Op.CheckPdf, "no hybrid PDF check (extraction and embedding only)");
});
