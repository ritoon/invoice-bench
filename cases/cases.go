// Package cases reads the conformance cases of invoice-bench (suites.tsv and
// cases.tsv), the outcomes of the harnesses on them (cases_*.jsonl) and
// computes the confidence score of each library by suite, syntax and country.
// See README.md for the format and the formula.
package cases

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Verdicts of a case.
const (
	Valid    = "valid"    // the reference raises no error
	Invalid  = "invalid"  // the reference raises at least one error
	Pending  = "-"        // not confirmed by the reference yet
	Excluded = "excluded" // the reference could not judge it (see the note)
)

// Suite is a set of cases checked by one reference validator.
type Suite struct {
	Name      string
	Syntaxes  []string // cii, ubl
	Countries []string // ISO 3166 codes of the countries the suite counts for
	Title     string
	RefLang   string // the harness and library that judge the cases
	RefLib    string
	External  []string // globs of fetched files, relative to the cases directory
	// Unchecked lists rules of the norm that the reference does not check
	// (BR-CO-25 is in none of the CEN artifacts): a valid case that a
	// library rejects under these rules alone is not held against it.
	Unchecked []string
}

// Case is one document and its expected outcome.
type Case struct {
	Suite   string
	File    string // relative to the cases directory, the key of the outcomes
	Syntax  string
	SHA     string // first 12 hex digits of the SHA-256 of the document
	Verdict string
	Rules   []string // rules the reference raises as errors, as it spells them
	Intent  []string // rules the case was written to break (warnings included)
	Note    string
}

// Index is the content of a cases directory.
type Index struct {
	Dir    string
	Suites []Suite
	Cases  []Case
}

// Header is the first line of cases.tsv.
const Header = "suite\tfile\tsyntax\tsha256\tverdict\trules\tintent\tnote"

// LoadIndex reads suites.tsv and cases.tsv of dir.
func LoadIndex(dir string) (*Index, error) {
	idx := &Index{Dir: dir}
	err := eachRecord(filepath.Join(dir, "suites.tsv"), func(f []string) error {
		if len(f) < 5 {
			return fmt.Errorf("suites.tsv: %d fields, want at least 5", len(f))
		}
		ref := strings.Fields(f[4])
		if len(ref) != 2 {
			return fmt.Errorf("suites.tsv: reference %q, want \"lang lib\"", f[4])
		}
		s := Suite{Name: f[0], Syntaxes: strings.Fields(f[1]), Countries: strings.Fields(f[2]), Title: f[3],
			RefLang: ref[0], RefLib: ref[1]}
		if len(f) > 5 {
			s.External = strings.Fields(f[5])
		}
		if len(f) > 6 {
			for _, r := range strings.Fields(f[6]) {
				s.Unchecked = append(s.Unchecked, Canonical(r))
			}
		}
		idx.Suites = append(idx.Suites, s)
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = eachRecord(filepath.Join(dir, "cases.tsv"), func(f []string) error {
		for len(f) < 8 {
			f = append(f, "")
		}
		idx.Cases = append(idx.Cases, Case{Suite: f[0], File: f[1], Syntax: f[2], SHA: f[3], Verdict: f[4],
			Rules: strings.Fields(f[5]), Intent: strings.Fields(f[6]), Note: f[7]})
		return nil
	})
	if os.IsNotExist(err) {
		err = nil
	}
	return idx, err
}

// eachRecord calls fn with the tab-separated fields of each line of a file
// but the first (the header) and the blank ones.
func eachRecord(path string, fn func([]string) error) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	first := true
	for sc.Scan() {
		line := sc.Text()
		if first {
			first = false
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		if err := fn(strings.Split(line, "\t")); err != nil {
			return err
		}
	}
	return sc.Err()
}

// WriteCases writes cases.tsv.
func (idx *Index) WriteCases() error {
	var b strings.Builder
	b.WriteString(Header + "\n")
	for _, c := range idx.Cases {
		fmt.Fprintf(&b, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", c.Suite, c.File, c.Syntax, c.SHA, c.Verdict,
			strings.Join(c.Rules, " "), strings.Join(c.Intent, " "), c.Note)
	}
	return os.WriteFile(filepath.Join(idx.Dir, "cases.tsv"), []byte(b.String()), 0o644)
}

// Suite returns the suite of that name.
func (idx *Index) Suite(name string) (Suite, bool) {
	for _, s := range idx.Suites {
		if s.Name == name {
			return s, true
		}
	}
	return Suite{}, false
}

// Countries returns the countries of the suites, in the order of suites.tsv.
func (idx *Index) Countries() []string {
	var out []string
	seen := map[string]bool{}
	for _, s := range idx.Suites {
		for _, c := range s.Countries {
			if !seen[c] {
				seen[c] = true
				out = append(out, c)
			}
		}
	}
	return out
}

// Outcome is what a library returned on one case.
type Outcome struct {
	Case   string   `json:"case"`
	Errors int      `json:"errors"`
	Rules  []string `json:"rules"`
	Error  string   `json:"error,omitempty"` // exception, timeout or crash
}

// Rejected reports whether the library rejected the document: at least one
// error, or an exception.
func (o Outcome) Rejected() bool { return o.Errors > 0 || o.Error != "" }

// Run is the line of a cases_*.jsonl file: one library on one suite.
type Run struct {
	Lang    string    `json:"lang"`
	Lib     string    `json:"lib"`
	Pkg     string    `json:"pkg"`
	Version string    `json:"version"`
	Suite   string    `json:"suite"`
	Results []Outcome `json:"results"`
	Note    string    `json:"note"`
}

// Key identifies a library across the files: lang/lib.
func (r Run) Key() string { return r.Lang + "/" + r.Lib }

// LoadRuns reads the cases_*.jsonl files of dir.
func LoadRuns(dir string) ([]Run, error) {
	paths, _ := filepath.Glob(filepath.Join(dir, "cases_*.jsonl"))
	sort.Strings(paths)
	var runs []Run
	for _, p := range paths {
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 64<<20), 64<<20)
		for sc.Scan() {
			if strings.TrimSpace(sc.Text()) == "" {
				continue
			}
			var r Run
			if err := json.Unmarshal(sc.Bytes(), &r); err != nil {
				f.Close()
				return nil, fmt.Errorf("%s: %v", p, err)
			}
			runs = append(runs, r)
		}
		f.Close()
		if err := sc.Err(); err != nil {
			return nil, fmt.Errorf("%s: %v", p, err)
		}
	}
	return runs, nil
}

