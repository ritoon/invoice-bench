package bench;

import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.IdentityHashMap;
import java.util.List;
import java.util.Map;

import bench.harness.Harness;
import bench.harness.Harness.Findings;
import bench.harness.Harness.Lazy;
import bench.harness.Harness.Lib;
import bench.harness.Harness.Op;
import de.kosit.validationtool.api.Check;
import de.kosit.validationtool.api.Configuration;
import de.kosit.validationtool.api.InputFactory;
import de.kosit.validationtool.api.Result;
import de.kosit.validationtool.api.XmlError;
import de.kosit.validationtool.impl.DefaultCheck;
import de.kosit.validationtool.impl.model.CustomFailedAssert;
import de.kosit.validationtool.impl.xml.ProcessorProvider;
import de.kosit.validationtool.model.scenarios.ErrorLevelType;
import net.sf.saxon.s9api.Processor;
import org.oclc.purl.dsdl.svrl.FailedAssert;

/**
 * KoSIT validator 1.6.3 (itplr-kosit/validator), the reference validator of
 * the German administration, driven by its official configuration
 * validator-configuration-xrechnung v2026-08-31 (XRechnung 3.0.2: XSD CII
 * D16B, CEN EN 16931 CII 1.3.16, KoSIT XRechnung schematron 2.6.0, custom
 * levels of the scenarios). Validator only.
 */
public final class Main {

    static final Lib LIB = new Lib("kosit-validator", "org.kosit:validator", "1.6.3");

    /** Unpacked by run.sh setup from the release zip (SHA-256 checked). */
    static final String CONFIG = "kosit-xrechnung-3.0.2-2026-08-31";

    /** Configuration loaded and stylesheets compiled once, on first use
     *  (inside first_call_ns), as a server using the validator would. */
    static final Lazy<Check> CHECK = Lazy.of(v -> {
        Path dir = Harness.artefacts().resolve(CONFIG);
        Path scenarios = dir.resolve("scenarios.xml");
        if (!Files.isRegularFile(scenarios)) {
            throw new IllegalStateException(scenarios + " missing: run ./run.sh setup");
        }
        Processor p = ProcessorProvider.getProcessor();
        Configuration cfg = Configuration.load(scenarios.toUri(), dir.toUri()).build(p);
        return new DefaultCheck(p, cfg);
    });

    static Result validate(byte[] xml) throws Exception {
        return CHECK.get().checkInput(InputFactory.read(xml, "invoice.xml"));
    }

    /**
     * Errors as the KoSIT report assesses them: schema violations of level
     * error or fatal, plus failed assertions whose effective level is error
     * (the scenario's custom level when it has one, the assertion's flag
     * otherwise), plus processing errors; a reject with none of them (no
     * matching scenario…) counts as one finding.
     */
    static Findings summarize(Result r) {
        List<String> rules = new ArrayList<>();
        for (XmlError e : r.getSchemaViolations() == null ? List.<XmlError>of() : r.getSchemaViolations()) {
            if (e.getSeverity() != XmlError.Severity.SEVERITY_WARNING) {
                rules.add("XSD");
            }
        }
        Map<FailedAssert, ErrorLevelType> custom = new IdentityHashMap<>();
        if (r.getCustomFailedAsserts() != null) {
            for (CustomFailedAssert c : r.getCustomFailedAsserts()) {
                custom.put(c.getFailedAssert(), c.getCustomLevelFlag());
            }
        }
        for (FailedAssert fa : r.getFailedAsserts()) {
            ErrorLevelType lvl = custom.get(fa);
            String flag = lvl != null ? lvl.value() : fa.getFlag() == null ? "error" : fa.getFlag();
            if (flag.equalsIgnoreCase("error") || flag.equalsIgnoreCase("fatal")) {
                rules.add(fa.getId() == null ? "?" : fa.getId());
            }
        }
        for (String e : r.getProcessingErrors() == null ? List.<String>of() : r.getProcessingErrors()) {
            rules.add("PROCESSING: " + e);
        }
        if (rules.isEmpty() && !r.isAcceptable()) {
            rules.add("REJECT: " + r.getAcceptRecommendation());
        }
        return Findings.of(rules);
    }

    public static void main(String[] args) {
        Harness.run(LIB, args, h -> {
            h.unsupported(Op.EXTRACT, "validator only");
            h.unsupported(Op.PARSE, "validator only");
            h.validate("DefaultCheck.checkInput with the official configuration validator-configuration-xrechnung v2026-08-31 "
                    + "(scenario chosen by BT-24: \"EN16931 (CII)\" = CII D16B XSD + CEN EN 16931 CII 1.3.16; "
                    + "XRechnung 3.0 = same + KoSIT XRechnung schematron 2.6.0, custom levels of the scenario); "
                    + "errors = XSD violations + assertions of effective level error/fatal; configuration loaded on the first call",
                    Main::validate, Main::summarize);
            h.unsupported(Op.WRITE, "validator only");
            h.unsupported(Op.EMBED, "validator only");
            h.unsupported(Op.CHECK_PDF, "validator only (XML only, no PDF parsing)");
        });
    }
}
