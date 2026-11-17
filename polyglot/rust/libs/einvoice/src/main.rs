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
            "extraction privée (get_embedded_xml_file) ; seule validate_zugferd_pdf(chemin de fichier) l'utilise, pas d'API sur octets",
        );
        h.parse("yaserde::de::from_str::<CrossIndustryInvoice> (modèle généré du XSD, celui de validate_invoice)", parse);
        h.unsupported(
            Op::Validate,
            "validate_invoice n'évalue aucune règle : les 777 règles CII générées depuis le Schematron sont des bouchons `if false` (src/cii/cii_business_rule_validator.rs)",
        );
        h.write("yaserde::ser::to_string, modèle lu par from_str", parse, write);
        h.unsupported(Op::Embed, "pas d'injection dans un PDF");
        h.unsupported(Op::CheckPdf, "validate_zugferd_pdf prend un chemin de fichier et ne contrôle pas le conteneur ; ses règles sont des bouchons");
    });
}
