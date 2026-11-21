<!-- Généré par bench/polyglot/report à partir de 2026-11-17 -->

#### Exactitude sur le corpus

Référence : les schematrons officiels (CEN EN 16931 et KoSIT) signalent 6 des 26 XML « valides » du corpus (BR-CL-25 deux fois, CII-SR-465 quatre fois).

| Librairie | Langage | XML valides signalés | XML lus | PDF extraits |
|---|---|---|---|---|
| **ZeRechnung (ce dépôt)** | Go | 6/26 | 26/26 | 104/104 |
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
| **ZeRechnung (ce dépôt)** | Go | **29,3 µs** | **33,3 µs** | 10,7 µs | **34,7 µs** | — |
| einvoice-validator | Rust | 32,5 µs (×1,1) | 41,7 µs (×1,3) | 22,9 µs (×2,1) | 37,1 µs (×1,1) | 1,44 ms |
| lopdf | Rust | 72,6 µs (×2,5) | 95,7 µs (×2,9) | 43,2 µs (×4,0) | 66,0 µs (×1,9) | 391 µs |
| pdfa | Rust | 78,1 µs (×2,7) | échec | **496 ns** (×0,0) | échec | 140 µs |
| zugferd-pdf | Python | 78,5 µs (×2,7) | 116 µs (×3,5) | 52,2 µs (×4,9) | 93,5 µs (×2,7) | 1,72 ms |
| mcp-einvoicing-de | Python | 78,7 µs (×2,7) | 115 µs (×3,5) | 51,9 µs (×4,9) | 228 µs (×6,6) | 421 µs |
| pdf-invoice | Rust | 79,4 µs (×2,7) | 102 µs (×3,1) | 47,9 µs (×4,5) | 74,5 µs (×2,1) | 768 µs |
| zugferd-rs | Rust | 91,4 µs (×3,1) | 116 µs (×3,5) | 54,0 µs (×5,1) | 84,1 µs (×2,4) | 898 µs |
| @attestwire/en16931 | TypeScript | 184 µs (×6,3) | 199 µs (×6,0) | 106 µs (×9,9) | 206 µs (×5,9) | 4,62 ms |
| pycheval | Python | 199 µs (×6,8) | échec | 176 µs (×17) | échec | 929 µs |
| getafix | Python | 225 µs (×7,7) | 376 µs (×11) | 210 µs (×20) | 207 µs (×6,0) | 1,07 ms |
| factur-x | Python | 241 µs (×8,2) | 375 µs (×11) | 193 µs (×18) | 266 µs (×7,7) | 2,46 ms |
| @lucaapp/zugferd-ts | TypeScript | 253 µs (×8,6) | 293 µs (×8,8) | 122 µs (×11) | 227 µs (×6,5) | 4,70 ms |
| @ekwo-ai/factur-x | TypeScript | 269 µs (×9,2) | 315 µs (×9,5) | 134 µs (×13) | 242 µs (×7,0) | 4,52 ms |
| @conformo/toolkit | TypeScript | 281 µs (×9,6) | 10,13 ms (×305) | 134 µs (×13) | 244 µs (×7,0) | 4,59 ms |
| @stackforge-eu/factur-x | TypeScript | 285 µs (×9,7) | 331 µs (×10,0) | 143 µs (×13) | 248 µs (×7,2) | 5,02 ms |
| @stafyniaksacha/facturx | TypeScript | 373 µs (×13) | 3,15 ms (×95) | 169 µs (×16) | 436 µs (×13) | 7,39 ms |
| en16931 | Rust | 388 µs (×13) | 703 µs (×21) | 248 µs (×23) | 281 µs (×8,1) | 1,56 ms |
| @fin.cx/einvoice | TypeScript | 402 µs (×14) | 492 µs (×15) | 135 µs (×13) | 431 µs (×12) | 50,70 ms |
| @kontor-mcp/core | TypeScript | 683 µs (×23) | 1,46 ms (×44) | 659 µs (×62) | 737 µs (×21) | 8,61 ms |
| @sweberdev/summand | TypeScript | 892 µs (×30) | 20,19 ms (×607) | 437 µs (×41) | 1,04 ms (×30) | 3,36 ms |

En gras : le plus rapide de la colonne. (×n) : temps rapporté à celui du package Go de ce dépôt. Tri sur la colonne « EN 16931, 149 Ko » ; « Premier appel » : durée du premier appel dans le processus, même colonne (non mesurée pour Go).

#### Lecture XML → modèle

