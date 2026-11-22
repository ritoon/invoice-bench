/**
 * AgreementMatrix — one row per validator, one column per valid invoice of
 * the corpus. Each cell compares what the validator flags with what the
 * official schematrons flag. States are told apart by shape first, colour
 * second, and the legend names them:
 *
 *   match  filled square  flags a file the official schematrons flag
 *   miss   hollow square  does not flag a file the official schematrons flag
 *   extra  red diamond    flags a file the official schematrons accept
 *   error  red cross      throws an exception on the file
 *   clean  dot            accepts a file the official schematrons accept
 *
 * The comparison is by file: libraries name the same check differently (the
 * FeRD schematron recodes BR-CL-25 as FX-SCH-A-000571), so the rules they
 * report are shown on hover, next to the official ones, not compared.
 *
 * Usage: AgreementMatrix.render(container, { files, reference, rows, groups, langName })
 *   rows: [{ id, name, lang, kind, cells: [{ state, flag }], flagged, verdict, group }]
 */
const AgreementMatrix = {
  states: {
    match: 'Flagged, like the official schematrons',
    miss: 'Not flagged, unlike the official schematrons',
    extra: 'False positive: flagged, the official schematrons accept it',
    error: 'Exception on a valid file',
    clean: 'Accepted, like the official schematrons',
  },

  legend() {
    return h('ul', { class: 'flex flex-wrap gap-x-5 gap-y-2 text-xs text-ink-2' },
      ['match', 'miss', 'extra', 'error', 'clean'].map((s) =>
        h('li', { class: 'inline-flex items-center gap-2' },
          h('span', { class: 'inline-flex w-3 justify-center' }, h('span', { class: `g g-${s}` })),
          this.states[s].replace(', like the official schematrons', '').replace(', unlike the official schematrons', ''))));
  },

  create({ files, reference, rows, groups, langName }) {
    const refRules = new Map(reference.flags.map((f) => [f.file, f.rules || []]));
    const prefixes = [];
    for (const f of files) {
      const p = f.split('_')[0];
      const last = prefixes[prefixes.length - 1];
      if (last && last.name === p) last.n++;
      else prefixes.push({ name: p, n: 1 });
    }
    const groupName = { EN16931: 'EN 16931 examples', XRECHNUNG: 'XRechnung' };
    const ncols = files.length + 3;

    const thead = h('thead', {},
      h('tr', {},
        h('th', { class: 'name text-left' }),
        prefixes.map((p) => h('th', { colspan: p.n, class: 'pb-1 text-[11px] font-medium text-muted text-left border-b border-axis', text: groupName[p.name] || p.name })),
        h('th', {}), h('th', {})),
      h('tr', {},
        h('th', { class: 'name text-left text-xs font-semibold text-muted pb-1', text: 'Validator' }),
        files.map((f, i) => h('th', { class: `col num text-[10px] font-normal text-muted${refRules.has(f) ? ' ref' : ''}`, 'data-file': i, scope: 'col', 'aria-label': f, text: String(i + 1) })),
        h('th', { class: 'pl-3 text-right text-xs font-semibold text-muted', text: 'Flagged' }),
        h('th', { class: 'pl-4 text-left text-xs font-semibold text-muted', text: 'Verdict' })));

    const tbody = h('tbody', {});
    const lookup = [];
    const addRow = (r, cls) => {
      const idx = lookup.push(r) - 1;
      tbody.append(h('tr', { class: cls, 'data-row': idx },
        h('td', { class: 'name py-1' },
          h('div', { class: 'flex items-baseline gap-2' },
            h('span', { class: `${r.kind === 'self' || r.kind === 'reference' ? 'font-semibold' : 'font-medium'} inline-block max-w-[9rem] truncate align-bottom sm:max-w-none`, title: r.name, text: r.name }),
            r.lang ? h('span', { class: 'tag lang', text: langName(r.lang) }) : null,
            r.kind === 'xslt' ? h('span', { class: 'hidden sm:inline' }, h('span', { class: 'tag', text: 'official XSLT' })) : null,
            r.sub ? h('span', { class: 'hidden text-xs text-muted sm:inline', text: r.sub }) : null)),
        r.cells.map((c, i) => h('td', {
          class: `cell${refRules.has(files[i]) ? ' ref' : ''}`, 'data-file': i,
          'aria-label': `${files[i]}: ${this.states[c.state]}`,
        }, h('span', { class: `g g-${c.state}` }))),
        h('td', { class: 'pl-3 text-right num whitespace-nowrap', text: `${r.flagged}/${files.length}` }),
        h('td', { class: 'pl-4 pr-2 text-xs text-ink-2 whitespace-nowrap', text: r.verdict })));
    };

    addRow(reference.row, 'reference');
    let lastGroup = null;
    for (const r of rows) {
      if (groups && r.groupLabel !== lastGroup) {
        lastGroup = r.groupLabel;
        tbody.append(h('tr', { class: 'group' }, h('td', { colspan: ncols, class: 'pb-1' }, h('span', { class: 'eyebrow sticky left-0', text: r.groupLabel }))));
      }
      addRow(r, r.kind === 'self' ? 'self' : null);
    }

    const table = h('table', { class: 'mx' }, thead, tbody);
    const tipFor = (td) => {
      const tr = td.closest('tr');
      const fi = Number(td.dataset.file);
      if (!tr || Number.isNaN(fi)) return null;
      if (!tr.dataset.row) {
        const rules = refRules.get(files[fi]);
        return { title: files[fi], sub: `File ${fi + 1} of ${files.length}`, rows: rules ? [['Official rules', rules.join(', ') || '—']] : [['Official schematrons', 'accept it']] };
      }
      const r = lookup[Number(tr.dataset.row)];
      const c = r.cells[fi];
      const rows = [];
      if (c.flag && c.flag.rules && c.flag.rules.length) rows.push(['Rules', c.flag.rules.join(', ')]);
      else if (c.flag && c.flag.errors) rows.push(['Errors', String(c.flag.errors)]);
      if (refRules.get(files[fi])) rows.push(['Official rules', refRules.get(files[fi]).join(', ')]);
      return { title: this.states[c.state], sub: `${r.name} · ${files[fi]}`, rows, note: c.flag && c.flag.error ? c.flag.error : null };
    };
    table.addEventListener('pointerover', (e) => {
      const td = e.target.closest('td.cell, th.col');
      if (!td) return;
      const tip = tipFor(td);
      if (tip) Tooltip.show(e, tip);
    });
    table.addEventListener('pointermove', (e) => {
      if (e.target.closest('td.cell, th.col')) Tooltip.place(e);
    });
    table.addEventListener('pointerout', (e) => {
      if (!e.relatedTarget || !table.contains(e.relatedTarget) || !e.relatedTarget.closest('td.cell, th.col')) Tooltip.hide();
    });

    return h('div', { class: 'overflow-x-auto pb-2' }, table);
  },

  render(container, props) {
    const el = this.create(props);
    container.appendChild(el);
    return el;
  },
};
