# Coverage: C++ ecosystem

C++ bench for the common protocol (`../PROTOCOL.md`). One executable per
library (`build/bench-<lib>`, source `libs/<lib>/main.cpp`), the C bench's
harness (`../c/harness/harness.c`, compiled as C) seen through its C++ facade
(`harness/harness.hpp`), CMake build in `Release` (CMake's default flags:
`-O3 -DNDEBUG`, no CPU-specific option), C++20.
Compiler: Apple clang 21.0.0, CMake 4.5.1, Ninja. The libraries come from
Homebrew (macOS arm64), except Kura (pinned source) and SaxonC-HE (pinned
official binary).

**The C++ ecosystem has no dedicated library that parses, validates or writes
CII.** Only one maintained library knows Factur-X: **Kura**
(bentopdf-kura, AGPL-3.0, ★25, version 1.4.1 released 2026-11-08), a PDF/A
engine that reads, attaches and checks the invoice **at the PDF container
level**, without ever parsing the XML. EN 16931 validation exists in C++ only
through **SaxonC-HE**, Saxonica's XSLT 3.0 processor compiled as a native
library, which runs the official CEN and KoSIT XSLTs: it is the bench's
reference validator, not an invoice library. Everything else is a
**floor**: qpdf, PoDoFo and poppler-cpp for PDF, Xerces-C (with the CII
XSD) and pugixml for XML, with a few dozen lines written in the
adapter.

The evidence cited is paths in the libraries' sources, at the tag of the
built version (GitHub repository, archive pinned in `.artefacts/`, or headers
installed by Homebrew in `/opt/homebrew/opt/<formula>/include`), and in this
bench's adapters. Popularity as of 2026-11-22: GitHub stars, date of the
measured version (tag), Homebrew installs over 365 days
(`formulae.brew.sh`, mostly as a dependency of other formulae). There is no
download registry for C++.

## Candidate search

- GitHub, search for C++ repositories: `zugferd`, `factur-x`, `facturx`,
  `xrechnung`, `en16931`, `e-rechnung`, `CrossIndustryInvoice`, `einvoice`,
  `e-invoice`. Code search in C++: `"urn:cen.eu:en16931:2017"`,
  `rsm:CrossIndustryInvoice`, `"factur-x.xml"`, `ZUGFeRD`. Results: Kura,
  the projects in the excluded table, applications (KDE Kraft,
  YACAppServer, UltraCanvas, Duncans-Suite, PDFE) and commercial SDK samples
  (Datalogics, PDF-XChange).
- vcpkg (2,917 ports) and Conan Center (1,982 recipes): no e-invoicing
  port; only `qpdf`, `podofo`, `poppler`, `xerces-c`, `pugixml`, `libxslt`
  (and `libmupdf` in vcpkg).
- Homebrew: no `zugferd`, `factur`, `xrechnung` or `einvoice` formula.
  The `saxon` formula is Saxon-HE for Java, not SaxonC. Debian: no source
  package with these names.

## Measured libraries

### Kura 1.4.1 (bentopdf-kura): dedicated, at the PDF container level

C++ engine for PDF/A, PDF/X and e-invoice conversion and checking, built on
qpdf. Pinned source: archive of the `v1.4.1` tag
(SHA-256 `18353ceb…43fb8` in `run.sh`), of which only `pdfa-engine/` is
unpacked; built with `add_subdirectory`, with the CLI, PDFium, zxing-cpp,
QuickJS, the fuzzers and the SDK test disabled (`CMakeLists.txt`), linked
against Homebrew's qpdf. It knows Factur-X, ZUGFeRD, XRechnung and Order-X at
the PDF level; it does not parse the CII (profile detected by text patterns,
`core/src/einvoice.cpp:92` `detectInvoice`).

- Extraction: `pdfa::readInvoice` (`core/src/einvoice_read.cpp:56`), first
  file in `/Names/EmbeddedFiles` under one of the names `factur-x.xml`,
  `zugferd-invoice.xml`, `ZUGFeRD-invoice.xml`, `xrechnung.xml`,
  `order-x.xml` (`:72-73`). The XMP is not consulted. 104 PDFs out of 104.
