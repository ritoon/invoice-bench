/**
 * perf.js — the "CPU & memory" page: resources pass of the cross-language campaign.
 */
const PerfApp = {
  init(root, data) {
    this.d = data;
    this.langName = (id) => (data.langs.find((l) => l.id === id) || { name: id }).name;
    this.libs = new Map(data.libs.map((l) => [l.id, l]));
    this.benches = data.ops.flatMap((o) => o.benches.map((b) => ({ ...b, op: o })));
    const fromHash = this.benches.find((b) => b.id === decodeURIComponent(location.hash.slice(1)));
    const start = fromHash || this.benches[0];
    this.state = {
      op: start.op.id,
      bench: start.id,
      langs: new Set(data.langs.map((l) => l.id)),
      floors: true,
    };

    root.replaceChildren(
      this.intro(),
      this.languageSection(),
      this.operationSection(),
      this.methodSection(),
    );
    this.update();

    // A link to another benchmark (#Parse_Large) opens it in place.
    window.addEventListener('hashchange', () => {
      const b = this.benches.find((x) => x.id === decodeURIComponent(location.hash.slice(1)));
      if (!b || b.id === this.state.bench) return;
      this.state.op = b.op.id;
      this.state.bench = b.id;
      this.renderFilters();
      this.update();
    });
  },

  intro() {
    const env = new Map(this.d.run.env);
    return h('section', { class: 'pt-6 pb-8 sm:pt-10 sm:pb-10' },
      h('p', { class: 'eyebrow', text: `Run of ${this.d.run.date} · ${env.get('cpu') || ''}` }),
      h('h1', { class: 'mt-2 text-3xl sm:text-4xl font-semibold tracking-tight', text: 'CPU time and memory' }),
      h('p', { class: 'mt-3 max-w-3xl text-base text-ink-2 leading-relaxed prose-site' },
        'ZeRechnung, the Go package, against ', String(this.d.libs.filter((l) => l.kind !== 'self' && l.kind !== 'floor').length),
        ' Factur-X, ZUGFeRD and EN 16931 libraries and ', String(this.d.libs.filter((l) => l.kind === 'floor').length),
        ' generic PDF and XML libraries in eight other languages. Each library runs in a fresh process that loads it and repeats a single benchmark for one second: ',
        h('strong', { class: 'text-ink font-semibold', text: 'CPU time' }), ' counts every thread of the process, ',
        h('strong', { class: 'text-ink font-semibold', text: 'peak memory' }), ' is its maximum resident set size, runtime included.'));
  },

  // --- Per language -------------------------------------------------------

  languageSection() {
    const rows = this.d.byLang;
    const mem = rows.map((s) => ({
      id: s.lang, label: this.langName(s.lang), sub: `${s.libraries} lib. · ${s.processes} runs`,
      min: s.rssMin, med: s.rssMedian, max: s.rssMax, emphasis: s.lang === 'go',
      tip: {
        title: this.langName(s.lang), sub: s.lang === 'go' ? 'ZeRechnung' : `${s.libraries} libraries, ${s.processes} processes`,
        rows: [['Median', Fmt.mib(s.rssMedian)], ['Lowest', Fmt.mib(s.rssMin)], ['Highest', Fmt.mib(s.rssMax)]],
      },
    }));
    const use = rows.map((s) => ({
      id: s.lang, label: this.langName(s.lang), sub: `${s.libraries} lib. · ${s.processes} runs`,
      min: s.useMin, med: s.useMedian, max: s.useMax, emphasis: s.lang === 'go',
      tip: {
        title: this.langName(s.lang), sub: 'CPU time / wall-clock time',
        rows: [['Median', Fmt.pct(s.useMedian)], ['Lowest', Fmt.pct(s.useMin)], ['Highest', Fmt.pct(s.useMax)]],
      },
    }));
    const memTicks = Scale.powers(Math.min(...rows.map((r) => r.rssMin)), Math.max(...rows.map((r) => r.rssMax)), 4);
    const useMax = Math.ceil(Math.max(...rows.map((r) => r.useMax)) + 0.001);
    const useTicks = Array.from({ length: useMax + 1 }, (_, i) => i);

    const memCard = h('div', { class: 'card p-5' },
      h('h3', { class: 'text-sm font-semibold', text: 'Peak memory of a process' }),
      h('p', { class: 'mt-1 text-xs text-ink-2', text: 'Median dot, lowest to highest across every library and benchmark. Log scale.' }));
    RangeChart.render(memCard, {
      rows: mem, scale: 'log', domain: [memTicks[0], memTicks[memTicks.length - 1]], ticks: memTicks,
      tickFormat: Fmt.mibTick, valueFormat: Fmt.mib, ariaLabel: 'Peak memory per language',
    }).classList.add('mt-4');

    const useCard = h('div', { class: 'card p-5' },
      h('h3', { class: 'text-sm font-semibold', text: 'CPU time ÷ wall-clock time' }),
      h('p', { class: 'mt-1 text-xs text-ink-2', text: '100% is one busy core. More means work on other threads: garbage collector, JIT.' }));
    RangeChart.render(useCard, {
      rows: use, scale: 'linear', domain: [0, useMax], ticks: useTicks,
      tickFormat: Fmt.pct, valueFormat: Fmt.pct, ariaLabel: 'CPU time over wall-clock time per language',
    }).classList.add('mt-4');

    const table = h('table', { class: 'data-table mt-3' },
      h('thead', {}, h('tr', {}, ['Language', 'Libraries', 'Processes', 'Peak memory min.', 'median', 'max.', 'CPU / wall median', 'max.']
        .map((t, i) => h('th', { class: i ? 'r' : null, scope: 'col', text: t })))),
      h('tbody', {}, rows.map((s) => h('tr', { class: s.lang === 'go' ? 'self' : null },
        h('td', { text: this.langName(s.lang) }),
        [s.libraries, s.processes].map((v) => h('td', { class: 'r num', text: Fmt.int(v) })),
        [s.rssMin, s.rssMedian, s.rssMax].map((v) => h('td', { class: 'r num', text: Fmt.mib(v) })),
        [s.useMedian, s.useMax].map((v) => h('td', { class: 'r num', text: Fmt.pct(v) }))))));

    return h('section', { class: 'pb-12', 'aria-labelledby': 'by-lang' },
      h('h2', { id: 'by-lang', class: 'text-xl font-semibold tracking-tight', text: 'By language' }),
      h('p', { class: 'mt-1 max-w-3xl text-sm text-ink-2', text: 'The runtime sets the floor: a few MiB for Rust, C and C++, tens for Go, Python and PHP, hundreds for Node, .NET and the JVM.' }),
      h('div', { class: 'mt-5 grid gap-4 lg:grid-cols-2' }, memCard, useCard),
      h('details', { class: 'table-view mt-3 card px-5 py-3' },
        h('summary', { text: 'Show as a table' }),
        h('div', { class: 'overflow-x-auto' }, table),
        h('p', { class: 'mt-3 text-xs text-muted', text: 'Without the Go round at GOMAXPROCS=1, nor the libraries that spawn child processes (their CPU and memory are not counted).' })));
  },

  // --- Per operation ------------------------------------------------------

  operationSection() {
    this.filters = h('div', { class: 'flex flex-wrap items-end gap-x-6 gap-y-4' });
    this.tiles = h('div', { class: 'grid grid-cols-2 gap-3 lg:grid-cols-4' });
    this.chartCard = h('div', { class: 'card p-5' });
    this.tableCard = h('div', { class: 'card p-2 sm:p-3' });
    this.renderFilters();
    return h('section', { class: 'pb-12', 'aria-labelledby': 'by-op' },
      h('h2', { id: 'by-op', class: 'text-xl font-semibold tracking-tight', text: 'By operation' }),
      h('p', { class: 'mt-1 max-w-3xl text-sm text-ink-2', text: 'Pick an operation and an input file. Every library that supports it, in its own process, on the same bytes of the official ZUGFeRD corpus.' }),
      h('div', { class: 'mt-5 sticky top-0 z-10 -mx-4 px-4 sm:-mx-6 sm:px-6 py-3 backdrop-blur border-b border-line', style: { background: 'color-mix(in srgb, var(--page) 92%, transparent)' } }, this.filters),
      h('div', { class: 'mt-5 space-y-4' }, this.tiles, this.chartCard, this.tableCard));
  },

  renderFilters() {
    const op = this.d.ops.find((o) => o.id === this.state.op);
    const counts = new Map();
    for (const r of this.d.rows) {
      if (r.bench !== this.state.bench || r.status !== 'ok') continue;
      const l = this.libs.get(r.lib);
      counts.set(l.lang, (counts.get(l.lang) || 0) + 1);
    }
    this.filters.replaceChildren(
      Field.create({
        label: 'Operation',
        control: Segmented.create({
          label: 'Operation', value: this.state.op,
          options: this.d.ops.map((o) => ({ value: o.id, text: o.title })),
          onChange: (v) => {
            this.state.op = v;
            this.state.bench = this.d.ops.find((o) => o.id === v).benches[0].id;
            this.renderFilters();
            this.update();
          },
        }),
      }),
      Field.create({
        label: 'Input',
        control: Segmented.create({
          label: 'Input file', value: this.state.bench, small: true,
          options: op.benches.map((b) => ({ value: b.id, text: b.label })),
          onChange: (v) => {
            this.state.bench = v;
            this.renderFilters();
            this.update();
          },
        }),
      }),
      Field.create({
        label: 'Languages',
        control: LangChips.create({
          langs: this.d.langs.map((l) => ({ ...l, count: counts.get(l.id) || 0 })),
          selected: this.state.langs,
          onChange: () => this.update(),
        }),
      }),
      Field.create({
        label: 'Floors',
        control: Toggle.create({ label: 'Show generic PDF / XML libraries', checked: this.state.floors, onChange: (v) => { this.state.floors = v; this.update(); } }),
      }),
    );
  },

  update() {
    const { bench } = this.state;
    const benchInfo = this.benches.find((b) => b.id === bench);
    history.replaceState(null, '', `#${bench}`);

    const all = this.d.rows.filter((r) => r.bench === bench).map((r) => ({ ...r, id: r.lib, lib: this.libs.get(r.lib) }));
    const rows = all.filter((r) => this.state.langs.has(r.lib.lang) && (this.state.floors || r.lib.kind !== 'floor'));
    const measured = (r) => r.status === 'ok' && r.cpu > 0 && r.rss > 0;
    const ref = all.find((r) => r.id === 'go/ZeRechnung' && measured(r));

    // Axes from every row of the benchmark, so filters never move them.
    const okAll = all.filter(measured);
    const cpuTicks = Scale.powers(Math.min(...okAll.map((r) => r.cpu)), Math.max(...okAll.map((r) => r.cpu)), 10);
    const rssTicks = Scale.powers(Math.min(...okAll.map((r) => r.rss)), Math.max(...okAll.map((r) => r.rss)), 4);
    const cpuDomain = [cpuTicks[0], cpuTicks[cpuTicks.length - 1]];
    const rssDomain = [rssTicks[0], rssTicks[rssTicks.length - 1]];

    this.renderTiles(rows.filter(measured), ref && rows.includes(ref) ? ref : null);

    // Scatter: child-process libraries are left out (their CPU is not counted).
    const points = rows.filter((r) => measured(r) && !r.child).map((r) => ({
      id: r.id, x: r.cpu, y: r.rss, kind: r.lib.kind,
      label: r.lib.kind === 'self' ? (r.id === 'go/ZeRechnung' ? 'ZeRechnung' : '1 core') : null,
      labelSide: r.id === 'go/ZeRechnung' ? 'right' : 'left',
      tip: {
        title: `${Fmt.ns(r.cpu)} · ${Fmt.mib(r.rss)}`,
        sub: `${r.lib.name} — ${this.langName(r.lib.lang)}${r.lib.kind === 'floor' ? ', floor' : ''}`,
        rows: [['CPU / wall', Fmt.pct(r.cpu / r.wall)], ['Wall time', Fmt.ns(r.wall)], ['Allocated', PerfTable.allocText(r)]]
          .concat(ref && r.id !== ref.id ? [['vs ZeRechnung', `${Fmt.ratio(r.cpu / ref.cpu)} CPU, ${Fmt.ratio(r.rss / ref.rss)} memory`]] : []),
      },
    }));

    const legend = h('div', { class: 'flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-ink-2' },
      this.legendDot('var(--accent)', 'ZeRechnung'),
      this.legendDot('var(--other)', 'Factur-X library'),
      this.state.floors ? this.legendDot('var(--other)', 'Floor: generic PDF or XML library', true) : null);
    this.chartCard.replaceChildren(
      h('div', { class: 'flex flex-wrap items-start justify-between gap-3' },
        h('div', {},
          h('h3', { class: 'text-sm font-semibold', text: `${benchInfo.op.title}, ${benchInfo.label}: CPU time and peak memory` }),
          h('p', { class: 'mt-1 text-xs text-ink-2' }, `${points.length} libraries, one process each. Lower left uses less of both. Log scales.`,
            benchInfo.op.id === 'validate' ? [' The rule sets differ widely: ', h('a', { class: 'link', href: 'compliance.html', text: 'what each validator checks' }), '.'] : null)),
        legend),
    );
    this.scatter = ScatterChart.render(this.chartCard, {
      points, xDomain: cpuDomain, yDomain: rssDomain, xTicks: cpuTicks, yTicks: rssTicks,
      xFormat: Fmt.nsTick, yFormat: Fmt.mibTick, xTitle: 'CPU time per call', yTitle: 'Peak memory',
      ariaLabel: `Scatter plot of CPU time per call against peak memory for ${benchInfo.op.title}, ${benchInfo.label}. The table below lists every value.`,
      onHover: (id) => this.table && this.table.highlight(id),
    });
    this.scatter.classList.add('mt-4');

    const failed = rows.filter((r) => r.status !== 'ok').length;
    this.tableCard.replaceChildren(
      h('div', { class: 'flex flex-wrap items-baseline justify-between gap-2 px-3 pt-2 pb-3' },
        h('h3', { class: 'text-sm font-semibold', text: 'Every library' }),
        h('p', { class: 'text-xs text-muted', text: `${rows.length - failed} measured${failed ? `, ${failed} failed` : ''} · bars on a log scale · (×n) relative to ZeRechnung · † child processes not counted` })),
    );
    this.table = PerfTable.render(this.tableCard, {
      rows, ref, cpuDomain, rssDomain, langName: this.langName,
      onHover: (id) => this.scatter && this.scatter.highlight(id),
    });
  },

  renderTiles(rows, ref) {
    const ranked = rows.filter((r) => r.id !== 'go/ZeRechnung, GOMAXPROCS=1' && !r.child);
    const byCpu = [...ranked].sort((a, b) => a.cpu - b.cpu);
    const byRss = [...ranked].sort((a, b) => a.rss - b.rss);
    const n = ranked.length;
    const rank = (list) => (ref ? list.findIndex((r) => r.id === ref.id) + 1 : 0);
    const ord = (k) => {
      if (k === 1) return 'Lowest';
      const suffix = k % 100 >= 11 && k % 100 <= 13 ? 'th' : { 1: 'st', 2: 'nd', 3: 'rd' }[k % 10] || 'th';
      return `${k}${suffix} lowest`;
    };
    const who = (r) => (r ? `${r.lib.name} · ${this.langName(r.lib.lang)}${r.lib.kind === 'floor' ? ' · floor' : ''}` : '');
    this.tiles.replaceChildren(
      StatTile.create({ label: 'ZeRechnung · CPU time per call', value: ref ? Fmt.ns(ref.cpu) : '—', sub: ref ? `${ord(rank(byCpu))} of ${n}` : 'Go is filtered out', emphasis: true }),
      StatTile.create({ label: 'ZeRechnung · peak memory', value: ref ? Fmt.mib(ref.rss) : '—', sub: ref ? `${ord(rank(byRss))} of ${n}` : 'Go is filtered out', emphasis: true }),
      StatTile.create({ label: 'Least CPU time', value: byCpu[0] ? Fmt.ns(byCpu[0].cpu) : '—', sub: who(byCpu[0]) }),
      StatTile.create({ label: 'Least peak memory', value: byRss[0] ? Fmt.mib(byRss[0].rss) : '—', sub: who(byRss[0]) }),
    );
  },

  legendDot(color, text, hollow = false) {
    return h('span', { class: 'inline-flex items-center gap-1.5' },
      svg('svg', { width: 12, height: 12, 'aria-hidden': 'true' },
        svg('circle', { cx: 6, cy: 6, r: hollow ? 4 : 5, fill: hollow ? 'var(--surface)' : color, stroke: hollow ? color : 'none', 'stroke-width': 2 })),
      text);
  },

  // --- Method -------------------------------------------------------------

  methodSection() {
    const li = (...c) => h('li', {}, ...c);
    return h('section', { class: 'pb-6', 'aria-labelledby': 'method' },
      h('div', { class: 'card p-5 sm:p-6' },
        h('h2', { id: 'method', class: 'text-base font-semibold', text: 'How to read these figures' }),
        h('ul', { class: 'mt-3 grid gap-x-8 gap-y-2 text-sm text-ink-2 leading-relaxed md:grid-cols-2 list-disc pl-5 prose-site' },
          li('CPU time is user + system time of the whole process over the timed series, divided by the number of calls. It can exceed the wall-clock time when other threads work: Go’s concurrent garbage collector, the JVM’s JIT and G1, the .NET thread pool.'),
          li(h('em', { text: 'ZeRechnung, 1 core' }), ' is the same package run again with ', h('code', { text: 'GOMAXPROCS=1' }), ': its garbage collector then shares the core of the measured thread.'),
          li('The peak includes the runtime. On the JVM it mostly measures the size of the G1 heap, which grows under allocation pressure, not the live data.'),
          li('A fresh process is not always at full speed: JIT or tiered compilation can still run during the one-second series. Response times measured in a warm process are in ', fileLink('README.md#response-time', 'README.md'), '.'),
          li('† The library spawns child processes on every call (poppler’s ', h('code', { text: 'pdfdetach' }), '); their CPU and memory are not counted, so it is left out of the chart and of the ranks.'),
          li('A floor is a generic PDF or XML library wired in a few lines: it knows nothing about Factur-X and gives the minimum cost of the operation in its language.'),
        ),
        h('p', { class: 'mt-4 text-xs text-muted prose-site' }, 'Full protocol: ', fileLink('PROTOCOL.md', 'PROTOCOL.md'), '. Every table per operation: ',
          fileLink(`results/${this.d.run.date}/tables.md`, 'tables.md'), '.')));
  },
};

document.addEventListener('DOMContentLoaded', () => Page.boot({ active: 'perf', data: 'perf', app: PerfApp }));
