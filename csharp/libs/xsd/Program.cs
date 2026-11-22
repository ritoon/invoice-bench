// Floor: XSD validation alone with the BCL (System.Xml.Schema.XmlSchemaSet, a
// validating XmlReader), against the UN/CEFACT CII D16B schemas that the KoSIT
// configuration ships (validator-configuration-xrechnung v2026-08-31,
// resources/cii/16b/xsd). Not a Factur-X library: it shows the cost of the
// schema layer that most libraries put in front of their rules. A UBL invoice
// or credit note (the conformance cases, ../cases) is validated against the
// UBL 2.1 schemas of the same configuration instead (resources/ubl/2.1/xsd),
// compiled on the first UBL document.

using System.Xml;
using System.Xml.Schema;
using Bench;

var lib = new Lib("xsd", "System.Xml.Schema (BCL)", Environment.Version.ToString());

var schemas = new Dictionary<string, XmlSchemaSet>(StringComparer.Ordinal);

XmlSchemaSet Schemas(string syntax)
{
    if (schemas.TryGetValue(syntax, out var hit))
    {
        return hit;
    }
    var res = Path.Combine(Corpus.Artefacts, "kosit-config-2026-08-31", "resources");
    string[] mains = syntax == "ubl"
        ? [Path.Combine(res, "ubl", "2.1", "xsd", "maindoc", "UBL-Invoice-2.1.xsd"), Path.Combine(res, "ubl", "2.1", "xsd", "maindoc", "UBL-CreditNote-2.1.xsd")]
        : [Path.Combine(res, "cii", "16b", "xsd", "CrossIndustryInvoice_100pD16B.xsd")];
    var set = new XmlSchemaSet { XmlResolver = new XmlUrlResolver() };
    if (syntax == "ubl")
    {
        // The XML-DSig schema that UBL imports opens with a DOCTYPE, which the
        // readers of XmlSchemaSet refuse: it is added first, read with DTD
        // processing on, and the import then finds its namespace in the set.
        var dsig = Path.Combine(res, "ubl", "2.1", "xsd", "common", "UBL-xmldsig-core-schema-2.1.xsd");
        using var r = XmlReader.Create(dsig, new XmlReaderSettings { DtdProcessing = DtdProcessing.Parse, XmlResolver = new XmlUrlResolver() });
        set.Add(null, r);
    }
    foreach (var main in mains)
    {
        if (!File.Exists(main))
        {
            throw new FileNotFoundException($"{main} missing: run ./run.sh setup");
        }
        set.Add(null, main);
    }
    set.Compile();
    return schemas[syntax] = set;
}

byte[] ublNs = "urn:oasis:names:specification:ubl:schema:xsd:"u8.ToArray();

List<string> Validate(byte[] xml)
{
    var errors = new List<string>();
    var settings = new XmlReaderSettings
    {
        ValidationType = ValidationType.Schema,
        Schemas = Schemas(xml.AsSpan(0, Math.Min(xml.Length, 65536)).IndexOf(ublNs) >= 0 ? "ubl" : "cii"),
        DtdProcessing = DtdProcessing.Prohibit,
        XmlResolver = null,
    };
    settings.ValidationEventHandler += (_, e) =>
    {
        if (e.Severity == XmlSeverityType.Error)
        {
            errors.Add("XSD: " + (e.Message.Length > 120 ? e.Message[..120] : e.Message));
        }
    };
    using var ms = new MemoryStream(xml, writable: false);
    using var reader = XmlReader.Create(ms, settings);
    while (reader.Read())
    {
    }
    return errors;
}

return Harness.Run(lib, args, h =>
{
    h.Unsupported(Op.Extract, "XSD validation floor only");
    h.Unsupported(Op.Parse, "XSD validation floor only");
    h.Validate(
        "floor: validating XmlReader (XmlSchemaSet compiled once) against the UN/CEFACT CII D16B XSDs (UBL 2.1 for a UBL document) of the KoSIT configuration v2026-08-31; no business rules",
        Validate,
        rules => Findings.Of(rules));
    foreach (var op in new[] { Op.Write, Op.Embed, Op.CheckPdf })
    {
        h.Unsupported(op, "XSD validation floor only");
    }
});
