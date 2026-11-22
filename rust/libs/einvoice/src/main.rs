//! einvoice 0.1.1 (bschneppe/einvoice-rs): CII and UBL models generated from
//! the schemas, (de)serialised with a yaserde fork; Kotlin bindings by uniffi.

use einvoice::CrossIndustryInvoice;
use harness::{Lib, Op};

const LIB: Lib = Lib { lib: "einvoice", pkg: "einvoice", version: "0.1.1" };

fn text(x: &[u8]) -> Result<&str, String> {
    std::str::from_utf8(x).map_err(|e| format!("not UTF-8: {e}"))
}

fn parse(x: &[u8]) -> Result<CrossIndustryInvoice, String> {
    einvoice_deps_yaserde::de::from_str::<CrossIndustryInvoice>(text(x)?)
}

fn write(inv: &CrossIndustryInvoice) -> Result<String, String> {
    einvoice_deps_yaserde::ser::to_string(inv)
}

fn main() {
    harness::run(LIB, |h| {
        h.unsupported(
            Op::Extract,
            "private extraction (get_embedded_xml_file); only validate_zugferd_pdf(file path) uses it, no API on bytes",
        );
        h.parse("yaserde::de::from_str::<CrossIndustryInvoice> (model generated from the XSD, the one validate_invoice uses)", parse);
        h.unsupported(
            Op::Validate,
            "validate_invoice evaluates no rule: the 777 CII rules generated from the Schematron are `if false` stubs (src/cii/cii_business_rule_validator.rs)",
        );
        h.write("yaserde::ser::to_string, model read by from_str", parse, write);
        h.unsupported(Op::Embed, "no embedding in a PDF");
        h.unsupported(Op::CheckPdf, "validate_zugferd_pdf takes a file path and does not check the container; its rules are stubs");
    });
}
