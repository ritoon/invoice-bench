<!-- Généré par bench/polyglot/report à partir de 2026-11-17 -->

#### Exactitude sur le corpus

Référence : les schematrons officiels (CEN EN 16931 et KoSIT) signalent 6 des 26 XML « valides » du corpus (BR-CL-25 deux fois, CII-SR-465 quatre fois).

| Librairie | Langage | XML valides signalés | XML lus | PDF extraits |
|---|---|---|---|---|
| **factur-x (ce dépôt)** | Go | 6/26 | 26/26 | 104/104 |
| @attestwire/en16931 | TypeScript | 2/26 | 26/26 | 104/104 |
| @conformo/toolkit | TypeScript | 6/26 | 26/26 | 104/104 |
| @ekwo-ai/factur-x | TypeScript | 3/26 | 26/26 | 104/104 |
| @fin.cx/einvoice | TypeScript | 0/26 | 26/26 | 104/104 |
| @kontor-mcp/core | TypeScript | 2/26 | 26/26 | 104/104 |
| @lucaapp/zugferd-ts | TypeScript | 8/26 | 19/26 | 104/104 |
| @stackforge-eu/factur-x | TypeScript | 0/26 | — | 104/104 |
| @stafyniaksacha/facturx | TypeScript | 12/26 | 26/26 | 100/104 |
| @sweberdev/summand | TypeScript | 6/26 | — | 104/104 |
| en16931 | TypeScript | 6/26 | — | — |
| factur-x-kit | TypeScript | 18/26 | 13/26 | — |
| drafthorse | Python | 0/26 | 26/26 | — |
| factur-x | Python | 0/26 | 26/26 | 104/104 |
| getafix | Python | 6/26 | 20/26 | 104/104 |
| mcp-einvoicing-de | Python | 2/26 | 20/26 | 104/104 |
| pycheval | Python | — | 7/26 | 65/104 |
| saxonche-schematron | Python | 6/26 | — | — |
| zugferd-pdf | Python | — | — | 104/104 |
| core-invoice | Rust | 7/26 | 23/26 | — |
| einvoice | Rust | — | 26/26 | — |
| einvoice-validator | Rust | 2/26 | — | 104/104 |
| en16931 | Rust | 2/26 | 26/26 | 104/104 |
| erechnung-validator | Rust | 6/26 | — | — |
| factur-x | Rust | 17/26 | 9/26 | — |
| lopdf | Rust | — | — | 104/104 |
| pdf-invoice | Rust | 20/26 | 26/26 | 101/104 |
| pdfa | Rust | — | — | 32/104 |
| zugferd-rs | Rust | 7/26 | 19/26 | 104/104 |

#### Extraction

| Librairie | Langage | EN 16931, 149 Ko | EN 16931, 1,1 Mo | MINIMUM FNFE, 89 Ko | EXTENDED ZF 2.0, 157 Ko | Premier appel |
|---|---|---|---|---|---|---|
| **factur-x (ce dépôt)** | Go | **28,9 µs** | **34,5 µs** | 10,0 µs | **33,4 µs** | — |
| einvoice-validator | Rust | 32,4 µs (×1,1) | 41,9 µs (×1,2) | 23,8 µs (×2,4) | 38,2 µs (×1,1) | 1,62 ms |
| lopdf | Rust | 71,6 µs (×2,5) | 94,5 µs (×2,7) | 42,9 µs (×4,3) | 65,0 µs (×1,9) | 678 µs |
| pdfa | Rust | 77,3 µs (×2,7) | échec | **499 ns** (×0,1) | échec | 154 µs |
| pdf-invoice | Rust | 79,5 µs (×2,7) | 105 µs (×3,0) | 48,5 µs (×4,9) | 74,6 µs (×2,2) | 682 µs |
| mcp-einvoicing-de | Python | 81,9 µs (×2,8) | 117 µs (×3,4) | 53,5 µs (×5,4) | 237 µs (×7,1) | 432 µs |
| zugferd-pdf | Python | 82,0 µs (×2,8) | 120 µs (×3,5) | 53,2 µs (×5,3) | 96,0 µs (×2,9) | 1,76 ms |
| zugferd-rs | Rust | 91,6 µs (×3,2) | 116 µs (×3,4) | 54,1 µs (×5,4) | 84,4 µs (×2,5) | 1,36 ms |
| @attestwire/en16931 | TypeScript | 182 µs (×6,3) | 205 µs (×6,0) | 111 µs (×11) | 215 µs (×6,5) | 4,43 ms |
| pycheval | Python | 205 µs (×7,1) | échec | 182 µs (×18) | échec | 1,02 ms |
| getafix | Python | 230 µs (×8,0) | 388 µs (×11) | 212 µs (×21) | 215 µs (×6,4) | 1,01 ms |
| @lucaapp/zugferd-ts | TypeScript | 251 µs (×8,7) | 305 µs (×8,8) | 124 µs (×12) | 235 µs (×7,1) | 4,64 ms |
| factur-x | Python | 252 µs (×8,7) | 390 µs (×11) | 202 µs (×20) | 280 µs (×8,4) | 2,47 ms |
| @ekwo-ai/factur-x | TypeScript | 276 µs (×9,5) | 317 µs (×9,2) | 141 µs (×14) | 255 µs (×7,6) | 4,98 ms |
| @stackforge-eu/factur-x | TypeScript | 288 µs (×10,0) | 335 µs (×9,7) | 145 µs (×15) | 251 µs (×7,5) | 4,72 ms |
| @conformo/toolkit | TypeScript | 288 µs (×10,0) | 10,45 ms (×303) | 139 µs (×14) | 261 µs (×7,8) | 4,58 ms |
| @stafyniaksacha/facturx | TypeScript | 375 µs (×13) | 3,12 ms (×90) | 173 µs (×17) | 437 µs (×13) | 7,50 ms |
| en16931 | Rust | 392 µs (×14) | 718 µs (×21) | 254 µs (×26) | 291 µs (×8,7) | 1,48 ms |
| @fin.cx/einvoice | TypeScript | 412 µs (×14) | 511 µs (×15) | 142 µs (×14) | 469 µs (×14) | 47,56 ms |
| @kontor-mcp/core | TypeScript | 743 µs (×26) | 1,53 ms (×44) | 704 µs (×71) | 820 µs (×25) | 9,18 ms |
| @sweberdev/summand | TypeScript | 946 µs (×33) | 20,93 ms (×607) | 456 µs (×46) | 1,13 ms (×34) | 3,23 ms |

