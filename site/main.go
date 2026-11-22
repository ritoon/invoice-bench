// Command site builds the benchmark micro-site published on GitHub Pages.
//
// It reads the raw outputs of the latest run of both campaigns (the
// cross-language one in results/<date>, the Go one in go/results/<date>),
// turns them into two JSON files, and writes them next to the static pages of
// web/:
//
//   - data/perf.json: CPU time, wall time, peak memory and allocations per
//     library and per benchmark (resources and alloc passes), shown by
//     index.html;
//   - data/compliance.json: the confidence score per country (conformance
//     cases, scored by the cases package), accuracy on the ZUGFeRD corpus
//     (surveys) and the enforced standards (standards.json, kept by hand),
//     shown by compliance.html.
//
// Usage, from this folder:
//
//	go run . -out ../_site        # build the site into ../_site
//	go run . -serve :8000         # build the data and serve web/ live
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ritoon/invoice-bench/cases"
)

func main() {
	root := flag.String("root", "..", "repository root")
	resDir := flag.String("results", "", "results folder of the cross-language campaign (default: latest under results)")
	goDir := flag.String("go", "", "results folder of the Go campaign (default: latest under go/results)")
	out := flag.String("out", "", "folder to write the site into")
	serve := flag.String("serve", "", "address to serve web/ and the generated data on, e.g. :8000")
	flag.Parse()
	if *out == "" && *serve == "" {
		fmt.Fprintln(os.Stderr, "usage: site -out <dir> | -serve <addr>")
		os.Exit(2)
	}
	if _, err := os.Stat(filepath.Join("web", "index.html")); err != nil {
		log.Fatal("web/index.html not found: run the command from the site folder")
	}

	var err error
	if *resDir == "" {
		if *resDir, err = latest(filepath.Join(*root, "results")); err != nil {
			log.Fatal(err)
		}
	}
	if *goDir == "" {
		if *goDir, err = latest(filepath.Join(*root, "go", "results")); err != nil {
			log.Fatal(err)
		}
	}

	env, err := loadEnv(filepath.Join(*resDir, "env.txt"))
	if err != nil {
		log.Fatal(err)
	}
	run := Run{Date: filepath.Base(*resDir), GoDate: filepath.Base(*goDir), Env: env}

	perf, err := buildPerf(*resDir, run)
	if err != nil {
		log.Fatal(err)
	}
	comp, err := buildCompliance(*root, *resDir, *goDir, run)
	if err != nil {
		log.Fatal(err)
	}
	data := map[string][]byte{}
	for name, v := range map[string]any{"perf.json": perf, "compliance.json": comp} {
		b, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}
		data[name] = b
	}

	if *out != "" {
		if err := write(*out, data); err != nil {
			log.Fatal(err)
		}
		log.Printf("site written to %s (results %s, go/results %s)", *out, run.Date, run.GoDate)
	}
	if *serve != "" {
		log.Printf("serving web/ on http://localhost%s (results %s, go/results %s)", *serve, run.Date, run.GoDate)
		log.Fatal(http.ListenAndServe(*serve, handler(data)))
	}
}

// latest returns the most recent <date> folder under dir.
func latest(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}
	date := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	best := ""
	for _, e := range entries {
		if e.IsDir() && date.MatchString(e.Name()) && e.Name() > best {
			best = e.Name()
		}
	}
	if best == "" {
		return "", fmt.Errorf("no results/<date> folder in %s", dir)
	}
	return filepath.Join(dir, best), nil
}

// write copies web/ into out and adds the generated data files. out must be
// missing, empty, or a previous build (it holds the .nojekyll marker): any
// other folder is left alone rather than wiped.
func write(out string, data map[string][]byte) error {
	if entries, err := os.ReadDir(out); err == nil && len(entries) > 0 {
		if _, err := os.Stat(filepath.Join(out, ".nojekyll")); err != nil {
			return fmt.Errorf("%s is not empty and is not a previous build of the site: choose another -out folder", out)
		}
		if err := os.RemoveAll(out); err != nil {
			return err
		}
	}
	err := filepath.WalkDir("web", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel("web", p)
		dst := filepath.Join(out, rel)
		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, b, 0o644)
	})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(out, "data"), 0o755); err != nil {
		return err
	}
	for name, b := range data {
		if err := os.WriteFile(filepath.Join(out, "data", name), b, 0o644); err != nil {
			return err
		}
	}
	// GitHub Pages would otherwise run Jekyll on the folder.
	return os.WriteFile(filepath.Join(out, ".nojekyll"), nil, 0o644)
}

