# PHP bench: selected and excluded libraries, and coverage of the standard

PHP bench of the common protocol (`../PROTOCOL.md`), compared with the Go
package `github.com/ritoon/zerechnung`. Surveys as of 2026-11-22.

- **Interpreter**: PHP 8.6.0 NTS (Homebrew, `/opt/homebrew/bin/php`),
  libxml2 2.9.13, libxslt 1.1.35, Composer 2.11.2. Extensions checked by
  `./run.sh setup`: dom, libxml, xml, simplexml, xmlreader, xmlwriter,
  mbstring, gd, zlib, bcmath, fileinfo, zip, iconv, ctype (required by the
  pinned packages), pcntl (call timeouts in the harness).
- **OPcache and JIT**: the CLI as installed, OPcache off
  (`opcache.enable_cli=0`, the CLI default), hence no JIT. This is what a PHP
  script launched from the command line gets, and the choice is written in
  the `runtime` field of each JSON line (`PHP 8.6.0 (NTS,
  OPcache off, JIT off, zend.assertions=-1)`). In a CLI process, OPcache
  does not avoid compiling the files (the cache dies with the process, and
  compilation lands in `first_call_ns`); it brings the opcode optimizer and
  makes the JIT possible. For a run with JIT (checked: `runtime` becomes
  `OPcache on, JIT tracing`):
  `PHP_OPTS="-d memory_limit=4G -d zend.assertions=-1 -d display_errors=stderr -d error_reporting=6143 -d opcache.enable_cli=1 -d opcache.jit=tracing -d opcache.jit_buffer_size=64M" ./run.sh bench`.
