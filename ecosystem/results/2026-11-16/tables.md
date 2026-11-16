<!-- Généré par bench/ecosystem/report à partir de 2026-11-16 -->

#### Exactitude sur le corpus

| Librairie | XML valides rejetés | XML parsés | PDF extraits |
|---|---|---|---|
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
| speedata | **74,5 µs** | 241 µs (×1,6) | **26,5 µs** | **84,6 µs** | 401 Ko | 318 |
| litho | 107 µs (×1,4) | **152 µs** | 52,8 µs (×2,0) | 116 µs (×1,4) | 585 Ko | 826 |
| papierleicht | 115 µs (×1,5) | 2,50 ms (×16,5) | 35,9 µs (×1,4) | 123 µs (×1,4) | 417 Ko | 490 |
| DCNWS | 216 µs (×2,9) | 2,44 ms (×16,1) | 111 µs (×4,2) | 208 µs (×2,5) | 184 Ko | 4 926 |
| pdfcpu | 332 µs (×4,5) | 2,47 ms (×16,3) | 228 µs (×8,6) | échec | 613 Ko | 3 628 |
| fX | 345 µs (×4,6) | 2,51 ms (×16,6) | 244 µs (×9,2) | échec | 628 Ko | 3 304 |
| kanjo | 361 µs (×4,9) | 2,79 ms (×18,4) | 259 µs (×9,8) | échec | 628 Ko | 3 307 |
| dotwave | 377 µs (×5,1) | 2,74 ms (×18,1) | 253 µs (×9,5) | échec | 628 Ko | 3 304 |
| angelo | 390 µs (×5,2) | 2,70 ms (×17,8) | 276 µs (×10,4) | échec | 628 Ko | 3 302 |
| gopdfattach | 392 µs (×5,3) | 557 µs (×3,7) | échec | 411 µs (×4,9) | 651 Ko | 3 135 |
| andee | 523 µs (×7,0) | 719 µs (×4,7) | 545 µs (×20,6) | échec | 660 Ko | 4 399 |
| taha | échec | échec | 2,85 ms (×107,6) | échec | — | — |

Mémoire et allocations : colonne « EN 16931, 149 Ko ». Tri sur cette même colonne.

#### Parsing

| Librairie | 13 Ko | 19 Ko | 1,3 Mo | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|
| fX | **126 µs** | **200 µs** | **6,95 ms** | 53 Ko | 1 165 |
| angelo | 128 µs (×1,0) | 204 µs (×1,0) | 7,22 ms (×1,0) | 56 Ko | 1 268 |
| taha | 147 µs (×1,2) | 239 µs (×1,2) | 9,30 ms (×1,3) | 209 Ko | 1 237 |
| dotwave | 151 µs (×1,2) | 224 µs (×1,1) | 7,31 ms (×1,1) | 64 Ko | 1 253 |
| DCNWS | 174 µs (×1,4) | 260 µs (×1,3) | 7,55 ms (×1,1) | 128 Ko | 2 079 |
| andee | 180 µs (×1,4) | 256 µs (×1,3) | 7,03 ms (×1,0) | 78 Ko | 1 326 |
| speedata | 225 µs (×1,8) | 365 µs (×1,8) | 7,56 ms (×1,1) | 331 Ko | 4 531 |
| kanjo | 257 µs (×2,0) | 392 µs (×2,0) | 13,80 ms (×2,0) | 99 Ko | 2 036 |
| gobl.de | 623 µs (×4,9) | 950 µs (×4,8) | 23,99 ms (×3,5) | 274 Ko | 6 251 |

Mémoire et allocations : colonne « 13 Ko ». Tri sur cette même colonne.

#### Validation

