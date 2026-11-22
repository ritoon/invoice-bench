// Securibox.FacturX 0.6.0-alpha-3 (Securibox, MIT): Factur-X import, export and
// validation. XmlSerializer models generated from the Factur-X 1.09 XSD (one
// per profile), the FeRD Factur-X 1.09 XSD and schematrons embedded, and its
// own Schematron engine (the .sch deserialised, asserts evaluated with the
// XPath2 1.1.5 library, no XSLT). PDF on PDFsharp 7.0.0-preview-1 (MIT) and
// XmpCore 6.1.10.1.
//
// No public call returns the raw XML or turns a model into XML without a PDF:
// Extract goes through ImportDataWithDeserialization (validation included) and
// Write replicates the serialisation FacturxExporter does internally.

using System.Text;
using System.Xml;
using System.Xml.Serialization;
using Bench;
using Securibox.FacturX;
using Securibox.FacturX.Models.Enums;
using Securibox.FacturX.SpecificationModels;
using SchValidation = Securibox.FacturX.Schematron.Helpers.ValidationResult;

var lib = new Lib("securibox", "Securibox.FacturX", "0.6.0-alpha-3");

return Harness.Run(lib, args, h =>
{
    h.Extract(
        "new FacturxImporter(stream).ImportDataWithDeserialization(): PDFsharp, XMP, XML from /AF or the name tree, then XSD + schematron validation (result ignored) and deserialisation; no public API that returns the bytes",
        Adapter.Extract,
        Adapter.Serialize);
    h.Parse("XmlSerializer on SpecificationModels.<profile>.CrossIndustryInvoice (the library's classes, profile detected by CrossIndustryInvoiceProfileProbe as FacturxImporter does)", Adapter.Parse);
    h.Validate(
        "FacturxXsdValidator.ValidateXml + FacturxSchematronValidator.ValidateXml at the detected profile: embedded FeRD Factur-X 1.09 XSD and schematrons (XRechnung and EN 16931 → EN 16931 profile), in-house Schematron engine on XPath2; the .sch is read again on every call",
        Adapter.Validate,
        r => Findings.Of(r.Sch._results.Where(x => x.IsError).Select(Adapter.RuleOf)
            .Concat(r.Xsd.Select(_ => "xsd"))));
    h.Write(
        "adapter: XmlSerializer.Serialize of the model with the namespaces of FacturxExporter.CreateFacturXStream (no public model → XML API without a PDF); model read by XmlSerializer",
        Adapter.Parse,
        Adapter.Serialize);
    h.Embed(
        "FacturxExporter.CreateFacturXStream(PDF, XML, EN 16931): XSD + schematron validation first (fails on error), then copy of the pages into a new PDFsharp document, factur-x.xml as associated file, XMP by XmpCore",
        Adapter.Embed,
        pdf => Adapter.Serialize(Adapter.Extract(pdf)));
    h.CheckPdf(
        "new FacturxImporter(stream).IsFacturXValid(): Factur-X XMP, embedded XML, Factur-X 1.09 XSD + schematron; no PDF/A check",
        Adapter.Check);
});

static partial class Adapter
{
    // The FeRD asserts carry no @id: the rule is the "[BR-CL-25]" that opens the message.
    [System.Text.RegularExpressions.GeneratedRegex(@"\[([A-Z]+(?:-[A-Z]+)*-\d+[a-z]?)\]")]
    private static partial System.Text.RegularExpressions.Regex RuleId();

    [System.Text.RegularExpressions.GeneratedRegex(@"cl\[@id=(\d+)\]")]
    private static partial System.Text.RegularExpressions.Regex CodeList();

    public static string RuleOf(global::Securibox.FacturX.Schematron.Helpers.ValidationReport r)
    {
        if (r.BusinessRuleCode is { Length: > 0 } id)
        {
            return id;
        }
        var m = RuleId().Match(r.Description ?? "");
        if (m.Success)
        {
            return m.Groups[1].Value;
        }
        // Code list asserts: "… document('…_codedb.xml')/codedb/cl[@id=17]/enumeration…", no rule id.
        var cl = CodeList().Match(r.Test ?? "");
        if (cl.Success)
        {
            return $"codedb cl{cl.Groups[1].Value}";
        }
        return r.Test is { } t && t.Length > 60 ? t[..60] : r.Test ?? "?";
    }

    // FacturxImporter.KnownGuidelineIds (private) and the fallback that follows it.
    private static readonly Dictionary<string, FacturXConformanceLevelType> KnownGuidelineIds = new(StringComparer.OrdinalIgnoreCase)
    {
        ["urn:factur-x.eu:1p0:minimum"] = FacturXConformanceLevelType.Minimum,
        ["urn:factur-x.eu:1p0:basicwl"] = FacturXConformanceLevelType.BasicWL,
        ["urn:cen.eu:en16931:2017#compliant#urn:factur-x.eu:1p0:basic"] = FacturXConformanceLevelType.Basic,
        ["urn:cen.eu:en16931:2017#compliant#urn:factur-x.eu:1p0:en16931"] = FacturXConformanceLevelType.EN16931,
        ["urn:cen.eu:en16931:2017#conformant#urn:factur-x.eu:1p0:extended"] = FacturXConformanceLevelType.Extended,
        ["urn:cen.eu:en16931:2017"] = FacturXConformanceLevelType.EN16931,
        ["urn:cen.eu:en16931:2017#compliant#urn:xeinkauf.de:kosit:xrechnung_3.0"] = FacturXConformanceLevelType.EN16931,
        ["urn:cen.eu:en16931:2017#compliant#urn:xoev-de:kosit:standard:xrechnung_3.0"] = FacturXConformanceLevelType.EN16931,
        ["urn:cen.eu:en16931:2017#compliant#urn:fdc:peppol.eu:2017:poacc:billing:3.0"] = FacturXConformanceLevelType.EN16931,
    };

