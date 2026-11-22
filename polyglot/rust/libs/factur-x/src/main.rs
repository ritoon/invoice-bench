//! factur-x 0.3.1 (cipansor/crateworks): EN 16931 model with CII and UBL
//! generation and parsing, profile detection, a handful of arithmetic rules.
//! PDF handling is the sibling crate `pdfa`, benchmarked separately.

use factur_x::{Invoice, Report, Syntax};
use harness::{Findings, Lib, Op};

const LIB: Lib = Lib { lib: "factur-x", pkg: "factur-x", version: "0.3.1" };

fn text(x: &[u8]) -> Result<&str, String> {
    std::str::from_utf8(x).map_err(|e| format!("not UTF-8: {e}"))
}

fn parse(x: &[u8]) -> Result<Invoice, String> {
    factur_x::parse(text(x)?, Syntax::Cii).map_err(|e| e.to_string())
}

fn validate(x: &[u8]) -> Result<Report, String> {
    parse(x).map(|inv| factur_x::validate(&inv))
}

fn summarize(r: &Report) -> Findings {
    let rules: Vec<String> = r.rules().into_iter().map(str::to_owned).collect();
    Findings { errors: rules.len(), rules }
}

fn write(inv: &Invoice) -> Result<String, String> {
    Ok(factur_x::generate(inv, Syntax::Cii))
}

fn main() {
    harness::run(LIB, |h| {
        h.unsupported(Op::Extract, "no PDF parsing (sister crate `pdfa`)");
        h.parse("parse(Syntax::Cii) → Invoice", parse);
        h.validate(
            "adapter: parse + validate(&Invoice); 16 rules on the model (BR-01 to BR-06, BR-09, BR-16, BR-CO-10, BR-CO-13 to BR-CO-18, BR-S-08); no XSD, no profile or national rules",
            validate,
            summarize,
        );
        h.write("generate(Syntax::Cii), model read by parse", parse, write);
        h.unsupported(Op::Embed, "no PDF writing (sister crate `pdfa`)");
        h.unsupported(Op::CheckPdf, "no PDF parsing");
    });
}
