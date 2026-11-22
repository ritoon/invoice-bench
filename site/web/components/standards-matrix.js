/**
 * StandardsMatrix — which standards each library actually enforces, as
 * surveyed in its code. The level is ordinal and drawn as a Harvey ball
 * (full, half, quarter, empty): shape carries it, one hue only.
 *
 * Usage: StandardsMatrix.render(container, { standards, rows, langName, details })
 */
const Harvey = {
  create(level, size = 16) {
    const c = size / 2, r = size / 2 - 1.5;
    const el = svg('svg', { class: `harvey ${level}`, width: size, height: size, viewBox: `0 0 ${size} ${size}`, 'aria-hidden': 'true' });
    if (level === 'full') el.append(svg('circle', { class: 'fill', cx: c, cy: c, r: r + 0.75 }));
    else {
      el.append(svg('circle', { class: 'ring', cx: c, cy: c, r }));
      if (level === 'partial') el.append(svg('path', { d: `M${c} ${c - r}A${r} ${r} 0 0 1 ${c} ${c + r}Z` }));
      if (level === 'external') el.append(svg('path', { d: `M${c} ${c}L${c} ${c - r}A${r} ${r} 0 0 1 ${c + r} ${c}Z` }));
    }
    return el;
  },
};

const StandardsMatrix = {
  create({ standards, rows, langName, details }) {
    const level = new Map(standards.levels.map((l) => [l.id, l.label]));
    const table = h('table', { class: `data-table std${details ? ' details' : ''}` },
      h('thead', {}, h('tr', {},
        h('th', { scope: 'col', text: 'Library' }),
        standards.columns.map((c) => h('th', { scope: 'col', class: 'text-center !whitespace-normal leading-tight', style: { minWidth: '5.5rem' }, text: c.label })),
        h('th', { scope: 'col', class: 'r', text: 'Enforced' }))),
      h('tbody', {}, rows.map((r) => h('tr', { class: r.self ? 'self' : null },
        h('td', { class: 'min-w-[11rem]' },
          h('div', { class: r.self ? 'font-semibold' : 'font-medium', text: r.name }),
          h('div', { class: 'text-xs text-muted' }, langName(r.lang), r.usage ? ` · ${r.usage}` : '')),
        r.cells.map(([lv, text], i) => h('td', {
          class: 'cell', tabindex: 0, 'data-tip': JSON.stringify([level.get(lv), `${r.name} · ${standards.columns[i].label}`, text]),
          'aria-label': `${standards.columns[i].label}: ${level.get(lv)}. ${text}`,
        }, Harvey.create(lv), h('span', { class: 'cell-text', text }))),
        h('td', { class: 'r num', text: `${r.cells.filter(([lv]) => lv === 'full').length}/${r.cells.length}` })))));

    const show = (td, anchor) => {
      const [title, sub, note] = JSON.parse(td.dataset.tip);
      Tooltip.show(anchor, { title, sub, note });
    };
    table.addEventListener('pointerover', (e) => {
      const td = e.target.closest('td.cell');
      if (td) show(td, e);
    });
    table.addEventListener('pointermove', (e) => {
      if (e.target.closest('td.cell')) Tooltip.place(e);
    });
    table.addEventListener('pointerleave', () => Tooltip.hide());
    table.addEventListener('focusin', (e) => {
      const td = e.target.closest('td.cell');
      if (td) show(td, td);
    });
    table.addEventListener('focusout', () => Tooltip.hide());
    return h('div', { class: 'overflow-x-auto' }, table);
  },

  legend(standards) {
    return h('ul', { class: 'flex flex-wrap gap-x-5 gap-y-2 text-xs text-ink-2' },
      standards.levels.map((l) => h('li', { class: 'inline-flex items-center gap-2' }, Harvey.create(l.id, 14), l.label)));
  },

  render(container, props) {
    const el = this.create(props);
    container.appendChild(el);
    return el;
  },
};