// handler serves web/ from disk, so that edits show on reload, and the
// generated data from memory.
func handler(data map[string][]byte) http.Handler {
	files := http.FileServer(http.Dir("web"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if name, ok := strings.CutPrefix(path.Clean(r.URL.Path), "/data/"); ok {
			if b, ok := data[name]; ok {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Cache-Control", "no-store")
				w.Write(b)
				return
			}
		}
		w.Header().Set("Cache-Control", "no-store")
		files.ServeHTTP(w, r)
	})
}

// ---------------------------------------------------------------------------
// Shared definitions

// Run describes where the data comes from.
type Run struct {
	Date      string     `json:"date"`   // results/<date>, the cross-language campaign
	GoDate    string     `json:"goDate"` // go/results/<date>, the Go campaign
	Env       [][]string `json:"env"`
	Generated string     `json:"generated"`
}

// Lang is a language, in the order of every table.
type Lang struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var langs = []Lang{
	{"go", "Go"}, {"typescript", "TypeScript"}, {"python", "Python"}, {"rust", "Rust"},
	{"java", "Java"}, {"csharp", "C#"}, {"c", "C"}, {"cpp", "C++"}, {"php", "PHP"},
}

func langOrder(id string) int {
	for i, l := range langs {
		if l.ID == id {
			return i
		}
	}
	return len(langs)
}

// Lib is one library of one language.
type Lib struct {
	ID      string `json:"id"` // lang/lib, as in the raw records
	Name    string `json:"name"`
	Lang    string `json:"lang"`
	Kind    string `json:"kind"` // library, floor, xslt (official XSLT run through Saxon), self
	Pkg     string `json:"pkg,omitempty"`
	Version string `json:"version,omitempty"`
}

// The ZeRechnung Go package and its second round at GOMAXPROCS=1.
const (
	goLib     = "ZeRechnung"
	goOneCore = "ZeRechnung, GOMAXPROCS=1"
)

// displayNames are the names the READMEs use, when they differ from the
// identifiers of the raw records.
var displayNames = map[string]string{
	"go/" + goOneCore:                 "ZeRechnung, 1 core",
	"c/mupdf":                         "MuPDF",
	"c/qpdf-c":                        "qpdf (C API)",
	"cpp/kura":                        "Kura",
	"cpp/podofo":                      "PoDoFo",
	"cpp/saxonc-he":                   "SaxonC-HE + official XSLT",
	"cpp/xerces-c":                    "Xerces-C++",
	"csharp/antiphon":                 "Antiphon",
	"csharp/balsoft":                  "Balsoft",
	"csharp/en16931":                  "En16931",
	"csharp/facturxdotnet":            "FacturXDotNet",
	"csharp/international-einvoicing": "International.EInvoicing",
	"csharp/officeimo":                "OfficeIMO",
	"csharp/pax-xrechnung":            "pax.XRechnung.NET",
	"csharp/pdfpig":                   "PdfPig",
	"csharp/pdfsharp":                 "PDFsharp",
	"csharp/saxoncs":                  "SaxonCS-HE + official XSLT",
	"csharp/securibox":                "Securibox",
	"csharp/xsd":                      "System.Xml.Schema",
	"csharp/zugferd-csharp":           "ZUGFeRD-csharp",
	"java/kosit-validator":            "KoSIT validator",
	"java/mustang":                    "Mustang",
	"java/paniolo":                    "Paniolo",
	"java/pdfbox":                     "PDFBox",
	"java/saxon-schematron":           "Saxon-HE + official XSLT",
	"python/saxonche-schematron":      "saxonche + official XSLT",
}

