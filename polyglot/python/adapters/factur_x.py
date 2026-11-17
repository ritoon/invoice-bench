"""factur-x (Akretion) : extraction, lecture en dictionnaire EN 16931, écriture CII, XSD, injection.

La validation Schematron de la librairie (xml_check_schematron) n'est pas
locale : elle envoie le XML en HTTP POST à un serveur Saxon
(http://localhost:5000/transform par défaut, facturx/facturx.py). Elle est
donc désactivée partout ici ; la validation hors ligne se limite au XSD.
"""

import facturx

PKG = "factur-x"

NOTES = {
    "extract": "check_xsd=False (par défaut l'extraction valide aussi le XML contre le XSD)",
    "parse": "parse_ubl_cii_xml, XSD incluse (défaut)",
    "validate": "XSD Factur-X 1.09.2 seulement (le schematron exige un serveur Saxon HTTP)",
    "write": "generate_cii_xml depuis le dictionnaire de parse_ubl_cii_xml, XSD incluse (défaut), schematron désactivé",
    "embed": "generate_from_binary, options par défaut (XSD incluse, passe par un fichier temporaire, XMP PDF/A-3)",
}

UNSUPPORTED = {
    "check_pdf": "pas de contrôle natif d'un PDF hybride",
}

SURVEY_NOTES = {
    "validate": "XSD Factur-X 1.09.2 seulement (xml_check_xsd)",
    "parse": "parse_ubl_cii_xml (XSD incluse)",
    "extract": "get_xml_from_pdf(check_xsd=False)",
}


def extract(pdf, _path=None):
    _name, xml = facturx.get_xml_from_pdf(pdf, check_xsd=False)
    if not xml:
        raise ValueError("aucune pièce jointe XML Factur-X/ZUGFeRD trouvée")
    return xml


def parse(xml, _path=None):
    return facturx.parse_ubl_cii_xml(xml)


def validate(xml, _path=None):
    try:
        facturx.xml_check_xsd(xml)
    except Exception as e:  # la librairie lève à la première erreur XSD
        return 1, ["XSD: " + str(e).split(". Here is the error")[0][:80]]
    return 0, []


def write(model):
    return facturx.generate_cii_xml(model, check_schematron=False)


def embed(host, xml):
    return facturx.generate_from_binary(host, xml)


def survey_extract_extra(files):
    """Combien de PDF l'extraction par défaut (check_xsd=True) garde-t-elle ?"""
    ok = 0
    for f in files:
        try:
            with open(f, "rb") as fh:
                _n, x = facturx.get_xml_from_pdf(fh.read())
            ok += bool(x)
        except Exception:
            pass
    return f"avec check_xsd=True (défaut) : {ok}/{len(files)}"
