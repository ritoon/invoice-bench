package cases

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// CountryNames are the names of the countries in the tables.
var CountryNames = map[string]string{
	"DE": "Germany", "FR": "France", "NL": "Netherlands", "IT": "Italy", "ES": "Spain", "PT": "Portugal",
	"BE": "Belgium", "AT": "Austria", "PL": "Poland",
}

// References describes the reference validators named in suites.tsv.
var References = map[string]string{
	"java/kosit-validator": "KoSIT validator 1.6.3 with the official XRechnung 3.0.2 configuration of 2026-08-31 " +
		"(CII D16B and UBL 2.1 XSD, CEN EN 16931 1.3.16, KoSIT XRechnung schematron 2.6.0, levels of its scenarios)",
}

// Options tells Markdown how to present the libraries.
type Options struct {
	// Name returns the name of a library in the tables (lang/lib key).
	Name func(lang, lib string) string
	// Lang returns the name of a language; nil leaves the Language column
	// out (a bench of one language).
	Lang func(lang string) string
	// Ops lists, by lang/lib key, the operations the bench pass measured
	// successfully ("parse", "write"…), for the Operations column; nil
	// leaves the column out.
	Ops map[string][]string
	// OpsSource says where Ops come from, for the text above the tables.
	OpsSource string
	// Unrated lists the libraries of the bench without a validation
	// operation, which are not rated.
	Unrated []string
	// Link is the relative link to cases/README.md.
	Link string
}

// Markdown writes the confidence tables: one section per country, then the
// rules raised by family.
func Markdown(w io.Writer, idx *Index, runs []Run, o Options) {
	libs := Libraries(runs)
	if len(libs) == 0 {
		return
	}
	fmt.Fprint(w, "\n#### Confidence by country\n\n")
	fmt.Fprintf(w, "Score out of 100 of each library on the conformance cases of a country ([`cases/`](%s)): "+
		"100 × √(S × D), where S is the share of valid cases the library accepts and D the mean detection of the expected rules, "+
		"each rule weighing the same: 1 when the library raises the rule on a case, 0.5 when it rejects the document under another rule "+
		"or with an exception, 0 when it accepts it. The expected verdict and rules are those of the reference validator. "+
		"The score of a country is the mean of the scores of its suites and syntaxes: a syntax the library does not read counts 0.\n", o.Link)
	if len(o.Unrated) > 0 {
		fmt.Fprintf(w, "\nNot rated, without a validation operation: %s.\n", strings.Join(o.Unrated, ", "))
	}
	for _, country := range idx.Countries() {
		parts := idx.Parts(country)
		if len(parts) == 0 {
			continue
		}
		type row struct {
			l  *Library
			cs CountryScore
		}
		var rows []row
		for _, l := range libs {
			cs := idx.ForCountry(l, country)
			if cs.Run {
				rows = append(rows, row{l, cs})
			}
		}
		if len(rows) == 0 {
			continue
		}
		sort.SliceStable(rows, func(i, j int) bool {
			if a, b := round(rows[i].cs.Note), round(rows[j].cs.Note); a != b {
				return a > b
			}
			return strings.ToLower(o.Name(rows[i].l.Lang, rows[i].l.Lib)) < strings.ToLower(o.Name(rows[j].l.Lang, rows[j].l.Lib))
		})
		refs := idx.references(country)

		fmt.Fprintf(w, "\n#### Confidence: %s (%s)\n\n", countryName(country), country)
		fmt.Fprintln(w, idx.describe(country))
		for _, r := range refs {
			fmt.Fprintf(w, "\nReference: %s; it scores 100 by construction.\n", referenceName(r))
		}
		unchecked := idx.unchecked(country)
		if len(unchecked) > 0 {
			fmt.Fprintf(w, "\nRules of the norm that the reference does not check: %s. A valid case that a library rejects under these rules alone "+
				"counts as accepted.\n", strings.Join(unchecked, ", "))
		}
		if o.Ops != nil && o.OpsSource != "" {
			fmt.Fprintf(w, "\nOperations: those the library performed successfully in %s (validation excepted).\n", o.OpsSource)
		}
		fmt.Fprintln(w)
		head := []string{"Library"}
		if o.Lang != nil {
			head = append(head, "Language")
		}
		head = append(head, "Score")
		for _, p := range parts {
			head = append(head, idx.partTitle(p))
		}
		head = append(head, "Valid accepted", "Invalid rejected", "Expected rule raised")
		if o.Ops != nil {
			head = append(head, "Operations")
		}
		fmt.Fprintln(w, "| "+strings.Join(head, " | ")+" |")
		fmt.Fprintln(w, "|"+strings.Repeat("---|", len(head)))
		for _, r := range rows {
			name := o.Name(r.l.Lang, r.l.Lib)
			if contains(refs, r.l.Lang+"/"+r.l.Lib) {
				name += " (reference)"
			}
			cells := []string{name}
			if o.Lang != nil {
				cells = append(cells, o.Lang(r.l.Lang))
			}
			cells = append(cells, fmt.Sprintf("**%d**", round(r.cs.Note)))
			for _, p := range parts {
				if s, ok := r.cs.Parts[p]; ok {
					cells = append(cells, fmt.Sprint(round(s.Note())))
				} else {
					cells = append(cells, "—")
				}
			}
			t := r.cs.Total
			cells = append(cells, fmt.Sprintf("%d/%d", t.Accepted, t.Valid), fmt.Sprintf("%d/%d", t.Rejected, t.Invalid),
				fmt.Sprintf("%d/%d", t.Exact, t.Invalid))
			if o.Ops != nil {
				cells = append(cells, orDash(strings.Join(o.Ops[r.l.Lang+"/"+r.l.Lib], ", ")))
			}
			fmt.Fprintln(w, "| "+strings.Join(cells, " | ")+" |")
		}
		var tolerated []string
		for _, r := range rows {
			if n := r.cs.Total.Tolerated; n > 0 {
				name := o.Name(r.l.Lang, r.l.Lib)
				if o.Lang != nil {
					name += " (" + o.Lang(r.l.Lang) + ")"
				}
				tolerated = append(tolerated, fmt.Sprintf("%s %d", name, n))
			}
		}
		if len(tolerated) > 0 {
			fmt.Fprintf(w, "\nValid cases counted as accepted although the library rejected them, under %s alone: %s.\n",
				strings.Join(unchecked, ", "), strings.Join(tolerated, ", "))
		}

		// Rules raised by family.
		total := map[string]map[string]bool{} // family → expected rules
		for _, p := range parts {
			for _, c := range idx.CasesOf(p) {
				if c.Verdict != Invalid {
					continue
				}
				for _, r := range c.Rules {
					f := Family(r)
					if total[f] == nil {
						total[f] = map[string]bool{}
					}
					total[f][Canonical(r)] = true
				}
			}
		}
		var fams []string
		for _, f := range Families {
			if len(total[f]) > 0 {
				fams = append(fams, f)
			}
		}
		fmt.Fprintf(w, "\n#### Rules raised: %s (%s)\n\n", countryName(country), country)
		fmt.Fprintln(w, "Expected rules that the library raises under their identifier on at least one of their cases, by family "+
			"(after normalization: case, prefixes, the equivalences of `cases/cases.go`). "+
			"XSD: the documents the reference rejects on its schema alone; raising the rule the case was written to break counts too. "+
			"A library that reports no rule identifiers scores here 0, whatever its verdicts.")
		fmt.Fprintln(w)
		head = []string{"Library"}
		if o.Lang != nil {
			head = append(head, "Language")
		}
		for _, f := range fams {
			head = append(head, fmt.Sprintf("%s (%d)", f, len(total[f])))
		}
		fmt.Fprintln(w, "| "+strings.Join(head, " | ")+" |")
		fmt.Fprintln(w, "|"+strings.Repeat("---|", len(head)))
		for _, r := range rows {
			cells := []string{o.Name(r.l.Lang, r.l.Lib)}
			if o.Lang != nil {
				cells = append(cells, o.Lang(r.l.Lang))
			}
			for _, f := range fams {
				n := 0
				for rule := range total[f] {
					if r.cs.Total.Raised[rule] {
						n++
					}
				}
				cells = append(cells, fmt.Sprint(n))
			}
			fmt.Fprintln(w, "| "+strings.Join(cells, " | ")+" |")
		}
	}
}

