"""Validateur de référence : saxonche (SaxonC-HE, sans Java) + XSLT officiels.

- CEN EN 16931 CII 1.3.17 (ConnectingEurope/eInvoicing-EN16931, xslt/EN16931-CII-validation.xslt)
  pour toutes les factures ;
- KoSIT XRechnung 3.0.2 schematron 2.6.0 (itplr-kosit/xrechnung-schematron,
  schematron/cii/XRechnung-CII-validation.xsl) en plus quand BT-24 désigne XRechnung,
  comme le scénario CII du validator-configuration-xrechnung.

Artefacts téléchargés par `run.sh setup` dans .artefacts/ (SHA-256 vérifiés).
Les feuilles sont compilées au premier appel (compris dans first_call_ns), puis
réutilisées.

Comptage comme le validateur KoSIT : est une erreur toute assertion de flag
fatal ou error. Pour le scénario XRechnung, les niveaux personnalisés du
validator-configuration-xrechnung v2026-08-31 (CHANGELOG « custom level ») sont
appliqués : CII-SR-465 et CII-SR-466 (warning dans le CEN) deviennent des
erreurs, CII-SR-475 et CII-SR-476 deviennent des informations. C'est ce qui
donne le résultat de référence 6/26 (BR-CL-25 ×2, CII-SR-465 ×4 sur XRECHNUNG).
Pas de XSD (non fourni par les artefacts CEN CII).
"""

import re
from pathlib import Path

from lxml import etree
from saxonche import PySaxonProcessor

PKG = "saxonche"

ART = Path(__file__).resolve().parent.parent / ".artefacts"
CEN_VERSION = "1.3.17"
XR_VERSION = "2.6.0"
CEN_XSLT = ART / f"cen-{CEN_VERSION}" / "xslt" / "EN16931-CII-validation.xslt"
XR_XSLT = ART / f"xrechnung-{XR_VERSION}" / "schematron" / "cii" / "XRechnung-CII-validation.xsl"

NOTES = {
    "validate": f"XSLT officiels CEN EN 16931 CII {CEN_VERSION} (+ KoSIT XRechnung 3.0.2 schematron {XR_VERSION} si BT-24 XRechnung) via saxonche, compilés au premier appel ; erreurs fatal/error avec les niveaux KoSIT",
}

UNSUPPORTED = {
    "extract": "validateur seul",
    "parse": "validateur seul",
    "write": "validateur seul",
    "embed": "validateur seul",
    "check_pdf": "validateur seul",
}

SURVEY_NOTES = {
    "validate": f"CEN {CEN_VERSION} + KoSIT XRechnung {XR_VERSION} (XRECHNUNG), erreurs fatal/error avec les niveaux personnalisés KoSIT v2026-08-31",
}

_BT24 = re.compile(rb"GuidelineSpecifiedDocumentContextParameter>\s*<(?:\w+:)?ID>([^<]+)<")
_SVRL = {"svrl": "http://purl.oclc.org/dsdl/svrl"}
# validator-configuration-xrechnung v2026-08-31, scénario CII : niveaux personnalisés
_KOSIT_CUSTOM_LEVEL = {"CII-SR-465": "error", "CII-SR-466": "error", "CII-SR-475": "information", "CII-SR-476": "information"}
_ERROR_FLAGS = ("fatal", "error")
_state = {}


def _compiled():
    if not _state:
        for p in (CEN_XSLT, XR_XSLT):
            if not p.exists():
                raise FileNotFoundError(f"{p} absent : lancer ./run.sh setup")
        proc = PySaxonProcessor(license=False)
        xslt = proc.new_xslt30_processor()
        _state["proc"] = proc
        _state["cen"] = xslt.compile_stylesheet(stylesheet_file=str(CEN_XSLT))
        _state["xr"] = xslt.compile_stylesheet(stylesheet_file=str(XR_XSLT))
    return _state["proc"], _state["cen"], _state["xr"]


def _errors(svrl_text, xrechnung):
    root = etree.fromstring(svrl_text.encode("utf-8"))
    out = []
    for el in root.xpath("//svrl:failed-assert | //svrl:successful-report", namespaces=_SVRL):
        rule, flag = el.get("id") or "?", (el.get("flag") or "error").lower()
        if xrechnung:
            flag = _KOSIT_CUSTOM_LEVEL.get(rule, flag)
        if flag in _ERROR_FLAGS:
            out.append(rule)
    return out


def validate(xml, _path=None):
    proc, cen, xr = _compiled()
    node = proc.parse_xml(xml_text=xml.decode("utf-8-sig"))
    m = _BT24.search(xml)
    xrechnung = bool(m and b"xrechnung" in m.group(1))
    rules = _errors(cen.transform_to_string(xdm_node=node), xrechnung)
    if xrechnung:
        rules += _errors(xr.transform_to_string(xdm_node=node), True)
    return len(rules), rules
