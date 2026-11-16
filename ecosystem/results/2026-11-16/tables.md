<!-- Généré par bench/ecosystem/report à partir de 2026-11-16 -->

#### Exactitude sur le corpus

| Librairie | XML valides rejetés | XML parsés | PDF extraits |
|---|---|---|---|
| **factur-x (ce dépôt)** | 6/26 | 26/26 | 104/104 |
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
| **factur-x (ce dépôt)** | **28,1 µs** | **33,8 µs** | **9,5 µs** | **32,6 µs** | 38 Ko | 212 |
| speedata | 74,5 µs (×2,6) | 241 µs (×7,1) | 26,5 µs (×2,8) | 84,6 µs (×2,6) | 401 Ko | 318 |
| litho | 107 µs (×3,8) | 152 µs (×4,5) | 52,8 µs (×5,6) | 116 µs (×3,5) | 585 Ko | 826 |
| papierleicht | 115 µs (×4,1) | 2,50 ms (×74,1) | 35,9 µs (×3,8) | 123 µs (×3,8) | 417 Ko | 490 |
| DCNWS | 216 µs (×7,7) | 2,44 ms (×72,3) | 111 µs (×11,7) | 208 µs (×6,4) | 184 Ko | 4 926 |
| pdfcpu | 332 µs (×11,8) | 2,47 ms (×73,1) | 228 µs (×24,0) | échec | 613 Ko | 3 628 |
| fX | 345 µs (×12,3) | 2,51 ms (×74,4) | 244 µs (×25,7) | échec | 628 Ko | 3 304 |
| kanjo | 361 µs (×12,9) | 2,79 ms (×82,7) | 259 µs (×27,3) | échec | 628 Ko | 3 307 |
| dotwave | 377 µs (×13,4) | 2,74 ms (×81,2) | 253 µs (×26,6) | échec | 628 Ko | 3 304 |
| angelo | 390 µs (×13,9) | 2,70 ms (×79,8) | 276 µs (×29,1) | échec | 628 Ko | 3 302 |
| gopdfattach | 392 µs (×13,9) | 557 µs (×16,5) | échec | 411 µs (×12,6) | 651 Ko | 3 135 |
| andee | 523 µs (×18,6) | 719 µs (×21,3) | 545 µs (×57,5) | échec | 660 Ko | 4 399 |
| taha | échec | échec | 2,85 ms (×300,5) | échec | — | — |

Mémoire et allocations : colonne « EN 16931, 149 Ko ». Tri sur cette même colonne.

#### Parsing

| Librairie | 13 Ko | 19 Ko | 1,3 Mo | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|
| **factur-x (ce dépôt)** | **14,7 µs** | **25,9 µs** | **1,12 ms** | 26 Ko | 80 |
| fX | 126 µs (×8,6) | 200 µs (×7,7) | 6,95 ms (×6,2) | 53 Ko | 1 165 |
| angelo | 128 µs (×8,7) | 204 µs (×7,9) | 7,22 ms (×6,4) | 56 Ko | 1 268 |
| taha | 147 µs (×10,0) | 239 µs (×9,2) | 9,30 ms (×8,3) | 209 Ko | 1 237 |
| dotwave | 151 µs (×10,3) | 224 µs (×8,6) | 7,31 ms (×6,5) | 64 Ko | 1 253 |
| DCNWS | 174 µs (×11,8) | 260 µs (×10,0) | 7,55 ms (×6,7) | 128 Ko | 2 079 |
| andee | 180 µs (×12,3) | 256 µs (×9,9) | 7,03 ms (×6,2) | 78 Ko | 1 326 |
| speedata | 225 µs (×15,3) | 365 µs (×14,1) | 7,56 ms (×6,7) | 331 Ko | 4 531 |
| kanjo | 257 µs (×17,5) | 392 µs (×15,1) | 13,80 ms (×12,3) | 99 Ko | 2 036 |
| gobl.de | 623 µs (×42,4) | 950 µs (×36,7) | 23,99 ms (×21,3) | 274 Ko | 6 251 |

Mémoire et allocations : colonne « 13 Ko ». Tri sur cette même colonne.

#### Validation

