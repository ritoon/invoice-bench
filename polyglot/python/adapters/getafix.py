"""getafix: typed BT/BG model, CII parsing/writing, subset of the BR rules, PDF via pypdf.

The validation (Document.validate / validate_internal) applies the BR rules
hand-coded in getafix/rules/*.py; ValidationWarning (plausibility)
does not count as an error (README). The XRECHNUNG profile does not exist in
getafix.schema.types.Profile: these invoices are rejected at parse time.
"""

import os
import tempfile
import xml.etree.ElementTree as ET
from pathlib import Path

from getafix.errors import ValidationWarning
from getafix.pdf import attach_xml, extract_xml
from getafix.schema.document import Document

PKG = "getafix"

NOTES = {
    "extract": "extract_xml(Path): the API only takes a path, reading the file (OS cache) is inside the measurement",
    "parse": "Document.from_xml(xml.etree.ElementTree.fromstring)",
    "validate": "parsing + getafix rules (89 of the 223 CEN BR rules, no BR-CL, + 16 BR-FXEXT); warnings not counted",
    "write": "Document.to_xml().render() from the parsed model",
    "embed": "attach_xml(bytes, bytes): generic attachment, no Factur-X XMP nor PDF/A-3 (documented)",
}

UNSUPPORTED = {
    "check_pdf": "no native check of a PDF (the CLI combines extraction + validation)",
}

SURVEY_NOTES = {
    "validate": "getafix BR rules (Document.validate), warnings not counted",
    "extract": "extract_xml(Path)",
}


def extract(pdf, path=None):
    if path is None:  # reading back a PDF produced by embed: go through a file
        fd, tmp = tempfile.mkstemp(suffix=".pdf")
        try:
            with os.fdopen(fd, "wb") as fh:
                fh.write(pdf)
            return _extract_path(tmp)
        finally:
            os.unlink(tmp)
    return _extract_path(path)


def _extract_path(path):
    x = extract_xml(Path(path))
    if x is None:
        raise ValueError("no attachment with a standard name")
    return x


def parse(xml, _path=None):
    return Document.from_xml(ET.fromstring(xml))


def validate(xml, _path=None):
    doc = Document.from_xml(ET.fromstring(xml))
    findings = doc.validate_internal(doc.context.guideline.id)
    errors = [e for e in findings if not isinstance(e, ValidationWarning)]
    return len(errors), [e.code for e in errors]


def write(model):
    return model.to_xml().render()


def embed(host, xml):
    return attach_xml(host, xml)
