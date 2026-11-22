"""drafthorse (pretix): 1:1 ZUGFeRD model, CII parsing/writing, PDF embedding.

No PDF extraction ("We do not support parsing PDF files"), no
schematron: the only validation is the bundled Factur-X 1.0.07 XSD
(drafthorse/utils.py validate_xml, called by Document.serialize).
"""

from lxml import etree

from drafthorse.models.document import Document
from drafthorse.pdf import attach_xml
from drafthorse.utils import validate_xml

PKG = "drafthorse"

NOTES = {
    "validate": "Factur-X 1.0.07 EN16931 XSD (drafthorse.utils.validate_xml, stops at the first error)",
    "write": "Document.serialize(schema='FACTUR-X_EN16931'), revalidates the output against the XSD",
    "embed": "attach_xml, default options (detected profile, PDF/A-3 XMP)",
}

UNSUPPORTED = {
    "extract": 'no PDF parsing (README: "We do not support parsing PDF files")',
    "check_pdf": "no PDF check",
}

SURVEY_NOTES = {
    "validate": "Factur-X 1.0.07 EN16931 XSD only",
    "parse": "Document.parse(strict=True)",
}


def parse(xml, _path=None):
    return Document.parse(xml)


def validate(xml, _path=None):
    try:
        validate_xml(xml, "FACTUR-X_EN16931")
    except etree.XMLSyntaxError as e:
        return 1, ["XSD: " + str(e)[:80]]
    return 0, []


def write(model):
    return model.serialize(schema="FACTUR-X_EN16931")


def embed(host, xml):
    return attach_xml(host, xml)
