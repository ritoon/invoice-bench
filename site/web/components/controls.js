/**
 * Segmented — a row of mutually exclusive buttons.
 * Usage: Segmented.render(container, { label, options: [{ value, text }], value, onChange, small })
 */
const Segmented = {
  create({ label, options, value, onChange, small = false }) {
    const group = h('div', { class: `segmented${small ? ' small' : ''}`, role: 'group', 'aria-label': label });
    const buttons = options.map((o) => {
      const b = h('button', { type: 'button', 'aria-pressed': String(o.value === value), text: o.text });
      b.addEventListener('click', () => {
        buttons.forEach((x) => x.setAttribute('aria-pressed', String(x === b)));
        onChange(o.value);
      });
      return b;
    });
    group.append(...buttons);
    return group;
  },

  render(container, props) {
    const el = this.create(props);
    container.appendChild(el);
    return el;
  },
};

/**
 * LangChips — one toggle per language; the colour never encodes the
 * language, the label does.
 * Usage: LangChips.render(container, { langs: [{ id, name, count }], selected: Set, onChange })
 */
const LangChips = {
  create({ langs, selected, onChange }) {
    const wrap = h('div', { class: 'flex flex-wrap items-center gap-1.5', role: 'group', 'aria-label': 'Languages' });
    const chips = langs.map((l) => {
      const c = h('button', { type: 'button', class: 'chip', 'aria-pressed': String(selected.has(l.id)) },
        svg('svg', { class: 'check', viewBox: '0 0 12 12', 'aria-hidden': 'true' }, svg('path', { d: 'M2.5 6.2 5 8.5l4.5-5' })),
        h('span', { text: l.name }),
        l.count != null ? h('span', { class: 'num text-muted', text: String(l.count) }) : null);
      c.addEventListener('click', () => {
        if (selected.has(l.id)) selected.delete(l.id);
        else selected.add(l.id);
        c.setAttribute('aria-pressed', String(selected.has(l.id)));
        onChange(selected);
      });
      c.dataset.lang = l.id;
      return c;
    });
    // One button selects every language, or none once they all are.
    const all = h('button', { type: 'button', class: 'text-xs font-medium text-ink-2 hover:text-ink px-2 py-1' });
    const label = () => {
      all.textContent = langs.every((l) => selected.has(l.id)) ? 'None' : 'All';
    };
    all.addEventListener('click', () => {
      const every = langs.every((l) => selected.has(l.id));
      langs.forEach((l) => (every ? selected.delete(l.id) : selected.add(l.id)));
      chips.forEach((c) => c.setAttribute('aria-pressed', String(selected.has(c.dataset.lang))));
      label();
      onChange(selected);
    });
    chips.forEach((c) => c.addEventListener('click', label));
    label();
    wrap.append(...chips, all);
    return wrap;
  },

  render(container, props) {
    const el = this.create(props);
    container.appendChild(el);
    return el;
  },
};

/**
 * Toggle — a labelled checkbox switch.
 * Usage: Toggle.render(container, { label, checked, onChange })
 */
const Toggle = {
  create({ label, checked, onChange }) {
    const input = h('input', { type: 'checkbox', class: 'h-4 w-4 accent-[var(--accent)]' });
    input.checked = checked;
    input.addEventListener('change', () => onChange(input.checked));
    return h('label', { class: 'inline-flex items-center gap-2 text-xs font-medium text-ink-2 cursor-pointer select-none' }, input, label);
  },

  render(container, props) {
    const el = this.create(props);
    container.appendChild(el);
    return el;
  },
};

/**
 * Field — a small label above a control, for the filter row.
 */
const Field = {
  create({ label, control }) {
    return h('div', { class: 'flex flex-col gap-1.5' }, h('span', { class: 'eyebrow', text: label }), control);
  },
};

/**
 * StatTile — label, value, and a line of context.
 * Usage: StatTile.render(container, { label, value, sub, emphasis })
 */
const StatTile = {
  create({ label, value, sub, emphasis = false }) {
    return h('div', { class: `card p-4 ${emphasis ? 'ring-1 ring-inset ring-[var(--accent-track)]' : ''}` },
      h('div', { class: 'text-xs font-medium text-ink-2', text: label }),
      h('div', { class: 'tile-value mt-1 truncate', text: value }),
      sub ? h('div', { class: 'mt-0.5 text-xs text-muted', text: sub }) : null);
  },

  render(container, props) {
    const el = this.create(props);
    container.appendChild(el);
    return el;
  },
};