| Librairie | Langage | 13 Ko | 19 Ko | 1,3 Mo | Premier appel |
|---|---|---|---|---|---|
| **ZeRechnung (ce dépôt)** | Go | **16,3 µs** | **28,8 µs** | 1,11 ms | — |
| zugferd-rs | Rust | 24,0 µs (×1,5) | 36,9 µs (×1,3) | échec | 2,07 ms |
| pdf-invoice | Rust | 25,6 µs (×1,6) | 38,6 µs (×1,3) | 1,63 ms (×1,5) | 760 µs |
| core-invoice | Rust | 44,4 µs (×2,7) | 66,5 µs (×2,3) | 3,65 ms (×3,3) | 312 µs |
| @ekwo-ai/factur-x | TypeScript | 69,6 µs (×4,3) | 114 µs (×4,0) | 23,70 ms (×21) | 1,73 ms |
| en16931 | Rust | 73,5 µs (×4,5) | 126 µs (×4,4) | 6,75 ms (×6,1) | 1,49 ms |
| @conformo/toolkit | TypeScript | 92,1 µs (×5,7) | 157 µs (×5,4) | 1,74 ms (×1,6) | 1,56 ms |
| @attestwire/en16931 | TypeScript | 135 µs (×8,3) | 242 µs (×8,4) | 3,27 ms (×3,0) | 2,30 ms |
| @lucaapp/zugferd-ts | TypeScript | 162 µs (×10,0) | 287 µs (×10,0) | échec | 2,44 ms |
| pycheval | Python | 221 µs (×14) | échec | échec | 888 µs |
| factur-x-kit | TypeScript | 268 µs (×16) | échec | échec | 8,39 ms |
| einvoice | Rust | 323 µs (×20) | 544 µs (×19) | 11,96 ms (×11) | 1,56 ms |
| mcp-einvoicing-de | Python | 335 µs (×21) | 477 µs (×17) | **702 µs** (×0,6) | 1,06 ms |
| drafthorse | Python | 586 µs (×36) | 982 µs (×34) | 1,07 ms (×1,0) | 3,21 ms |
| @stafyniaksacha/facturx | TypeScript | 947 µs (×58) | 1,69 ms (×59) | 5,55 ms (×5,0) | 8,09 ms |
| @kontor-mcp/core | TypeScript | 1,08 ms (×66) | 1,67 ms (×58) | 4,35 ms (×3,9) | 12,34 ms |
| getafix | Python | 1,22 ms (×75) | 2,08 ms (×72) | échec | 2,72 ms |
| factur-x | Python | 1,40 ms (×86) | 1,98 ms (×69) | 11,79 ms (×11) | 4,99 ms |
| @fin.cx/einvoice | TypeScript | 2,09 ms (×128) | 3,72 ms (×129) | 7,11 ms (×6,4) | 13,51 ms |
| factur-x | Rust | échec | échec | 1,61 ms (×1,5) | — |

En gras : le plus rapide de la colonne. (×n) : temps rapporté à celui du package Go de ce dépôt. Tri sur la colonne « 13 Ko » ; « Premier appel » : durée du premier appel dans le processus, même colonne (non mesurée pour Go).

#### Lecture et validation