| Librairie | Valides rejetés /26 | 13 Ko | 19 Ko | 1,3 Mo | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|---|
| fX | 26/26 | **126 µs** | **202 µs** | **6,90 ms** | 58 Ko | 1 225 |
| pdf0 | 6/26 | 132 µs (×1,0) | 208 µs (×1,0) | 7,40 ms (×1,1) | 67 Ko | 1 462 |
| taha | 15/26 | 150 µs (×1,2) | 242 µs (×1,2) | 9,41 ms (×1,4) | 209 Ko | 1 246 |
| angelo | 4/26 | 159 µs (×1,3) | 236 µs (×1,2) | 7,01 ms (×1,0) | 91 Ko | 1 912 |
| andee | 22/26 | 211 µs (×1,7) | 310 µs (×1,5) | 7,25 ms (×1,1) | 100 Ko | 2 166 |
| kanjo | 9/26 | 297 µs (×2,3) | 433 µs (×2,1) | 13,74 ms (×2,0) | 120 Ko | 2 106 |
| DCNWS | 2/26 | 798 µs (×6,3) | 1,31 ms (×6,5) | 8,15 ms (×1,2) | 624 Ko | 15 279 |
| speedata | 0/26 | 1,34 ms (×10,6) | 2,30 ms (×11,4) | 9,21 ms (×1,3) | 1,9 Mo | 46 462 |
| gobl.de | 20/26 | 1,88 ms (×14,9) | 2,99 ms (×14,8) | 26,07 ms (×3,8) | 1018 Ko | 17 558 |
| papierleicht | 6/26 | 15,01 ms (×118,8) | 26,36 ms (×130,7) | 15,75 ms (×2,3) | 65,6 Mo | 254 579 |

Mémoire et allocations : colonne « 13 Ko ». Tri sur cette même colonne.

#### Écriture

| Librairie | 13 Ko | 19 Ko | Mémoire/op | Allocations/op |
|---|---|---|---|---|
| fX | **20,7 µs** | **37,7 µs** | 32 Ko | 175 |
| dotwave | 23,4 µs (×1,1) | 38,5 µs (×1,0) | 34 Ko | 361 |
| speedata | 32,4 µs (×1,6) | 56,5 µs (×1,5) | 44 Ko | 988 |
| DCNWS | 37,4 µs (×1,8) | 75,7 µs (×2,0) | 72 Ko | 1 130 |
| kanjo | 95,9 µs (×4,6) | 173 µs (×4,6) | 547 Ko | 826 |
| andee | 315 µs (×15,2) | 577 µs (×15,3) | 165 Ko | 2 980 |
| taha | 657 µs (×31,7) | 1,31 ms (×34,8) | 812 Ko | 341 |

Mémoire et allocations : colonne « 13 Ko ». Tri sur cette même colonne.

#### Injection

| Librairie | hôte à xref classique | hôte à flux xref | Mémoire/op | Allocations/op |
|---|---|---|---|---|
| kanjo | **355 µs** | échec | 739 Ko | 3 211 |
| DCNWS | 425 µs (×1,2) | 1,90 ms (×3,5) | 2,5 Mo | 6 157 |
| pdf0 | 464 µs (×1,3) | 545 µs (×1,0) | 912 Ko | 9 390 |
| pdfcpu | 506 µs (×1,4) | **543 µs** | 1,4 Mo | 3 844 |
| andee | 619 µs (×1,7) | échec | 1,8 Mo | 6 918 |
| gopdfattach | 1,06 ms (×3,0) | 1,13 ms (×2,1) | 4,2 Mo | 7 997 |
| dotwave | 1,30 ms (×3,7) | 1,36 ms (×2,5) | 3,9 Mo | 8 813 |
| litho | 35,68 ms (×100,5) | 35,96 ms (×66,2) | 35,3 Mo | 392 773 |

Mémoire et allocations : colonne « hôte à xref classique ». Tri sur cette même colonne.

#### Contrôle

| Librairie | 149 Ko | 1,1 Mo | Mémoire/op | Allocations/op |
|---|---|---|---|---|
| DCNWS | **1,18 ms** | 3,34 ms (×1,2) | 896 Ko | 24 617 |
| pdf0 | 2,37 ms (×2,0) | **2,70 ms** | 1,8 Mo | 7 536 |
| papierleicht | 15,29 ms (×13,0) | 18,19 ms (×6,7) | 66,0 Mo | 255 058 |
| litho | 20,61 ms (×17,5) | 21,73 ms (×8,1) | 16,7 Mo | 298 498 |

Mémoire et allocations : colonne « 149 Ko ». Tri sur cette même colonne.

#### Planchers

| Librairie | RawToken 13 Ko | RawToken 19 Ko | RawToken 1,3 Mo | memcpy 149 Ko | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|---|
| plancher | **94,2 µs** | **149 µs** | **6,47 ms** | **2,0 µs** | 31 Ko | 711 |

Mémoire et allocations : colonne « RawToken 13 Ko ». Tri sur cette même colonne.
