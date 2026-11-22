# Coverage: TypeScript / JavaScript libraries

Bench run according to `../PROTOCOL.md` on 2026-11-17 with Node v22.15.0. Versions are frozen in `package.json` and `package-lock.json`. Weekly downloads come from `api.npmjs.org/downloads/point/last-week`, week of 2026-11-09 to 2026-11-15. Licenses and latest versions come from `npm view`.

Every claim about a rule, an artifact, the network or Java was checked in the installed code. The paths cited are relative to `node_modules/` and match the frozen versions.

## Summary

13 libraries are measured. The bench columns give, for each family of operations:
- `ok`: the operation is measured;
- `—`: the library cannot do it (`unsupported`);
- `fail`: the library throws an exception or returns a wrong result.

The survey (`survey_typescript.jsonl`) counts:
- validate: the valid files the library flags, out of 26 (official schematron reference: 6);
- parse: the CII files parsed, out of 26;
- extract: the PDFs whose XML is extracted, out of 104.

| Library | Extract ×4 | Parse ×3 | ParseValidate ×3 | Write Small / Lines | Embed ×2 | CheckPDF ×2 | validate | parse | extract |
|---|---|---|---|---|---|---|---|---|---|
| @stackforge-eu/factur-x 1.4.7 | ok | — | ok (XSD only) | ok (hand-built) / — | ok | — | 0/26 | — | 104/104 |
| @e-invoice-eu/core 3.4.4 | — | — | — | ok (hand-built) / — | — | — | — | — | — |
| node-zugferd 0.1.1-beta.1 | — | — | — | ok (hand-built) / — | — | — | — | — | — |
| @stafyniaksacha/facturx 0.6.1 | ok | ok | ok | ok / ok | **fail** | — | 12/26 | 26/26 | 100/104 |
| factur-x-kit 0.3.2 | — | ok / **fail** ×2 | ok / **fail** ×2 | ok / **fail** | — | — | 18/26 | 13/26 | — |
| @lucaapp/zugferd-ts 0.1.0 | ok | ok ×2 / **fail** (Large) | ok | ok / ok | ok | — | 8/26 | 19/26 | 104/104 |
| @attestwire/en16931 0.14.1 | ok | ok | ok | ok / ok | — | ok | 2/26 | 26/26 | 104/104 |
| @ekwo-ai/factur-x 0.27.0 | ok | ok | ok (8 rules) | ok (hand-built) / — | ok | — | 3/26 | 26/26 | 104/104 |
| @conformo/toolkit 0.2.4 | ok | ok | ok | ok / ok | ok | ok | **6/26** | 26/26 | 104/104 |
| @kontor-mcp/core 1.0.4 | ok | ok | ok | — | — | — | 2/26 | 26/26 | 104/104 |
| @sweberdev/summand 2.0.1 | ok | — | ok | — | — | ok | **6/26** | — | 104/104 |
| en16931 0.1.0 | — | — | ok | — | — | — | **6/26** | — | — |
| @fin.cx/einvoice 11.3.1 | ok | ok | ok | ok / ok | ok (no XMP) | — | 0/26 | 26/26 | 104/104 |

Three libraries match the reference of 6 flagged files exactly: conformo, summand and en16931 (BR-CL-25 twice, CII-SR-465 on the four XRECHNUNG files). en16931 also flags CII-SR-475 and CII-SR-476 on two of those files: it does not read severities.

## Method specific to the TypeScript bench

- **One process per library and per benchmark.** `run.sh` starts `node bench.mjs <lib> <Benchmark>` for each pair, so `first_call_ns` is a true cold first call: lazy loading, schematron compilation, JIT. The survey runs in one process per library (`node survey.mjs <lib>`).
- **Console silenced.** Several libraries write to the console on every call:
  - @fin.cx/einvoice prints four lines per extraction;
  - stackforge warns about the missing OutputIntent;
  - node-zugferd reports that `strict` is disabled.

  The console is therefore replaced with empty functions in the measurement processes (`lib/harness.mjs`, `silenceConsole`). The JSONL results go straight to stdout.
