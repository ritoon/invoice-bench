"""zugferd-pdf (module « zugferd ») : génération XML/PDF/A-3U, extraction, contrôle de PDF.

Pas de parseur XML vers modèle. validate_zugferd_pdf (zugferd/validator.py)
est un contrôle maison : en-tête PDF, XMP (pdfaid:part=3 cherché sous forme
d'élément seulement, conformance, schéma d'extension fx), OutputIntent et
profil ICC, /AF, présence de la pièce jointe, BT-24 dans une liste, BR-DEC-23
sur les lignes et décimales de BT-112. Aucun schematron, malgré le README.
"""

from datetime import date
from decimal import Decimal

import zugferd
from zugferd.models import Invoice, LineItem, PostalAddress, TradeParty

PKG = "zugferd-pdf"

NOTES = {
    "write": "generate_facturx_xml (totaux et TVA recalculés depuis les lignes)",
    "embed": "package_zugferd_pdfa3u : ajoute OutputIntent sRGB et XMP PDF/A-3U, date de création = maintenant",
    "check_pdf": "validate_zugferd_pdf : contrôles maison du conteneur (XMP, OutputIntent, /AF) + BT-24 + 2 contrôles de décimales, pas de schematron",
}

UNSUPPORTED = {
    "parse": "pas de lecture XML vers modèle",
    "validate": "pas de validation d'un XML seul (validate_zugferd_pdf ne prend qu'un PDF)",
}

SURVEY_NOTES = {
    "extract": "extract_zugferd_xml (pikepdf)",
}


def extract(pdf, _path=None):
    return zugferd.extract_zugferd_xml(pdf)


def hand_model():
    """EN16931_Einfach.cii.xml recopié à la main dans le modèle de zugferd-pdf.

    Le modèle n'a ni identifiants globaux, ni notes REG, ni conditions de paiement
    structurées ; l'échéance (absente du XML) reprend la date du texte de paiement.
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
