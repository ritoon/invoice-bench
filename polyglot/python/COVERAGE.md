# Python bench: selected libraries, excluded libraries and coverage of the standard

Python bench of the common protocol (`../PROTOCOL.md`), compared with the Go
package `github.com/ritoon/zerechnung`. Surveys as of 2026-11-17.

- **Interpreter**: CPython 3.13.18 (Homebrew), venv `.venv/`, no uv.
  Dependencies frozen in `requirements.txt` (output of `pip freeze`),
  installed as wheels only (`--only-binary=:all:`) except `tagic`
  (a pure-Python dependency of getafix, published as an sdist only; pyproject
  reviewed: poetry-core backend, no `setup.py`, no code executed at
  install time).
- **Official artifacts** downloaded by `./run.sh setup` into `.artefacts/`
  (ignored by git), SHA-256 checked in `run.sh`:
  CEN EN 16931 CII **1.3.17** (`en16931-cii-1.3.17.zip`, release
  `validation-1.3.17` of ConnectingEurope/eInvoicing-EN16931) and KoSIT
  **XRechnung 3.0.2 / schematron 2.6.0** (`xrechnung-3.0.2-schematron-2.6.0.zip`,
  itplr-kosit/xrechnung-schematron).
- **No Java**: no adapter starts a JVM. The only XSLT 2.0 engine
  is `saxonche` (SaxonC-HE compiled as a GraalVM native library:
  `saxonche.dylibs/libsaxonc-core-he.13.1.1.dylib`, 97 MB).
- **No network** during `survey` and `bench`: the only network calls
  found in the code of the measured libraries (factur-x: schematron through a
  Saxon HTTP server; mcp-einvoicing-de: remote KoSIT validator) are never
  exercised (see the library notes).

How to run: `./run.sh setup`, then `./run.sh survey [lib…]` and
`./run.sh bench [lib…]` (variables `OUT`, `COUNT`, `BENCHTIME`; also
`CHECK_TIMEOUT` = 30 s for a benchmark's check call,
`SURVEY_TIMEOUT` = 10 s per survey file, `CHILD_TIMEOUT` = 1800 s per
child process). Each benchmark runs in a **fresh child process**, so
`first_call_ns` is a true first call (imports already done, library caches
cold). The survey starts one process per library.

The `package/file.py:line` paths below are relative to
`.venv/lib/python3.13/site-packages/`.

## Survey results (`./run.sh survey`)

| Library | validate (flagged / 26) | parse (/ 26) | extract (/ 104) |
|---|---|---|---|
| factur-x 7.10 | 0 (XSD only) | 26 | 104 (97 with `check_xsd=True`, the default) |
| drafthorse 2025.2.0 | 0 (XSD only) | 26 | — |
| pycheval 0.4.0 | — | 7 | 65 (39 infinite loops cut off) |
| getafix 0.5.0 | 6 (six exceptions, no rule detected) | 20 | 104 |
| zugferd-pdf 0.1.2 | — | — | 104 |
| erechnung 0.3.0 | — | — | — |
| mcp-einvoicing-de 0.11.4 | 2 (FX-SCH-A-000571 ×2) | 20 | 104 |
| saxonche + official XSLT | **6** (BR-CL-25 ×2, CII-SR-465 ×4) | — | — |

"—": the operation does not exist in the library (nothing is written to
`survey_python.jsonl`).

