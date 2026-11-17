//! pdfa 0.1.10 (cipansor/crateworks): XMP metadata, embedded files and the
//! Factur-X container, with its own PDF reader and writer (no PDF library).

use harness::{Lib, Op};
use pdfa::{DocumentInfo, EmbeddedFile, Level, WriteOptions, Xmp};

const LIB: Lib = Lib { lib: "pdfa", pkg: "pdfa", version: "0.1.10" };

const STAMP: &str = "2026-11-16T12:00:00+00:00";

fn extract(p: &[u8]) -> Result<Vec<u8>, String> {
    pdfa::embedded_file(p).ok_or_else(|| "no embedded file read".to_owned())
}

fn embed(host: &[u8], xml: &[u8]) -> Result<Vec<u8>, String> {
    let info = DocumentInfo {
        title: "Invoice".into(),
        author: "bench".into(),
        subject: "Invoice".into(),
        producer: "bench".into(),
        creator_tool: "bench".into(),
        timestamp: STAMP.into(),
    };
    let opts = WriteOptions { xmp: Xmp::factur_x(Level::En16931, info), timestamp: STAMP.into(), compress: false };
    pdfa::embed(host, &EmbeddedFile::factur_x(xml.to_vec()), &opts).map_err(|e| e.to_string())
}

fn main() {
    harness::run(LIB, |h| {
        h.extract("embedded_file : premier flux /Type /EmbeddedFile trouvé par recherche d'octets", extract);
        h.unsupported(Op::Parse, "pas de modèle de facture (crate sœur `factur-x`)");
        h.unsupported(Op::Validate, "pas de validation du XML");
        h.unsupported(Op::Write, "pas de génération de XML (crate sœur `factur-x`)");
        h.embed(
            "embed (Level::En16931, horodatage fixe) : réécrit tout le fichier avec une table xref classique et un nouveau XMP PDF/A-3B ; refuse les hôtes à flux d'objets / flux xref",
            embed,
            extract,
        );
        h.unsupported(
            Op::CheckPdf,
            "check_factur_x ne contrôle que le conteneur (XMP PDF/A-3, /Metadata, nom de la pièce jointe, /AFRelationship, somme de contrôle), pas le XML",
        );
    });
}
