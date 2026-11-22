/**
 * ConfidenceTable — the confidence score of each library for one country,
 * ranked, with an inline bar on the 0–100 scale (one hue, ZeRechnung in the
 * accent), the two factors of the score and the score of each part. It is
 * also the table view: every figure is printed.
 *
 * Usage: ConfidenceTable.render(container, { country, libs, langName, limit })
 */
const ConfidenceTable = {
  create({ country, libs, langName, limit = 25 }) {
    const state = { sort: 'score', dir: -1, all: false };
    const wrap = h('div', {});
    const table = h('table', { class: 'data-table sticky-first' });
    const more = h('button', { type: 'button', class: 'mt-3 ml-3 text-sm font-medium text-accent-ink hover:underline' });
    const paint = () => {
      table.replaceChildren(this.head(country, state, paint), this.body(country, libs, langName, state, limit));
      const hidden = libs.length - limit;
      more.hidden = hidden <= 0;
      more.textContent = state.all ? 'Show fewer' : `Show all ${libs.length} libraries`;
    };
    more.addEventListener('click', () => {
      state.all = !state.all;
      paint();
    });
    paint();
    wrap.append(h('div', { class: 'overflow-x-auto' }, table), more);
    return wrap;
  },

  columns(country) {
    // "XRechnung 3.0.2" → "XRechnung", "EN 16931 core rules" → "EN 16931".
    const short = (p) => `${p.title.replace(' core rules', '').replace(/\s+\d+(\.\d+)+$/, '')} ${p.syntax.toUpperCase()}`;
    return [
      { key: 'name', text: 'Library' },
      { key: 'score', text: `${country.name}, out of 100`, title: '100 × √(S × D), mean of the parts' },
      { key: 's', text: 'Valid accepted', align: 'r', title: 'S: share of the valid cases the library accepts' },
      { key: 'd', text: 'Detection', align: 'r', title: 'D: mean credit of the expected rules (1 raised, 0.5 rejected otherwise, 0 accepted)' },
      ...country.parts.map((p, i) => ({ key: `p${i}`, text: short(p), align: 'r', title: `${p.title}, ${p.syntax.toUpperCase()}` })),
    ];
  },

  value(l, key) {
    if (key === 'name') return l.name.toLowerCase();
    if (key.startsWith('p')) return l.parts[Number(key.slice(1))].score;
    return l[key];
  },

  head(country, state, paint) {
    return h('thead', {}, h('tr', {}, this.columns(country).map((c) => {
      const sorted = state.sort === c.key;
      const btn = h('button', { type: 'button', class: 'sort', title: c.title || null },
        c.text, h('span', { class: 'arrow', 'aria-hidden': 'true', text: sorted && state.dir < 0 ? '▼' : '▲' }));
      btn.addEventListener('click', () => {
        if (state.sort === c.key) state.dir = -state.dir;
        else {
          state.sort = c.key;
          state.dir = c.key === 'name' ? 1 : -1;
        }
        paint();
      });
      return h('th', { class: c.align === 'r' ? 'r' : null, scope: 'col', 'aria-sort': sorted ? (state.dir > 0 ? 'ascending' : 'descending') : null }, btn);
    })));
  },

  body(country, libs, langName, state, limit) {
    const rows = [...libs].sort((a, b) => {
      const va = this.value(a, state.sort), vb = this.value(b, state.sort);
      return (va < vb ? -1 : va > vb ? 1 : 0) * state.dir || b.score - a.score;
    });
    const shown = state.all ? rows : rows.slice(0, limit);
    const pct = (v) => Fmt.decimal(v * 100, 0) + '%';
    return h('tbody', {}, shown.map((l) => {
      const self = l.kind === 'self';
      return h('tr', { class: self ? 'self' : null },
        h('td', { class: 'min-w-[10rem] sm:min-w-[13rem]' },
          h('div', { class: 'flex flex-wrap items-baseline gap-x-2 gap-y-0.5' },
            h('span', { class: self ? 'font-semibold' : 'font-medium', text: l.name }),
            h('span', { class: 'tag lang', text: langName(l.lang) }),
            l.kind === 'floor' ? h('span', { class: 'tag', text: 'floor' }) : null,
            l.kind === 'xslt' ? h('span', { class: 'tag', text: 'official XSLT' }) : null)),
        h('td', { class: 'min-w-[9rem]' }, h('div', { class: 'flex items-center gap-3' },
          h('div', { class: 'ibar', role: 'meter', 'aria-valuemin': 0, 'aria-valuemax': 100, 'aria-valuenow': Math.round(l.score), 'aria-label': `${l.name}: ${Math.round(l.score)} out of 100` },
            h('span', { style: { width: `${Math.max(0.5, l.score)}%` } })),
          h('span', { class: `num w-8 text-right shrink-0 ${self ? 'font-semibold' : ''}`, text: Fmt.decimal(l.score, 0) }))),
        h('td', { class: 'r num', text: pct(l.s) }),
        h('td', { class: 'r num', text: pct(l.d) }),
        l.parts.map((p) => h('td', { class: `r num ${p.run ? (p.score < 0.5 ? 'text-muted' : '') : 'text-muted'}`, title: p.run ? `S ${pct(p.s)} · D ${pct(p.d)}` : 'Not run on this suite: counts 0', text: p.run ? Fmt.decimal(p.score, 0) : '—' })));
    }));
  },

  render(container, props) {
    const el = this.create(props);
    container.appendChild(el);
    return el;
  },
};