**Reference 6/26 reproduced** by `saxonche-schematron`, and explained: with
CEN 1.3.17 alone, the four `XRECHNUNG_*.cii.xml` raise only CII-SR-465,
which is `warning` level in the CEN schematron. The KoSIT validator
(validator-configuration-xrechnung v2026-08-31, CHANGELOG "custom level
error for CII-SR-465 and CII-SR-466") raises CII-SR-465/466 to `error`
level and lowers CII-SR-475/476 to `information` level in the CII
scenario. The adapter applies these levels when BT-24 designates XRechnung
(`adapters/saxonche_schematron.py`, `_KOSIT_CUSTOM_LEVEL`); the Go package
also counts CII-SR-465 as an error (`internal/rules/en16931/syntax.go`).
The two BR-CL-25 (electronic address scheme `EM` not in the EAS list) are
on `EN16931_ElektronischeAdresse` and `EN16931_RechnungsUebertragung`.
Sensitivity check: a modified GrandTotal yields BR-CO-15 and BR-CO-16.
The corpus file `not_validating_full_invoice…` raises no fatal/error with
the official XSLT (nor with the Factur-X XSD of factur-x and
drafthorse): its name comes from the corpus.

Details of the deviations from the reference:

- **getafix 6/26**: the six flags are parsing exceptions, not rules. Four
  `XRECHNUNG_*` rejected ("not a valid Profile"),
  `EN16931_Betriebskostenabrechnung` (empty `ValueError`),
  `EN16931_Physiotherapeut` (`TypeError: ProcuringProject.__init__() missing
  … 'id' and 'name'`). The two BR-CL-25 are not seen (no BR-CL rule in
  getafix).
- **mcp-einvoicing-de 2/26**: FX-SCH-A-000571 (FeRD equivalent of BR-CL-25)
  on the same two files as the reference. The XRECHNUNG files pass:
  CII-SR-465 stays a warning, and the library counts only errors (the
  library's verdict, consistent with its documentation).
- **factur-x, drafthorse 0/26**: XSD validation only, and all 26 files
  conform to the XSD.
- **parse**: pycheval rejects EXTENDED and XRECHNUNG (`pycheval/parse.py:356-361`),
  TypeCodes 204, 387, 389, 575, SubjectCode ADU, the `EM` scheme, and
  empty elements (`BuyerAssignedID`, `LineTwo`,
  `AttachmentBinaryObject`). mcp-einvoicing-de rejects 5 sample German VAT
  numbers with a wrong check digit (DE123456789, DE234567891) and a
  negative allowance (strict pydantic model). getafix: the same six files
  as for validate.
- **extract**: pycheval loops forever as soon as the first attachment is
  not named `factur-x.xml` (`pycheval/pdf_extract.py:30-32`,
  `while doc: if doc[0] == FACTURX_FILENAME: break` without advancing); the
  harness cuts the call off after `SURVEY_TIMEOUT`. factur-x with its
  default options (`check_xsd=True`) rejects 7 extracted XML files that do
  not conform to the XSD.

## Library notes: selected libraries

### factur-x 7.10 (Akretion)

BSD, latest release 2026-10-30, 411,142 downloads last month
(pypistats). The reference library of the ecosystem (Odoo).

- Parsing: `parse_ubl_cii_xml` (`facturx/parse_xml.py:49`, XSD by default
  at line 56) into a dictionary; writing: `generate_cii_xml`
  (`facturx/generate_xml.py:2351`) and `generate_ubl_xml` (line 3607).
- PDF: `get_xml_from_pdf` (`facturx/facturx.py:662`),
  `generate_from_binary` (line 1425, PDF/A-3 XMP, no PDF/A conversion,
  goes through a temporary file).
- Bundled XSD: Factur-X **1.09** (5 profiles,
  `facturx/xsd_and_schematron/facturx-*/Factur-X_*.xsd`), UBL 2.1,
  ZUGFeRD 1.0, CII D22B (extended CTC FR).
- Bundled schematrons: Factur-X **1.09.2** (FNFE-MPE/FeRD, 5 profiles;
  the EN 16931 XSLT counts 200 BR identifiers, 14 CII-SR and 218 FX-SCH),
  BR-FR-Flux2 CII and UBL (124 BR-FR identifiers), CEN UBL **1.3.16**.
  But `xml_check_schematron` runs them **only through an HTTP POST to a
  Saxon server** (`facturx/facturx.py:495` `requests.post`, default URL
  `http://localhost:5000/transform` at line 179; the codedb files are served
  from raw.githubusercontent.com, lines 182-186). That server is not
  provided (it needs a JVM): offline, validation comes down to the XSD.
- `generate_cii_xml` has `check_schematron="base"` by default
  (`generate_xml.py:2355`): the adapter passes `check_schematron=False`
  so that nothing is sent over the network.
- No native hybrid PDF check (`CheckPDF` not supported).

### drafthorse 2025.2.0 (pretix)

Apache-2.0, latest release 2025-09-15, 156,671 downloads/month.

- Full object model, `Document.parse` and `Document.serialize`
  (`drafthorse/models/elements.py:70-78`: serialization revalidates the
  output against the XSD, which counts toward `Write_*`).
- Factur-X **1.0.07** XSD for the 5 profiles (`drafthorse/schema/Factur-X_1.0.07_*.xsd`,
  `drafthorse/utils.py:14-15`), no schematron.
- `attach_xml` (`drafthorse/pdf.py:47`): PDF/A-3 XMP, reuses the host's
  OutputIntents; no PDF reading (README).

### PyCheval 0.4.0

Apache-2.0, latest release 2026-11-16, 3,561 downloads/month.

- The README announces Factur-X 1.09.2, the code declares
  `FACTURX_VERSION = "1.0.07"` (`pycheval/__init__.py:19`).
- Models MINIMUM, BASIC WL, BASIC, EN 16931; EXTENDED and XRECHNUNG
  rejected at parse time (`pycheval/parse.py:356-361`). 52 structural
  checks `raise ModelError` in `pycheval/model.py`, no BR rules.
- Extraction: the infinite loop described above; also filters on
  `/Subtype /text/xml`.
- Check of the AFRelationship by country (fr, de) when reading a PDF
  (`pycheval/pdf_parse.py:85-110`): the only container check.
- Embedding: the public API (`embed_facturx_file_in_pdf`,
  `pycheval/pdf_embed.py:77`) reads the XML from a file; the adapter
  calls `_embed` (line 133) with the same default options
  (Data relationship) to stay in memory.

### getafix 0.5.0

Apache-2.0, latest release 2026-11-10, 2,092 downloads/month.

- Typed BT/BG model, CII parsing and writing (stdlib ElementTree).
- Profiles: enum "Factur-X 1.08 / ZUGFeRD 2.4" (`getafix/schema/types.py:6-7`),
  5 profiles, no XRECHNUNG.
- Hand-coded rules: 90 distinct BR-* identifiers in the package
  (counted from the AST), including **89 of the 223 BR rules of CEN 1.3.17**
  (after normalizing BR-05/BR-5 and BR-IG/BR-AF), **no BR-CL**, no
  CII-SR/CII-DT; plus 16 BR-FXEXT-* (EXTENDED overlay) and 11 GETAFIX-*
  (`getafix/rules/*.py`). `ValidationWarning`s do not count.
- PDF: `extract_xml(pdf: Path)` (`getafix/pdf.py:92`) only accepts a
  path (reading the file is inside the measurement); `attach_xml`
  (line 45) adds a generic attachment and, according to its documentation
  (lines 18-20), does not make the document PDF/A-3.

### zugferd-pdf 0.1.2 (module `zugferd`)

MIT, latest release 2026-08-25, 97 downloads/month.

- Writing from an in-house model (`generate_facturx_xml`), extraction
  (`zugferd/extractor.py:13`), PDF/A-3U packaging
  (`zugferd/pdfa_packager.py:140`: adds an sRGB OutputIntent and a PDF/A-3U
  XMP, AFRelationship Alternative).
- `validate_zugferd_pdf` (`zugferd/validator.py:78`) is the only native
  "CheckPDF" in the bench, but it is home-made: XMP, OutputIntent,
  /AF, BT-24 list, BR-DEC-23 and BT-112 decimals; no schematron
  despite the README. False positives: `pdfaid:part` is only looked for as
  an element (line 124), so the official `EN16931_Einfach.pdf`
  (XMP in attributes) is flagged "Missing pdfaid:part" and
  "Missing /OutputIntents" (line 184).
- No XML-to-model parsing: `Write_*` starts from a hand-built model
  (`adapters/zugferd_pdf.py`, `hand_model`).

### erechnung 0.3.0

MIT, latest release 2026-08-28, 116 downloads/month.

- Generation only: DIN 5008 PDF drawn by reportlab, CII XML EN 16931
  or XRechnung 3.0 (`erechnung/xml.py:16-18`), then attached by factur-x
  (`erechnung/generator.py:7`, 72).
- Measured: `calculate` + `xml_bytes`, the XML part of the public path
  `Invoice.zugferd_file`, on a hand-built model
  (`adapters/erechnung.py`). `erechnung/validation.py` only checks the
  Python input data (IBAN, country, dates…).

### mcp-einvoicing-de 0.11.4 (+ mcp-einvoicing-core 1.35.0)

Apache-2.0, latest release 2026-09-12, 287 downloads/month (core:
1,582). An MCP server; only the library layer is used.

- Parsing / writing: `ZUGFeRDCIIParser`, `ZUGFeRDCIISerializer`,
  `XRechnungUBLParser/Serializer` (`mcp_einvoicing_de/serializers.py:26-95`),
  strict pydantic model (check digit of German VAT numbers). The
  serializer drops attachments (`Write_Large` cannot be compared
  byte for byte).
- Local validation `invoice_validate`: bundled XSLT run by
  saxonche (`mcp_einvoicing_de/validators/schematron.py:69-84`,
  `tools/invoice_validate.py:172-195`). EN 16931 profile: FeRD XSLT
  Factur-X **1.09.2** (file header; the docstring at line 10 still says
  1.08). XRECHNUNG: CEN EN 16931 CII (compiled by SchXslt 1.10.1, from the
  KoSIT 2026-08-31 configuration, hence CEN **1.3.16**) then
  XRechnung-CII-validation.xsl "Schematron Version 2.6.0 / XRechnung
  3.0.2", same rule identifiers as the official 2.6.0 artifact
  (39 BR-DE, 7 BR-DEX, 6 BR-TMP, 22 PEPPOL-EN16931) but the file is not
  byte-for-byte identical. Only fatal/error results count.
- Cost: `_run_one_stylesheet` (`tools/invoice_validate.py:278`) builds
  a new validator on every call, hence a new `PySaxonProcessor`
  and a recompilation of the stylesheet (`mcp_einvoicing_core/schematron.py:381-384`):
  ~0.1 to 0.2 s per validation, which is the real cost of the API.
- Network: remote KoSIT validation (`validators/kosit.py`, httpx) has been
  disabled by default since 0.11 (`cloud_validate=False`,
  `tools/invoice_validate.py:40` reads `EINVOICING_DE_KOSIT_ENABLE`).
  The adapter forces this variable to `0` before the import and checks
  `validator_used == "local_schematron"` on every call.
- No use of the XSD found in the DE package; no hybrid PDF
  check (`PDFEmbedder.identify` only reads the XMP).

### saxonche 13.1.1 + official XSLT (reference validator)

SaxonC-HE, MPL-2.0 according to Saxonica (empty license field on PyPI),
latest release 2026-10-27, 1,978,505 downloads/month. Native, no JVM.

- CEN EN 16931 CII 1.3.17 XSLT: 223 BR identifiers (including 23 BR-CL and
  21 BR-DEC), 486 CII-SR, 103 CII-DT. Plus the KoSIT XRechnung 2.6.0 XSLT
  when BT-24 designates XRechnung, with the KoSIT custom levels
  (see above).
- Both stylesheets are compiled on the first call (0.15 to 0.3 s, included
  in `first_call_ns`), then 1 to 15 ms per bench invoice.
- Only `ParseValidate_*` (validator only).

## Benchmark status

Trial run `COUNT=1 BENCHTIME=0.2 ./run.sh bench` (2 min 30 s for the
8 libraries × 16 benchmarks; the durations are not meant to be quoted):

| Library | ok | unsupported | fail |
|---|---|---|---|
| factur-x | 14 | CheckPDF ×2 | — |
| drafthorse | 10 | Extract ×4, CheckPDF ×2 | — |
| pycheval | 6 | ParseValidate ×3, CheckPDF ×2 | Extract_Large, Extract_Extended (infinite loop, cut off at 30 s), Parse_Lines and Write_Lines (`BuyerAssignedID has no text`), Parse_Large (XRECHNUNG profile rejected) |
| getafix | 12 | CheckPDF ×2 | Parse_Large, ParseValidate_Large (XRECHNUNG profile rejected) |
| zugferd-pdf | 9 | Parse ×3, ParseValidate ×3, Write_Lines | — |
| erechnung | 1 (Write_Small) | 15 | — |
| mcp-einvoicing-de | 14 | CheckPDF ×2 | — |
| saxonche-schematron | 3 (ParseValidate) | 13 | — |

Things to know when reading the numbers:

- `Write_*` starts from the model produced by the library's parser (parsing
  outside the measurement). zugferd-pdf and erechnung have no parser: `Write_Small`
  starts from a hand-built model based on `EN16931_Einfach`
  (`note: "modèle construit à la main"`), `Write_Lines` is `unsupported`.
- `Embed_*` is read back with the library itself when it can extract,
  otherwise with factur-x (`note: "relu avec factur-x"`, the drafthorse case).
- factur-x `Parse_*` includes XSD validation (default option), hence
  `Parse_Large` being slower than `ParseValidate_Large` (XSD only).
- drafthorse `Write_*` revalidates the output against the XSD (behavior of
  `serialize`).
- getafix `Extract_*` reads the file from disk inside the measurement
  (path-based API, OS cache).
- mcp-einvoicing-de `ParseValidate_*` recompiles the XSLT stylesheet on every
  call (~0.12 to 0.2 s/op); saxonche-schematron compiles once.
- The only native `CheckPDF` is zugferd-pdf's, with home-made checks.

## Library notes: excluded libraries

| Package | Reason |
|---|---|
| en16931 0.2 (b2brouter) | UBL only; validation through the remote validex.net API |
| zugferd 0.1.1, xrechnung 0.1.2 | clients of the remote xrechnung.dev API |
| beliq 0.7.3, factora 0.1.2, invowerk 1.0.19, facturx-engine 0.1.0 | clients of remote APIs or of Docker containers |
| python-invoice-lint 1.0.0 | wrapper around saxonche + CEN downloaded at run time into `~/.cache` (duplicate of the reference validator), path-based API, recompiles on every call, "your-org" project URLs left in place, "Other/Proprietary" license; 26 downloads/month |
| pyfactx 0.7.0 | generation only, pins `saxonche<13`, self-validation through the Factur-X 1.07.2 XSLT; 20/month |
| facturxlib 1.0.1 | generation only, the simplified model handles only one VAT rate (EN16931_Einfach has two); 25/month; installed for a trial, then removed |
| pydantic-factur-x 0.2.0 | generation only, latest release 2023-12; 70/month |
| factur-x-nolog 1.7.3 | old fork of factur-x |
| pyschematron 1.2.5 | generic pure-Python Schematron engine (GPLv3+, 7,401/month), not a Factur-X library. Off-bench trial with the preprocessed CEN 1.3.17 `.sch`: correct verdicts (BR-CL-25 found), but 12 s for `EN16931_Einfach`, 33 s for `EN16931_ElektronischeAdresse` and 97 s for the full CEN file, against a few ms for saxonche: unusable for this bench |
| einvoice 0.0.1 | Indian invoice QR code |
| invoice2data | template-based text extraction from PDFs, no XML |
| python-ubl, peppol-py | UBL/Peppol only, or AS4 transport |
| Rust/C bindings | none found on PyPI (`core-invoice`, `zugferd-rs`, `facturx-rs`…: 404 on the PyPI JSON API; nothing in web and GitHub searches) |

## Coverage of the standard

| Library | Version | License | Latest release | Downloads | CII parsing | CII writing | UBL | PDF extraction | PDF embedding | XSD | EN 16931 (CII) | Factur-X profiles (MINIMUM, BASIC WL, EXTENDED) | BR-FR | BR-DE (XRechnung) | Peppol BIS 3 | PDF/A-3 container check | External dependencies |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| factur-x | 7.10 | BSD | 2026-11-19 | 411,142/month | yes (dictionary) | yes | yes (parsing, writing; CEN UBL 1.3.16 XSLT via Saxon server) | yes | yes (PDF/A-3 XMP, no conversion) | yes (Factur-X 1.09, 5 profiles) | FNFE-MPE Factur-X 1.09.2 XSLT bundled, run only by a Saxon HTTP server; offline: no | XSD for the 5 profiles; schematron for the 5 profiles via Saxon server | BR-FR-Flux2 XSLT (124 rules) via Saxon server | no | no | no | lxml, pypdf, requests; Saxon HTTP server (JVM) for any schematron |
| drafthorse | 2025.2.0 | Apache-2.0 | 2025-09-15 | 156,671/month | yes | yes (revalidated against the XSD) | no | no | yes (PDF/A-3 XMP) | yes (Factur-X 1.0.07, 5 profiles) | no | XSD for the 5 profiles | no | no | no | no | lxml, pypdf |
| pycheval | 0.4.0 | Apache-2.0 | 2026-11-16 | 3,561/month | partial (7/26: rejects EXTENDED, XRECHNUNG, several TypeCodes) | yes | no | partial (65/104, infinite loop) | yes (file-based public API) | no | no (52 structural checks) | MINIMUM, BASIC WL, BASIC, EN 16931; not EXTENDED | no | no (XRECHNUNG rejected) | no | partial (AFRelationship by country) | pypdf |
| getafix | 0.5.0 | Apache-2.0 | 2026-11-10 | 2,092/month | yes (20/26) | yes | no | yes (path only) | partial (attachment without XMP or PDF/A-3) | no | partial (89 of the 223 CEN BR rules, no BR-CL or CII-SR) | 5 profiles (Factur-X 1.08); EXTENDED: 16 BR-FXEXT rules | no | no | no | no | tagic, pypdf |
| zugferd-pdf | 0.1.2 | MIT | 2026-08-25 | 97/month | no | yes (in-house model) | no | yes | yes (adds an sRGB OutputIntent and a PDF/A-3U XMP) | no | no (BR-DEC-23, BT-112 decimals) | BT-24 list only | no | no | no | partial (home-made checks, false positives on XMP in attributes) | pikepdf, pypdf, lxml, reportlab |
| erechnung | 0.3.0 | MIT | 2026-08-28 | 116/month | no | yes (EN 16931, XRechnung 3.0) | no | no | no (draws its own PDF, attached by factur-x) | no | no (input data only) | no | no | no | no | no | reportlab, factur-x, pikepdf, lxml, pillow |
| mcp-einvoicing-de | 0.11.4 | Apache-2.0 | 2026-09-12 | 287/month (core 1,582) | yes (20/26, strict model) | yes (lossy) | yes (parser, serializer, CEN UBL and XRechnung UBL XSLT) | yes | yes (pikepdf, pdfaid XMP, sRGB OutputIntent) | no | FeRD Factur-X 1.09.2 XSLT (EN 16931) or CEN CII 1.3.16 (XRECHNUNG) via saxonche, offline | FeRD XSLT for the 5 profiles via saxonche | no | KoSIT XRechnung 3.0.2 XSLT, schematron 2.6.0 via saxonche | no (22 PEPPOL-EN16931 rules carried over by XRechnung) | no | saxonche (97 MB native), pikepdf, pydantic, httpx, fastmcp, cryptography; remote KoSIT optional, disabled |
| saxonche + official XSLT (reference) | 13.1.1; CEN 1.3.17; KoSIT 2.6.0 | MPL-2.0 | 2026-10-27 | 1,978,505/month | — | — | no (CII XSLT only) | — | — | no (CEN does not provide a CII XSD) | official CEN 1.3.17 XSLT via Saxon (223 BR, 589 CII-SR/CII-DT) | no | no | official KoSIT 2.6.0 XSLT via Saxon, KoSIT levels | no (22 PEPPOL-EN16931 rules via XRechnung) | no | native saxonche (no JVM); artifacts downloaded by `run.sh setup` |

Downloads: pypistats, last 30 days as of 2026-11-17. Licenses and
dates: PyPI JSON API.
