"""One adapter per library.

Each module exposes PKG (PyPI name) and, for each operation it can do:

    extract(pdf: bytes, path) -> bytes | str        PDF -> XML of the invoice
    parse(xml: bytes, path) -> model                XML -> model of the library
    validate(xml: bytes, path) -> (n, [rules])      full validation, n = number of errors
    write(model) -> bytes | str                     model -> CII XML
    hand_model() -> model                           failing a parser: EN16931_Einfach by hand
    embed(host: bytes, xml: bytes) -> bytes         host PDF + XML -> Factur-X PDF
    check_pdf(pdf: bytes, path) -> (n, [rules])     native check of a hybrid PDF

A missing operation counts as "unsupported" (reason in UNSUPPORTED[op]).
NOTES[op or benchmark] is copied into the note field of the JSONL.
"""
