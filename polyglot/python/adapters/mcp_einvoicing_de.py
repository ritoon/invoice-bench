"""mcp-einvoicing-de: library layer of the MCP server, in local mode only.

Validation: tools/invoice_validate.invoice_validate (async), which runs the
bundled XSLT with saxonche (extra [xslt2]): FeRD Factur-X 1.09.2 for the
EN 16931 profile, CEN EN 16931 CII + KoSIT XRechnung for XRECHNUNG.
Remote KoSIT validation has been disabled by default since 0.11
(cloud_validate=False, EINVOICING_DE_KOSIT_ENABLE unset); we force it to
"0" below to be safe. _run_one_stylesheet recreates a PySaxonProcessor
and recompiles the XSLT stylesheet on every call: that is the real cost of the API.
"""

import asyncio
import os

os.environ["EINVOICING_DE_KOSIT_ENABLE"] = "0"  # never send anything to the remote KoSIT validator

from mcp_einvoicing_de.serializers import ZUGFeRDCIIParser, ZUGFeRDCIISerializer  # noqa: E402
from mcp_einvoicing_de.tools import invoice_validate as _iv  # noqa: E402
from mcp_einvoicing_de.tools.invoice_parse import _extract_xml_from_pdf  # noqa: E402
from mcp_einvoicing_de.utils.pdf import embed_xml_in_pdf  # noqa: E402

assert not _iv._KOSIT_ENABLE_ENV, "remote KoSIT validation enabled"

PKG = "mcp-einvoicing-de"

NOTES = {
    "extract": "_extract_xml_from_pdf of the invoice_parse module (pikepdf, one pass per known file name)",
    "parse": "ZUGFeRDCIIParser (EN 16931 pydantic model)",
    "validate": "local invoice_validate: FeRD Factur-X 1.09.2 XSLT (EN 16931) or CEN + KoSIT XRechnung 3.0.2 via saxonche, recompiled on every call; warnings not counted",
    "write": "ZUGFeRDCIISerializer from the parsed model (lossy model: attachments not carried over)",
    "embed": "utils.pdf.embed_xml_in_pdf (pikepdf, AFRelationship Alternative, XMP Factur-X)",
}

UNSUPPORTED = {
    "check_pdf": "no native check of a hybrid PDF (PDFEmbedder.identify only reads the XMP)",
}

SURVEY_NOTES = {
    "validate": "bundled XSLT via saxonche, the library's verdict (fatal/error errors; warnings not counted)",
}

_loop = asyncio.new_event_loop()


def extract(pdf, _path=None):
    return _extract_xml_from_pdf(pdf)


def parse(xml, _path=None):
    return ZUGFeRDCIIParser().parse(xml)


def validate(xml, _path=None):
    r = _loop.run_until_complete(_iv.invoice_validate(xml_content=xml.decode("utf-8"), cloud_validate=False))
    if r.get("validator_used") != "local_schematron":
        raise RuntimeError(f"unexpected validator: {r.get('validator_used')!r} ({r})"[:200])
    rules = [f"{e['rule_id']}" for e in r["errors"]]
    return r["error_count"], rules


def write(model):
    return ZUGFeRDCIISerializer().serialize(model)


def embed(host, xml):
    return embed_xml_in_pdf(host, xml)
