# Coverage: Java ecosystem

Java bench for the common protocol (`../PROTOCOL.md`). One Maven module per
library (`libs/<lib>`), a shared harness (`harness`), versions pinned in each
`pom.xml`, built with `mvn -C` (strict checksums). JDK: OpenJDK 27.0.1.1
(Homebrew), target `release 21`; Maven 3.10.1. Each library runs in its own
JVM, with the default options.

The evidence cited is paths inside the jars under `~/.m2/repository/…` and
their source jars (`-sources.jar`), or inside the official artifacts that
`./run.sh setup` downloads into `.artefacts/`. Popularity: GitHub stars
(GitHub API, as of 2026-11-22), latest version and date (Maven Central,
`maven-metadata.xml`), dependents of the measured version (deps.dev). Maven
Central does not publish download counts: the "Downloads" column of the
table gives stars and dependents instead.

## Measured libraries

### Mustang 2.26.0 (ZUGFeRD/mustangproject)

`org.mustangproject:library` and `org.mustangproject:validator`, the reference
library of the ZUGFeRD community. Invoice model, CII and UBL parsing, CII
writing for all profiles, embedding into a PDF/A, validator (XSD,
schematron, veraPDF).

- Validator: `XMLValidator.java` picks the rules from BT-24.
  - Factur-X profiles: the profile's XSD
    (`schema/ZF_250/<profile>/FACTUR-X_*.xsd`, `XMLValidator.java:256-277`)
    and the profile's Factur-X 1.09.2 schematron
    (`xslt/ZF_250/FACTUR-X_*.xslt`, 433 assertions for EN16931). The CEN
    rules are recoded there under `FX-SCH-*` identifiers. BR-CL-25 becomes
    FX-SCH-A-000571 there (code list for `@schemeID`,
    `FACTUR-X_EN16931.xslt:6420`).
  - XRechnung: Factur-X EN16931 XSD, then CEN CII 1.3.16
    (`xslt/en16931schematron/EN16931-CII-validation.xslt`, 806 assertions,
    `XMLValidator.java:401-404`), then KoSIT 2.4.0 for XRechnung 3.0.2
    (`xslt/XR_30/XRechnung-CII-validation.xslt`, 70 assertions, `:574`), all
    as errors. On the BASIC and EN16931 profiles, the XRechnung schematron
    runs too, but its findings stay at `notice`.
  - BR-FR Flux 2 V1.3.0
    (`xslt/XP_Z12_012/20260216_BR-FR-Flux2-Schematron-CII_V1.3.0.xsl`)
    only if the seller and the buyer are in France (`:233`, `:393`). That is
    never the case in this corpus.
  - Arithmetic check of the totals (`:412`).
  - CII-SR-465 keeps CEN's `warning` flag. Mustang therefore flags 2 files
    out of 26 (FX-SCH-A-000571 on the two BR-CL-25 files) and does not flag
    the 4 XRECHNUNG files.
- Native CheckPDF: `ZUGFeRDValidator` on a PDF runs veraPDF
  (`PDFValidator.java:29-34`), checks the XMP and the attachment, then
  validates the extracted XML.
- Parsing: `ZUGFeRDImporter` also parses the XML into the model at import time
  (`ZUGFeRDInvoiceImporter.java:292-313`, `355-368`). `Extract_*` therefore
  includes that parsing, as the note says. `extractInvoice` rejects 3 files
  out of 26 whose amount due it cannot reproduce (ArithmeticException:
  `EN16931_1_Teilrechnung`, `EN16931_Kraftfahrversicherung_Bruttopreise`,
  `EN16931_Sachversicherung_berechneter_Steuersatz`).
- Writing: `ZUGFeRD2PullProvider`, EN16931 profile. UBL: parsing only; UBL
  writing only serves the despatch advice (`UBLDAPullProvider`).