- **String or bytes.** The XML is passed in whatever form each library's API takes:
  - text (decoded as UTF-8 outside the measurement) for stafyniaksacha, lucaapp, ekwo, conformo, en16931, fin.cx, stackforge;
  - bytes (`Buffer`) for attestwire, kontor, summand, factur-x-kit.

  PDFs are always passed as bytes.
- **Asynchronous.** A function that returns a promise is awaited in the loop (`await`). The cost of a microtask turn is therefore included, as it would be for a real caller.
- **Hand-built models.** `Write_*` starts from the model produced by the library's parser when there is one. Four libraries have no reader, or have a reader whose model cannot be written back: stackforge, e-invoice-eu, node-zugferd and ekwo. For them, the `EN16931_Einfach` model is written by hand in `lib/handmodels.mjs`, with the note `"hand-built model"`.
  - Only `Write_Small` is measured; `Write_Lines` is `unsupported`.
  - The four outputs were run through summand: 0 errors against the XSD and the CEN schematron.
- **Reading back `Embed`.** The PDF produced is read back with the library's own extractor. Every library that embeds has one.
- **SaxonJS.** conformo, kontor and stafyniaksacha run schematrons compiled to SEF with `saxon-js` 2.7, under a proprietary Saxonica license: free of charge, but not open source. It is the only dependency of this kind.
- **Install scripts.** `npm ci --ignore-scripts`, also enforced by `.npmrc` (which also sets `legacy-peer-deps`, the setting the lockfile was built with). No measured package needs an install script. Two are blocked:
  - the `postinstall` of @fin.cx/einvoice, which downloads schematrons from GitHub `master` branches (see its library note);
  - that of `xsd-schema-validator`, an optional dependency of node-zugferd, which compiles a Java class with `javac`.

## Library notes: measured libraries

### @stackforge-eu/factur-x 1.4.7 (EUPL-1.2, 94,466/week)

Generates Factur-X / ZUGFeRD / XRechnung XML from its own input structure (`FacturXInvoiceInput`, `buildXml`). It also embeds into PDF/A-3 and extracts the XML, with `pdf-lib`. It has no CII reader.

- **Validation.** `validateXsd(xml, Profile.EN16931)`: XSD only, run by `libxml2-wasm`. The schema is re-read and recompiled on every call (`@stackforge-eu/factur-x/dist/index.js:1139`). The XSDs are those of Factur-X 1.09 / ZUGFeRD 2.5 for the five profiles (`@stackforge-eu/factur-x/schema/README.md`).
  - No schematron.
  - `validateInput` checks the input, not an XML document.
- **Embedding.** `embedFacturX({ pdf, xml, profile })` adds the attachment, the AFRelationship and the PDF/A-3b XMP. It adds an OutputIntent only if an ICC profile is supplied (`rgbIccProfile`); otherwise it warns that the PDF will not be PDF/A-3 compliant (`@stackforge-eu/factur-x/dist/index.js:1473`). The bench keeps the default options.
- **Network and Java.** No network calls, no Java.

### @e-invoice-eu/core 3.4.4 (WTFPL, 63,013/week)

Generation only: a JSON modeled on UBL (checked by an ajv schema) produces UBL, CII, XRechnung or Factur-X.

- **Parsing and validation.** No XML reader, no XML validation.
- **Factur-X.** The `Factur-X-*` formats attach the generated XML to a supplied PDF. When no PDF is supplied, they launch LibreOffice (`spawn`, `@e-invoice-eu/core/dist/e-invoice-eu.esm.js:7103`).
- **Measurement.** `Write_Small` starts from a hand-written JSON, and the ajv check is included in the call. `InvoiceService.generate` is slow: about 50 ms per call, versus 15 µs for stackforge.

### node-zugferd 0.1.1-beta.1 (MIT, 62,514/week)

