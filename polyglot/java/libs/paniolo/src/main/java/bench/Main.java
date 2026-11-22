package bench;

import java.util.ArrayList;
import java.util.List;

import bench.harness.Harness;
import bench.harness.Harness.Findings;
import bench.harness.Harness.Lib;
import io.github.nalukit.paniolo.Paniolo;
import io.github.nalukit.paniolo.config.FacturXProfile;
import io.github.nalukit.paniolo.config.InvoiceSyntax;
import io.github.nalukit.paniolo.model.Invoice;
import io.github.nalukit.paniolo.validation.Severity;
import io.github.nalukit.paniolo.validation.ValidationMessage;
import io.github.nalukit.paniolo.validation.ValidationResult;

/**
 * Paniolo 1.0.0 (NaluKit/paniolo): EN 16931 model, UBL / CII read and write,
 * validation with the official artefacts it bundles (KoSIT XRechnung
 * configuration, Factur-X / ZUGFeRD schematrons) on Saxon-HE; paniolo-facturx
 * for the PDF side (PDFBox), paniolo-pdfa for PDF/A conformance (veraPDF).
 */
public final class Main {

    static final Lib LIB = new Lib("paniolo", "io.github.nalukit:paniolo+paniolo-facturx+paniolo-pdfa", "1.0.0");

    /** Thread safe and meant to be shared; rules are compiled lazily by the
     *  first validation (inside first_call_ns). */
    static final Paniolo P = Paniolo.create();

    /** Errors: messages of severity ERROR or FATAL, all stages. */
    static Findings summarize(ValidationResult r) {
        List<String> rules = new ArrayList<>();
        for (ValidationMessage m : r.messages()) {
            if (m.severity() == Severity.ERROR || m.severity() == Severity.FATAL) {
                rules.add(m.ruleId() == null || m.ruleId().isEmpty() ? String.valueOf(m.stage()) : m.ruleId());
            }
        }
        return Findings.of(rules);
    }

    static byte[] embed(byte[] host, byte[] xml) {
        return P.embedXml(host, xml, FacturXProfile.EN16931).pdf();
    }

    public static void main(String[] args) {
        Harness.run(LIB, args, h -> {
            h.extract("Paniolo.extractXml(pdf) (paniolo-facturx)", P::extractXml);
            h.parse("Paniolo.read(byte[]) in STRICT mode → Invoice (EN 16931 model; rejects amount discrepancies)",
                    P::read);
            h.validate("Paniolo.validate(byte[]): XSD + bundled official schematrons chosen by BT-24 "
                    + "(XRechnung 3.0: KoSIT configuration validator-configuration-xrechnung v2026-08-31, i.e. CEN CII "
                    + "1.3.16 + KoSIT 2.6.0 and its levels; plain EN 16931: the same rules, findings of the XRechnung CIUS "
                    + "(BR-DE, BR-TMP, PEPPOL) downgraded to information; Factur-X profiles: Factur-X 1.09.2 schematrons) "
                    + "+ amount check; "
                    + "errors = ERROR/FATAL; rules compiled on the first call", P::validate, Main::summarize);
            h.write("Paniolo.toBytes(Invoice, InvoiceSyntax.CII), model read by read", P::read,
                    (Invoice inv) -> P.toBytes(inv, InvoiceSyntax.CII));
            h.embed("Paniolo.embedXml(pdf, xml, FacturXProfile.EN16931).pdf(); checks BT-24, embeds then revalidates the "
                    + "resulting PDF (container, PDF/A by veraPDF, XML), revalidation included in the measurement; does not "
                    + "convert to PDF/A; rejects a host that already contains an invoice: host.pdf and host_classic.pdf keep "
                    + "in /AF a factur-x.xml file specification with its /UF stream (only the EmbeddedFiles tree was removed), "
                    + "which Paniolo reads (PdfBoxFacturXSupport.allFiles)", Main::embed, P::extractXml);
            h.checkPdf("Paniolo.validate on the PDF: Factur-X container (attachment, relationship, XMP), PDF/A by "
                    + "veraPDF (paniolo-pdfa), then full validation of the embedded XML for the declared profile",
                    P::validate);
        });
    }
}
