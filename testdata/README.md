# Shared inputs

The inputs that every language measures, byte for byte (see [`../PROTOCOL.md`](../PROTOCOL.md#input-files)). None of them is versioned: `go/run.sh setup` (or `./run.sh setup` at the root) fetches or builds them here.

| Path | Contents |
|---|---|
| `corpus/` | Sparse checkout of the official [ZUGFeRD corpus](https://github.com/ZUGFeRD/corpus) (Apache-2.0), pinned at commit `84f8d8f` |
| `host.pdf` | The Small PDF of the corpus without its invoice: a PDF/A-3 document with no attachment, no `/AF` and no Factur-X XMP, built by `TestMakeHost` in [`go/libs/pdfcpu-pdfcpu`](../go/libs/pdfcpu-pdfcpu) |
| `host_classic.pdf` | The same, rewritten with a classic xref table and no object streams (`TestMakeHostClassic`) |

`FX_TESTDATA` points every harness to another folder with the same contents; `FX_CORPUS` replaces the corpus alone.