// floors are the generic PDF or XML libraries wired in a few lines, which give
// the minimum cost of an operation in their language (the COVERAGE.md file
// of each language lists them).
var floors = map[string]bool{
	"c/libxml2": true, "c/libxslt": true, "c/mupdf": true, "c/qpdf-c": true,
	"cpp/qpdf": true, "cpp/podofo": true, "cpp/poppler-cpp": true, "cpp/xerces-c": true, "cpp/pugixml": true,
	"csharp/xsd": true, "csharp/pdfpig": true, "csharp/pdfsharp": true,
	"java/pdfbox":          true,
	"php/smalot/pdfparser": true, "php/ext-dom": true,
	"rust/lopdf": true,
}

// xslt are the official CEN and KoSIT XSLT run through Saxon, in each
// language that has a Saxon: the reference implementation of the rules.
var xslt = map[string]bool{
	"java/saxon-schematron": true, "csharp/saxoncs": true, "cpp/saxonc-he": true, "python/saxonche-schematron": true,
}

func newLib(lang, lib, pkg, version string) Lib {
	id := lang + "/" + lib
	l := Lib{ID: id, Name: lib, Lang: lang, Kind: "library", Pkg: pkg, Version: version}
	if n, ok := displayNames[id]; ok {
		l.Name = n
	}
	switch {
	case lang == "go" && strings.HasPrefix(lib, goLib):
		l.Kind = "self"
	case floors[id]:
		l.Kind = "floor"
	case xslt[id]:
		l.Kind = "xslt"
	}
	if l.Pkg == lib {
		l.Pkg = ""
	}
	return l
}

func sortLibs(libs []Lib) {
	sort.Slice(libs, func(i, j int) bool {
		a, b := libs[i], libs[j]
		if a.Lang != b.Lang {
			return langOrder(a.Lang) < langOrder(b.Lang)
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
}

// loadEnv reads env.txt ("key: value" lines).
func loadEnv(p string) ([][]string, error) {
	var out [][]string
	err := eachLine(p, func(line string) {
		k, v, ok := strings.Cut(line, ": ")
		if ok {
			out = append(out, []string{strings.TrimSpace(k), strings.TrimSpace(v)})
		}
	})
	return out, err
}

func eachLine(p string, fn func(string)) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64*1024), 8*1024*1024)
	for sc.Scan() {
		fn(sc.Text())
	}
	return sc.Err()
}

// eachJSON decodes every non-empty line of a JSONL file into a fresh T.
func eachJSON[T any](p string, fn func(T)) error {
	return eachLine(p, func(line string) {
		if strings.TrimSpace(line) == "" {
			return
		}
		var v T
		if json.Unmarshal([]byte(line), &v) == nil {
			fn(v)
		}
	})
}

func median(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	v = slices.Clone(v)
	slices.Sort(v)
	n := len(v)
	if n%2 == 1 {
		return v[n/2]
	}
	return (v[n/2-1] + v[n/2]) / 2
}

// ---------------------------------------------------------------------------
// CPU time and memory

// Op is a measured operation and its benchmarks.
type Op struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Benches []Bench `json:"benches"`
}

// Bench is one benchmark: an operation on one input file.
type Bench struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

var ops = []Op{
	{"extract", "Extraction", []Bench{{"Extract_Small", "EN 16931, 149 KB"}, {"Extract_Large", "EN 16931, 1.1 MB"}, {"Extract_Minimum", "MINIMUM, 89 KB"}, {"Extract_Extended", "EXTENDED, 157 KB"}}},
	{"parse", "Parsing", []Bench{{"Parse_Small", "13 KB"}, {"Parse_Lines", "19 KB"}, {"Parse_Large", "1.3 MB"}}},
	{"validate", "Validation", []Bench{{"ParseValidate_Small", "13 KB"}, {"ParseValidate_Lines", "19 KB"}, {"ParseValidate_Large", "1.3 MB"}}},
	{"write", "Writing", []Bench{{"Write_Small", "13 KB"}, {"Write_Lines", "19 KB"}}},
	{"embed", "Embedding", []Bench{{"Embed_Small", "xref-stream host"}, {"EmbedClassic_Small", "classic xref host"}}},
	{"check", "PDF check", []Bench{{"CheckPDF_Small", "149 KB"}, {"CheckPDF_Large", "1.1 MB"}}},
}

// Perf is data/perf.json.
type Perf struct {
	Run    Run           `json:"run"`
	Langs  []Lang        `json:"langs"`
	Ops    []Op          `json:"ops"`
	Libs   []Lib         `json:"libs"`
	Rows   []Row         `json:"rows"`
	ByLang []LangSummary `json:"byLang"`
}

