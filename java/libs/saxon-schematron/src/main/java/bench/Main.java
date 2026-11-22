package bench;

import java.io.ByteArrayInputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;

import javax.xml.transform.stream.StreamSource;

import bench.harness.Harness;
import bench.harness.Harness.Findings;
import bench.harness.Harness.Lazy;
import bench.harness.Harness.Lib;
import bench.harness.Harness.Op;
import net.sf.saxon.s9api.DocumentBuilder;
import net.sf.saxon.s9api.Processor;
import net.sf.saxon.s9api.QName;
import net.sf.saxon.s9api.XPathCompiler;
import net.sf.saxon.s9api.XPathExecutable;
import net.sf.saxon.s9api.XPathSelector;
import net.sf.saxon.s9api.XdmDestination;
import net.sf.saxon.s9api.XdmItem;
import net.sf.saxon.s9api.XdmNode;
import net.sf.saxon.s9api.XsltExecutable;
import net.sf.saxon.s9api.XsltTransformer;

/**
 * Reference floor: Saxon-HE 13.1 (s9api) running the official validation
 * stylesheets, the Java port of python/adapters/saxonche_schematron.py.
 *
 * <ul>
 * <li>CEN EN 16931 1.3.17 (ConnectingEurope/eInvoicing-EN16931):
 * xslt/EN16931-CII-validation.xslt for a CII invoice,
 * xslt/EN16931-UBL-validation.xslt for a UBL invoice or credit note;</li>
 * <li>KoSIT XRechnung 3.0.2 schematron 2.6.0 (itplr-kosit/xrechnung-schematron,
 * schematron/cii/XRechnung-CII-validation.xsl or
 * schematron/ubl/XRechnung-UBL-validation.xsl) as well when BT-24 names
 * XRechnung, as the scenarios of validator-configuration-xrechnung.</li>
 * </ul>
 *
 * <p>An error is any failed assertion or successful report of flag fatal or
 * error. For XRechnung, the custom levels of the scenarios of
 * validator-configuration-xrechnung v2026-08-31 apply (see {@link #level}).
 * No XSD (the CEN artefacts do not ship one). The CII path is the one of the
 * bench; the UBL path serves the conformance cases (../cases).
 */
public final class Main {

    static final Lib LIB = new Lib("saxon-schematron", "net.sf.saxon:Saxon-HE", "13.1");

    static final String CEN_VERSION = "1.3.17";
    static final String XR_VERSION = "2.6.0";

    static final QName ID = new QName("id");
    static final QName FLAG = new QName("flag");

    record Compiled(Processor proc, DocumentBuilder builder, XsltExecutable cen, XsltExecutable xr,
                    XPathExecutable findings) {
    }

    static Lazy<Compiled> compiled(String cen, String xr) {
        return Lazy.of(v -> {
            Path c = Harness.artefacts().resolve(cen);
            Path x = Harness.artefacts().resolve(xr);
            for (Path p : List.of(c, x)) {
                if (!Files.isRegularFile(p)) {
                    throw new IllegalStateException(p + " missing: run ./run.sh setup");
                }
            }
            Processor proc = new Processor(false);
            var comp = proc.newXsltCompiler();
            XPathCompiler xp = proc.newXPathCompiler();
            xp.declareNamespace("svrl", "http://purl.oclc.org/dsdl/svrl");
            return new Compiled(proc, proc.newDocumentBuilder(),
                    comp.compile(new StreamSource(c.toFile())), comp.compile(new StreamSource(x.toFile())),
                    xp.compile("//svrl:failed-assert | //svrl:successful-report"));
        });
    }

    /** The CII stylesheets, compiled once on first use (inside first_call_ns). */
    static final Lazy<Compiled> CII = compiled("cen-" + CEN_VERSION + "/xslt/EN16931-CII-validation.xslt",
            "xrechnung-" + XR_VERSION + "/schematron/cii/XRechnung-CII-validation.xsl");
    /** The UBL stylesheets, compiled on the first UBL document. */
    static final Lazy<Compiled> UBL = compiled("cen-ubl-" + CEN_VERSION + "/xslt/EN16931-UBL-validation.xslt",
            "xrechnung-" + XR_VERSION + "/schematron/ubl/XRechnung-UBL-validation.xsl");

