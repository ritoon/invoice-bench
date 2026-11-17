//! Floor: lopdf 0.47.0, the general-purpose PDF library most Rust Factur-X
//! crates are built on (zugferd-rs, en16931-formats, erechnung-validator,
//! einvoice), with the invoice extraction written by hand here: load the
//! document, walk /Root/Names/EmbeddedFiles (with /Kids), take the first file
//! under a standard invoice name, inflate it. Not a Factur-X library.

use harness::{Lib, Op};
use lopdf::{Dictionary, Document, Object};

const LIB: Lib = Lib { lib: "lopdf", pkg: "lopdf", version: "0.47.0" };

const NAMES: [&str; 5] = ["factur-x.xml", "zugferd-invoice.xml", "xrechnung.xml", "order-x.xml", "cii.xml"];

fn dict<'a>(doc: &'a Document, o: &'a Object) -> Option<&'a Dictionary> {
    doc.dereference(o).ok()?.1.as_dict().ok()
}

fn walk(doc: &Document, node: &Dictionary, depth: u32, out: &mut Vec<(String, Object)>) {
    if depth > 32 {
        return;
    }
    if let Ok(names) = node.get(b"Names").and_then(Object::as_array) {
        for pair in names.chunks(2) {
            if let [k, v] = pair {
                if let Ok((_, Object::String(s, _))) = doc.dereference(k) {
                    out.push((String::from_utf8_lossy(s).into_owned(), v.clone()));
                }
            }
        }
    }
    if let Ok(kids) = node.get(b"Kids").and_then(Object::as_array) {
        for k in kids {
            if let Some(d) = dict(doc, k) {
                walk(doc, d, depth + 1, out);
            }
        }
    }
}

fn extract(pdf: &[u8]) -> Result<Vec<u8>, String> {
    let doc = Document::load_mem(pdf).map_err(|e| e.to_string())?;
    let cat = doc.catalog().map_err(|e| e.to_string())?;
    let tree = cat
        .get(b"Names")
        .ok()
        .and_then(|n| dict(&doc, n))
        .and_then(|n| n.get(b"EmbeddedFiles").ok())
        .and_then(|e| dict(&doc, e))
        .ok_or("no /Names/EmbeddedFiles")?;
    let mut files = Vec::new();
    walk(&doc, tree, 0, &mut files);
    let (_, spec) = files
        .iter()
        .find(|(n, _)| NAMES.iter().any(|s| s.eq_ignore_ascii_case(n)))
        .ok_or("no attachment under a standard invoice name")?;
    let spec = dict(&doc, spec).ok_or("file specification is not a dictionary")?;
    let ef = spec.get(b"EF").ok().and_then(|e| dict(&doc, e)).ok_or("no /EF")?;
    let f = ef.get(b"UF").or_else(|_| ef.get(b"F")).map_err(|e| e.to_string())?;
    let stream = doc.dereference(f).map_err(|e| e.to_string())?.1.as_stream().map_err(|e| e.to_string())?;
    if stream.dict.get(b"Filter").is_ok() {
        stream.decompressed_content().map_err(|e| e.to_string())
    } else {
        Ok(stream.content.clone())
    }
}

fn main() {
    harness::run(LIB, |h| {
        h.extract(
            "plancher : extraction écrite à la main sur lopdf (load_mem complet, arbre /Names/EmbeddedFiles avec /Kids, inflate) ; pas une librairie Factur-X",
            extract,
        );
        for op in [Op::Parse, Op::Validate, Op::Write, Op::Embed, Op::CheckPdf] {
            h.unsupported(op, "plancher d'extraction seulement");
        }
    });
}
