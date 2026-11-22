# TypeScript, Python, Rust, Java, C#, C, C++ and PHP benchmarks

This folder compares [ZeRechnung](https://github.com/ritoon/zerechnung), the Go package, with the Factur-X / ZUGFeRD / EN 16931 libraries of the TypeScript, Python, Rust, Java, C# (.NET), C, C++ and PHP ecosystems: response time, CPU time and memory, allocations per call, accuracy on the official ZUGFeRD corpus, and the standards actually enforced. It complements [`ecosystem`](../ecosystem), which does the same for the Go libraries.

Each language has its own self-contained folder:

| Folder | Contents |
|---|---|
| [`PROTOCOL.md`](PROTOCOL.md) | Common protocol: input files, operations, timing, output format |
| [`typescript/`](typescript) | 13 npm libraries, Node 22, `.mjs` harness; library notes in [`typescript/COVERAGE.md`](typescript/COVERAGE.md) |
| [`python/`](python) | 8 PyPI libraries, CPython 3.13; library notes in [`python/COVERAGE.md`](python/COVERAGE.md) |
| [`rust/`](rust) | 11 crates, rustc 1.93, one binary per crate; library notes in [`rust/COVERAGE.md`](rust/COVERAGE.md) |
| [`java/`](java) | 10 Maven Central libraries (including Saxon-HE with the official XSLT and the PDFBox floor), OpenJDK 27, one Maven module and one JVM per library; library notes in [`java/COVERAGE.md`](java/COVERAGE.md) |
| [`csharp/`](csharp) | 13 NuGet packages (including SaxonCS-HE with the official XSLT and three floors), .NET 11, one project per library; library notes in [`csharp/COVERAGE.md`](csharp/COVERAGE.md) |
| [`c/`](c) | no dedicated C library: 4 floors (MuPDF, qpdf, libxml2, libxslt), CMake, Apple clang 21; library notes in [`c/COVERAGE.md`](c/COVERAGE.md) |
| [`cpp/`](cpp) | Kura (PDF/A container only), SaxonC-HE with the official XSLT, 5 floors (qpdf, PoDoFo, poppler-cpp, Xerces-C, pugixml); library notes in [`cpp/COVERAGE.md`](cpp/COVERAGE.md) |
| [`php/`](php) | 19 Packagist packages (including the smalot/pdfparser and ext-dom floors), PHP 8.6 CLI, one process per benchmark; library notes in [`php/COVERAGE.md`](php/COVERAGE.md) |
| `report/` | Tool that merges the outputs into Markdown tables |
| `results/<date>/` | Raw outputs (`bench_*`, `resources_*`, `alloc_*`, `survey_*`), environment (`env.txt`) and tables (`tables.md`) |

## Method

The operations, the input files and the measurement rules are those of the Go bench: same bytes of the ZUGFeRD corpus (commit `84f8d8f`), same benchmark names (`Extract_*`, `Parse_*`, `ParseValidate_*`, `Write_*`, `Embed_*`, `CheckPDF_*`), result checked before any measurement.

The timing reproduces `go test -bench . -count 3 -benchtime 1s`: N grows in steps 1, 2, 5, 10… until a series lasts 1 s, then 3 series of N calls; the reported value is the median. Measurements run in process, on a single thread, with each language's monotonic clock. The "First call" column gives the duration of the very first call: lazy loading, compilation of the schematrons.

A separate pass, `resources`, measures the **CPU time** per operation (user + system, all threads of the process, over the same window as the timer) and the **peak resident memory** of a process that loads one library and runs one benchmark only, runtime included (details in [`PROTOCOL.md`](PROTOCOL.md)).

Allocations are counted separately (`./run.sh alloc`): after the same result check and a warm-up, a series of about 0.2 s gives the bytes allocated and the number of allocations of one call. Go takes them from `go test -benchmem`; Rust from a counting global allocator; Python from a `malloc` counter injected into the process ([`python/alloccount.c`](python/alloccount.c), with `PYTHONMALLOC=malloc`); TypeScript from the growth of the V8 heap, without a number of allocations, which V8 does not count. The Java, C#, C, C++ and PHP harnesses do not have a counter yet. The details are in [`PROTOCOL.md`](PROTOCOL.md#allocations).

For validation, each library applies the most complete rule set it has offline, without French rules (the corpus is German). Remote API clients, libraries that require a JVM outside the Java bench, and empty packages are left out; the reason is given in each `COVERAGE.md`.

A **floor** is a generic library (PDF or XML) that knows nothing about Factur-X, wired in with a few lines in the adapter: it gives the minimum cost of an operation in the language. In C, there are only floors.

## Usage

```bash
./run.sh setup     # corpus (through ecosystem), npm ci, venv, cargo, Maven, dotnet, CMake, Composer, official artifacts
./run.sh survey    # accuracy on the corpus, all nine languages
./run.sh bench     # benchmarks, one language after the other
./run.sh alloc     # allocations per call (TypeScript, Python, Rust; Go reports them with bench)
./run.sh resources # CPU time and peak memory, one process per library and per benchmark
./run.sh report    # results/<date>/tables.md and env.txt
./run.sh all       # all six in a row
```

`survey`, `bench`, `alloc` and `resources` accept a list of languages (`./run.sh bench rust go`); each folder also has its own `run.sh`, which accepts a list of libraries. Variables: `OUT`, `COUNT` (3), `BENCHTIME` (1, in seconds), for `resources` `RES_COUNT` (1) and `RES_BENCHTIME` (1), and `ZERECHNUNG_DIR`, which measures a local clone of ZeRechnung instead of the version `../ecosystem/facturx/go.mod` requires (see [`../ecosystem`](../ecosystem/README.md#usage)).

Requirements:

- Go 1.27, Node 22 and npm, Python 3.13 (Homebrew) and a C compiler for its allocation counter, Rust 1.93, and for zugferd-rs `libxml2` and `libclang`;
- Java: JDK 21 or newer (measured with Homebrew's OpenJDK 27; `/usr/bin/java` is only a stub on macOS, `java/run.sh` uses `JAVA_HOME` or `/opt/homebrew/opt/openjdk`) and Maven 3;
- C#: .NET 11 SDK (`brew install dotnet`); no Java;
- C and C++: CMake, Ninja, pkg-config and `brew install mupdf qpdf libxml2 libxslt podofo poppler xerces-c pugixml freetype openjpeg jpeg-turbo little-cms2`; SaxonC-HE (official binary) and Kura (sources) are downloaded and pinned by `cpp/run.sh`, for macOS arm64 only;
- PHP: PHP 8.4 or newer as CLI, with the extensions listed in [`php/COVERAGE.md`](php/COVERAGE.md), Composer 2, and `pdfdetach` (poppler);
- curl, unzip, shasum and network access during `setup`.

On an Apple M4 Pro, the accuracy surveys take about ten minutes, the timing measurements about 1 h 25 min (including 17 minutes for PHP and 16 to 18 for TypeScript, because of the slow validators), the allocations an hour and a half (including 85 minutes for TypeScript, one process per benchmark) and the CPU and memory pass about 55 minutes. Run the measurements on an idle machine, on mains power or in high performance mode, lid open (`run.sh` prevents idle sleep with `caffeinate -i`, not lid-closed sleep): another process that keeps a core busy, or a sleep, skews the times.

## Results of 2026-11-22

Apple M4 Pro on battery, high performance mode (`powermode 2`). Go 1.27.4, Node 22.15.0, CPython 3.13.18, rustc 1.93.1, OpenJDK 27.0.1.1, .NET 11.0.0, Apple clang 21, PHP 8.6.0 without OPcache or JIT; full versions in [`results/2026-11-22/env.txt`](results/2026-11-22/env.txt). Full tables: [`results/2026-11-22/tables.md`](results/2026-11-22/tables.md); those of the 2026-11-17 run (v0.1.0, TypeScript, Python and Rust) remain in [`results/2026-11-17/`](results/2026-11-17/tables.md).

The results combine two runs from the same day, on the same ZeRechnung code (commits [`3a3b11a`](https://github.com/ritoon/zerechnung/commit/3a3b11a) and [`8d728bd`](https://github.com/ritoon/zerechnung/commit/8d728bd): only the website changes between the two):

- **Go, TypeScript, Python and Rust**: surveys, times and allocations at commit `3a3b11a` (30 commits after v0.1.0), with the latest versions of the libraries;
- **Java, C#, C, C++ and PHP**: surveys and times at commit `8d728bd`, with Docker Desktop running in the background (about 17% of one core);
- **CPU time and memory**, all nine languages: `resources` pass of the second run, where five libraries still had their previous version (@ekwo-ai/factur-x 0.25.2, @stackforge-eu/factur-x 1.4.6, factur-x 7.9 in Python, mcp-einvoicing-core 1.35.0 for mcp-einvoicing-de, pdf-invoice 1.4.0).

### Response time

The fastest of the other languages for each operation. A floor (generic PDF or XML library) gives the minimum cost of the operation in its language; the fastest Factur-X library is given next to it when it differs.

| Operation | ZeRechnung (Go) | Fastest elsewhere | Ratio |
|---|---|---|---|
| XML extraction (149 KB PDF) | 31.5 µs | qpdf (C, floor): 27.9 µs; Kura (C++): 28.9 µs | ×0.89; ×0.92 |
| XML extraction (1.1 MB PDF) | 32.9 µs | qpdf (C, floor): 32.8 µs; Kura (C++): 34.3 µs | ×1.00; ×1.04 |
| Parsing XML → model (13 KB) | 18.1 µs | pugixml (C++, floor, DOM tree without an invoice model): 5.6 µs; zugferd-rs (Rust): 24.5 µs | ×0.31; ×1.4 |
| Parsing XML → model (1.3 MB, base64 attachment) | 1.17 ms | pugixml (C++, floor): 238 µs; john-wink/en16931-php (PHP): 646 µs | ×0.20; ×0.55 |
| Validation, all libraries (13 KB) | 24.8 µs | pdf-invoice (Rust): 26.3 µs, 14 rules, 20 false positives out of 26 | ×1.1 |
| Full EN 16931 validation (13 KB) | 24.8 µs | official CEN and KoSIT XSLT with Saxon-HE (Java): 419 µs | ×17 |
| Writing model → XML (13 KB) | 3.3 µs | pugixml (C++, floor): 3.7 µs; zugferd-rs (Rust): 7.0 µs | ×1.1; ×2.1 |
| Embedding into a PDF with a classic xref table | 96.2 µs | pdfa (Rust): 88.6 µs, fails on a host with an xref stream | ×0.92 |
| Embedding into a PDF with an xref stream | 96.7 µs | MuPDF (C, floor): 123 µs; pdf-invoice (Rust): 180 µs | ×1.3; ×1.9 |
| Full check of a hybrid PDF (149 KB) | 811 µs | zugferd-pdf (Python): 213 µs, container checks only | ×0.26 |

Full EN 16931 validation (CEN schematrons, with or without KoSIT), the fastest in each language, on the 13 KB XML:

| Language | Library | Time | First call | Valid flagged /26 |
|---|---|---|---|---|
| Go | **ZeRechnung**, compiled rules | 24.8 µs | — | 6 |
| Java | Saxon-HE + official XSLT | 419 µs (×17) | 629 ms | 6 |
| Rust | einvoice-validator (CEN XSLT run in Rust) | 483 µs (×19) | 23 ms | 2 |
| C# | SaxonCS-HE + official XSLT | 528 µs (×21) | 707 ms | 6 |
| TypeScript | @sweberdev/summand | 854 µs (×34) | 18 ms | 6 |
| C++ | SaxonC-HE + official XSLT | 862 µs (×35) | 269 ms | 6 |
| Python | saxonche + official XSLT | 900 µs (×36) | 137 ms | 6 |
| PHP | dealerweb/einvoice (CEN XSLT compiled to PHP) | 8.03 ms (×323) | 24 ms | 6 |

- The qpdf (C, C++) and pugixml (C++) floors beat the Go package on extraction and parsing: they only do the generic part, finding the stream or building a DOM tree. On extraction, qpdf ties with it on the 1.1 MB PDF and is behind it on the MINIMUM PDF (14.9 µs against 11.1 µs). On writing, the Go package beats pugixml (3.3 µs against 3.7 µs; 5.3 µs against 5.5 µs on 19 KB).
- Among the Factur-X libraries, the Go package is the fastest on 7 columns out of 16: 1.1 MB extraction, 13 and 19 KB parsing, 13 KB validation, 13 and 19 KB writing, embedding into a host with an xref stream. Against the TypeScript, Python and Rust libraries alone, it is the fastest on 10 columns (13 on 2026-11-17).
- The nine other columns:
  - extraction: Kura (C++) reads the 149 KB PDF in 28.9 µs (×0.92) and the EXTENDED PDF in 28.2 µs (×0.78); pdfa (Rust) reads the MINIMUM PDF in 504 ns by searching the bytes for the first `/EmbeddedFile` stream, without reading the PDF structure, which misses 72 of the 104 PDFs of the corpus;
  - parsing the 1.3 MB XML: the Go package decodes the base64 attachment (3.5 MB allocated per parse); john-wink/en16931-php (PHP) parses the file in 646 µs, mcp-einvoicing-de (Python, lxml) in 704 µs while keeping the attachment as text;
  - validation: pdf-invoice (Rust, 14 rules) on 19 KB (40.1 µs against 45.4 µs), pat-o-dev/factur-x (PHP, a few hand-written rules) on 1.3 MB (466 µs against 794 µs);
  - embedding into a host with a classic xref table: pdfa (Rust), 88.6 µs against 96.2 µs, which rejects hosts with an xref stream;
  - full check of the 149 KB and 1.1 MB PDFs: zugferd-pdf (Python, 213 and 269 µs), @attestwire/en16931 (TS, 368 and 403 µs) and en16931 (Rust, 489 and 817 µs) read neither the page content nor the fonts (next point).
- **The hybrid PDF check went from 75.0 µs (2026-11-17, v0.1.0) to 811 µs.** Since PR ritoon/zerechnung#73, `Check` verifies the PDF/A requirements across the whole file (every object, page content, font programs, JPEG 2000 images): the CPU profile shows the time spent decompressing and reading the TrueType font programs (`pdfaScan.font`), and a call allocates 10 times more bytes (700 KB instead of 69 KB). The faster libraries only check the container and the XMP.
- The validators that run the official schematron cost from 0.4 ms to 101 ms per invoice, except @lucaapp/zugferd-ts (from 2.7 to 9.2 s). The slowest: mustang (Java) 11 ms, erechnung-validator (Rust) 32 ms, securibox (C#) 38 ms, @kontor-mcp/core (TS) 48 ms, vampire-php (PHP) 64 ms, @conformo/toolkit (TS) 72 ms, mcp-einvoicing-de (Python) 101 ms, which recompiles its XSLT on every call. Their first call compiles the stylesheets: from 18 ms (summand) to 2.1 s (KoSIT validator). The Go package's rules are compiled code: there is nothing to load.
- Embedding: Mustang (Java) takes 5.2 ms (`ZUGFeRDExporterFromPDFA`, PDF saved without compression); several libraries reject the bench's hosts or leave two `factur-x.xml` in them (see "Known specifics").

### Allocations per call

Bytes allocated · number of allocations, for the library of each language that allocates the least among those that succeed (TypeScript: bytes only; Java, C#, C, C++ and PHP: no counter yet):

| Operation | ZeRechnung (Go) | Other Go | TypeScript | Python | Rust |
|---|---|---|---|---|---|
| Extraction (149 KB PDF) | 39 KB · 228 | DCNWS 184 KB · 4,926 | @sweberdev/summand 235 KB | mcp-einvoicing-de 186 KB · 1,553 | **pdfa 36 KB · 18**, byte search |
| Parsing (13 KB) | **28 KB · 111** | fX 53 KB · 1,165 | @ekwo-ai/factur-x 164 KB | pycheval 374 KB · 3,302 | zugferd-rs 81 KB · 77 |
| Full EN 16931 validation (13 KB) | **23 KB · 82** | pdf0 67 KB · 1,462 | @sweberdev/summand 3.1 MB | saxonche 321 KB · 749, GraalVM heap not counted | einvoice-validator 397 KB · 7,188 |
| Writing (13 KB) | **8 KB · 1** | fX 32 KB · 175 | @stackforge-eu/factur-x 44 KB | pycheval 122 KB · 1,254 | core-invoice 24 KB · 396 |
| Embedding (host with xref stream) | **245 KB · 456** | pdf0 1.3 MB · 9,593 | @ekwo-ai/factur-x 1.0 MB | mcp-einvoicing-de 3.4 MB · 5,847 | pdf-invoice 949 KB · 615 |
| Check (149 KB PDF) | 700 KB · 2,383 | DCNWS 896 KB · 24,617 | @attestwire/en16931 755 KB | **zugferd-pdf 466 KB · 4,008**, container only | en16931 1.5 MB · 3,834 |

- The Go package allocates the least to parse, validate, write and embed. The two libraries that allocate less than it to extract and check do less work: pdfa fails on 72 of the 104 PDFs, zugferd-pdf only checks the container.
- pdf-invoice (64) and zugferd-rs (77) make fewer allocations than it to parse (111), but for three times as many bytes.
- Python and TypeScript allocate far more for the same operation: parsing 13 KB costs from 374 KB to 2.2 MB in Python, in thousands of allocations, and from 164 KB to 7.4 MB in TypeScript. The heaviest full validators allocate 151 MB (erechnung-validator, 850,719 allocations), at least 60.9 MB (@conformo/toolkit) and at least 13.2 GB (@lucaapp/zugferd-ts) per 13 KB invoice.

### CPU time and memory

`resources` pass of the second run, from 12:17 to 13:10, on mains power: one process per library and per benchmark, one series of one second. In each language, the median wall-clock time of this pass is within 5% of that of the `bench` pass of the same run; for Go, TypeScript, Python and Rust, the response times above come from the other run. Full tables per operation in `tables.md` ("CPU time: …" and "Peak memory: …").

| Language | Libraries | Processes | Peak memory min. | median | max. | CPU / wall time median | max. |
|---|---|---|---|---|---|---|---|
| Go | 1 | 16 | 16.8 MiB | 20.9 MiB | 25.5 MiB | 163% | 206% |
| TypeScript | 13 | 119 | 65.4 MiB | 174 MiB | 509 MiB | 102% | 433% |
| Python | 8 | 69 | 39.8 MiB | 59.5 MiB | 387 MiB | 100% | 100% |
| Rust | 11 | 74 | 1.9 MiB | 7.2 MiB | 503 MiB | 100% | 111% |
| Java | 10 | 63 | 259 MiB | 602 MiB | 2,129 MiB | 125% | 399% |
| C# | 13 | 111 | 69.6 MiB | 90.0 MiB | 509 MiB | 100% | 506% |
| C | 3 | 18 | 2.2 MiB | 5.8 MiB | 11.7 MiB | 100% | 100% |
| C++ | 7 | 40 | 1.7 MiB | 9.1 MiB | 278 MiB | 100% | 100% |
| PHP | 19 | 135 | 32.0 MiB | 43.0 MiB | 1,173 MiB | 100% | 100% |

Without the Go series at `GOMAXPROCS=1` or the two PHP libraries that spawn child processes. Peak memory: maximum resident set size (max RSS) of the process, runtime included. CPU / wall time: 100% = one core busy during the timed series.

Per operation, the Go package and the fastest Factur-X library of each language in this pass (a floor if there is none); for validation, the fastest one that applies the full EN 16931 schematron. In each cell: CPU time per operation, its share of wall-clock time, peak memory of the process.

| Language | Extraction (149 KB) | Parsing XML → model (13 KB) | Full EN 16931 validation (13 KB) | Embedding (host with xref stream) | PDF check (149 KB) |
|---|---|---|---|---|---|
| **Go** | 46.9 µs (152%), 21.1 MiB | 32.0 µs (174%), 19.2 MiB | 34.5 µs (141%), 18.7 MiB | 193 µs (201%), 23.9 MiB | 1.15 ms (133%), 22.8 MiB |
| **Go, `GOMAXPROCS=1`** | 33.0 µs (99%), 17.4 MiB | 22.3 µs (99%), 16.4 MiB | 29.2 µs (99%), 16.7 MiB | 106 µs (98%), 18.6 MiB | 892 µs (100%), 19.9 MiB |
| TypeScript | @attestwire/en16931: 194 µs (102%), 142 MiB | @ekwo-ai/factur-x: 72.2 µs (100%), 104 MiB | @sweberdev/summand: 1.00 ms (111%), 154 MiB | @stackforge-eu/factur-x: 1.34 ms (109%), 185 MiB | @attestwire/en16931: 377 µs (102%), 157 MiB |
| Python | zugferd-pdf: 79.2 µs (100%), 40.1 MiB | pycheval: 226 µs (100%), 42.5 MiB | saxonche + XSLT: 893 µs (100%), 214 MiB | mcp-einvoicing-de: 667 µs (100%), 96.8 MiB | zugferd-pdf: 211 µs (100%), 40.3 MiB |
| Rust | einvoice-validator: 32.6 µs (100%), 3.1 MiB | zugferd-rs: 24.4 µs (100%), 6.6 MiB | einvoice-validator: 479 µs (100%), 26.7 MiB | pdf-invoice: 181 µs (104%), 9.0 MiB | en16931: 479 µs (100%), 5.5 MiB |
| Java | Paniolo: 119 µs (100%), 272 MiB | imixs-e-invoice: 72.1 µs (100%), 536 MiB | Saxon-HE + XSLT: 1.96 ms (399%), 341 MiB | Mustang: 8.66 ms (154%), 295 MiB | Paniolo: 19.16 ms (165%), 612 MiB |
| C# | Antiphon: 78.6 µs (100%), 73.6 MiB | Antiphon: 27.9 µs (100%), 74.3 MiB | SaxonCS-HE + XSLT: 527 µs (100%), 142 MiB | Balsoft: 1.07 ms (115%), 86.8 MiB | OfficeIMO: 229 µs (100%), 98.0 MiB |
| C | qpdf (floor): 28.5 µs (100%), 3.6 MiB | libxml2 (floor): 53.9 µs (100%), 2.8 MiB | — | MuPDF (floor): 127 µs (100%), 11.5 MiB | — |
| C++ | Kura: 28.3 µs (100%), 4.4 MiB | pugixml (floor): 5.3 µs (100%), 1.7 MiB | SaxonC-HE + XSLT: 823 µs (100%), 201 MiB | Kura: 13.03 ms (100%), 12.1 MiB | Kura: 1.54 ms (100%), 4.7 MiB |
| PHP | quittio/core: 72.2 µs (100%), 39.2 MiB | tiime/cross-industry-invoice: 226 µs (100%), 40.5 MiB | dealerweb/einvoice: 7.61 ms (100%), 51.4 MiB | dealerweb/einvoice: 2.03 ms (100%), 54.4 MiB | — |

- **The Go package uses more than one core**: from 125 to 206% of wall-clock time, 163% median. The measured thread is alone, but Go's concurrent garbage collector works on other cores (`GOMAXPROCS` is 14 on this machine). Limited to one core, the package uses 5 to 47% less CPU for a wall-clock time 4 to 38% longer: 33 µs per extraction, 22 µs per parse, 29 µs per validation, 4.4 µs per write, 107 µs per embedding. It then remains the fastest of the Factur-X libraries on parsing 13 KB, writing, and embedding into a host with an xref stream.
- **Single-threaded runtimes stay at 100%**: Python, PHP, C, C++ and Rust; in Rust, lopdf keeps a rayon thread next to the measured thread (104 to 111% for zugferd-rs and pdf-invoice).
- **JVM**: 125% median, up to 399% (Saxon-HE on validation, 231% for imixs-e-invoice on a parse). In a fresh JVM, JIT compilation and the G1 garbage collector are still running during the one-second series; in the `bench` pass, with a warm JVM, PDFBox is at 101 or 102%.
- **.NET**: 100% median. ZUGFeRD-csharp goes up to 438 and 506% on embedding: on every call, about twenty threads of the .NET pool look for fonts in the system folders.
- **Node**: 102% median; @stafyniaksacha/facturx goes up to 433% on writing, with the V8 threads.
- **Peak memory of the Go package: 17 to 26 MiB** depending on the operation, runtime included (16 to 24 MiB on one core). That is 2 to 10 times more than Rust, C and C++ (2 to 13 MiB, excluding validators and the 1.3 MB XML), less than Python (at least 40 MiB, 60 median) and PHP (at least 32 MiB, 43 median), and far less than Node (65 to 509 MiB), .NET (70 to 509 MiB) and the JVM (259 MiB to 2.1 GiB).
- **Java**: the peak mostly measures the size of the G1 heap, not the live data. With the default options, the maximum heap is a quarter of the machine's memory (6 GiB here) and G1 grows it under allocation pressure: PDFBox keeps about 2 MB live in a 448 MB heap. The KoSIT validator reaches 2.1 GiB on the 1.3 MB XML, Mustang 1.8 GiB on a parse.
- **Validating with the official XSLT costs 140 to 340 MiB** whatever the language (Saxon in Java, C#, C++ and Python), against 27 MiB for einvoice-validator (Rust) and 19 MiB for the Go package.
- **Other peaks**: quittio/core (PHP) 1.1 GiB when parsing the 1.3 MB XML, En16931 (C#) and @stafyniaksacha/facturx (TS) 509 MiB, erechnung-validator (Rust) 503 MiB, mcp-einvoicing-de (Python) 387 MiB.
- **Child processes**: mahdiabderraouf/facturx-php and php-erechnung-toolkit run `pdfdetach` (poppler) on every extraction; their CPU (about 22 and 26 ms per call) and their memory are not counted (marked † in `tables.md`).
- **A fresh process does not always run at the same speed** as the `bench` pass: in .NET, several libraries are 2 to 4 times slower in a one-second process (tiered compilation not finished); Kura, alone in its process, embeds and checks twice as fast.

### Accuracy on the corpus

Of the 26 valid CII XML files of the corpus, the official schematrons (CEN EN 16931 1.3.17 and KoSIT XRechnung) flag 6: BR-CL-25 twice, CII-SR-465 four times on the XRECHNUNG examples. CII-SR-465 is a warning in the CEN artifacts; the KoSIT configuration raises it to an error. A validator faithful to CEN alone therefore flags 2.

| Flagged | Libraries |
|---|---|
| 6, like the reference | **ZeRechnung (Go)**; TypeScript: @conformo/toolkit, @sweberdev/summand, en16931 (without telling warnings from errors); Python: saxonche + official XSLT; getafix too, but through 6 parsing exceptions; Rust: erechnung-validator; Java: KoSIT validator, Saxon-HE + official XSLT, Paniolo; C#: SaxonCS-HE + official XSLT, pax.XRechnung.NET and Securibox, but for other reasons (pax.XRechnung.NET: a FeRD rule on 2 files and the rejection of the 4 XRECHNUNG files; Securibox: CII-SR-465 and code list checks instead of BR-CL-25); C++: SaxonC-HE; PHP: dealerweb/einvoice, intragile/php-einvoicing, vampire-php, quittio/core |
| 2, like CEN alone | @attestwire/en16931, @kontor-mcp/core, mcp-einvoicing-de, einvoice-validator, en16931 (Rust), Mustang, phive, International.EInvoicing, Antiphon, En16931 (C#), pat-o-dev/factur-x |
| 1 to 5 | php-erechnung-toolkit (1), @ekwo-ai/factur-x (3, IBAN keys), ZUGFeRD-csharp (3), atgp/factur-x, tiime/factur-x and facturx-php (4, Factur-X XSD only), Balsoft (5) |
| 7 to 26, false positives | core-invoice (7), zugferd-rs (7), @lucaapp/zugferd-ts (8), @stafyniaksacha/facturx (12), FacturXDotNet (15), factur-x Rust (17), factur-x-kit (18), pdf-invoice (20), OfficeIMO (23), invoice4j (26) |
| 0 | @stackforge-eu/factur-x, @fin.cx/einvoice, factur-x (Python), drafthorse, ph-cii, horstoeko/zugferd, horstoeko/invoicesuite, easybill/zugferd-php, john-wink/en16931-php, stann/factur-x, pn-invoice; the XSD floors (System.Xml.Schema, libxml2, Xerces-C, ext-dom) |

All 104 hybrid PDFs are extracted by the Go package and by most of the libraries that can extract, floors included. The others: ZUGFeRD-csharp and MuPDF 103, @stafyniaksacha/facturx, atgp/factur-x and facturx-php 100, pdf-invoice 101, horstoeko/invoicesuite 88, Securibox 74, FacturXDotNet and tiime/factur-x 70, pycheval 65, pdfa 32, pat-o-dev/factur-x 1.

### Enforced standards

Surveyed in the libraries' code, not in their README (files and lines in each `COVERAGE.md`). A selection: the most downloaded and the most complete in each language.

| Library | Language | XSD | EN 16931 (CII) | Factur-X profiles | BR-FR | BR-DE (XRechnung) | Peppol BIS 3 in CII | PDF/A-3 container | UBL |
|---|---|---|---|---|---|---|---|---|---|
| **ZeRechnung** | Go | Factur-X 1.09.2, CII D16B (XRECHNUNG), UBL 2.1 | 793 assertions, CEN 1.3.17 | 544 rules (schematrons 1.09 and 1.09.2) | FNFE-MPE 1.4.0.04 | KoSIT 3.0.2, Leitweg-ID and GEBA | 3.0.21 (90 assertions) | 69 container checks, page content and fonts read; partial conversion | parsing, writing, EN 16931, BR-FR, BR-DE, Peppol validation |
| @stackforge-eu/factur-x (94,000/week) | TS | Factur-X 1.09 | no | XSD only | no | no | no | no | no |
| @e-invoice-eu/core (63,000/week) | TS | no | no (generation) | generation | no | generation | no | no | writing |
| @conformo/toolkit | TS | yes | CEN 1.3.17 (SaxonJS) | XSD only | Flux 2, if `country: "FR"` | KoSIT 3.0.2 | UBL only | partial | yes |
| @sweberdev/summand | TS | yes | CEN 1.3.16 | XMP compared with BT-24 | no | KoSIT 2.6.0 | no | partial | validation |
| @stafyniaksacha/facturx | TS | Factur-X 1.09 | through the Factur-X 1.09 schematron (SaxonJS) | yes, without severities | no | no | no | no | no |
| factur-x, Akretion (411,000/month) | Python | Factur-X 1.09 | only through a Saxon HTTP server | XSD; schematron through a server | 124 rules, through a server | no | no | no | parsing, writing |
| drafthorse (157,000/month) | Python | Factur-X 1.0.07 | no | XSD only | no | no | no | no | no |
| mcp-einvoicing-de | Python | no | FeRD 1.09.2 or CEN 1.3.16 (saxonche) | yes (saxonche) | no | KoSIT 3.0.2 | no | no | yes |
| en16931 + en16931-formats (9,700 downloads) | Rust | no | CEN 1.3.17 rewritten, 447 of the 522 CII-SR and CII-DT | partial | no | KoSIT 2.6.0 | yes, 3.0.21 | partial | yes |
| einvoice-validator | Rust | yes | CEN 1.3.17 XSLT run in Rust | only with the `fetch` feature (download) | same | same | same | no | validation |
| erechnung-validator | Rust | yes | CEN 1.3.16 | Factur-X 1.09.2, 5 profiles | no | KoSIT 3.0.2 | no | no | validation |
| zugferd-rs | Rust | Factur-X 1.09.2 | partial (XPath 1.0) | XSD | no | no | no | 21 veraPDF rules | yes |
| Mustang (Java) | Java | Factur-X 1.09.2 | through the Factur-X EN16931 schematron | 5 profiles (Factur-X 1.09.2) | Flux 2 1.3.0, between French companies | KoSIT 2.4.0 | no | veraPDF, XMP, attachment | parsing |
| phive + phive-rules | Java | CII D16B, Factur-X 1.09.2 | CEN 1.3.17 XSLT | 5 profiles (ZUGFeRD 2.5.2) | separate package | KoSIT 2.6.0, KoSIT levels | separate package | no | validation |
| KoSIT validator | Java | CII D16B | CEN 1.3.16 XSLT | no | no | KoSIT 2.6.0, KoSIT levels | no | no | validation |
| Paniolo | Java | CII D16B, Factur-X 1.09.2 | CEN 1.3.16 XSLT (KoSIT configuration) | 5 profiles | EXTENDED-CTC-FR profile only | KoSIT XRechnung 3.0.2 configuration | no | veraPDF, XMP, relationship | parsing, writing |
| ZUGFeRD-csharp (1.8 M downloads) | C# | no | totals recomputed | parsing and writing, without rules | no | no | no | no | XRechnung UBL |
| International.EInvoicing | C# | CII D16B | CEN 1.3.16, in-house engine | not offline | not loaded | partial (52 assertions) | separate package | no | yes |
| Antiphon (commercial) | C# | CII D16B | CEN 1.3.16 precompiled | generation only | no | KoSIT 2.6.0, without levels | separate package | no | separate packages |
| Kura | C++ | no | no | profile inferred for the XMP | no | no | no | partial (payload, /AF, XMP) | no |
| horstoeko/zugferd (1.25 M/month) | PHP | Factur-X, 5 profiles | not offline (KoSIT through Java) | XSD only | no | KoSIT through Java only | no | veraPDF through Java only | no |
| dealerweb/einvoice | PHP | CII D16B, UBL | CEN 1.3.16 XSLT compiled to PHP | FeRD 1.09.2 schematron | FNFE-MPE 1.4.0 (option) | XRechnung 3.0.2, KoSIT levels | BIS 3.0.21 rewritten | no | yes |
| intragile/php-einvoicing | PHP | CII D16B | CEN 1.3.16 schematron precompiled to PHP | Factur-X 1.09.2 schematron | no | XRechnung 2.6.0, KoSIT levels | not installed | no | yes |

No library, in any of the eight other languages, combines the full EN 16931 rules, the Factur-X profile rules, the French rules, the German rules, Peppol in CII and the container check. The closest:

- **dealerweb/einvoice** (PHP): CEN, FeRD, KoSIT, BR-FR as an option, Peppol rewritten, UBL; no PDF container check; 8 ms per validation.
- **Paniolo** (Java): CEN, KoSIT, Factur-X schematrons, veraPDF for the container; BR-FR for the EXTENDED-CTC-FR profile only, no Peppol; 0.8 ms per validation, 12 ms per PDF check.
- **@conformo/toolkit** (TS): CEN, KoSIT and BR-FR, UBL, PDF/A-3 creation; no Factur-X profile schematron, Peppol in UBL only, depends on SaxonJS (Saxonica license); 72 ms per validation.
- **erechnung-validator** (Rust): CEN, KoSIT and Factur-X schematrons for the 5 profiles; no BR-FR, no Peppol, no container, no parsing into a model; 32 ms.

The most downloaded libraries do not validate the business rules: @stackforge-eu/factur-x, factur-x (Python), atgp/factur-x and horstoeko/zugferd offline stop at the XSD; ZUGFeRD-csharp recomputes the totals; @e-invoice-eu/core and node-zugferd only generate.

**PDF/A-3 conversion** of the host PDF, which zugferd-rs, @lucaapp/zugferd-ts, @conformo/toolkit and Kura do (mcp-einvoicing-de adds an sRGB OutputIntent), also exists in the Go package: it is partial and rejects what would require a PDF engine, such as embedding a font; the bench does not measure it. Mustang and Paniolo hand the PDF/A check to **veraPDF**, the reference tool; the Go package has its own checks. It also parses, writes and validates the **UBL** syntax, like several of these libraries (@conformo/toolkit, en16931, zugferd-rs, mcp-einvoicing-de, @lucaapp/zugferd-ts…), but the bench only measures CII.

## Known specifics

- **pdfa (Rust)** extracts by byte search: 504 ns and 36 KB on the MINIMUM PDF, but only 32 PDFs out of 104.
- **Securibox (C#)** fails on `Extract_Extended` in the `bench` pass ("Schema namespace URI and prefix mismatch", XMP exception) but succeeds in the `resources` pass, alone in its process: the failure depends on the benchmarks already run in the process.
- **erechnung-validator** appears in `CheckPDF` with `beide_pdf`, its own function, which extracts and validates the XML without checking the container.
- **en16931 (Rust)** refuses to embed into the two host PDFs, because their XMP already contains `fx:` properties. **@stafyniaksacha/facturx** fails on any host whose trailer has an `/ID`.
- **Hand-built models** for `Write_Small` when the library cannot parse CII: @stackforge-eu/factur-x, @e-invoice-eu/core, node-zugferd, @ekwo-ai/factur-x, zugferd-pdf, erechnung, zugferd (Rust). Their `Write_Lines` is not measured.
- **One process per benchmark** in TypeScript, Python and PHP: "First call" there is a cold call. In Rust, Java, C#, C and C++, one process per library: in Java and C#, the first call therefore includes class loading and JIT compilation.
- **The host PDFs are not empty.** `host.pdf` and `host_classic.pdf` no longer have an `/EmbeddedFiles` tree, but their catalog keeps in `/AF` the file specification of the old `factur-x.xml` (with its stream) and a Factur-X XMP. The `Embed` check only verifies that a CII can be read back; the observed effects:
  - refusal to embed a second invoice: Paniolo (Java), intragile/php-einvoicing and pat-o-dev/factur-x (PHP);
  - PDF produced with two `factur-x.xml` in `/AF`: Mustang (Java), Kura (C++), whose PDF is flagged PDF-28 by `zerechnung validate`;
  - the C and C++ floors, OfficeIMO (C#) and the other embedders remove or replace the old entry.
- **Java**: Mustang's `Extract` includes parsing into the model (no raw extraction API). Mustang and phive apply their own severity levels and leave CII-SR-465 as a warning. No JVM options: the 512 MiB stack comes from a thread created by the harness.
- **C#**: the programs run through `dotnet <dll>` (the apphost does not find a Homebrew .NET without `DOTNET_ROOT`); invariant culture set in the code (`InvariantGlobalization` breaks SaxonCS); 256 MiB stack for the official XSLT. Antiphon is closed commercial software, measured under an evaluation license.
- **C and C++**: `CLOCK_UPTIME_RAW` clock on macOS (the one behind `Instant` in Rust), because `CLOCK_MONOTONIC` advances there in steps of 1 µs. The floors do not add an OutputIntent and do not convert to PDF/A; Kura converts to PDF/A-3b, hence its 12.6 ms (host with xref stream) and 25 ms (classic xref) per embedding. libxslt only does XSLT 1.0: the official XSLT (XSLT 2.0) do not compile with it.
- **PHP**: PHP CLI as installed, OPcache off (`opcache.enable_cli=0`), hence no JIT; `memory_limit=4G` and `zend.assertions=-1` passed through `PHP_OPTS`. No XSLT 2.0 engine is usable in PHP 8.6 (ext-xsl is XSLT 1.0, the SaxonC 13.1.1 extension does not compile with PHP 8.6): the full validators are rewrites in PHP. FPDI's free parser does not read xref streams, hence the `Embed_Small` failures of the FPDI or mPDF embedders. The installation passes `--no-scripts --no-plugins`.
- **Installation**: `npm ci --ignore-scripts` (the postinstall of @fin.cx/einvoice downloads unpinned schematrons), Python wheels only (except `tagic`, pure Python, reviewed), `cargo build --ignore-rust-version` for einvoice-validator, `RAYON_NUM_THREADS=1` for lopdf.
- **Allocations**: the Python counter does not see the memory that a library manages itself through `mmap`, like the GraalVM heap of saxonche, whose figures are underestimated. In TypeScript, "≥" means that a garbage collection ran during the series: the `ArrayBuffer`s it freed are missing, the value is a minimum.
- The full check of a PDF by the Go package measures 811 µs (2,383 allocations) at commit `3a3b11a`, against 75.0 µs (504 allocations) in v0.1.0: since PR ritoon/zerechnung#73, `zerechnung.Check` verifies the PDF/A requirements across the whole file (every object, page content, font programs, JPEG 2000 images).

## Licenses

The code in this folder is under the MIT license, like the repository. The dependencies are installed by `run.sh setup` and are not imported by the `zerechnung` package. Some have restrictive licenses:

- AGPL-3.0: pdf-invoice (Rust), MuPDF (C), Kura (C++), php-erechnung-toolkit and php-pdf-toolkit (PHP);
- GPL: poppler (C++), vampire-php (PHP), veraPDF (GPL-3.0 or MPL-2.0, pulled in by Mustang and Paniolo in Java);
- LGPL: PoDoFo (C++), smalot/pdfparser (PHP);
- EUPL-1.2: @stackforge-eu/factur-x, En16931 (C#), intragile/einvoicing-rules-en16931 (PHP), CEN artifacts;
- MPL-2.0: saxonche, Saxon-HE (Java), SaxonCS-HE (C#), SaxonC-HE (C++); Saxonica license: SaxonJS, used by @conformo, @kontor-mcp and @stafyniaksacha;
- closed commercial license: Antiphon (C#), measured under an evaluation license (free, no key, "to evaluate before purchase"); the C# bench pins PDFsharp 6.2.4 and 7.0.0-preview-1 (International.EInvoicing, Balsoft, Antiphon, Securibox, floor), MIT versions: PDFsharp 6.5.0 and later, and 7.0.0-preview-2, are AGPL-3.0;
- no declared license: stann/factur-x (PHP).

The official artifacts downloaded by the Python, Java, C# and C++ benches (CEN 1.3.17, KoSIT schematron 2.6.0 and XRechnung configuration) are verified by SHA-256 and are not versioned.