En gras : le plus rapide de la colonne. (×n) : temps rapporté à celui du package Go de ce dépôt. Tri sur la colonne « EN 16931, 149 Ko » ; « Premier appel » : durée du premier appel dans le processus, même colonne (non mesurée pour Go).

#### Lecture XML → modèle

| Librairie | Langage | 13 Ko | 19 Ko | 1,3 Mo | Premier appel |
|---|---|---|---|---|---|
| **factur-x (ce dépôt)** | Go | **16,0 µs** | **28,2 µs** | 1,14 ms | — |
| zugferd-rs | Rust | 23,9 µs (×1,5) | 37,0 µs (×1,3) | échec | 1,11 ms |
| pdf-invoice | Rust | 25,5 µs (×1,6) | 38,8 µs (×1,4) | 1,63 ms (×1,4) | 1,06 ms |
| core-invoice | Rust | 45,1 µs (×2,8) | 69,3 µs (×2,5) | 3,65 ms (×3,2) | 283 µs |
| @ekwo-ai/factur-x | TypeScript | 73,6 µs (×4,6) | 116 µs (×4,1) | 26,84 ms (×24) | 1,79 ms |
| en16931 | Rust | 76,5 µs (×4,8) | 131 µs (×4,6) | 7,28 ms (×6,4) | 2,34 ms |
| @conformo/toolkit | TypeScript | 93,4 µs (×5,8) | 169 µs (×6,0) | 1,81 ms (×1,6) | 1,57 ms |
| @attestwire/en16931 | TypeScript | 137 µs (×8,5) | 245 µs (×8,7) | 3,46 ms (×3,0) | 2,18 ms |
| @lucaapp/zugferd-ts | TypeScript | 162 µs (×10) | 294 µs (×10) | échec | 2,48 ms |
| pycheval | Python | 229 µs (×14) | échec | échec | 1,06 ms |
| factur-x-kit | TypeScript | 272 µs (×17) | échec | échec | 8,72 ms |
| einvoice | Rust | 321 µs (×20) | 547 µs (×19) | 11,89 ms (×10) | 2,00 ms |
| mcp-einvoicing-de | Python | 341 µs (×21) | 491 µs (×17) | **713 µs** (×0,6) | 1,68 ms |
| drafthorse | Python | 597 µs (×37) | 1,00 ms (×35) | 1,10 ms (×1,0) | 3,24 ms |
| @stafyniaksacha/facturx | TypeScript | 968 µs (×60) | 1,69 ms (×60) | 5,56 ms (×4,9) | 7,99 ms |
| @kontor-mcp/core | TypeScript | 1,18 ms (×74) | 1,96 ms (×69) | 4,69 ms (×4,1) | 13,91 ms |
| getafix | Python | 1,24 ms (×77) | 2,14 ms (×76) | échec | 2,87 ms |
| factur-x | Python | 1,47 ms (×92) | 2,07 ms (×73) | 12,08 ms (×11) | 7,48 ms |
| @fin.cx/einvoice | TypeScript | 2,19 ms (×137) | 4,01 ms (×142) | 7,34 ms (×6,5) | 13,06 ms |
| factur-x | Rust | échec | échec | 1,63 ms (×1,4) | — |

En gras : le plus rapide de la colonne. (×n) : temps rapporté à celui du package Go de ce dépôt. Tri sur la colonne « 13 Ko » ; « Premier appel » : durée du premier appel dans le processus, même colonne (non mesurée pour Go).

#### Lecture et validation