    public static FacturXConformanceLevelType Level(byte[] xml)
    {
        using var ms = new MemoryStream(xml, writable: false);
        var probe = (CrossIndustryInvoiceProfileProbe?)new XmlSerializer(typeof(CrossIndustryInvoiceProfileProbe)).Deserialize(ms);
        var id = probe?.ExchangedDocumentContext?.GuidelineSpecifiedDocumentContextParameter?.ID?.Value?.Trim();
        if (string.IsNullOrEmpty(id))
        {
            throw new InvalidOperationException("GuidelineSpecifiedDocumentContextParameter/ID not found.");
        }
        if (KnownGuidelineIds.TryGetValue(id, out var level))
        {
            return level;
        }
        var token = id.Split('#')[^1].Split(':')[^1].ToLowerInvariant();
        return token switch
        {
            _ when token.StartsWith("minimum", StringComparison.Ordinal) => FacturXConformanceLevelType.Minimum,
            _ when token.StartsWith("basicwl", StringComparison.Ordinal) => FacturXConformanceLevelType.BasicWL,
            _ when token.StartsWith("basic", StringComparison.Ordinal) => FacturXConformanceLevelType.Basic,
            _ when token.StartsWith("en16931", StringComparison.Ordinal) => FacturXConformanceLevelType.EN16931,
            _ when token.StartsWith("extended", StringComparison.Ordinal) => FacturXConformanceLevelType.Extended,
            _ when id.Contains("en16931", StringComparison.OrdinalIgnoreCase) => FacturXConformanceLevelType.EN16931,
            _ => throw new NotSupportedException($"Unknown Factur-X guideline ID: {id}"),
        };
    }

    private static Type ModelOf(FacturXConformanceLevelType level) => level.Id switch
    {
        1 => typeof(global::Securibox.FacturX.SpecificationModels.Minimum.CrossIndustryInvoice),
        2 => typeof(global::Securibox.FacturX.SpecificationModels.Basic.CrossIndustryInvoice),
        3 => typeof(global::Securibox.FacturX.SpecificationModels.BasicWL.CrossIndustryInvoice),
        4 => typeof(global::Securibox.FacturX.SpecificationModels.EN16931.CrossIndustryInvoice),
        _ => typeof(global::Securibox.FacturX.SpecificationModels.Extended.CrossIndustryInvoice),
    };

    public static object Parse(byte[] xml)
    {
        var type = ModelOf(Level(xml));
        using var ms = new MemoryStream(xml, writable: false);
        return new XmlSerializer(type).Deserialize(ms) ?? throw new InvalidDataException("XmlSerializer returned null");
    }

    // The namespaces and settings of FacturxExporter.CreateFacturXStream(Stream, ICrossIndustryInvoice).
    public static byte[] Serialize(object model)
    {
        var namespaces = new XmlSerializerNamespaces();
        namespaces.Add("qdt", "urn:un:unece:uncefact:data:standard:QualifiedDataType:100");
        namespaces.Add("ram", "urn:un:unece:uncefact:data:standard:ReusableAggregateBusinessInformationEntity:100");
        namespaces.Add("rsm", "urn:un:unece:uncefact:data:standard:CrossIndustryInvoice:100");
        namespaces.Add("udt", "urn:un:unece:uncefact:data:standard:UnqualifiedDataType:100");
        namespaces.Add("xsi", "http://www.w3.org/2001/XMLSchema-instance");
        namespaces.Add("", "");
        var settings = new XmlWriterSettings { Encoding = new UTF8Encoding(), Indent = true };
        using var ms = new MemoryStream();
        using (var w = XmlWriter.Create(ms, settings))
        {
            new XmlSerializer(model.GetType()).Serialize(w, model, namespaces);
        }
        return ms.ToArray();
    }

    public static (List<string> Xsd, SchValidation Sch) Validate(byte[] xml)
    {
        var level = Level(xml);
        var xsd = new List<string>();
        using (var ms = new MemoryStream(xml, writable: false))
        {
            FacturxXsdValidator.ValidateXml(ms, level, xsd);
        }
        using var sch = new MemoryStream(xml, writable: false);
        return (xsd, FacturxSchematronValidator.ValidateXml(sch, level));
    }

    public static object Extract(byte[] pdf)
    {
        using var ms = new MemoryStream(pdf, writable: false);
        using var importer = new FacturxImporter(ms);
        return importer.ImportDataWithDeserialization() ?? throw new InvalidDataException("ImportDataWithDeserialization returned null");
    }

    public static byte[] Embed(byte[] host, byte[] xml)
    {
        using var pdf = new MemoryStream(host, writable: false);
        using var x = new MemoryStream(xml, writable: false);
        using var output = new FacturxExporter().CreateFacturXStream(pdf, x, FacturXConformanceLevelType.EN16931);
        using var copy = new MemoryStream();
        output.Position = 0;
        output.CopyTo(copy);
        return copy.ToArray();
    }

    public static bool Check(byte[] pdf)
    {
        using var ms = new MemoryStream(pdf, writable: false);
        using var importer = new FacturxImporter(ms);
        return importer.IsFacturXValid();
    }
}
