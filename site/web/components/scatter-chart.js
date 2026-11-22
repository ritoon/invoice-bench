/**
 * ScatterChart — two measures of the same calls, each on its own log axis
 * (one x, one y: never two scales on one axis). Emphasis form: ZeRechnung in
 * the accent with a direct label, every other library in grey; floors are
 * hollow rings (shape, not colour). The pointer only has to be closest to a
 * point, within 24px, not on it.
 *
 * Usage: const chart = ScatterChart.render(container, {
 *   points: [{ id, x, y, kind: 'self' | 'floor' | 'library' | 'xslt', label, tip }],
 *   xDomain, yDomain, xTicks, yTicks, xFormat, yFormat, xTitle, yTitle, ariaLabel,
 *   onHover(id | null) });
 *   chart.highlight(id | null)   // from the table
 */
const ScatterChart = {
  create(props) {
    const wrap = h('div', { class: 'chart relative' });
    const state = { active: null, layout: null, width: 0 };
    const redraw = () => wrap.replaceChildren(this.draw(state, props));
    onResize(wrap, (w) => {
      state.width = w;
      redraw();
    });
    wrap.highlight = (id) => {
      state.active = id;
      this.paintActive(wrap, state, props);
    };
    return wrap;
  },

  draw(state, props) {
    const { points, xDomain, yDomain, xTicks, yTicks, xFormat, yFormat, xTitle, yTitle, ariaLabel, onHover } = props;
    const width = state.width;
    const height = Math.round(Math.max(300, Math.min(480, width * 0.56)));
    const m = { l: 62, r: 18, t: 28, b: 42 };
    const x = Scale.log(xDomain[0], xDomain[1], m.l, width - m.r);
    const y = Scale.log(yDomain[0], yDomain[1], height - m.b, m.t);
    state.layout = { x, y, width, height, m };

    const root = svg('svg', { width, height, role: 'img', 'aria-label': ariaLabel });
    // Grid and ticks
    const gap = xTicks.length > 1 ? (width - m.l - m.r) / (xTicks.length - 1) : Infinity;
    const stride = Math.max(1, Math.ceil(52 / gap));
    xTicks.forEach((t, i) => {
      const tx = Math.round(x(t)) + 0.5;
      root.append(svg('line', { class: 'grid', x1: tx, x2: tx, y1: m.t, y2: height - m.b }));
      if (i % stride === 0) root.append(svg('text', { x: tx, y: height - m.b + 16, 'text-anchor': 'middle', text: xFormat(t) }));
    });
    for (const t of yTicks) {
      const ty = Math.round(y(t)) + 0.5;
      root.append(svg('line', { class: 'grid', x1: m.l, x2: width - m.r, y1: ty, y2: ty }),
        svg('text', { x: m.l - 8, y: ty + 4, 'text-anchor': 'end', text: yFormat(t) }));
    }
    root.append(
      svg('line', { class: 'baseline', x1: m.l, x2: width - m.r, y1: height - m.b + 0.5, y2: height - m.b + 0.5 }),
      svg('text', { class: 'label', x: width - m.r, y: height - 6, 'text-anchor': 'end', text: `${xTitle} →` }),
      svg('text', { class: 'label', x: m.l - 8, y: 12, 'text-anchor': 'start', text: `↑ ${yTitle}` }),
    );

    // Marks: context first, emphasis last so it sits on top.
    const order = { library: 0, xslt: 0, floor: 1, self: 2 };
    const sorted = [...points].sort((a, b) => order[a.kind] - order[b.kind]);
    const marks = svg('g', {});
    for (const p of sorted) {
      const cx = x(p.x), cy = y(p.y);
      if (p.kind === 'self') {
        marks.append(svg('circle', { cx, cy, r: 6.5, fill: 'var(--accent)', stroke: 'var(--surface)', 'stroke-width': 2 }));
      } else if (p.kind === 'floor') {
        marks.append(svg('circle', { cx, cy, r: 4.5, fill: 'var(--surface)', stroke: 'var(--other)', 'stroke-width': 2 }));
      } else {
        marks.append(svg('circle', { cx, cy, r: 5, fill: 'var(--other)', stroke: 'var(--surface)', 'stroke-width': 2 }));
      }
    }
    root.append(marks);

    // Direct labels on the emphasized points only. The one-core round always
    // uses less CPU and memory, so it is labelled below left, the main one
    // above right.
    for (const p of points.filter((q) => q.kind === 'self' && q.label)) {
      const left = p.labelSide === 'left';
      root.append(svg('text', {
        class: 'label-strong', x: x(p.x) + (left ? -11 : 11), y: y(p.y) + (left ? 16 : -9),
        'text-anchor': left ? 'end' : 'start', text: p.label,
      }));
    }

    const halo = svg('circle', { class: 'halo', r: 9, fill: 'none', stroke: 'var(--ink)', 'stroke-width': 2, visibility: 'hidden', 'pointer-events': 'none' });
    root.append(halo);

    // Nearest-point hover layer.
    const overlay = svg('rect', { class: 'hit', x: m.l, y: m.t, width: Math.max(0, width - m.l - m.r), height: Math.max(0, height - m.t - m.b) });
    const nearest = (evt) => {
      const r = root.getBoundingClientRect();
      const px = evt.clientX - r.left, py = evt.clientY - r.top;
      let best = null, bd = 24 * 24;
      for (const p of points) {
        const d = (x(p.x) - px) ** 2 + (y(p.y) - py) ** 2;
        if (d < bd) {
          bd = d;
          best = p;
        }
      }
      return best;
    };
    overlay.addEventListener('pointermove', (e) => {
      const p = nearest(e);
      state.active = p ? p.id : null;
      this.paintActive(root.parentNode, state, props);
      if (p) Tooltip.show(e, p.tip);
      else Tooltip.hide();
      if (onHover) onHover(state.active);
    });
    overlay.addEventListener('pointerleave', () => {
      state.active = null;
      this.paintActive(root.parentNode, state, props);
      Tooltip.hide();
      if (onHover) onHover(null);
    });
    root.append(overlay);
    requestAnimationFrame(() => this.paintActive(root.parentNode, state, props));
    return root;
  },

  paintActive(wrap, state, props) {
    if (!wrap || !state.layout) return;
    const halo = wrap.querySelector('circle.halo');
    if (!halo) return;
    const p = props.points.find((q) => q.id === state.active);
    if (!p) {
      halo.setAttribute('visibility', 'hidden');
      return;
    }
    halo.setAttribute('cx', state.layout.x(p.x));
    halo.setAttribute('cy', state.layout.y(p.y));
    halo.setAttribute('visibility', 'visible');
  },

  render(container, props) {
    const el = this.create(props);
    container.appendChild(el);
    return el;
  },
};