- Embedding: `ZUGFeRDExporterFromPDFA` accepts a PDF/A-1 or PDF/A-3 host and
  rejects the others (`ZUGFeRDExporterFromPDFA.java:59-75`). The file is
  saved without compression (`ZUGFeRDExporterFromA3.java:361`). The resulting
  PDF is PDF/A-3U and passes the Mustang validator. Yet it keeps the host's
  leftover `factur-x.xml` file specification in `/AF` next to the new one
  (see the notes): Paniolo rejects it (FX-MULTIPLE-INVOICES).
- Classpath of 78 jars: Saxon-HE 12.10, PDFBox 3.0.8 and preflight, veraPDF
  validation-model-jakarta 1.30.2 (GPL-3.0 or MPL-2.0), ph-schematron 8.0.6,
  jaxb-runtime 4.0.2.
- Network: none.

### phive 12.3.0 and phive-rules 4.7.1 (phax)

`com.helger.phive:phive-xml` and the rule packages
`com.helger.phive.rules:phive-rules-{en16931,xrechnung,zugferd}`. A generic
validation engine, with the official XSDs and XSLTs precompiled, run by
ph-schematron 10.3.2 on Saxon-HE 12.10. It is a validator only.

- Rule sets registered by `EN16931Validation.initEN16931`,
  `XRechnungValidation.initXRechnung` and `ZugferdValidation.initZugferd`.
  The adapter picks the set from BT-24.
- EN 16931: `VID_CII_1317` (`EN16931Validation.java:122`), i.e. the CII D16B
  XSD (`:598`) and `external/schematron/1.3.17/cii/EN16931-CII-validation.xslt`
  (810 assertions).
- XRechnung: `VID_XRECHNUNG_CII_302`, i.e. the D16B XSD, CEN CII 1.3.16 and
  `external/schematron/3.0.2/XRechnung-CII-validation.xslt` (KoSIT 2.6.0, 74
  assertions). phive applies its own levels
  (`XRechnungValidation.java:1042-1063`: BR-CL-23 and BR-CL-21 as warnings,
  CII-SR-452/453/454 as errors), not those of the KoSIT configuration.
  CII-SR-465 therefore stays a warning, and the 4 XRECHNUNG files pass.
- Factur-X: ZUGFeRD 2.5.2 / Factur-X 1.09.2 XSDs and schematrons for the five
  profiles (`external/schemas/2.5.2/…`, `external/schematron/2.5.2/…`,
  `ZugferdValidation.java:211-223`). No file in the corpus uses them.
- Packages not loaded: `phive-rules-peppol` and `phive-rules-france` 4.7.1
  (French rules, excluded by the protocol), among others.
- ph-schematron produces its SVRL through JAXB. The module therefore adds
  `jaxb-runtime` 4.0.9: without it, every file fails.
- Network: none.

### KoSIT validator 1.6.3 and validator-configuration-xrechnung v2026-08-31

`org.kosit:validator`, KoSIT's generic engine, driven by a scenario
configuration. It is the reference for XRechnung validation in Germany.

- Official configuration downloaded by `./run.sh setup` (SHA-256 checked),
  file `scenarios.xml`:
  - XRechnung CII scenario (`:335-400`, selected by BT-24 at `:357`): CII
    D16B XSD (`resources/cii/16b/xsd`), CEN EN16931-CII 1.3.16 (806
    assertions), then KoSIT XRechnung 2.6.0
    (`resources/xrechnung/3.0.2/xsl`, 74 assertions: 25 BR-DE, 7 BR-DEX,
    5 BR-DE-CVD, 5 BR-TMP, 20 PEPPOL-EN16931).
  - Custom levels of the scenario (`:384-400`): CII-SR-465 and CII-SR-466
    are raised to error, CII-SR-475 and CII-SR-476 lowered to information.
  - The "EN16931 (CII)" scenario (`:652`, selection at `:672`) only accepts
    BT-24 = `urn:cen.eu:en16931:2017` and applies XSD and CEN. No scenario
    recognizes a Factur-X profile (MINIMUM, BASIC…): no Factur-X rules.
- `DefaultCheck` applies the custom levels (`DefaultCheck.java:145-162`).
  The adapter counts the XSD violations that are not warnings and the
  assertions whose effective level is error or fatal: 6 files out of 26, the
  reference.