Per-profile generation (MINIMUM, BASIC WL, BASIC, EN16931, EXTENDED) from a typed input. `embedInPdf` only attaches the XML the library has just generated.

- **Validation.** In `strict` mode, only the Factur-X 1.07.3 XSD (`node-zugferd/dist/Factur-X_1.07.3_*.xsd`) is checked, by `xsd-schema-validator`, which requires Java (`node-zugferd/dist/profile.js:40`).
  - Since Java is not installed, the measurement uses `strict: false`.
  - No schematron.
- **Versions.** The npm `latest` tag is 0.1.1-beta.1, marked deprecated in favor of 1.0.0-beta.1 (released on 2026-11-16 under the `beta` tag, not measured).

### @stafyniaksacha/facturx 0.6.1 (MIT, 7,339/week)

TypeScript port of akretion/factur-x (Python). It parses and writes CII to and from a typed model (`xmlToInvoice` / `invoiceToXml`), and extracts and embeds with `pdf-lib`. It also supports Order-X.

- **Validation.** `check({ xml, level, schematron: true })` runs:
  - the profile's Factur-X 1.09 XSD, through `libxml2-wasm`;
  - the profile's official Factur-X 1.09 schematron, compiled to SEF and run by SaxonJS.

  There is one SEF per profile, from MINIMUM to EXTENDED (`@stafyniaksacha/facturx/dist/shared/xsd/facturx/<profile>/FACTUR-X_<PROFILE>.sef.json.gz`).
  - This schematron is the one from the Factur-X package: it takes over the CEN rules and adds `PEPPOL-EN16931-R008` (empty elements) and the `FACTUR-X_EN16931_codedb.xml` code list.
  - Results carry no severity: every failed assertion counts. Hence 12/26 in the survey:
    - R008 on 6 files;
    - the code list on 6 files: electronic address scheme (`URIID`, the equivalent of BR-CL-25) on 2, XRechnung BT-24 missing from the Factur-X list on the 4 XRECHNUNG files, in addition to CII-SR-465.
  - The level is forced to `en16931`, because automatic detection throws "Unknown level: kosit" on the XRECHNUNG files.
- **Embedding.** `generate({ pdf, xml })` fails on both hosts: "Not implemented yet: Document ID already set". The error is thrown as soon as the host PDF's trailer has an `/ID` (`@stafyniaksacha/facturx/dist/shared/facturx.BQpPzQjc.mjs:184`), which is the case for most real-world PDFs.
- **Extraction.** Only recognizes `factur-x.xml` and `zugferd-invoice.xml`, hence the failure on the 4 XRECHNUNG PDFs (`xrechnung.xml`).

### factur-x-kit 0.3.2 (MIT, 5,067/week)

CII parsing and writing to and from a zod model, for the MINIMUM, BASIC WL, BASIC and EN16931 profiles. EXTENDED falls back to EN16931; there is no XRechnung profile.

- **Validation.** `validate()` applies the profile's zod schema and the BR rules rewritten in TypeScript: 131 distinct `BR-*` identifiers are cited in `factur-x-kit/dist/index.js`, plus in-house `BR-OWN-*` rules. No XSD or schematron.
- **Parsing.** The reader fails on 13 of the 26 CII files:
  - "Received Invalid XML Input" on 5 files, including `EN16931_Rabatte` (the Lines XML);
  - "Received Invalid Data Value" on 2 files;
  - a structural problem in the parsed data on 2 files;
  - "Profile not yet implemented" on the 4 XRECHNUNG files.
- **Extraction and embedding.** `FacturXPdf` is not exported: `FacturX.fromPDF` chains extraction and parsing (and only looks for `factur-x.xml`). `getPDF` rewrites the XML from the model. Extract and Embed are therefore `unsupported`.
- **Network.** The only `fetch` is used to decode a `data:` URL for a header image (`factur-x-kit/dist/index.js:5106`).

