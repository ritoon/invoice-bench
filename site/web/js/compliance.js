/**
 * compliance.js — the "Compliance" page: confidence per country (conformance
 * cases), accuracy on the official ZUGFeRD corpus (surveys of both
 * campaigns) and the standards each library enforces.
 */
const ComplianceApp = {
  init(root, data) {
    this.d = data;
    this.langName = (id) => (data.langs.find((l) => l.id === id) || { name: id }).name;
    this.state = { langs: new Set(data.langs.map((l) => l.id)), sort: 'agreement', details: false };
    this.ref = this.reference();
    this.validators = data.libs.filter((l) => l.validate).map((l) => this.analyze(l));

    this.country = (data.countries || [])[0] || null;

    this.filters = h('div', { class: 'flex flex-wrap items-end gap-x-6 gap-y-4' });
    this.confidence = h('div', { class: 'mt-2' });
    this.matrixBody = h('div', { class: 'mt-4' });
    this.coverage = h('div', { class: 'grid gap-4 lg:grid-cols-2' });
    this.standards = h('div', { class: 'mt-4' });

    root.replaceChildren(...[
      this.intro(),
      h('div', { class: 'sticky top-0 z-10 -mx-4 px-4 sm:-mx-6 sm:px-6 py-3 backdrop-blur border-b border-line', style: { background: 'color-mix(in srgb, var(--page) 92%, transparent)' } }, this.filters),
      this.country ? this.confidenceSection() : null,
      this.matrixSection(),
      this.coverageSection(),
      this.standardsSection(),
      this.notes(),
    ].filter(Boolean));
    this.renderFilters();
    this.update();
  },

  // --- Analysis -----------------------------------------------------------

  /** The official schematrons' verdict, and its CEN-only subset. */
  reference() {
    const flags = new Map(this.d.reference.flags.map((f) => [f.file, f]));
    // CII-SR-465 is a warning in the CEN artifacts; KoSIT raises it to an
    // error. A validator faithful to CEN alone flags the other files only.
    const cen = new Set([...flags].filter(([, f]) => (f.rules || []).some((r) => r !== 'CII-SR-465')).map(([file]) => file));
    const lib = this.d.libs.find((l) => l.id === this.d.reference.lib);
    return { flags, cen, lib };
  },

  analyze(lib) {
    const files = this.d.files;
    const flags = new Map((lib.validate.flags || []).map((f) => [f.file, f]));
    const n = { match: 0, miss: 0, extra: 0, error: 0, clean: 0 };
    const cells = files.map((file) => {
      const flag = flags.get(file);
      const official = this.ref.flags.get(file);
      let state;
      if (flag) {
        const rules = flag.rules || [];
        if (flag.error && !rules.length) state = 'error';
        else state = official ? 'match' : 'extra';
      } else state = official ? 'miss' : 'clean';
      n[state]++;
      return { state, flag };
    });
    const flagged = flags.size;
    const wrong = n.extra + n.error;
    let group, verdict;
    const sameAs = (set) => flagged === set.size && [...set].every((f) => flags.has(f));
    if (!wrong && !n.miss) {
      group = 0;
      verdict = 'Same files';
    } else if (!wrong && sameAs(this.ref.cen)) {
      group = 1;
      verdict = 'Same files as CEN alone';
    } else {
      group = flagged === 0 ? 3 : 2;
      const plural = (k, one, many) => `${k} ${k === 1 ? one : many}`;
      verdict = [
        n.extra && plural(n.extra, 'false positive', 'false positives'),
        n.error && plural(n.error, 'exception', 'exceptions'),
        n.miss && `${n.miss} missed`,
      ].filter(Boolean).join(' · ');
    }
    return { id: lib.id, name: lib.name, lang: lib.lang, kind: lib.kind, cells, flagged, n, group, verdict, score: wrong + n.miss };
  },

  groupLabels: ['Flag the same files as the official schematrons', 'Flag the same files as CEN alone', 'Differ from the official schematrons', 'Flag nothing'],

  // --- Sections -----------------------------------------------------------

  intro() {
    const v = this.validators;
    const same = v.filter((r) => r.group === 0);
    const pdfTotal = Math.max(0, ...this.d.libs.filter((l) => l.extract).map((l) => l.extract.total));
    const pdfAll = this.d.libs.filter((l) => l.extract && l.extract.ok >= l.extract.total).length;
    const ct = this.country;
    const zr = ct && ct.libs.find((l) => l.kind === 'self');
    const high = ct ? ct.libs.filter((l) => Math.round(l.score) >= 90).length : 0;
    const tiles = h('div', { class: 'mt-8 grid grid-cols-2 gap-3 lg:grid-cols-4' },
      ct ? StatTile.create({ label: `ZeRechnung · confidence for ${ct.name}`, value: zr ? `${Fmt.decimal(zr.score, 0)} / 100` : '—', sub: 'An upper bound: the cases come from its tests', emphasis: true }) : null,
      ct ? StatTile.create({ label: `Libraries rated for ${ct.name}`, value: String(ct.libs.length), sub: `${high} score 90 or more` }) : null,
      StatTile.create({ label: 'Same corpus files flagged as the official schematrons', value: `${same.length} of ${v.length}`, sub: same.some((r) => r.kind === 'self') ? 'ZeRechnung among them' : '', emphasis: !ct }),
      StatTile.create({ label: 'Hybrid PDFs of the corpus', value: String(pdfTotal), sub: `${pdfAll} libraries extract every one` }));
    return h('section', { class: 'pt-6 pb-8 sm:pt-10' },
      h('p', { class: 'eyebrow', text: `Run of ${this.d.run.date} · conformance cases and official ZUGFeRD corpus` }),
      h('h1', { class: 'mt-2 text-3xl sm:text-4xl font-semibold tracking-tight', text: 'Compliance' }),
      h('p', { class: 'mt-3 max-w-3xl text-base text-ink-2 leading-relaxed', text: `Do the ${this.d.libs.length} libraries enforce the rules of the standard? Three views: a confidence score per country, from documents whose verdict the reference validator gives; what each validator flags on the valid invoices of the official ZUGFeRD corpus; and which standards each library actually enforces, surveyed in its code rather than in its README.` }),
      tiles);
  },

  confidenceSection() {
    const ct = this.country;
    return h('section', { class: 'pt-10 pb-12', 'aria-labelledby': 'confidence' },
      h('h2', { id: 'confidence', class: 'text-xl font-semibold tracking-tight', text: `Confidence for ${ct.name}` }),
      h('div', { class: 'mt-2 max-w-3xl space-y-2 text-sm text-ink-2 leading-relaxed prose-site' },
        h('p', {}, `Each library validates ${Fmt.int(ct.valid)} valid and ${Fmt.int(ct.invalid)} invalid documents (${ct.parts.map((p) => `${p.title} ${p.syntax.toUpperCase()}`).join(', ')}), with the rule set of the survey. The expected verdicts and rules come from the reference: ${ct.reference}.`),
        h('p', {}, h('strong', { class: 'text-ink font-semibold', text: 'Score = 100 × √(S × D)' }),
          `, S the share of valid cases accepted, D the mean detection of the ${Fmt.int(ct.rules)} expected rules (1 raised, 0.5 rejected under another rule, 0 accepted). Rejecting everything or accepting everything both score 0. ${ct.name} is the mean of its ${ct.parts.length} parts; a syntax the library does not read counts 0.`),
        h('p', {}, h('strong', { class: 'text-ink font-semibold', text: 'Limits. ' }),
          'The mutation cases come from ZeRechnung’s own tests, so its score is an upper bound. Only errors count, XML only: no hybrid PDF, no Peppol, no French rules yet. Details in ', fileLink('cases/README.md#limits', 'cases/README.md'), '.')),
      h('div', { class: 'card mt-5 p-2 sm:p-3' }, this.confidence));
  },

  renderFilters() {
    const counts = new Map();
    for (const l of this.d.libs) counts.set(l.lang, (counts.get(l.lang) || 0) + 1);
    this.filters.replaceChildren(
      Field.create({
        label: 'Languages',
        control: LangChips.create({
          langs: this.d.langs.map((l) => ({ ...l, count: counts.get(l.id) || 0 })),
          selected: this.state.langs,
          onChange: () => this.update(),
        }),
      }),
    );
  },

  matrixSection() {
    const rulesCount = new Map();
    for (const f of this.ref.flags.values()) for (const r of f.rules || []) rulesCount.set(r, (rulesCount.get(r) || 0) + 1);
    const ruleText = [...rulesCount].map(([r, k]) => `${r} on ${k} file${k > 1 ? 's' : ''}`).join(', ');
    const refName = this.ref.lib ? `${this.ref.lib.name} (${this.langName(this.ref.lib.lang)})` : this.d.reference.lib;
    this.sortControl = Segmented.create({
      label: 'Order', value: this.state.sort, small: true,
      options: [{ value: 'agreement', text: 'Closest first' }, { value: 'language', text: 'By language' }],
      onChange: (v) => { this.state.sort = v; this.update(); },
    });
    return h('section', { class: 'pt-10 pb-12', 'aria-labelledby': 'flagged' },
      h('h2', { id: 'flagged', class: 'text-xl font-semibold tracking-tight', text: 'Valid invoices of the corpus flagged' }),
      h('div', { class: 'mt-2 max-w-3xl space-y-2 text-sm text-ink-2 leading-relaxed' },
        h('p', {}, `All ${this.d.files.length} files are valid examples of the corpus. The official schematrons, run through ${refName}, still flag ${this.ref.flags.size} of them: ${ruleText}.`),
        rulesCount.has('CII-SR-465') ? h('p', { text: `CII-SR-465 is a warning in the CEN artifacts, which the KoSIT configuration for XRechnung raises to an error: a validator faithful to CEN alone flags ${this.ref.cen.size}. Both count as agreeing below. A file flagged by an exception, rather than by a rule, is not.` }) : null),
      h('div', { class: 'card mt-5 p-5' },
        h('div', { class: 'flex flex-wrap items-center justify-between gap-3' }, AgreementMatrix.legend(), this.sortControl),
        this.matrixBody));
  },

  coverageSection() {
    return h('section', { class: 'pb-12', 'aria-labelledby': 'coverage' },
      h('h2', { id: 'coverage', class: 'text-xl font-semibold tracking-tight', text: 'Files read' }),
      h('p', { class: 'mt-1 max-w-3xl text-sm text-ink-2', text: 'Before any rule: can the library get the XML out of every hybrid PDF of the corpus, and parse every valid CII invoice into its model?' }),
      h('div', { class: 'mt-5' }, this.coverage));
  },

  standardsSection() {
    const s = this.d.standards;
    const toggle = Toggle.create({ label: 'Show what each cell covers', checked: this.state.details, onChange: (v) => { this.state.details = v; this.update(); } });
    return h('section', { class: 'pb-12', 'aria-labelledby': 'standards' },
      h('h2', { id: 'standards', class: 'text-xl font-semibold tracking-tight', text: 'Enforced standards' }),
      h('p', { class: 'mt-1 max-w-3xl text-sm text-ink-2 prose-site' },
        'Surveyed in the libraries’ code, not in their README: files and lines are in each language’s ', h('code', { text: 'COVERAGE.md' }),
        '. A selection: the most downloaded and the most complete libraries of each language.'),
      h('div', { class: 'card mt-5 p-2 sm:p-3' },
        h('div', { class: 'flex flex-wrap items-center justify-between gap-3 px-3 pt-2' }, StandardsMatrix.legend(s), toggle),
        this.standards));
  },

  notes() {
    const li = (...c) => h('li', {}, ...c);
    return h('section', { class: 'pb-6', 'aria-labelledby': 'notes' },
      h('div', { class: 'card p-5 sm:p-6' },
        h('h2', { id: 'notes', class: 'text-base font-semibold', text: 'How these figures were obtained' }),
        h('ul', { class: 'mt-3 grid gap-x-8 gap-y-2 text-sm text-ink-2 leading-relaxed md:grid-cols-2 list-disc pl-5 prose-site' },
          li('Each validator applies the most complete rule set it has offline, without the French rules (the corpus is German), with the same adapter for the corpus and the conformance cases. Remote API clients are left out.'),
          li('The confidence scores are computed by the ', fileLink('cases/cases.go', 'cases'), ' package, the one the report tools use: they match the tables of ', fileLink(`results/${this.d.run.date}/tables.md`, 'tables.md'), ' and ', fileLink(`go/results/${this.d.run.goDate}/tables.md`, 'go/results'), '.'),
          li('The Go libraries other than ZeRechnung come from the ', fileLink('go/README.md', 'Go campaign'), ', one module per library; their surveys give the files they flag, not the rules.'),
          li('The comparison is by file. Libraries name the same check differently (the FeRD schematron recodes BR-CL-25 as FX-SCH-A-000571): the rules each one reports are shown on hover, next to the official ones.'),
          li('The standards grid is kept by hand from ', fileLink('README.md#enforced-standards', 'README.md'), '; the levels (enforced, partial, external) are a reading of its cells.')),
        h('p', { class: 'mt-4 text-xs text-muted prose-site' }, 'Raw surveys: ', fileLink(`results/${this.d.run.date}`, `results/${this.d.run.date}`), ' (', h('code', { text: 'survey_*.jsonl' }), ') and ',
          fileLink(`go/results/${this.d.run.goDate}`, `go/results/${this.d.run.goDate}`), ' (', h('code', { text: 'survey_*.txt' }), ').')));
  },

  // --- Rendering ----------------------------------------------------------

  update() {
    const show = (lang) => this.state.langs.has(lang);
    const langOrder = (id) => this.d.langs.findIndex((l) => l.id === id);

    // Confidence
    if (this.country) {
      const libs = this.country.libs.filter((l) => show(l.lang));
      this.confidence.replaceChildren(libs.length
        ? ConfidenceTable.create({ country: this.country, libs, langName: this.langName })
        : h('p', { class: 'py-8 text-center text-sm text-muted', text: 'No library in the selected languages.' }));
    }

    // Matrix
    let rows = this.validators.filter((r) => show(r.lang));
    if (this.state.sort === 'agreement') {
      rows.sort((a, b) => a.group - b.group || (b.kind === 'self') - (a.kind === 'self') || a.score - b.score
        || langOrder(a.lang) - langOrder(b.lang) || a.name.localeCompare(b.name));
      rows = rows.map((r) => ({ ...r, groupLabel: `${this.groupLabels[r.group]}` }));
    } else {
      rows = rows.map((r) => ({ ...r, groupLabel: this.langName(r.lang) }));
    }
    const refRow = {
      name: 'Official schematrons', sub: 'CEN EN 16931 + KoSIT, official XSLT', kind: 'reference',
      cells: this.d.files.map((f) => ({ state: this.ref.flags.has(f) ? 'match' : 'clean', flag: this.ref.flags.get(f) })),
      flagged: this.ref.flags.size, verdict: 'Reference',
    };
    const counts = new Map();
    for (const r of rows) counts.set(r.groupLabel, (counts.get(r.groupLabel) || 0) + 1);
    rows = rows.map((r) => ({ ...r, groupLabel: `${r.groupLabel} · ${counts.get(r.groupLabel)}` }));
    this.matrixBody.replaceChildren(
      rows.length
        ? AgreementMatrix.create({ files: this.d.files, reference: { flags: this.d.reference.flags, row: refRow }, rows, groups: true, langName: this.langName })
        : h('p', { class: 'py-8 text-center text-sm text-muted', text: 'No validator in the selected languages.' }));

    // Coverage
    const items = (op) => this.d.libs.filter((l) => l[op] && show(l.lang)).map((l) => ({ id: l.id, name: l.name, lang: l.lang, kind: l.kind, ok: l[op].ok, total: l[op].total }));
    const ex = items('extract'), pa = items('parse');
    const total = (list, fallback) => (list.length ? Math.max(...list.map((i) => i.total)) : fallback);
    this.coverage.replaceChildren(
      CoverageList.create({ title: 'XML extracted from the hybrid PDFs', unit: 'PDFs', total: total(ex, 104), items: ex, langName: this.langName }),
      CoverageList.create({ title: 'Valid CII invoices parsed into a model', unit: 'files', total: total(pa, 26), items: pa, langName: this.langName }));

    // Standards
    const s = this.d.standards;
    const srows = s.rows.filter((r) => show(r.lang));
    this.standards.replaceChildren(srows.length
      ? StandardsMatrix.create({ standards: s, rows: srows, langName: this.langName, details: this.state.details })
      : h('p', { class: 'py-8 text-center text-sm text-muted', text: 'No library in the selected languages.' }));
  },
};

document.addEventListener('DOMContentLoaded', () => Page.boot({ active: 'compliance', data: 'compliance', app: ComplianceApp }));
