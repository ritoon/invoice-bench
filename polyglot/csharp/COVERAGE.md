# Coverage: C# / .NET ecosystem

C# bench for the common protocol (`../PROTOCOL.md`). One console project per
library (`libs/<lib>`, assembly `bench-<lib>`), a shared harness (`harness`),
and a solution that groups them (`Bench.slnx`). Versions are pinned with
`[x.y.z]`, and each project has a `packages.lock.json`, restored in locked
mode. Toolchain: .NET SDK 11.0.100 (runtime Microsoft.NETCore.App 11.0.0),
target `net11.0`, Release configuration.

The evidence cited is paths in the repositories' sources at the tag of the
published version (ZUGFeRD-csharp 18.0.1, International.EInvoicing v1.1.0,
Balsoft v1.1.0, Securibox V0.6.0-alpha-3, FacturXDotNet v0.5.9-alpha,
pax.XRechnung.NET at the last commit before 0.3.3). For OfficeIMO and
En16931 (sources not cloned) and for Antiphon (closed-source software), the
evidence comes from the contents of the `.nupkg` (README, resources, XML
documentation) and from the public API read from the assemblies' metadata
(MetadataLoadContext, without executing any code). License, latest version
and downloads come from the NuGet API (registration and search, as of
2026-11-22). Downloads are the total across all versions: NuGet does not
publish a 90-day figure.

## Measured libraries

### ZUGFeRD-csharp + ZUGFeRD.PDF-csharp 18.0.1 (Stephan Stapel)

The most downloaded .NET library in the field. `InvoiceDescriptor` model,
ZUGFeRD 1.0, 2.0 and 2.3 (CII) readers and writers, as well as XRechnung UBL.
The `.PDF` package reads and writes the PDF on top of PDFsharp-extended 1.3.0
(`ZUGFeRD.PDF/ZUGFeRD.PDF.csproj:45`).

- End of maintenance announced in the NuGet description: support until
  30 November 2026, "Legacy" package from 1 December 2026.