### @lucaapp/zugferd-ts 0.1.0 (Apache-2.0, 4,395/week)

Pure TypeScript: generation, PDF/A-3b embedding, CII and UBL parsing, validation.

- **Validation.** `validateXml` checks two things:
  - the profile's Factur-X 1.08 XSD, through `xmllint-wasm` (`@lucaapp/zugferd-ts/schemas/cii/Factur-X_1.08_*.xsd`);
  - the CEN EN 16931 CII 1.3.16 schematron, dated 2026-04-10. Its 806 `assert`s (`@lucaapp/zugferd-ts/schemas/schematron/en16931-cii.sch`) are run by `node-schematron` (fontoxpath).

  Things to know:
  - The library classifies every failed assertion as `error`, without reading the CEN's `flag="warning"` (`@lucaapp/zugferd-ts/dist/index.js:5828`). CII-SR-465/475/476 therefore count as errors.
  - False positives specific to the engine: BR-CO-16 on both Betriebskostenabrechnung files, BR-S-08 on Rabatte. Neither en16931 (same schematron, exact decimals) nor Saxon raises them.
  - Slowness: 3.5 to 10 s per validation. The schema is compiled only once (`getSchema`, cached); it is the evaluation that is slow.
  - UBL: only the UBL 2.1 XSD (`validateUblXml`). No KoSIT schematron (`serializeCiiXrechnung` writes XRechnung without validating it).
- **Container check.** `validatePdfA` applies 21 rules (`PDFA_RULES`) to the container alone, without the XML, hence CheckPDF `unsupported`.
- **Parsing.** Fails on 7 CII files whose document type (BT-3) is 204, 387, 575 or 877 (`@lucaapp/zugferd-ts/dist/index.js:4740`), including XRECHNUNG_Elektron (Parse_Large).

### @attestwire/en16931 0.14.1 (MIT, 2,538/week)

No dependencies: in-house XML parser, PDF reader and DEFLATE decoder. It parses CII and UBL, writes CII (`generateCii`) and cannot write PDFs.

- **Rules.** Rewritten by hand in TypeScript, with no schematron or XSD. The README announces "294 rules" (`@attestwire/en16931/README.md:75`). Distinct identifiers found in `dist/`:
  - EN 16931: 58 BR-nn, 19 BR-CO, 23 BR-CL, 96 VAT category rules, 21 BR-DEC;
  - XRechnung: 24 BR-DE;
  - Peppol: 27 PEPPOL-EN16931-R, applied to UBL only (`dist/rules.js:143`);
  - 27 in-house ATW-* rules (IBAN, BIC, SIREN/SIRET, Leitweg-ID…).
- **Hybrid PDF.** `validate(pdf)` adds six families of `AW-PDF-*` findings on the container: attachment, /AF, AFRelationship, MIME type, XMP, XMP profile / BT-24 (`dist/facturx-findings.js`). These are warnings, which do not count as errors.
- **Measurement.** Validation is fast (0.4 ms on Small). In the survey, 2/26 (BR-CL-25): it does not flag CII-SR-465, a CEN warning.

### @ekwo-ai/factur-x 0.27.0 (MIT, 521/week)

CII generation from a plain object, for five profiles. It has a CII reader (`readCii`), plus embedding and extraction with `pdf-lib` (entry point `/pdf`).

- **Validation.** `readCii(xml).violations` combines:
  - parsing anomalies: dates, IBANs with a wrong check digit (`@ekwo-ai/factur-x/dist/cii-read.js:161`);
  - 8 arithmetic checks BR-CO-10 to BR-CO-17 (`dist/received-checks.js`).

  In the survey, 3/26: sample IBANs with an invalid check digit.
- **Writing.** The model read by `readCii` (`ReceivedInvoice`) is not the one taken by `generateCiiXml` (`Invoice`). `Write_Small` therefore starts from a hand-written `Invoice`, which has only one note and no item GTIN.

### @conformo/toolkit 0.2.4 and its packages (Apache-2.0, 497/week)

