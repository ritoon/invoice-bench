//! zugferd-rs 0.2.2 (philipp-berger, port of zugferd-ts): CII / UBL
//! generation and parsing, PDF/A-3b embedding and extraction, validation by
//! libxml2 (XSD + the XPath 1.0 subset of the CEN Schematron).
//!
//! The schemas ship in the crate's `assets/schemas`; `run.sh` points
//! `ZUGFERD_SCHEMAS_DIR` at them.

use harness::{Findings, Lib, Op};
use zugferd::{
    EmbedOptions, Invoice, ValidateXmlOptions, ValidationResult, ValidationSeverity, ZugferdProfile,
};

const LIB: Lib = Lib { lib: "zugferd-rs", pkg: "zugferd-rs", version: "0.2.2" };

fn text(x: &[u8]) -> Result<&str, String> {
    std::str::from_utf8(x).map_err(|e| format!("not UTF-8: {e}"))
}

fn extract(p: &[u8]) -> Result<Vec<u8>, String> {
    zugferd::extract_xml(p).map(String::into_bytes).map_err(|e| e.to_string())
}

fn parse(x: &[u8]) -> Result<Invoice, String> {
    zugferd::parse_cii_xml(text(x)?, None).map_err(|e| e.to_string())
}

fn validate(x: &[u8]) -> Result<ValidationResult, String> {
    zugferd::validate_xml(x, ValidateXmlOptions::default()).map_err(|e| e.to_string())
}

fn summarize(r: &ValidationResult) -> Findings {
    let rules: Vec<String> = r
        .issues
        .iter()
        .filter(|i| i.severity == ValidationSeverity::Error)
        .map(|i| i.rule_id.clone().unwrap_or_else(|| "?".into()))
        .collect();
    Findings { errors: rules.len(), rules }
}

fn write(inv: &Invoice) -> Result<String, String> {
    Ok(zugferd::serialize_to_xml(inv, ZugferdProfile::En16931))
}

fn embed(host: &[u8], xml: &[u8]) -> Result<Vec<u8>, String> {
    let opts = EmbedOptions { profile: Some(ZugferdProfile::En16931), ..EmbedOptions::default() };
    zugferd::embed_xml(host, text(xml)?, opts).map_err(|e| e.to_string())
}

fn main() {
    if std::env::var_os("ZUGFERD_SCHEMAS_DIR").is_none() {
        eprintln!("warning: ZUGFERD_SCHEMAS_DIR is not set (run.sh sets it to the crate's assets/schemas)");
    }
    harness::run(LIB, |h| {
        h.extract("extract_xml (lopdf)", extract);
        h.parse("parse_cii_xml → Invoice", parse);
        h.validate(
            "validate_xml (défaut) : XSD Factur-X 1.09.2 du profil déclaré + Schematron CEN EN 16931 évalué par libxml2 (XPath 1.0 seulement : les règles XPath 2.0 sont sautées, résultat déclaré incomplet) ; pas de BR-DE",
            validate,
            summarize,
        );
        h.write("serialize_to_xml(ZugferdProfile::En16931), modèle lu par parse_cii_xml", parse, write);
        h.embed(
            "embed_xml (profil EN 16931) : réécrit le PDF en PDF/A-3b (XMP Factur-X, OutputIntent sRGB) ; le XML est injecté tel quel",
            embed,
            extract,
        );
        h.unsupported(
            Op::CheckPdf,
            "pas de contrôle complet natif : validate_pdf_a ne vérifie que le conteneur (~21 règles PDF/A-3b), pas le XML",
        );
    });
}
