# Benchmarks TypeScript, Python et Rust

Ce dossier compare le package Go de ce dépôt aux librairies Factur-X / ZUGFeRD / EN 16931 des écosystèmes TypeScript, Python et Rust : temps de réponse, exactitude sur le corpus officiel ZUGFeRD, et normes réellement appliquées. Il complète [`bench/ecosystem`](../ecosystem), qui fait la même chose pour les librairies Go.

Le code de ce dossier ne fait pas partie du package `factur-x`. Chaque langage a son propre dossier, autonome :

| Dossier | Contenu |
|---|---|
| [`PROTOCOL.md`](PROTOCOL.md) | Protocole commun : fichiers d'entrée, opérations, chronométrage, format des sorties |
| [`typescript/`](typescript) | 13 librairies npm, Node 22, harnais `.mjs` ; fiches dans [`typescript/COVERAGE.md`](typescript/COVERAGE.md) |
| [`python/`](python) | 8 librairies PyPI, CPython 3.13 ; fiches dans [`python/COVERAGE.md`](python/COVERAGE.md) |
| [`rust/`](rust) | 11 crates, rustc 1.93, un binaire par crate ; fiches dans [`rust/COVERAGE.md`](rust/COVERAGE.md) |
| `report/` | Outil qui fusionne les sorties en tableaux Markdown |
| `results/<date>/` | Sorties brutes (`bench_*`, `survey_*`), environnement (`env.txt`) et tableaux (`tables.md`) |

## Méthode

Les opérations, les fichiers d'entrée et les règles de mesure sont ceux du banc Go : mêmes octets du corpus ZUGFeRD (commit `84f8d8f`), mêmes noms de benchmarks (`Extract_*`, `Parse_*`, `ParseValidate_*`, `Write_*`, `Embed_*`, `CheckPDF_*`), résultat vérifié avant toute mesure.

Le chronométrage reproduit `go test -bench . -count 3 -benchtime 1s` : N croît par paliers 1, 2, 5, 10… jusqu'à ce qu'une série dure 1 s, puis 3 séries de N appels ; la valeur retenue est la médiane. Les mesures se font en processus, sur un seul fil, avec l'horloge monotone de chaque langage. La colonne « Premier appel » donne la durée du tout premier appel : chargement paresseux, compilation des schematrons.

Pour la validation, chaque librairie applique le jeu de règles le plus complet dont elle dispose hors ligne, sans règles françaises (le corpus est allemand). Les clients d'API distantes, les librairies qui exigent Java et les paquets vides sont écartés ; la raison est donnée dans chaque `COVERAGE.md`.

## Utilisation

```bash
./run.sh setup     # corpus (via bench/ecosystem), npm ci, venv + artefacts XSLT officiels, cargo build --release
./run.sh survey    # exactitude sur le corpus, les quatre langages
./run.sh bench     # benchmarks, un langage après l'autre
./run.sh report    # results/<date>/tables.md et env.txt
./run.sh all       # les quatre à la suite
```

`survey` et `bench` acceptent une liste de langages (`./run.sh bench rust go`) ; chaque dossier a aussi son propre `run.sh`, qui accepte une liste de librairies. Variables : `OUT`, `COUNT` (3), `BENCHTIME` (1, en secondes).

Prérequis : Go 1.26, Node 22 et npm, Python 3.13 (Homebrew), Rust 1.93, et pour zugferd-rs `libxml2` et `libclang`. Les mesures durent environ 45 minutes sur un Apple M4 Pro (dont 18 pour TypeScript, à cause des validateurs lents), les relevés d'exactitude une dizaine. Lancer les mesures sur une machine au repos, sur secteur ou en mode hautes performances, capot ouvert (`run.sh` empêche la veille sur inactivité avec `caffeinate -i`, pas celle du capot fermé) : un autre processus qui occupe un cœur, ou une veille, fausse les temps.

## Résultats du 17/11/2026