// Row is one library on one benchmark in the resources pass: a fresh process
// that loads the library and runs this benchmark only.
type Row struct {
	Lib    string  `json:"lib"`
	Bench  string  `json:"bench"`
	Status string  `json:"status"`           // ok, fail
	Wall   float64 `json:"wall,omitempty"`   // median wall-clock ns per call
	CPU    float64 `json:"cpu,omitempty"`    // median CPU ns per call, all threads
	RSS    float64 `json:"rss,omitempty"`    // peak resident set size of the process, bytes
	First  float64 `json:"first,omitempty"`  // first call in the process, ns
	Bytes  float64 `json:"bytes,omitempty"`  // bytes allocated per call
	Allocs float64 `json:"allocs,omitempty"` // allocations per call (0: not counted)
	Lower  bool    `json:"lower,omitempty"`  // a GC ran during the count: bytes is a minimum
	Child  bool    `json:"child,omitempty"`  // child processes not counted
	Note   string  `json:"note,omitempty"`
}

// LangSummary is the spread of peak memory and CPU use over all the processes
// of one language, as in the "Resources per language" table of tables.md.
type LangSummary struct {
	Lang      string  `json:"lang"`
	Libraries int     `json:"libraries"`
	Processes int     `json:"processes"`
	RSSMin    float64 `json:"rssMin"`
	RSSMedian float64 `json:"rssMedian"`
	RSSMax    float64 `json:"rssMax"`
	UseMin    float64 `json:"useMin"` // CPU / wall time
	UseMedian float64 `json:"useMedian"`
	UseMax    float64 `json:"useMax"`
}

type rawRecord struct {
	Lang, Lib, Pkg, Version, Bench, Status, Note string
	NsPerOp                                      []float64 `json:"ns_per_op"`
	CPUNsPerOp                                   []float64 `json:"cpu_ns_per_op"`
	FirstCallNs                                  float64   `json:"first_call_ns"`
	MaxRSSBytes                                  float64   `json:"max_rss_bytes"`
	BytesPerOp                                   float64   `json:"bytes_per_op"`
	AllocsPerOp                                  float64   `json:"allocs_per_op"`
	// alloc pass
	Bytes  *float64 `json:"bytes"`
	Allocs *float64 `json:"allocs"`
	Gcs    float64  `json:"gcs"`
}

func buildPerf(dir string, run Run) (*Perf, error) {
	files, _ := filepath.Glob(filepath.Join(dir, "resources_*.jsonl"))
	if len(files) == 0 {
		return nil, fmt.Errorf("no resources_*.jsonl in %s: run ./run.sh resources", dir)
	}
	type allocKey struct{ lib, bench string }
	allocs := map[allocKey]rawRecord{}
	afiles, _ := filepath.Glob(filepath.Join(dir, "alloc_*.jsonl"))
	for _, f := range afiles {
		err := eachJSON(f, func(r rawRecord) {
			if r.Status == "ok" && r.Bytes != nil {
				allocs[allocKey{r.Lang + "/" + r.Lib, r.Bench}] = r
			}
		})
		if err != nil {
			return nil, err
		}
	}

	p := &Perf{Run: run, Langs: langs, Ops: ops}
	p.Run.Generated = time.Now().UTC().Format(time.RFC3339)
	libs := map[string]Lib{}
	for _, f := range files {
		err := eachJSON(f, func(r rawRecord) {
			if r.Status == "unsupported" {
				return
			}
			id := r.Lang + "/" + r.Lib
			if _, ok := libs[id]; !ok {
				version := r.Version
				if r.Lang == "go" {
					version = "" // the Go version, not the package's
				}
				libs[id] = newLib(r.Lang, r.Lib, r.Pkg, version)
			}
			row := Row{Lib: id, Bench: r.Bench, Status: r.Status, Note: r.Note,
				Child: strings.Contains(r.Note, "child process")}
			if r.Status == "ok" && len(r.NsPerOp) > 0 && len(r.CPUNsPerOp) > 0 {
				row.Wall, row.CPU = median(r.NsPerOp), median(r.CPUNsPerOp)
				row.RSS, row.First = r.MaxRSSBytes, r.FirstCallNs
			} else if r.Status == "ok" {
				row.Status = "fail"
			}
			if r.BytesPerOp > 0 {
				row.Bytes, row.Allocs = r.BytesPerOp, r.AllocsPerOp
			} else if a, ok := allocs[allocKey{id, r.Bench}]; ok && row.Status == "ok" {
				row.Bytes = *a.Bytes
				if a.Allocs != nil {
					row.Allocs = *a.Allocs
				}
				row.Lower = a.Gcs > 0
			}
			p.Rows = append(p.Rows, row)
		})
		if err != nil {
			return nil, err
		}
	}
	for _, l := range libs {
		p.Libs = append(p.Libs, l)
	}
	sortLibs(p.Libs)
	sort.SliceStable(p.Rows, func(i, j int) bool {
		if p.Rows[i].Lib != p.Rows[j].Lib {
			return p.Rows[i].Lib < p.Rows[j].Lib
		}
		return p.Rows[i].Bench < p.Rows[j].Bench
	})
	p.ByLang = summarize(p.Rows, libs)
	return p, nil
}

