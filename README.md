# invoice-bench

Benchmarks and a conformance comparison of Factur-X / ZUGFeRD / EN 16931 e-invoicing libraries, in nine languages. They measure [ZeRechnung](https://github.com/ritoon/zerechnung), a Go package, against the other libraries, on the same files from the official ZUGFeRD corpus and with the same protocol.

| Folder | Contents |
|---|---|
| [`ecosystem/`](ecosystem) | ZeRechnung and 14 other Go libraries: response time, memory and accuracy for extraction, parsing, validation, writing, embedding into a PDF and checking a hybrid PDF |
| [`polyglot/`](polyglot) | The TypeScript, Python, Rust, Java, C#, C, C++ and PHP libraries, with the same protocol ([`PROTOCOL.md`](polyglot/PROTOCOL.md)): time, CPU time, memory, allocations, accuracy, and the standards each library actually enforces (one `COVERAGE.md` file per language) |

## Usage

```bash
cd ecosystem && ./run.sh all        # Go bench: corpus, surveys, measurements, tables
cd polyglot && ./run.sh all         # the other languages (requirements in polyglot/README.md)
```

The Go bench needs Go 1.26 or newer, git and network access; a full run takes about ten minutes on an Apple M4 Pro. The multi-language bench needs each language's toolchain and takes several hours. The details are in each folder's README.

## Measured ZeRechnung version

[`ecosystem/facturx/go.mod`](ecosystem/facturx/go.mod) pins the ZeRechnung version that both benches measure. The `ZERECHNUNG_DIR` variable measures a local clone instead, for example a branch before it is merged:

```bash
ZERECHNUNG_DIR=~/src/zerechnung ./ecosystem/run.sh bench facturx
```

## Results

Each README gives the results of its latest run, and `results/<date>/` holds the raw outputs and the full tables. The analysis is published on the ZeRechnung wiki, in French: [Benchmarks](https://github.com/ritoon/zerechnung/wiki/Benchmarks) for the Go libraries, [Benchmarks autres langages](https://github.com/ritoon/zerechnung/wiki/Benchmarks-autres-langages) for the others.

## History

These benches lived in the ZeRechnung repository, under `bench/`, until 2026-11-22 (ritoon/zerechnung#114). Their git history was kept, and the pull request numbers in the old commit messages refer to ritoon/zerechnung. The raw outputs produced before the move show the module paths `github.com/ritoon/zerechnung/bench/…`; the report tools still read these outputs. The harness notes in the committed multi-language records were written in French and were translated into English afterwards, without touching the measurements; all the tables were regenerated from these records.

## License

The code in this repository is under the MIT license ([LICENSE](LICENSE)). The benches download third-party libraries under their own licenses, some of them restrictive (AGPL-3.0, GPL, BSL 1.1, commercial evaluation license): the list is in the Licenses section of [`ecosystem/README.md`](ecosystem/README.md#licenses) and of [`polyglot/README.md`](polyglot/README.md#licenses). The ZUGFeRD corpus (Apache-2.0) and the official artifacts are downloaded and verified, not versioned.