| Librairie | Langage | Valides signalés /26 | 13 Ko | 19 Ko | 1,3 Mo | Premier appel |
|---|---|---|---|---|---|---|
| **factur-x (ce dépôt)** | Go | 6/26 | **23,3 µs** | 40,8 µs | **755 µs** | — |
| pdf-invoice | Rust | 20/26 | 25,7 µs (×1,1) | **39,3 µs** (×1,0) | 1,68 ms (×2,2) | 53,4 µs |
| @ekwo-ai/factur-x | TypeScript | 3/26 | 71,1 µs (×3,1) | 116 µs (×2,9) | 21,20 ms (×28) | 1,67 ms |
| core-invoice | Rust | 7/26 | 77,4 µs (×3,3) | 120 µs (×2,9) | 5,42 ms (×7,2) | 739 µs |
| en16931 | Rust | 2/26 | 80,0 µs (×3,4) | 130 µs (×3,2) | 6,93 ms (×9,2) | 1,40 ms |
| @attestwire/en16931 | TypeScript | 2/26 | 203 µs (×8,7) | 349 µs (×8,6) | 7,04 ms (×9,3) | 8,72 ms |
| factur-x-kit | TypeScript | 18/26 | 284 µs (×12) | échec | échec | 10,67 ms |
| drafthorse | Python | 0/26 | 510 µs (×22) | 548 µs (×13) | 8,38 ms (×11) | 3,27 ms |
| einvoice-validator | Rust | 2/26 | 510 µs (×22) | 839 µs (×21) | 4,34 ms (×5,7) | 25,89 ms |
| factur-x | Python | 0/26 | 518 µs (×22) | 548 µs (×13) | 5,41 ms (×7,2) | 824 µs |
| @stackforge-eu/factur-x | TypeScript | 0/26 | 819 µs (×35) | 891 µs (×22) | 10,18 ms (×13) | 34,15 ms |
| @sweberdev/summand | TypeScript | 6/26 | 866 µs (×37) | 1,45 ms (×36) | 7,83 ms (×10) | 16,74 ms |
| saxonche-schematron | Python | 6/26 | 903 µs (×39) | 1,48 ms (×36) | 13,34 ms (×18) | 257,28 ms |
| getafix | Python | 6/26 | 1,59 ms (×68) | 2,66 ms (×65) | échec | 3,13 ms |
| zugferd-rs | Rust | 7/26 | 2,02 ms (×86) | 2,79 ms (×69) | 5,21 ms (×6,9) | 5,07 ms |
| @fin.cx/einvoice | TypeScript | 0/26 | 3,08 ms (×132) | 5,52 ms (×135) | 10,99 ms (×15) | 17,67 ms |
| en16931 | TypeScript | 6/26 | 15,06 ms (×646) | 27,64 ms (×678) | 41,81 ms (×55) | 170,62 ms |
| erechnung-validator | Rust | 6/26 | 34,05 ms (×1 461) | 56,30 ms (×1 380) | 84,82 ms (×112) | 170,79 ms |
| @kontor-mcp/core | TypeScript | 2/26 | 47,74 ms (×2 048) | 54,14 ms (×1 328) | 74,70 ms (×99) | 144,50 ms |
| @stafyniaksacha/facturx | TypeScript | 12/26 | 60,07 ms (×2 577) | 77,06 ms (×1 890) | 78,74 ms (×104) | 133,37 ms |
| @conformo/toolkit | TypeScript | 6/26 | 72,53 ms (×3 111) | 80,08 ms (×1 964) | 126,70 ms (×168) | 116,81 ms |
| mcp-einvoicing-de | Python | 2/26 | 104,01 ms (×4 462) | 106,26 ms (×2 606) | 168,66 ms (×223) | 288,17 ms |
| @lucaapp/zugferd-ts | TypeScript | 8/26 | 2,85 s (×122 445) | 9,34 s (×229 106) | 4,35 s (×5 764) | 3,00 s |
| factur-x | Rust | 17/26 | échec | échec | 1,64 ms (×2,2) | — |

En gras : le plus rapide de la colonne. (×n) : temps rapporté à celui du package Go de ce dépôt. Tri sur la colonne « 13 Ko » ; « Premier appel » : durée du premier appel dans le processus, même colonne (non mesurée pour Go).

#### Écriture modèle → XML

| Librairie | Langage | 13 Ko | 19 Ko | Premier appel |
|---|---|---|---|---|
| **factur-x (ce dépôt)** | Go | **3,0 µs** | **4,8 µs** | — |
| zugferd-rs | Rust | 6,7 µs (×2,3) | 11,6 µs (×2,4) | 401 µs |
| core-invoice | Rust | 10,8 µs (×3,7) | 17,1 µs (×3,6) | 1,58 ms |
| @ekwo-ai/factur-x | TypeScript | 11,0 µs (×3,7) | — | 450 µs |
| @stackforge-eu/factur-x | TypeScript | 14,5 µs (×4,9) | — | 536 µs |
| @conformo/toolkit | TypeScript | 22,4 µs (×7,6) | 38,5 µs (×8,0) | 1,66 ms |
| zugferd | Rust | 28,4 µs (×9,6) | — | 240 µs |
| einvoice | Rust | 36,1 µs (×12) | 62,8 µs (×13) | 1,06 ms |
| @attestwire/en16931 | TypeScript | 43,2 µs (×15) | 82,5 µs (×17) | 1,02 ms |
| en16931 | Rust | 47,0 µs (×16) | 81,3 µs (×17) | 1,06 ms |
| @stafyniaksacha/facturx | TypeScript | 74,9 µs (×25) | 122 µs (×25) | 1,55 ms |
| pycheval | Python | 78,7 µs (×27) | échec | 151 µs |
| @lucaapp/zugferd-ts | TypeScript | 111 µs (×37) | 200 µs (×42) | 9,40 ms |
| mcp-einvoicing-de | Python | 111 µs (×37) | 189 µs (×39) | 178 µs |
| zugferd-pdf | Python | 124 µs (×42) | — | 297 µs |
| erechnung | Python | 130 µs (×44) | — | 513 µs |
| factur-x-kit | TypeScript | 172 µs (×58) | échec | 5,79 ms |
| @fin.cx/einvoice | TypeScript | 352 µs (×119) | 490 µs (×102) | 3,47 ms |
| getafix | Python | 407 µs (×138) | 711 µs (×148) | 895 µs |
| factur-x | Python | 718 µs (×243) | 861 µs (×180) | 1,57 ms |
| node-zugferd | TypeScript | 921 µs (×312) | — | 5,90 ms |
| drafthorse | Python | 1,00 ms (×340) | 1,34 ms (×279) | 1,50 ms |
| @e-invoice-eu/core | TypeScript | 39,37 ms (×13 324) | — | 81,43 ms |
| factur-x | Rust | échec | échec | — |
| pdf-invoice | Rust | échec | échec | — |

En gras : le plus rapide de la colonne. (×n) : temps rapporté à celui du package Go de ce dépôt. Tri sur la colonne « 13 Ko » ; « Premier appel » : durée du premier appel dans le processus, même colonne (non mesurée pour Go).

#### Injection dans un PDF