| Librairie | Langage | Valides signalés /26 | 13 Ko | 19 Ko | 1,3 Mo | Premier appel |
|---|---|---|---|---|---|---|
| **ZeRechnung (ce dépôt)** | Go | 6/26 | **23,1 µs** | 40,9 µs | **800 µs** | — |
| pdf-invoice | Rust | 20/26 | 27,1 µs (×1,2) | **40,2 µs** (×1,0) | 1,65 ms (×2,1) | 46,2 µs |
| @ekwo-ai/factur-x | TypeScript | 3/26 | 69,8 µs (×3,0) | 114 µs (×2,8) | 21,66 ms (×27) | 1,71 ms |
| core-invoice | Rust | 7/26 | 75,1 µs (×3,2) | 117 µs (×2,9) | 5,45 ms (×6,8) | 585 µs |
| en16931 | Rust | 2/26 | 75,4 µs (×3,3) | 128 µs (×3,1) | 6,70 ms (×8,4) | 1,66 ms |
| @attestwire/en16931 | TypeScript | 2/26 | 202 µs (×8,7) | 331 µs (×8,1) | 6,84 ms (×8,6) | 9,27 ms |
| factur-x-kit | TypeScript | 18/26 | 289 µs (×13) | échec | échec | 10,26 ms |
| einvoice-validator | Rust | 2/26 | 473 µs (×20) | 776 µs (×19) | 4,10 ms (×5,1) | 22,09 ms |
| factur-x | Python | 0/26 | 486 µs (×21) | 513 µs (×13) | 5,14 ms (×6,4) | 781 µs |
| drafthorse | Python | 0/26 | 492 µs (×21) | 529 µs (×13) | 7,98 ms (×10,0) | 1,82 ms |
| @sweberdev/summand | TypeScript | 6/26 | 802 µs (×35) | 1,36 ms (×33) | 7,49 ms (×9,4) | 16,68 ms |
| @stackforge-eu/factur-x | TypeScript | 0/26 | 828 µs (×36) | 878 µs (×21) | 10,06 ms (×13) | 32,98 ms |
| saxonche-schematron | Python | 6/26 | 903 µs (×39) | 1,46 ms (×36) | 12,61 ms (×16) | 133,18 ms |
| getafix | Python | 6/26 | 1,55 ms (×67) | 2,57 ms (×63) | échec | 3,01 ms |
| zugferd-rs | Rust | 7/26 | 1,96 ms (×85) | 2,73 ms (×67) | 5,15 ms (×6,4) | 5,79 ms |
| @fin.cx/einvoice | TypeScript | 0/26 | 3,09 ms (×134) | 5,28 ms (×129) | 10,57 ms (×13) | 17,20 ms |
| en16931 | TypeScript | 6/26 | 14,19 ms (×614) | 25,91 ms (×633) | 39,72 ms (×50) | 159,55 ms |
| erechnung-validator | Rust | 6/26 | 31,16 ms (×1 347) | 54,77 ms (×1 339) | 82,99 ms (×104) | 150,36 ms |
| @kontor-mcp/core | TypeScript | 2/26 | 44,31 ms (×1 916) | 51,53 ms (×1 260) | 71,43 ms (×89) | 142,98 ms |
| @stafyniaksacha/facturx | TypeScript | 12/26 | 60,76 ms (×2 627) | 76,86 ms (×1 879) | 78,84 ms (×99) | 131,89 ms |
| @conformo/toolkit | TypeScript | 6/26 | 71,49 ms (×3 091) | 78,08 ms (×1 909) | 121,80 ms (×152) | 117,26 ms |
| mcp-einvoicing-de | Python | 2/26 | 102,14 ms (×4 416) | 105,90 ms (×2 589) | 155,68 ms (×195) | 318,56 ms |
| @lucaapp/zugferd-ts | TypeScript | 8/26 | 2,67 s (×115 305) | 8,84 s (×216 176) | 4,15 s (×5 195) | 2,89 s |
| factur-x | Rust | 17/26 | échec | échec | 1,61 ms (×2,0) | — |

En gras : le plus rapide de la colonne. (×n) : temps rapporté à celui du package Go de ce dépôt. Tri sur la colonne « 13 Ko » ; « Premier appel » : durée du premier appel dans le processus, même colonne (non mesurée pour Go).

#### Écriture modèle → XML

| Librairie | Langage | 13 Ko | 19 Ko | Premier appel |
|---|---|---|---|---|
| **ZeRechnung (ce dépôt)** | Go | **3,1 µs** | **5,0 µs** | — |
| zugferd-rs | Rust | 6,7 µs (×2,2) | 11,6 µs (×2,3) | 59,3 µs |
| core-invoice | Rust | 10,6 µs (×3,4) | 16,9 µs (×3,4) | 52,3 µs |
| @ekwo-ai/factur-x | TypeScript | 10,9 µs (×3,5) | — | 431 µs |
| @stackforge-eu/factur-x | TypeScript | 14,3 µs (×4,6) | — | 488 µs |
| @conformo/toolkit | TypeScript | 21,5 µs (×6,9) | 37,3 µs (×7,5) | 1,65 ms |
| zugferd | Rust | 28,3 µs (×9,2) | — | 293 µs |
| einvoice | Rust | 36,2 µs (×12) | 70,7 µs (×14) | 1,03 ms |
| @attestwire/en16931 | TypeScript | 41,5 µs (×13) | 80,2 µs (×16) | 964 µs |
| en16931 | Rust | 45,8 µs (×15) | 79,3 µs (×16) | 1,55 ms |
| @stafyniaksacha/facturx | TypeScript | 72,5 µs (×23) | 120 µs (×24) | 1,18 ms |
| pycheval | Python | 76,6 µs (×25) | échec | 130 µs |
| @lucaapp/zugferd-ts | TypeScript | 106 µs (×34) | 190 µs (×38) | 5,03 ms |
| mcp-einvoicing-de | Python | 109 µs (×35) | 183 µs (×37) | 180 µs |
| zugferd-pdf | Python | 120 µs (×39) | — | 260 µs |
| erechnung | Python | 124 µs (×40) | — | 528 µs |
| factur-x-kit | TypeScript | 171 µs (×55) | échec | 4,97 ms |
| @fin.cx/einvoice | TypeScript | 336 µs (×109) | 454 µs (×91) | 3,47 ms |
| getafix | Python | 399 µs (×129) | 693 µs (×139) | 853 µs |
| factur-x | Python | 665 µs (×215) | 820 µs (×165) | 1,44 ms |
| node-zugferd | TypeScript | 928 µs (×301) | — | 4,84 ms |
| drafthorse | Python | 961 µs (×311) | 1,31 ms (×263) | 1,31 ms |
| @e-invoice-eu/core | TypeScript | 39,44 ms (×12 772) | — | 82,44 ms |
| factur-x | Rust | échec | échec | — |
| pdf-invoice | Rust | échec | échec | — |

