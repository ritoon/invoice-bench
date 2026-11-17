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
        h.unsupported(Op::Extract, "extraction interne (module `pdf` privé) : seulement au travers de kosit_pdf / beide_pdf");
        h.unsupported(Op::Parse, "pas de modèle de facture : la librairie valide le XML sans le lire dans un modèle");
        h.validate(
            "Validator::beide(Szenarien::Automatisch) : KoSIT (XSD CII D16B, CEN EN 16931 1.3.16 en Schematron, XRechnung 3.0.2 / schematron 2.6.0 pour les XRECHNUNG, configuration KoSIT, CHANGELOG au 2026-01-31) + Mustang 2.26.0 (XSD et Schematron Factur-X du profil déclaré, XR 3.0 pour les XRECHNUNG) ; erreurs des deux outils additionnées",
            |x| both(v.beide(text(x)?, Szenarien::Automatisch).map_err(|e| e.to_string())?),
            summarize,
        );
        h.unsupported(Op::Write, "pas de génération de XML");
        h.unsupported(Op::Embed, "pas d'injection dans un PDF");
        h.check_pdf(
            "Validator::beide_pdf : extraction du XML puis KoSIT + Mustang ; aucun contrôle du conteneur PDF/A-3 (partie PDF de Mustang/veraPDF non reprise)",
            |p| v.beide_pdf(p, Szenarien::Automatisch).map_err(|e| e.to_string()).and_then(|(_, r)| both(r)),
        );
    });
}
