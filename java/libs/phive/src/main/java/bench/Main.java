package bench;

import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;

import bench.harness.Harness;
import bench.harness.Harness.Findings;
import bench.harness.Harness.Lazy;
import bench.harness.Harness.Lib;
import bench.harness.Harness.Op;
import com.helger.diagnostics.error.IError;
import com.helger.diver.api.coord.DVRCoordinate;
import com.helger.io.resource.inmemory.ReadableResourceByteArray;
import com.helger.phive.api.execute.ValidationExecutionManager;
import com.helger.phive.api.executorset.IValidationExecutorSet;
import com.helger.phive.api.executorset.ValidationExecutorSetRegistry;
import com.helger.phive.api.result.ValidationResult;
import com.helger.phive.api.result.ValidationResultList;
import com.helger.phive.api.validity.IValidityDeterminator;
import com.helger.phive.en16931.EN16931Validation;
import com.helger.phive.xml.source.IValidationSourceXML;
import com.helger.phive.xml.source.ValidationSourceXML;
import com.helger.phive.xrechnung.XRechnungValidation;
import com.helger.phive.zugferd.ZugferdValidation;

/**
 * phive 12.3.0 (phax/phive) with the rule packages phive-rules-en16931,
 * phive-rules-xrechnung and phive-rules-zugferd 4.7.1: the official XSD and
 * pre-compiled schematron XSLT run by ph-schematron on Saxon-HE. Validator
 * only.
 */
public final class Main {

    static final Lib LIB = new Lib("phive", "com.helger.phive:phive-xml+phive-rules-{en16931,xrechnung,zugferd}",
            "12.3.0/4.7.1");

    static final Lazy<ValidationExecutorSetRegistry<IValidationSourceXML>> REGISTRY = Lazy.of(v -> {
        ValidationExecutorSetRegistry<IValidationSourceXML> reg = new ValidationExecutorSetRegistry<>();
        EN16931Validation.initEN16931(reg);
        XRechnungValidation.initXRechnung(reg);
        ZugferdValidation.initZugferd(reg);
        return reg;
    });

    /**
     * The validation executor set (VES) for the syntax and the specification
     * identifier (BT-24): XRechnung 3.0.2 for XRechnung (its extension set
     * for the XRechnung extension), the ZUGFeRD 2.5.2 / Factur-X 1.09.2 set
     * of the profile for the Factur-X profiles, the CEN EN 16931 1.3.17 set
     * otherwise; in CII, or in UBL for a UBL invoice or credit note (the
     * conformance cases).
     */
    static DVRCoordinate vesFor(byte[] xml) {
        String id = Harness.specId(xml);
        boolean ext = id.contains("kosit:extension:xrechnung");
        if (Harness.isUbl(xml)) {
            boolean credit = Harness.rootName(xml).equals("CreditNote");
            if (id.contains("xrechnung")) {
                return ext && !credit ? XRechnungValidation.VID_XRECHNUNG_EXTENSION_UBL_INVOICE_302
                        : credit ? XRechnungValidation.VID_XRECHNUNG_UBL_CREDITNOTE_302
                        : XRechnungValidation.VID_XRECHNUNG_UBL_INVOICE_302;
            }
            return credit ? EN16931Validation.VID_UBL_CREDIT_NOTE_1317 : EN16931Validation.VID_UBL_INVOICE_1317;
        }
        if (id.contains("xrechnung")) {
            return ext ? XRechnungValidation.VID_XRECHNUNG_EXTENSION_CII_302 : XRechnungValidation.VID_XRECHNUNG_CII_302;
        }
        if (id.contains("factur-x.eu:1p0:minimum")) {
            return ZugferdValidation.VID_ZUGFERD_2_5_2_MINIMUM;
        }
        if (id.contains("factur-x.eu:1p0:basicwl")) {
            return ZugferdValidation.VID_ZUGFERD_2_5_2_BASIC_WL;
        }
        if (id.contains("factur-x.eu:1p0:basic")) {
            return ZugferdValidation.VID_ZUGFERD_2_5_2_BASIC;
        }
        if (id.contains("factur-x.eu:1p0:extended")) {
            return ZugferdValidation.VID_ZUGFERD_2_5_2_EXTENDED;
        }
        return EN16931Validation.VID_CII_1317;
    }

    static ValidationResultList validate(byte[] xml) throws Exception {
        IValidationExecutorSet<IValidationSourceXML> ves = REGISTRY.get().getOfID(vesFor(xml));
        if (ves == null) {
            throw new IllegalStateException("no validation executor set for " + vesFor(xml));
        }
        return ValidationExecutionManager.executeValidation(IValidityDeterminator.createDefault(), ves,
                ValidationSourceXML.create(new ReadableResourceByteArray("invoice.xml", xml)));
    }

    /** Errors: every finding of level ERROR or above in every layer (XSD,
     *  schematron), after the set's custom levels; a finding without a rule
     *  id is named after its layer (XSD, or an internal schematron error). */
    static Findings summarize(ValidationResultList r) {
        List<String> rules = new ArrayList<>();
        for (ValidationResult vr : r) {
            for (IError e : vr.getErrorList()) {
                if (!e.isError()) {
                    continue;
                }
                String id = e.getErrorID();
                if (id != null && !id.isEmpty()) {
                    rules.add(id);
                } else {
                    String text = e.getErrorText(Locale.ROOT);
                    String layer = vr.getValidationArtefact().getValidationType().getID();
                    rules.add(layer.contains("xsd") ? "XSD"
                            : layer + ": " + (text == null ? "" : text.substring(0, Math.min(text.length(), 120))));
                }
            }
        }
        return Findings.of(rules);
    }

    public static void main(String[] args) {
        Harness.run(LIB, args, h -> {
            h.unsupported(Op.EXTRACT, "validator only");
            h.unsupported(Op.PARSE, "validator only (the JAXB model is ph-cii's, measured separately)");
            h.validate("ValidationExecutionManager.executeValidation, set chosen by BT-24: XRechnung CII 3.0.2 "
                    + "(D16B XSD + CEN CII 1.3.16 + KoSIT XRechnung 2.6.0, phive's own levels) for XRechnung "
                    + "(the extension set for the XRechnung extension; the UBL sets for a UBL case), "
                    + "the profile's ZUGFeRD 2.5.2 (XSD + Factur-X 1.09.2 schematron) for MINIMUM/BASIC WL/BASIC/EXTENDED, "
                    + "EN 16931 CII 1.3.17 (D16B XSD + CEN schematron) otherwise, UBL 1.3.17 for a UBL case; DOM parsing included; "
                    + "errors = level ERROR and above; registry initialized and XSLT loaded on the first call",
                    Main::validate, Main::summarize);
            h.unsupported(Op.WRITE, "validator only");
            h.unsupported(Op.EMBED, "validator only");
            h.unsupported(Op.CHECK_PDF, "validator only (XML only)");
        });
    }
}