- `ProcessorProvider` forbids resolving external resources.
  `jaxb-runtime` 4.0.8 must be supplied: it is an optional dependency of the
  KoSIT POM.
- Version 2.0.0-tp2 exists, but it is a technical preview: we measure
  1.6.3.
- Network: none at run time.

### Saxon-HE 13.1 and official XSLTs (floor)

The "official artifacts" floor, hand-written in
`libs/saxon-schematron/src/main/java/bench/Main.java`. It is a port of the
Python bench's saxonche adapter.

- Saxon-HE (s9api) compiles `cen-1.3.17/xslt/EN16931-CII-validation.xslt`
  and, if BT-24 designates XRechnung,
  `xrechnung-2.6.0/schematron/cii/XRechnung-CII-validation.xsl`.
- The adapter counts the `svrl:failed-assert` and `svrl:successful-report`
  with a fatal or error flag, after applying the custom levels of KoSIT
  v2026-08-31, copied from `scenarios.xml:384-400`.
- No XSD. Result: 6 files out of 26.

### Paniolo 1.0.0 (NaluKit/paniolo)

`io.github.nalukit:paniolo`, `paniolo-facturx` and `paniolo-pdfa`, released on
2026-11-18. EN 16931 model, CII and UBL parsing and writing, validation with
the bundled official artifacts (Saxon-HE 13.1), PDF through `paniolo-facturx`
(PDFBox 3.0.8), PDF/A through `paniolo-pdfa` (veraPDF validation-model
1.30.4, under MPL-2.0 according to the POM).

- Artifacts bundled in `paniolo-rules-1.0.0.jar`
  (`io/github/nalukit/paniolo/rules/…`):
  - `xrechnung-3.0.2-validator-configuration-2026-08-31.zip`, byte for byte
    the official zip (same SHA-256 as in `run.sh`), and the older
    configurations from 1.2.2 to 2.3.1;
  - `factur-x-1.09.2-zugferd-2.5.2.zip`;
  - `fnfe-france-rfe-1.4.0.04-cii-extended-ctc-fr.zip`;
  - Order-X and ZUGFeRD 1.0.
- Validation (`Paniolo.validate`, `Paniolo.java:1710`), chosen from BT-24
  (`ValidationProfile.java:573-601`):
  - XRechnung: the KoSIT scenario with its levels, like the KoSIT
    validator: 6 files out of 26.
  - Plain EN 16931: the same KoSIT configuration. Findings specific to the
    CIUS (BR-DE, BR-TMP, PEPPOL-EN16931) are downgraded to information
    (`Paniolo.java:531-560`, `:1590-1602`).
  - Factur-X profiles: Factur-X 1.09.2 XSDs and schematrons.
  - EXTENDED-CTC-FR: FNFE-MPE 1.4.0.04 rules (EXTENDED-CTC-FR and BR-FR
    Flux 2), only for that profile (`Paniolo.java:3706`). This profile
    never appears in the corpus.
  - Always: an amount check (`CalculationCheck`).
- Parsing: `read` in STRICT mode (`Paniolo.java:1209`) rejects 2 files out of
  26, in which a line's net amount differs from quantity × price
  (`EN16931_Kraftfahrversicherung_Bruttopreise`,
  `EN16931_Sachversicherung_berechneter_Steuersatz`).
- Native CheckPDF: `validate` on a PDF checks the attachment, the
  relationship and the XMP, the PDF/A through veraPDF, then the XML against
  the declared profile.
- Embedding: `embedXml` (`Paniolo.java:2078-2130`) checks BT-24, embeds the
  XML, then revalidates the resulting PDF (this revalidation is included in
  the measurement). No PDF/A conversion.
  - It rejects a PDF that already contains an invoice
    (`PdfBoxFacturXSupport.java:1190-1192`) and, to find out, also counts
    the files in `/AF` (`:248-291`).
  - Both protocol hosts keep a `factur-x.xml` file specification in `/AF`:
    hence the `fail` on `Embed_Small` and `EmbedClassic_Small`.
  - Replacement (`replaceInvoice`) exists only in the SPI.
