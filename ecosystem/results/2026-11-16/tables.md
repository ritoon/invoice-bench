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
| **factur-x (ce dépôt)** | **29,0 µs** | **33,9 µs** | **9,5 µs** | **32,8 µs** | 38 Ko | 212 |
| speedata | 74,5 µs (×2,6) | 241 µs (×7,1) | 26,5 µs (×2,8) | 84,6 µs (×2,6) | 401 Ko | 318 |
| litho | 107 µs (×3,7) | 152 µs (×4,5) | 52,8 µs (×5,5) | 116 µs (×3,5) | 585 Ko | 826 |
| papierleicht | 115 µs (×4,0) | 2,50 ms (×73,7) | 35,9 µs (×3,8) | 123 µs (×3,7) | 417 Ko | 490 |
| DCNWS | 216 µs (×7,4) | 2,44 ms (×72,0) | 111 µs (×11,6) | 208 µs (×6,3) | 184 Ko | 4 926 |
| pdfcpu | 332 µs (×11,4) | 2,47 ms (×72,8) | 228 µs (×23,9) | échec | 613 Ko | 3 628 |
| fX | 345 µs (×11,9) | 2,51 ms (×74,1) | 244 µs (×25,6) | échec | 628 Ko | 3 304 |
| kanjo | 361 µs (×12,5) | 2,79 ms (×82,3) | 259 µs (×27,1) | échec | 628 Ko | 3 307 |
| dotwave | 377 µs (×13,0) | 2,74 ms (×80,8) | 253 µs (×26,5) | échec | 628 Ko | 3 304 |
| angelo | 390 µs (×13,4) | 2,70 ms (×79,5) | 276 µs (×28,9) | échec | 628 Ko | 3 302 |
| gopdfattach | 392 µs (×13,5) | 557 µs (×16,4) | échec | 411 µs (×12,5) | 651 Ko | 3 135 |
| andee | 523 µs (×18,0) | 719 µs (×21,2) | 545 µs (×57,2) | échec | 660 Ko | 4 399 |
| taha | échec | échec | 2,85 ms (×298,6) | échec | — | — |

Mémoire et allocations : colonne « EN 16931, 149 Ko ». Tri sur cette même colonne.

#### Parsing

| Librairie | 13 Ko | 19 Ko | 1,3 Mo | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|
| **factur-x (ce dépôt)** | **120 µs** | **186 µs** | 8,90 ms (×1,3) | 56 Ko | 1 250 |
| fX | 126 µs (×1,1) | 200 µs (×1,1) | **6,95 ms** | 53 Ko | 1 165 |
| angelo | 128 µs (×1,1) | 204 µs (×1,1) | 7,22 ms (×1,0) | 56 Ko | 1 268 |
| taha | 147 µs (×1,2) | 239 µs (×1,3) | 9,30 ms (×1,3) | 209 Ko | 1 237 |
| dotwave | 151 µs (×1,3) | 224 µs (×1,2) | 7,31 ms (×1,1) | 64 Ko | 1 253 |
| DCNWS | 174 µs (×1,5) | 260 µs (×1,4) | 7,55 ms (×1,1) | 128 Ko | 2 079 |
| andee | 180 µs (×1,5) | 256 µs (×1,4) | 7,03 ms (×1,0) | 78 Ko | 1 326 |
| speedata | 225 µs (×1,9) | 365 µs (×2,0) | 7,56 ms (×1,1) | 331 Ko | 4 531 |
| kanjo | 257 µs (×2,2) | 392 µs (×2,1) | 13,80 ms (×2,0) | 99 Ko | 2 036 |
| gobl.de | 623 µs (×5,2) | 950 µs (×5,1) | 23,99 ms (×3,5) | 274 Ko | 6 251 |

Mémoire et allocations : colonne « 13 Ko ». Tri sur cette même colonne.

#### Validation