Downloads of the other packages: validate 528, parse 345, pdf 160, mcp-server 560. The library is split as follows:
- `@conformo/parse`: CII and UBL parsing into the EN 16931 model, extraction with `pdf-lib` and `pdfjs-dist`;
- `@conformo/validate`: XSD through `xmllint-wasm`, schematrons compiled to SEF and run by SaxonJS;
- `@conformo/pdf`: PDF/A-3 finalization;
- `@conformo/toolkit`: facade, CII and UBL `serialize`, `validate` of an XML or a PDF.

The CLI and the MCP server are just wrappers.

- **Artifacts** (`@conformo/validate/artefacts/rulesets.json`):
  - CEN EN 16931 1.3.17 (CII and UBL);
  - KoSIT validator-configuration-xrechnung 2026-08-31 (XRechnung 3.0.2), with an `iban-mod97-decimal` fix;
  - Peppol BIS Billing 3.0.22 (UBL only);
  - BR-FR "French CTC Flux 2", source "factur-x 6.8 (FNFE-MPE)".
- **XSD**: CII D16B, Factur-X MINIMUM / BASIC WL / EXTENDED, UBL 2.1 (`artefacts/xsd/`).
- **Rule set selection** (`@conformo/validate/src/engine.js:7`):
  - EN 16931 except for MINIMUM, BASIC WL and EXTENDED;
  - XRechnung if BT-24 declares it;
  - BR-FR only if `country: "FR"`, which is not supplied here.
- **Severities.** Those of the KoSIT configuration (`levelOverrides`): CII-SR-465 is raised to `error` there, CII-SR-475/476 set to `information`. Hence exactly the reference 6/26.
- **Container.** `validate(pdf)` adds XMP / file name / profile notices, with no effect on `valid` (`@conformo/toolkit/src/container.js`). No PDF/A rules.
- **Embedding.** The `createPdf` recipe applied to the host: `finalizePDFA` with the toolkit's sRGB ICC profile and a fixed date.
- **Network.** `@conformo/validate/src/page.js` loads the SEF files via `fetch`, but only in the browser; that code is not used under Node.

### @kontor-mcp/core 1.0.4 (Apache-2.0, 337/week; @kontor-mcp/rules 554)

Library behind the Kontor MCP server, without the MCP part: detection, CII and UBL parsing, validation, conversion, PDF generation.

- **Validation.** CII D16B and UBL 2.1 XSD through `xmllint-wasm`, then the schematrons compiled to SEF and run by SaxonJS (`@kontor-mcp/rules/artifacts/MANIFEST.json`):
  - CEN EN 16931 1.3.16;
  - XRechnung 3.0.2 (schematron 2.5.0, KoSIT configuration 2026-01-31), with the D-019 fix on BR-DE-19.

  No Peppol or BR-FR. The in-house `plausibility` layer (22 `KONTOR-PLAUS-*` rules) is excluded from the count, just as the library itself excludes it from `valid`.
- **Writing and embedding.** `modelToCii` is not exported (exports limited to `.`), hence Write `unsupported`. `generateZugferdPdf` draws its own PDF from an input structure; it does not embed into a host.
- **Measurement.** In the survey, 2/26: it does not raise CII-SR-465.

### @sweberdev/summand 2.0.1 (MIT, 323/week)

Validator with no dependencies at all: in-house XSD validator, XPath 2.0 engine and Schematron engine. Parsing and writing are reserved for "Summand Pro".

- **Artifacts** (`@sweberdev/summand/NOTICE.md`, `src/rules/`):
  - CEN 1.3.16 converted to JSON, 806 asserts in `src/rules/en16931-cii.json`;
  - XRechnung 2.6.0, 52 asserts in `src/rules/xrechnung-cii.json`;
  - XRechnung 2.5.0 shipped, but not applied by default;
  - compiled CII and UBL XSD (`src/schemas/`).

  Severities from the KoSIT configuration: the library matches the reference 6/26 exactly.
