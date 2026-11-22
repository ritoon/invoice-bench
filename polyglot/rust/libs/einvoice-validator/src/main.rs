//! einvoice-validator 0.4.0 (with einvoice-pdf and einvoice-xsd 0.4.0): the
//! official CEN validation XSLT executed by a pure-Rust XSLT 2.0 subset, after
//! XML Schema validation; PDF extraction.

use std::sync::OnceLock;

use einvoice_validator::{Report, Severity, Validator, pdf};
use harness::{Findings, Lib, Op};

const LIB: Lib = Lib { lib: "einvoice-validator", pkg: "einvoice-validator", version: "0.4.0" };

fn extract(p: &[u8]) -> Result<Vec<u8>, String> {
    pdf::extract(p).map(|e| e.xml.data).map_err(|e| e.to_string())
}

/// Built once per process, on the first call (schemas loaded, stylesheet
/// prepared): that cost is in `first_call_ns` of ParseValidate_Small.
fn validator() -> Result<&'static Validator, String> {
    static V: OnceLock<Result<Validator, String>> = OnceLock::new();
    V.get_or_init(|| Validator::cen_cii().map_err(|e| e.to_string())).as_ref().map_err(Clone::clone)
}

fn validate(x: &[u8]) -> Result<Report, String> {
    validator()?.validate(x).map_err(|e| e.to_string())
}

/// Errors: schema errors, fatal findings, and a missing scenario.
fn summarize(r: &Report) -> Findings {
    let mut rules: Vec<String> = r.schema_errors.iter().map(|_| "XSD".to_owned()).collect();
    rules.extend(
        r.findings()
            .filter(|f| f.severity == Severity::Fatal)
            .map(|f| f.id.clone().unwrap_or_else(|| f.test.clone())),
    );
    if r.no_scenario.is_some() {
        rules.push("NO-SCENARIO".into());
    }
    Findings { errors: rules.len(), rules }
}

fn main() {
    harness::run(LIB, |h| {
        h.extract("pdf::extract (einvoice-pdf, default Limits)", extract);
        h.unsupported(Op::Parse, "no invoice model: the library validates the XML without reading it into a model");
        h.validate(
            "Validator::cen_cii(): UN/CEFACT CII D16B XSD + official CEN EN 16931 CII 1.3.17 XSLT executed in Rust; the only rule set bundled offline (XRechnung, Factur-X and Peppol need the `fetch` feature and a download)",
            validate,
            summarize,
        );
        h.unsupported(Op::Write, "no XML generation");
        h.unsupported(Op::Embed, "no embedding in a PDF");
        h.unsupported(Op::CheckPdf, "no native check of a hybrid PDF (extraction then validation, to be combined by the caller)");
    });
}
