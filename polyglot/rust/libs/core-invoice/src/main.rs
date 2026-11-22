//! core-invoice + core-invoice-formats 2.0.3 (proxeon): EN 16931 and Peppol
//! PINT semantic model with its rules in Rust, UBL 2.1 and CII D16B codecs.
//! Built with the optional `xrechnung` feature (BR-DE rules).

use core_invoice::{Invoice, Report, Severity};
use harness::{Findings, Lib, Op};

const LIB: Lib = Lib { lib: "core-invoice", pkg: "core-invoice-formats", version: "2.0.3" };

fn text(x: &[u8]) -> Result<&str, String> {
    std::str::from_utf8(x).map_err(|e| format!("not UTF-8: {e}"))
}

fn parse(x: &[u8]) -> Result<Invoice, String> {
    core_invoice_formats::read(text(x)?).map_err(|e| e.to_string())
}

fn validate(x: &[u8]) -> Result<Report, String> {
    core_invoice_formats::validate_xml(text(x)?, None).map_err(|e| e.to_string())
}

fn summarize(r: &Report) -> Findings {
    let rules: Vec<String> =
        r.findings.iter().filter(|f| f.severity == Severity::Fatal).map(|f| f.id.to_owned()).collect();
    Findings { errors: rules.len(), rules }
}

fn write(inv: &Invoice) -> Result<String, String> {
    core_invoice_formats::cii::write_unchecked(inv).map_err(|e| e.to_string())
}

fn main() {
    harness::run(LIB, |h| {
        h.unsupported(Op::Extract, "no PDF parsing");
        h.parse("read → Invoice (semantic model)", parse);
        h.validate(
            "validate_xml(profile from BT-24): EN 16931 rules reimplemented in Rust (aligned with CEN validation-1.3.16), currency lists, Peppol syntax if Peppol profile; BR-DE (feature `xrechnung`) on the XRECHNUNG files; no XSD or CII-SR rules",
            validate,
            summarize,
        );
        h.write("cii::write_unchecked, model read by read", parse, write);
        h.unsupported(Op::Embed, "no embedding in a PDF");
        h.unsupported(Op::CheckPdf, "no PDF parsing");
    });
}