| Librairie | Langage | hôte à xref classique | hôte à flux xref | Premier appel |
|---|---|---|---|---|
| **factur-x (ce dépôt)** | Go | **79,4 µs** | **81,9 µs** | — |
| pdfa | Rust | 87,6 µs (×1,1) | échec | 137 µs |
| pdf-invoice | Rust | 155 µs (×2,0) | 167 µs (×2,0) | 180 µs |
| zugferd-rs | Rust | 200 µs (×2,5) | 216 µs (×2,6) | 205 µs |
| mcp-einvoicing-de | Python | 780 µs (×9,8) | 679 µs (×8,3) | 1,53 ms |
| zugferd-pdf | Python | 890 µs (×11) | 797 µs (×9,7) | 5,29 ms |
| @stackforge-eu/factur-x | TypeScript | 1,14 ms (×14) | 1,15 ms (×14) | 15,20 ms |
| @fin.cx/einvoice | TypeScript | 1,93 ms (×24) | 1,94 ms (×24) | 52,51 ms |
| drafthorse | Python | 2,12 ms (×27) | 2,12 ms (×26) | 3,32 ms |
| getafix | Python | 2,16 ms (×27) | 2,31 ms (×28) | 2,96 ms |
| @lucaapp/zugferd-ts | TypeScript | 2,21 ms (×28) | 2,29 ms (×28) | 12,84 ms |
| @ekwo-ai/factur-x | TypeScript | 2,22 ms (×28) | 2,22 ms (×27) | 13,56 ms |
| @conformo/toolkit | TypeScript | 2,26 ms (×28) | 2,56 ms (×31) | 16,86 ms |
| pycheval | Python | 2,57 ms (×32) | 2,65 ms (×32) | 3,52 ms |
| factur-x | Python | 4,29 ms (×54) | 4,01 ms (×49) | 6,38 ms |
| @stafyniaksacha/facturx | TypeScript | échec | échec | — |
| en16931 | Rust | échec | échec | — |

En gras : le plus rapide de la colonne. (×n) : temps rapporté à celui du package Go de ce dépôt. Tri sur la colonne « hôte à xref classique » ; « Premier appel » : durée du premier appel dans le processus, même colonne (non mesurée pour Go).

#### Contrôle d'un PDF hybride

| Librairie | Langage | 149 Ko | 1,1 Mo | Premier appel |
|---|---|---|---|---|
| **factur-x (ce dépôt)** | Go | **85,6 µs** | **106 µs** | — |
| zugferd-pdf | Python | 218 µs (×2,6) | 280 µs (×2,6) | 517 µs |
| @attestwire/en16931 | TypeScript | 376 µs (×4,4) | 421 µs (×4,0) | 11,79 ms |
| en16931 | Rust | 482 µs (×5,6) | 818 µs (×7,7) | 497 µs |
| @sweberdev/summand | TypeScript | 1,95 ms (×23) | 22,49 ms (×211) | 19,29 ms |
| erechnung-validator | Rust | 32,10 ms (×375) | 36,28 ms (×341) | 32,94 ms |
| @conformo/toolkit | TypeScript | 76,88 ms (×898) | 102,65 ms (×965) | 128,97 ms |

En gras : le plus rapide de la colonne. (×n) : temps rapporté à celui du package Go de ce dépôt. Tri sur la colonne « 149 Ko » ; « Premier appel » : durée du premier appel dans le processus, même colonne (non mesurée pour Go).

#### Notes des harnais

Appels mesurés et particularités, pour les opérations réussies ou en échec (les raisons des opérations non prises en charge sont dans les fichiers `COVERAGE.md`).