// ruleID finds the rule identifiers of the official artefacts in what a
// library reports: EN 16931 (BR, BR-CO, BR-CL, BR-DEC, VAT categories,
// CII-SR, CII-DT, UBL-CR, UBL-SR, UBL-DT), XRechnung (BR-DE, BR-DEX,
// BR-TMP), Peppol, and the code list assertions of the Factur-X schematrons.
var ruleID = regexp.MustCompile(`(?i)\b(PEPPOL-EN16931-[A-Z]{0,3}\d+(?:-\d+)?|FX-SCH-A-\d+|(?:BR|CII|UBL)(?:-[A-Z]+)*-\d+(?:-[A-Z0-9]+)?)\b`)

// schemaError recognizes a schema violation reported without a rule
// identifier: KoSIT reports it as XSD, the libraries as XSD or schema errors,
// Xerces as cvc-…, libxml2 as "Element '…': This element is not expected".
var schemaError = regexp.MustCompile(`(?i)\bxsd\b|schema|\bcvc-|element is not expected|missing child element`)

// equivalents maps an identifier to the one of the reference that checks the
// same thing:
//   - the code list assertions of the Factur-X 1.09.2 schematron of EN 16931,
//     to the CEN rule on the same element (as ZeRechnung's parity tool,
//     conformance/parity/compare.go);
//   - the IGIC and IPSI rules in CII, which CEN 1.3.17 renamed from BR-AF-*
//     to BR-IP-* (IGIC, category L) and from BR-AG-* to BR-IG-* (IPSI,
//     category M); the KoSIT configuration runs CEN 1.3.16.
var equivalents = map[string]string{
	"FX-SCH-A-000571": "BR-CL-25", // scheme of BT-49 (EAS)
	"FX-SCH-A-000581": "BR-CL-11", // scheme of BT-30 (ICD)
	"FX-SCH-A-000612": "BR-CL-23", // unit of BT-149
	"FX-SCH-A-000613": "BR-CL-23", // unit of BT-130
}

// Canonical returns the identifier under which a rule is compared: upper
// case, equivalents applied.
func Canonical(id string) string {
	id = strings.ToUpper(id)
	if e, ok := equivalents[id]; ok {
		return e
	}
	for _, p := range [][2]string{{"BR-AF-", "BR-IP-"}, {"BR-AG-", "BR-IG-"}} {
		if rest, ok := strings.CutPrefix(id, p[0]); ok {
			return p[1] + rest
		}
	}
	return id
}