En gras : le plus rapide de la colonne. (×n) : temps rapporté à celui du package Go de ce dépôt. Tri sur la colonne « 13 Ko » ; « Premier appel » : durée du premier appel dans le processus, même colonne (non mesurée pour Go).

#### Injection dans un PDF

| Librairie | Langage | hôte à xref classique | hôte à flux xref | Premier appel |
|---|---|---|---|---|
| **ZeRechnung (ce dépôt)** | Go | **84,6 µs** | **87,8 µs** | — |
| pdfa | Rust | 88,4 µs (×1,0) | échec | 104 µs |
| pdf-invoice | Rust | 156 µs (×1,8) | 178 µs (×2,0) | 172 µs |
| zugferd-rs | Rust | 200 µs (×2,4) | 214 µs (×2,4) | 220 µs |
| mcp-einvoicing-de | Python | 765 µs (×9,0) | 662 µs (×7,5) | 1,47 ms |
| zugferd-pdf | Python | 855 µs (×10) | 747 µs (×8,5) | 5,20 ms |
| @stackforge-eu/factur-x | TypeScript | 1,14 ms (×13) | 1,14 ms (×13) | 15,74 ms |
| @fin.cx/einvoice | TypeScript | 2,01 ms (×24) | 1,95 ms (×22) | 49,30 ms |
| drafthorse | Python | 2,07 ms (×24) | 2,09 ms (×24) | 3,20 ms |
| getafix | Python | 2,09 ms (×25) | 2,12 ms (×24) | 2,91 ms |
| pycheval | Python | 2,51 ms (×30) | 2,56 ms (×29) | 3,39 ms |
| @ekwo-ai/factur-x | TypeScript | 2,64 ms (×31) | 2,64 ms (×30) | 12,75 ms |
| @lucaapp/zugferd-ts | TypeScript | 2,72 ms (×32) | 2,64 ms (×30) | 12,94 ms |
| @conformo/toolkit | TypeScript | 2,87 ms (×34) | 2,98 ms (×34) | 16,23 ms |
| factur-x | Python | 4,05 ms (×48) | 3,85 ms (×44) | 6,72 ms |
| @stafyniaksacha/facturx | TypeScript | échec | échec | — |
| en16931 | Rust | échec | échec | — |

En gras : le plus rapide de la colonne. (×n) : temps rapporté à celui du package Go de ce dépôt. Tri sur la colonne « hôte à xref classique » ; « Premier appel » : durée du premier appel dans le processus, même colonne (non mesurée pour Go).

#### Contrôle d'un PDF hybride

| Librairie | Langage | 149 Ko | 1,1 Mo | Premier appel |
|---|---|---|---|---|
| **ZeRechnung (ce dépôt)** | Go | **75,0 µs** | **85,4 µs** | — |
| zugferd-pdf | Python | 212 µs (×2,8) | 269 µs (×3,2) | 464 µs |
| @attestwire/en16931 | TypeScript | 359 µs (×4,8) | 408 µs (×4,8) | 11,45 ms |
| en16931 | Rust | 479 µs (×6,4) | 812 µs (×9,5) | 509 µs |
| @sweberdev/summand | TypeScript | 1,84 ms (×25) | 21,50 ms (×252) | 17,88 ms |
| erechnung-validator | Rust | 31,74 ms (×423) | 35,68 ms (×418) | 32,48 ms |
| @conformo/toolkit | TypeScript | 73,69 ms (×983) | 96,17 ms (×1 126) | 118,80 ms |

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
