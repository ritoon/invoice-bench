package bench;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.util.ArrayList;
import java.util.List;

import bench.harness.Harness;
import bench.harness.Harness.Findings;
import bench.harness.Harness.Lib;
import org.mustangproject.Invoice;
import org.mustangproject.ZUGFeRD.Profiles;
import org.mustangproject.ZUGFeRD.ZUGFeRD2PullProvider;
import org.mustangproject.ZUGFeRD.ZUGFeRDExporterFromPDFA;
import org.mustangproject.ZUGFeRD.ZUGFeRDImporter;
import org.mustangproject.ZUGFeRD.ZUGFeRDInvoiceImporter;
import org.mustangproject.validator.ESeverity;
import org.mustangproject.validator.ValidationResultItem;
import org.mustangproject.validator.ZUGFeRDValidator;

/**
 * Mustang 2.26.0 (ZUGFeRD/mustangproject): org.mustangproject:library (read,
 * write, embed on PDFBox) and org.mustangproject:validator (XSD + Factur-X /
 * XRechnung schematrons through ph-schematron and Saxon-HE, PDF/A through
 * veraPDF), in one JVM.
 */
public final class Main {

    static final Lib LIB = new Lib("mustang", "org.mustangproject:library+validator", "2.26.0");

    /** The report string and a snapshot of the context's items (the context
     *  is reused by the next call). */
    record Report(String xml, List<ValidationResultItem> items) {
    }

    static byte[] extract(byte[] pdf) {
        return new ZUGFeRDImporter(new ByteArrayInputStream(pdf)).getRawXML();
    }

    /** One pass: setRawXML(…, false) builds the DOM, extractInvoice fills the
     *  model. fromXML + extractInvoice (the documented route) would build the
     *  model twice. */
    static Invoice parse(byte[] xml) throws Exception {
        ZUGFeRDInvoiceImporter zii = new ZUGFeRDInvoiceImporter();
        zii.setRawXML(xml, false);
        return zii.extractInvoice();
    }

    static Report validate(byte[] doc, String name) {
        ZUGFeRDValidator zv = new ZUGFeRDValidator();
        String xml = zv.validate(doc, name);
        return new Report(xml, List.copyOf(zv.getContext().getResults()));
    }

    /** Errors: items of severity error, fatal or exception (notices and
     *  warnings do not invalidate, as in ValidationContext.isValid). */
    static Findings summarize(Report r) {
        List<String> rules = new ArrayList<>();
        for (ValidationResultItem i : r.items()) {
            ESeverity s = i.getSeverity();
            if (s == ESeverity.error || s == ESeverity.fatal || s == ESeverity.exception) {
                String id = i.getID();
                rules.add(id != null && !id.isEmpty() ? id : "section " + i.getSection());
            }
        }
        return Findings.of(rules);
    }

    static byte[] write(Invoice inv) {
        ZUGFeRD2PullProvider p = new ZUGFeRD2PullProvider();
        p.setProfile(Profiles.getByName("EN16931"));
        p.generateXML(inv);
        return p.getXML();
    }

    static byte[] embed(byte[] host, byte[] xml) throws Exception {
        ZUGFeRDExporterFromPDFA e = new ZUGFeRDExporterFromPDFA();
        e.load(host);
        e.setXML(xml);
        ByteArrayOutputStream out = new ByteArrayOutputStream(host.length + xml.length + 65536);
        e.export(out);
        return out.toByteArray();
    }

    public static void main(String[] args) {
        Harness.run(LIB, args, h -> {
            h.extract("new ZUGFeRDImporter(InputStream).getRawXML(); the import also parses the XML "
                    + "into an Invoice model (setRawXML with automatic parsing), included in the measurement", Main::extract);
            h.parse("ZUGFeRDInvoiceImporter.setRawXML(xml, false) + extractInvoice → Invoice "
                    + "(rejects an invoice whose totals it cannot reproduce)", Main::parse);
            h.validate("ZUGFeRDValidator.validate(byte[], name): XSD + the profile's Factur-X 1.09.2 / ZUGFeRD 2.5.2 schematron "
                    + "(EN 16931 for the EN16931 samples), XRechnung 3.0.2 (KoSIT schematron 2.4.0 + CEN CII 1.3.16) if BT-24 "
                    + "is XRechnung, arithmetic check; BR-FR rules only if seller and buyer are French; "
                    + "errors = error/fatal/exception items; XSLT compiled on the first call",
                    x -> validate(x, "invoice.xml"), Main::summarize);
            h.write("ZUGFeRD2PullProvider.generateXML (EN16931 profile), model read by setRawXML + extractInvoice",
                    Main::parse, Main::write);
            h.embed("ZUGFeRDExporterFromPDFA.load + setXML + export (default EN16931 profile); PDF/A-3B hosts, "
                    + "hence ZUGFeRDExporterFromA3; a non-PDF/A host is rejected (PDF/A-1 converted to A-3 at the "
                    + "metadata level); saved without compression; the resulting PDF (PDF/A-3U) keeps the host's "
                    + "leftover factur-x.xml file specification in /AF next to the new one: the Mustang validator "
                    + "accepts it, Paniolo rejects it (FX-MULTIPLE-INVOICES)", Main::embed, Main::extract);
            h.checkPdf("ZUGFeRDValidator.validate on the PDF: PDF/A-3 by veraPDF, Factur-X XMP and attachment, "
                    + "then full validation of the extracted XML (as ParseValidate)", x -> validate(x, "invoice.pdf"));
        });
    }
}