// summarize mirrors printResourcesByLang in report/: the Go series at
// GOMAXPROCS=1 and the libraries that spawn child processes are left out.
func summarize(rows []Row, libs map[string]Lib) []LangSummary {
	rss := map[string][]float64{}
	use := map[string][]float64{}
	names := map[string]map[string]bool{}
	for _, r := range rows {
		l := libs[r.Lib]
		if r.Status != "ok" || r.RSS <= 0 || r.Child || r.Lib == "go/"+goOneCore {
			continue
		}
		rss[l.Lang] = append(rss[l.Lang], r.RSS)
		if r.Wall > 0 {
			use[l.Lang] = append(use[l.Lang], r.CPU/r.Wall)
		}
		if names[l.Lang] == nil {
			names[l.Lang] = map[string]bool{}
		}
		names[l.Lang][r.Lib] = true
	}
	var out []LangSummary
	for _, l := range langs {
		v := rss[l.ID]
		if len(v) == 0 {
			continue
		}
		out = append(out, LangSummary{
			Lang: l.ID, Libraries: len(names[l.ID]), Processes: len(v),
			RSSMin: slices.Min(v), RSSMedian: median(v), RSSMax: slices.Max(v),
			UseMin: slices.Min(use[l.ID]), UseMedian: median(use[l.ID]), UseMax: slices.Max(use[l.ID]),
		})
	}
	return out
}

// ---------------------------------------------------------------------------
// Compliance

// Compliance is data/compliance.json.
type Compliance struct {
	Run       Run             `json:"run"`
	Langs     []Lang          `json:"langs"`
	Files     []string        `json:"files"`     // the valid CII XML files of the corpus
	Reference Reference       `json:"reference"` // what the official schematrons flag
	Libs      []SurveyLib     `json:"libs"`
	Standards json.RawMessage `json:"standards"`
	Countries []Country       `json:"countries"` // confidence per country (conformance cases)
}

// Reference is the result of the official CEN and KoSIT XSLT on the corpus.
type Reference struct {
	Lib   string `json:"lib"`
	Flags []Flag `json:"flags"`
}

// SurveyLib is one library's accuracy on the corpus.
type SurveyLib struct {
	Lib
	Validate *Survey `json:"validate,omitempty"`
	Parse    *Survey `json:"parse,omitempty"`
	Extract  *Survey `json:"extract,omitempty"`
}

// Survey is one operation run on the whole corpus. For validate, Ok counts
// the valid files the library flags, and Flags lists them.
type Survey struct {
	Ok    int    `json:"ok"`
	Total int    `json:"total"`
	Flags []Flag `json:"flags,omitempty"`
	Note  string `json:"note,omitempty"`
}

// Flag is one valid file a validator rejects: by rules, or by an exception.
type Flag struct {
	File   string   `json:"file"`
	Errors int      `json:"errors,omitempty"`
	Rules  []string `json:"rules,omitempty"`
	Error  string   `json:"error,omitempty"`
}

