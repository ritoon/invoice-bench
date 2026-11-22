"""erechnung: generation only (DIN 5008 PDF drawn by reportlab + CII XML).

The public API (Invoice.zugferd_file) produces the full PDF and writes it to
disk. We measure the XML part of the same path (erechnung/generator.py
_zugferd_file): calculation.calculate then xml.xml_bytes. No parsing,
no extraction, no embedding into a host PDF (the PDF is always
drawn by the library, then attached with factur-x).
"""

from datetime import date
from decimal import Decimal

from erechnung import Customer, CustomerType, Invoice, LineItem, PaymentMethod, Seller, Skonto
from erechnung.calculation import calculate
from erechnung.xml import xml_bytes

PKG = "erechnung"

NOTES = {
    "write": "erechnung.calculation.calculate + erechnung.xml.xml_bytes (the public API only produces the full PDF)",
}

UNSUPPORTED = {
    "extract": "generation only",
    "parse": "generation only (no XML parsing)",
    "validate": "generation only (validation of the Python input data only)",
    "embed": "no host PDF as input: the PDF is drawn by the library (reportlab) then attached with factur-x",
    "check_pdf": "generation only",
}


def hand_model():
    """EN16931_Einfach.cii.xml copied by hand into the erechnung model.

    The model requires separate street and house number; it has no global
    identifiers (GTIN, GLN) nor REG note; the payment means (missing from the
    XML) is set to "cash" so as not to invent an IBAN.
    """
    seller = Seller(
        name="Lieferant GmbH", street="Lieferantenstraße", house_number="20", postal_code="80333",
        city="München", country="DE", tax_number="201/113/40209", vat_id="DE123456789",
        identifier="549910", registration="H A 123",
    )
    customer = Customer(
        company="Kunden AG Mitte", street="Kundenstraße", house_number="15", postal_code="69876",
        city="Frankfurt", country="DE", customer_type=CustomerType.B2B, customer_number="GE2020211",
    )
    return Invoice(
        seller=seller,
        customer=customer,
        line_items=[
            LineItem(description="Trennblätter A4", quantity=Decimal("20.0000"), unit_price=Decimal("9.9000"),
                     tax_rate=Decimal("19"), unit="H87"),
            LineItem(description="Joghurt Banane", quantity=Decimal("50.0000"), unit_price=Decimal("5.5000"),
                     tax_rate=Decimal("7"), unit="H87"),
        ],
        invoice_number="471102",
        issue_date=date(2018, 3, 5),
        service_date=date(2018, 3, 5),
        due_date=date(2018, 4, 4),
        skonto=Skonto(days=10, percent=Decimal("3")),
        note="Rechnung gemäß Bestellung vom 01.03.2018.",
        payment_method=PaymentMethod.CASH,
    )


def write(model):
    return xml_bytes(model, calculate(model))