- Extraction: `InvoicePdfProcessor.LoadFromPdf`
  (`ZUGFeRD.PDF/InvoicePdfProcessor.cs:46`) extracts the attachment and parses
  the model. Names looked up: zugferd.xml, zugferd-invoice.xml, factur-x.xml,
  xrechnung.xml (`ZUGFeRD.PDF/InvoiceDescriptorPdfLoader.cs:103`). No public
  API returns the bytes: `Extract_*` measures extraction plus parsing. The
  check rewrites the model outside the timed section. 103 PDFs out of 104:
  `Avoir_FR_type381_BASIC.pdf` is rejected ("No ZUGFeRD invoice reader was
  able to parse this stream").
- Validation: `InvoiceValidator.Validate` (`ZUGFeRD/InvoiceValidator.cs:50`)
  only recomputes the totals: sum of the lines (l. 69-108), BR-CO-17 and
  BR-CO-14 (l. 148-182), totals and amount due. No XSD or schematron. The
  messages mix the calculation trace with the failures. The adapter keeps the
  failure lines: those that carry a BR- or BT- identifier, or a "berechnet /
  tatsächlich" discrepancy. This filter is heuristic.
- The 3 flagged files (Teilrechnung, Kraftfahrversicherung,
  Sachversicherung) are false positives on `lineTotal` and `grandTotal`.
  Balsoft produces the same ones, for the same reason.
- Writing: `InvoiceDescriptor.Save(Version23, profile read, CII)`.
- Embedding: `InvoicePdfProcessor.SaveToPdf(Version23, Comfort, CII)`
  rewrites the XML from the model and copies the pages into a new PDFsharp
  document. It adds a Factur-X XMP and an sRGB OutputIntent
  (`ZUGFeRD.PDF/InvoiceDescriptorPdfSaver.cs:321-326`), and embeds fonts
  taken from the system folders (`:329`, `:345-`,
  `ZUGFeRD.PDF/FontInfoProvider.cs:89-103`). Dates = `DateTime.UtcNow`
  (`InvoiceDescriptorPdfSaver.cs:185`, `:230`), not configurable.
- No hybrid PDF check: CheckPDF is `unsupported`. Network: none (no
  `HttpClient` in the sources).

### International.EInvoicing 1.1.0 (Epic Dice Studio)

EN 16931 model, UBL, CII, Factur-X / ZUGFeRD and XRechnung readers and
writers, with an in-house Schematron engine: an XPath subset evaluated on the
preprocessed official `.sch` files, without XSLT. Packages used, all at
1.1.0: the facade, `.FacturX.PdfSharp` (PDFsharp 6.2.4, MIT), `.Countries.Germany`,
`.Validation.XRechnung` and `.Validation.Xsd`.

- EN 16931: CEN 1.3.16 artifacts
  (`src/International.EInvoicing.Validation.En16931/En16931Rules.cs:18`),
  806 assertions in the preprocessed CII `.sch`. `AddDefaults` loads UBL, CII,
  CDAR, Factur-X, the EN 16931 rules and the facade
  (`src/International.EInvoicing/EInvoicingBuilderExtensions.cs:27-38`).
- XRechnung: `ArtefactVersion = "3.0"`
  (`src/International.EInvoicing.Validation.XRechnung/XRechnungRules.cs:18`),
  52 assertions. The KoSIT schematron was taken from `master` on 2026-08-29,
  not from a published release (`specs/xrechnung/PROVENANCE.md:6-7`). The
  custom levels of the KoSIT configuration do not apply: CII-SR-465 stays a
  warning, hence 2 flagged files out of 26 instead of 6.
- XSD: UN/CEFACT CII D16B
  (`src/International.EInvoicing.Validation.Xsd/XsdBuilderExtensions.cs:37`).
- Measured configuration: `AddDefaults().AddGermany().AddXRechnungRules().AddCiiSchema()`,
  that is `GermanEInvoicing.Create`
  (`src/International.EInvoicing.Countries.Germany/GermanEInvoicing.cs:42-43`)
  plus the XSD. The XRechnung rules apply only if BT-24 declares it.
- Nothing else is available offline. The Factur-X (FeRD) rules are not
  shipped and are loaded from a folder (`AddFacturXRulesFrom`,
  `src/International.EInvoicing.FacturX/FacturXRules.cs:43`). The French
  rules come from `AddFrance`
  (`src/International.EInvoicing.Countries.France/FranceServiceCollectionExtensions.cs:18`)
  and Peppol from a separate package (`International.EInvoicing.Peppol`). None
  of the three is loaded: no example in the corpus declares Peppol.
- Extraction: `PdfSharpAttachmentReader.FindAttachment` with
  `FacturXAttachment.KnownFileNames`
  (`src/International.EInvoicing.FacturX/Pdf/FacturXAttachment.cs:21-`) and the
  default limit of 64 MiB. It returns the raw bytes.
- Embedding: `PdfSharpAttachmentWriter.Attach` adds factur-x.xml as an
  associated file and adds or completes the Factur-X XMP. ModDate =
  `DateTimeOffset.UtcNow`
  (`src/International.EInvoicing.FacturX.PdfSharp/PdfSharpAttachmentWriter.cs:59`).
- No CheckPDF: reading a PDF compares the XMP with the XML (warnings),
  without PDF/A.
- Network: none in the packages used. The only `HttpClient` in the sources is
  in `web/Playground`, which is not published.

### Balsoft.Hive.EInvoice + Balsoft.Hive.EInvoice.Pdf 1.1.0 (Balsoft GmbH)

Hand-written EN 16931 model, CII and UBL readers and writers, a validator
coded in C# on the model, hybrid PDF creation on PDFsharp 6.2.4.

- Rules: 115 distinct BR identifiers in
  `src/Balsoft.Hive.EInvoice/Validation/InvoiceValidator.cs`, including 26
  BR-DE, plus 4 PEPPOL-EN16931 rules. No XSD or schematron.
- The model does not keep the declared totals. `ReadResult.TotalsDiscrepancies`
  (`src/Balsoft.Hive.EInvoice/Reading/ReadResult.cs`) compares them with the
  recomputed totals, and the adapter adds these discrepancies (BT
  identifiers) to the errors. Result: BR-CL-25 on the 2 expected files, and
  BT-106, BT-109, BT-112 and BT-115 discrepancies on 3 valid examples. These
  are the same false positives as ZUGFeRD-csharp.
- Writing: `CiiWriter.Write` validates before writing by default
  (`src/Balsoft.Hive.EInvoice/Cii/CiiWriter.cs:12-22`). The adapter passes
  `ValidateBeforeWriting=false`, so as not to measure validation twice.
- Extraction: `HybridPdf.ExtractXml`
  (`src/Balsoft.Hive.EInvoice.Pdf/HybridPdf.cs:127`) looks for factur-x.xml,
  xrechnung.xml, zugferd-invoice.xml, otherwise the first `.xml`.
- Embedding: `HybridPdf.Create` (`HybridPdf.cs:83`) copies the pages into a
  new document and adds factur-x.xml as an associated file and an sRGB
  OutputIntent. Missing fonts are taken from the system folders. The XMP is
  replaced through an incremental update. Timestamp = `DateTimeOffset.Now` by
  default.
- No CheckPDF. Network: none.

### Securibox.FacturX 0.6.0-alpha-3 (Securibox)

Factur-X import, export and validation. One `XmlSerializer` model per profile,
generated from the Factur-X 1.09 XSDs. The FeRD Factur-X 1.09 XSDs and
schematrons for the five profiles are embedded (`Securibox.FacturX/Resources/`,
424 assertions in `Factur-X_1.09_EN16931.sch`). In-house Schematron engine:
the `.sch` is deserialized and the assertions are evaluated by XPath2 1.1.5,
without XSLT. PDF on PDFsharp 7.0.0-preview-1 (MIT; preview-2 is AGPL), XMP
on XmpCore 6.1.10.1.

- The csproj at the tag declares `0.6.0-alpha-2`
  (`Securibox.FacturX/Securibox.FacturX.csproj:8`), whereas the tag and the
  package are at 0.6.0-alpha-3.
- Measured validation: `FacturxXsdValidator.ValidateXml` +
  `FacturxSchematronValidator.ValidateXml` with the detected profile.
  The adapter detects the profile the way `FacturxImporter` does
  (`FacturxImporter.cs:574-627`): XRechnung and EN 16931 both map to the
  EN 16931 profile. The `.sch` is reread and deserialized on every call
  (`FacturxSchematronValidator.cs:11-60`), hence about 100 ms per validation.
- The FeRD assertions have no `@id`, so `BusinessRuleCode` is empty
  (`Schematron/Helpers/ValidationReportMapper.cs`). The adapter takes the
  identifier from the start of the message ("[CII-SR-465]"). For code-list
  assertions, which have none, it records `codedb clN`.
- Survey: 6 flagged files, the same as the reference, but not for the same
  reasons:
  - CII-SR-465 on the 4 XRECHNUNG files: the FeRD assertion has no flag and
    therefore counts as an error;
  - `codedb cl1` on the same 4 files: list 1 of the EN 16931 profile contains
    only `urn:cen.eu:en16931:2017`, and mapping XRechnung to this profile
    creates the error;
  - `codedb cl17` on the 2 BR-CL-25 files (list of identifier schemes).
- Extraction: no public API returns the bytes.
  `ImportDataWithDeserialization` (`FacturxImporter.cs:143`) calls
  `IsFacturXValid`, ignores the result, then deserializes. The check
  reserializes outside the timed section. Result: 74 PDFs out of 104. 29 fail
  with "XmpException: Schema namespace URI and prefix mismatch" (XMP of
  ZUGFeRD 2.0 PDFs, `Extract_Extended` included), 1 with "Page has no
  MediaBox".
- Writing: no public API turns a model into XML without a PDF. The adapter
  reuses `XmlSerializer.Serialize` with the namespaces of `FacturxExporter`.
- Embedding: `FacturxExporter.CreateFacturXStream` validates first and fails
  on error (`FacturxExporter.cs:222-`). It then copies the pages and writes
  the XMP.
- CheckPDF: `IsFacturXValid` (`FacturxImporter.cs:97-140`) checks the XMP,
  the embedded XML, the XSD and the schematron. Errors and warnings both
  count (`:119`). No PDF/A check.
- Network: none.

### OfficeIMO.Invoicing + .Validation + .Pdf 3.6.2 (Evotec)

Typed invoice model, bounded CII and UBL readers and writers (EN 16931
1.3.16, XRechnung 3.0.2, Factur-X 1.09.2 / ZUGFeRD 2.5.2, Peppol BIS 3.0.21),
model checks (semantic, arithmetic), and OfficeIMO.Pdf, a PDF engine written
in C#.

- Validation: `OfficeIMO.Invoicing.Validation` runs the XSD stage in .NET,
  from the KoSIT `validator-configuration-xrechnung` v2026-08-31 archive that
  it pins (`run.sh setup` downloads it). The business rules stage requires
  Java and SaxonJ-HE 12.10 (package README, l. 23 and 77-95:
  `SaxonInvoiceRulesRunner("tools/saxon/saxon-he-12.10.jar")`). It is left
  out, since there is no JVM. What is measured: model checks
  (`InvoiceModelValidator.ValidateSource`) and XSD (CII D16B, or XRechnung
  3.0.2 if BT-24 declares it).
- The model checks are strict and are not EN 16931 rules. They flag 23 files
  out of 26: INV-WHITESPACE (55 occurrences, leading or trailing whitespace
  in text), INV-PRICE (8), INV-ADJUSTMENT-BASE (6), INV-DECLARED-AMOUNT,
  INV-ATTACHMENT-URI… The XSD stage flags nothing.
- Writing: `InvoiceReadResult.Write` refuses to write a model that fails
  these checks. `Write_Small` and `Write_Lines` are therefore `fail`
  (INV-WHITESPACE).
- Extraction: `PdfDocument.Load` + `ExtractAttachments` (`/Names/EmbeddedFiles`
  and `/AF`, decoded streams).
- Embedding: the Factur-X helpers (associated file and XMP) exist only on
  `PdfCreateOptions`, for the PDFs that OfficeIMO generates. The classes that
  the XML documentation mentions (`PdfAttachmentExtractor`,
  `PdfComplianceAnalyzer`) are internal: the public API was read from the
  metadata. For an existing PDF, the adapter goes through
  `Attachments.Edit`, with a generic associated file, without Factur-X XMP.
  - On `host.pdf`, the result is `PdfMutationBlockedException`: xref stream
    and object streams not rewritten, incremental writing not implemented
    for attachments.
  - On `host_classic.pdf`, embedding succeeds.
- CheckPDF: `AssessCompliance(PdfComplianceProfile.FacturX)` rereads the file
  for the PDF/A-3 and Factur-X requirements (XMP, associated file,
  OutputIntent, fonts…). The XML is not validated. It is the only container
  check in the ecosystem that looks at PDF/A.
- Network: none (the KoSIT archive is read locally).

### FacturXDotNet 0.5.9-alpha (FacturX-NET)

Factur-X documents (PDF, XMP and factur-x.xml). The CII model is limited to
the MINIMUM profile header: seller, buyer, totals. TurboXml reader, writer,
rules in C#, PDF on the FacturXDotNet.PDFsharp 6.2.0-preview-3 fork.
Asynchronous API: the adapter waits for each task. Latest version on
2025-04-21.

- Rules: 19 CII rules
  (`src/FacturXDotNet/Validation/BusinessRules/CII/`: BR-01 to BR-09, BR-13
  to BR-15, BR-CO-09, BR-CO-26, BR-DEC-12, 13, 14 and 18, BR-S-01) and 16
  hybrid rules (BR-HYBRID-01 to 15, BR-HYBRID-FR-01).
- The reader rejects 15 of the 26 examples. For 11, the `schemeID` attribute
  of the seller's tax identifier ("FC") falls outside its enumeration. For 4,
  it is the XRechnung BT-24 (`ArgumentOutOfRangeException`). `Parse`,
  `ParseValidate`, `Write` and `Embed` therefore fail on the protocol's
  examples.
- Extraction: factur-x.xml only, the default name
  (`FacturXDocument.cs:83`). Result: 70 PDFs out of 104. The other 34
  (ZUGFeRD 2.0, XRechnung, Gnuaccounting) give their attachment a different
  name, `Extract_Extended` included.
- CheckPDF: `FacturXValidator`, hybrid and CII rules, with BR-HYBRID-FR-01
  (French rule) left out. The XML of `EN16931_Einfach.pdf` cannot be parsed,
  hence `fail` on `CheckPDF_Small`. `CheckPDF_Large` passes.
- Network: none.

### pax.XRechnung.NET 0.3.3 (Philipp Hetzner)

XRechnung UBL model, writer and validators. For CII, the only feature is
`ZugferdXmlValidator`, a local ZUGFeRD 2.3.3 / Factur-X 1.07.3 EN 16931
validator. It applies the Factur-X 1.07.3 EN16931 XSD
(`src/pax.XRechnung.NET/ZugferdXmlValidator.cs:111`), then the FeRD XSLT
`FACTUR-X_EN16931.xslt` (`:134`), run by SaxonCS-HE 13.0.0. Both are
embedded.

- Accepts only BT-24 = `urn:cen.eu:en16931:2017` (`:17`, `:44`). The 4
  XRECHNUNG files are rejected with a message (the survey counts them as
  flagged), and `ParseValidate_Large` is skipped.
- Flagged: FX-SCH-A-000031 on the 2 BR-CL-25 files.
- No CII reader or writer, no PDF.
- License: MIT on NuGet, GPL-3.0 in the repository (`LICENSE`).
- Network: the KoSIT validation server client
  (`KositValidator/InvoiceKositValidator.cs`, `HttpClient`) is not called.

### En16931 0.1.0 (Jonas Fassbender)

EN 16931 semantic model in C#, parsed and written in CII and UBL through an
intermediate representation. The XSLTs `IR/cii2ir.xslt` and `IR/ir2cii.xslt`
are shipped as `contentFiles`. Validation runs the official XSLTs with
SaxonCS-HE 13.1.0 (package README, l. 11-16,
`contentFiles/any/any/En16931.Resources.Extern/`):

- CEN 1.3.16;
- KoSIT XRechnung 3.0.2 / schematron 2.6.0, with the extension and CVD
  variants;
- FeRD Factur-X 1.09.2 BASIC and EN 16931.

The CII D16B and D22B and UBL 2.1 XSDs are shipped too.

- Specification chosen from BT-24: `FacturXEn16931` for
  `urn:cen.eu:en16931:2017` (CEN and Factur-X EN 16931 rules), `XRechnung`
  for XRechnung 3.0. No MINIMUM, BASIC WL or EXTENDED. No PDF.
- Parsing: 18 files out of 26.
  - "XmlException: 'Element' is an invalid XmlNodeType" (3), including
    `EN16931_Rabatte`, hence `Parse_Lines` and `Write_Lines` as `fail`.
  - "ReadContentAsString … not supported" (3).
  - `ValidationException` (2): `Parser.Parse` validates first.
- Validation: FX-SCH-A-000571 on the 2 BR-CL-25 files. The XRECHNUNG files
  go through the XRechnung specification without the KoSIT custom levels:
  CII-SR-465 stays a warning.
- Published on 2026-11-05 (a single version, 72 downloads).
- Network: none.

### Antiphon.Invoicing 1.0.0 (Clockwork Otter Foundry)

Commercial, closed-source software. The evaluation license is free and
complete, with no key and no time limit (`Antiphon-Evaluation-License-Notice.md`
in the package), and restricted to evaluation before purchase. All versions
were published on 2026-09-18. No sources: the evidence comes from the public
API read from the metadata and from the packages' XML documentation.

- Package split: `Antiphon.Invoicing` contains only the model (30 types).
  The adapter also uses:
  - `Parsing.Cii` (`CiiInvoiceParser.Parse`) and `Generation.Cii`
    (`CiiInvoiceGenerator.GenerateXml`);
  - `Validation.Pipeline` (`LayeredValidator`), `Rules.Cen.Cii`,
    `Rules.XRechnung.Cii` and `Xsd.Cii`;
  - `Pdf` (`FacturXExtractor`, `FacturXEmbedder`, on PDFsharp 6.2.4).

  `Rules.Peppol.Cii` and the UBL packages are not used here.
- Validation: well-formedness, then CII D16B XSD, then the Schematron sets
  that `FormatDetector` derives from BT-24 (`cii-en16931` → `cen-cii`,
  `cii-xrechnung` → `cen-cii` and `xrechnung-cii`). The sets ship precompiled
  in an intermediate representation and are evaluated by an XPath evaluator
  written in C#, without XSLT or Java. They are deserialized once, on the
  first call.
- Sets shipped: CEN `validation-1.3.16` (806 assertions) and KoSIT XRechnung
  schematron v2.6.0 (74 assertions), according to the `RuleSetInfo`
  constants. There is no Factur-X (FeRD) set and no French set. The KoSIT
  custom levels do not apply: CII-SR-465 stays a warning, hence 2 flagged
  files out of 26 (BR-CL-25 ×2).
- Errors counted: `OfficialFatalViolations`, that is the fatal violations of
  the XSD and of the official sets.
- Parsing: 26 out of 26. Elements that are not mapped are returned
  separately (`UnmappedElements`) rather than lost.
- Writing: `GenerateXml(Invoice, Comfort)` from the parsed model.
- Extraction: `FacturXExtractor.ExtractBytes`, names factur-x.xml,
  zugferd-invoice.xml, ZUGFeRD-invoice.xml, xrechnung.xml and order-x.xml,
  then case-insensitively. 104 PDFs out of 104.
- Embedding: `FacturXEmbedder.Embed`. The documented contract requires a host
  that is already PDF/A: the library only adds the attachment layer, `/AF`
  and PDF/A-3 + Factur-X XMP, and does not check the host. It succeeds on
  both hosts.
- No container check: CheckPDF is `unsupported`.
- License: only the plain entry points are called, not the `*WithLicense`
  variants, and no key is read.
- Network: no reference to `System.Net.Http`, `HttpClient` or sockets in the
  assemblies. Only `Antiphon.Xsd.Core` references `System.Net`, for the
  `Credentials` property of its `XmlResolver`, which reads the embedded XSDs.

### SaxonCS-HE 13.1.0 (Saxonica), reference

SaxonJ-HE transpiled to C#: neither Java nor IKVM. It is the offline XSLT 2.0+
reference outside Java, the equivalent of saxonche in the Python bench.

- Official CEN EN 16931 CII 1.3.17 XSLTs for all invoices. KoSIT XRechnung
  3.0.2 schematron 2.6.0 is added when BT-24 declares XRechnung, with the
  custom levels of `validator-configuration-xrechnung` v2026-08-31:
  CII-SR-465 and 466 become errors, CII-SR-475 and 476 become information.
- Result: 6 files out of 26 (BR-CL-25 ×2, CII-SR-465 ×4), identical to the
  protocol's reference.
- Artifacts downloaded by `run.sh setup`, SHA-256 checked, same archives as
  the Python bench. The stylesheets are compiled on the first call, which is
  included in `first_call_ns`.
- `InvariantGlobalization` breaks SaxonCS. The invariant culture is therefore
  set in the harness code.

### System.Xml.Schema (XSD floor)

The BCL's `XmlSchemaSet`, compiled once, and a validating `XmlReader`,
against the UN/CEFACT CII D16B XSDs of the KoSIT v2026-08-31 configuration
(`resources/cii/16b/xsd`). Not a Factur-X library: it shows the cost of the
schema layer that most libraries put in front of their rules. None of the 26
examples is flagged.

### PdfPig 0.1.16 (floor)

The most downloaded C# PDF reader on NuGet. The extraction is written in
`libs/pdfpig/Program.cs`: `PdfDocument.Open`, then
`Advanced.TryGetEmbeddedFiles` (`/Names/EmbeddedFiles` tree with `/Kids`,
decoded streams). It takes the first file that carries a standard invoice
name, with the same list of names as the Rust floor.

### PDFsharp 6.2.4 (floor)

The PDF library under ZUGFeRD.PDF-csharp (PDFsharp-extended fork), Balsoft,
International.EInvoicing, FacturXDotNet (fork) and Securibox (7.0 preview).
The `/Root/Names/EmbeddedFiles` walk with `/Kids` is hand-written in
`libs/pdfsharp/Program.cs`. 6.2.4 is the last MIT version: 6.5.0
(2026-10-14) and 7.0.0-preview-2 are under AGPL-3.0-only.

## Excluded libraries

| Package | Version (downloads) | Reason |
|---|---|---|
| ZUGFeRD.NET | 1.0.1 (22,533) | Konik Java library compiled by IKVM |
| Saxon-HE | 10.9.0 (2,203,102) | Saxon Java on IKVM, for .NET Framework |
| SaxonHE10Net31, SaxonHE12s9apiExtensions | 10.9.21 (101,120), 12.10.1 (103,595) | Saxon Java compiled by IKVM; SaxonCS-HE covers the need without Java |
| Lefty.Schematron | 0.0.16 (1,635) | Schematron tool that relies on Saxon Java (IKVM) |
| Eunormia.Sdk | 1.0.1 (489) | Bundles a JRE |
| OfficeIMO.Invoicing.Validation, business rules stage | 3.6.2 | Java and SaxonJ-HE 12.10; only the XSD stage is measured |
| FactoorSharp.FacturX (+PDF, Validator) | 21.0.0 (8,140) | Commercial; the license contacts a server |
| Facturexa.* | 7.6.1 (700) | Commercial: license key or 30-day trial |
| AloahaZUGFeRDToolkit | 11.0.0.16 (100,402) | Commercial, with no evaluation license usable offline |
| itext7.pdfinvoice | 2.0.1 (10,700) | Commercial (iText), license key; latest version in 2018 |
| Access.It.Pdf.Facturx | 10.12.1 (12,006) | No declared license, hence no right of use |
| FacturX-csharp | 1.0.6 (9,319) | Derived from ZUGFeRD-csharp, already measured |
| PROBIS-ZUGFeRD (+PDF, CDAR) | 19.0.4 (307) | Fork of ZUGFeRD-csharp, already measured |
| Lombiq.EInvoiceValidator | 2.0.1 (6,435) | Requires Node.js |
| InvoiceXml.FacturX, InvoiceXml.Zugferd.XRechnung | 1.5.0 (738, 769) | Clients of a remote API (excluded by the protocol) |
| StanzaApi.FacturXValidator, StanzaApi.PeppolValidator | 1.0.0 (141, 139) | Clients of a remote API |
| Xingen.Sdk, Fakturnaut.Sdk, EInvoice.Sdk, Falco.Api | 0.2.0 (370), 1.0.1 (214), 0.1.0 (178), 1.0.4 (1,186) | Clients of a remote API |
| Recommand.Client, EConnect.Psb | 0.4.2 (4,882), 0.9.35 (30,168) | Clients of remote Peppol APIs |
| SW-ZUGFeRD | 1.0.0 (430) | Empty package (net461) |
| csv2zugferd | 1.3.1 (218) | Command-line tool |
| PdfPinata.EInvoice | 0.6.1-alpha.26 (1,274) | Embedding only, into a PdfPinata document, in an alpha prerelease published on the day of the survey; no parsing or validation |
| Schematron, XSchematron, Xsl3, Bosak.Xslt, OutSmart.DAXon, XPath2 | 1.0.0 (28,247), 0.2.3 (942)… | Generic engines with no invoice artifacts (XPath2 is measured through Securibox) |
| RemitExtract | 0.1.1 (294) | Data extraction from PDF invoices and QR-bills, off topic |
| National packages (OioUbl, ZATCA, FatturaPA…) and UBL-only packages | — | Off topic for a Factur-X / XRechnung CII corpus |
| Antiphon.Invoicing.Rules.Peppol.Cii, Antiphon's UBL packages | 1.0.0 (130 to 210) | Antiphon components that the corpus does not use: no example declares Peppol, and the corpus is CII |

Java is used only by the excluded candidates and by OfficeIMO's rules stage.
SaxonCS-HE, used by pax.XRechnung.NET, En16931 and the reference, is C#.

## Accuracy surveys (`./run.sh survey`)

| Library | validate (flagged / 26) | parse (/ 26) | extract (/ 104) |
|---|---|---|---|
| zugferd-csharp | 3 (false positives on the totals: Teilrechnung, Kraftfahrversicherung, Sachversicherung) | 26 | 103 (with model parsing) |
| international-einvoicing | 2 (BR-CL-25 ×2) | 26 | 104 |
| balsoft | 5 (BR-CL-25 ×2; BT-106/109/112/115 discrepancies on 3 valid examples) | 26 | 104 |
| securibox | **6** (same files as the reference: CII-SR-465 and `codedb cl1` ×4 on the XRECHNUNG files, `codedb cl17` ×2) | 26 | 74 (XmpException ×29) |
| officeimo | 23 (model checks: INV-WHITESPACE, INV-PRICE…; XSD without errors) | 26 | 104 |
| facturxdotnet | 15 (parse errors: schemeID "FC", XRechnung BT-24) | 11 | 70 (factur-x.xml only) |
| pax-xrechnung | 6 (FX-SCH-A-000031 ×2; 4 XRECHNUNG rejected) | — | — |
| en16931 | 2 (FX-SCH-A-000571 ×2) | 18 | — |
| antiphon | 2 (BR-CL-25 ×2) | 26 | 104 |
| saxoncs (reference) | **6** (BR-CL-25 ×2, CII-SR-465 ×4) | — | — |
| xsd (floor) | 0 | — | — |
| pdfpig (floor) | — | — | 104 |
| pdfsharp (floor) | — | — | 104 |

Reference (PROTOCOL): 6 out of 26. SaxonCS-HE matches it exactly. Securibox
flags the same files, but as an artifact (see its section).

## Benchmark status (trial run `COUNT=1 BENCHTIME=0.05`)

`o` ok, `F` fail (message in the note of the JSONL line), `-` unsupported.
Columns in order:

1. Extract S/L/Min/Ext
2. Parse S/Li/La
3. ParseValidate S/Li/La
4. Write S/Li
5. Embed, EmbedClassic
6. CheckPDF S/L

| Library | Extract | Parse | ParseValidate | Write | Embed | CheckPDF |
|---|---|---|---|---|---|---|
| zugferd-csharp | o o o o | o o o | o o o | o o | o o | - - |
| international-einvoicing | o o o o | o o o | o o o | o o | o o | - - |
| balsoft | o o o o | o o o | o o o | o o | o o | - - |
| securibox | o o o F | o o o | o o o | o o | o o | o o |
| officeimo | o o o o | o o o | o o o | F F | F o | o o |
| facturxdotnet | o o o F | F F F | F F F | F F | F F | F o |
| pax-xrechnung | - - - - | - - - | o o - (skipped) | - - | - - | - - |
| en16931 | - - - - | o F o | o o o | o F | - - | - - |
| antiphon | o o o o | o o o | o o o | o o | o o | - - |
| saxoncs | - - - - | - - - | o o o | - - | - - | - - |
| xsd | - - - - | - - - | o o o | - - | - - | - - |
| pdfpig | o o o o | - - - | - - - | - - | - - | - - |
| pdfsharp | o o o o | - - - | - - - | - - | - - | - - |

The protocol's two hosts already carry an `/AF` pointing to a factur-x.xml
(Filespec dated 2022-02-18). OfficeIMO refused to add a second file with the
same name: the adapter replaces it (`Replace`). The other libraries rewrite or
complete the entry.

## Coverage table

| Library | Version | License | Latest version | Downloads | CII parsing | CII writing | UBL | PDF extraction | PDF embedding | XSD | EN 16931 (CII) | Factur-X profiles (MINIMUM, BASIC WL, EXTENDED) | BR-FR | BR-DE (XRechnung) | Peppol BIS 3 | PDF/A-3 container check | External dependencies |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| ZUGFeRD-csharp (+PDF) | 18.0.1 | Apache-2.0 | 18.0.1 (2026-10-09); end of maintenance on 2026-11-30 | 1,842,187; PDF 359,076 | yes (ZUGFeRD 1.0, 2.0, 2.3) | yes (all profiles) | yes (XRechnung UBL, parsing and writing) | yes (103/104, with model parsing) | yes (rewrite, XMP, OutputIntent, system fonts) | no | partial: totals recomputation (lines, BR-CO-14, BR-CO-17…) | parsing and writing of all profiles; no profile rules | no | no (XRechnung writing only) | no | no | PDFsharp-extended 1.3.0; system fonts at embedding |
| International.EInvoicing | 1.1.0 | MIT | 1.1.0 (2026-09-07) | 10,807; FacturX.PdfSharp 5,014 | yes | yes | yes | yes (104/104) | yes (associated file, XMP added or completed) | yes: CII D16B | yes: preprocessed CEN 1.3.16 `.sch` (806 assertions), in-house engine | not offline; FeRD rules must be loaded from a folder | no (`AddFrance`, not loaded) | partial: KoSIT "3.0" schematron taken from `master` (52 assertions), without custom levels | separate package, not loaded | no (XMP compared with the XML) | PDFsharp 6.2.4 (MIT) |
| Balsoft.Hive.EInvoice (+Pdf) | 1.1.0 | Apache-2.0 | 1.1.0 (2026-10-23) | 626; Pdf 563 | yes | yes | yes | yes (104/104) | yes (rewrite, OutputIntent, incremental XMP) | no | partial: rules coded in C# (115 BR identifiers), recomputed totals | partial: profiles checked in C# | no | partial: 26 BR-DE coded in C# | partial: 4 PEPPOL-EN16931 rules | no | PDFsharp 6.2.4; system fonts at embedding |
| Securibox.FacturX | 0.6.0-alpha-3 | MIT (license file) | 0.6.0-alpha-3 (2026-09-04) | 20,644 | yes (XmlSerializer, one model per profile) | partial (public serialization only through the PDF) | no | partial (74/104; validation included) | yes (rewrite, prior validation) | yes: Factur-X 1.09 XSD per profile | yes: through the FeRD 1.09 schematrons | yes: FeRD 1.09 schematrons for the 5 profiles (424 assertions in EN 16931) | no | no | no | partial: `IsFacturXValid` (XMP and XML), no PDF/A | PDFsharp 7.0.0-preview-1, XmpCore 6.1.10.1 (Adobe license), XPath2 1.1.5 (MS-PL) |
| OfficeIMO.Invoicing (+Validation, +Pdf) | 3.6.2 | MIT | 3.6.2 (2026-11-18) | 2,524; Validation 1,766; Pdf 213,774 | yes (bounded) | partial (refuses a model that fails its checks) | yes | yes (104/104) | partial (generic associated file without XMP; refused on `host.pdf`) | yes: XSD from the KoSIT v2026-08-31 archive | no (requires Java); model checks only | no (requires Java) | no | no (requires Java) | no (requires Java) | partial: `AssessCompliance(FacturX)` (XMP, associated file, OutputIntent, fonts), without the XML | KoSIT archive downloaded; Java and SaxonJ-HE for the rules (not used) |
| FacturXDotNet | 0.5.9-alpha | MIT | 0.5.9-alpha (2025-04-21) | 6,411 | partial (MINIMUM header, 11/26) | partial (MINIMUM) | no | partial (factur-x.xml only, 70/104) | yes (XMP generated from the XML) | no | partial (19 rules of the MINIMUM profile) | MINIMUM only | partial (BR-HYBRID-FR-01, left out here) | no | no | partial: 15 hybrid rules (XMP, file name and relationship), no PDF/A | FacturXDotNet.PDFsharp 6.2.0-preview-3 fork |
| pax.XRechnung.NET | 0.3.3 | MIT (NuGet); GPL-3.0 (repository) | 0.3.3 (2026-10-10) | 3,143 | no | no | yes (XRechnung UBL model) | no | no | yes: Factur-X 1.07.3 EN16931 XSD | yes: through the FeRD 1.07.3 EN16931 XSLT | EN 16931 only (BT-24 `urn:cen.eu:en16931:2017` required) | no | no for CII | no | no | SaxonCS-HE 13.0.0 (MPL-2.0), ICU4N |
| En16931 | 0.1.0 | EUPL-1.2 | 0.1.0 (2026-11-05) | 72 | partial (18/26) | partial | yes (through the intermediate representation) | no | no | XSDs shipped (CII D16B and D22B, UBL 2.1) | yes: official CEN 1.3.16 XSLT | partial: FeRD 1.09.2 BASIC and EN 16931 | no | yes: KoSIT 3.0.2 / 2.6.0 (and extension, CVD), without custom levels | no | no | SaxonCS-HE 13.1.0 (MPL-2.0) |
| Antiphon.Invoicing (+ Parsing.Cii, Generation.Cii, Validation.Pipeline, Rules.Cen.Cii, Rules.XRechnung.Cii, Xsd.Cii, Pdf) | 1.0.0 | commercial, closed-source; free evaluation without a key | 1.0.0 (2026-09-18) | 202; sub-packages 130 to 144 | yes (unmapped elements reported) | yes (Comfort, Basic, BasicWithoutLines, Minimum) | yes (separate packages, not measured) | yes (104/104) | yes (associated file, PDF/A-3 XMP; PDF/A host assumed) | yes: CII D16B | yes: precompiled CEN 1.3.16 Schematron (806 assertions), in-house evaluator | profile generation; no FeRD rules | no | yes: KoSIT schematron v2.6.0 (74 assertions), without custom levels | separate package (Rules.Peppol.Cii), not loaded | no | PDFsharp 6.2.4 (MIT), System.Security.Cryptography.Pkcs (license) |
| SaxonCS-HE (reference) | 13.1.0 | MPL-2.0 | 13.1.0 (2026-10-15) | 69,586 | no | no | no | no | no | no | yes: official CEN 1.3.17 XSLT | no | no | yes: KoSIT 2.6.0 and v2026-08-31 levels | no | no | official artifacts downloaded |
| System.Xml.Schema (floor) | BCL .NET 11.0.0 | MIT | — | — | no | no | no | no | no | yes: CII D16B (KoSIT archive) | no | no | no | no | no | no | none |
| PdfPig (floor) | 0.1.16 | Apache-2.0 | 0.1.16 (2026-08-22); 0.1.17-alpha prereleases | 47,810,678 | no | no | no | hand-written floor | no | no | no | no | no | no | no | no | none |
| PDFsharp (floor) | 6.2.4 | MIT (6.5.0 and later: AGPL-3.0-only) | 6.5.0 (2026-10-14, AGPL); 6.2.4 (2026-01-06), last MIT | 85,920,430 | no | no | no | hand-written floor | no | no | no | no | no | no | no | no | none |

## Build and measurement notes

- One project per library, restored (`dotnet restore --locked-mode`) and
  built (`dotnet build -c Release`) separately. That way a project that does
  not build does not stop the others, and each process loads only its own
  dependencies. Outputs go to `artifacts/` (`UseArtifactsOutput`), the
  official artifacts to `.artefacts/`.
- Run through the muxer (`dotnet artifacts/bin/<lib>/release/bench-<lib>.dll`)
  rather than the apphost, which does not find a .NET installed by Homebrew
  without `DOTNET_ROOT`.
- Runtime settings: the .NET defaults (tiered compilation, dynamic PGO,
  concurrent workstation GC). The only deviation: the invariant culture is
  set in the harness code (`harness/Harness.cs:743-746`).
  `InvariantGlobalization` would have the same effect on numbers, but breaks
  SaxonCS.
- The body of each program runs in a thread with a 256 MiB stack
  (`harness/Harness.cs:726`): the official stylesheets recurse deeply. A
  single measurement thread, no forced GC. Each result goes through
  `Sink.Consume`, which prevents dead-code elimination.
- CPU time (`cpu_ns_per_op`): `Environment.CpuUsage.TotalTime` (`getrusage`
  on macOS, all threads of the process), read just before and just after
  each series. Short measurement from the `bench` pass
  (`COUNT=2 BENCHTIME=0.5`): CPU stays within ±20% of wall-clock time.
  Deviations up to ×1.6 (Securibox, En16931, Balsoft…), due to background
  tiered compilation and to the GC. The only clear exception is
  ZUGFeRD-csharp on `Embed_Small` and `EmbedClassic_Small`: ×7 to ×8. On each
  call, about twenty thread-pool threads (`.NET TP Worker`) walk and open
  files in the system folders in parallel, looking for the fonts to embed
  (`open`, `lstat`, `fstat`, `flock`, `pread` in a `sample` trace). This work
  counts in the process's CPU, not in that of the measured thread.
- `./run.sh resources`: one process per library and per benchmark,
  `dotnet artifacts/bin/<lib>/release/bench-<lib>.dll resources <Bench>`. It
  is the `bench` command, launched directly under `/usr/bin/time -l` (peak
  RSS in bytes). The muxer loads the runtime into its own process (hostfxr,
  hostpolicy, coreclr). Homebrew's `dotnet` is a script that `exec`s the
  muxer: the measured PID is therefore the runtime's. No child process. Only
  two assemblies reference `ProcessStartInfo`: OfficeIMO.Invoicing.Validation,
  for the Java runner of its rules stage, which the adapter does not
  configure, and PdfSharp.Quality, for development helpers. Neither is
  called.
- `first_call_ns` includes JIT compilation and lazy compilations:
  stylesheets (SaxonCS-HE, pax, En16931), schemas (`XmlSchemaSet`,
  OfficeIMO) and rules (International.EInvoicing, Antiphon).
- Security review before building: MSBuild targets imported by the packages
  (`obj/*.nuget.g.targets`).
  - Microsoft compatibility targets;
    `Microsoft.Extensions.Configuration.Binder`; `ICU4N.Resources.targets`
    (a pax dependency).
  - The latter imports `CopyPatch.targets`, an inline RoslynCodeTaskFactory
    task. It copies the satellite resource DLLs and updates `.deps.json`
    with the Newtonsoft.Json it bundles. No network access; with
    `SatelliteResourceLanguages=en`, it copies nothing.
  - The `build/` folder of DotNet.Bundle only activates on an explicit
    `BundleApp` target.
  - No other package runs code at build time.
- The protocol's two hosts (`host.pdf`, `host_classic.pdf`) are PDF/A-3
  documents without an invoice: no attachment, no `/AF`, no Factur-X XMP
  (until 2026-11-22, they kept an `/AF` pointing to a factur-x.xml). The
  `Embed_*` check rereads the embedded XML and compares it with the CII
  profile, not with the exact content.

## Conformance cases (`./run.sh cases`)

The adapters validate the documents of [`../../cases`](../../cases/README.md) with the same call and rule set as `survey` and `bench`; a UBL document (root element read in the first 64 KiB) goes to the library's UBL rules when it has them, and a library without UBL support gets the UBL documents like any other. The CII path, the one the benchmarks time, is unchanged. Adapter changes made for the cases:

- **SaxonCS-HE + official XSLT**: CEN UBL 1.3.17 XSLT for a UBL document, KoSIT XRechnung UBL schematron when BT-24 declares XRechnung, all the custom levels of the KoSIT scenarios. Survey unchanged (6/26).
- **System.Xml.Schema floor**: the UBL 2.1 Invoice and CreditNote schemas for a UBL document; the XML-DSig schema they import has a DOCTYPE and is loaded first, with `DtdProcessing.Parse`. Survey unchanged (0/26).