// goLibs are the modules of go/libs, named after the module they measure
// (and naming the raw outputs), their display name (as in go/report) and the
// module path of the library.
var goLibs = []struct{ dir, name, module string }{
	{"speedata-einvoice", "speedata", "github.com/speedata/einvoice"},
	{"andeedotnet-go-xinvoice", "andee", "github.com/andeedotnet/go-xinvoice"},
	{"dotwavehq-go-einvoice", "dotwave", "github.com/dotwavehq/go-einvoice"},
	{"invopop-gobl.de.xinvoice", "gobl.de", "github.com/invopop/gobl.de.xinvoice"},
	{"angelodlfrtr-go-invoice-generator", "angelo", "github.com/angelodlfrtr/go-invoice-generator"},
	{"papierleicht-go-e-rechnung-validator", "papierleicht", "github.com/papierleicht/go-e-rechnung-validator"},
	{"dcnws-invoicekit.go", "DCNWS", "github.com/dcnws/invoicekit.go"},
	{"mgilbir-pdf0", "pdf0", "github.com/mgilbir/pdf0"},
	{"gotenberg-litho", "litho", "github.com/gotenberg/litho"},
	{"cyprienbrisset-kanjo", "kanjo", "github.com/cyprienbrisset/kanjo"},
	{"tahayassine-facturx-go", "taha", "github.com/tahayassine/facturx-go"},
	{"solarpush-fx", "fX", "github.com/solarpush/fX"},
	{"marlinkuhn-gopdfattach", "gopdfattach", "github.com/MarlinKuhn/gopdfattach"},
	{"pdfcpu-pdfcpu", "pdfcpu", "github.com/pdfcpu/pdfcpu"},
}

// zerechnungModule is the module of go/libs that measures ZeRechnung: it
// names its raw outputs in both campaigns.
const zerechnungModule = "ritoon-zerechnung"

// referenceLibs are tried in order to find the official XSLT run.
var referenceLibs = []string{"java/saxon-schematron", "python/saxonche-schematron", "csharp/saxoncs", "cpp/saxonc-he"}

var (
	goSurvey = regexp.MustCompile(`^SURVEY (validate|extract|parse): (\d+)/(\d+)`)
	goFlag   = regexp.MustCompile(`^FLAG (\S+): (\d+) errors, err=(.*)$`)
)

