"""erechnung : génération seule (PDF DIN 5008 dessiné par reportlab + XML CII).

L'API publique (Invoice.zugferd_file) produit le PDF complet et l'écrit sur
disque. On mesure la partie XML du même chemin (erechnung/generator.py
_zugferd_file) : calculation.calculate puis xml.xml_bytes. Pas de lecture,
pas d'extraction, pas d'injection dans un PDF hôte (le PDF est toujours
dessiné par la librairie, puis attaché avec factur-x).
"""

from datetime import date
from decimal import Decimal

from erechnung import Customer, CustomerType, Invoice, LineItem, PaymentMethod, Seller, Skonto
from erechnung.calculation import calculate
from erechnung.xml import xml_bytes

PKG = "erechnung"

NOTES = {
    "write": "erechnung.calculation.calculate + erechnung.xml.xml_bytes (l'API publique ne produit que le PDF complet)",
}

UNSUPPORTED = {
    "extract": "génération seule",
    "parse": "génération seule (pas de lecture XML)",
    "validate": "génération seule (validation des données d'entrée Python seulement)",
    "embed": "pas d'hôte PDF en entrée : le PDF est dessiné par la librairie (reportlab) puis attaché avec factur-x",
    "check_pdf": "génération seule",
}


def hand_model():
    """EN16931_Einfach.cii.xml recopié à la main dans le modèle d'erechnung.

    Le modèle impose rue et numéro séparés ; il n'a pas d'identifiants globaux
    (GTIN, GLN) ni de note REG ; le moyen de paiement (absent du XML) est fixé
    à « espèces » pour ne pas inventer d'IBAN.
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
