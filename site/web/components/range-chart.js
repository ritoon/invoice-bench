/**
 * RangeChart — one row per category: a hairline from min to max and a dot
 * at the median, on a log or linear axis. Emphasis, not categorical colour:
 * the emphasized row takes the accent, the others the de-emphasis grey.
 *
 * Usage: RangeChart.render(container, {
 *   rows: [{ id, label, sub, min, med, max, emphasis, tip: { title, sub, rows } }],
 *   scale: 'log' | 'linear', domain: [lo, hi], ticks: [..], tickFormat, valueFormat,
 *   ariaLabel })
 */
const RangeChart = {
  create(props) {
    const wrap = h('div', { class: 'chart' });
    onResize(wrap, (w) => wrap.replaceChildren(this.draw(w, props)));
    return wrap;
  },

  draw(width, { rows, scale, domain, ticks, tickFormat, valueFormat, ariaLabel }) {
    const narrow = width < 520;
    const labelW = narrow ? 92 : 132;
    const valueW = narrow ? 64 : 84;
    const rowH = 34;
    const top = 26;
    const plotL = labelW, plotR = width - valueW - 8;
    const height = top + rows.length * rowH + 6;
    const x = (scale === 'log' ? Scale.log : Scale.linear)(domain[0], domain[1], plotL, plotR);

    const root = svg('svg', { width, height, role: 'list', 'aria-label': ariaLabel });
    // Axis: ticks on top, hairline gridlines down the plot. Labels thin out
    // when the ticks get closer than a label's width.
    const gap = ticks.length > 1 ? (plotR - plotL) / (ticks.length - 1) : Infinity;
    const stride = Math.max(1, Math.ceil(50 / gap));
    ticks.forEach((t, i) => {
      const tx = Math.round(x(t)) + 0.5;
      root.append(svg('line', { class: 'grid', x1: tx, x2: tx, y1: top - 6, y2: height - 4 }));
      if (i % stride === 0) root.append(svg('text', { x: tx, y: top - 12, 'text-anchor': 'middle', text: tickFormat(t) }));
    });
    root.append(svg('text', { x: width - 8, y: top - 12, 'text-anchor': 'end', text: 'median' }));

    rows.forEach((r, i) => {
      const y = top + i * rowH + rowH / 2;
      const color = r.emphasis ? 'var(--accent)' : 'var(--other)';
      const g = svg('g', { class: 'row', tabindex: 0, role: 'listitem', 'aria-label': `${r.label}: median ${valueFormat(r.med)}, from ${valueFormat(r.min)} to ${valueFormat(r.max)}` });
      g.append(
        svg('rect', { class: 'row-hover', x: 0, y: y - rowH / 2 + 1, width, height: rowH - 2, rx: 6 }),
        svg('text', { class: r.emphasis ? 'label-strong' : 'label', x: 8, y: y - (r.sub && !narrow ? 2 : -4), text: r.label }),
        r.sub && !narrow ? svg('text', { x: 8, y: y + 11, text: r.sub }) : null,
        svg('line', { x1: x(r.min), x2: x(r.max), y1: y, y2: y, stroke: color, 'stroke-width': 2, 'stroke-linecap': 'round' }),
        svg('rect', { x: x(r.min) - 1, y: y - 5, width: 2, height: 10, rx: 1, fill: color }),
        svg('rect', { x: x(r.max) - 1, y: y - 5, width: 2, height: 10, rx: 1, fill: color }),
        svg('circle', { cx: x(r.med), cy: y, r: r.emphasis ? 6 : 5, fill: color, stroke: 'var(--surface)', 'stroke-width': 2 }),
        svg('text', { class: r.emphasis ? 'label-strong' : 'label', x: width - 8, y: y + 4, 'text-anchor': 'end', text: valueFormat(r.med) }),
      );
      const show = (e) => Tooltip.show(e, r.tip);
      g.addEventListener('pointermove', show);
      g.addEventListener('pointerleave', () => Tooltip.hide());
      g.addEventListener('focus', () => Tooltip.show(g, r.tip));
      g.addEventListener('blur', () => Tooltip.hide());
      root.append(g);
    });
    return root;
  },

  render(container, props) {
    const el = this.create(props);
    container.appendChild(el);
    return el;
  },
};
