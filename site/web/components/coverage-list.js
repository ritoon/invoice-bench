/**
 * CoverageList — how many files of the corpus a library handles. The
 * libraries that handle every file are listed by name; the others get a
 * meter each (a ratio against the same whole), worst last.
 *
 * Usage: CoverageList.render(container, { title, unit, total, items: [{ id, name, lang, kind, ok }], langName })
 */
const CoverageList = {
  create({ title, unit, total, items, langName }) {
    const complete = items.filter((i) => i.ok >= total);
    const partial = items.filter((i) => i.ok < total).sort((a, b) => b.ok - a.ok || a.name.localeCompare(b.name));
    const pill = (i) => h('li', { class: `inline-flex items-baseline gap-1.5 rounded-full px-2.5 py-0.5 text-xs ${i.kind === 'self' ? 'bg-accent-wash text-accent-ink font-semibold' : 'bg-surface-2 text-ink'}` },
      h('span', { text: i.name }), h('span', { class: 'text-muted font-normal', text: langName(i.lang) }));

    return h('div', { class: 'card p-5' },
      h('h3', { class: 'text-sm font-semibold', text: title }),
      h('p', { class: 'mt-1 text-xs text-ink-2', text: `${complete.length} of ${items.length} libraries handle all ${total} ${unit}.` }),
      complete.length ? h('div', { class: 'mt-4' },
        h('p', { class: 'eyebrow', text: `All ${total}` }),
        h('ul', { class: 'mt-2 flex flex-wrap gap-1.5' }, complete.map(pill))) : null,
      partial.length ? h('div', { class: 'mt-5' },
        h('p', { class: 'eyebrow', text: 'Fewer' }),
        h('ul', { class: 'mt-2 space-y-2' }, partial.map((i) =>
          h('li', { class: 'grid grid-cols-[minmax(0,11rem)_1fr_auto] items-center gap-3 text-sm' },
            h('span', { class: 'truncate', title: i.name },
              h('span', { class: i.kind === 'self' ? 'font-semibold' : 'font-medium', text: i.name }), ' ',
              h('span', { class: 'tag lang', text: langName(i.lang) })),
            h('div', { class: 'meter', role: 'meter', 'aria-valuemin': 0, 'aria-valuemax': total, 'aria-valuenow': i.ok, 'aria-label': `${i.name}: ${i.ok} of ${total}` },
              h('span', { style: { width: `${(100 * i.ok) / total}%` } })),
            h('span', { class: 'num text-xs text-ink-2 w-14 text-right', text: `${i.ok}/${total}` }))))) : null);
  },

  render(container, props) {
    const el = this.create(props);
    container.appendChild(el);
    return el;
  },
};