- **Hybrid PDF.** `validateInvoice(pdf)` extracts the XML, compares the XMP profile with BT-24, then validates the XML. This is a native CheckPDF; there are no PDF/A rules.

### en16931 0.1.0 (MIT AND EUPL-1.2, 360/week)

Runs the CEF schematrons unmodified, on `fontoxpath` with exact decimals and `slimdom`. Validation only, CII and UBL.

- **Artifacts.** Copied from commit `b6c9e06` of the `validation-1.3.16` branch (`en16931/cef/VENDORED.md`). The XSDs and `codelist/` were removed, since the preprocessed schematron includes the lists.
- **Measurement.** The result (`fired`) carries no severity, so CII-SR-475/476 (CEN warnings) count. Compilation (`new SchematronRunner(loadCiiSchematron())`, about 170 ms) happens on the first call and lands in `first_call_ns`.

### @fin.cx/einvoice 11.3.1 (MIT, 323/week)

CII and UBL parsing and writing (Factur-X, ZUGFeRD, XRechnung, UBL…), with `@xmldom/xmldom`, `xpath` and `pdf-lib`.

- **Validation.** `EInvoice.validate(BUSINESS, { featureFlags })` runs only TypeScript validators (`@fin.cx/einvoice/dist_ts/einvoice.js:456`). Identifiers found in `dist_ts/`: 65 BR-nn, 23 BR-CO, 26 BR-CL, 84 VAT rules, 12 BR-DE, 2 PEPPOL.
  - In the survey, 0/26.
  - No XSD.
- **Schematron layer and network.** `IntegratedValidator` is separate from `validate()` and depends on files that the `postinstall` (`dist_ts_install/index.js`) downloads from GitHub `master` branches, with no pinned version (`dist_ts/formats/validation/schematron.downloader.js`). This postinstall is not run and this layer is not measured.
- **Embedding.** `PDFEmbedder.embedXml` only does a `pdf-lib` `attach` (AFRelationship Alternative), with no XMP or PDF/A metadata (`dist_ts/formats/pdf/pdf.embedder.js:25`).
- **Writing.** On XRECHNUNG_Elektron, `exportXml('facturx')` drops the base64 attachment: the output is 6.5 KB for 1.3 MB of input.

## Excluded

### Remote API clients or MCP wrappers

All these entries send the invoice to a service and cannot be measured offline.

| Package | Downloads/week | Reason |
|---|---|---|
| @invoicexml/zugferd-xrechnung 1.4.0 | 48 | HTTP client for `https://api.invoicexml.com` (`dist/index.js:51`). |
| @facturxapi/sdk | 412 | API client. |
| @invowerk-dev/sdk | 67 | API client. |
| zugferd-client | 13 | Client for the "E-Rechnung API". |
| e-invoice-api | 2,480 | API client. |
| xrechnung-check | 15 | API client. |
| @normapi/client 0.11.1 | 616 | Validates "via the NormAPI service". |
| @conformo/cli, @conformo/mcp-server | 419, 560 | Wrappers around @conformo/toolkit, which is measured. |

### Installed and read, then excluded

| Package | Downloads/week | Reason |
|---|---|---|
| zugferd-ts 0.1.0 (MIT) | 776 | Generator only, from its own input structure: CII EXTENDED, embedding of the XML it generates, TS rules on the input. No parsing, extraction or validation of an XML. It merely duplicates stackforge and node-zugferd. |
| factur-x 0.0.2 | 242 | Unbuilt package (`main: src/index.ts`), MINIMUM profile only ("unknown profile" otherwise), "under development". |
| @e-rechnung-inbox/xrechnung-core 0.1.1 | 33 | In-house structural checks. Extraction by a "simple substring scanner" (README), with no PDF reader. |
| @normwerk/einvoice-cii / -pdfa 0.1.0 | 21 / 93 | Serializer and PDF/A assembly from their own model. Very little used, and they duplicate the measured generators. |
| zugferd-generator 2.1.0 | 613 | Generator from its own input structure (pdf-lib). Duplicate of the measured generators. |
| pdf-to-einvoice 0.3.0 | 746 | Application (CLI, MCP, browser) built on @e-invoice-eu/core, which is measured. |