    /**
     * The level of a rule under the scenario of validator-configuration-xrechnung
     * v2026-08-31 that BT-24 selects (CIUS, extension or CVD), its flag
     * otherwise: the customLevel elements of scenarios.xml, where they change
     * whether the rule is an error.
     */
    static String level(String rule, String flag, String spec, String root) {
        if (!spec.contains("xrechnung")) {
            return flag;
        }
        boolean ext = spec.contains(":extension:"), cvd = spec.contains(":cvd_");
        boolean ubl = !root.equals("CrossIndustryInvoice"), creditNote = root.equals("CreditNote");
        switch (rule) {
            case "BR-CL-21", "BR-CL-23", "CII-SR-475", "CII-SR-476":
                return "warning";
            case "CII-SR-452", "CII-SR-453", "CII-SR-454", "CII-SR-465", "CII-SR-466":
                return "error";
            case "BR-CL-13":
                return cvd ? "information" : flag;
            case "BR-CL-10", "BR-CL-11", "BR-CL-24", "BR-CL-25", "BR-CL-26":
                return ext ? "information" : flag;
            case "BR-CO-16", "UBL-CR-470": // UBL extension scenario only
                return ext && ubl ? "information" : flag;
            case "UBL-CR-646":
                return ext ? "information" : cvd || !creditNote ? "error" : flag;
            default:
                return flag;
        }
    }

    static List<String> errors(Compiled c, XsltExecutable x, XdmNode doc, String spec, String root)
            throws Exception {
        XsltTransformer t = x.load();
        t.setInitialContextNode(doc);
        XdmDestination d = new XdmDestination();
        t.setDestination(d);
        t.transform();
        XPathSelector sel = c.findings().load();
        sel.setContextItem(d.getXdmNode());
        List<String> out = new ArrayList<>();
        for (XdmItem it : sel) {
            XdmNode n = (XdmNode) it;
            String rule = n.getAttributeValue(ID) == null ? "?" : n.getAttributeValue(ID);
            String flag = n.getAttributeValue(FLAG) == null ? "error" : n.getAttributeValue(FLAG).toLowerCase();
            flag = level(rule, flag, spec, root);
            if (flag.equals("fatal") || flag.equals("error")) {
                out.add(rule);
            }
        }
        return out;
    }

    static Findings validate(byte[] xml) throws Exception {
        String root = Harness.rootName(xml);
        Compiled c = Harness.isUbl(xml) ? UBL.get() : CII.get();
        XdmNode doc = c.builder().build(new StreamSource(new ByteArrayInputStream(xml)));
        String spec = Harness.specId(xml);
        List<String> rules = errors(c, c.cen(), doc, spec, root);
        if (spec.contains("xrechnung")) {
            rules.addAll(errors(c, c.xr(), doc, spec, root));
        }
        return Findings.of(rules);
    }

    public static void main(String[] args) {
        Harness.run(LIB, args, h -> {
            h.unsupported(Op.EXTRACT, "validator only");
            h.unsupported(Op.PARSE, "validator only");
            h.validate("official CEN EN 16931 " + CEN_VERSION + " XSLT, CII or UBL by root (+ KoSIT XRechnung 3.0.2 "
                    + "schematron " + XR_VERSION + " if BT-24 is XRechnung) via Saxon-HE 13.1 (s9api), compiled on the "
                    + "first call; fatal/error errors with the custom levels of the KoSIT v2026-08-31 scenarios; no XSD",
                    Main::validate, f -> f);
            h.unsupported(Op.WRITE, "validator only");
            h.unsupported(Op.EMBED, "validator only");
            h.unsupported(Op.CHECK_PDF, "validator only");
        });
    }
}