Apple M4 Pro sur batterie, mode hautes performances (`powermode 2`), Go 1.26.4, Node 22.15.0, CPython 3.13.18, rustc 1.93.1. Package Go mesuré au commit `4194b98`, avant l'arrivée de l'UBL (#34, #38, #39) et de #37 sur `main` ; une nouvelle mesure du package après ces changements n'a pas pu être faite au calme (un autre processus de profilage tournait), mais ses relevés d'exactitude n'ont pas changé (6/26, 26/26, 104/104). Tableaux complets : [`results/2026-11-17/tables.md`](results/2026-11-17/tables.md).

Le Mac s'est mis en veille pendant la campagne (capot fermé, de 10 h 27 à 12 h 19) : les deux librairies TypeScript mesurées pendant ce temps, factur-x-kit et @lucaapp/zugferd-ts, ont été remesurées après coup ; toutes les autres mesures ont été prises machine éveillée.

### Temps de réponse

Le plus rapide des autres langages pour chaque opération, et en particulier parmi les validateurs qui appliquent tout le schematron EN 16931 :

| Opération | factur-x (Go) | Plus rapide ailleurs | Écart |
|---|---|---|---|
| Extraction du XML (PDF de 149 Ko) | 28,9 µs | einvoice-validator (Rust) : 32,4 µs | ×1,1 |
| Lecture XML → modèle (13 Ko) | 16,0 µs | zugferd-rs (Rust) : 23,9 µs | ×1,5 |
| Lecture XML → modèle (1,3 Mo, pièce jointe base64) | 1,14 ms | mcp-einvoicing-de (Python, lxml) : 713 µs | ×0,6 |
| Validation, toutes librairies (13 Ko) | 23,3 µs | pdf-invoice (Rust) : 25,7 µs, 14 règles, 20 faux positifs sur 26 | ×1,1 |
| Validation EN 16931 complète (13 Ko) | 23,3 µs | einvoice-validator (Rust, XSLT CEN 1.3.17) : 510 µs | ×22 |
| Validation EN 16931 complète, en TypeScript | 23,3 µs | @sweberdev/summand : 866 µs | ×37 |
| Validation EN 16931 complète, en Python | 23,3 µs | saxonche + XSLT officiels CEN et KoSIT : 903 µs | ×39 |
| Écriture modèle → XML (13 Ko) | 3,0 µs | zugferd-rs (Rust) : 6,7 µs | ×2,3 |
| Injection dans un PDF à xref classique | 79,4 µs | pdfa (Rust) : 87,6 µs, échoue sur un hôte à flux xref | ×1,1 |
| Injection dans un PDF à flux xref | 81,9 µs | pdf-invoice (Rust) : 167 µs | ×2,0 |
| Contrôle complet d'un PDF hybride (149 Ko) | 85,6 µs | zugferd-pdf (Python) : 218 µs, contrôles du conteneur seulement | ×2,6 |

- Le package Go est le plus rapide dans 13 colonnes sur 16. Les trois exceptions :
  - l'extraction du PDF MINIMUM : pdfa (Rust) cherche le premier flux `/EmbeddedFile` dans les octets sans lire la structure du PDF, ce qui marche sur ce fichier (499 ns) mais pas sur 72 des 104 PDF du corpus ;
  - la lecture du XML de 1,3 Mo : les lecteurs Python fondés sur lxml (libxml2, en C) gardent la pièce jointe base64 en texte ;
  - la validation de 19 Ko : pdf-invoice (Rust) est à égalité (39,3 µs contre 40,8 µs) avec 14 règles.
- Les validateurs qui exécutent le schematron officiel coûtent de 0,5 ms à 100 ms par facture : einvoice-validator (Rust) 510 µs, summand (TS) 866 µs, saxonche (Python) 903 µs, erechnung-validator (Rust) 34 ms, @kontor-mcp/core (TS) 48 ms, @conformo/toolkit (TS) 73 ms, mcp-einvoicing-de (Python) 104 ms, qui recompile ses XSLT à chaque appel. @lucaapp/zugferd-ts met de 2,9 à 9,3 s.
- Leur premier appel compile les feuilles de style : 26 ms (einvoice-validator) à 290 ms (mcp-einvoicing-de). Les règles du package Go sont du code compilé : il n'y a rien à charger.

### Exactitude sur le corpus

Sur les 26 XML CII valides du corpus, les schematrons officiels (CEN EN 16931 1.3.17 et KoSIT XRechnung) en signalent 6 : BR-CL-25 deux fois, CII-SR-465 quatre fois sur les exemples XRECHNUNG. CII-SR-465 est un avertissement dans les artefacts du CEN ; la configuration KoSIT l'élève au rang d'erreur. Un validateur fidèle au CEN seul en signale donc 2.

| Signalés | Librairies |
|---|---|
| 6, comme la référence | **factur-x (Go)**, @conformo/toolkit, @sweberdev/summand, en16931 (TS, sans distinguer avertissements et erreurs), saxonche + XSLT officiels, erechnung-validator ; getafix aussi, mais par 6 exceptions de lecture et non par des règles |
| 2, comme le CEN seul | @attestwire/en16931, @kontor-mcp/core, mcp-einvoicing-de, einvoice-validator, en16931 (Rust) |
| 3 à 20, faux positifs | @ekwo-ai/factur-x (3, des clés IBAN), core-invoice (7), zugferd-rs (7), @lucaapp/zugferd-ts (8), @stafyniaksacha/facturx (12), factur-x Rust (17), factur-x-kit (18), pdf-invoice (20) |
| 0 | @stackforge-eu/factur-x, @fin.cx/einvoice, factur-x (Python), drafthorse : schéma XSD seul ou règles trop partielles |

Les 104 PDF hybrides sont tous extraits par le package Go et par 15 des 19 autres librairies qui savent extraire (et par lopdf, le plancher Rust). Les autres : @stafyniaksacha/facturx 100, pdf-invoice 101, pycheval 65, pdfa 32.

### Normes appliquées

Relevé dans le code des librairies, pas dans leur README (fichiers et lignes dans chaque `COVERAGE.md`). Une sélection : les plus téléchargées et les plus complètes de chaque langage.

| Librairie | Langage | XSD | EN 16931 (CII) | Profils Factur-X | BR-FR | BR-DE (XRechnung) | Peppol BIS 3 en CII | Conteneur PDF/A-3 | UBL |
|---|---|---|---|---|---|---|---|---|---|
| **factur-x (ce dépôt)** | Go | Factur-X 1.09.2 | 793 assertions, CEN 1.3.17 | 229 règles (1.09, 1.09.2) | FNFE-MPE 1.4.0.04 | KoSIT 3.0.2 | 3.0.21 (90 assertions) | 40 contrôles | non au commit mesuré ; oui depuis #34, #38 et #39 (lecture, écriture, validation EN 16931, BR-FR, BR-DE, Peppol) |
| @stackforge-eu/factur-x (94 000/sem.) | TS | Factur-X 1.09 | non | XSD seul | non | non | non | non | non |
| @e-invoice-eu/core (63 000/sem.) | TS | non | non (génération) | génération | non | génération | non | non | écriture |
| @conformo/toolkit | TS | oui | CEN 1.3.17 (SaxonJS) | XSD seul | Flux 2, si `country: "FR"` | KoSIT 3.0.2 | UBL seulement | partiel | oui |
| @sweberdev/summand | TS | oui | CEN 1.3.16 | XMP comparé à BT-24 | non | KoSIT 2.6.0 | non | partiel | validation |
| @stafyniaksacha/facturx | TS | Factur-X 1.09 | via le schematron Factur-X 1.09 (SaxonJS) | oui, sans gravités | non | non | non | non | non |
| factur-x, Akretion (411 000/mois) | Python | Factur-X 1.09 | seulement par un serveur Saxon HTTP | XSD ; schematron par serveur | 124 règles, par serveur | non | non | non | lecture, écriture |
| drafthorse (157 000/mois) | Python | Factur-X 1.0.07 | non | XSD seul | non | non | non | non | non |
| mcp-einvoicing-de | Python | non | FeRD 1.09.2 ou CEN 1.3.16 (saxonche) | oui (saxonche) | non | KoSIT 3.0.2 | non | non | oui |
| en16931 + en16931-formats (9 700 tél.) | Rust | non | CEN 1.3.17 réécrit, 447 des 522 CII-SR et CII-DT | partiel | non | KoSIT 2.6.0 | oui, 3.0.21 | partiel | oui |
| einvoice-validator | Rust | oui | XSLT CEN 1.3.17 exécuté en Rust | seulement avec la feature `fetch` (téléchargement) | idem | idem | idem | non | validation |
| erechnung-validator | Rust | oui | CEN 1.3.16 | Factur-X 1.09.2, 5 profils | non | KoSIT 3.0.2 | non | non | validation |
| zugferd-rs | Rust | Factur-X 1.09.2 | partiel (XPath 1.0) | XSD | non | non | non | 21 règles veraPDF | oui |

Aucune librairie TypeScript, Python ou Rust ne réunit les règles EN 16931 complètes, les règles de profil Factur-X, les règles françaises, les règles allemandes, Peppol en CII et le contrôle du conteneur. Les plus proches :

- **@conformo/toolkit** (TS) : CEN, KoSIT et BR-FR, UBL, création de PDF/A-3 ; pas de schematron de profil Factur-X, Peppol en UBL seulement, dépend de SaxonJS (licence Saxonica) ; 73 ms par validation.
- **erechnung-validator** (Rust) : CEN, KoSIT et schematrons Factur-X des 5 profils ; ni BR-FR, ni Peppol, ni conteneur, ni lecture vers un modèle ; 34 ms.

Les librairies les plus téléchargées ne valident pas les règles métier : @stackforge-eu/factur-x et factur-x (Python) hors ligne s'arrêtent au XSD, @e-invoice-eu/core et node-zugferd génèrent seulement.

Ce que d'autres font et que le package Go ne fait pas : la **conversion en PDF/A-3** du PDF hôte (zugferd-rs, @lucaapp/zugferd-ts, @conformo/toolkit ; mcp-einvoicing-de ajoute un OutputIntent sRGB). La syntaxe **UBL**, que plusieurs librairies traitent (@conformo/toolkit, en16931, zugferd-rs, mcp-einvoicing-de, @lucaapp/zugferd-ts…), manquait au commit mesuré ; le package la lit, l'écrit et la valide depuis #34, #38 et #39. Le banc ne mesure que le CII.

## Particularités connues

- **pdfa (Rust)** extrait par recherche d'octets : 499 ns sur le PDF MINIMUM, mais 32 PDF sur 104 seulement.
- **erechnung-validator** apparaît dans `CheckPDF` avec `beide_pdf`, sa propre fonction, qui extrait et valide le XML sans contrôler le conteneur.
- **en16931 (Rust)** refuse d'injecter dans les deux PDF hôtes, parce que leur XMP contient déjà des propriétés `fx:`. **@stafyniaksacha/facturx** échoue sur tout hôte dont le trailer a un `/ID`.
- **Modèles construits à la main** pour `Write_Small` quand la librairie ne sait pas lire le CII : @stackforge-eu/factur-x, @e-invoice-eu/core, node-zugferd, @ekwo-ai/factur-x, zugferd-pdf, erechnung, zugferd (Rust). Leur `Write_Lines` n'est pas mesuré.
- **Un processus par benchmark** en TypeScript et en Python : « Premier appel » y est un appel à froid. En Rust, un processus par librairie.
- **Installation** : `npm ci --ignore-scripts` (le postinstall de @fin.cx/einvoice télécharge des schematrons non figés), wheels Python seulement (sauf `tagic`, pur Python, relu), `cargo build --ignore-rust-version` pour einvoice-validator, `RAYON_NUM_THREADS=1` pour lopdf.
- Le contrôle complet d'un PDF par le package Go mesure 85,6 µs le 17/11, contre 54,8 µs dans la campagne Go du 16/11 (1 326 allocations au lieu de 301) : le code de `facturx.Check` a changé entre les deux.

## Licences

Le code de ce dossier est sous licence MIT, comme le dépôt. Les dépendances sont installées par `run.sh setup` et ne sont pas importées par le package `factur-x`. Certaines ont des licences restrictives : AGPL-3.0 (pdf-invoice), EUPL-1.2 (@stackforge-eu/factur-x, artefacts CEN), licence Saxonica (SaxonJS, utilisé par @conformo, @kontor-mcp et @stafyniaksacha), MPL-2.0 (saxonche). Les artefacts officiels téléchargés par le banc Python (CEN 1.3.17, KoSIT 2.6.0) sont vérifiés par SHA-256 et ne sont pas versionnés.