- Embedding: `kura_convert(host, "3b", &options)` with `invoice_xml`
  (`core/include/kura/kura.h:114`, `capi/kura_c.cpp:203`), the SDK entry
  point. It converts the whole host to PDF/A-3b (fonts, color spaces,
  metadata) and attaches the invoice (`attachInvoiceXml`,
  `core/src/sanitize.cpp:990`, called at `:1347`) with `/AF`,
  `/AFRelationship` and the `fx:` XMP extension schema; the profile and the
  file name come from `detectInvoice`. The whole PDF/A conversion is inside
  the measurement (11.7 ms on the xref-stream host and 12.4 ms on the
  classic-xref host, versus 0.4 to 0.5 ms for the floors, which convert
  nothing).
- Check of the output PDF (outside the bench, `FX_DUMP_EMBED`): a single
  `factur-x.xml`; `zerechnung validate` reports one error, PDF-18 (the
  attachment stream has no `/Params /ModDate`), and no warning.
- Quirk observed: Kura re-registers in `/Names/EmbeddedFiles` the files
  that `/AF` lists without naming them (`sanitize.cpp:825-870`), then
  `attachAssociatedFile` (`:934`) adds the new invoice without removing an
  invoice already present. Until 2026-11-22, the protocol hosts kept in
  `/AF` the `factur-x.xml` of the intarsys PDF they were derived from: the
  output PDF then contained two `factur-x.xml` files (PDF-28 for
  `zerechnung validate`), and `readInvoice` read back the first, old one.
  The fixed hosts no longer have an `/AF`.
- CheckPDF: `kura_read_invoice` (`kura.h:123`, `kura_c.cpp:268`), which calls
  `pdfa::checkInvoice` (`core/src/einvoice_check.cpp:8-95`): specification
  URN, presence of the identifier, type code and date, more than one
  attached XML, MINIMUM / BASIC WL warning, attachment name,
  `/AFRelationship`, presence in `/AF`, and in the XMP `DocumentFileName`,
  `ConformanceLevel`, `DocumentType`. These checks rely on `std::regex`
  (`einvoice_check.cpp:31`), hence about 1.4 ms per call. Not a PDF/A
  validation: an inconsistency it finds goes into the report, not into a
  failure.
- No XML parsing, writing or validation.
- Network: none.

### SaxonC-HE 13.1.1 (Saxonica): reference validator

Saxon-HE compiled as a native library by GraalVM `native-image` (no JVM),
with a C++ API (`saxonc/SaxonProcessor.h`, `Xslt30Processor.h`). Official
binary for macOS arm64 (`SaxonCHE-macos-arm64-13-1-1.zip`, files dated
2026-10-26). Saxonica does not publish a checksum: the one in `run.sh`
(`8f164115…8630281`) was recorded at the first download. It is the same
library as `saxonche`, the Python binding used by
`../python/adapters/saxonche_schematron.py`, with the same artifacts and the
same counting.

- ParseValidate (`libs/saxonc-he/main.cpp:155-187`): stylesheets compiled
  once on the first call (`compile`, `:59-87`, cost included in the
  `first_call_ns` of `ParseValidate_Small`, about 140 ms), then on each call
  `parseXmlFromString` and `transformToString`:
  - the official CEN XSLT `EN16931-CII-validation.xslt` from
    `validation-1.3.17` (ConnectingEurope/eInvoicing-EN16931, XSLT 2.0, 810
    rules: CII-SR 484, CII-DT 103, BR 58, BR-CL 23, BR-CO 23, BR-DEC 21,
    BR-IC 12, BR-O 14, BR-S, BR-Z, BR-E, BR-AE, BR-IG, BR-IP, BR-G 10 each,
    BR-B 2; 339 flagged `fatal`, 471 `warning`) for every invoice;
  - in addition, when BT-24 contains `xrechnung` (`is_xrechnung`, `:90-110`,
    same test as the Python adapter), the KoSIT stylesheet
    `XRechnung-CII-validation.xsl` from `xrechnung-schematron` 2.6.0
    (XRechnung 3.0.2, 74 rules: 31 BR-DE, 7 BR-DEX, 7 BR-DE-CVD, 22
    PEPPOL-EN16931-R, 7 transitional), as in the CII scenario of the
    KoSIT configuration.
- Counting (`errors`, `:119-149`): every `failed-assert` or
  `successful-report` in the SVRL flagged `fatal` or `error`, after the levels
  specific to the KoSIT configuration of 2026-08-31 in the XRechnung scenario
  (CII-SR-465 and CII-SR-466 raised to `error`, CII-SR-475 and CII-SR-476
  lowered to `information`), applied to both reports as the Python adapter
  does. The SVRL is read with pugixml inside the measured call.
