"""PyCheval : lecture/écriture Factur-X MINIMUM à EN 16931, extraction et injection PDF.

Pas de validateur autonome : le parseur contrôle la structure du modèle
(ModelError, pycheval/model.py) mais n'applique aucune règle BR.
EXTENDED et XRECHNUNG sont refusés (UnsupportedProfileError, pycheval/parse.py).

Attention : extract_facturx_from_pdf boucle indéfiniment quand la première
entrée de /EmbeddedFiles n'est pas « factur-x.xml » (pycheval/pdf_extract.py,
`while doc: if doc[0] == FACTURX_FILENAME: break`) ; le harnais coupe l'appel
après CHECK_TIMEOUT secondes.
"""

from io import BytesIO

import pycheval
from pycheval.pdf_embed import _embed
from pycheval.pdf_extract import extract_facturx_from_pdf

PKG = "pycheval"

NOTES = {
    "extract": "extract_facturx_from_pdf(BytesIO), renvoie une str",
    "parse": "parse_xml(BytesIO)",
    "write": "generate_xml depuis le modèle de parse_xml",
    "embed": "fonction interne _embed (l'API publique embed_facturx_file_in_pdf lit le XML depuis un fichier), relation Data (défaut)",
}

UNSUPPORTED = {
    "validate": "pas de validateur (seulement des contrôles de structure du modèle à la lecture)",
    "check_pdf": "pas de contrôle de PDF (parse_pdf lit sans rapport d'erreurs)",
}

SURVEY_NOTES = {
    "extract": "boucle infinie si factur-x.xml n'est pas la première pièce jointe (coupée par le harnais)",
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