| Librairie | Valides rejetés /26 | 13 Ko | 19 Ko | 1,3 Mo | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|---|
| fX | 26/26 | **126 µs** | **202 µs** | **6,90 ms** | 58 Ko | 1 225 |
| **factur-x (ce dépôt)** | 6/26 | 127 µs (×1,0) | 204 µs (×1,0) | 12,27 ms (×1,8) | 53 Ko | 1 245 |
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
| **factur-x (ce dépôt)** | **2,8 µs** | **4,5 µs** | 8 Ko | 1 |
| fX | 20,7 µs (×7,3) | 37,7 µs (×8,4) | 32 Ko | 175 |
| dotwave | 23,4 µs (×8,3) | 38,5 µs (×8,6) | 34 Ko | 361 |
| speedata | 32,4 µs (×11,4) | 56,5 µs (×12,6) | 44 Ko | 988 |
| DCNWS | 37,4 µs (×13,2) | 75,7 µs (×16,9) | 72 Ko | 1 130 |
| kanjo | 95,9 µs (×33,8) | 173 µs (×38,6) | 547 Ko | 826 |
| andee | 315 µs (×111,2) | 577 µs (×128,6) | 165 Ko | 2 980 |
| taha | 657 µs (×231,7) | 1,31 ms (×292,5) | 812 Ko | 341 |

Mémoire et allocations : colonne « 13 Ko ». Tri sur cette même colonne.

#### Injection

| Librairie | hôte à xref classique | hôte à flux xref | Mémoire/op | Allocations/op |
|---|---|---|---|---|
| **factur-x (ce dépôt)** | **79,8 µs** | **82,7 µs** | 229 Ko | 289 |
| kanjo | 355 µs (×4,4) | échec | 739 Ko | 3 211 |
| DCNWS | 425 µs (×5,3) | 1,90 ms (×23,0) | 2,5 Mo | 6 157 |
| pdf0 | 464 µs (×5,8) | 545 µs (×6,6) | 912 Ko | 9 390 |
| pdfcpu | 506 µs (×6,3) | 543 µs (×6,6) | 1,4 Mo | 3 844 |
| andee | 619 µs (×7,8) | échec | 1,8 Mo | 6 918 |
| gopdfattach | 1,06 ms (×13,2) | 1,13 ms (×13,7) | 4,2 Mo | 7 997 |
| dotwave | 1,30 ms (×16,3) | 1,36 ms (×16,4) | 3,9 Mo | 8 813 |
| litho | 35,68 ms (×447,0) | 35,96 ms (×435,1) | 35,3 Mo | 392 773 |

Mémoire et allocations : colonne « hôte à xref classique ». Tri sur cette même colonne.

#### Contrôle

| Librairie | 149 Ko | 1,1 Mo | Mémoire/op | Allocations/op |
|---|---|---|---|---|
| **factur-x (ce dépôt)** | **165 µs** | **204 µs** | 98 Ko | 1 470 |
| DCNWS | 1,18 ms (×7,2) | 3,34 ms (×16,4) | 896 Ko | 24 617 |
| pdf0 | 2,37 ms (×14,4) | 2,70 ms (×13,2) | 1,8 Mo | 7 536 |
| papierleicht | 15,29 ms (×92,9) | 18,19 ms (×89,1) | 66,0 Mo | 255 058 |
| litho | 20,61 ms (×125,2) | 21,73 ms (×106,5) | 16,7 Mo | 298 498 |

Mémoire et allocations : colonne « 149 Ko ». Tri sur cette même colonne.

#### Planchers

| Librairie | RawToken 13 Ko | RawToken 19 Ko | RawToken 1,3 Mo | memcpy 149 Ko | Mémoire/op | Allocations/op |
|---|---|---|---|---|---|---|
| plancher | **94,2 µs** | **149 µs** | **6,47 ms** | **2,0 µs** | 31 Ko | 711 |

Mémoire et allocations : colonne « RawToken 13 Ko ». Tri sur cette même colonne.
