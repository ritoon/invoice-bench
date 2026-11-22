"""Reference validator: saxonche (SaxonC-HE, no Java) + official XSLT.

- CEN EN 16931 1.3.17 (ConnectingEurope/eInvoicing-EN16931): xslt/EN16931-CII-validation.xslt
  for a CII invoice, xslt/EN16931-UBL-validation.xslt for a UBL invoice or credit note;
- KoSIT XRechnung 3.0.2 schematron 2.6.0 (itplr-kosit/xrechnung-schematron,
  schematron/cii/XRechnung-CII-validation.xsl or schematron/ubl/XRechnung-UBL-validation.xsl)
  in addition when BT-24 designates XRechnung, like the scenarios of
  validator-configuration-xrechnung.

Artifacts downloaded by `run.sh setup` into .artefacts/ (SHA-256 checked).
The stylesheets of a syntax are compiled on its first call (included in
first_call_ns for CII, the syntax of the bench), then reused.

Counting like the KoSIT validator: an error is any assertion whose flag is
fatal or error. For XRechnung, the custom levels of the scenarios of
validator-configuration-xrechnung v2026-08-31 apply (see _level). That is what
gives the reference result 6/26 (BR-CL-25 ×2, CII-SR-465 ×4 on XRECHNUNG).
No XSD (not provided by the CEN artifacts). The CII path is the one of the
bench; the UBL path serves the conformance cases (../../cases).
"""

import re
from pathlib import Path

from lxml import etree
from saxonche import PySaxonProcessor

PKG = "saxonche"

ART = Path(__file__).resolve().parent.parent / ".artefacts"
CEN_VERSION = "1.3.17"
XR_VERSION = "2.6.0"
XSLT = {
    "cii": (ART / f"cen-{CEN_VERSION}" / "xslt" / "EN16931-CII-validation.xslt",
            ART / f"xrechnung-{XR_VERSION}" / "schematron" / "cii" / "XRechnung-CII-validation.xsl"),
    "ubl": (ART / f"cen-ubl-{CEN_VERSION}" / "xslt" / "EN16931-UBL-validation.xslt",
            ART / f"xrechnung-{XR_VERSION}" / "schematron" / "ubl" / "XRechnung-UBL-validation.xsl"),
}

NOTES = {
    "validate": f"official CEN EN 16931 {CEN_VERSION} XSLT, CII or UBL by root (+ KoSIT XRechnung 3.0.2 schematron {XR_VERSION} if BT-24 is XRechnung) via saxonche, compiled on the first call; fatal/error errors with the levels of the KoSIT scenarios",
}

UNSUPPORTED = {
    "extract": "validator only",
    "parse": "validator only",
    "write": "validator only",
    "embed": "validator only",
    "check_pdf": "validator only",
}

SURVEY_NOTES = {
    "validate": f"CEN {CEN_VERSION} + KoSIT XRechnung {XR_VERSION} (XRECHNUNG), fatal/error errors with the custom levels of the KoSIT v2026-08-31 scenarios",
}

_ROOT = re.compile(rb"<(?:[\w.-]+:)?(CrossIndustryInvoice|Invoice|CreditNote)[\s>/]")
_SPEC = re.compile(rb"GuidelineSpecifiedDocumentContextParameter>\s*<(?:[\w.-]+:)?ID>([^<]+)<"
                   rb"|<(?:[\w.-]+:)?CustomizationID>([^<]+)<")
_SVRL = {"svrl": "http://purl.oclc.org/dsdl/svrl"}
_ERROR_FLAGS = ("fatal", "error")
_state = {}


def _level(rule, flag, spec, root):
    """The level of a rule under the scenario of validator-configuration-xrechnung
    v2026-08-31 that BT-24 selects (CIUS, extension or CVD), its flag otherwise:
    the customLevel elements of scenarios.xml, where they change whether the
    rule is an error."""
    if "xrechnung" not in spec:
        return flag
    ext, cvd = ":extension:" in spec, ":cvd_" in spec
    if rule in ("BR-CL-21", "BR-CL-23", "CII-SR-475", "CII-SR-476"):
        return "warning"
    if rule in ("CII-SR-452", "CII-SR-453", "CII-SR-454", "CII-SR-465", "CII-SR-466"):
        return "error"
    if rule == "BR-CL-13":
        return "information" if cvd else flag
    if rule in ("BR-CL-10", "BR-CL-11", "BR-CL-24", "BR-CL-25", "BR-CL-26"):
        return "information" if ext else flag
    if rule in ("BR-CO-16", "UBL-CR-470"):  # UBL extension scenario only
        return "information" if ext and root != "CrossIndustryInvoice" else flag
    if rule == "UBL-CR-646":
        if ext:
            return "information"
        return "error" if cvd or root != "CreditNote" else flag
    return flag


def _compiled(syntax):
    if syntax not in _state:
        for p in XSLT[syntax]:
            if not p.exists():
                raise FileNotFoundError(f"{p} missing: run ./run.sh setup")
        if "proc" not in _state:
            _state["proc"] = PySaxonProcessor(license=False)
        xslt = _state["proc"].new_xslt30_processor()
        cen, xr = XSLT[syntax]
        _state[syntax] = (xslt.compile_stylesheet(stylesheet_file=str(cen)),
                          xslt.compile_stylesheet(stylesheet_file=str(xr)))
    return _state["proc"], _state[syntax]


def _errors(svrl_text, spec, root):
    doc = etree.fromstring(svrl_text.encode("utf-8"))
    out = []
    for el in doc.xpath("//svrl:failed-assert | //svrl:successful-report", namespaces=_SVRL):
        rule, flag = el.get("id") or "?", (el.get("flag") or "error").lower()
        if _level(rule, flag, spec, root) in _ERROR_FLAGS:
            out.append(rule)
    return out


def validate(xml, _path=None):
    m = _ROOT.search(xml[:65536])  # after the licence comments of the samples
    root = m.group(1).decode() if m else "CrossIndustryInvoice"
    proc, (cen, xr) = _compiled("cii" if root == "CrossIndustryInvoice" else "ubl")
    m = _SPEC.search(xml[:65536])
    spec = (m.group(1) or m.group(2)).decode("utf-8", "replace").strip() if m else ""
    node = proc.parse_xml(xml_text=xml.decode("utf-8-sig"))
    rules = _errors(cen.transform_to_string(xdm_node=node), spec, root)
    if "xrechnung" in spec:
        rules += _errors(xr.transform_to_string(xdm_node=node), spec, root)
    return len(rules), rules