- Result: 6 invoices flagged out of 26 (BR-CL-25 ×2, CII-SR-465 ×4), the
  protocol reference. On a file corrupted by hand (wrong total amount),
  BR-CO-15 and BR-CO-16 are flagged as expected.
- No XSD (absent from the CEN CII artifacts); no French rules; no UBL or
  Peppol BIS measured (the CEN UBL XSLTs are not loaded).
- Network: none (local stylesheets; `setcwd` to `.artefacts`).

### qpdf 12.4.2 (PDF floor)

The "content-preserving" PDF library of the C++ world, with no knowledge
of Factur-X.

- Extraction (`libs/qpdf/main.cpp:37-53`): `processMemoryFile`,
  `QPDFEmbeddedFileDocumentHelper::getEmbeddedFiles()` (the whole
  `/Names/EmbeddedFiles` tree), first file under a standard invoice name,
  stream decoded by `getStreamData(qpdf_dl_generalized)`. 104 PDFs out of 104.
- Embedding (`:55-100`): `QPDFEFStreamObjectHelper::createEFStream`,
  `QPDFFileSpecObjectHelper::createFileSpec`, `replaceEmbeddedFile`; by
  hand: `/AFRelationship` (`Alternative`, or `Data` in MINIMUM and BASIC WL,
  `:69`), `/AF` rebuilt without any invoice file already present
  (`:75-86`), new `/Metadata` with the harness's Factur-X XMP (`:91-93`);
  full write by `QPDFWriter` in memory with
  `setNewlineBeforeEndstream(true)` (`:97`, required by PDF/A). No PDF/A
  conversion.
- Check of the output PDF (outside the bench, `FX_DUMP_EMBED`): `zerechnung
  validate` reports no error (warnings: no PDF/A output intent, `/Info`
  without an XMP equivalent), Kura's `checkInvoice` finds it consistent,
  `qpdf --check` reports no error.

### PoDoFo 1.1.3 (PDF floor)

General-purpose C++17 PDF library, with no knowledge of Factur-X.

- Extraction (`libs/podofo/main.cpp:24-45`):
  `PdfMemDocument::LoadFromBuffer`, `PdfEmbeddedFiles` tree walked with
  `ToDictionary` (including `/Kids`), first file under a standard invoice
  name, `PdfFileSpec::GetEmbeddedData`. 104 PDFs out of 104.
- Embedding (`:47-106`): `CreateFileSpec` and `SetEmbeddedData`, added to
  the `PdfEmbeddedFiles` tree (`:71`); by hand: `/AFRelationship` per profile
  (`:56`), `/AF` rebuilt without any invoice file already present
  (`:75-92`), new Factur-X XMP `/Metadata` (`:94-99`);
  `PdfMemDocument::Save` with `NoMetadataUpdate` (`:104`), so that PoDoFo
  keeps this XMP and does not write the current time into it. No PDF/A
  conversion. Same checks outside the bench as for qpdf, same results. The
  host's `/AF` is read with `FindKey` (`:77`): `FindKeyAs<PdfArray *>`, used
  until 2026-11-22, throws `ObjectNotFound` when the key is missing, which
  made embedding fail on the fixed hosts, which have no `/AF`.

### poppler-cpp 26.11.0 (extraction floor)

C++ facade of the Poppler rendering library. Extraction only
(`libs/poppler-cpp/main.cpp:22-38`): `document::load_from_raw_data`,
`embedded_files()`, first file under a standard invoice name,
`embedded_file::data()`. 104 PDFs out of 104. poppler-cpp has no writing
API: Embed is `unsupported`. It is the slowest of the extractors
(0.4 to 0.7 ms versus 15 to 45 µs): opening a Poppler document prepares
rendering.

### Xerces-C++ 3.3.0 (validating XML floor)

Apache's validating XML parser, with no knowledge of Factur-X. Both
parsers live for the whole process and are reused, as in a long-running
service.

- Parse (`libs/xerces-c/main.cpp:90-110`): `XercesDOMParser` with
  namespaces, without validation (`Val_Never`), on a `MemBufInputSource`; the
  DOM is adopted and then released on each call.
- ParseValidate (`:112-136` and `:164-175`): a second parser with XML Schema validation
  (`setDoSchema`, `Val_Always`) against the UN/CEFACT CII D16B XSD
  (`CrossIndustryInvoice_100pD16B.xsd` and its 3 imported schemas) from the
  KoSIT XRechnung 3.0.2 configuration of 2026-08-31, loaded once by
  `loadGrammar` and cached (`useCachedGrammarInParse`, `:133`); errors
  counted by an `ErrorHandler`. Schema only: no business rules.
  The 26 valid examples pass; an unknown element (`ram:Bogus`) added by
  hand gives 2 schema errors.
