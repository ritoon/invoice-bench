// Command tool maintains cases.tsv (see ../README.md).
//
//	go run . index                  list the documents of each suite in ../cases.tsv
//	go run . confirm RESULTS_DIR    fill in the verdicts and the expected rules
//	                                from the outcomes of the reference validators
//
// index keeps the verdict of a document whose content has not changed (same
// SHA-256), and marks the others "-" (to confirm). confirm reads the
// cases_<lang>.jsonl files of RESULTS_DIR, takes for each suite the run of its
// reference (suites.tsv) and records, for each document, the rules the
// reference raises as errors: none, valid; some, invalid. A document the
// reference cannot judge (exception, processing error, no scenario) is
// excluded, with the reason in the note column.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ritoon/invoice-bench/cases"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	var err error
	switch os.Args[1] {
	case "index":
		err = index("..")
	case "confirm":
		if len(os.Args) != 3 {
			usage()
		}
		err = confirm("..", os.Args[2])
	default:
		usage()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "tool:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: go run . index | confirm RESULTS_DIR")
	os.Exit(2)
}

func index(dir string) error {
	idx, err := cases.LoadIndex(dir)
	if err != nil {
		return err
	}
	old := map[string]cases.Case{}
	for _, c := range idx.Cases {
		old[c.File] = c
	}
	var out []cases.Case
	for _, s := range idx.Suites {
		intent, err := intents(filepath.Join(dir, s.Name, "intent.tsv"))
		if err != nil {
			return err
		}
		files, _ := filepath.Glob(filepath.Join(dir, s.Name, "*.xml"))
		sort.Strings(files)
		for _, g := range s.External {
			ext, _ := filepath.Glob(filepath.Join(dir, g))
			if len(ext) == 0 {
				return fmt.Errorf("suite %s: no file matches %s (run ./fetch.sh)", s.Name, g)
			}
			sort.Strings(ext)
			files = append(files, ext...)
		}
		n := 0
		for _, f := range files {
			data, err := os.ReadFile(f)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(dir, f)
			rel = filepath.ToSlash(rel)
			syn := syntax(data)
			if syn == "" {
				return fmt.Errorf("%s: neither a CII nor a UBL invoice", rel)
			}
			sum := sha256.Sum256(data)
			c := cases.Case{Suite: s.Name, File: rel, Syntax: syn, SHA: hex.EncodeToString(sum[:6]),
				Verdict: cases.Pending, Intent: intent[filepath.Base(f)]}
			if o, ok := old[rel]; ok && o.SHA == c.SHA {
				c.Verdict, c.Rules, c.Note = o.Verdict, o.Rules, o.Note
			}
			out = append(out, c)
			n++
		}
		fmt.Printf("== %s: %d documents\n", s.Name, n)
	}
	idx.Cases = out
	return idx.WriteCases()
}

// intents reads the expected.tsv of the mutation tests: file, then the rules
// the case was written to break, space separated.
func intents(path string) (map[string][]string, error) {
	m := map[string][]string{}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return m, nil
	}
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		f, rules, ok := strings.Cut(line, "\t")
		if ok {
			m[f] = strings.Fields(rules)
		}
	}
	return m, nil
}

// syntax tells a CII invoice from a UBL one by its root element.
func syntax(doc []byte) string {
	i := 0
	for {
		j := bytes.IndexByte(doc[i:], '<')
		if j < 0 {
			return ""
		}
		i += j + 1
		if i < len(doc) && (doc[i] == '?' || doc[i] == '!') {
			continue
		}
		end := bytes.IndexAny(doc[i:], " \t\r\n/>")
		if end < 0 {
			return ""
		}
		name := string(doc[i : i+end])
		if k := strings.IndexByte(name, ':'); k >= 0 {
			name = name[k+1:]
		}
		switch name {
		case "CrossIndustryInvoice":
			return "cii"
		case "Invoice", "CreditNote":
			return "ubl"
		}
		return ""
	}
}

func confirm(dir, results string) error {
	idx, err := cases.LoadIndex(dir)
	if err != nil {
		return err
	}
	runs, err := cases.LoadRuns(results)
	if err != nil {
		return err
	}
	ref := map[string]map[string]cases.Outcome{} // suite → case → outcome of its reference
	for _, s := range idx.Suites {
		for _, r := range runs {
			if r.Lang == s.RefLang && r.Lib == s.RefLib && r.Suite == s.Name {
				m := map[string]cases.Outcome{}
				for _, o := range r.Results {
					m[o.Case] = o
				}
				ref[s.Name] = m
			}
		}
		if ref[s.Name] == nil {
			return fmt.Errorf("suite %s: no outcome of its reference %s %s in %s", s.Name, s.RefLang, s.RefLib, results)
		}
	}
	count := map[string]int{}
	for i, c := range idx.Cases {
		o, ok := ref[c.Suite][c.File]
		if !ok {
			return fmt.Errorf("%s: no outcome of the reference (run it on the current cases.tsv)", c.File)
		}
		c.Rules, c.Note = nil, ""
		switch {
		case o.Error != "":
			c.Verdict, c.Note = cases.Excluded, "reference: "+o.Error
		default:
			var rules []string
			seen := map[string]bool{}
			for _, r := range o.Rules {
				if strings.HasPrefix(r, "PROCESSING:") || strings.HasPrefix(r, "REJECT:") {
					c.Note = "reference: " + strings.Join(strings.Fields(r), " ")
					continue
				}
				if !seen[r] {
					seen[r] = true
					rules = append(rules, r)
				}
			}
			sort.Strings(rules)
			switch {
			case c.Note != "":
				c.Verdict = cases.Excluded
			case len(rules) == 0:
				c.Verdict = cases.Valid
			default:
				c.Verdict, c.Rules = cases.Invalid, rules
			}
		}
		idx.Cases[i] = c
		count[c.Suite+" "+c.Verdict]++
	}
	keys := make([]string, 0, len(count))
	for k := range count {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("== %s: %d\n", k, count[k])
	}
	return idx.WriteCases()
}
