# Coverage: C ecosystem

C bench for the common protocol (`../PROTOCOL.md`). One executable per library
(`build/bench-<lib>`, source `libs/<lib>/main.c`), a shared harness
(`harness/harness.c`, reused by the C++ bench), CMake build in `Release`
(CMake's default flags: `-O3 -DNDEBUG`, no CPU-specific option). Compiler:
Apple clang 21.0.0, CMake 4.5.1, Ninja. The libraries come from Homebrew
(macOS arm64).

**There is no C library dedicated to Factur-X, ZUGFeRD, XRechnung or
EN 16931.** This bench therefore mostly measures **floors**: general-purpose
PDF or XML libraries, with a few dozen lines written in the adapter. The only
ZUGFeRD-specific function in the C ecosystem is MuPDF's `pdf_zugferd_xml()`,
which can read the attachment of a hybrid PDF and nothing more. No C library
validates the EN 16931 rules: libxslt, the XSLT engine of the C world, only
does XSLT 1.0 and cannot run the official XSLTs (XSLT 2.0). The reference
validation with the official artifacts is in the C++ bench (SaxonC-HE).

The evidence cited is paths in the libraries' sources, at the tag of the
version built (GitHub repository or headers installed by Homebrew in
`/opt/homebrew/opt/<formula>/include`), and in this bench's adapters.
Popularity as of 2026-11-22: GitHub stars, date of the measured version
(tag), Homebrew installs over 365 days (`formulae.brew.sh`, mostly as a
dependency of other formulae). There is no download registry for C.

## Candidate search

- GitHub, repository search in the C language: `zugferd`, `factur-x`,
  `facturx`, `xrechnung`, `en16931`, `e-rechnung`, `CrossIndustryInvoice`,
  `einvoice`, `e-invoice`. A single relevant result,
  `miyako/4d-plugin-factur-x` (4D plugin, see the excluded libraries).
- GitHub, code search in C: `"urn:cen.eu:en16931:2017"`,
  `rsm:CrossIndustryInvoice` and `"factur-x.xml"` only find MuPDF
  (`source/pdf/pdf-zugferd.c`) and its embedded copies (SumatraPDF and
  others), plus libharu's demo `demo/pdf_a_conformance.c` and its copies.
- vcpkg (2,917 ports) and Conan Center (1,982 recipes): no e-invoicing port;
  only the general-purpose libraries (`libmupdf`, `qpdf`, `libxslt`,
  `poppler`…).
- Homebrew: no `zugferd`, `factur`, `xrechnung` or `einvoice` formula (the
  `invoice` formula is a PDF invoice generator written in Go, unrelated).
  Debian: no source package by those names (search on
  `sources.debian.org`).

## Measured libraries

### MuPDF 1.28.6 (Artifex): ZUGFeRD reading by the library, embedding as a floor

A general-purpose PDF library (rendering, editing). It contains a small
ZUGFeRD module, `source/pdf/pdf-zugferd.c` (277 lines at tag `1.28.6`),
declared in `include/mupdf/pdf/zugferd.h:28-47`.

- Extraction (library function): `pdf_zugferd_xml()`
  (`pdf-zugferd.c:216-257`). It reads the XMP, looks for an
  `rdf:Description` that carries the attribute `xmlns:zf` =
  `urn:ferd:pdfa:CrossIndustryDocument:invoice:1p0#` (ZUGFeRD 1.0,
  `pdf-zugferd.c:74`), `xmlns:fx` = `urn:zugferd:pdfa:CrossIndustryDocument:invoice:2p0#`
  (ZUGFeRD 2.x, `:113`) or `xmlns:fx` =
  `urn:factur-x:pdfa:CrossIndustryDocument:invoice:1p0#` (Factur-X, `:155`),
  takes the name declared in `DocumentFileName`, then loads that file from
  the document's `/AF` associated files (`:232-247`). No XMP declaration, or
  a file missing from `/AF`: no invoice.
- The lookup is on the prefix, not on the namespace: the Mustang PDF
  `MustangGnuaccountingBeispielRE-20190610_507.pdf` binds the ZUGFeRD 2.0
  namespace to the `zf:` prefix (`xmlns:zf="urn:zugferd:pdfa:CrossIndustryDocument:invoice:2p0#"`),
  and MuPDF does not recognize it: 103 PDFs out of 104.
- `pdf_zugferd_profile()` (`:211`) only reads the profile declared in the
  XMP (ZUGFeRD 1.0 COMFORT, BASIC, EXTENDED; 2.x MINIMUM, BASIC WL;
  XRECHNUNG). No check: CheckPDF is `unsupported`.
- Embedding: no library function. A floor written in
  `libs/mupdf/main.c:94-183` with MuPDF's primitives:
  `pdf_add_embedded_file` (`:123`), sorted insertion into
  `/Names/EmbeddedFiles`, `/AF` (an invoice file already listed is
  replaced), `/AFRelationship` `Alternative` (or `Data` for MINIMUM and
  BASIC WL), `/Desc`, a new `/Metadata` stream with a Factur-X XMP written by
  the harness (`fx_facturx_xmp`), full rewrite by `pdf_write_document`
  (default options, `:164`). No PDF/A conversion. Read back by
  `pdf_zugferd_xml()`, which requires the XMP and `/AF`: the read-back
  therefore proves the container.
- Check of the resulting PDF (outside the bench, `FX_DUMP_EMBED`):
  `zerechnung validate` finds no errors (warnings: no PDF/A output intent,
  `/Info` without an XMP equivalent); Kura's `checkInvoice` deems it
  consistent; `qpdf --check` reports no errors.
- No parsing, writing or validation of the XML.
- Network: none.

### qpdf 12.4.2, C API (extraction floor)

qpdf is written in C++; it publishes a C API (`qpdf/qpdf-c.h`,
`qpdf/qpdf-c_objects.h`). An extraction floor written in
`libs/qpdf-c/main.c`: `qpdf_read_memory` (`:54`), walk of
`/Names/EmbeddedFiles` with `/Kids` (`find_in_tree`, `:18`, depth capped at
32), first file under a standard invoice name (`/UF` then `/F`), stream
decoded by `qpdf_oh_get_stream_data(…, qpdf_dl_generalized)` (`:81`). 104
PDFs out of 104. The C API does not expose the attachment helpers
(`QPDFEmbeddedFileDocumentHelper`): embedding is measured with the C++ API
in the C++ bench.

### libxml2 2.15.4 (XML floor)

The XML library of the C world (and of lxml). No knowledge of Factur-X.

- Parse: DOM through `xmlCtxtReadMemory` with `XML_PARSE_NONET`
  (`libs/libxml2/main.c:53-70`), one parser context per call.
- Write: `xmlDocDumpMemory` of the DOM (`:88-100`).
- ParseValidate: the same parsing, then `xmlSchemaValidateDoc` (`:146`)
  against the UN/CEFACT CII D16B XSD (`CrossIndustryInvoice_100pD16B.xsd`
  and its 3 imported schemas) from the KoSIT XRechnung 3.0.2 configuration of
  2026-08-31, an archive pinned by SHA-256 in `run.sh`. The schema is
  compiled on the first call (`compiled`, `:107-131`; cost in
  `first_call_ns` of `ParseValidate_Small`). The schema only: no EN 16931
  business rules. The 26 valid samples pass; a hand-corrupted file
  (`ram:Bogus` element) is correctly flagged.
- Network: none (`XML_PARSE_NONET`, local schemas).

### libxslt 1.1.45 (verified: cannot run the official XSLTs)

libxslt implements only XSLT 1.0 (with EXSLT). The official XSLTs are
XSLT 2.0. The `bench-libxslt` executable checks this at every launch:
`xsltParseStylesheetFile` (`libs/libxslt/main.c:73`) on
`cen-1.3.17/xslt/EN16931-CII-validation.xslt` (version `2.0`, 810 rules)
returns 36 compilation errors, the first being
`xsl:when : could not compile test expression 'xs:decimal(ram:LineTotalAmount) = …'`;
on the KoSIT stylesheet `XRechnung-CII-validation.xsl` (schematron 2.6.0),
170 errors. ParseValidate is therefore `unsupported`, with these numbers in
the message; the other operations are `unsupported` ("XSLT engine
only").

## Excluded libraries

| Project | Stars, activity, license | Reason |
|---|---|---|
| miyako/4d-plugin-factur-x | ★2, 2026-06-18, MIT | Plugin for the 4D database (`4DPluginAPI` API, `C_BLOB` and `C_TEXT` types): not a library usable outside 4D. The C++ variants (PoDoFo, QPDF) are in the C++ candidate search |
| libharu 2.4.6 | ★2,005, 2026-03-26, Zlib | Writes new PDFs only, neither reads nor modifies an existing PDF: no extraction and no embedding into a host. Its demo `demo/pdf_a_conformance.c` attaches a `factur-x.xml` to a document it creates |
| PDFio | ★362, 2026-11-17, Apache-2.0 | General-purpose PDF library, no Homebrew formula and no vcpkg/Conan port, no attachment helper: one more floor, like qpdf-c |
| Ghostscript (`lib/zugferd.ps`) | ghostpdl ★304, AGPL-3.0 | PostScript program for `gs` that produces a ZUGFeRD PDF/A-3 by re-distilling a file; not an in-memory API, `gs` not installed |
| PDFium | Google, BSD-3-Clause | Attachment extraction API (`FPDFDoc_GetAttachment`), but neither a Homebrew formula nor a port of a stable release; the extraction floors are covered by qpdf-c and MuPDF |
| Duncans-Suite (`docs/FacturX.h`) | ★2, 2026-11-21, no license | Internal module of a Win32 application (C++), see the C++ candidate search |
| Commercial SDKs (PDFlib, Datalogics APDFL, Apryse, PDF-XChange, VPE) | proprietary | ZUGFeRD / Factur-X features advertised by the vendors, not verified; paid licenses, outside the protocol |

## Accuracy surveys (`./run.sh survey`, trial run on 2026-11-22)

| Library | validate (flagged / 26) | parse (/ 26) | extract (/ 104) |
|---|---|---|---|
| mupdf | — | — | 103 (`zf:` prefix on the ZUGFeRD 2.0 namespace of a Mustang PDF) |
| qpdf-c | — | — | 104 |
| libxml2 | 0 (XSD only) | 26 | — |
| libxslt | — (`unsupported`) | — | — |

Reference (PROTOCOL): 6 out of 26. The XSD alone catches none of the six:
they are Schematron rules (BR-CL-25, CII-SR-465).

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
| mupdf | o o o o | - - - | - - - | - - | o o (floor) | - - |
| qpdf-c | o o o o | - - - | - - - | - - | - - | - - |
| libxml2 | - - - - | o o o | o o o (XSD) | o o | - - | - - |
| libxslt | - - - - | - - - | - - - | - - | - - | - - |

## Coverage table

| Library | Version | License | Latest version | Downloads | CII parsing | CII writing | UBL | PDF extraction | PDF embedding | XSD | EN 16931 (CII) | Factur-X profiles (MINIMUM, BASIC WL, EXTENDED) | BR-FR | BR-DE (XRechnung) | Peppol BIS 3 | PDF/A-3 container check | External dependencies |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| MuPDF | 1.28.6 | AGPL-3.0-or-later or commercial (Artifex) | 1.28.6 (2026-10-29) | ★3,066; Homebrew 16,036 / year | no | no | no | yes: `pdf_zugferd_xml` (XMP fx:/zf: then `/AF`), 103/104 | hand-written floor (no library function) | no | no | reads the profile declared in the XMP only | no | no | no | no | none at call time (libmupdf; freetype, harfbuzz, openjpeg, jbig2dec, gumbo… at build time) |
| qpdf-c (floor) | 12.4.2 | Apache-2.0 | 12.4.2 (2026-09-26) | ★5,606; Homebrew 85,709 / year | no | no | no | hand-written floor, 104/104 | no (see the C++ bench) | no | no | no | no | no | no | no | libqpdf (C++), jpeg-turbo, openssl |
| libxml2 (floor) | 2.15.4 | MIT | 2.15.4 (2026-09-04) | ★771 (GNOME mirror); Homebrew 562,672 / year | DOM only | DOM serialization | no | no | no | yes: CII D16B (KoSIT configuration 2026-08-31) | no | no | no | no | no | no | none |
| libxslt | 1.1.45 | X11 | 1.1.45 (2025-11-30) | ★75 (GNOME mirror); Homebrew 101,579 / year | no | no | no | no | no | no | no: XSLT 1.0, the official XSLTs (2.0) do not compile | no | no | no | no | no | libxml2, libgcrypt |

## Build and measurement notes

- `./run.sh setup` downloads three official archives verified by SHA-256
  (KoSIT XRechnung 3.0.2 configuration of 2026-08-31, CEN
  `validation-1.3.17` CII, KoSIT `xrechnung-schematron` 2.6.0; the same as
  the Python bench), unpacks them into `.artefacts/`, then builds each target
  separately. A missing library is skipped at configure time.
- libxml2 and libxslt are keg-only in Homebrew: `run.sh` passes their
  prefixes in `CMAKE_PREFIX_PATH`. `libqpdf.pc` adds `-std=c++20` to its
  `Cflags`, which the C compiler rejects: qpdf is found with `find_path` and
  `find_library`.
- Clock: `clock_gettime(CLOCK_UPTIME_RAW)` on macOS, `CLOCK_MONOTONIC`
  elsewhere. On macOS, `CLOCK_MONOTONIC` only advances in microseconds
  (measured resolution 1,000 ns, against 42 ns for `CLOCK_UPTIME_RAW`, the
  clock behind Rust's `std::time::Instant` on this platform).
- As in the Rust harness: an unmeasured check call, then N in steps of 1, 2,
  5, 10… until a series lasts `BENCHTIME`, then `COUNT` series; a single
  thread, the body of each executable runs in a thread with a 512 MiB stack;
  results are freed inside the measured loop and a `volatile` sink keeps the
  compiler from removing the call.
- CPU time and peak memory (harness shared with the C++ bench):
  `getrusage(RUSAGE_SELF)`, user + system time of all threads, read before
  and after each series (`cpu_ns_per_op`); `ru_maxrss` (bytes on macOS) read
  at the end of the series in the `resources` pass, where `run.sh` directly
  launches `build/bench-<lib> resources <Bench>`, one process per benchmark.
  None of the four libraries starts a thread or a child process.
- A library error becomes a `fail` with its message. If a library kills the
  process (SIGSEGV, SIGBUS, SIGABRT…), a signal handler writes the `fail`
  line (or the survey line) prepared before the call, then exits with code
  70: the JSONL line stays complete.
- The embedding check reads the resulting PDF back with the same library's
  extraction and requires a CII XML, like the Rust harness. With
  `FX_DUMP_EMBED=<dir>`, the PDF from that check call is kept for
  inspection with other tools (this is how MuPDF's PDFs were run through
  `zerechnung validate`, Kura and `qpdf --check`).
- The protocol's host PDFs (`host.pdf`, `host_classic.pdf`) are PDF/A-3
  documents without an invoice: no attachment, no `/AF`, no Factur-X XMP.
  The MuPDF floor still removes every invoice file from `/AF` before adding
  its own: on a PDF that already carries one, `zerechnung validate` would
  otherwise report two invoice attachments (PDF-28).