| Librairie | Valides rejetés /26 | 13 Ko | 19 Ko | 1,3 Mo | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|---|
| **factur-x (ce dépôt)** | 6/26 | **22,8 µs** | **40,4 µs** | **744 µs** | 23 Ko | 75 |
| fX | 26/26 | 126 µs (×5,5) | 202 µs (×5,0) | 6,90 ms (×9,3) | 58 Ko | 1 225 |
| pdf0 | 6/26 | 132 µs (×5,8) | 208 µs (×5,1) | 7,40 ms (×9,9) | 67 Ko | 1 462 |
| taha | 15/26 | 150 µs (×6,6) | 242 µs (×6,0) | 9,41 ms (×12,6) | 209 Ko | 1 246 |
| angelo | 4/26 | 159 µs (×7,0) | 236 µs (×5,9) | 7,01 ms (×9,4) | 91 Ko | 1 912 |
| andee | 22/26 | 211 µs (×9,3) | 310 µs (×7,7) | 7,25 ms (×9,7) | 100 Ko | 2 166 |
| kanjo | 9/26 | 297 µs (×13,0) | 433 µs (×10,7) | 13,74 ms (×18,5) | 120 Ko | 2 106 |
| DCNWS | 2/26 | 798 µs (×35,0) | 1,31 ms (×32,5) | 8,15 ms (×11,0) | 624 Ko | 15 279 |
| speedata | 0/26 | 1,34 ms (×58,8) | 2,30 ms (×57,0) | 9,21 ms (×12,4) | 1,9 Mo | 46 462 |
| gobl.de | 20/26 | 1,88 ms (×82,6) | 2,99 ms (×74,1) | 26,07 ms (×35,0) | 1018 Ko | 17 558 |
| papierleicht | 6/26 | 15,01 ms (×657,6) | 26,36 ms (×653,2) | 15,75 ms (×21,2) | 65,6 Mo | 254 579 |

Mémoire et allocations : colonne « 13 Ko ». Tri sur cette même colonne.

#### Écriture

| Librairie | 13 Ko | 19 Ko | Mémoire/op | Allocations/op |
|---|---|---|---|---|
| **factur-x (ce dépôt)** | **2,7 µs** | **4,5 µs** | 8 Ko | 1 |
| fX | 20,7 µs (×7,6) | 37,7 µs (×8,4) | 32 Ko | 175 |
| dotwave | 23,4 µs (×8,6) | 38,5 µs (×8,6) | 34 Ko | 361 |
| speedata | 32,4 µs (×11,8) | 56,5 µs (×12,6) | 44 Ko | 988 |
| DCNWS | 37,4 µs (×13,7) | 75,7 µs (×16,9) | 72 Ko | 1 130 |
| kanjo | 95,9 µs (×35,1) | 173 µs (×38,7) | 547 Ko | 826 |
| andee | 315 µs (×115,1) | 577 µs (×128,9) | 165 Ko | 2 980 |
| taha | 657 µs (×240,1) | 1,31 ms (×293,2) | 812 Ko | 341 |

Mémoire et allocations : colonne « 13 Ko ». Tri sur cette même colonne.

#### Injection

| Librairie | hôte à xref classique | hôte à flux xref | Mémoire/op | Allocations/op |
|---|---|---|---|---|
| **factur-x (ce dépôt)** | **77,1 µs** | **79,9 µs** | 229 Ko | 289 |
| kanjo | 355 µs (×4,6) | échec | 739 Ko | 3 211 |
| DCNWS | 425 µs (×5,5) | 1,90 ms (×23,8) | 2,5 Mo | 6 157 |
| pdf0 | 464 µs (×6,0) | 545 µs (×6,8) | 912 Ko | 9 390 |
| pdfcpu | 506 µs (×6,6) | 543 µs (×6,8) | 1,4 Mo | 3 844 |
| andee | 619 µs (×8,0) | échec | 1,8 Mo | 6 918 |
| gopdfattach | 1,06 ms (×13,7) | 1,13 ms (×14,1) | 4,2 Mo | 7 997 |
| dotwave | 1,30 ms (×16,8) | 1,36 ms (×17,0) | 3,9 Mo | 8 813 |
| litho | 35,68 ms (×462,8) | 35,96 ms (×449,8) | 35,3 Mo | 392 773 |

Mémoire et allocations : colonne « hôte à xref classique ». Tri sur cette même colonne.

#### Contrôle

| Librairie | 149 Ko | 1,1 Mo | Mémoire/op | Allocations/op |
|---|---|---|---|---|
| **factur-x (ce dépôt)** | **54,8 µs** | **70,5 µs** | 68 Ko | 301 |
| DCNWS | 1,18 ms (×21,5) | 3,34 ms (×47,3) | 896 Ko | 24 617 |
| pdf0 | 2,37 ms (×43,2) | 2,70 ms (×38,3) | 1,8 Mo | 7 536 |
| papierleicht | 15,29 ms (×278,9) | 18,19 ms (×257,9) | 66,0 Mo | 255 058 |
| litho | 20,61 ms (×375,9) | 21,73 ms (×308,0) | 16,7 Mo | 298 498 |

Mémoire et allocations : colonne « 149 Ko ». Tri sur cette même colonne.

#### Planchers

| Librairie | RawToken 13 Ko | RawToken 19 Ko | RawToken 1,3 Mo | memcpy 149 Ko | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|---|
| plancher | **94,2 µs** | **149 µs** | **6,47 ms** | **2,0 µs** | 31 Ko | 711 |

Mémoire et allocations : colonne « RawToken 13 Ko ». Tri sur cette même colonne.