### Not installed: little used, beyond the cap of about ten libraries

These packages were not read, apart from their npm page. They remain to be looked at if the bench grows.

| Package | Downloads/week | Reason and remark |
|---|---|---|
| @jasy/e-invoice 1.0.0-beta.4 | 720 | CII/UBL + PDF/A-3 generator. |
| @crewhaus/tool-einvoice | 396 | Fixed-format financial files, not a Factur-X library. |
| e-rechnung-parser 0.1.0 | 283 | Reader only (fast-xml-parser). |
| facturx-sdk 1.10.0 | 177 | Advertises built-in BR-FR rules: the first one to check when looking for BR-FR outside a schematron. |
| xrechnung-generator 2.1.0 | 153 | XRechnung UBL generator only. |
| @aifind/mustangjs 0.4.0 | 142 | Port of Mustang with SaxonJS. |
| @itegro/einvoice 0.1.0 | 119 | Generation + SaxonJS schematron + @react-pdf. |
| @financica/facturx 0.2.0 | 89 | Parsing/writing/embedding. |
| @audoora-ext/einvoice 1.2.2 | 20 | GPL-3.0, reader. |
| einvoice-kit, factur-x-ts, @canopynexus/factur-x, @openbim/xrechnung… | < 20 | Practically unused. |

No package requiring Java was measured. The only Java dependency encountered is `xsd-schema-validator`, optional in node-zugferd and bypassed with `strict: false`.

## Coverage table

Legend for the "rules" columns:
- "rewritten": rules written by hand in the library's language, rather than the official schematron;
- "SEF/SaxonJS": official schematron compiled, run by SaxonJS.

