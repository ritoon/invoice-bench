"""PyCheval: Factur-X MINIMUM to EN 16931 parsing/writing, PDF extraction and embedding.

No standalone validator: the parser checks the structure of the model
(ModelError, pycheval/model.py) but applies no BR rule.
EXTENDED and XRECHNUNG are rejected (UnsupportedProfileError, pycheval/parse.py).

Warning: extract_facturx_from_pdf loops forever when the first
entry of /EmbeddedFiles is not "factur-x.xml" (pycheval/pdf_extract.py,
`while doc: if doc[0] == FACTURX_FILENAME: break`); the harness cuts the call
off after CHECK_TIMEOUT seconds.
"""

from io import BytesIO

import pycheval
from pycheval.pdf_embed import _embed
from pycheval.pdf_extract import extract_facturx_from_pdf

PKG = "pycheval"

NOTES = {
    "extract": "extract_facturx_from_pdf(BytesIO), returns a str",
    "parse": "parse_xml(BytesIO)",
    "write": "generate_xml from the model of parse_xml",
    "embed": "internal function _embed (the public API embed_facturx_file_in_pdf reads the XML from a file), Data relationship (default)",
}

UNSUPPORTED = {
    "validate": "no validator (only structural checks of the model at parse time)",
    "check_pdf": "no PDF check (parse_pdf reads without an error report)",
}

SURVEY_NOTES = {
    "extract": "infinite loop if factur-x.xml is not the first attachment (cut off by the harness)",
}


def extract(pdf, _path=None):
    xml, _rel = extract_facturx_from_pdf(BytesIO(pdf))
    return xml


def parse(xml, _path=None):
    return pycheval.parse_xml(BytesIO(xml))


def write(model):
    return pycheval.generate_xml(model)


def embed(host, xml):
    return _embed(BytesIO(host), xml.decode("utf-8"), "EN 16931", relationship=pycheval.FileRelationship.DATA)
