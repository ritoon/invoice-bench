"""zugferd-pdf ("zugferd" module): XML/PDF/A-3U generation, extraction, PDF check.

No XML-to-model parser. validate_zugferd_pdf (zugferd/validator.py)
is a home-made check: PDF header, XMP (pdfaid:part=3 looked for as an
element only, conformance, fx extension schema), OutputIntent and
ICC profile, /AF, presence of the attachment, BT-24 in a list, BR-DEC-23
on the lines and BT-112 decimals. No schematron, despite the README.
"""

from datetime import date
from decimal import Decimal

import zugferd
from zugferd.models import Invoice, LineItem, PostalAddress, TradeParty

PKG = "zugferd-pdf"

NOTES = {
    "write": "generate_facturx_xml (totals and VAT recalculated from the lines)",
    "embed": "package_zugferd_pdfa3u: adds an sRGB OutputIntent and a PDF/A-3U XMP, creation date = now",
    "check_pdf": "validate_zugferd_pdf: home-made checks of the container (XMP, OutputIntent, /AF) + BT-24 + 2 decimal checks, no schematron",
}

UNSUPPORTED = {
    "parse": "no XML-to-model parsing",
    "validate": "no validation of a standalone XML (validate_zugferd_pdf only takes a PDF)",
}

SURVEY_NOTES = {
    "extract": "extract_zugferd_xml (pikepdf)",
}


def extract(pdf, _path=None):
    return zugferd.extract_zugferd_xml(pdf)


def hand_model():
    """EN16931_Einfach.cii.xml copied by hand into the zugferd-pdf model.

    The model has no global identifiers, no REG notes, no structured payment
    terms; the due date (missing from the XML) takes the date of the payment text.
    """
    seller = TradeParty(
        name="Lieferant GmbH",
        address=PostalAddress(postcode="80333", line_one="Lieferantenstraße 20", city_name="München", country_id="DE"),
        vat_id="DE123456789",
        tax_number="201/113/40209",
    )
    buyer = TradeParty(
        name="Kunden AG Mitte",
        address=PostalAddress(postcode="69876", line_one="Kundenstraße 15", city_name="Frankfurt", country_id="DE"),
    )
    return Invoice(
        invoice_number="471102",
        issue_date=date(2018, 3, 5),
        due_date=date(2018, 4, 4),
        currency="EUR",
        seller=seller,
        buyer=buyer,
        items=[
            LineItem(line_id="1", name="Trennblätter A4", quantity=Decimal("20.0000"), unit_code="H87",
                     net_unit_price=Decimal("9.9000"), vat_rate=Decimal("19.00")),
            LineItem(line_id="2", name="Joghurt Banane", quantity=Decimal("50.0000"), unit_code="H87",
                     net_unit_price=Decimal("5.5000"), vat_rate=Decimal("7.00")),
        ],
        notes=["Rechnung gemäß Bestellung vom 01.03.2018."],
    )


def write(model):
    return zugferd.generate_facturx_xml(model)


def embed(host, xml):
    return zugferd.package_zugferd_pdfa3u(host, xml)


def check_pdf(pdf, _path=None):
    r = zugferd.validate_zugferd_pdf(pdf)
    return len(r.errors), list(r.errors)