| Library | Version | License | Latest release | Downloads | CII parsing | CII writing | UBL | PDF extraction | PDF embedding | XSD | EN 16931 (CII) | Factur-X profiles (MINIMUM, BASIC WL, EXTENDED) | BR-FR | BR-DE (XRechnung) | Peppol BIS 3 | PDF/A-3 container check | External dependencies |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| @stackforge-eu/factur-x | 1.4.7 | EUPL-1.2 | 1.4.7 (2026-11-19) | 94,466 | no | yes (own input structure) | no | yes | yes (PDF/A-3b XMP) | yes, Factur-X 1.09, 5 profiles (libxml2-wasm) | no | XSD only | no | no | no | no | libxml2-wasm (WASM) |
| @e-invoice-eu/core | 3.4.4 | WTFPL | 3.4.4 (2026-11-16) | 63,013 | no | yes (UBL-style JSON) | writing | no | partial (generated XML only) | no | no | generation only | no | generation only | no | no | LibreOffice if no PDF is supplied |
| node-zugferd | 0.1.1-beta.1 | MIT | 0.1.1-beta.1 (`latest`); 1.0.0-beta.1 (`beta`, 2026-11-16) | 62,514 | no | yes (per-profile input) | no | no | partial (generated XML only) | Factur-X 1.07.3, via Java | no | generation for 5 profiles + XSD (Java) | no | no | no | no | Java for validation (optional) |
| @stafyniaksacha/facturx | 0.6.1 | MIT | 0.6.1 (2026-08-28) | 7,339 | yes | yes | no | yes (except `xrechnung.xml`) | fails if the host has an /ID | yes, Factur-X 1.09, 5 profiles | official Factur-X 1.09 schematron (CEN rules), SEF/SaxonJS, no severities | yes (XSD + schematron per profile) | no | no | no (R008 only) | no | SaxonJS (Saxonica license), libxml2-wasm |
| factur-x-kit | 0.3.2 | MIT | 0.3.2 (2026-07-31) | 5,067 | partial (13/26) | yes | no | partial (inside `fromPDF`) | partial (rewrites the XML) | no | partial, rewritten (131 BR ids cited, zod) | MINIMUM, BASIC WL; EXTENDED falls back to EN16931 | no | no | no | no | none |
| @lucaapp/zugferd-ts | 0.1.0 | Apache-2.0 | 0.1.0 (2026-06-13) | 4,395 | partial (19/26) | yes | yes (parsing, writing, XSD) | yes | yes (PDF/A-3b) | yes, Factur-X 1.08 + UBL 2.1 (xmllint-wasm) | yes, CEN 1.3.16, 806 asserts (node-schematron, severities ignored) | XSD per profile | no | no (writes XRechnung without KoSIT rules) | no | partial (21 rules, container only) | xmllint-wasm (WASM) |
| @attestwire/en16931 | 0.14.1 | MIT | 0.14.1 (2026-10-04) | 2,538 | yes | yes | yes (parsing) | yes (in-house reader) | no | no | rewritten (≈ 217 EN 16931 ids; README: 294 rules in total) | partial (AW-PROFILE-*) | no (SIREN/SIRET only) | rewritten (24 BR-DE) | UBL only, rewritten (27 R) | partial (6 AW-PDF-* families, warnings) | none |
| @ekwo-ai/factur-x | 0.27.0 | MIT | 0.27.0 (2026-11-19) | 521 | yes (`readCii`) | yes (own model) | no | yes | yes (XMP) | no | partial (8 BR-CO) | generation for 5 profiles | no | no | no | no | none |
| @conformo/toolkit (+ validate, parse, pdf) | 0.2.4 | Apache-2.0 | 0.2.4 (2026-11-12) | 497 (validate 528) | yes | yes | yes | yes | yes (PDF/A-3, ICC) | yes, CII D16B, Factur-X MINIMUM/BASIC WL/EXTENDED, UBL 2.1 (xmllint-wasm) | yes, CEN 1.3.17, SEF/SaxonJS | XSD only (no schematron, by design) | yes, FNFE-MPE factur-x 6.8, Flux 2 (if `country: "FR"`) | yes, KoSIT 3.0.2 (2026-08-31) | UBL only (3.0.22) | partial (XMP / file name notices) | SaxonJS (Saxonica license), xmllint-wasm, pdfjs-dist |
| @kontor-mcp/core | 1.0.4 | Apache-2.0 | 1.0.4 (2026-08-26) | 337 (rules 554) | yes | partial (`convertInvoice`) | yes | yes | no (draws its own PDF) | yes, CII D16B + UBL 2.1 (xmllint-wasm) | yes, CEN 1.3.16, SEF/SaxonJS | no | no | yes, XRechnung 3.0.2 (schematron 2.5.0) | no | no | SaxonJS (Saxonica license), xmllint-wasm |
| @sweberdev/summand | 2.0.1 | MIT | 2.0.1 (2026-10-13) | 323 | no (Pro version) | no (Pro version) | yes (validation) | yes | no | yes, CII + UBL (in-house validator) | yes, CEN 1.3.16, 806 asserts (in-house engine) | partial (XMP profile compared with BT-24) | no | yes, XRechnung 2.6.0 (52 asserts) | no | partial (XMP profile) | none |
| en16931 | 0.1.0 | MIT AND EUPL-1.2 | 0.1.0 (2026-07-28) | 360 | no | no | yes (validation) | no | no | no | yes, CEN 1.3.16 copied (fontoxpath, exact decimals, no severities) | no | no | no | no | no | none |
| @fin.cx/einvoice | 11.3.1 | MIT | 11.3.1 (2026-09-30) | 323 | yes | yes | yes | yes | partial (attachment without XMP) | no | partial, rewritten (≈ 200 ids) | no | no | partial, rewritten (12 BR-DE) | no (2 ids) | no | network postinstall (not run) |
