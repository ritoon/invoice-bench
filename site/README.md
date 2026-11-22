# Results site

A two-page micro-site, published on GitHub Pages at [ritoon.github.io/invoice-bench](https://ritoon.github.io/invoice-bench/), that shows the results of the latest run of both campaigns (the cross-language one in `results/`, the Go one in `go/results/`):

| Page | Contents | Source records |
|---|---|---|
| `index.html`, **CPU & memory** | Peak memory and CPU / wall time per language; for each benchmark, CPU time per call against peak memory (scatter), and every library in a sortable table with wall time and allocations | `results/<date>/resources_*.jsonl`, `alloc_*.jsonl`, `env.txt` |
| `compliance.html`, **Compliance** | The confidence score per country (Germany for now) with its factors and parts; the 26 valid CII invoices of the corpus, file by file, as each validator judges them against the official CEN and KoSIT schematrons; PDFs extracted and XML parsed; the standards each library enforces | `results/<date>/cases_*.jsonl` and `survey_*`, `go/results/<date>/cases_*.jsonl` and `survey_*.txt`, `cases/*.tsv`, `standards.json` |

## Layout

| Path | Role |
|---|---|
| `main.go` | Generator: reads the latest `results/<date>` and `go/results/<date>` folders, writes `data/perf.json` and `data/compliance.json`, copies `web/`. Standard library only, plus the [`cases`](../cases) package of this repository, which scores the conformance cases |
| `standards.json` | The "Enforced standards" table of [`../README.md`](../README.md#enforced-standards), with a level per cell (enforced, partial, external, no). Kept by hand: update it with the README |
| `web/` | The pages: HTML5, Tailwind CSS through its CDN, vanilla JavaScript. `css/main.css` holds the colour tokens (light and dark), `js/common.js` the shared helpers, `components/` one object per component (`create()` / `render()`), `js/perf.js` and `js/compliance.js` the two pages |

The generator mirrors the rules of [`report/`](../report): medians of the series, the Go round at `GOMAXPROCS=1` and the libraries that spawn child processes left out of the per-language summary. The confidence scores come from the `cases` package, which both report tools use: ZeRechnung and the other languages from `results/`, the other Go libraries from `go/results/`. The figures therefore match the `tables.md` files.

## Preview

From this folder:

```bash
go run . -serve :8000
```

then open http://localhost:8000. The data is generated at start-up; `web/` is served from disk, so edits to the pages show on reload. To build the static site as it is published:

```bash
go run . -out ../_site
```

`-results` and `-go` pick another results folder than the latest one. `-out` refuses a non-empty folder that is not a previous build.

## Publishing

`.github/workflows/pages.yml` builds the site and deploys it on every push to `main` that touches `site/`, the results of either campaign or the conformance cases, and on demand (Actions tab, "Run workflow"). It needs GitHub Pages to be enabled once, with **Settings → Pages → Build and deployment → Source: GitHub Actions**. The CI workflow also checks that the site builds from the committed records.

A new run committed under `results/<date>/` is picked up without any change here, as long as the record formats stay the same. Only `standards.json` must be updated by hand when the README table changes.