// Normalize returns the canonical identifiers found in the rules a library
// reported, and XSD for a schema violation reported without one.
func Normalize(reported []string) map[string]bool {
	out := map[string]bool{}
	for _, r := range reported {
		ids := ruleID.FindAllString(r, -1)
		for _, id := range ids {
			out[Canonical(id)] = true
		}
		if len(ids) == 0 && schemaError.MatchString(r) {
			out["XSD"] = true
		}
	}
	return out
}

// Score is the result of one library on the cases of one suite and syntax
// (or of several, for the totals).
type Score struct {
	Valid, Accepted int // valid cases, and those the library accepted
	Invalid         int // invalid cases
	Rejected        int // invalid cases the library rejected (any error or exception)
	Exact           int // invalid cases where it raised at least one expected rule
	Missing         int // cases without an outcome
	// Tolerated counts the valid cases the library rejected under rules the
	// reference does not check only (Suite.Unchecked), counted as accepted.
	Tolerated int
	// Credit per expected rule: mean over the cases that expect it of 1
	// (rule raised), 0.5 (rejected otherwise) or 0 (accepted).
	Credit map[string]float64
	// Raised is the set of expected rules the library raised on at least one
	// of their cases.
	Raised map[string]bool
}

// S is the specificity: the share of valid cases accepted (1 without valid case).
func (s Score) S() float64 {
	if s.Valid == 0 {
		return 1
	}
	return float64(s.Accepted) / float64(s.Valid)
}

// D is the detection: the mean credit of the expected rules (1 without invalid case).
func (s Score) D() float64 {
	if len(s.Credit) == 0 {
		return 1
	}
	t := 0.0
	for _, c := range s.Credit {
		t += c
	}
	return t / float64(len(s.Credit))
}

// Note is the score out of 100: 100 × √(S × D).
func (s Score) Note() float64 { return 100 * sqrt(s.S()*s.D()) }

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 50; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}

// Judge scores the outcomes of a library on cases (confirmed ones only).
func (idx *Index) Judge(cases []Case, outcomes map[string]Outcome) Score {
	s := Score{Credit: map[string]float64{}, Raised: map[string]bool{}}
	sum, n := map[string]float64{}, map[string]int{}
	for _, c := range cases {
		if c.Verdict != Valid && c.Verdict != Invalid {
			continue
		}
		o, ok := outcomes[c.File]
		if !ok {
			s.Missing++
			o = Outcome{Error: "no outcome"}
		}
		if c.Verdict == Valid {
			s.Valid++
			switch {
			case !o.Rejected():
				s.Accepted++
			case idx.onlyUnchecked(c.Suite, o):
				s.Accepted++
				s.Tolerated++
			}
			continue
		}
		s.Invalid++
		got := Normalize(o.Rules)
		// KoSIT does not run the schematrons on a document its XSD rejects:
		// when the reference stops at the schema, raising a rule the case
		// was written to break counts as raising the expected one.
		intent := false
		if len(c.Rules) == 1 && strings.EqualFold(c.Rules[0], "XSD") {
			for _, r := range c.Intent {
				intent = intent || got[Canonical(r)]
			}
		}
		exact := false
		for _, r := range c.Rules {
			r = Canonical(r)
			cr := 0.0
			switch {
			case got[r] || intent:
				cr, exact = 1, true
				s.Raised[r] = true
			case o.Rejected():
				cr = 0.5
			}
			sum[r] += cr
			n[r]++
		}
		if o.Rejected() {
			s.Rejected++
		}
		if exact {
			s.Exact++
		}
	}
	for r := range sum {
		s.Credit[r] = sum[r] / float64(n[r])
	}
	return s
}

// onlyUnchecked reports whether every rule of a rejection is one the
// reference of the suite does not check. Each reported rule must carry such
// an identifier: an exception, a message without identifier or any other rule
// is held against the library.
func (idx *Index) onlyUnchecked(suite string, o Outcome) bool {
	s, ok := idx.Suite(suite)
	if !ok || len(s.Unchecked) == 0 || o.Error != "" || len(o.Rules) == 0 {
		return false
	}
	for _, r := range o.Rules {
		ids := Normalize([]string{r})
		if len(ids) == 0 {
			return false
		}
		for id := range ids {
			if !contains(s.Unchecked, id) {
				return false
			}
		}
	}
	return true
}

// Part is one suite and one syntax: the unit the notes are computed on.
type Part struct{ Suite, Syntax string }

