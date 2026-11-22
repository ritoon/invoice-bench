//! en16931 + en16931-formats 0.10.1 (hupe1980/en16931): EN 16931 model and
//! business rules re-implemented in Rust, UBL / CII / ZUGFeRD read and write.

use std::time::{Duration, SystemTime};

use en16931::profiles::En16931;
use en16931::{Severity, Validated};
use en16931_formats::zugferd::{self, EmbedOptions, Limits};
use en16931_formats::{Policy, Received, Recommendation, Syntax};
use harness::{Findings, Lib};

const LIB: Lib = Lib { lib: "en16931", pkg: "en16931-formats", version: "0.10.1" };

fn extract(pdf: &[u8]) -> Result<Vec<u8>, String> {
    zugferd::extract(pdf, &Limits::default()).map(|x| x.xml).map_err(|e| e.to_string())
}

fn parse(x: &[u8]) -> Result<en16931::Invoice, String> {
    en16931_formats::read(x).map(|r| r.invoice).map_err(|e| e.to_string())
}

/// `receive` with the default policy: read, classify BT-24, validate under the
/// declared profile (XRechnung 3.0 for the XRECHNUNG samples, EN 16931 core
/// otherwise) and recommend. Never fails: an unreadable document is a reject.
fn receive(x: &[u8]) -> Result<Received, String> {
    Ok(en16931_formats::receive(x, &Policy::default()))
}

/// Errors: the fatal entries of the container, syntax and schema layers and the
/// fatal findings of the semantic layer; a reject with none of them (unknown
/// BT-24…) counts its reasons.
fn summarize(r: &Received) -> Findings {
    let rep = &r.report_declared;
    let mut rules: Vec<String> = rep
        .entries()
        .filter(|e| e.level() == Severity::Fatal)
        .map(|e| e.rule.clone().unwrap_or_else(|| format!("{:?}", e.kind)))
        .collect();
    if let Some(sem) = &rep.semantic.report {
        rules.extend(sem.fatal().map(|f| f.rule.clone()));
    }
    if rules.is_empty() && r.recommendation == Recommendation::Reject {
        rules.extend(r.reasons.iter().map(|w| format!("REJECT: {w}")));
    }
    Findings { errors: rules.len(), rules }
}

fn write(inv: &en16931::Invoice) -> Result<String, String> {
    en16931_formats::write_unchecked(inv, Syntax::Cii).map(|w| w.xml).map_err(|e| e.to_string())
}

fn embed(host: &[u8], xml: &[u8]) -> Result<Vec<u8>, String> {
    let inv = parse(xml)?;
    let proof = Validated::<En16931>::new(inv).map_err(|e| format!("not a valid EN 16931 invoice: {e}"))?;
    let opts = EmbedOptions {
        modified: Some(SystemTime::UNIX_EPOCH + Duration::from_secs(1_794_830_400)), // 2026-11-16T12:00:00Z
        ..EmbedOptions::default()
    };
    zugferd::embed(host, &proof, &opts).map_err(|e| e.to_string())
}

/// Native check of a hybrid PDF: `receive` on the PDF reads the container
/// (XMP, /AF, the embedded-files tree, their agreement with BT-24) and
/// validates the payload, as for an XML document.
fn check(pdf: &[u8]) -> Result<usize, String> {
    let r = en16931_formats::receive(pdf, &Policy::default());
    if r.invoice.is_none() {
        return Err(format!("no invoice read: {}", r.reasons.join("; ")));
    }
    Ok(summarize(&r).errors)
}

fn main() {
    harness::run(LIB, |h| {
        h.extract("zugferd::extract (default Limits)", extract);
        h.parse("read → Invoice (EN 16931 semantic model)", parse);
        h.validate(
            "receive(Policy::default()): syntax and schema layers (CII-SR/CII-DT rules generated from the CEN artifacts) + EN 16931 business rules reimplemented in Rust, profile declared by BT-24 (XRechnung 3.0: BR-DE; ZUGFeRD EXTENDED/BASIC: CEN core rules); CEN validation-1.3.17, KoSIT xrechnung-schematron v2.6.0",
            receive,
            summarize,
        );
        h.write("write_unchecked(Syntax::Cii), model read by read", parse, write);
        h.embed(
            "adapter: read + Validated::<En16931>::new + zugferd::embed (the library rewrites the XML from its model; requires a PDF/A-3 host without Factur-X metadata, does not convert to PDF/A)",
            embed,
            extract,
        );
        h.check_pdf(
            "receive on the PDF: container layer (XMP, /AF, EmbeddedFiles tree, agreement with BT-24) + full validation of the XML; no PDF/A-3 check in the ISO 19005 sense",
            check,
        );
    });
}
