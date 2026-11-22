"""Reference validator: saxonche (SaxonC-HE, no Java) + official XSLT.

- CEN EN 16931 CII 1.3.17 (ConnectingEurope/eInvoicing-EN16931, xslt/EN16931-CII-validation.xslt)
  for all invoices;
- KoSIT XRechnung 3.0.2 schematron 2.6.0 (itplr-kosit/xrechnung-schematron,
  schematron/cii/XRechnung-CII-validation.xsl) in addition when BT-24 designates XRechnung,
  like the CII scenario of validator-configuration-xrechnung.

Artifacts downloaded by `run.sh setup` into .artefacts/ (SHA-256 checked).
The stylesheets are compiled on the first call (included in first_call_ns), then
reused.

Counting like the KoSIT validator: an error is any assertion whose flag is
fatal or error. For the XRechnung scenario, the custom levels of
validator-configuration-xrechnung v2026-08-31 (CHANGELOG "custom level") are
applied: CII-SR-465 and CII-SR-466 (warning in the CEN) become
errors, CII-SR-475 and CII-SR-476 become information. That is what
gives the reference result 6/26 (BR-CL-25 ×2, CII-SR-465 ×4 on XRECHNUNG).
No XSD (not provided by the CEN CII artifacts).
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
    "validate": f"official CEN EN 16931 CII {CEN_VERSION} XSLT (+ KoSIT XRechnung 3.0.2 schematron {XR_VERSION} if BT-24 is XRechnung) via saxonche, compiled on the first call; fatal/error errors with the KoSIT levels",
}

UNSUPPORTED = {
    "extract": "validator only",
    "parse": "validator only",
    "write": "validator only",
    "embed": "validator only",
    "check_pdf": "validator only",
}

SURVEY_NOTES = {
    "validate": f"CEN {CEN_VERSION} + KoSIT XRechnung {XR_VERSION} (XRECHNUNG), fatal/error errors with the KoSIT v2026-08-31 custom levels",
}

_BT24 = re.compile(rb"GuidelineSpecifiedDocumentContextParameter>\s*<(?:\w+:)?ID>([^<]+)<")
_SVRL = {"svrl": "http://purl.oclc.org/dsdl/svrl"}
# validator-configuration-xrechnung v2026-08-31, CII scenario: custom levels
_KOSIT_CUSTOM_LEVEL = {"CII-SR-465": "error", "CII-SR-466": "error", "CII-SR-475": "information", "CII-SR-476": "information"}
_ERROR_FLAGS = ("fatal", "error")
_state = {}


def _compiled():
    if not _state:
        for p in (CEN_XSLT, XR_XSLT):
            if not p.exists():
                raise FileNotFoundError(f"{p} missing: run ./run.sh setup")
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
