"""mcp-einvoicing-de : couche bibliothèque du serveur MCP, en mode local uniquement.

Validation : tools/invoice_validate.invoice_validate (async), qui exécute les
XSLT embarqués avec saxonche (extra [xslt2]) : FeRD Factur-X 1.09.2 pour le
profil EN 16931, CEN EN 16931 CII + KoSIT XRechnung pour XRECHNUNG.
La validation distante KoSIT est désactivée par défaut depuis la 0.11
(cloud_validate=False, EINVOICING_DE_KOSIT_ENABLE absent) ; on la force à
« 0 » ci-dessous par sécurité. _run_one_stylesheet recrée un PySaxonProcessor
et recompile la feuille XSLT à chaque appel : c'est le coût réel de l'API.
"""

import asyncio
import os

os.environ["EINVOICING_DE_KOSIT_ENABLE"] = "0"  # jamais d'envoi au validateur KoSIT distant

from mcp_einvoicing_de.serializers import ZUGFeRDCIIParser, ZUGFeRDCIISerializer  # noqa: E402
from mcp_einvoicing_de.tools import invoice_validate as _iv  # noqa: E402
from mcp_einvoicing_de.tools.invoice_parse import _extract_xml_from_pdf  # noqa: E402
from mcp_einvoicing_de.utils.pdf import embed_xml_in_pdf  # noqa: E402

assert not _iv._KOSIT_ENABLE_ENV, "validation KoSIT distante activée"

PKG = "mcp-einvoicing-de"

NOTES = {
    "extract": "_extract_xml_from_pdf du module invoice_parse (pikepdf, un passage par nom de fichier connu)",
    "parse": "ZUGFeRDCIIParser (modèle pydantic EN 16931)",
    "validate": "invoice_validate local : XSLT FeRD Factur-X 1.09.2 (EN 16931) ou CEN + KoSIT XRechnung 3.0.2 via saxonche, recompilé à chaque appel ; avertissements non comptés",
    "write": "ZUGFeRDCIISerializer depuis le modèle lu (modèle avec pertes : pièces jointes non reprises)",
    "embed": "utils.pdf.embed_xml_in_pdf (pikepdf, AFRelationship Alternative, XMP Factur-X)",
}

UNSUPPORTED = {
    "check_pdf": "pas de contrôle natif d'un PDF hybride (PDFEmbedder.identify ne lit que le XMP)",
}

SURVEY_NOTES = {
    "validate": "XSLT embarqués via saxonche, verdict de la librairie (erreurs fatal/error ; avertissements non comptés)",
}

_loop = asyncio.new_event_loop()


def extract(pdf, _path=None):
    return _extract_xml_from_pdf(pdf)


def parse(xml, _path=None):
    return ZUGFeRDCIIParser().parse(xml)


def validate(xml, _path=None):
    r = _loop.run_until_complete(_iv.invoice_validate(xml_content=xml.decode("utf-8"), cloud_validate=False))
    if r.get("validator_used") != "local_schematron":
        raise RuntimeError(f"validateur inattendu : {r.get('validator_used')!r} ({r})"[:200])
    rules = [f"{e['rule_id']}" for e in r["errors"]]
    return r["error_count"], rules


def write(model):
    return ZUGFeRDCIISerializer().serialize(model)


def embed(host, xml):
    return embed_xml_in_pdf(host, xml)
