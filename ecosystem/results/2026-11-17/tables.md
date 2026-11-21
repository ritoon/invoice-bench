<!-- Généré par bench/ecosystem/report à partir de 2026-11-17 -->

#### Exactitude sur le corpus

| Librairie | XML valides rejetés | XML parsés | PDF extraits |
|---|---|---|---|
| **ZeRechnung (ce dépôt)** | 6/26 | 26/26 | 104/104 |
| DCNWS | 2/26 | 26/26 | 104/104 |
| andee | 22/26 | 26/26 | 75/104 |
| angelo | 4/26 | 26/26 | 75/104 |
| dotwave | — | 26/26 | 75/104 |
| fX | 26/26 | 26/26 | 75/104 |
| gobl.de | 20/26 | 26/26 | — |
| gopdfattach | — | — | 97/104 |
| kanjo | 9/26 | 26/26 | 75/104 |
| litho | — | — | 104/104 |
| papierleicht | 6/26 | — | 104/104 |
| pdf0 | 6/26 | — | 104/104 |
| pdfcpu | — | — | 75/104 |
| speedata | 0/26 | — | 104/104 |
| taha | 15/26 | 26/26 | 48/104 |

#### Extraction

| Librairie | EN 16931, 149 Ko | EN 16931, 1,1 Mo | MINIMUM FNFE, 89 Ko | EXTENDED ZF 2.0, 157 Ko | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|---|
| **ZeRechnung (ce dépôt)** | **30,2 µs** | **33,6 µs** | **11,0 µs** | **35,6 µs** | 39 Ko | 227 |
| speedata | 75,2 µs (×2,5) | 239 µs (×7,1) | 26,9 µs (×2,5) | 84,9 µs (×2,4) | 401 Ko | 318 |
| litho | 106 µs (×3,5) | 149 µs (×4,4) | 50,8 µs (×4,6) | 113 µs (×3,2) | 585 Ko | 826 |
| papierleicht | 113 µs (×3,7) | 2,52 ms (×75,1) | 39,0 µs (×3,6) | 137 µs (×3,8) | 417 Ko | 490 |
| DCNWS | 209 µs (×6,9) | 2,46 ms (×73,1) | 109 µs (×9,9) | 185 µs (×5,2) | 184 Ko | 4 926 |
| pdfcpu | 317 µs (×10,5) | 2,60 ms (×77,4) | 240 µs (×21,9) | échec | 613 Ko | 3 629 |
| fX | 332 µs (×11,0) | 2,56 ms (×76,1) | 239 µs (×21,8) | échec | 628 Ko | 3 305 |
| dotwave | 343 µs (×11,3) | 2,71 ms (×80,6) | 237 µs (×21,6) | échec | 628 Ko | 3 304 |
| angelo | 356 µs (×11,8) | 2,69 ms (×80,0) | 296 µs (×27,0) | échec | 628 Ko | 3 303 |
| kanjo | 386 µs (×12,8) | 2,87 ms (×85,4) | 256 µs (×23,4) | échec | 628 Ko | 3 308 |
| gopdfattach | 396 µs (×13,1) | 543 µs (×16,2) | échec | 424 µs (×11,9) | 651 Ko | 3 135 |
| andee | 516 µs (×17,1) | 704 µs (×21,0) | 545 µs (×49,8) | échec | 660 Ko | 4 399 |
| taha | échec | échec | 2,84 ms (×259,3) | échec | — | — |

Mémoire et allocations : colonne « EN 16931, 149 Ko ». Tri sur cette même colonne.

#### Parsing

| Librairie | 13 Ko | 19 Ko | 1,3 Mo | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|
| **ZeRechnung (ce dépôt)** | **16,9 µs** | **29,8 µs** | **1,07 ms** | 28 Ko | 108 |
| fX | 124 µs (×7,3) | 194 µs (×6,5) | 7,00 ms (×6,5) | 53 Ko | 1 165 |
| angelo | 131 µs (×7,7) | 210 µs (×7,0) | 7,13 ms (×6,7) | 56 Ko | 1 268 |
| dotwave | 144 µs (×8,5) | 214 µs (×7,2) | 6,76 ms (×6,3) | 64 Ko | 1 253 |
| taha | 147 µs (×8,7) | 235 µs (×7,9) | 9,33 ms (×8,7) | 209 Ko | 1 237 |
| DCNWS | 165 µs (×9,8) | 252 µs (×8,4) | 6,91 ms (×6,5) | 128 Ko | 2 079 |
| andee | 182 µs (×10,7) | 254 µs (×8,5) | 7,39 ms (×6,9) | 78 Ko | 1 326 |
| speedata | 225 µs (×13,3) | 368 µs (×12,4) | 7,80 ms (×7,3) | 331 Ko | 4 531 |
| kanjo | 253 µs (×14,9) | 389 µs (×13,1) | 14,08 ms (×13,1) | 99 Ko | 2 036 |
| gobl.de | 564 µs (×33,3) | 897 µs (×30,1) | 21,97 ms (×20,5) | 275 Ko | 6 251 |

Mémoire et allocations : colonne « 13 Ko ». Tri sur cette même colonne.

#### Validation

