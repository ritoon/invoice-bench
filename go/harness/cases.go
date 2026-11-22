package harness

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"
	"testing"
	"time"
)

// CasesDir is the directory of the conformance cases: CASES_DIR, or the
// cases directory at the root of the repository.
func CasesDir() string {
	if v := os.Getenv("CASES_DIR"); v != "" {
		return v
	}
	_, f, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(f), "..", "..", "cases")
}

// caseTimeout bounds one validation: a library stuck on a malformed case
// is recorded as an error of that case, and the run goes on.
const caseTimeout = 60 * time.Second

type caseResult struct {
	Case   string   `json:"case"`
	Errors int      `json:"errors"`
	Rules  []string `json:"rules"`
	Error  string   `json:"error,omitempty"`
}

type casesLine struct {
	Lang    string       `json:"lang"`
	Lib     string       `json:"lib"`
	Pkg     string       `json:"pkg"`
	Version string       `json:"version"`
	Suite   string       `json:"suite"`
	Results []caseResult `json:"results"`
	Note    string       `json:"note"`
}

// Cases validates the conformance cases of cases/cases.tsv (see
// cases/README.md) and writes one JSON line per suite to the file CASES_OUT,
// in the format of the harnesses of the other languages: lang "go", lib the name of the
// bench module, pkg the module path of the library (its version read from the
// build information), suite, results (case, errors, rules, error) and note.
// validate returns the number of errors and their rule identifiers, or
// messages that contain them. A panic, an error or a call that does not end
// within a minute is recorded as an error of the case. Without CASES_OUT the
// test is skipped.
func Cases(t *testing.T, lib, pkg, note string, validate func([]byte) (int, []string, error)) {
	out := os.Getenv("CASES_OUT")
	if out == "" {
		t.Skip("CASES_OUT not set (run ./run.sh cases)")
	}
	dir := CasesDir()
	data, err := os.ReadFile(filepath.Join(dir, "cases.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	var suites []string
	files := map[string][]string{}
	for i, line := range strings.Split(string(data), "\n") {
		f := strings.Split(line, "\t")
		if i == 0 || len(f) < 2 || strings.TrimSpace(f[0]) == "" {
			continue
		}
		if files[f[0]] == nil {
			suites = append(suites, f[0])
		}
		files[f[0]] = append(files[f[0]], f[1])
	}
	w, err := os.Create(out)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	for _, s := range suites {
		line := casesLine{Lang: "go", Lib: lib, Pkg: pkg, Version: version(pkg), Suite: s, Note: note}
		rejected := 0
		for _, f := range files[s] {
			x, err := os.ReadFile(filepath.Join(dir, f))
			if err != nil {
				t.Fatal(err)
			}
			r := caseResult{Case: f, Rules: []string{}}
			n, rules, err := guarded(validate, x)
			if err != nil {
				r.Error = shorten(err.Error())
			} else {
				r.Errors = n
				r.Rules = append(r.Rules, rules...)
				slices.Sort(r.Rules)
				r.Rules = slices.Compact(r.Rules)
			}
			if r.Errors > 0 || r.Error != "" {
				rejected++
			}
			line.Results = append(line.Results, r)
		}
		t.Logf("CASES %s: %d/%d rejected", s, rejected, len(files[s]))
		if err := enc.Encode(line); err != nil {
			t.Fatal(err)
		}
	}
}

// guarded runs validate on x in its own goroutine, turning a panic into an
// error and giving up after caseTimeout.
func guarded(validate func([]byte) (int, []string, error), x []byte) (int, []string, error) {
	type res struct {
		n     int
		rules []string
		err   error
	}
	done := make(chan res, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				done <- res{err: fmt.Errorf("panic: %v", p)}
			}
		}()
		n, rules, err := validate(x)
		done <- res{n, rules, err}
	}()
	select {
	case r := <-done:
		return r.n, r.rules, r.err
	case <-time.After(caseTimeout):
		return 0, nil, fmt.Errorf("no result after %v", caseTimeout)
	}
}

// version is the version of module pkg in the test binary: its replacement's
// when it is replaced by another version, "local" when by a directory.
func version(pkg string) string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "?"
	}
	for _, d := range bi.Deps {
		if d.Path != pkg {
			continue
		}
		if d.Replace != nil {
			if d.Replace.Version != "" {
				return d.Replace.Version
			}
			return "local"
		}
		return d.Version
	}
	return "?"
}

func shorten(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if r := []rune(s); len(r) > 400 {
		return string(r[:400]) + "…"
	}
	return s
}
