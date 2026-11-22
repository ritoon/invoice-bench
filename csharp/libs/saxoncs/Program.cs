// Reference validator: SaxonCS-HE 13.1.0 (Saxonica, MPL-2.0, SaxonJ-HE transpiled
// to C#: no Java, no IKVM) running the official XSLT, as the Python bench does
// with saxonche:
//
// - CEN EN 16931 1.3.17 (ConnectingEurope/eInvoicing-EN16931):
//   xslt/EN16931-CII-validation.xslt for a CII invoice,
//   xslt/EN16931-UBL-validation.xslt for a UBL invoice or credit note;
// - KoSIT XRechnung 3.0.2 schematron 2.6.0 (itplr-kosit/xrechnung-schematron,
//   schematron/cii/XRechnung-CII-validation.xsl or
//   schematron/ubl/XRechnung-UBL-validation.xsl) as well when BT-24 names
//   XRechnung, as the scenarios of validator-configuration-xrechnung do.
//
// Artefacts downloaded by `run.sh setup` into .artefacts/ (SHA-256 checked).
// The stylesheets of a syntax are compiled on its first call (part of
// first_call_ns for CII, the syntax of the bench) and reused. Counting as the
// KoSIT validator does: an error is any assertion with flag fatal or error;
// for XRechnung the custom levels of the scenarios of
// validator-configuration-xrechnung v2026-08-31 apply (see Level). No XSD
// (the CEN artefacts do not ship one). The CII path is the one of the bench;
// the UBL path serves the conformance cases (../cases).

using System.Text.RegularExpressions;
using Bench;
using Saxon.Api;

var lib = new Lib("saxoncs", "SaxonCS-HE", "13.1.0");

return Harness.Run(lib, args, h =>
{
    h.Unsupported(Op.Extract, "validator only");
    h.Unsupported(Op.Parse, "validator only");
    h.Validate(
        $"official CEN EN 16931 {Reference.CenVersion} XSLT, CII or UBL by root (+ KoSIT XRechnung 3.0.2 schematron {Reference.XrVersion} if BT-24 is XRechnung) via SaxonCS-HE, compiled on the first call; fatal/error errors with the levels of the KoSIT v2026-08-31 scenarios",
        Reference.Validate,
        rules => Findings.Of(rules));
    h.Unsupported(Op.Write, "validator only");
    h.Unsupported(Op.Embed, "validator only");
    h.Unsupported(Op.CheckPdf, "validator only");
});

static partial class Reference
{
    public const string CenVersion = "1.3.17";
    public const string XrVersion = "2.6.0";

    /// <summary>
    /// The level of a rule under the scenario of validator-configuration-xrechnung
    /// v2026-08-31 that BT-24 selects (CIUS, extension or CVD), its flag
    /// otherwise: the customLevel elements of scenarios.xml, where they change
    /// whether the rule is an error.
    /// </summary>
    private static string Level(string rule, string flag, string spec, string root)
    {
        if (!spec.Contains("xrechnung", StringComparison.Ordinal))
        {
            return flag;
        }
        bool ext = spec.Contains(":extension:", StringComparison.Ordinal), cvd = spec.Contains(":cvd_", StringComparison.Ordinal);
        return rule switch
        {
            "BR-CL-21" or "BR-CL-23" or "CII-SR-475" or "CII-SR-476" => "warning",
            "CII-SR-452" or "CII-SR-453" or "CII-SR-454" or "CII-SR-465" or "CII-SR-466" => "error",
            "BR-CL-13" => cvd ? "information" : flag,
            "BR-CL-10" or "BR-CL-11" or "BR-CL-24" or "BR-CL-25" or "BR-CL-26" => ext ? "information" : flag,
            // UBL extension scenario only
            "BR-CO-16" or "UBL-CR-470" => ext && root != "CrossIndustryInvoice" ? "information" : flag,
            "UBL-CR-646" => ext ? "information" : cvd || root != "CreditNote" ? "error" : flag,
            _ => flag,
        };
    }

    private sealed record Compiled(Processor Processor, DocumentBuilder Builder, XsltExecutable Cen, XsltExecutable Xr, XPathExecutable Asserts);