- Network: none.

### ph-cii 10.1.0 (phax)

`com.helger.cii:ph-cii-d16b`: JAXB classes generated from the CII D16B XSD
(`external/schemas/d16b/…/CrossIndustryInvoice_100pD16B.xsd`,
`CCIID16B.java:36,63`) and the marshaller
`CIID16BCrossIndustryInvoiceTypeMarshaller` (`:39-40`). It is a syntactic
model, not a semantic one.

- Parse and Write: the marshaller without a schema (`setUseSchema(false)`).
- Validation: the same marshaller with its XSD, enabled by default
  (`GenericJAXBMarshaller.java:93`), and `setCollectErrors`. XSD only, no
  business rules: 0 files flagged out of 26.
- `jaxb-runtime` 4.0.6 must be supplied (the version in the ph-commons-parent
  12.1.0 BOM).
- No PDF. UBL: separate project (ph-ubl), not measured.

### invoice4j 0.5.0 (jcodeforge)

`io.github.jcodeforge:invoice4j-zugferd`: ZUGFeRD / Factur-X model, reader
and writer, PDF reading and writing.

- The profile's Factur-X 1.09 XSD
  (`/xsd/cii/en16931/Factur-X_1.09_EN16931.xsd`,
  `CiiEn16931XsdValidator.java:20`; `/xsd/{minimum,basic,basicwl,extended}/…`).
- Rules: `ZugferdBasicValidator`, 3 rules (BR-CO-15, BR-CO-26, BR-S-2),
  stopping at the first error. The CEN XSLT (`xsl/EN16931-CII-validation.xsl`
  in invoice4j-base) and the Factur-X schematrons
  (`schematron/FACTUR-X_*.sch`) are in the jars, but no code references them.
- The `SellerParser` reader (invoice4j-base,
  `cii/parser/SellerParser.java:16-31`) reads neither the seller's
  identifier, nor its VAT number, nor its registration number. Consequences:
  - BR-CO-26 (`BrCo26Validator.java:34-47`) fails on the 11 files parsed;
  - the output of Write violates BR-S-02 and BR-CO-26 under the CEN
    schematron.
- The reader rejects 15 files out of 26, even without XSD:
  - unknown document types 204, 387 and 575 (`DocumentTypeCode.java:118`);
  - identifier scheme 9958;
  - in-house rules: BT-70 required, negative BT-115 and BT-116 amounts
    rejected, bank account required for a SEPA direct debit;
  - the 4 XRECHNUNG files ("Unsupported ZUGFeRD profile",
    `ZugferdProfileDetector`).
- PDF: there are only `ZugferdPdfReader.read(File)` and
  `ZugferdPdfWriter.write(File, File)`, no API on in-memory bytes (protocol
  rule 4). Extract and Embed are therefore `unsupported`.
- UBL: reader and writer in `invoice4jbase/ubl`, not measured.

### Imixs e-invoice 1.0.3 (imixs/e-invoice)

`org.imixs.util:imixs-e-invoice`: a read-and-modify wrapper over the DOM of a
CII, UBL or KSeF (Poland) invoice, written for the Imixs workflow engine.

- `EInvoiceModelFactory.read` detects the format and parses the DOM, then
  `EInvoiceModelCII.parseContent` (`:81`) reads the header, the parties and
  the lines.
- `getContent` reserializes the DOM (indented): Write is a reserialization,
  not a generation.
