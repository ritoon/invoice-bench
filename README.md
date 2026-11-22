# invoice-bench

Bancs de mesure et comparatif de conformité des librairies de facture électronique Factur-X / ZUGFeRD / EN 16931, dans neuf langages. Ils mesurent [ZeRechnung](https://github.com/ritoon/zerechnung), un package Go, face aux autres librairies, sur les mêmes fichiers du corpus officiel ZUGFeRD et avec le même protocole.

| Dossier | Contenu |
|---|---|
| [`ecosystem/`](ecosystem) | ZeRechnung et 14 autres librairies Go : temps de réponse, mémoire et exactitude pour l'extraction, la lecture, la validation, l'écriture, l'injection dans un PDF et le contrôle d'un PDF hybride |
| [`polyglot/`](polyglot) | Les librairies TypeScript, Python, Rust, Java, C#, C, C++ et PHP, avec le même protocole ([`PROTOCOL.md`](polyglot/PROTOCOL.md)) : temps, temps CPU, mémoire, allocations, exactitude, et normes que chaque librairie applique réellement (fiches `COVERAGE.md` de chaque langage) |

## Utilisation

```bash
cd ecosystem && ./run.sh all        # banc Go : corpus, relevés, mesures, tableaux
cd polyglot && ./run.sh all         # les autres langages (prérequis dans polyglot/README.md)
```

Le banc Go demande Go 1.26 ou plus récent, git et un accès réseau ; une campagne complète prend une dizaine de minutes sur un Apple M4 Pro. Le banc multi-langages demande les chaînes d'outils de chaque langage et prend plusieurs heures. Le détail est dans le README de chaque dossier.

## Version de ZeRechnung mesurée

[`ecosystem/facturx/go.mod`](ecosystem/facturx/go.mod) fixe la version de ZeRechnung que mesurent les deux bancs. La variable `ZERECHNUNG_DIR` mesure un clone local à la place, par exemple une branche avant sa fusion :

```bash
ZERECHNUNG_DIR=~/src/zerechnung ./ecosystem/run.sh bench facturx
```

## Résultats

Chaque README donne les résultats de sa dernière campagne, et `results/<date>/` les sorties brutes et les tableaux complets. L'analyse est publiée sur le wiki de ZeRechnung, en français : [Benchmarks](https://github.com/ritoon/zerechnung/wiki/Benchmarks) pour les librairies Go, [Benchmarks autres langages](https://github.com/ritoon/zerechnung/wiki/Benchmarks-autres-langages) pour les autres.

## Historique

Ces bancs ont vécu dans le dépôt de ZeRechnung, sous `bench/`, jusqu'au 22 novembre 2026 (ritoon/zerechnung#114). Leur historique git a été gardé, et les numéros de pull request des anciens messages de commit renvoient à ritoon/zerechnung. Les sorties brutes produites avant le déplacement affichent les chemins de modules `github.com/ritoon/zerechnung/bench/…`, et leurs tableaux le libellé « ZeRechnung (ce dépôt) » ; les outils de rapport lisent encore ces sorties.

## Licence

Le code de ce dépôt est sous licence MIT ([LICENSE](LICENSE)). Les bancs téléchargent des librairies tierces sous leur propre licence, dont certaines sont restrictives (AGPL-3.0, GPL, BSL 1.1, licence commerciale d'évaluation) : la liste est dans la section Licences de [`ecosystem/README.md`](ecosystem/README.md#licences) et de [`polyglot/README.md`](polyglot/README.md#licences). Le corpus ZUGFeRD (Apache-2.0) et les artefacts officiels sont téléchargés et vérifiés, pas versionnés.
