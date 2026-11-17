"""Un adaptateur par librairie.

Chaque module expose PKG (nom PyPI) et, pour chaque opération qu'il sait faire :

    extract(pdf: bytes, path) -> bytes | str        PDF -> XML de la facture
    parse(xml: bytes, path) -> modèle               XML -> modèle de la librairie
    validate(xml: bytes, path) -> (n, [règles])     validation complète, n = nombre d'erreurs
    write(modèle) -> bytes | str                    modèle -> XML CII
    hand_model() -> modèle                          à défaut de parseur : EN16931_Einfach à la main
    embed(host: bytes, xml: bytes) -> bytes         PDF hôte + XML -> PDF Factur-X
    check_pdf(pdf: bytes, path) -> (n, [règles])    contrôle natif d'un PDF hybride

Une opération absente vaut « unsupported » (raison dans UNSUPPORTED[op]).
NOTES[op ou benchmark] est recopié dans le champ note du JSONL.
"""