// references returns the lang/lib keys of the references of a country's suites.
func (idx *Index) references(country string) []string {
	var out []string
	for _, s := range idx.Suites {
		k := s.RefLang + "/" + s.RefLib
		if contains(s.Countries, country) && !contains(out, k) {
			out = append(out, k)
		}
	}
	return out
}

// unchecked returns the rules that the references of a country's suites do
// not check.
func (idx *Index) unchecked(country string) []string {
	var out []string
	for _, s := range idx.Suites {
		if !contains(s.Countries, country) {
			continue
		}
		for _, r := range s.Unchecked {
			if !contains(out, r) {
				out = append(out, r)
			}
		}
	}
	return out
}

// describe sums up the cases of a country.
func (idx *Index) describe(country string) string {
	var parts []string
	valid, invalid, excluded := 0, 0, 0
	rules := map[string]bool{}
	for _, s := range idx.Suites {
		if !contains(s.Countries, country) {
			continue
		}
		var bySyn []string
		for _, syn := range s.Syntaxes {
			n := 0
			for _, c := range idx.Cases {
				if c.Suite != s.Name || c.Syntax != syn {
					continue
				}
				switch c.Verdict {
				case Valid:
					valid++
					n++
				case Invalid:
					invalid++
					n++
					for _, r := range c.Rules {
						rules[Canonical(r)] = true
					}
				case Excluded:
					excluded++
				}
			}
			if n > 0 {
				bySyn = append(bySyn, fmt.Sprintf("%s %d", strings.ToUpper(syn), n))
			}
		}
		parts = append(parts, fmt.Sprintf("%s (`%s`: %s)", s.Title, s.Name, strings.Join(bySyn, ", ")))
	}
	out := fmt.Sprintf("Suites: %s; %d valid cases, %d invalid, %d expected rules.", strings.Join(parts, "; "), valid, invalid, len(rules))
	if excluded > 0 {
		out += fmt.Sprintf(" %d cases are left out: the reference could not judge them (see the note column of `cases/cases.tsv`).", excluded)
	}
	return out
}

func (idx *Index) partTitle(p Part) string {
	s, _ := idx.Suite(p.Suite)
	return s.Title + ", " + strings.ToUpper(p.Syntax)
}

func countryName(c string) string {
	if n, ok := CountryNames[c]; ok {
		return n
	}
	return c
}

func referenceName(k string) string {
	if n, ok := References[k]; ok {
		return n
	}
	return k
}

func round(x float64) int { return int(x + 0.5) }

func orDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}
