"""factur-x (Akretion): extraction, parsing into an EN 16931 dictionary, CII writing, XSD, embedding.

The library's Schematron validation (xml_check_schematron) is not
local: it sends the XML by HTTP POST to a Saxon server
(http://localhost:5000/transform by default, facturx/facturx.py). It is
therefore disabled everywhere here; offline validation is limited to the XSD.
"""

import facturx

PKG = "factur-x"

NOTES = {
    "extract": "check_xsd=False (by default the extraction also validates the XML against the XSD)",
    "parse": "parse_ubl_cii_xml, XSD included (default)",
    "validate": "Factur-X 1.09.2 XSD only (the schematron requires a Saxon HTTP server)",
    "write": "generate_cii_xml from the dictionary of parse_ubl_cii_xml, XSD included (default), schematron disabled",
    "embed": "generate_from_binary, default options (XSD included, goes through a temporary file, PDF/A-3 XMP)",
}

UNSUPPORTED = {
    "check_pdf": "no native check of a hybrid PDF",
}

SURVEY_NOTES = {
    "validate": "Factur-X 1.09.2 XSD only (xml_check_xsd)",
    "parse": "parse_ubl_cii_xml (XSD included)",
    "extract": "get_xml_from_pdf(check_xsd=False)",
}


def extract(pdf, _path=None):
    _name, xml = facturx.get_xml_from_pdf(pdf, check_xsd=False)
    if not xml:
        raise ValueError("no Factur-X/ZUGFeRD XML attachment found")
    return xml


def parse(xml, _path=None):
    return facturx.parse_ubl_cii_xml(xml)


def validate(xml, _path=None):
    try:
        facturx.xml_check_xsd(xml)
    except Exception as e:  # the library raises at the first XSD error
        return 1, ["XSD: " + str(e).split(". Here is the error")[0][:80]]
    return 0, []


def write(model):
    return facturx.generate_cii_xml(model, check_schematron=False)


def embed(host, xml):
    return facturx.generate_from_binary(host, xml)


def survey_extract_extra(files):
    """How many PDFs does the default extraction (check_xsd=True) keep?"""
    ok = 0
    for f in files:
        try:
            with open(f, "rb") as fh:
                _n, x = facturx.get_xml_from_pdf(fh.read())
            ok += bool(x)
        except Exception:
            pass
    return f"with check_xsd=True (default): {ok}/{len(files)}"
