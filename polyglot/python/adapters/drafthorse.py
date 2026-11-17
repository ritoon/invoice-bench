"""drafthorse (pretix) : modèle ZUGFeRD 1:1, lecture/écriture CII, injection PDF.

Pas d'extraction PDF (« We do not support parsing PDF files »), pas de
schematron : la seule validation est le XSD Factur-X 1.0.07 embarqué
(drafthorse/utils.py validate_xml, appelée par Document.serialize).
"""

from lxml import etree

from drafthorse.models.document import Document
from drafthorse.pdf import attach_xml
from drafthorse.utils import validate_xml

PKG = "drafthorse"

NOTES = {
    "validate": "XSD Factur-X 1.0.07 EN16931 (drafthorse.utils.validate_xml, arrêt à la première erreur)",
    "write": "Document.serialize(schema='FACTUR-X_EN16931'), revalide la sortie contre le XSD",
    "embed": "attach_xml, options par défaut (profil détecté, XMP PDF/A-3)",
}

UNSUPPORTED = {
    "extract": "pas de lecture de PDF (README : « We do not support parsing PDF files »)",
    "check_pdf": "pas de contrôle de PDF",
}

SURVEY_NOTES = {
    "validate": "XSD Factur-X 1.0.07 EN16931 seulement",
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