- **@attestwire/en16931** (TypeScript) : extractFacturX [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; parseCiiInvoice [Parse_Small, Parse_Lines, Parse_Large] ; parseCiiInvoice hors mesure, puis generateCii(invoice) [Write_Small, Write_Lines] ; validate(octets du PDF) : extraction + règles sur le XML + constats du conteneur (AW-PDF-*, avertissements seulement) ; erreurs de sévérité fatal ; résultat : 0 erreur(s) [CheckPDF_Small, CheckPDF_Large] ; validate(octets) : règles EN 16931 + CIUS XRechnung (profil lu dans BT-24) + codes BR-CL, réécrites en TypeScript (pas de schematron, pas de XSD) ; erreurs de sévérité fatal ; résultat : 0 erreur(s) [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large]
- **@conformo/toolkit** (TypeScript) : @conformo/parse detectXmlFormat (lecture CII/UBL vers le modèle EN 16931 de @conformo/core) [Parse_Small, Parse_Lines, Parse_Large] ; @conformo/parse extractEmbeddedXml [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; @conformo/toolkit validate(octets du PDF) : extraction + XSD + schematrons + notices du conteneur (XMP, nom de fichier ; n’affectent pas valid) ; résultat : 0 erreur(s) [CheckPDF_Small, CheckPDF_Large] ; @conformo/validate validate(xml) (défauts) : XSD CII D16B (xmllint-wasm) + schematron CEN EN 16931 CII 1.3.17, et KoSIT XRechnung 3.0.2 si BT-24 le déclare, compilés en SEF et exécutés par SaxonJS ; sévérités error/fatal après les re-niveaux de la configuration KoSIT ; BR-FR non activé (country non fourni) ; résultat : 0 erreur(s) [ParseValidate_Small, ParseValidate_Lines] ; @conformo/validate validate(xml) (défauts) : XSD CII D16B (xmllint-wasm) + schematron CEN EN 16931 CII 1.3.17, et KoSIT XRechnung 3.0.2 si BT-24 le déclare, compilés en SEF et exécutés par SaxonJS ; sévérités error/fatal après les re-niveaux de la configuration KoSIT ; BR-FR non activé (country non fourni) ; résultat : 1 erreur(s) (CII-SR-465) [ParseValidate_Large] ; detectXmlFormat hors mesure, puis @conformo/toolkit serialize(invoice, "cii") [Write_Small, Write_Lines] ; recette de createPdf appliquée à l’hôte : PDFDocument.load (pdf-lib) + @conformo/pdf finalizePDFA (profil en16931, ICC sRGB du toolkit, date fixe) + save({ useObjectStreams: false }) [Embed_Small, EmbedClassic_Small]
- **@e-invoice-eu/core** (TypeScript) : new InvoiceService(logger).generate(json façon UBL de la librairie, { format: "CII", lang: "de-de" }), contrôle ajv du JSON inclus ; modèle construit à la main [Write_Small]
- **@ekwo-ai/factur-x** (TypeScript) : @ekwo-ai/factur-x/pdf embedFacturX(pdf, xml, { profile: "en16931" }) : pdf-lib, pièce jointe + XMP [Embed_Small, EmbedClassic_Small] ; @ekwo-ai/factur-x/pdf extractFacturX [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; generateCiiXml(invoice, { profile: "en16931" }) ; le modèle lu par readCii n’est pas celui de l’écriture, d’où un modèle Invoice à la main (une seule note, pas de GTIN d’article) ; modèle construit à la main [Write_Small] ; readCii(xml) [Parse_Small, Parse_Lines, Parse_Large] ; readCii(xml).violations : anomalies de lecture (dates, IBAN…) + 8 contrôles arithmétiques BR-CO-10..17 (checkReceived) ; pas de XSD ni de schematron ; résultat : 0 erreur(s) [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large]
- **@fin.cx/einvoice** (TypeScript) : EInvoice.fromXml [Parse_Small, Parse_Lines, Parse_Large] ; EInvoice.fromXml hors mesure, puis exportXml("facturx") [Write_Small, Write_Lines] ; EInvoice.fromXml puis validate(BUSINESS, { featureFlags: [EN16931_BUSINESS_RULES, CODE_LIST_VALIDATION] }) : règles EN 16931 / XRechnung / Factur-X réécrites en TypeScript (pas de XSD ; la couche schematron dépend de fichiers téléchargés au postinstall et n’est pas utilisée par validate) ; résultat : 0 erreur(s) [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large] ; new PDFEmbedder().embedXml(pdf, xml, "factur-x.xml", …) : pdf-lib attach + save, sans XMP ni métadonnées PDF/A [Embed_Small, EmbedClassic_Small] ; new PDFExtractor().extractXml (une instance réutilisée ; ses console.log sont neutralisés) [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended]
- **@kontor-mcp/core** (TypeScript) : extractEmbeddedXml [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; parseInvoice(octets) [Parse_Small, Parse_Lines, Parse_Large] ; validateInvoice(octets) (défauts) : XSD CII D16B (xmllint-wasm) + schematron CEN EN 16931 CII 1.3.16, plus KoSIT XRechnung 3.0.2 (schematron 2.5.0) si BT-24 le déclare, compilés en SEF et exécutés par SaxonJS ; erreurs error/fatal hors couche « plausibility » (règles maison KONTOR-PLAUS-*, que la librairie exclut elle-même de valid) ; résultat : 0 erreur(s) [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large]
- **@lucaapp/zugferd-ts** (TypeScript) : embedXml(pdf, xml) (défauts : factur-x.xml, surcouche PDF/A-3b) [Embed_Small, EmbedClassic_Small] ; extractXml [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; parseCiiXml [Parse_Small, Parse_Lines] ; parseCiiXml ; CII parse: unknown ExchangedDocument.TypeCode "877" [Parse_Large] ; parseCiiXml puis detectProfile hors mesure, serializeToXml(modèle, profil) [Write_Small, Write_Lines] ; validateXml (défauts) : XSD Factur-X 1.08 du profil détecté (xmllint-wasm) + schematron CEN EN 16931 CII 1.3.16 (node-schematron, JS pur), issues de sévérité error ; la librairie classe toute assertion en échec en error (les flag="warning" du CEN ne sont pas lus) ; résultat : 0 erreur(s) [ParseValidate_Small] ; validateXml (défauts) : XSD Factur-X 1.08 du profil détecté (xmllint-wasm) + schematron CEN EN 16931 CII 1.3.16 (node-schematron, JS pur), issues de sévérité error ; la librairie classe toute assertion en échec en error (les flag="warning" du CEN ne sont pas lus) ; résultat : 1 erreur(s) (BR-S-08) [ParseValidate_Lines] ; validateXml (défauts) : XSD Factur-X 1.08 du profil détecté (xmllint-wasm) + schematron CEN EN 16931 CII 1.3.16 (node-schematron, JS pur), issues de sévérité error ; la librairie classe toute assertion en échec en error (les flag="warning" du CEN ne sont pas lus) ; résultat : 3 erreur(s) (CII-SR-465, CII-SR-475, CII-SR-476) [ParseValidate_Large]
- **@stackforge-eu/factur-x** (TypeScript) : buildXml(saisie FacturXInvoiceInput, Profile.EN16931), sans lecteur CII ; modèle construit à la main [Write_Small] ; embedFacturX({ pdf, xml, profile: EN16931 }) : pdf-lib, pièce jointe AFRelationship + XMP PDF/A-3b [Embed_Small, EmbedClassic_Small] ; extractXml [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; validateXsd(xml, Profile.EN16931) : XSD Factur-X EN16931 seul (libxml2-wasm ; schéma relu et compilé à chaque appel par la librairie) ; pas de schematron EN 16931 ; résultat : 0 erreur(s) [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large]
- **@stafyniaksacha/facturx** (TypeScript) : check({ xml, level: "en16931", schematron: true }) : XSD Factur-X 1.09 EN16931 (libxml2-wasm) + schematron officiel Factur-X 1.09 EN16931 compilé en SEF, exécuté par SaxonJS ; le résultat ne porte pas de sévérité, toute assertion en échec compte ; niveau forcé à en16931 car la détection automatique lève « Unknown level: kosit » sur les XRECHNUNG ; résultat : 0 erreur(s) [ParseValidate_Small] ; check({ xml, level: "en16931", schematron: true }) : XSD Factur-X 1.09 EN16931 (libxml2-wasm) + schematron officiel Factur-X 1.09 EN16931 compilé en SEF, exécuté par SaxonJS ; le résultat ne porte pas de sévérité, toute assertion en échec compte ; niveau forcé à en16931 car la détection automatique lève « Unknown level: kosit » sur les XRECHNUNG ; résultat : 2 erreur(s) (CII-SR-465, codedb ID) [ParseValidate_Large] ; check({ xml, level: "en16931", schematron: true }) : XSD Factur-X 1.09 EN16931 (libxml2-wasm) + schematron officiel Factur-X 1.09 EN16931 compilé en SEF, exécuté par SaxonJS ; le résultat ne porte pas de sévérité, toute assertion en échec compte ; niveau forcé à en16931 car la détection automatique lève « Unknown level: kosit » sur les XRECHNUNG ; résultat : 4 erreur(s) (PEPPOL-EN16931-R008) [ParseValidate_Lines] ; extract({ pdf }) [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; generate({ pdf, xml }) (options par défaut) ; Not implemented yet: Document ID already set [Embed_Small, EmbedClassic_Small] ; xmlToInvoice [Parse_Small, Parse_Lines, Parse_Large] ; xmlToInvoice hors mesure, puis invoiceToXml + toString (le XmlDocument libxml2-wasm est libéré après chaque appel) [Write_Small, Write_Lines]
- **@sweberdev/summand** (TypeScript) : extractInvoiceXml [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; validateInvoice(octets du PDF) : extraction, comparaison du profil XMP, XSD et schematrons sur le XML ; résultat : 0 erreur(s) [CheckPDF_Small, CheckPDF_Large] ; validateInvoice(octets) (défauts) : XSD CII D16B (validateur XSD maison) + schematron CEN EN 16931 CII 1.3.16, et KoSIT XRechnung 2.6.0 si BT-24 le déclare, avec les sévérités de la configuration KoSIT ; moteur XPath 2.0 maison, sans Java ; résultat : 0 erreur(s) [ParseValidate_Small, ParseValidate_Lines] ; validateInvoice(octets) (défauts) : XSD CII D16B (validateur XSD maison) + schematron CEN EN 16931 CII 1.3.16, et KoSIT XRechnung 2.6.0 si BT-24 le déclare, avec les sévérités de la configuration KoSIT ; moteur XPath 2.0 maison, sans Java ; résultat : 1 erreur(s) (CII-SR-465) [ParseValidate_Large]
- **en16931** (TypeScript) : new SchematronRunner(loadCiiSchematron()) une fois (au premier appel), puis runner.validate(xml) : schematron CEN EN 16931 CII 1.3.16 seul (pas de XSD, pas de XRechnung) ; le résultat ne distingue pas les sévérités, toute assertion déclenchée compte ; résultat : 0 erreur(s) [ParseValidate_Small, ParseValidate_Lines] ; new SchematronRunner(loadCiiSchematron()) une fois (au premier appel), puis runner.validate(xml) : schematron CEN EN 16931 CII 1.3.16 seul (pas de XSD, pas de XRechnung) ; le résultat ne distingue pas les sévérités, toute assertion déclenchée compte ; résultat : 3 erreur(s) (CII-SR-465, CII-SR-475, CII-SR-476) [ParseValidate_Large]
- **factur-x-kit** (TypeScript) : FacturX.fromXML hors mesure, puis getXML() [Write_Small] ; FacturX.fromXML hors mesure, puis getXML() ; Received Invalid XML Input [Write_Lines] ; FacturX.fromXML(octets) [Parse_Small] ; FacturX.fromXML(octets) ; Profile not yet implemented: urn:cen.eu:en16931:2017#compliant#urn:xeinkauf.de:kosit:xrechnung_3.0 [Parse_Large] ; FacturX.fromXML(octets) ; Received Invalid XML Input [Parse_Lines] ; FacturX.fromXML(octets) puis validate() : schéma zod du profil détecté et règles BR-* réécrites en TypeScript (MINIMUM à EN16931) ; pas de XSD ni de schematron ; Profile not yet implemented: urn:cen.eu:en16931:2017#compliant#urn:xeinkauf.de:kosit:xrechnung_3.0 [ParseValidate_Large] ; FacturX.fromXML(octets) puis validate() : schéma zod du profil détecté et règles BR-* réécrites en TypeScript (MINIMUM à EN16931) ; pas de XSD ni de schematron ; Received Invalid XML Input [ParseValidate_Lines] ; FacturX.fromXML(octets) puis validate() : schéma zod du profil détecté et règles BR-* réécrites en TypeScript (MINIMUM à EN16931) ; pas de XSD ni de schematron ; résultat : 0 erreur(s) [ParseValidate_Small]
- **node-zugferd** (TypeScript) : zugferd({ profile: EN16931, strict: false }).create(saisie) hors mesure, puis toXML() ; strict: false car la validation XSD de la librairie appelle Java ; modèle construit à la main [Write_Small]
- **drafthorse** (Python) : Document.serialize(schema='FACTUR-X_EN16931'), revalide la sortie contre le XSD [Write_Small, Write_Lines] ; XSD Factur-X 1.0.07 EN16931 (drafthorse.utils.validate_xml, arrêt à la première erreur) [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large] ; attach_xml, options par défaut (profil détecté, XMP PDF/A-3) ; relu avec factur-x [Embed_Small, EmbedClassic_Small]
- **erechnung** (Python) : erechnung.calculation.calculate + erechnung.xml.xml_bytes (l'API publique ne produit que le PDF complet) ; modèle construit à la main [Write_Small]
- **factur-x** (Python) : XSD Factur-X 1.09.2 seulement (le schematron exige un serveur Saxon HTTP) [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large] ; check_xsd=False (par défaut l'extraction valide aussi le XML contre le XSD) [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; generate_cii_xml depuis le dictionnaire de parse_ubl_cii_xml, XSD incluse (défaut), schematron désactivé [Write_Small, Write_Lines] ; generate_from_binary, options par défaut (XSD incluse, passe par un fichier temporaire, XMP PDF/A-3) [Embed_Small, EmbedClassic_Small] ; parse_ubl_cii_xml, XSD incluse (défaut) [Parse_Small, Parse_Lines, Parse_Large]
- **getafix** (Python) : Document.from_xml(xml.etree.ElementTree.fromstring) [Parse_Small, Parse_Lines] ; Document.from_xml(xml.etree.ElementTree.fromstring) ; ValueError: 'urn:cen.eu:en16931:2017#compliant#urn:xeinkauf.de:kosit:xrechnung_3.0' is not a valid Profile [Parse_Large] ; Document.to_xml().render() depuis le modèle lu [Write_Small, Write_Lines] ; attach_xml(bytes, bytes) : pièce jointe générique, pas de XMP Factur-X ni de PDF/A-3 (documenté) [Embed_Small, EmbedClassic_Small] ; extract_xml(Path) : l'API ne prend qu'un chemin, la lecture du fichier (cache OS) est dans la mesure [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; lecture + règles de getafix (89 des 223 règles BR du CEN, sans BR-CL, + 16 BR-FXEXT) ; avertissements non comptés [ParseValidate_Small, ParseValidate_Lines] ; lecture + règles de getafix (89 des 223 règles BR du CEN, sans BR-CL, + 16 BR-FXEXT) ; avertissements non comptés ; ValueError: 'urn:cen.eu:en16931:2017#compliant#urn:xeinkauf.de:kosit:xrechnung_3.0' is not a valid Profile [ParseValidate_Large]
- **mcp-einvoicing-de** (Python) : ZUGFeRDCIIParser (modèle pydantic EN 16931) [Parse_Small, Parse_Lines, Parse_Large] ; ZUGFeRDCIISerializer depuis le modèle lu (modèle avec pertes : pièces jointes non reprises) [Write_Small, Write_Lines] ; _extract_xml_from_pdf du module invoice_parse (pikepdf, un passage par nom de fichier connu) [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; invoice_validate local : XSLT FeRD Factur-X 1.09.2 (EN 16931) ou CEN + KoSIT XRechnung 3.0.2 via saxonche, recompilé à chaque appel ; avertissements non comptés [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large] ; utils.pdf.embed_xml_in_pdf (pikepdf, AFRelationship Alternative, XMP Factur-X) [Embed_Small, EmbedClassic_Small]
- **pycheval** (Python) : extract_facturx_from_pdf(BytesIO), renvoie une str [Extract_Small, Extract_Minimum] ; extract_facturx_from_pdf(BytesIO), renvoie une str ; Timeout: timeout (30 s) : appel interrompu (boucle infinie ?) [Extract_Large, Extract_Extended] ; fonction interne _embed (l'API publique embed_facturx_file_in_pdf lit le XML depuis un fichier), relation Data (défaut) [Embed_Small, EmbedClassic_Small] ; generate_xml depuis le modèle de parse_xml [Write_Small] ; generate_xml depuis le modèle de parse_xml ; InvalidXMLError: Element ./{urn:un:unece:uncefact:data:standard:ReusableAggregateBusinessInformationEntity:100}BuyerAssignedID has no text [Write_Lines] ; parse_xml(BytesIO) [Parse_Small] ; parse_xml(BytesIO) ; InvalidXMLError: Element ./{urn:un:unece:uncefact:data:standard:ReusableAggregateBusinessInformationEntity:100}BuyerAssignedID has no text [Parse_Lines] ; parse_xml(BytesIO) ; UnsupportedProfileError: Unsupported profile: urn:cen.eu:en16931:2017#compliant#urn:xeinkauf.de:kosit:xrechnung_3.0 [Parse_Large]
- **saxonche-schematron** (Python) : XSLT officiels CEN EN 16931 CII 1.3.17 (+ KoSIT XRechnung 3.0.2 schematron 2.6.0 si BT-24 XRechnung) via saxonche, compilés au premier appel ; erreurs fatal/error avec les niveaux KoSIT [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large]
- **zugferd-pdf** (Python) : generate_facturx_xml (totaux et TVA recalculés depuis les lignes) ; modèle construit à la main [Write_Small] ; package_zugferd_pdfa3u : ajoute OutputIntent sRGB et XMP PDF/A-3U, date de création = maintenant [Embed_Small, EmbedClassic_Small] ; validate_zugferd_pdf : contrôles maison du conteneur (XMP, OutputIntent, /AF) + BT-24 + 2 contrôles de décimales, pas de schematron [CheckPDF_Small, CheckPDF_Large]
- **core-invoice** (Rust) : cii::write_unchecked, modèle lu par read [Write_Small, Write_Lines] ; read → Invoice (modèle sémantique) [Parse_Small, Parse_Lines, Parse_Large] ; validate_xml(profil d'après BT-24) : règles EN 16931 réimplémentées en Rust (alignées sur CEN validation-1.3.16), listes de devises, syntaxe Peppol si profil Peppol ; BR-DE (feature `xrechnung`) sur les XRECHNUNG ; pas de XSD ni de règles CII-SR [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large]
- **einvoice** (Rust) : yaserde::de::from_str::<CrossIndustryInvoice> (modèle généré du XSD, celui de validate_invoice) [Parse_Small, Parse_Lines, Parse_Large] ; yaserde::ser::to_string, modèle lu par from_str [Write_Small, Write_Lines]
- **einvoice-validator** (Rust) : Validator::cen_cii() : XSD UN/CEFACT CII D16B + XSLT officiel CEN EN 16931 CII 1.3.17 exécuté en Rust ; seule règle hors ligne intégrée (XRechnung, Factur-X et Peppol demandent la feature `fetch` et un téléchargement) [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large] ; pdf::extract (einvoice-pdf, Limits par défaut) [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended]
- **en16931** (Rust) : read → Invoice (modèle sémantique EN 16931) [Parse_Small, Parse_Lines, Parse_Large] ; receive sur le PDF : couche conteneur (XMP, /AF, arbre EmbeddedFiles, accord avec BT-24) + validation complète du XML ; pas de contrôle PDF/A-3 au sens ISO 19005 [CheckPDF_Small, CheckPDF_Large] ; receive(Policy::default()) : couches syntaxe et schéma (règles CII-SR/CII-DT générées des artefacts CEN) + règles métier EN 16931 réimplémentées en Rust, profil déclaré par BT-24 (XRechnung 3.0 : BR-DE ; ZUGFeRD EXTENDED/BASIC : règles CEN cœur) ; CEN validation-1.3.17, KoSIT xrechnung-schematron v2.6.0 [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large] ; the PDF already carries invoice metadata in the namespace "urn:factur-x:pdfa:CrossIndustryDocument:invoice:1p0#"; embedding a second invoice is refused — adaptateur : read + Validated::<En16931>::new + zugferd::embed (la librairie réécrit le XML depuis son modèle ; exige un hôte PDF/A-3 sans métadonnées Factur-X, ne convertit pas en PDF/A) [Embed_Small, EmbedClassic_Small] ; write_unchecked(Syntax::Cii), modèle lu par read [Write_Small, Write_Lines] ; zugferd::extract (Limits par défaut) [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended]
- **erechnung-validator** (Rust) : Validator::beide(Szenarien::Automatisch) : KoSIT (XSD CII D16B, CEN EN 16931 1.3.16 en Schematron, XRechnung 3.0.2 / schematron 2.6.0 pour les XRECHNUNG, configuration KoSIT, CHANGELOG au 2026-01-31) + Mustang 2.26.0 (XSD et Schematron Factur-X du profil déclaré, XR 3.0 pour les XRECHNUNG) ; erreurs des deux outils additionnées [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large] ; Validator::beide_pdf : extraction du XML puis KoSIT + Mustang ; aucun contrôle du conteneur PDF/A-3 (partie PDF de Mustang/veraPDF non reprise) [CheckPDF_Small, CheckPDF_Large]
- **factur-x** (Rust) : adaptateur : parse + validate(&Invoice) ; 16 règles sur le modèle (BR-01 à BR-06, BR-09, BR-16, BR-CO-10, BR-CO-13 à BR-CO-18, BR-S-08) ; pas de XSD, pas de règles de profil ni nationales [ParseValidate_Large] ; amount has an invalid value "3.3333" — adaptateur : parse + validate(&Invoice) ; 16 règles sur le modèle (BR-01 à BR-06, BR-09, BR-16, BR-CO-10, BR-CO-13 à BR-CO-18, BR-S-08) ; pas de XSD, pas de règles de profil ni nationales [ParseValidate_Lines] ; amount has an invalid value "3.3333" — parse(Syntax::Cii) → Invoice [Parse_Lines] ; amount has an invalid value "9.9000" — adaptateur : parse + validate(&Invoice) ; 16 règles sur le modèle (BR-01 à BR-06, BR-09, BR-16, BR-CO-10, BR-CO-13 à BR-CO-18, BR-S-08) ; pas de XSD, pas de règles de profil ni nationales [ParseValidate_Small] ; amount has an invalid value "9.9000" — parse(Syntax::Cii) → Invoice [Parse_Small] ; parse(Syntax::Cii) → Invoice [Parse_Large] ; parse: amount has an invalid value "3.3333" [Write_Lines] ; parse: amount has an invalid value "9.9000" [Write_Small]
- **lopdf** (Rust) : plancher : extraction écrite à la main sur lopdf (load_mem complet, arbre /Names/EmbeddedFiles avec /Kids, inflate) ; pas une librairie Factur-X [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended]
- **pdf-invoice** (Rust) : ZugferdInvoice::from_xml (montants en f64) [Parse_Small, Parse_Lines, Parse_Large] ; adaptateur : Document::load_mem (pdfluent-lopdf) + extract_xml_attachment essayé pour 5 noms standard [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; adaptateur : Document::load_mem + embed_xml_attachment (Alternative) + add_zugferd_xmp + save_to ; réécrit tout le fichier, remplace le XMP sans conversion PDF/A [Embed_Small, EmbedClassic_Small] ; adaptateur : from_xml + validate (exigences du profil) + validate_en16931 (14 règles : BR-01 à BR-08, BR-CL-01, BR-CL-04, BR-CO-10, BR-CO-13, BR-CO-15, BR-S-08) ; pas de XSD [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large] ; profile validation: buyer_reference is required for EN16931/Extended; payment_means.type_code is required for EN16931/Extended — to_xml (valide d'abord le profil et les 14 règles EN 16931, refuse d'écrire sinon), modèle lu par from_xml [Write_Small, Write_Lines]
- **pdfa** (Rust) : embed (Level::En16931, horodatage fixe) : réécrit tout le fichier avec une table xref classique et un nouveau XMP PDF/A-3B ; refuse les hôtes à flux d'objets / flux xref [EmbedClassic_Small] ; embedded_file : premier flux /Type /EmbeddedFile trouvé par recherche d'octets [Extract_Small, Extract_Minimum] ; no embedded file read — embedded_file : premier flux /Type /EmbeddedFile trouvé par recherche d'octets [Extract_Large, Extract_Extended] ; the document has no the trailer dictionary — embed (Level::En16931, horodatage fixe) : réécrit tout le fichier avec une table xref classique et un nouveau XMP PDF/A-3B ; refuse les hôtes à flux d'objets / flux xref [Embed_Small]
- **zugferd** (Rust) : modèle construit à la main (contenu de EN16931_Einfach, niveau BASIC : le crate ne génère pas EN 16931) ; Invoice::to_xml_string (serde + quick-xml) [Write_Small]
- **zugferd-rs** (Rust) : embed_xml (profil EN 16931) : réécrit le PDF en PDF/A-3b (XMP Factur-X, OutputIntent sRGB) ; le XML est injecté tel quel [Embed_Small, EmbedClassic_Small] ; extract_xml (lopdf) [Extract_Small, Extract_Large, Extract_Minimum, Extract_Extended] ; invalid invoice type code: 877 — parse_cii_xml → Invoice [Parse_Large] ; parse_cii_xml → Invoice [Parse_Small, Parse_Lines] ; serialize_to_xml(ZugferdProfile::En16931), modèle lu par parse_cii_xml [Write_Small, Write_Lines] ; validate_xml (défaut) : XSD Factur-X 1.09.2 du profil déclaré + Schematron CEN EN 16931 évalué par libxml2 (XPath 1.0 seulement : les règles XPath 2.0 sont sautées, résultat déclaré incomplet) ; pas de BR-DE [ParseValidate_Small, ParseValidate_Lines, ParseValidate_Large]
