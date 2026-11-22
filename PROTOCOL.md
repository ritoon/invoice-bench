# Common protocol: TypeScript, Python, Rust, Java, C#, C, C++ and PHP benches

The cross-language campaign measures the Factur-X / ZUGFeRD / EN 16931 libraries of the TypeScript, Python, Rust, Java, C# (.NET), C, C++ and PHP ecosystems with **the same protocol** as the Go bench in `go/` (see `go/common/common.go` and `harness.go`), to compare their times with those of [ZeRechnung](https://github.com/ritoon/zerechnung), the Go package.

One folder per language at the root of the repository: `go/`, `typescript/`, `python/`, `rust/`, `java/`, `csharp/`, `c/`, `cpp/`, `php/`. Each one is self-contained (its dependencies, its harness) and touches nothing else in the repository, apart from reading `testdata/` and writing to `results/`.

## Input files

Paths relative to `testdata/` (`FX_TESTDATA` overrides it; the corpus and the host PDFs are fetched and built by `go/run.sh setup`, ZUGFeRD/corpus commit `84f8d8f`). Use these exact bytes, without copying or modifying them.

| Name | File |
|---|---|
| PDF Small (149 KB, EN 16931) | `corpus/XML-Rechnung/FX/EN16931_Einfach.pdf` |
| PDF Large (1.1 MB, EN 16931) | `corpus/XML-Rechnung/FX/EN16931_Elektron.pdf` |
| PDF Minimum (89 KB, FNFE) | `corpus/ZUGFeRDv2/correct/FNFE-factur-x-examples/Facture_FR_MINIMUM.pdf` |
| PDF Extended (157 KB, ZUGFeRD 2.0) | `corpus/ZUGFeRDv2/correct/intarsys/EXTENDED/zugferd_2p0_EXTENDED_Warenrechnung.pdf` |
| XML Small (13 KB) | `corpus/XML-Rechnung/CII/EN16931_Einfach.cii.xml` |
| XML Lines (19 KB, allowances and charges) | `corpus/XML-Rechnung/CII/EN16931_Rabatte.cii.xml` |
| XML Large (1.3 MB, base64 attachment) | `corpus/XML-Rechnung/CII/XRECHNUNG_Elektron.cii.xml` |
| Host with xref stream | `host.pdf` (136 KB): the Small PDF without its invoice, a PDF/A-3 document with no attachment, no `/AF` and no Factur-X XMP |
| Host with classic xref table | `host_classic.pdf` (140 KB): the same, rewritten with a classic xref table and no object streams |

Corpus of the accuracy surveys:

- **26 valid CII XML files**: `corpus/XML-Rechnung/CII/*.cii.xml`, except `not_validating_full_invoice_based_onTest_EeISI_300_CENfullmodel.cii.xml`.
- **104 valid hybrid PDFs**: every `.pdf` under `corpus/ZUGFeRDv2/correct/` (recursive) and `corpus/XML-Rechnung/FX/*.pdf`.

## Measured operations

Same names as the Go bench:

