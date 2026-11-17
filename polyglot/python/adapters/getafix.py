"""getafix : modèle typé BT/BG, lecture/écriture CII, sous-ensemble de règles BR, PDF via pypdf.

La validation (Document.validate / validate_internal) applique les règles BR
codées à la main dans getafix/rules/*.py ; ValidationWarning (plausibilité)
ne compte pas comme erreur (README). Le profil XRECHNUNG n'existe pas dans
getafix.schema.types.Profile : ces factures sont refusées à la lecture.
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
    "extract": "extract_xml(Path) : l'API ne prend qu'un chemin, la lecture du fichier (cache OS) est dans la mesure",
    "parse": "Document.from_xml(xml.etree.ElementTree.fromstring)",
    "validate": "lecture + règles de getafix (89 des 223 règles BR du CEN, sans BR-CL, + 16 BR-FXEXT) ; avertissements non comptés",
    "write": "Document.to_xml().render() depuis le modèle lu",
    "embed": "attach_xml(bytes, bytes) : pièce jointe générique, pas de XMP Factur-X ni de PDF/A-3 (documenté)",
}

UNSUPPORTED = {
    "check_pdf": "pas de contrôle natif d'un PDF (la CLI compose extraction + validation)",
}

SURVEY_NOTES = {
    "validate": "règles BR de getafix (Document.validate), avertissements non comptés",
    "extract": "extract_xml(Path)",
}


def extract(pdf, path=None):
    if path is None:  # relecture d'un PDF produit par embed : passer par un fichier
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
        raise ValueError("aucune pièce jointe aux noms standard")
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
