# Go Factur-X ecosystem benchmarks

This folder measures the existing Go libraries on the official ZUGFeRD corpus: response time, memory, and accuracy (extraction, parsing, validation, embedding). The results and their analysis are published on the ZeRechnung wiki, in French:

- [Benchmarks](https://github.com/ritoon/zerechnung/wiki/Benchmarks)
- [Grille de la norme](https://github.com/ritoon/zerechnung/wiki/Grille-de-la-norme) (coverage of the standard)

Each subfolder is an independent Go module. The `facturx` module measures [ZeRechnung](https://github.com/ritoon/zerechnung) at the version its `go.mod` requires, or a local clone (see `ZERECHNUNG_DIR` below).

## Requirements

- Go 1.26 or newer. With `GOTOOLCHAIN=auto` (set by the script), Go downloads 1.27.0 for DCNWS and pdf0, and 1.27.3 for litho.
- git and network access, to fetch the corpus, the dependencies and the clone of solarpush/fX.

## Usage

```bash
./run.sh setup     # ZUGFeRD corpus (sparse) and host PDFs in ../testdata, clone of solarpush/fX
./run.sh survey    # accuracy on the whole corpus
./run.sh cases     # conformance cases of ../cases: confidence score per country
./run.sh bench     # benchmarks (COUNT=3, BENCHTIME=1s by default)
./run.sh report    # Markdown tables in results/<date>/tables.md
./run.sh all       # all five in a row
```

Each command accepts a list of modules (`./run.sh bench pdf0 litho`). The `OUT`, `COUNT` and `BENCHTIME` variables change the output folder, the number of runs and the duration of each benchmark.

`ZERECHNUNG_DIR` measures a local clone of ZeRechnung instead of the version `facturx/go.mod` requires, for example a branch before it is merged: `ZERECHNUNG_DIR=~/src/zerechnung ./run.sh bench facturx`. [`zerechnung.sh`](zerechnung.sh) then writes `facturx/local.mod`, a copy of `facturx/go.mod` with a `replace` directive pointing to that clone, and passes it to the go command with `-modfile`. `bench_facturx.txt` starts with the measured version, as `go list -m` prints it. To measure a newly released version: `cd facturx && go get github.com/ritoon/zerechnung@<version>`. A full run takes about ten minutes on an Apple M4 Pro.

## Layout

| Folder | Contents |
|---|---|
| `common/` | Shared input files and measurement functions (`Extract`, `Bytes`, `Write`, `Embed`, `Survey*`, `Cases`) |
| `floor/` | Floors: tokenization with `encoding/xml`, memory copy |
| `<library>/` | One module per library: benchmarks, `TestSurvey`, `TestCases` for the libraries that validate, sometimes `TestDetail` |
| `report/` | Tool that builds the wiki tables from the raw outputs |
| `zerechnung.sh` | Choice of the measured ZeRechnung version (`ZERECHNUNG_DIR`), shared with the cross-language campaign (`../run.sh`) |
| `results/<date>/` | Raw outputs (`bench_*.txt`, `survey_*.txt`, `cases_*.jsonl`) and tables (`tables.md`) |
| `.deps/` | Clone of solarpush/fX, fetched by `./run.sh setup`, not versioned |

The corpus and the host PDFs are shared with the other languages: `./run.sh setup` writes them to [`../testdata`](../testdata), and `FX_TESTDATA` (or `FX_CORPUS` for the corpus alone) points the harness elsewhere.

A single module for all the libraries is impossible: pdfcpu broke its API between versions 0.11 and 0.16, and these libraries use different versions of it.

## Measured operations

| Benchmark | Measures |
|---|---|
| `Extract_*` | PDF → invoice XML |
| `Parse_*` | XML → Go model |
| `ParseValidate_*` | XML → parsing and full validation |
| `Write_*` | Go model → CII XML (parsing done once, outside the measurement) |
| `Embed_Small` | host PDF with xref stream + XML → Factur-X PDF, then the result is read back |
| `EmbedClassic_Small` | same with the host with a classic xref table |
| `CheckPDF_*` | full check of a Factur-X PDF |

Each measurement first checks the result (the extracted XML is indeed a CII, the produced PDF can be read back). If the library cannot do the operation, the benchmark is skipped; if it returns a wrong result, the benchmark fails.

The accuracy surveys (`TestSurvey`) cover 26 valid CII XML files (`XML-Rechnung/CII`) and 104 valid hybrid PDFs (`ZUGFeRDv2/correct` and `XML-Rechnung/FX`).

The conformance cases (`TestCases`, `./run.sh cases`) run each validating library on the 561 documents of [`../cases`](../cases/README.md), valid and invalid, CII and UBL, whose verdict and expected rules come from the KoSIT validator, and give it a confidence score per country. The validation call is the one of `ParseValidate_*`, rule identifiers read as in `TestDetail`; pdf0 is called through `formalis.Detect`, which picks the validator of the document (EN 16931 or XRechnung, CII or UBL), where the bench uses the Factur-X EN16931 profile. Without `CASES_OUT`, which `run.sh` sets, `TestCases` is skipped.

## How each library is called

| Module | Library | Calls |
|---|---|---|
| `speedata` | speedata/einvoice v0.0.23 + pdfdisassembler v0.0.9 | `einvoice.ParseReader`, `inv.Validate`, `inv.Write`; extraction through `pdf.Open` → `EmbeddedFiles` |
| `andee` | andeedotnet/go-xinvoice + go-xinvoice-pdf v0.1.4 | `xinvoicepdf.ExtractXML`, `xinvoice.ParseXML`, `xinvoice.Validate`, `inv.ToXML`, `xinvoicepdf.Embed` |
| `dotwave` | dotwavehq/go-einvoice v0.7.1 | `zugferd.ExtractXML`, `einvoice.Parse`, `cii.NewSerializer().Serialize`, `zugferd.EmbedXML` (through temporary files) |
| `goblde` | invopop/gobl.de.xinvoice v0.0.10 | `xinvoice.Parse`, `Envelope.Validate`, `xinvoice.Convert` |
| `angelo` | angelodlfrtr/go-invoice-generator v0.8.0 | `facturx.Extract`, `facturx.Parse`, `Parsed.Verify` |
| `papierleicht` | papierleicht/go-e-rechnung-validator v0.1.0 | `validator.ExtractXML`, `Validator.Validate`, `Validator.ValidatePDF` |
| `dcnws` | DCNWS/invoicekit.go v0.2.13 | `facturx.Extract`, `einvoice.Parse`, `einvoice.Validate`, `einvoice.GenerateCII`, `facturx.Embed`, `hybrid.Check` |
| `pdf0` | mgilbir/pdf0 v0.9.1 + formalis v0.3.1 | `pdf0.Read` + `pdf0.ValidateFacturX`, `formalis.Validate`, `pdf0.EmbedFacturX` + `Document.Write` |
| `litho` | gotenberg/litho v0.0.18 | `pdf.Open` + `EmbeddedFiles`, `ValidateFacturXReport`, `Writer.RemediateFacturX` |
| `kanjo` | cyprienbrisset/kanjo v0.3.2 | `pdfa.ExtractInvoiceXML`, `read.ReadBytes`, `rules.Validate(doc, "en16931")`, `write.WriteBytes`, `pdfa.EmbedXML` |
| `taha` | tahayassine/facturx-go v0.1.0 | `facturx.ExtractXML`, `cii.Decode`, `en16931.Validate`, `cii.Encode` |
| `fx` | solarpush/fX v1.4.0 + 6 commits (`1314f3d`) | `pdf.ExtractXML`, `cii.Parse`, `invoice.ValidateReport`, `cii.Generate` |
| `gopdfattach` | MarlinKuhn/gopdfattach v0.0.3 | `gopdfattach.Extract`, `gopdfattach.AttachFacturX` |
| `pdfcpu` | pdfcpu v0.16.4 | `api.ExtractAttachmentsRaw`; embedding through a low-level recipe written for this bench (FileSpec, `/AF`, XMP) |
| `facturx` | ZeRechnung, at the version in `facturx/go.mod` or from `ZERECHNUNG_DIR` | `facturx.Extract`, `cii.Parse`, `validate.XML`, `cii.Marshal`, `facturx.Embed`, `facturx.Check` |

## Known specifics

- **`Embed_Small` fails for andee and kanjo, on purpose.** On the host with object streams, these libraries return a 23 to 25 KB PDF that cannot be read back (the host is 136 KB). `TestEmbedOutputSize` prints the input and output sizes. The script carries on and records the rest of the measurements.
- **fX is built from a local clone.** `go get` rejects the module archive of solarpush/fX (a file name contains a carriage return); `run.sh setup` clones the repository into `.deps/` and `fx/go.mod` references it with `replace`.
- **pdf0 does not expose extraction on its own.** The XML is only reachable through `ValidateFacturX`, which checks the whole PDF: pdf0 therefore appears in `CheckPDF`, not in `Extract`.
- **kanjo is restricted to the `en16931` rule set.** The default set includes the French rules and rejects the 26 German examples.
- **pdfcpu writes a `config.yml`** and certificates into the user's configuration directory (`~/Library/Application Support/pdfcpu` on macOS), unless `model.ConfigPath = "disable"`. Each module that depends on pdfcpu sets it in an `init()` (`noconfig_test.go`). gopdfattach tries to do it itself, but with the value `"disabled"`, which pdfcpu does not recognize.
- **papierleicht compiles its rules on the first call** (about one second) into a temporary cache folder; the benchmark creates the validator only once.
- `TestDetail` (in several modules) lists the rules that reject valid XML files.
- **What the conformance cases show** ([`../cases`](../cases/README.md)): andee applies the BR-DE rules to every document, so it rejects every plain EN 16931 invoice; kanjo has no XRechnung rule set (`cius.fr`, `en16931`, `kanjo`, `orderx`); papierleicht answers `UNSUPPORTED` for a plain EN 16931 UBL invoice; taha and angelo read CII only; fX rejects any document without a Factur-X profile (`PROFILE REQUIRED`); gobl.de rejects the test VAT identifiers of the cases (`GOBL-DE-TAX-IDENTITY-01`).

## Results of 2026-11-22

Apple M4 Pro on battery, high performance mode (`powermode 2`), Go 1.27.4 for every module. ZeRechnung is measured at commit [`3a3b11a`](https://github.com/ritoon/zerechnung/commit/3a3b11a), over 6 runs; the other libraries over 3. litho moves from v0.0.13 to v0.0.18. The tables are in [`results/2026-11-22/tables.md`](results/2026-11-22/tables.md), with, for each operation, the bytes and the number of allocations of one call (`go test -benchmem`); the analysis is on the [Benchmarks](https://github.com/ritoon/zerechnung/wiki/Benchmarks) wiki page.

- The package is the fastest on every operation, and it allocates the least on each one: 28 KB in 104 allocations to parse a 13 KB invoice, when tokenization alone with `encoding/xml` takes 31 KB and 711; a single 8 KB allocation to write it; 226 KB to embed it into a PDF, against 1.2 MB for pdf0.
- The full check of a PDF takes 832 µs, against 54.8 µs on 2026-11-16: since PR ritoon/zerechnung#73, it checks the PDF/A requirements across the whole file (every object, page content, fonts, JPEG 2000 images). DCNWS, the next one, takes 1.11 ms.

Both embedding benchmarks (`Embed_Small`, `EmbedClassic_Small`) were measured again the same day, after the host PDFs were fixed, still in the ZeRechnung tree (commit [`8d728bd`](https://github.com/ritoon/zerechnung/commit/8d728bd), same package code as `3a3b11a`). `host.pdf` and `host_classic.pdf` had lost their `/EmbeddedFiles` tree, but their catalog kept in `/AF` the file specification of the former `factur-x.xml`, and their XMP the Factur-X schema. They are now PDF/A-3 documents with no attachment, no `/AF` and no Factur-X XMP, which `TestHostsStrict` (pdfcpu module) checks. The embedding lines of the `bench_*.txt` files come from this new run (local time, UTC+1): the package over 6 runs of its whole module, from 19:35 to 19:37; litho alone, at 19:22; the other modules in a full pass, from 17:56 to 18:05. The other benchmarks do not use the hosts.

On the fixed hosts, the package embeds in 92.0 µs (xref stream) and 89.4 µs (classic xref), against 98.7 and 96.2 µs, and allocates 226 KB in 376 allocations instead of 245 KB in 456 (xref stream). The other embedders also allocate a little less (pdf0 1.2 MB instead of 1.3 MB, litho 383,412 allocations instead of 394,412) and are up to 10% faster (kanjo, pdf0), except DCNWS (+5%): the host is smaller (136 KB instead of 144), but these measurements were taken at other times than the campaign, on a machine where other work was running (disturbed steps were run again). andee and kanjo still fail on `Embed_Small`.

### Confidence by country

Conformance cases of [`../cases`](../cases/README.md), run on 2026-11-22 evening with the module versions of this campaign (ZeRechnung at the version of `facturx/go.mod`, commit [`724b388`](https://github.com/ritoon/zerechnung/commit/724b388); the cases come from its commit `9b9f489`). Score out of 100: 100 × √(S × D), S the share of valid cases accepted, D the mean detection of the expected rules (1 raised, 0.5 rejected otherwise, 0 accepted); Germany is the mean of four parts, EN 16931 and XRechnung 3.0.2, each in CII and UBL. Reference: KoSIT validator 1.6.3 with the XRechnung 3.0.2 configuration of 2026-08-31. Details, rules raised by family: [`results/2026-11-22/tables.md`](results/2026-11-22/tables.md).

| Library | Germany | EN 16931, CII | EN 16931, UBL | XRechnung, CII | XRechnung, UBL | Valid accepted | Invalid rejected |
|---|---|---|---|---|---|---|---|
| **ZeRechnung** | **99** | 96 | 100 | 100 | 100 | 207/208 | 351/351 |
| pdf0 | **93** | 92 | 93 | 92 | 95 | 204/208 | 301/351 |
| DCNWS | **81** | 98 | 100 | 64 | 63 | 193/208 | 236/351 |
| papierleicht | **63** | 83 | 0 | 88 | 82 | 153/208 | 348/351 |
| kanjo | **61** | 91 | 58 | 47 | 47 | 163/208 | 193/351 |
| andee | **37** | 0 | 0 | 74 | 74 | 164/208 | 265/351 |
| gobl.de | **24** | 73 | 0 | 11 | 13 | 19/208 | 341/351 |
| taha | **20** | 39 | 0 | 40 | 0 | 82/208 | 209/351 |
| angelo | **11** | 0 | 0 | 43 | 0 | 29/208 | 350/351 |
| fX | **0** | 0 | 0 | 0 | 0 | 0/208 | 351/351 |

- **The cases come from ZeRechnung's own mutation tests**: its score is an upper bound. It rejects one valid case, a CII invoice with `TestIndicator` (CII-SR-002, a warning for the reference), which its Factur-X EN16931 schema does not allow.
- **pdf0** (formalis) runs the official rules for each syntax and CIUS and raises the BR-DE and Peppol rules; it misses part of the BR-CO and BR-CL rules.
- **DCNWS** applies EN 16931 in both syntaxes almost like the reference, but no BR-DE rule: it scores 63 to 64 on XRechnung.
- kanjo and gobl.de reject the 13 and 14 valid EN 16931 CII cases under BR-CO-25, which no CEN artifact checks; these rejections are not counted against them (see [`../cases`](../cases/README.md#the-score)).

## Results of 2026-11-16

The `bench_*.txt` files in `results/2026-11-16/` were produced before the harness moved into the ZeRechnung repository: the module path they print is `example.local/fxbench/…`, the measured code is the same. The `survey_*.txt` files were regenerated from this folder and give the same figures. For angelo, kanjo, fX and gopdfattach, the measurements predate turning off pdfcpu's `config.yml`; a new measurement of `Extract_Small` after the fix differs by less than 7%, without changing the ranking.

The `bench_facturx.txt` and `survey_facturx.txt` files measure the first version of ZeRechnung, on the same machine (Apple M4 Pro) and the same day. Its validator applies by default the schema, the EN 16931 rules and the Factur-X profile rules, without the French rules; the 6 XML files it rejects are the ones papierleicht also rejects, which runs the official schematrons (CEN and KoSIT): BR-CL-25 twice, CII-SR-465 four times (on the XRECHNUNG examples).

## Adding a library

1. Create a folder with a `go.mod` that depends on `github.com/ritoon/invoice-bench/go/common` through `replace … => ../common`.
2. Write `extract`, `parse`, `validate`, `write` or `embed` functions according to what the library can do, and pass them to the functions of `common`.
3. Add a `TestSurvey`, then the folder name to `MODULES` in `run.sh`.

## Licenses

The code in this folder is under the MIT license, like the repository. The modules download third-party libraries under their own licenses, including AGPL-3.0 (dotwave, litho, fX), BSL 1.1 (kanjo) and MPL-2.0 (papierleicht); they are only used by these benchmarks and are not imported by the `zerechnung` package. The ZUGFeRD corpus (Apache-2.0) is downloaded, not versioned here.