// Parts returns the suites and syntaxes of a country that have confirmed cases.
func (idx *Index) Parts(country string) []Part {
	var out []Part
	for _, s := range idx.Suites {
		if !contains(s.Countries, country) {
			continue
		}
		for _, syn := range s.Syntaxes {
			for _, c := range idx.Cases {
				if c.Suite == s.Name && c.Syntax == syn && (c.Verdict == Valid || c.Verdict == Invalid) {
					out = append(out, Part{s.Name, syn})
					break
				}
			}
		}
	}
	return out
}

// CasesOf returns the cases of a part.
func (idx *Index) CasesOf(p Part) []Case {
	var out []Case
	for _, c := range idx.Cases {
		if c.Suite == p.Suite && c.Syntax == p.Syntax {
			out = append(out, c)
		}
	}
	return out
}

func contains(xs []string, x string) bool {
	for _, y := range xs {
		if y == x {
			return true
		}
	}
	return false
}

// Library is the results of one library on every suite it was run on.
type Library struct {
	Lang, Lib, Pkg, Version, Note string
	Outcomes                      map[string]map[string]Outcome // suite → case → outcome
}

// Libraries groups runs by library, in the order they first appear.
func Libraries(runs []Run) []*Library {
	var out []*Library
	byKey := map[string]*Library{}
	for _, r := range runs {
		l := byKey[r.Key()]
		if l == nil {
			l = &Library{Lang: r.Lang, Lib: r.Lib, Pkg: r.Pkg, Version: r.Version, Note: r.Note,
				Outcomes: map[string]map[string]Outcome{}}
			byKey[r.Key()] = l
			out = append(out, l)
		}
		m := map[string]Outcome{}
		for _, o := range r.Results {
			m[o.Case] = o
		}
		l.Outcomes[r.Suite] = m
	}
	return out
}

// CountryScore is the result of a library for a country.
type CountryScore struct {
	Parts map[Part]Score
	Total Score   // all the parts together, for the counts
	Note  float64 // mean of the notes of the parts; a part not run counts 0
	Run   bool    // the library was run on at least one suite of the country
}

// ForCountry scores a library on the parts of a country.
func (idx *Index) ForCountry(l *Library, country string) CountryScore {
	cs := CountryScore{Parts: map[Part]Score{}}
	parts := idx.Parts(country)
	var all []Case
	total := map[string]Outcome{}
	sum := 0.0
	for _, p := range parts {
		outs, ok := l.Outcomes[p.Suite]
		if !ok {
			continue
		}
		cs.Run = true
		cases := idx.CasesOf(p)
		s := idx.Judge(cases, outs)
		cs.Parts[p] = s
		sum += s.Note()
		all = append(all, cases...)
		for k, v := range outs {
			total[k] = v
		}
	}
	if len(parts) > 0 {
		cs.Note = sum / float64(len(parts))
	}
	cs.Total = idx.Judge(all, total)
	return cs
}

// Family groups the rules for the detail tables.
func Family(rule string) string {
	r := strings.ToUpper(rule)
	switch {
	case r == "XSD":
		return "XSD"
	case strings.HasPrefix(r, "BR-DEX-"):
		return "BR-DEX"
	case strings.HasPrefix(r, "BR-DE-"):
		return "BR-DE"
	case strings.HasPrefix(r, "BR-TMP-"):
		return "BR-TMP"
	case strings.HasPrefix(r, "PEPPOL-"):
		return "Peppol"
	case strings.HasPrefix(r, "BR-CO-"):
		return "BR-CO"
	case strings.HasPrefix(r, "BR-CL-"):
		return "BR-CL"
	case strings.HasPrefix(r, "BR-DEC-"):
		return "BR-DEC"
	case strings.HasPrefix(r, "CII-"):
		return "CII-SR, CII-DT"
	case strings.HasPrefix(r, "UBL-"):
		return "UBL-CR, UBL-SR, UBL-DT"
	case regexp.MustCompile(`^BR-(S|Z|E|AE|K|G|O|L|M|IC|IG|IP|AF|AG|B)-\d`).MatchString(r):
		return "VAT categories"
	case regexp.MustCompile(`^BR-\d`).MatchString(r):
		return "BR"
	}
	return "other"
}

// Families lists the families of rules, in the order of the detail tables.
var Families = []string{"BR", "BR-CO", "BR-CL", "BR-DEC", "VAT categories", "CII-SR, CII-DT",
	"UBL-CR, UBL-SR, UBL-DT", "BR-DE", "BR-DEX", "BR-TMP", "Peppol", "XSD", "other"}
