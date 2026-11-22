//! erechnung-validator 0.1.1 (Promaris): the KoSIT and Mustang validation
//! artefacts (Schematron and compiled XSLT, XSD) executed by the pure-Rust
//! engine oxixml, embedded in the crate.

use erechnung_validator::{Bericht, Fehler, Stufe, Szenarien, Validator};
use harness::{Findings, Lib, Op};

const LIB: Lib = Lib { lib: "erechnung-validator", pkg: "erechnung-validator", version: "0.1.1" };

type Both = (Bericht, Bericht);

fn text(x: &[u8]) -> Result<&str, String> {
    std::str::from_utf8(x).map_err(|e| format!("not UTF-8: {e}"))
}

fn both(r: (Result<Bericht, Fehler>, Result<Bericht, Fehler>)) -> Result<Both, String> {
    let k = r.0.map_err(|e| format!("KoSIT: {e}"))?;
    let m = r.1.map_err(|e| format!("Mustang: {e}"))?;
    Ok((k, m))
}

fn errors_of(b: &Bericht, tool: &str, rules: &mut Vec<String>) {
    rules.extend(b.xsd_fehler.iter().map(|_| format!("{tool}:XSD")));
    if !b.xsd_gueltig && b.xsd_fehler.is_empty() {
        rules.push(format!("{tool}:XSD"));
    }
    rules.extend(b.befunde.iter().filter(|f| f.stufe == Stufe::Fehler).map(|f| format!("{tool}:{}", f.regel)));
}

/// Errors of both tools: XSD errors and findings of level `Fehler`.
fn summarize(r: &Both) -> Findings {
    let mut rules = Vec::new();
    errors_of(&r.0, "KoSIT", &mut rules);
    errors_of(&r.1, "Mustang", &mut rules);
    Findings { errors: rules.len(), rules }
}

fn main() {
    harness::run(LIB, |h| {
        // One validator per thread; the artefacts are compiled on first use,
        // so that cost lands in first_call_ns of ParseValidate_Small.
        let v = Validator::new();
        h.unsupported(Op::Extract, "internal extraction (private `pdf` module): only through kosit_pdf / beide_pdf");
        h.unsupported(Op::Parse, "no invoice model: the library validates the XML without reading it into a model");
        h.validate(
            "Validator::beide(Szenarien::Automatisch): KoSIT (CII D16B XSD, CEN EN 16931 1.3.16 as Schematron, XRechnung 3.0.2 / schematron 2.6.0 for the XRECHNUNG files, KoSIT configuration, CHANGELOG as of 2026-01-31) + Mustang 2.26.0 (Factur-X XSD and Schematron of the declared profile, XR 3.0 for the XRECHNUNG files); errors of both tools added up",
            |x| both(v.beide(text(x)?, Szenarien::Automatisch).map_err(|e| e.to_string())?),
            summarize,
        );
        h.unsupported(Op::Write, "no XML generation");
        h.unsupported(Op::Embed, "no embedding in a PDF");
        h.check_pdf(
            "Validator::beide_pdf: extraction of the XML then KoSIT + Mustang; no check of the PDF/A-3 container (PDF part of Mustang/veraPDF not carried over)",
            |p| v.beide_pdf(p, Szenarien::Automatisch).map_err(|e| e.to_string()).and_then(|(_, r)| both(r)),
        );
    });
}