    private static readonly Dictionary<string, Compiled> Cache = new(StringComparer.Ordinal);

    private static Compiled Compile(string syntax)
    {
        if (Cache.TryGetValue(syntax, out var hit))
        {
            return hit;
        }
        var art = Corpus.Artefacts;
        var (cen, xr) = syntax == "ubl"
            ? (Path.Combine(art, $"cen-ubl-{CenVersion}", "xslt", "EN16931-UBL-validation.xslt"),
                Path.Combine(art, $"xrechnung-{XrVersion}", "schematron", "ubl", "XRechnung-UBL-validation.xsl"))
            : (Path.Combine(art, $"cen-{CenVersion}", "xslt", "EN16931-CII-validation.xslt"),
                Path.Combine(art, $"xrechnung-{XrVersion}", "schematron", "cii", "XRechnung-CII-validation.xsl"));
        foreach (var p in new[] { cen, xr })
        {
            if (!File.Exists(p))
            {
                throw new FileNotFoundException($"{p} missing: run ./run.sh setup");
            }
        }
        var processor = new Processor(false);
        var compiler = processor.NewXsltCompiler();
        XsltExecutable Load(string path)
        {
            using var s = File.OpenRead(path);
            return compiler.Compile(s, new Uri(path));
        }
        var xpath = processor.NewXPathCompiler();
        xpath.DeclareNamespace("svrl", "http://purl.oclc.org/dsdl/svrl");
        var compiled = new Compiled(
            processor,
            processor.NewDocumentBuilder(),
            Load(cen),
            Load(xr),
            xpath.Compile("//svrl:failed-assert | //svrl:successful-report"));
        Cache[syntax] = compiled;
        return compiled;
    }

    [GeneratedRegex(@"<(?:[\w.-]+:)?(CrossIndustryInvoice|Invoice|CreditNote)[\s>/]")]
    private static partial Regex Root();

    [GeneratedRegex(@"GuidelineSpecifiedDocumentContextParameter>\s*<(?:[\w.-]+:)?ID>([^<]+)<|<(?:[\w.-]+:)?CustomizationID>([^<]+)<")]
    private static partial Regex Spec();

    private static string Head(byte[] xml, int n) => System.Text.Encoding.UTF8.GetString(xml, 0, Math.Min(n, xml.Length));

    private static List<string> Errors(Compiled c, XsltExecutable xslt, XdmNode doc, string spec, string root)
    {
        var dest = new XdmDestination();
        var transformer = xslt.Load30();
        // The KoSIT stylesheet reads the document from global variables.
        transformer.GlobalContextItem = doc;
        transformer.ApplyTemplates(doc, dest);
        var sel = c.Asserts.Load();
        sel.ContextItem = dest.XdmNode;
        var output = new List<string>();
        foreach (var item in sel)
        {
            var el = (XdmNode)item;
            var rule = el.GetAttributeValue("id") ?? "?";
            var flag = (el.GetAttributeValue("flag") ?? "error").ToLowerInvariant();
            if (Level(rule, flag, spec, root) is "fatal" or "error")
            {
                output.Add(rule);
            }
        }
        return output;
    }

    public static List<string> Validate(byte[] xml)
    {
        // The samples open with long licence comments: the root and BT-24 are
        // looked for in the first 64 KiB.
        var head = Head(xml, 65536);
        var rm = Root().Match(head);
        var root = rm.Success ? rm.Groups[1].Value : "CrossIndustryInvoice";
        var c = Compile(root == "CrossIndustryInvoice" ? "cii" : "ubl");
        XdmNode doc;
        using (var ms = new MemoryStream(xml, writable: false))
        {
            doc = c.Builder.Build(ms);
        }
        var sm = Spec().Match(head);
        var spec = sm.Success ? (sm.Groups[1].Success ? sm.Groups[1].Value : sm.Groups[2].Value).Trim() : "";
        var rules = Errors(c, c.Cen, doc, spec, root);
        if (spec.Contains("xrechnung", StringComparison.Ordinal))
        {
            rules.AddRange(Errors(c, c.Xr, doc, spec, root));
        }
        return rules;
    }
}