| Benchmark | Input | Measures |
|---|---|---|
| `Extract_Small`, `Extract_Large`, `Extract_Minimum`, `Extract_Extended` | PDF | PDF → bytes of the invoice XML |
| `Parse_Small`, `Parse_Lines`, `Parse_Large` | XML | XML → the library's model |
| `ParseValidate_Small`, `ParseValidate_Lines`, `ParseValidate_Large` | XML | XML → full validation, number of errors |
| `Write_Small`, `Write_Lines` | XML | model → CII XML; the model is built **once, outside the measurement** (by the library's parser; failing that, by hand from the content of `EN16931_Einfach`, with `note: "hand-built model"`) |
| `Embed_Small` | `host.pdf` + XML Small | host PDF + XML → Factur-X PDF |
| `EmbedClassic_Small` | `host_classic.pdf` + XML Small | same, host with classic xref table |
| `CheckPDF_Small`, `CheckPDF_Large` | PDF | full check of a hybrid PDF by a **native** function of the library (container + XML); not a composition of extraction + validation |

Rules, identical to the Go bench:

1. **Check before measuring.** A first call outside the measurement checks the result: the extracted XML contains `CrossIndustryInvoice` or `CrossIndustryDocument`; the PDF produced by `Embed` can be read back and the invoice extracted from it (with the library itself, or failing that another library of the same language). If the library cannot do the operation: `status: "unsupported"`. If it returns a wrong result or throws an exception: `status: "fail"` with the message.
2. **Validation**: the most complete rule set the library applies offline to these German invoices, **without** French national rules (the corpus is German). Record which one in `note`.
3. **Embedding**: default options, fixed modification date if the library asks for one. If the library also converts to PDF/A or rejects a non-PDF/A host, record it.
4. Measurement **in process**, on bytes already in memory (files are read outside the measurement).

## Timing method

Modelled on `go test -bench . -count 3 -benchtime 1s`:

- For each benchmark: find N in increasing steps (1, 2, 5, 10, 20, 50, 100…) until a series of N calls lasts at least 1 s (`BENCHTIME`, in seconds). This ramp-up serves as a warm-up (V8 JIT, caches).
- Then run `COUNT` series (3 by default) of N calls, timed with the monotonic clock (`process.hrtime.bigint()`, `time.perf_counter_ns()`, `std::time::Instant`, `System.nanoTime()`, `Stopwatch.GetTimestamp()`, `clock_gettime(CLOCK_MONOTONIC)`, `std::chrono::steady_clock`, `hrtime(true)`). `ns_per_op` = series duration / N.
- Also record `first_call_ns`: the duration of the very first call of the operation in the process (rule compilation, lazy loading…).
- No forced GC, no parallelism: a single thread.
- Rust: Cargo's default `release` profile, without CPU-specific options. Node: the installed version (`node --version`), without options. Python: CPython, version recorded.
- Java: one JVM process per library, default JVM options (a stack or heap option is only added if a validator requires it, and it is recorded). C#: `Release` build, default .NET runtime settings (tiered compilation). C and C++: CMake `Release` build (`-O3 -DNDEBUG`), without CPU-specific options. PHP: the installed PHP CLI, settings recorded in `php/COVERAGE.md`.

## CPU time and memory

Two measurements complement wall-clock time, in every language:

- **CPU time per operation** (`cpu_ns_per_op`, next to `ns_per_op`): for each timed series, the CPU time of the process (user + system, **all threads**: garbage collector, JIT, runtime threads) consumed over the same window as the monotonic clock, divided by N. The clock is read once before and once after the series: `process.cpuUsage()`, `time.process_time_ns()`, `getrusage(RUSAGE_SELF)` (Rust, C, C++, PHP, Go), `OperatingSystemMXBean.getProcessCpuTime()`, `Environment.CpuUsage`. The CPU / wall time ratio is 100% for one busy core; above that, work runs in parallel with the measured thread.
- **Peak memory** (`max_rss_bytes`): the maximum resident set size (max RSS) of a process that loads the library and runs **only one** benchmark (check, first call, calibration, series), runtime included. Read at the end of the process with `getrusage(RUSAGE_SELF).ru_maxrss`, or measured from outside with `/usr/bin/time -l` on the process itself (never on a launcher such as npm, mvn, `dotnet run` or `cargo run`). For Go, each benchmark runs in its own test binary (`go test -c`), and `FX_CPU=1` makes the helpers of `go/common` report `cpu-ns/op`. A second Go series runs with `GOMAXPROCS=1` (row "ZeRechnung, GOMAXPROCS=1"): Go's concurrent garbage collector otherwise uses other cores, and this series shows the package limited to one core.

These two measurements come from a separate pass, `resources`: one process per library and per benchmark, `COUNT=1` and `BENCHTIME=1` by default (`RES_COUNT`, `RES_BENCHTIME` in the top-level `run.sh`). The times in the main tables remain those of the `bench` pass. The work of a child process spawned by a library (a JVM, a command-line tool) is not counted; the harness records it in `note`.

## Allocations

`run.sh alloc` measures what **one call** of each operation allocates: the bytes requested from the allocator and, when the language counts them, the number of allocations. It is all the memory the call requests, including what it frees before returning; not the memory the process occupies.

- Same result check as for the times (first call, `first_call_ns`), then a warm-up: at least 3 calls, up to 100 calls or 0.2 s.
- Then a counted series of about 0.2 s, from 1 to 1,000 calls. Reported value: the average per call.
- **Go**: `go test -benchmem` (B/op, allocs/op) gives the figures with the benchmarks; `run.sh alloc go` does nothing.
- **Rust**: the harness installs a `#[global_allocator]` that wraps `System`. During the series, each `alloc`, `alloc_zeroed` and `realloc` counts as one allocation of the requested size; outside the series, the allocator only reads an atomic flag.
- **Python**: the measurement process runs with `PYTHONMALLOC=malloc` (every object goes through `malloc`) and loads [`python/alloccount.c`](python/alloccount.c) through `DYLD_INSERT_LIBRARIES` (macOS) or `LD_PRELOAD` (Linux). This counter intercepts `malloc`, `calloc`, `realloc`, `posix_memalign` and `aligned_alloc` across the whole process, C extensions included (libxml2, qpdf, pydantic-core). Memory managed directly through `mmap` (saxonche's GraalVM heap) is not seen. The process is started with `-s -P` and without the other `PYTHON*` variables (`-I` would ignore `PYTHONMALLOC`); a preliminary check verifies that the counter does see a 4 KB allocation.
- **TypeScript**: V8 does not count its allocations. The process runs with `--expose-gc` and a 64 MB young generation (`--min-semi-space-size=64 --max-semi-space-size=64`). Bytes = growth of `heapUsed` and `arrayBuffers` during the series + what each garbage collection freed (`v8.GCProfiler`). The series is sized to allocate at most 32 MB; if a garbage collection still runs, it is shortened (down to 1 call), then the value is recorded as a minimum (`gcs` > 0, "≥" in the tables). No number of allocations (`allocs: null`). One process per benchmark.
- **Java, C#, C, C++, PHP**: no counter yet; the top-level `run.sh alloc` skips them.

## Outputs

`run.sh` in each language folder, with the same subcommands as `go/run.sh`:

```bash
./run.sh setup              # dependencies (npm, venv, cargo, Maven, dotnet, CMake, Composer)
./run.sh survey    [lib…]   # accuracy surveys
./run.sh cases     [lib…]   # conformance cases (../cases from a language folder)
./run.sh bench     [lib…]   # benchmarks
./run.sh alloc     [lib…]   # allocations per call (TypeScript, Python, Rust)
./run.sh resources [lib…]   # CPU time and peak memory, one process per benchmark
```

Variables: `OUT` (default `../results/$(date +%F)`), `COUNT` (3), `BENCHTIME` (1). Files written to `$OUT`:

`bench_<lang>.jsonl`, one line per library and benchmark:

```json
{"lang":"python","lib":"factur-x","pkg":"factur-x","version":"7.9","bench":"Extract_Small","status":"ok","ns_per_op":[51234.1,50988.0,51100.7],"cpu_ns_per_op":[51190.4,50950.2,51080.9],"iters":20000,"first_call_ns":812000,"note":""}
```

`resources_<lang>.jsonl`: the same fields, one series by default, plus `max_rss_bytes` (bytes).

`alloc_<lang>.jsonl`, one line per library and benchmark (`gcs` in TypeScript only):

```json
{"lang":"rust","lib":"en16931","pkg":"en16931-formats","version":"0.10.1","bench":"Extract_Small","status":"ok","allocs":2199.0,"bytes":1390014.0,"calls":370,"first_call_ns":2006333,"note":""}
```

`survey_<lang>.jsonl`, one line per library and operation (`validate`: valid XML files flagged out of 26; `parse`: XML files parsed out of 26; `extract`: PDFs extracted out of 104):

```json
{"lang":"python","lib":"factur-x","op":"extract","ok":104,"total":104,"failures":[{"file":"x.pdf","error":"..."}],"note":""}
```

For `validate`, `ok` = number of files **flagged** (at least one error, or an exception), `failures` = the list of those files with the number of errors and the rule identifiers. Reference: the official schematrons (CEN EN 16931 CII + KoSIT for the XRECHNUNG files) flag 6 files out of 26 (BR-CL-25 twice, CII-SR-465 four times on the XRECHNUNG examples).

`cases_<lang>.jsonl`, one line per library and suite of conformance cases: the outcome of each document (number of errors, rule identifiers or messages, exception). Format and scoring in [`cases/README.md`](cases/README.md).

## Conformance cases

`./run.sh cases` runs each library's validation on the documents of [`cases/cases.tsv`](cases/cases.tsv) (`CASES_DIR` to read another folder), with the same adapter and the same rule set as `survey` and `bench`, in one process per library (PHP: survey child processes, restarted after a crash). Each document is validated once, outside any timing; an exception or a crash is recorded as an error of the document, and so is a call that exceeds the timeout where the harness has one (Go, Python); the run goes on. Rules:

- the adapter routes a UBL document to the library's UBL rules when the library has them (by the root element, read in the first 64 KiB of the document); a library without UBL support gets the UBL documents like any other, and its rejections or failures count;
- `errors` counts the findings the library reports as errors (fatal errors, `error` severity, or every finding for a library without severities), as in `survey`; `rules` gives their identifiers, or their messages when the library has no identifiers;
- no case is skipped and none is retried: the outcome is that of the single call.

The report tool turns these outcomes into a confidence score per country: section "Confidence by country" of `tables.md`, formula in [`cases/README.md`](cases/README.md#the-score).

## Coverage of the standard

`<lang>/COVERAGE.md`: for each selected or excluded library, a library note, then a table with these columns:

| Library | Version | License | Latest version | Downloads | CII parsing | CII writing | UBL | PDF extraction | PDF embedding | XSD | EN 16931 (CII) | Factur-X profiles (MINIMUM, BASIC WL, EXTENDED) | BR-FR | BR-DE (XRechnung) | Peppol BIS 3 | PDF/A-3 container check | External dependencies |

Expected values: `yes`, `no`, `partial (n rules)`, `official XSLT via Saxon`, etc. **Check in the code** (file and line, or version of the embedded artifacts), not only in the README: number of rules actually implemented, version of the official artifacts (CEN 1.3.x, Factur-X 1.0x, XRechnung 3.0.x, FNFE-MPE 1.x), network calls, need for Java.

Exclude (with the reason in `COVERAGE.md`): remote API clients, libraries that require a JVM outside the Java bench (the Java bench measures those libraries directly), empty or name-squatting packages.