func buildCompliance(root, resDir, goDir string, run Run) (*Compliance, error) {
	c := &Compliance{Run: run, Langs: langs}
	c.Run.Generated = time.Now().UTC().Format(time.RFC3339)
	byID := map[string]*SurveyLib{}
	get := func(l Lib) *SurveyLib {
		if s, ok := byID[l.ID]; ok {
			return s
		}
		s := &SurveyLib{Lib: l}
		byID[l.ID] = s
		return s
	}

	// ZeRechnung, as measured by the cross-language campaign.
	zr, err := loadGoSurvey(filepath.Join(resDir, "survey_"+zerechnungModule+".txt"))
	if err != nil {
		return nil, err
	}
	s := get(newLib("go", goLib, "github.com/ritoon/zerechnung", ""))
	s.Validate, s.Parse, s.Extract = zr["validate"], zr["parse"], zr["extract"]

	// The other Go libraries.
	for _, m := range goLibs {
		sv, err := loadGoSurvey(filepath.Join(goDir, "survey_"+m.dir+".txt"))
		if errors.Is(err, fs.ErrNotExist) {
			continue
		} else if err != nil {
			return nil, err
		}
		l := newLib("go", m.dir, m.module, goModVersion(filepath.Join(root, "go", "libs", m.dir, "go.mod"), m.module))
		l.Name = m.name
		s := get(l)
		s.Validate, s.Parse, s.Extract = sv["validate"], sv["parse"], sv["extract"]
	}

	// The other languages.
	files, _ := filepath.Glob(filepath.Join(resDir, "survey_*.jsonl"))
	for _, f := range files {
		err := eachJSON(f, func(r struct {
			Lang, Lib, Pkg, Version, Op, Note string
			Ok, Total                         int
			Failures                          []Flag
		}) {
			s := get(newLib(r.Lang, r.Lib, r.Pkg, r.Version))
			sv := &Survey{Ok: r.Ok, Total: r.Total}
			switch r.Op {
			case "validate":
				sv.Note = r.Note
				for _, fl := range r.Failures {
					fl.File = path.Base(fl.File)
					fl.Error = truncate(fl.Error, 240)
					sv.Flags = append(sv.Flags, fl)
				}
				s.Validate = sv
			case "parse":
				s.Parse = sv
			case "extract":
				s.Extract = sv
			}
		})
		if err != nil {
			return nil, err
		}
	}

	// The valid files: every file some validator flags (fX and invoice4j
	// flag all 26).
	seen := map[string]bool{}
	for _, s := range byID {
		if s.Validate == nil {
			continue
		}
		for _, f := range s.Validate.Flags {
			seen[f.File] = true
		}
	}
	for f := range seen {
		c.Files = append(c.Files, f)
	}
	sort.Strings(c.Files)

	for _, id := range referenceLibs {
		if s, ok := byID[id]; ok && s.Validate != nil {
			c.Reference = Reference{Lib: id, Flags: s.Validate.Flags}
			break
		}
	}
	if c.Reference.Lib == "" {
		return nil, errors.New("no official XSLT survey to use as the reference")
	}

	for _, s := range byID {
		c.Libs = append(c.Libs, *s)
	}
	sort.Slice(c.Libs, func(i, j int) bool {
		a, b := c.Libs[i].Lib, c.Libs[j].Lib
		if a.Lang != b.Lang {
			return langOrder(a.Lang) < langOrder(b.Lang)
		}
		if (a.Kind == "self") != (b.Kind == "self") {
			return a.Kind == "self"
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})

	if c.Countries, err = buildConfidence(filepath.Join(root, "cases"), resDir, goDir); err != nil {
		return nil, err
	}

	std, err := os.ReadFile("standards.json")
	if err != nil {
		return nil, err
	}
	if !json.Valid(std) {
		return nil, errors.New("standards.json is not valid JSON")
	}
	c.Standards = std
	return c, nil
}

// loadGoSurvey reads the survey_<module>.txt of a Go module.
func loadGoSurvey(p string) (map[string]*Survey, error) {
	out := map[string]*Survey{}
	var flags []Flag
	err := eachLine(p, func(line string) {
		if m := goFlag.FindStringSubmatch(line); m != nil {
			n, _ := strconv.Atoi(m[2])
			fl := Flag{File: path.Base(m[1]), Errors: n}
			if m[3] != "<nil>" {
				fl.Error = truncate(m[3], 240)
			}
			flags = append(flags, fl)
		}
		if m := goSurvey.FindStringSubmatch(line); m != nil {
			ok, _ := strconv.Atoi(m[2])
			total, _ := strconv.Atoi(m[3])
			out[m[1]] = &Survey{Ok: ok, Total: total}
		}
	})
	if v := out["validate"]; v != nil {
		v.Flags = flags
	}
	return out, err
}

// goModVersion returns the version of module that a go.mod requires.
func goModVersion(gomod, module string) string {
	version := ""
	eachLine(gomod, func(line string) {
		f := strings.Fields(strings.TrimPrefix(strings.TrimSpace(line), "require "))
		if version == "" && len(f) >= 2 && strings.EqualFold(f[0], module) {
			version = f[1]
		}
	})
	return version
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if r := []rune(s); len(r) > n {
		return string(r[:n-1]) + "…"
	}
	return s
}

// ---------------------------------------------------------------------------
// Confidence per country

// Country is the confidence score of every library rated for one country,
// computed by the cases package exactly as the report tools do.
type Country struct {
	Code      string          `json:"code"`
	Name      string          `json:"name"`
	Reference string          `json:"reference"`
	Parts     []CountryPart   `json:"parts"`
	Valid     int             `json:"valid"`   // valid cases, all parts together
	Invalid   int             `json:"invalid"` // invalid cases
	Rules     int             `json:"rules"`   // expected rules
	Libs      []ConfidenceLib `json:"libs"`
}

// CountryPart is one suite in one syntax: the unit a score is computed on.
type CountryPart struct {
	Suite  string `json:"suite"`
	Syntax string `json:"syntax"`
	Title  string `json:"title"`
}

// ConfidenceLib is one library's result for a country.
type ConfidenceLib struct {
	Lib
	Score     float64     `json:"score"` // mean of the parts, out of 100
	S         float64     `json:"s"`     // valid cases accepted, all parts together
	D         float64     `json:"d"`     // detection of the expected rules, all parts together
	Accepted  int         `json:"accepted"`
	Rejected  int         `json:"rejected"`
	Exact     int         `json:"exact"`
	Tolerated int         `json:"tolerated,omitempty"`
	Parts     []PartScore `json:"parts"`
}

// PartScore is a library's result on one part; Run is false when the
// library was not run on its suite (the part then counts 0).
type PartScore struct {
	Score float64 `json:"score"`
	S     float64 `json:"s"`
	D     float64 `json:"d"`
	Run   bool    `json:"run"`
}

// buildConfidence scores the cases_*.jsonl files of both campaigns: from the
// cross-language one, every library but the Go ones (ZeRechnung aside),
// which the Go campaign rates.
func buildConfidence(casesDir, resDir, goDir string) ([]Country, error) {
	idx, err := cases.LoadIndex(casesDir)
	if err != nil {
		return nil, err
	}
	cross, err := cases.LoadRuns(resDir)
	if err != nil {
		return nil, err
	}
	gor, err := cases.LoadRuns(goDir)
	if err != nil {
		return nil, err
	}
	var runs []cases.Run
	for _, r := range cross {
		if r.Lang != "go" || r.Lib == zerechnungModule {
			runs = append(runs, r)
		}
	}
	for _, r := range gor {
		if r.Lib != zerechnungModule {
			runs = append(runs, r)
		}
	}
	if len(runs) == 0 {
		return nil, nil
	}
	goNames := map[string]string{}
	goModules := map[string]string{}
	for _, m := range goLibs {
		goNames[m.dir], goModules[m.dir] = m.name, m.module
	}

	var out []Country
	for _, code := range idx.Countries() {
		ct := Country{Code: code, Name: code}
		if n, ok := cases.CountryNames[code]; ok {
			ct.Name = n
		}
		parts := idx.Parts(code)
		for _, p := range parts {
			title := p.Suite
			if s, ok := idx.Suite(p.Suite); ok {
				title = s.Title
				if ct.Reference == "" {
					ref := s.RefLang + "/" + s.RefLib
					ct.Reference = ref
					if d, ok := cases.References[ref]; ok {
						ct.Reference = d
					}
				}
			}
			ct.Parts = append(ct.Parts, CountryPart{p.Suite, p.Syntax, title})
		}
		for _, l := range cases.Libraries(runs) {
			cs := idx.ForCountry(l, code)
			if !cs.Run {
				continue
			}
			var lib Lib
			switch {
			case l.Lang == "go" && l.Lib == zerechnungModule:
				lib = newLib("go", goLib, "github.com/ritoon/zerechnung", "")
			case l.Lang == "go":
				lib = newLib("go", l.Lib, goModules[l.Lib], l.Version)
				if n, ok := goNames[l.Lib]; ok {
					lib.Name = n
				}
			default:
				lib = newLib(l.Lang, l.Lib, l.Pkg, l.Version)
			}
			t := cs.Total
			cl := ConfidenceLib{Lib: lib, Score: cs.Note, S: t.S(), D: t.D(),
				Accepted: t.Accepted, Rejected: t.Rejected, Exact: t.Exact, Tolerated: t.Tolerated}
			for _, p := range parts {
				ps, ok := cs.Parts[p]
				if !ok {
					cl.Parts = append(cl.Parts, PartScore{})
					continue
				}
				cl.Parts = append(cl.Parts, PartScore{Score: ps.Note(), S: ps.S(), D: ps.D(), Run: true})
			}
			if ct.Valid == 0 {
				ct.Valid, ct.Invalid, ct.Rules = t.Valid, t.Invalid, len(t.Credit)
			}
			ct.Libs = append(ct.Libs, cl)
		}
		sort.SliceStable(ct.Libs, func(i, j int) bool {
			a, b := ct.Libs[i], ct.Libs[j]
			if a.Score != b.Score {
				return a.Score > b.Score
			}
			if a.Lang != b.Lang {
				return langOrder(a.Lang) < langOrder(b.Lang)
			}
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)
		})
		out = append(out, ct)
	}
	return out, nil
}