- No validation, no PDF.
- License: the POM declares the GPL (gpl-2.0 URL, "All the results of the
  imixs workflow project are subject to the GPL"); the GitHub repository,
  MIT.

### gallop 3.0.0 (codebarista-de/gallop)

`de.codebarista:gallop`: a generator of XRechnung and Factur-X EN 16931 CII
from a fluent model, on the JDK DOM. No reader, no validation, no PDF.

- `Write_Small` uses a hand-built model based on `EN16931_Einfach`
  (`einfach()` in `libs/gallop/src/main/java/bench/Main.java`).
  - gallop always writes a deliver-to party BG-13
    (`CIIXMLEInvoiceWriter.java:261-270`) and a payment means BG-16
    (`:297`, `:392-395`). The model puts the buyer's address and code 1
    there, otherwise BR-57 and BR-CL-16 fail.
  - No field for GTINs.
  - The output passes the CEN 1.3.17 schematron.
- `Write_Lines` is skipped.
- Formats: `EInvoiceFormat.XRECHNUNG` and `EN16931_CORE`, of which `ZUGFERD`
  and `FACTURX` are aliases.

### PDFBox 3.0.8 (floor)

Apache's general-purpose PDF library, the one Mustang, Paniolo and invoice4j
use. It serves as the extraction floor, hand-written in
`libs/pdfbox/src/main/java/bench/Main.java`:

- `Loader.loadPDF`, then a walk of the `/EmbeddedFiles` tree with `/Kids`;
- first file under a standard name, case-insensitive (`factur-x.xml`,
  `zugferd-invoice.xml`, `xrechnung.xml`, `order-x.xml`, `cii.xml`), like the
  Rust bench's lopdf floor;
- `/UF` stream, otherwise `/F`.

## Excluded libraries

| Library | Version | Reason |
|---|---|---|
| Konik | 1.3.2 | Archived, ZUGFeRD 1 only, last commit in 2022 (★60) |
| iText pdfInvoice (`com.itextpdf:zugferd`) | 2.0.1 (2018-08-29) | Archived, ZUGFeRD 1.0, AGPL (★6) |
| SRSoftware document.zugferd | 1.0.5 | Plugin for the SRSoftware document engine (`DocumentFactory`, `render(Map)`): generates from a template and embeds through a copy of Mustang's `A3Exporter`; neither parsing nor validation |
| codinux e-invoice-jvm | 0.8.1 | Kotlin layer over Mustang 2.15.1 (Mustang is measured directly); proprietary license "codinux GmbH & Co. KG" |
| Yona einvoice-java (`com.useyona`) | 0.2.3 | Client for the remote Yona API (excluded by the protocol) |
| en16931-basics (phax) | 1.1.0 | Shared code lists and utilities, no invoice input/output (★1) |
| en16931-cii2ubl (phax) | — | CII → UBL converter; its CII parsing is that of ph-cii, already measured (★48) |
| OpenXRechnungToolbox | — | GUI application under GPL-3.0, not a library (★172) |
| veraPDF | 1.30.x | General-purpose PDF/A validator, measured through Mustang and Paniolo (★352) |

All the libraries kept are pure Java, with no native dependency, and none
calls the network at run time.

## Accuracy surveys (`./run.sh survey`)

Survey taken with a temporary `OUT`, outside `results/`.

| Library | validate (flagged / 26) | parse (/ 26) | extract (/ 104) |
|---|---|---|---|
| mustang | 2 (FX-SCH-A-000571 ×2, recoded BR-CL-25; CII-SR-465 stays a warning) | 23 (3 ArithmeticException on the amount due) | 104 |
| phive | 2 (BR-CL-25 ×2; phive's own levels, CII-SR-465 as a warning) | — | — |
| kosit-validator | **6** (BR-CL-25 ×2, CII-SR-465 ×4 on the XRECHNUNG files) | — | — |
| saxon-schematron | **6** (same) | — | — |
| paniolo | **6** (same) | 24 (2 line net amount mismatches, STRICT mode) | 104 |
| ph-cii | 0 (XSD only) | 26 | — |
| invoice4j | 26 (11 BR-CO-26 false positives, 15 parse rejections) | 11 | — |
| imixs-e-invoice | — | 26 | — |
| gallop | — | — | — |
| pdfbox | — | — | 104 |

Reference (PROTOCOL): 6 out of 26.

## Benchmark status (trial run `COUNT=1 BENCHTIME=0.05`)

`o` ok, `F` fail (message in the JSONL line), `-` unsupported. Columns in
order:

1. Extract S/L/Min/Ext
2. Parse S/Li/La
3. ParseValidate S/Li/La
4. Write S/Li
5. Embed, EmbedClassic
6. CheckPDF S/L

| Library | Extract | Parse | ParseValidate | Write | Embed | CheckPDF |
|---|---|---|---|---|---|---|
| mustang | o o o o | o o o | o o o | o o | o o | o o |
| phive | - - - - | - - - | o o o | - - | - - | - - |
| kosit-validator | - - - - | - - - | o o o | - - | - - | - - |
| saxon-schematron | - - - - | - - - | o o o | - - | - - | - - |
| paniolo | o o o o | o o o | o o o | o o | F F | o o |
| ph-cii | - - - - | o o o | o o o | o o | - - | - - |
| invoice4j | - - - - | o o F | o o F | o o | - - | - - |
| imixs-e-invoice | - - - - | o o o | - - - | o o | - - | - - |
| gallop | - - - - | - - - | - - - | o (skipped) | - - | - - |
| pdfbox | o o o o | - - - | - - - | - - | - - | - - |

## Coverage table

| Library | Version | License | Latest version | Downloads | CII parsing | CII writing | UBL | PDF extraction | PDF embedding | XSD | EN 16931 (CII) | Factur-X profiles (MINIMUM, BASIC WL, EXTENDED) | BR-FR | BR-DE (XRechnung) | Peppol BIS 3 | PDF/A-3 container check | External dependencies |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Mustang (library + validator) | 2.26.0 | Apache-2.0 (veraPDF: GPL-3.0 or MPL-2.0) | 2.26.0 (2026-08-25) | ★497; dependents: library 3, validator 1 | partial (rejects 3/26 whose amount due it cannot reproduce) | yes (all profiles, XRechnung) | parsing; writing: despatch advice only | yes (also parses into the model) | yes (PDF/A-1 or A-3 host required, PDF/A-3U produced; keeps the host's leftover `/AF`) | yes: Factur-X 1.09.2 per profile | yes: CEN rules recoded in the Factur-X EN16931 schematron; CEN CII 1.3.16 XSLT on top for XRechnung | yes: Factur-X 1.09.2 schematrons for the 5 profiles | yes: BR-FR Flux 2 V1.3.0, only if seller and buyer are French | yes: XRechnung 3.0.2, KoSIT 2.4.0 schematron (70 assertions); as notice on BASIC and EN16931 | no | yes: veraPDF, XMP, attachment | none (pure Java: PDFBox, Saxon-HE 12.10, veraPDF) |
| phive + phive-rules | 12.3.0 / 4.7.1 | Apache-2.0 (phive-rules: "Apache 2" in the POM, NOASSERTION on GitHub; CEN artifacts EUPL-1.2) | 12.3.0 (2026-11-15) / 4.7.1 (2026-11-17) | ★79 / ★55; dependents: not indexed | no (validator only) | no | validation (UBL EN 16931 and XRechnung sets) | no | no | yes: CII D16B; Factur-X 1.09.2 per profile | yes: official CEN 1.3.17 XSLT (810 assertions) via ph-schematron and Saxon-HE | yes: ZUGFeRD 2.5.2 / Factur-X 1.09.2, 5 profiles | not loaded (separate phive-rules-france package) | yes: XRechnung 3.0.2, KoSIT 2.6.0 schematron (74 assertions), phive's own levels (CII-SR-465 as a warning) | not loaded (separate phive-rules-peppol package) | no | none (pure Java; jaxb-runtime must be supplied) |
| KoSIT validator + XRechnung configuration | 1.6.3 / v2026-08-31 | Apache-2.0 (engine and configuration) | 1.6.3 (2026-08-20); 2.0.0-tp2 as a preview | ★199 (configuration ★120); dependents: 0 | no (validator only) | no | validation (UBL scenarios of the configuration) | no | no | yes: CII D16B (configuration) | yes: CEN CII 1.3.16 XSLT (806 assertions) from the configuration, via Saxon-HE 12.9 | no (no Factur-X scenario) | no | yes: KoSIT 2.6.0 schematron and custom levels (25 BR-DE, 7 BR-DEX, 5 BR-DE-CVD) | no (20 PEPPOL-EN16931 rules taken over by XRechnung) | no | official configuration downloaded by `setup` (SHA-256); jaxb-runtime must be supplied |
| Saxon-HE + official XSLTs (floor) | 13.1 | MPL-2.0 (CEN artifacts EUPL-1.2, KoSIT Apache-2.0) | 13.1 (2026-10-15) | ★178; dependents: 30 | no | no | no | no | no | no | official CEN CII 1.3.17 XSLT via Saxon | no | no | official KoSIT 2.6.0 XSLT via Saxon, KoSIT v2026-08-31 levels | no | no | official artifacts downloaded by `setup` (SHA-256) |
| Paniolo (+ facturx, pdfa) | 1.0.0 | Apache-2.0 (veraPDF validation-model: MPL-2.0) | 1.0.0 (2026-11-18) | ★1; dependents: not indexed | partial (STRICT rejects 2/26: line amount mismatches) | yes | yes (parsing, writing) | yes | partial: no PDF/A conversion, revalidates the output; rejects a host that already carries an invoice (both protocol hosts) | yes: CII D16B (KoSIT configuration), Factur-X 1.09.2 per profile | yes: CEN CII 1.3.16 XSLT from the bundled KoSIT v2026-08-31 configuration | yes: Factur-X 1.09.2 / ZUGFeRD 2.5.2 schematrons | yes, for the EXTENDED-CTC-FR profile only (FNFE-MPE 1.4.0.04, BR-FR Flux 2) | yes: KoSIT XRechnung 3.0.2 v2026-08-31 configuration with its levels; findings downgraded to information outside XRechnung | no | yes: veraPDF (paniolo-pdfa), attachment, relationship, XMP | none (pure Java: Saxon-HE 13.1, PDFBox 3.0.8, veraPDF 1.30.4) |
| ph-cii (d16b) | 10.1.0 | Apache-2.0 | 10.1.0 (2026-11-14) | ★15; dependents: 0 | yes (JAXB, D16B syntactic model) | yes (JAXB) | no (separate ph-ubl project) | no | no | yes: CII D16B | no | no | no | no | no | no | none (jaxb-runtime must be supplied) |
| invoice4j-zugferd | 0.5.0 | Apache-2.0 | 0.5.0 (2026-09-18) | ★4; dependents: 0 | partial (11/26: types 204/387/575, scheme 9958, in-house rules, XRechnung rejected) | partial (seller identifiers lost on parsing) | yes (invoice4j-base, not measured) | partial (file-based API only) | partial (file-based API only) | yes: Factur-X 1.09 per profile | partial (3 rules: BR-CO-15, BR-CO-26, BR-S-2); bundled CEN XSLT never run | XSD only (bundled Factur-X schematrons never run) | no | no (XRechnung profile rejected) | no | no | none (PDFBox 3.0.8) |
| Imixs e-invoice | 1.0.3 | GPL in the POM, MIT on GitHub | 1.0.3 (2026-09-28) | ★4; dependents: 0 | yes (DOM wrapper, main fields) | partial (reserializes the parsed or modified DOM) | yes (parsing; KSeF too) | no | no | no | no | no | no | no | no | no | none |
| gallop | 3.0.0 | MIT | 3.0.0 (2026-09-25) | ★0; dependents: 0 | no | yes (EN 16931 and XRechnung from a fluent model; always writes BG-13 and BG-16) | no | no | no | no | no | no (EN 16931 only) | no | XRechnung writing only | no | no | none (JDK DOM) |
| PDFBox (floor) | 3.0.8 | Apache-2.0 | 3.0.8 (2026-07-08) | ★3,143; dependents: 392 (119 direct) | no | no | no | hand-written floor | no | no | no | no | no | no | no | no | none |

## Build and measurement notes

- One Maven module per library: the classpaths stay separate (Mustang pulls
  in Saxon-HE 12.10, KoSIT 12.9, Paniolo 13.1). `run.sh` launches one JVM per
  library, with the default options: neither `-Xss` nor `-Xmx`.
- The harness runs the operations in a thread with a 512 MiB stack created by
  the program (`new Thread(…, 512L << 20)`), like the Rust bench. No validator
  needed a JVM option.
- `mvn -B -C -q package`: `-C` makes the build fail on an invalid checksum.
  The `package` phase writes each module's classpath to
  `target/classpath.txt` (`maven-dependency-plugin:build-classpath`), which
  `run.sh` passes to `java -cp`. Only Apache Maven plugins are used, at the
  versions pinned in the parent POM.
- JAXB: KoSIT, ph-cii and phive (through ph-schematron) need an
  implementation of Jakarta XML Binding that their POMs do not pull in.
  `org.glassfish.jaxb:jaxb-runtime` is therefore added with `runtime` scope,
  at the version each project manages: 4.0.8 for KoSIT, 4.0.6 for ph-cii,
  4.0.9 for phive.
- Logs: there is no SLF4J provider, hence a NOP warning at startup. In the
  Mustang module, JDK 27 also warns that a final field is modified through
  reflection (JAXB on veraPDF classes). Neither changes the results.
- `first_call_ns` includes the lazy compilation of the XSLTs (Mustang, phive,
  KoSIT, Saxon, Paniolo) and class loading: from 0.6 to 2 s for a
  validator.
- CPU time (`cpu_ns_per_op`): `OperatingSystemMXBean.getProcessCpuTime()`
  (`com.sun.management`), i.e. `getrusage` on macOS, all JVM threads, read
  just before and just after each series. No library creates a thread or a
  process. A `jcmd Thread.print` dump taken while each library was measured
  shows only the `bench` thread, `main` and the JVM's own threads: 2 to 4
  C1/C2 compilers, up to 11 parallel and 3 concurrent G1 threads, and Java2D
  Disposer in Mustang (embedding). CPU time therefore exceeds wall-clock time
  while the JIT is still compiling or while G1 collects on several threads.
  In the `bench` pass, PDFBox is at +1 or 2%. In a fresh JVM (`resources`
  pass), it rises to ×3.5 with `BENCHTIME=0.2` and comes back to +2% with
  `BENCHTIME=1`; at 1 s, Saxon and Mustang still range from ×1.1 to ×3.4
  depending on the benchmark.
- `./run.sh resources`: one JVM per library and per benchmark
  (`bench.Main resources <Bench>`), same command line and same options as
  `bench`, launched directly under `/usr/bin/time -l` (peak RSS in bytes).
  With the default options (G1, 8 MiB initial heap, 6 GiB maximum heap, a
  quarter of the machine's 24 GiB), this peak mostly tracks how G1 sizes the
  heap under heavy allocation, not the library's live data. Example: PDFBox,
  2 MiB live in a 448 MiB heap, 300 to 570 MiB of RSS; Mustang, up to 2 GiB.
- The harness output is formatted with `Locale.ROOT`.
- Write outputs checked again with the Saxon floor (CEN 1.3.17) on
  `EN16931_Einfach` and `EN16931_Rabatte`: Mustang, Paniolo, ph-cii, Imixs
  and gallop raise no errors; invoice4j violates BR-S-02 and BR-CO-26.
- Network: none at run time. `./run.sh setup` downloads over HTTPS, verifies
  by SHA-256 and unpacks into `.artefacts/` the three official artifacts: CEN
  `validation-1.3.17` (CII), KoSIT `xrechnung-schematron` v2.6.0 and KoSIT
  `validator-configuration-xrechnung` v2026-08-31. Maven downloads the jars.
- Hosts: `host.pdf` and `host_classic.pdf` no longer have an `/EmbeddedFiles`
  tree. Yet their catalog keeps in `/AF` a `factur-x.xml` file specification
  (`/AFRelationship /Alternative`) with its `/UF` stream, as well as a
  Factur-X XMP in PDF/A-3. Consequences:
  - Paniolo refuses to embed an invoice into them;
  - Mustang produces a PDF that contains two `factur-x.xml` in `/AF`.
