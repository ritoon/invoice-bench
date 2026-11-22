/**
 * Tooltip — one floating readout shared by every chart. Values lead,
 * labels follow; text is always set with textContent.
 * Usage: Tooltip.show(anchor, { title, sub, rows: [[label, value]], note })
 *        anchor is a pointer event or an element (keyboard focus).
 */
const Tooltip = {
  el: null,

  create() {
    const el = h('div', { class: 'tooltip', role: 'tooltip', 'aria-hidden': 'true' });
    document.body.appendChild(el);
    return el;
  },

  show(anchor, { title, sub, rows = [], note } = {}) {
    if (!this.el) this.el = this.create();
    // replaceChildren would print a null as "null": keep the nodes only.
    this.el.replaceChildren(...[
      title ? h('div', { class: 'tt-title', text: title }) : null,
      sub ? h('div', { class: 'tt-sub', text: sub }) : null,
      rows.length ? h('div', { class: 'tt-rows' }, rows.map(([k, v]) => [h('span', { class: 'tt-k', text: k }), h('span', { class: 'tt-v', text: v })])) : null,
      note ? h('div', { class: 'tt-note', text: note }) : null,
    ].filter(Boolean));
    this.el.classList.add('is-visible');
    this.place(anchor);
  },

  place(anchor) {
    const pad = 14;
    let x, y;
    if (anchor && 'clientX' in anchor) {
      x = anchor.clientX;
      y = anchor.clientY;
    } else if (anchor && anchor.getBoundingClientRect) {
      const r = anchor.getBoundingClientRect();
      x = r.left + r.width / 2;
      y = r.top;
    } else return;
    const tw = this.el.offsetWidth, th = this.el.offsetHeight;
    let left = x + pad, top = y - th - pad;
    if (left + tw > window.innerWidth - 8) left = x - tw - pad;
    if (left < 8) left = 8;
    if (top < 8) top = y + pad;
    this.el.style.left = `${left}px`;
    this.el.style.top = `${top}px`;
  },

  hide() {
    if (this.el) this.el.classList.remove('is-visible');
  },
};