- **Settings forced** by `run.sh` in `PHP_OPTS`, passed on to every child
  process: `memory_limit=4G` (128M by default; set so that no library fails
  for lack of memory on the large files or during the survey),
  `zend.assertions=-1` (the installed php.ini is the development one,
  `zend.assertions=1`, which would execute the libraries' `assert()` calls),
  `display_errors=stderr` and `error_reporting=6143` (E_ALL without
  deprecations: several libraries trigger them under PHP 8.6, and printing
  them would slow down the measurement). None of these settings changes the
  code that runs.
- **Two Composer environments**, with exactly pinned versions and committed
  `composer.lock` files: `composer.json` (96 packages installed) and
  `invoicesuite/composer.json` (36 packages) for horstoeko/invoicesuite,
  which declares a conflict with horstoeko/zugferd. Installed with
  `composer install --no-scripts --no-plugins --ignore-platform-req=php+
  --optimize-autoloader`: no third-party script or plugin runs at install
  time; `php+` only ignores upper bounds (php-erechnung-toolkit,
  php-pdf-toolkit, the dschuppelius/* packages and
  mahdiabderraouf/facturx-php declare `php <8.6`). `vendor/` weighs 516 MB
  (including 120 MB of vampire-php's precompiled rules),
  `invoicesuite/vendor/` 134 MB.
- **No Java**: no adapter launches a JVM. Validators that need one are not
  measured: KoSIT and veraPDF in horstoeko/zugferd, KoSIT in
  horstoeko/invoicesuite, quittio and php-erechnung-toolkit, the Schematron
  layer of pn-invoice (see the library notes).
- **No network** during `survey` and `bench`. Network code found and not
  exercised: download of the KoSIT validator and of veraPDF
  (horstoeko/zugferd, horstoeko/invoicesuite), sending the invoice to the
  Docuflair service (horstoeko/invoicesuite), download of the Peppol BIS
  rules (intragile), SMP lookup and Peppol transport
  (php-erechnung-toolkit).
- **External tool**: poppler's `pdfdetach` (`/opt/homebrew/bin`), launched by
  mahdiabderraouf/facturx-php and by php-pdf-toolkit (extraction for
  php-erechnung-toolkit). Process launch time is part of the measurement.
- **XSLT 2.0 impossible**: ext-xsl (libxslt 1.1.35) rejects the official
  XSLTs. `XSLTProcessor::importStylesheet` returns `false` on CEN 1.3.17's
  `EN16931-CII-validation.xslt` as well as on
  `horstoeko/zugferd/src/schema/xslt/FACTUR-X_EN16931.xslt`, with
  "xsl:version: only 1.1 features are supported" followed by "Invalid
  expression". The only XSLT 2.0 engine for PHP is the SaxonC-HE extension:
  Saxonica's 13.1.1 macOS arm64 archive (no brew) does not compile against
  PHP 8.6, with 20 errors in `php8_SchemaValidator.cpp` (`XtOffsetOf`
  undefined, `WRONG_PARAM_COUNT` removed from the Zend API). The extension
  was not patched: no PHP bench runs the official XSLTs. The full
  validators measured here replace them with PHP rewrites (vampire-php's
  XPath 2.0 interpreter, intragile's and dealerweb's rules compiled to
  PHP).

Running: `./run.sh setup`, then `./run.sh survey [lib…]` and
`./run.sh bench [lib…]` (variables `OUT`, `COUNT`, `BENCHTIME`; also `PHP`,
`COMPOSER`, `PHP_OPTS`, `FX_TMP` = directory for the temporary files of
path-based APIs, `<system temp>/fx-bench-php` by default, `CHECK_TIMEOUT` =
30 s for a benchmark's check call, `SURVEY_TIMEOUT` = 10 s per survey file,
`CHILD_TIMEOUT` = 1800 s per child process). Each benchmark runs in a
**fresh child process** (`proc_open`), so `first_call_ns` is a true first
call (Composer autoloader loaded, library classes not yet loaded, cold
caches). The survey launches one process per library and per operation; a
child that dies (fatal error, memory) is restarted after the file it died
on. Measurement follows the protocol: N = 1, 2, 5, 10, 20, 50… until one
series lasts at least `BENCHTIME`, then `COUNT` series, `hrtime(true)`, a
single thread.

The `package/file.php:line` paths below are relative to `vendor/`,
except those of horstoeko/invoicesuite and prinsfrank/pdfparser, which are
relative to `invoicesuite/vendor/`.

## Survey results (`./run.sh survey`)

| Library | validate (flagged / 26) | parse (/ 26) | extract (/ 104) |
|---|---|---|---|
| horstoeko/zugferd 1.0.135 | 0 (XSD only) | 26 | 104 |
| horstoeko/invoicesuite 0.0.44 | 0 (XSD only) | 26 | 88 |
| atgp/factur-x 3.5.0 | 4 (XRECHNUNG rejected, XSD only) | — | 100 (90 with `validateXsd: true`, the default) |
| tiime/factur-x 0.3.2 | 4 (XRECHNUNG rejected, XSD only) | — | 70 |
| tiime/cross-industry-invoice 0.3.11.1 | — | 23 | — |
| easybill/zugferd-php 6.0.0 | 0 (XSD only) | 26 | — |
| easybill/e-invoicing 0.4.0 | — | 24 | — |
| mahdiabderraouf/facturx-php 2.4.0 | 4 (XRECHNUNG rejected, XSD only) | — | 100 |
| php-erechnung-toolkit 0.23.1 | 1 (TK-ATT-EMPTY) | 26 | 104 |
| john-wink/en16931-php 0.4.1 | 0 | 26 | — |
| invoicepress/vampire-php 1.3.0 | **6** (BR-CL-25 ×2, CII-SR-465 ×4) | — | — |
| intragile/php-einvoicing 1.0.2 | **6** (BR-CL-25 ×2, CII-SR-465 ×4) | 26 | 104 |
| dealerweb/einvoice 1.2.3 | **6** (BR-CL-25 ×2, CII-SR-465 ×4) | 26 | 104 |
| quittio/core 1.1.0 | 6, not the same ones (QUITTIO-IBAN ×3, BR-CL-25 ×2, DE-USTG-03) | 26 | 104 |
| pnscripts/pn-invoice 0.2.1 | 0 | — | — |
| pat-o-dev/factur-x 0.8.0 | 2 (has-optional-party-name) | — | 1 |
| stann/factur-x 2.0.0 | 0 (XSD only) | — | — |
| smalot/pdfparser 2.13.1 (floor) | — | — | 104 |
| ext-dom (XSD floor) | 0 (XSD only) | — | — |

"—": the operation does not exist in the library (nothing is written to
`survey_php.jsonl`). An exception during validation counts as a flag, as in
the other benches.

**Reference 6/26 reproduced** by three libraries, on the same files and
with the same rules: BR-CL-25 on `EN16931_ElektronischeAdresse` and
`EN16931_RechnungsUebertragung`, CII-SR-465 on the four `XRECHNUNG_*`.
All three apply the KoSIT validator's scenario table, with the custom
levels that make CII-SR-465 an error for XRechnung: vampire-php (CEN 1.3.17
schematron), intragile and dealerweb (CEN 1.3.16, the version pinned by the
KoSIT configuration; intragile ships the `scenarios-2026-08-31` table).

Deviations from the reference in detail:

- **quittio 6/26, but not the same ones**: both BR-CL-25 are seen; the four
  CII-SR-465 are not (no offline schematron, level L5 is the KoSIT
  validator in Java). The four other flags come from quittio's own checks,
  outside EN 16931: QUITTIO-IBAN (IBAN check digits, L7 plausibility level,
  `quittio/core/src/Validation/Rules/PlausibilityRules.php:33`) on
  `EN16931_AbweichenderZahlungsempf`, `EN16931_Innergemeinschaftliche_Lieferungen`
  and `EN16931_OEPNV` (sample IBAN in the corpus), and DE-USTG-03 (delivery
  date required by § 14 UStG,
  `quittio/country-de/src/Rules/UstgRules.php:82`) on
  `EN16931_Physiotherapeut`.
- **john-wink 0/26**: BR-CL-25 is coded at warning level
  (`john-wink/en16931-php/src/RuleSets.php:227-228`, `Severity::Warning`),
  and CII-SR-465 is missing from `resources/syntax-rules.json` (the
  library's reference is validator-configuration-xrechnung 3.0.2 of
  2025-03-21, which predates the custom level). Only fatal violations are
  counted (the library's verdict).
- **atgp, tiime/factur-x, mahdiabderraouf 4/26**: the four `XRECHNUNG_*` are
  rejected (exception "Invalid Factur-X URN" or "invalid profile found"),
  since the XRechnung identifier is not a Factur-X profile for these
  libraries. The other 22 pass the XSD.
- **php-erechnung-toolkit 1/26**: TK-ATT-EMPTY (empty attachment, a toolkit
  check) on `EN16931_Betriebskostenabrechnung`.
- **pat-o-dev 2/26**: has-optional-party-name (a structural rule of the
  library) on both `*_Betriebskostenabrechnung`.
- **horstoeko/zugferd, invoicesuite, easybill/zugferd-php, stann, ext-dom
  0/26**: XSD only, and all 26 files conform to the XSD. pn-invoice 0/26:
  its 104 business rules cover neither BR-CL nor CII-SR.
- **parse**: tiime/cross-industry-invoice rejects the `EM` scheme ("Wrong
  schemeID", 2 files) and a `SpecifiedProcuringProject` without a name
  (`EN16931_Physiotherapeut`). easybill/e-invoicing rejects EAS code 9958,
  missing from its enumeration (the same two files as BR-CL-25).
- **extract**: atgp and mahdiabderraouf do not look for the `xrechnung.xml`
  attachment (the four `FX/XRECHNUNG_*.pdf`); tiime/factur-x only looks for
  `factur-x.xml` (`tiime/factur-x/src/Reader.php:44`, 34 failures); atgp
  with its default options also validates the XSD and rejects 10 more XML
  files. invoicesuite (prinsfrank/pdfparser) fails on 16 PDFs: unreadable
  `gzuncompress` streams, root not found, and embedded XRechnung whose
  content no format provider recognizes. pat-o-dev reads only one: 62
  unresolved /Names or /AF references, 38 indirect `/Length`, 3 Filespecs
  without /EF (the library advertises classic xref only).

## Library notes: selected libraries

### horstoeko/zugferd 1.0.135

MIT, latest release 2026-11-08, 10,006,092 downloads in total, 1,252,800
last month. The reference PHP library (ZUGFeRD 1 and 2, Factur-X,
XRechnung CII).

- Parsing `ZugferdDocumentReader::readAndGuessFromContent`
  (`horstoeko/zugferd/src/ZugferdDocumentReader.php:306`), JMS model
  generated from the XSDs; writing `ZugferdDocument::serializeAsXml`.
- PDF: `ZugferdDocumentPdfReader::getXmlFromContent`
  (`src/ZugferdDocumentPdfReader.php:97`, smalot/pdfparser); embedding
  `ZugferdDocumentPdfMerger` (`src/ZugferdDocumentPdfMerger.php:56`): FPDI
  re-imports the pages, PDF/A-3B XMP.
- Offline validation: `ZugferdXsdValidator`
  (`src/ZugferdXsdValidator.php:66`, on the re-serialized document, lines
  164-182) with the profile's Factur-X XSD; the XRechnung profile points to
  the EN16931 XSD (`src/ZugferdProfiles.php:260-266`).
  `ZugferdDocumentValidator` (`src/ZugferdDocumentValidator.php:67`) only
  checks Symfony constraints (`xsd_rules` group).
- Schematron and XSLT shipped (`src/schema/schematron`, `src/schema/xslt`),
  declared in `src/ZugferdProfiles.php:109-157` but executed by no code;
  ext-xsl rejects the (2.0) XSLT anyway.
- `ZugferdKositValidator` downloads the KoSIT validator 1.5.0 and the
  XRechnung 3.0.2 configuration (`src/ZugferdKositValidator.php:96`,
  `:103`) and launches `java` (`:786`, `:999`); `ZugferdPdfValidator`
  downloads veraPDF (`src/ZugferdPdfValidator.php:158`) and launches Java
  (`:636`). Not measured.

### horstoeko/invoicesuite 0.0.44

MIT, latest release 2026-11-20, 21,563 in total, 15,658 last month.
Successor to horstoeko/zugferd by the same author: one DTO for many formats
(ZUGFeRD/Factur-X, XRechnung CII and UBL, Peppol, FatturaPA…). Separate
Composer environment (declared conflict with horstoeko/zugferd).

- Parsing `InvoiceSuiteDocumentReader::createFromContent` (guessed format
  provider, JMS); writing `copyToBuilder()->getContent()`.
- PDF: `InvoiceSuitePdfDocumentReader::createFromContent`
  (`horstoeko/invoicesuite/src/InvoiceSuitePdfDocumentReader.php:108`)
  reads the PDF with prinsfrank/pdfparser 4.4.0 and **also deserializes the
  XML it finds** (included in `Extract_*`); badly compressed streams fail in
  `prinsfrank/pdfparser/src/Document/Filter/Decode/FlateDecode.php:136`.
  Embedding `InvoiceSuitePdfDocumentBuilder` (FPDI/TCPDF).
- Offline validation: `InvoiceSuiteXsdDocumentValidator`
  (`src/validators/InvoiceSuiteXsdDocumentValidator.php:200` XSD selection,
  `:236` and `:248` `schemaValidate`). The two other validators are not
  measured: `InvoiceSuiteKositDocumentValidator` downloads the KoSIT 1.6.2
  jar (`src/validators/InvoiceSuiteKositDocumentValidator.php:54`) and
  launches Java; `InvoiceSuiteDocuflairDocumentValidator` sends the invoice
  to `https://invoice.docuflair.com` (`:39`, curl at lines 132 and 181).

### atgp/factur-x 3.5.0

MIT, latest release 2026-08-18, 1,122,993 in total, 61,110 last month.
Model-less Factur-X tool: extraction, embedding, XSD.

- `Reader::extractXML` (`atgp/factur-x/src/Reader.php:40`,
  smalot/pdfparser), called with `validateXsd: false` (by default extraction
  also validates the XSD); `Writer::generate` (`src/Writer.php:86`) with
  FPDI, PDF/A-3 XMP, XSD included by default.
- `XsdValidator` (`src/XsdValidator.php:21-28`): Factur-X **1.09.2** XSDs
  for the five profiles, nothing else; XRechnung is not a recognized
  profile.

### tiime/factur-x 0.3.2

MIT, latest release 2026-10-07, 280,508 in total, 4,784 last month.
Rewrite of atgp/factur-x for PHP 8.3.

- `Reader::extractXML` (`tiime/factur-x/src/Reader.php:20`) only looks for
  `factur-x.xml` (`:44`); `Writer::generate` with FPDI.
- Factur-X **1.07.3** XSD (`src/XsdProcessor.php:7-13`). The schematron
  shipped under `schematron/` is called by no code in `src/`.

### tiime/cross-industry-invoice 0.3.11.1

MIT, latest release 2025-10-27, 269,684 in total, 808 last month.

- CII model typed per profile (Minimum, BasicWL, Basic, EN16931 and a
  French "Flux1" model), `CrossIndustryInvoice::fromXML` and `toXML`
  (`tiime/cross-industry-invoice/src/EN16931/CrossIndustryInvoice.php:77`,
  `:42`). The EN 16931 semantic model lives in tiime/en-16931 0.10.5
  (dependency).
- The XSDs and schematrons under `xsd/` only serve the package's tests: no
  validation API. No PDF.

### easybill/zugferd-php 6.0.0

MIT, latest release 2026-03-30, 407,329 in total, 22,919 last month.

- JMS model for ZUGFeRD 1 and ZUGFeRD 2 / Factur-X: `Reader::transform`
  (`easybill/zugferd-php/src/zugferd2/Reader.php:15`), `Builder::transform`
  (`src/zugferd2/Builder.php:22`).
- `Validator::validateAgainstXsd` (`src/zugferd2/Validator.php:11-15`):
  Factur-X **1.08.0** XSD, schema passed in by the caller (the adapter picks
  it from BT-24, EN16931 for XRechnung). No PDF.

### easybill/e-invoicing 0.4.0

MIT, latest release 2026-09-07, 198,338 in total, 9,426 last month.

- JMS models for CII and UBL (EN 16931 and its CIUSs): `Reader::read`
  (`easybill/e-invoicing/src/Reader.php:16`), `Transformer::transformToXml`
  (`src/Transformer.php:16`). No validation, no PDF.

### mahdiabderraouf/facturx-php 2.4.0

MIT, latest release 2026-09-11, 12,121 in total, 1,708 last month.

- Extraction through poppler's `pdfdetach`, launched twice via `exec`
  (`mahdiabderraouf/facturx-php/src/Helpers/XmlExtractor.php:40` `-list`,
  `:74` `-save` to a temporary file): path-based API, process launch time
  is measured.
- XML generator for MINIMUM, BASIC WL and BASIC only (`readme.md:24`), with
  no reader: `Write_Small` starts from a hand-built model based on
  `EN16931_Einfach`, written at the BASIC profile.
- `Validator::validate` (`src/Validator.php:35`): Factur-X **1.09.2** XSD
  for the profile read from BT-24; the `Profile` enumeration has no
  XRechnung case, hence the rejection of the `XRECHNUNG_*` files.
- `Generator::generate` validates the XML against the XSD before embedding
  (`src/Generator.php:88`), then PDF/A-3b via FPDI/tFPDF (`:94`).

### daniel-jorg-schuppelius/php-erechnung-toolkit 0.23.1 (+ php-pdf-toolkit 0.33.0)

AGPL-3.0-or-later (both packages), latest release 2026-11-14,
7,295 in total, 1,468 last month (php-pdf-toolkit: 7,671 and 1,555).
German toolkit: XRechnung/ZUGFeRD, UBL and CII, Order-X, GAEB, Peppol.

- Parsing `ERechnungParser::parse`, writing
  `ERechnungGenerator::generateCii`.
- Extraction: `ZugferdPdfParser::extractXml`
  (`daniel-jorg-schuppelius/php-erechnung-toolkit/src/Parsers/ZugferdPdfParser.php:98`)
  delegates to php-pdf-toolkit, whose `ZugferdReader`
  (`daniel-jorg-schuppelius/php-pdf-toolkit/src/Readers/ZugferdReader.php:73`,
  `:103`) launches `pdfdetach` (or pdftk). Path-based API.
- Offline validation: `CiiSchemaValidator`
  (`src/Validators/CiiSchemaValidator.php:63`, CII D16B XSD) then
  `DocumentRulesValidator` (`src/Validators/DocumentRulesValidator.php:36`):
  BR-29, BR-30, BR-52, BR-CL-24, BR-CO-19, BR-CO-20, BR-DE-22 and the
  toolkit's own TK-* checks. The EN 16931 / XRechnung schematron goes
  through `KositValidator` (`src/Validators/KositValidator.php:95`,
  `Java::execute` at line 160, command in
  `src/Helper/KositCommandHelper.php:22-56`): the README says so
  (`README.md:436`, `:487`: XSLT 2.0, hence Java). Not measured.
- No embedding into a host PDF (the toolkit draws its own PDF). The library
  writes configuration errors (ocrmypdf missing) to stderr at load time: no
  effect on the results.

### john-wink/en16931-php 0.4.1

MIT, latest release 2026-10-16, 1,292 in total, 579 last month. Native PHP
port of the KoSIT validation chain, without Java or XSLT.

- `En16931Validator::xrechnung()` or `::en16931()`
  (`john-wink/en16931-php/src/En16931Validator.php:51`, `:59`), `validate`
  (`:105`): CII D16B XSD (`src/Syntax/SchemaValidator.php:65`), then 1276
  syntax rules (CII-SR/DT, UBL-CR/SR) read from
  `resources/syntax-rules.json`, then 274 EN 16931 business rules and
  XRechnung CIUS rules coded in PHP (the package's `COVERAGE.md:15-16`).
  Stated reference: validator-configuration-xrechnung 3.0.2.
- Observed deviations: BR-CL-25 as a warning, CII-SR-465 missing (see the
  survey).
- Parsing `CiiInvoiceReader::read` (and `UblInvoiceReader`); no writing, no
  PDF.

### invoicepress/vampire-php 1.3.0

GPL-2.0-or-later, latest release 2026-10-26, 377 in total, 373 last month.

- The official schematrons precompiled to JSON (`rulesets/`: EN16931-CII
  and UBL **1.3.17**, XRechnung CII and UBL **3.0.2**, Factur-X **1.09.2**
  for the five profiles, Peppol EN16931 UBL **3.0.21**, CIUS-RO, CIUS-HR…,
  120 MB), executed by an XPath 2.0 interpreter written in PHP
  (`README.md:7`), after the XSD.
- The scenario is chosen the way the KoSIT validator does it
  (`rulesets/scenarios.json`, 22 scenarios; XRechnung CII with the custom
  levels, including CII-SR-465 as an error). `Validator::validateString`
  (`src/Validator.php:123`); errors = messages that block acceptance
  (`src/Result/Result.php:248`).
- Validation only. The slowest of the full validators (66 to 107 ms per
  bench invoice).

### intragile/php-einvoicing 1.0.2

MIT (rules: rules-en16931 EUPL-1.2, rules-xrechnung and rules-facturx
Apache-2.0), latest release 2026-11-17, 69 in total, 69 last month
(package published in November 2026).

- EN 16931 model, CII and UBL readers and writers: `CiiReader::read`
  (`intragile/php-einvoicing/src/Syntax/Cii/CiiReader.php:58`),
  `CiiWriter::write` (`src/Syntax/Cii/CiiWriter.php:74`).
- PDF written without FPDI: `FacturXPdfReader::read`
  (`src/Syntax/FacturX/FacturXPdfReader.php:60`); `FacturXPdfWriter::embed`
  (`src/Syntax/FacturX/FacturXPdfWriter.php:68`) by incremental update,
  without PDF/A conversion, and it refuses a PDF that already declares an
  invoice (`refuseExistingInvoice`, line 172).
- Validation: official schematrons precompiled to PHP
  (`intragile/einvoicing-rules-*` packages: CEN 1.3.16 and 1.3.17,
  XRechnung 2.6.0, Factur-X 1.09.2 for the five profiles), in-house
  XPath 2.0 evaluator. `XRechnungValidator` reproduces the KoSIT scenarios
  (`src/Validation/XRechnung/XRechnungValidator.php:110-130`, CEN `1.3.16`
  at line 117, `scenarios-2026-08-31` table); `FacturXValidator`
  (`src/Validation/FacturX/FacturXValidator.php:72`, `:102`) for the other
  Factur-X profiles.
- The Peppol BIS rules are downloaded from
  `raw.githubusercontent.com/OpenPEPPOL/peppol-bis-invoice-3/3.0.22`
  (`src/Validation/Peppol/PeppolSources.php:18`): not installed, not
  measured.

### dealerweb/einvoice 1.2.3

MIT, latest release 2026-10-30, 200 in total, 148 last month.

- EN 16931 model read by the KoSIT syntax → model transformation, compiled
  to PHP: `Invoice::fromXml`; `Invoice::xmlFromPdf`
  (`dealerweb/einvoice/src/Model/Behavior/InvoiceBehavior.php:78`); writing
  of all profiles, CII and UBL.
- Validation: official XSLTs compiled to PHP from Saxon's execution plan
  (`README.md:426`: CEN 1.3.16, XRechnung 3.0.2, FeRD 1.09.2), with the
  KoSIT validator's scenarios and levels (`README.md:358`).
  `Validator` has `frenchRules = false` by default
  (`src/Validation/Validator.php:49`): the BR-FR rules (FNFE-MPE France_RFE
  1.4.0.04, `resources/france/SOURCE.md`) do not run. Peppol BIS 3.0 rules
  rewritten for UBL and CII (`resources/peppol/SOURCE.md`).
- `toXml()` and `toPdf()` (`InvoiceBehavior.php:137`, `:155`) revalidate
  their output: the adapter goes through `Generator` with `validate: false`
  (`src/Generation/Generator.php:89`) so that `Write_*` does not measure a
  validation.
- Its own PDF layer (classic xref tables and xref streams), the document is
  rewritten as a single revision; no PDF/A check (`README.md:524`).

### quittio/core 1.1.0 (+ quittio/country-de 1.1.0)

MIT, latest release 2026-11-17, 361 in total, 319 last month
(country-de: 320 and 301). Complete invoicing framework: model,
calculation, CII and UBL readers and writers, ZUGFeRD/Factur-X rendering
with mPDF.

- `Quittio::create` assembles the validators (`quittio/core/src/Quittio.php:183-240`):
  `PhpRulesDriver` (`:210`, levels L1 to L3 and L7), `XsdValidator` (`:215`,
  CII D16B XSD, `src/Validation/Drivers/XsdValidator.php:113`),
  `KositValidator` (`:219`, only if configured: `java -jar`,
  `src/Validation/Drivers/KositValidator.php:145`) and `VeraPdfValidator`
  (`:231`, Docker or binary). Levels described in
  `src/Validation/ValidationLevel.php:12-25`.
- Rules in PHP: 68 BR identifiers in the core, 23 BR-DE in country-de
  (`quittio/country-de/src/Rules/XRechnungRules.php`), plus plausibility
  checks (IBAN, VAT identifiers, Leitweg-ID) and German-law checks
  (§ 14 UStG).
- Extraction `PdfEmbeddedFileExtractor::extractFile`
  (`src/Pdf/PdfEmbeddedFileExtractor.php:72`); embedding
  `MpdfPostProcessor::process` (`src/Pdf/PostProcessing/MpdfPostProcessor.php:34`),
  a rendering step that re-imports the pages via FPDI into mPDF.

### pnscripts/pn-invoice 0.2.1

MIT, latest release 2026-10-16, 187 in total, 184 last month.

- EN 16931 model and CII and UBL writers; no XML reader
  (`README.md:168`): `Write_Small` starts from a hand-built model based on
  `EN16931_Einfach`, `Write_Lines` is `unsupported`.
- `InvoiceValidator` (`src/Validation/InvoiceValidator.php:28`, `:36`):
  XSD (UBL 2.1, CII D16B) then 104 business rules written in PHP
  (`README.md:114`; no BR-CL, no BR-DEC, no syntax rules). The schematron
  layer is an external command (Saxon on Java, `README.md:138`): not wired
  in.
- No PDF (announced for a commercial package).

### pat-o-dev/factur-x 0.8.0

MIT, latest release 2026-07-31, 84 in total, 1 last month.
France-oriented Factur-X EN 16931 generator.

- Writing `FacturXGenerator::generateXml` from its own model, with no
  reader: `Write_Small` starts from a hand-built model (unit C62, since H87
  does not exist in its enumeration; note without a REG subject).
- `InvoiceValidator`: a selection of "hand-picked" rules (structure,
  BR-CO-15, BR-CO-25, BR-CO-5-6) and French rules BR-FR-MAP-12/14
  (`src/Validation/InvoiceValidator.php:118`, `:123`), not counted here.
  Neither XSD nor schematron.
- PDF: classic xref only (`README.md:31`); `PdfA3Attacher::attach`
  (`src/Pdf/PdfA3Attacher.php:35`) by incremental update, which refuses a
  catalog that already carries /AF or /Names;
  `EmbeddedXmlExtractor::extract` (`src/Pdf/EmbeddedXmlExtractor.php:27`).

### stann/factur-x 2.0.0

No license declared (empty field on Packagist), latest release 2025-12-23,
6,078 in total, 564 last month. Factur-X generator for Chorus Pro.

- Header-only model (`stann/factur-x/src/CrossIndustryInvoice.php:13-22`:
  no lines, no VAT breakdown): `Write_*` not measured.
- `FacturX::validate` (`src/FacturX.php:57`): Factur-X **1.0.06** XSD
  (`src/Schema.php:9-13`). The library leaves the libxml error buffer
  filled when it throws an exception: the adapter clears it before each
  call, otherwise the errors from one invoice would count against the
  following ones.
- `FacturX::createFile` (`src/FacturX.php:20`): mPDF in PDF/A-3B (`:26`),
  pages re-imported via FPDI; mPDF writes its cache to
  `vendor/mpdf/mpdf/tmp` (no `tempDir` option).

### smalot/pdfparser 2.13.1 (extraction floor)

LGPL-3.0, latest release 2026-10-13, 60,612,162 in total, 6,173,943 last
month. General-purpose PDF parser, used by horstoeko/zugferd, atgp and
tiime for extraction. Floor written here (`adapters/smalot-pdfparser.php`):
`Parser::parseContent`, then the first /Filespec named `factur-x.xml`,
`zugferd-invoice.xml` or `xrechnung.xml`, stream /EF /F. 104/104.

### ext-dom (XSD validation floor)

PHP 8.6.0 extension (libxml2 2.9.13). Floor written here
(`adapters/ext-dom.php`): `DOMDocument::loadXML` then `schemaValidate` with
the profile's Factur-X 1.09.2 XSD, with the files taken from
`vendor/atgp/factur-x/xsd` (atgp's PHP code is not called).

## Benchmark status

Trial run `COUNT=1 BENCHTIME=0.05 ./run.sh bench` (19 libraries × 16
benchmarks, 304 lines; the timings are not meant to be used):

| Library | ok | unsupported | fail |
|---|---|---|---|
| horstoeko/zugferd | 13 | CheckPDF ×2 | Embed_Small (FPDI) |
| horstoeko/invoicesuite | 12 | CheckPDF ×2 | Extract_Large (prinsfrank: `gzuncompress`), Embed_Small (FPDI) |
| atgp/factur-x | 7 | Parse ×3, Write ×2, CheckPDF ×2 | ParseValidate_Large (XRechnung profile rejected), Embed_Small (FPDI) |
| tiime/factur-x | 6 | Parse ×3, Write ×2, CheckPDF ×2 | Extract_Extended (`zugferd-invoice.xml` attachment), ParseValidate_Large (XRechnung profile rejected), Embed_Small (FPDI) |
| tiime/cross-industry-invoice | 5 | Extract ×4, ParseValidate ×3, Embed ×2, CheckPDF ×2 | — |
| easybill/zugferd-php | 8 | Extract ×4, Embed ×2, CheckPDF ×2 | — |
| easybill/e-invoicing | 5 | Extract ×4, ParseValidate ×3, Embed ×2, CheckPDF ×2 | — |
| mahdiabderraouf/facturx-php | 8 | Parse ×3, Write_Lines, CheckPDF ×2 | ParseValidate_Large (XRechnung profile rejected), Embed_Small (FPDI) |
| php-erechnung-toolkit | 12 | Embed ×2, CheckPDF ×2 | — |
| john-wink/en16931-php | 6 | Extract ×4, Write ×2, Embed ×2, CheckPDF ×2 | — |
| invoicepress/vampire-php | 3 (ParseValidate) | 13 | — |
| intragile/php-einvoicing | 14 | CheckPDF ×2 | — |
| dealerweb/einvoice | 14 | CheckPDF ×2 | — |
| quittio/core | 13 | CheckPDF ×2 | Embed_Small (FPDI via mPDF) |
| pnscripts/pn-invoice | 4 (ParseValidate ×3, Write_Small) | Extract ×4, Parse ×3, Write_Lines, Embed ×2, CheckPDF ×2 | — |
| pat-o-dev/factur-x | 5 (ParseValidate ×3, Write_Small, EmbedClassic_Small) | Parse ×3, Write_Lines, CheckPDF ×2 | Extract ×4 (unresolved /Names or /AF references), Embed_Small (xref stream) |
| stann/factur-x | 4 (ParseValidate ×3, EmbedClassic_Small) | Extract ×4, Parse ×3, Write ×2, CheckPDF ×2 | Embed_Small (FPDI via mPDF) |
| smalot/pdfparser | 4 (Extract) | 12 | — |
| ext-dom | 3 (ParseValidate) | 13 | — |

No native offline `CheckPDF` in PHP: the only container checks found
launch veraPDF (Java, Docker or an external binary).

Things to know when reading the numbers:

- **FPDI and xref streams**: `host.pdf` (`Embed_Small`) has a
  cross-reference stream (PDF 1.5); the free parser shipped with FPDI
  cannot read it ("compression technique which is not supported by the free
  parser"). All the FPDI- or mPDF-based embedders (horstoeko, invoicesuite,
  atgp, tiime, mahdiabderraouf, quittio, stann) fail on `Embed_Small` and
  succeed on `host_classic.pdf`. dealerweb (its own PDF layer) and
  intragile (no FPDI) succeed on both.
- **Hosts without an invoice**: `host.pdf` and `host_classic.pdf` have no
  attachment, no `/AF` and no Factur-X XMP. Until 2026-11-22, their catalog
  kept an `/AF` pointing to a `factur-x.xml` Filespec (pdfcpu
  `RemoveAttachments` only removed the name tree): intragile, which refuses
  to add a second invoice (`InvoiceAlreadyEmbedded`, a choice made by the
  library), and pat-o-dev, which refuses a catalog that already has `/AF`
  or `/Names`, failed. On the fixed hosts, intragile embeds into both;
  pat-o-dev embeds into `host_classic.pdf` and still fails on `host.pdf`
  (xref streams not supported).
- `Embed_*` is read back with the library itself when it can extract,
  otherwise with horstoeko/zugferd (`note: "read back with horstoeko/zugferd"`,
  the case of stann).
- `Write_*` starts from the model produced by the library's parser (parsing
  outside the measurement). Without a parser, `Write_Small` starts from a
  hand-built model based on `EN16931_Einfach`
  (`note: "hand-built model"`) and `Write_Lines` is
  `unsupported`: pn-invoice, pat-o-dev, and mahdiabderraouf (written at the
  BASIC profile, since the generator lacks EN16931). These three
  `Write_Small` results are not comparable with the others.
- dealerweb `Write_*` and `Embed_*` pass `validate: false` to `Generator`;
  with `toXml()` or `toPdf()`, the output would be revalidated.
- mahdiabderraouf `EmbedClassic_Small` includes the XSD validation that
  `Generator::generate` performs before embedding (native behavior).
- invoicesuite `Extract_*` includes deserializing the extracted XML (the
  PDF reader always does it).
- atgp and tiime `EmbedClassic_Small` keep their default options (XSD
  validation of the XML included).
- mahdiabderraouf and php-erechnung-toolkit `Extract_*` launch `pdfdetach`
  (two processes per call, about 25 to 40 ms): their time is mostly process
  launch time. Their APIs take a path: the harness passes the sample's
  path, and for the `Embed` read-back a temporary file in `FX_TMP`.
- vampire-php `ParseValidate_*`: 66 to 107 ms per invoice (XPath 2.0
  interpreter in PHP), versus 8 to 16 ms for intragile and dealerweb, which
  check the same rules.
- stann: the libxml buffer is cleared by the adapter before each
  validation (see its profile); mPDF writes its cache to `vendor/`.
- quittio and mPDF: mPDF's cache and quittio's temporary files go to
  `FX_TMP`.

## Library notes: excluded libraries

| Package | Reason |
|---|---|
| josemmo/einvoicing 0.3.1 | UBL only (no CII); 575,090 in total, 67,004/month |
| num-num/ubl-invoice 2.4.1, cleverit, invoiceninja/ubl_invoice | UBL only; num-num: 1,188,255 in total, 75,659/month |
| doards/xrechnung-validator 1.0.1 | requires the SaxonC extension for PHP (which does not compile under PHP 8.6), UBL and Peppol only; 60 in total |
| silarhi/zugferd-factur-x 2.0.0, pyrexx/zugferd-php 1.0.6 | old forks of easybill/zugferd-php (2024-02 and 2019-06) |
| ephenyxdigital/einvoice 2.2.7, charkoulis/einvoicing 1.0.0 | forks of josemmo/einvoicing; 44 and 6 downloads in total |
| code4business/zugferd-php, einfacharchiv/zugferd, lianhua/phacturx, inwendo/zugferd | abandoned (flagged as such on Packagist), PHP 7 |
| kinulab/facturx 1.0.13, tualo/zugferd 1.0.44, youniwemi/digital-invoice 0.4.3 | thin wrappers around the measured libraries: kinulab around atgp/factur-x ^1.0, tualo around easybill/zugferd-php, youniwemi around easybill/zugferd-php, atgp/factur-x and josemmo/einvoicing (dev-master branch) |
| horstoeko/zugferd-laravel, intragile/laravel-zugferd, quittio/quittio, tmonier/sylius-facturx-plugin | Laravel or Sylius integrations of the measured libraries |
| horstoeko/zugferdublbridge, horstoeko/zugferdvisualizer, horstoeko/zugferdmail | CII ↔ UBL conversion, HTML/PDF rendering, email reading: outside the protocol |
| horstoeko/ubl 1.0.2 | old, PHP 7, latest release 2021-02 |
| invoicepress/vampire-zugferd, invoicepress/vampire-invoicesuite | bridges between vampire-php and horstoeko (duplicates of what is already measured) |
| vinelabs-de/xrechnung-kit, vineethkrishnan/xrechnung-kit-core 2.1.0 | XRechnung UBL template generator + XSD; 100 and 16 in total |
| d3/xrechnungproject | application (OXID module), GPL-3.0, 22 downloads |
| invoiceninja/einvoice 1.1.0 | schemas under the Elastic-2.0 license, model generated for Invoice Ninja |
| tiime/en-16931 0.10.5 | EN 16931 semantic model without XML, a dependency of tiime/cross-industry-invoice and tiime/factur-x, measured through them |
| scell/sdk, logalysm/valide-facture, factpulse/sdk, xmlvalid, ventnet/xinvoice-client, einvoicing/sdk, invoicexml/*, facturino, peppol-sh, selectpdf, gotenberg, e-invoice/e-invoice-api | clients for remote APIs or containers: validation or generation happens elsewhere |
| genkgo/xsl 1.7.1 | generic XSLT 2.0 to XSLT 1.0 transpiler (not a Factur-X library), declares `php ~8.3 \|\| ~8.4 \|\| ~8.5`, not tried |
| milo/schematron 1.0.1 | generic Schematron engine in XPath 1.0 (latest release 2014-04): cannot read the official schematrons (XPath 2.0) |
| jdomenechb/xslt2processor | partial XSLT 2.0 engine, abandoned |
| SaxonC-HE 13.1.1 extension | the only native XSLT 2.0 engine for PHP, does not compile under PHP 8.6 (see above) |

## Coverage of the standard

| Library | Version | License | Latest version | Downloads | CII parsing | CII writing | UBL | PDF extraction | PDF embedding | XSD | EN 16931 (CII) | Factur-X profiles (MINIMUM, BASIC WL, EXTENDED) | BR-FR | BR-DE (XRechnung) | Peppol BIS 3 | PDF/A-3 container check | External dependencies |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| horstoeko/zugferd | 1.0.135 | MIT | 2026-11-08 | 10,006,092; 1,252,800/month | yes | yes | no (separate zugferdublbridge bridge) | yes (104/104) | yes (FPDI, PDF/A-3B XMP; no xref streams) | yes (Factur-X, 5 profiles; XRechnung on EN16931) | not offline (KoSIT via Java) | XSD for the 5 profiles; schematron shipped, never executed | no | KoSIT via Java only | no | veraPDF via Java only | smalot/pdfparser, setasign/fpdi, jms/serializer, symfony/validator; Java for KoSIT and veraPDF (downloaded) |
| horstoeko/invoicesuite | 0.0.44 | MIT | 2026-11-20 | 21,563; 15,658/month | yes | yes | yes (XRechnung UBL, Peppol, other formats) | partial (88/104) | yes (FPDI/TCPDF; no xref streams) | yes (XSD of the detected format) | not offline (KoSIT via Java, remote Docuflair) | XSD | no | KoSIT via Java only | no | no | prinsfrank/pdfparser, setasign/fpdi, mpdf, jms/serializer; separate Composer environment |
| atgp/factur-x | 3.5.0 | MIT | 2026-08-18 | 1,122,993; 61,110/month | no | no | no | yes (100/104; not `xrechnung.xml`) | yes (FPDI, PDF/A-3 XMP; no xref streams) | yes (Factur-X 1.09.2, 5 profiles) | no | XSD for the 5 profiles (XRechnung rejected) | no | no | no | no | smalot/pdfparser, setasign/fpdi |
| tiime/factur-x | 0.3.2 | MIT | 2026-10-07 | 280,508; 4,784/month | no | no | no | partial (70/104: `factur-x.xml` only) | yes (FPDI; no xref streams) | yes (Factur-X 1.07.3) | no (schematron shipped, never called) | XSD for the 5 profiles (XRechnung rejected) | no | no | no | no | smalot/pdfparser, setasign/fpdi, tiime/en-16931 |
| tiime/cross-industry-invoice | 0.3.11.1 | MIT | 2025-10-27 | 269,684; 808/month | yes (23/26) | yes | no | no | no | no (test XSDs only) | no | MINIMUM, BASIC WL, BASIC, EN 16931 models; no EXTENDED | "Flux1" model, no rules | no | no | no | tiime/en-16931, twig |
| easybill/zugferd-php | 6.0.0 | MIT | 2026-03-30 | 407,329; 22,919/month | yes | yes | no | no | no | yes (Factur-X 1.08.0, schema chosen by the caller) | no | XSD for the 5 profiles; ZUGFeRD 1 | no | no | no | no | jms/serializer |
| easybill/e-invoicing | 0.4.0 | MIT | 2026-09-07 | 198,338; 9,426/month | yes (24/26) | yes | yes (parsing, writing) | no | no | no | no | no | no | no | no | no | jms/serializer |
| mahdiabderraouf/facturx-php | 2.4.0 | MIT | 2026-09-11 | 12,121; 1,708/month | no | partial (MINIMUM, BASIC WL, BASIC) | no | yes (100/104, via `pdfdetach`) | yes (PDF/A-3b via FPDI/tFPDF, XSD checked first; no xref streams) | yes (Factur-X 1.09.2) | no | XSD for the 5 profiles (XRechnung rejected) | no | no | no | no | setasign/fpdi, fpdf; `pdfdetach` (poppler); declares `php <8.6` |
| daniel-jorg-schuppelius/php-erechnung-toolkit | 0.23.1 (+ php-pdf-toolkit 0.33.0) | AGPL-3.0-or-later | 2026-11-14 | 7,295; 1,468/month | yes | yes | yes (parser, generator, UBL XSD) | yes (104/104, via `pdfdetach`) | no (draws its own PDF) | yes (CII D16B, UBL) | partial offline (BR-29, BR-30, BR-52, BR-CL-24, BR-CO-19/20); schematron via KoSIT (Java) | MINIMUM to EXTENDED profiles in the model, no rules of their own | no | BR-DE-22 offline; the rest via KoSIT (Java) | no (SMP client and Peppol transport, network) | no | php-pdf-toolkit (AGPL: dompdf, tcpdf, fpdi), `pdfdetach`; Java for KoSIT; declares `php <8.6` |
| john-wink/en16931-php | 0.4.1 | MIT | 2026-10-16 | 1,292; 579/month | yes | no | yes (parsing, validation) | no | no | yes (CII D16B, UBL 2.1) | yes, rewritten in PHP (274 business rules, 1276 syntax rules); BR-CL-25 as a warning, CII-SR-465 missing | no | no | yes (XRechnung CIUS in PHP) | no | no | none (bcmath) |
| invoicepress/vampire-php | 1.3.0 | GPL-2.0-or-later | 2026-10-26 | 377; 373/month | — | — | yes (CEN, XRechnung, Peppol, CIUS-RO, CIUS-HR… validation) | — | — | yes (the scenario's XSD) | yes (CEN 1.3.17 schematron precompiled, XPath 2.0 in PHP) | Factur-X 1.09.2 schematron for the 5 profiles | no | yes (XRechnung 3.0.2, KoSIT levels) | yes (PEPPOL-EN16931 3.0.21, UBL only) | no | none (bcmath); 120 MB of rules |
| intragile/php-einvoicing | 1.0.2 | MIT; rules EUPL-1.2 and Apache-2.0 | 2026-11-17 | 69; 69/month | yes | yes | yes (parsing, writing, CEN UBL schematron) | yes (104/104) | partial (incremental update without PDF/A; refuses a PDF that already declares an invoice) | yes (CII D16B) | yes (CEN 1.3.16 and 1.3.17 schematrons precompiled to PHP; 1.3.16 via the KoSIT scenarios) | Factur-X 1.09.2 schematron for the 5 profiles | no | yes (XRechnung 2.6.0, KoSIT levels 2026-08-31) | not installed (OpenPEPPOL 3.0.22 rules downloaded from GitHub) | no | brick/math, intragile/einvoicing-rules-* |
| dealerweb/einvoice | 1.2.3 | MIT | 2026-10-30 | 200; 148/month | yes | yes (all profiles) | yes (parsing, writing, validation) | yes (104/104) | yes (own PDF layer, document rewritten; classic xref and xref streams) | yes (CII D16B, UBL) | yes (CEN 1.3.16 XSLTs compiled to PHP, KoSIT scenarios) | FeRD 1.09.2 schematron for the 5 profiles | yes (FNFE-MPE XP Z12-012 1.4.0), `frenchRules` option, off here | yes (XRechnung 3.0.2, KoSIT levels) | yes (BIS 3.0.21 rules rewritten, UBL and CII) | no | dompdf |
| quittio/core | 1.1.0 (+ country-de 1.1.0) | MIT | 2026-11-17 | 361; 319/month (country-de 320; 301/month) | yes | yes | yes (parsing, writing) | yes (104/104) | yes (mPDF + FPDI; no xref streams) | yes (CII D16B) | partial (68 BR rules in PHP, no offline schematron; CII-SR-465 not seen) | MINIMUM to EXTENDED profiles | no | partial (23 BR-DE in PHP; schematron via KoSIT, Java) | partial (PeppolBis3 profile, a few rules in PHP) | veraPDF (Docker or binary) only | mpdf, setasign/fpdi, brick/money, nesbot/carbon, league/commonmark; Java for KoSIT |
| pnscripts/pn-invoice | 0.2.1 | MIT | 2026-10-16 | 187; 184/month | no | yes | yes (writing, UBL 2.1 XSD) | no | no | yes (CII D16B, UBL 2.1) | partial (104 BR rules in PHP; no BR-CL, no BR-DEC, no syntax rules) | no | no | no | no | no | none (bcmath); Saxon on Java for the optional schematron layer |
| pat-o-dev/factur-x | 0.8.0 | MIT | 2026-07-31 | 84; 1/month | no | yes (EN 16931) | no | partial (1/104: classic xref only) | partial (incremental update; refuses /AF, /Names and xref streams) | no | partial (structure, BR-CO-15, BR-CO-25, BR-CO-5-6) | EN 16931 only | partial (BR-FR-MAP-12/14, not counted) | no | no | no | none |
| stann/factur-x | 2.0.0 | none declared | 2025-12-23 | 6,078; 564/month | no | partial (header only) | no | no | yes (mPDF + FPDI, PDF/A-3B; no xref streams) | yes (Factur-X 1.0.06, 5 profiles) | no | XSD for the 5 profiles | no | no | no | no | mpdf (cache in `vendor/`) |
| smalot/pdfparser (floor) | 2.13.1 | LGPL-3.0 | 2026-10-13 | 60,612,162; 6,173,943/month | — | — | — | yes (104/104) | no | — | — | — | — | — | — | no | symfony/polyfill-mbstring |
| ext-dom (floor) | PHP 8.6.0, libxml2 2.9.13 | PHP-3.01 | — | — | no | no | no | no | no | yes (atgp's Factur-X 1.09.2 XSD) | no | XSD for the 5 profiles | no | no | no | no | none |

Downloads: Packagist API (`downloads.total` and `downloads.monthly`, last
30 days), as of 2026-11-22; licenses and dates: metadata of the latest
release published on Packagist.