| Librairie | Valides rejetés /26 | 13 Ko | 19 Ko | 1,3 Mo | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|---|
| **ZeRechnung (ce dépôt)** | 6/26 | **24,3 µs** | **41,1 µs** | **742 µs** | 23 Ko | 79 |
| fX | 26/26 | 123 µs (×5,1) | 202 µs (×4,9) | 7,09 ms (×9,6) | 58 Ko | 1 225 |
| pdf0 | 6/26 | 126 µs (×5,2) | 199 µs (×4,8) | 6,81 ms (×9,2) | 67 Ko | 1 462 |
| taha | 15/26 | 147 µs (×6,1) | 239 µs (×5,8) | 9,21 ms (×12,4) | 209 Ko | 1 246 |
| angelo | 4/26 | 172 µs (×7,1) | 251 µs (×6,1) | 7,33 ms (×9,9) | 91 Ko | 1 912 |
| andee | 22/26 | 208 µs (×8,6) | 311 µs (×7,6) | 6,95 ms (×9,4) | 100 Ko | 2 166 |
| kanjo | 9/26 | 293 µs (×12,1) | 423 µs (×10,3) | 13,59 ms (×18,3) | 120 Ko | 2 106 |
| DCNWS | 2/26 | 754 µs (×31,0) | 1,24 ms (×30,2) | 7,45 ms (×10,0) | 624 Ko | 15 279 |
| speedata | 0/26 | 1,28 ms (×52,7) | 2,37 ms (×57,7) | 9,27 ms (×12,5) | 1,9 Mo | 46 462 |
| gobl.de | 20/26 | 1,77 ms (×73,0) | 2,88 ms (×70,2) | 24,48 ms (×33,0) | 1017 Ko | 17 557 |
| papierleicht | 6/26 | 14,65 ms (×603,2) | 25,74 ms (×626,9) | 15,43 ms (×20,8) | 65,6 Mo | 254 579 |

Mémoire et allocations : colonne « 13 Ko ». Tri sur cette même colonne.

#### Écriture

| Librairie | 13 Ko | 19 Ko | Mémoire/op | Allocations/op |
|---|---|---|---|---|
| **ZeRechnung (ce dépôt)** | **3,0 µs** | **5,1 µs** | 8 Ko | 1 |
| fX | 20,2 µs (×6,6) | 38,6 µs (×7,5) | 32 Ko | 175 |
| dotwave | 22,1 µs (×7,2) | 35,0 µs (×6,8) | 34 Ko | 361 |
| speedata | 30,5 µs (×10,0) | 54,9 µs (×10,7) | 44 Ko | 988 |
| DCNWS | 36,2 µs (×11,9) | 72,0 µs (×14,0) | 72 Ko | 1 130 |
| kanjo | 98,3 µs (×32,2) | 179 µs (×34,9) | 547 Ko | 826 |
| andee | 311 µs (×102,1) | 562 µs (×109,4) | 165 Ko | 2 980 |
| taha | 652 µs (×213,7) | 1,26 ms (×245,7) | 812 Ko | 341 |

Mémoire et allocations : colonne « 13 Ko ». Tri sur cette même colonne.

#### Injection

| Librairie | hôte à xref classique | hôte à flux xref | Mémoire/op | Allocations/op |
|---|---|---|---|---|
| **ZeRechnung (ce dépôt)** | **86,7 µs** | **89,2 µs** | 231 Ko | 301 |
| kanjo | 361 µs (×4,2) | échec | 739 Ko | 3 213 |
| pdf0 | 452 µs (×5,2) | 522 µs (×5,8) | 912 Ko | 9 390 |
| DCNWS | 486 µs (×5,6) | 1,66 ms (×18,6) | 2,5 Mo | 6 157 |
| pdfcpu | 514 µs (×5,9) | 546 µs (×6,1) | 1,4 Mo | 3 844 |
| andee | 623 µs (×7,2) | échec | 1,8 Mo | 6 918 |
| gopdfattach | 1,02 ms (×11,8) | 1,12 ms (×12,5) | 4,2 Mo | 7 998 |
| dotwave | 1,21 ms (×13,9) | 1,24 ms (×13,9) | 3,9 Mo | 8 813 |
| litho | 35,06 ms (×404,1) | 36,16 ms (×405,4) | 35,3 Mo | 392 772 |

Mémoire et allocations : colonne « hôte à xref classique ». Tri sur cette même colonne.

#### Contrôle

| Librairie | 149 Ko | 1,1 Mo | Mémoire/op | Allocations/op |
|---|---|---|---|---|
| **ZeRechnung (ce dépôt)** | **74,2 µs** | **86,0 µs** | 69 Ko | 508 |
| DCNWS | 1,14 ms (×15,3) | 3,46 ms (×40,2) | 896 Ko | 24 617 |
| pdf0 | 2,20 ms (×29,6) | 2,54 ms (×29,6) | 1,8 Mo | 7 536 |
| papierleicht | 15,78 ms (×212,6) | 18,34 ms (×213,2) | 66,0 Mo | 255 059 |
| litho | 19,87 ms (×267,7) | 21,43 ms (×249,2) | 16,7 Mo | 298 497 |

Mémoire et allocations : colonne « 149 Ko ». Tri sur cette même colonne.

#### Planchers

| Librairie | RawToken 13 Ko | RawToken 19 Ko | RawToken 1,3 Mo | memcpy 149 Ko | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|---|
| plancher | **100 µs** | **158 µs** | **7,14 ms** | **1,6 µs** | 31 Ko | 711 |

Mémoire et allocations : colonne « RawToken 13 Ko ». Tri sur cette même colonne.
