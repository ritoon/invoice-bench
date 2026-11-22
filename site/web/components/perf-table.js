/**
 * PerfTable — the table view of the scatter, and more: every library that
 * runs the benchmark, sortable, with inline bars on a log scale. It is the
 * way to read every value without hovering.
 *
 * Usage: const table = PerfTable.render(container, {
 *   rows: [{ id, lib, status, wall, cpu, rss, bytes, allocs, lower, child, note }],
 *   ref,            // the ZeRechnung row, for the (×n) ratios
 *   cpuDomain, rssDomain, langName, onHover(id | null) });
 *   table.highlight(id | null)
 */
const PerfTable = {
  columns: [
    { key: 'name', text: 'Library', align: 'l' },
    { key: 'cpu', text: 'CPU time per call', align: 'l', title: 'User + system time of the process, all threads, divided by the number of calls' },
    { key: 'use', text: 'CPU / wall', align: 'r', title: '100% = one core busy during the timed series; more = parallel work (garbage collector, JIT)' },
    { key: 'rss', text: 'Peak memory', align: 'l', title: 'Maximum resident set size of the process, runtime included' },
    { key: 'wall', text: 'Wall time', align: 'r', title: 'Median wall-clock time per call, in the same process' },
    { key: 'bytes', text: 'Allocated per call', align: 'r', title: 'Bytes · allocations of one call (alloc pass; Go from the benchmark). Not counted yet in Java, C#, C, C++ and PHP' },
  ],

  create(props) {
    const state = { sort: 'cpu', dir: 1, active: null };
    const table = h('table', { class: 'data-table sticky-first' });
    const wrap = h('div', { class: 'overflow-x-auto' }, table);
    const paint = () => table.replaceChildren(this.head(state, paint), this.body(state, props));
    paint();
    wrap.highlight = (id) => {
      state.active = id;
      table.querySelectorAll('tbody tr').forEach((tr) => tr.classList.toggle('is-active', tr.dataset.id === id));
    };
    return wrap;
  },

  value(r, key) {
    switch (key) {
      case 'name': return r.lib.name.toLowerCase();
      case 'use': return r.cpu / r.wall;
      default: return r[key] || Infinity;
    }
  },

  head(state, paint) {
    return h('thead', {}, h('tr', {}, this.columns.map((c) => {
      const sorted = state.sort === c.key;
      const btn = h('button', { type: 'button', class: 'sort', title: c.title || null },
        c.text, h('span', { class: 'arrow', 'aria-hidden': 'true', text: sorted && state.dir < 0 ? '▼' : '▲' }));
      btn.addEventListener('click', () => {
        if (state.sort === c.key) state.dir = -state.dir;
        else {
          state.sort = c.key;
          state.dir = 1;
        }
        paint();
      });
      return h('th', {
        class: c.align === 'r' ? 'r' : null, scope: 'col',
        'aria-sort': sorted ? (state.dir > 0 ? 'ascending' : 'descending') : null,
      }, btn);
    })));
  },

  body(state, { rows, ref, cpuDomain, rssDomain, langName, onHover }) {
    const ok = rows.filter((r) => r.status === 'ok');
    const failed = rows.filter((r) => r.status !== 'ok');
    ok.sort((a, b) => {
      const va = this.value(a, state.sort), vb = this.value(b, state.sort);
      return (va < vb ? -1 : va > vb ? 1 : 0) * state.dir;
    });
    failed.sort((a, b) => a.lib.name.localeCompare(b.lib.name));
    const bar = (v, dom) => {
      const s = Scale.log(dom[0], dom[1], 0, 100);
      return h('div', { class: 'ibar', 'aria-hidden': 'true' }, h('span', { style: { width: `${Math.max(2, Math.min(100, s(v)))}%` } }));
    };
    const ratio = (v, rv) => (ref && rv > 0 && v > 0 && ref.id !== undefined ? h('span', { class: 'ml-1 text-xs text-muted', text: Fmt.ratio(v / rv) }) : null);

    const tbody = h('tbody', {});
    for (const r of [...ok, ...failed]) {
      const isRef = ref && r.id === ref.id;
      const self = r.lib.kind === 'self';
      const tr = h('tr', { class: self ? 'self' : null, tabindex: 0, 'data-id': r.id },
        h('td', { class: 'min-w-[10rem] sm:min-w-[12rem]' },
          h('div', { class: 'flex flex-wrap items-baseline gap-x-2 gap-y-0.5' },
            h('span', { class: self ? 'font-semibold' : 'font-medium', text: r.lib.name + (r.child ? ' †' : '') }),
            h('span', { class: 'tag lang', text: langName(r.lib.lang) }),
            r.lib.kind === 'floor' ? h('span', { class: 'tag', text: 'floor' }) : null,
            r.lib.kind === 'xslt' ? h('span', { class: 'tag', text: 'official XSLT' }) : null,
            r.lib.version ? h('span', { class: 'text-xs text-muted num', text: r.lib.version }) : null)),
      );
      if (r.status === 'ok') {
        tr.append(
          h('td', { class: 'min-w-[13rem]' }, h('div', { class: 'flex items-center gap-3' },
            bar(r.cpu, cpuDomain),
            h('div', { class: 'num whitespace-nowrap text-right w-28 shrink-0' }, h('span', { text: Fmt.ns(r.cpu) }), isRef ? null : ratio(r.cpu, ref && ref.cpu)))),
          h('td', { class: 'r num', text: Fmt.pct(r.cpu / r.wall) }),
          h('td', { class: 'min-w-[13rem]' }, h('div', { class: 'flex items-center gap-3' },
            bar(r.rss, rssDomain),
            h('div', { class: 'num whitespace-nowrap text-right w-28 shrink-0' }, h('span', { text: Fmt.mib(r.rss) }), isRef ? null : ratio(r.rss, ref && ref.rss)))),
          h('td', { class: 'r num whitespace-nowrap', text: Fmt.ns(r.wall) }),
          h('td', { class: 'r num whitespace-nowrap', text: this.allocText(r) }),
        );
      } else {
        tr.append(h('td', { colspan: 5, class: 'text-xs text-ink-2' },
          h('span', { class: 'font-semibold text-critical', text: 'Failed' }), r.note ? ` — ${r.note}` : ''));
      }
      const tip = this.tip(r, langName);
      tr.addEventListener('pointerenter', () => onHover && onHover(r.id));
      tr.addEventListener('pointerleave', () => onHover && onHover(null));
      tr.addEventListener('focus', () => {
        if (onHover) onHover(r.id);
        if (r.note) Tooltip.show(tr.firstChild, tip);
      });
      tr.addEventListener('blur', () => {
        if (onHover) onHover(null);
        Tooltip.hide();
      });
      if (r.note) tr.title = r.note;
      if (r.id === state.active) tr.classList.add('is-active');
      tbody.append(tr);
    }
    return tbody;
  },

  allocText(r) {
    if (!(r.bytes > 0)) return '—';
    let s = (r.lower ? '≥ ' : '') + Fmt.bytes(r.bytes);
    if (r.allocs > 0) s += ' · ' + Fmt.int(r.allocs);
    return s;
  },

  tip(r, langName) {
    return { title: r.lib.name, sub: langName(r.lib.lang), note: r.note };
  },

  render(container, props) {
    const el = this.create(props);
    container.appendChild(el);
    return el;
  },
};