- Write (`:138-162`): `DOMLSSerializer` from the DOM to a
  `MemBufFormatTarget` in UTF-8.
- Network: none (`setLoadExternalDTD(false)`, local schemas).

### pugixml 1.16 (XML floor)

The fastest in-memory DOM in the C++ world, with no validation and no
knowledge of Factur-X. Parse: `xml_document::load_buffer` with
`parse_default` (`libs/pugixml/main.cpp:18`). Write: `xml_document::save`
with `format_default` into a string (`:33`). ParseValidate is
`unsupported`.

## Excluded libraries

| Project | Stars, activity, license | Reason |
|---|---|---|
| LASTRADA-Software/XRechnung | ★3, last commit 2024-10-16, MIT | XRechnung model in C++20 that only writes UBL (`XRechnungSerializer.cpp:2101`, `ubl:Invoice`; no `ram:` element); no parsing; `XRechnungValidator.hpp` is only an abstract interface. No activity for two years, none of the protocol's CII operations |
| dragotin/kraft (KDE Kraft) | ★80, v2.1.0 released 2026-10-22, GPL-2.0 | Quotes and invoicing application (Qt, KDE): `src/exportxrechnung.cpp` fills a UBL XRechnung template (`reports/xrechnung.xrtmpl`). Not a library, no CII |
| classix/zugferd-cpp | ★0, 2025-01-31, no license | Training project (IHK): a command-line tool that prints the embedded XML, with a copy of PoDoFo in the repository. No license, not a library; the operation is measured with PoDoFo |
| miyako/4d-plugin-factur-x, -PoDoFo, -QPDF | ★2, ★0, ★0, 2026-11-10, MIT | Plugins for the 4D database (`4DPluginAPI`, `C_BLOB`) that assemble a Factur-X with PoDoFo or qpdf. Not usable outside 4D; the PDF operations are measured directly with PoDoFo and qpdf |
| Neox365/Duncans-Suite (`docs/FacturX.h`) | ★2, 2026-11-21, no license | Module of a Win32 application that writes an EN 16931 CII from its own invoice structure and inserts it into the PDF it generates; its header states that it has not been validated with the KoSIT validator. No license, not a library |
| jw23578/YACAppServer (`xrechnunggenerator.cpp`) | ★0, 2025-11-19, no license | XRechnung CII generator (pugixml) internal to an application server. No license, not a library |
| ULTRA-OS-Project/UltraCanvas (UltraFIBU) | ★4, 2026-11-22, MIT | Invoice PDF report from a demo application; its header explicitly says it does not do ZUGFeRD |
| ipswitch8/PDFE-public | ★0, 2026-11-20, AGPL-3.0 | PDF editor for Windows (Qt6); ZUGFeRD only appears in its PDF 2.0 feature tests. An application, not a library |
| aspose-pdf-foss/Aspose.PDF-FOSS-for-Cpp | ★2, 2026-10-22, MIT | Recent general-purpose PDF library (C++20); `ZUGFeRD` is only a value of the `PdfFormat` enumeration (`include/aspose/pdf/pdf_format.hpp:33`) carried over from Aspose.PDF, with no Factur-X code in the repository. Neither a Homebrew formula nor a port: just another floor |
| ksef_fa3 | ★0, 2026-04-08 | Polish KSeF FA(3) e-invoice, not EN 16931 |
| FriNagy/flex2pdf | ★0, 2026-02-26, no license | Windows command-line tool on top of the commercial VPE SDK |
| PDFium | Google, BSD-3-Clause | Attachment API (`FPDFDoc_GetAttachment`, `FPDFDoc_AddAttachment`), but neither a Homebrew formula nor a packaged stable release (Kura's PDFium option is also disabled here). The PDF floors are covered by qpdf, PoDoFo and poppler-cpp |
| MuPDF | ★3,066, AGPL-3.0 | C API: measured in the C bench (`../c/COVERAGE.md`) |
| Commercial SDKs (Datalogics APDFL, PDFlib, Apryse, PDF-XChange, VPE) | proprietary | ZUGFeRD and Factur-X samples published by the vendors (`ConvertToZUGFeRD.cpp`, `ConvertToFactur-X.cpp` from Datalogics), not verified; paid licenses, outside the protocol |

## Accuracy surveys (`./run.sh survey`, trial run on 2026-11-22)

| Library | validate (flagged / 26) | parse (/ 26) | extract (/ 104) |
|---|---|---|---|
| kura | — | — | 104 |
| qpdf | — | — | 104 |
| podofo | — | — | 104 |
| poppler-cpp | — | — | 104 |
| xerces-c | 0 (XSD only) | 26 | — |
| pugixml | — | 26 | — |
| saxonc-he | 6: BR-CL-25 ×2, CII-SR-465 ×4 | — | — |

Reference (PROTOCOL): 6 out of 26. SaxonC-HE matches it exactly, with the
same rules as the Python bench. The XSD alone catches none of the six: they
are Schematron rules.

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
| kura | o o o o | - - - | - - - | - - | o o (PDF/A-3b conversion) | o o |
| qpdf | o o o o | - - - | - - - | - - | o o (floor) | - - |
| podofo | o o o o | - - - | - - - | - - | o o (floor) | - - |
| poppler-cpp | o o o o | - - - | - - - | - - | - - | - - |
| xerces-c | - - - - | o o o | o o o (XSD) | o o | - - | - - |
| pugixml | - - - - | o o o | - - - | o o | - - | - - |
| saxonc-he | - - - - | - - - | o o o | - - | - - | - - |

## Coverage table

| Library | Version | License | Latest version | Downloads | CII parsing | CII writing | UBL | PDF extraction | PDF embedding | XSD | EN 16931 (CII) | Factur-X profiles (MINIMUM, BASIC WL, EXTENDED) | BR-FR | BR-DE (XRechnung) | Peppol BIS 3 | PDF/A-3 container check | External dependencies |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Kura (bentopdf-kura) | 1.4.1 | AGPL-3.0 or commercial | 1.4.1 (2026-11-08) | ★25; no package | no (profile detected by patterns) | no | no | yes: `readInvoice`, 104/104 | yes: `kura_convert` level 3b, full conversion to PDF/A-3b (does not remove an invoice already listed in `/AF`) | no | no | profile inferred from the URN for the XMP; MINIMUM / BASIC WL warning | no | no | no | partial: consistency of payload, name, `/AFRelationship`, `/AF`, XMP (`checkInvoice`); not a PDF/A validation | qpdf, freetype, jpeg-turbo, openjpeg, little-cms2, zlib |
| SaxonC-HE (reference) | 13.1.1 | MPL-2.0 (CEN artifacts: EUPL-1.2; KoSIT: Apache-2.0) | 13.1.1 (2026-10-26) | Saxonica binary, no package | no | no | no (UBL XSLTs not loaded) | no | no | no | yes: official CEN 1.3.17 XSLTs, 810 rules | no | no | yes: KoSIT schematron 2.6.0 (XRechnung 3.0.2) if BT-24 is XRechnung | no (only the PEPPOL-EN16931-R rules of the KoSIT stylesheet) | no | none (embedded native GraalVM); macOS arm64 pinned |
| qpdf (floor) | 12.4.2 | Apache-2.0 | 12.4.2 (2026-09-26) | ★5,606; Homebrew 85,709 / year | no | no | no | floor, 104/104 | hand-written floor | no | no | no | no | no | no | no | jpeg-turbo, openssl |
| PoDoFo (floor) | 1.1.3 | LGPL-2.0-or-later or MPL-2.0 | 1.1.3 (2026-11-16) | ★616; Homebrew 918 / year | no | no | no | floor, 104/104 | hand-written floor | no | no | no | no | no | no | no | freetype, fontconfig, jpeg-turbo, libpng, libtiff, openssl |
| poppler-cpp (floor) | 26.11.0 | GPL-2.0-only or GPL-3.0-only | 26.11.0 (November 2026) | freedesktop.org repository (not on GitHub); Homebrew 800,191 / year | no | no | no | floor, 104/104 | no (no writing API) | no | no | no | no | no | no | no | cairo, glib, freetype, fontconfig, nss, gpgme… |
| Xerces-C++ (floor) | 3.3.0 | Apache-2.0 | 3.3.0 (2024-10-14) | ★161; Homebrew 40,081 / year | DOM only | DOM serialization | no | no | no | yes: CII D16B (KoSIT configuration 2026-08-31) | no | no | no | no | no | no | none |
| pugixml (floor) | 1.16 | MIT | 1.16 (2026-06-15) | ★4,674; Homebrew 34,961 / year | DOM only | DOM serialization | no | no | no | no | no | no | no | no | no | no | none |

## Build and measurement notes

- `./run.sh setup` downloads the same three official archives as the C
  bench, verified by SHA-256 (KoSIT XRechnung 3.0.2 configuration of
  2026-08-31, CEN `validation-1.3.17` CII, KoSIT `xrechnung-schematron`
  2.6.0), plus SaxonC-HE 13.1.1 for macOS arm64 (elsewhere, `saxonc-he` is
  skipped) and the archive of the Kura tag `v1.4.1`, of which only
  `pdfa-engine/` is extracted. It then builds each target separately; a
  missing library is skipped at configure time.
- SaxonC-HE is imported by hand in `CMakeLists.txt` (libraries and headers
  from the archive), without the CMake package files it contains: the
  configuration runs none of its scripts.
- Kura is built with `add_subdirectory` from its pinned source: the CMake
  configuration therefore runs its `CMakeLists.txt`. The dependencies
  (qpdf, freetype, jpeg-turbo, openjpeg, little-cms2, zlib) are checked
  beforehand, because Kura's `CMakeLists.txt` stops the whole configuration if
  one is missing; its optional parts (CLI, PDFium, zxing-cpp, QuickJS,
  fuzzers, SDK test) are disabled.
- The harness is the C bench's (`../c/harness/harness.c`), with the same
  protocol: an unmeasured check call, N in steps of 1, 2, 5, 10…,
  `COUNT` series, a single thread (512 MiB stack), freeing inside the
  measured loop, a `volatile` sink. `CLOCK_UPTIME_RAW` clock on macOS (the
  one behind `Instant` in Rust; `CLOCK_MONOTONIC` only advances in
  microseconds there). A C++ exception becomes a `fail` with its message; a
  fatal signal (SIGSEGV, SIGABRT…) writes the `fail` line prepared before the
  call.
- CPU time and peak memory, through the same harness: `getrusage(RUSAGE_SELF)`,
  user + system time of all threads, read before and after each
  series (`cpu_ns_per_op`); `ru_maxrss` (bytes on macOS) read at the end of
  the series in the `resources` pass, where `run.sh` runs
  `build/bench-<lib> resources <Bench>` directly, one process per benchmark.
  No library starts a child process. SaxonC-HE keeps its GraalVM isolate
  in the process, with its own threads (five threads in total during
  validation): their work and the GraalVM heap are counted; CPU time
  stays at about 100% of wall-clock time. Kura only starts threads to
  convert the colors of images of at least 2 Mpx (`colorx.cpp`,
  `colorconv.cpp`), which the host PDFs do not require.
- The embedding check reads the output PDF back with the same library's
  extraction and requires a CII XML, as in the Rust harness. With
  `FX_DUMP_EMBED=<dir>`, the PDF from the check call is kept; this is
  how the qpdf, PoDoFo and Kura PDFs were run through `zerechnung
  validate`, Kura's `checkInvoice` and `qpdf --check`.
- The embedding floors (qpdf, PoDoFo, and MuPDF in the C bench) write
  `/AFRelationship /Alternative`, or `/Data` for MINIMUM and BASIC WL
  (Factur-X 1.0.07, § 6.2.2), remove from `/AF` any invoice file already
  listed (the fixed hosts have no `/AF`) and replace the host's XMP with a
  Factur-X packet written by the harness. None of them converts to PDF/A
  (no output intent); Kura is the only one that does, hence the gap of a
  factor of 25 to 32 on Embed.
- SaxonC-HE starts a GraalVM isolate: the first call (about 140 ms)
  includes this startup and the compilation of both stylesheets. PoDoFo
  once showed a first call of 141 ms (one-time initialization), not
  reproduced afterwards.

## Conformance cases (`./run.sh cases`)

The adapters validate the documents of [`../../cases`](../../cases/README.md) with the same call and rule set as `survey` and `bench`; a UBL document (root element read in the first 64 KiB) goes to the library's UBL rules when it has them, and a library without UBL support gets the UBL documents like any other. The CII path, the one the benchmarks time, is unchanged. Adapter changes made for the cases:

- **SaxonC-HE + official XSLT**: CEN UBL 1.3.17 XSLT for a UBL document, KoSIT XRechnung UBL schematron when BT-24 declares XRechnung, all the custom levels of the KoSIT scenarios. Survey unchanged (6/26).
- **Xerces-C floor**: the UBL 2.1 grammars, loaded on the first UBL document. Survey unchanged (0/26).
