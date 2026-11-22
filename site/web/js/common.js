/**
 * common.js — helpers shared by both pages: DOM and SVG builders,
 * number formatting (same rules as report/), scales, data loading,
 * theme. Loaded before the components.
 */

/* Tailwind Play CDN reads this global: colours are CSS custom properties,
 * so one class (bg-surface, text-ink-2…) works in both themes. */
window.tailwind = window.tailwind || {};
tailwind.config = {
  theme: {
    extend: {
      colors: {
        page: 'var(--page)',
        surface: 'var(--surface)',
        'surface-2': 'var(--surface-2)',
        ink: 'var(--ink)',
        'ink-2': 'var(--ink-2)',
        muted: 'var(--muted)',
        line: 'var(--line)',
        grid: 'var(--grid)',
        axis: 'var(--axis)',
        accent: 'var(--accent)',
        'accent-ink': 'var(--accent-ink)',
        'accent-wash': 'var(--accent-wash)',
        critical: 'var(--critical)',
      },
      maxWidth: { site: '76rem' },
    },
  },
};

const REPO = 'https://github.com/ritoon/invoice-bench';

/**
 * h builds an HTML element. attrs: class, text, html (trusted only),
 * style (object), on<Event> handlers, and any other attribute.
 * Children: nodes, strings (as text) or arrays of them; null is skipped.
 */
function h(tag, attrs = {}, ...children) {
  const el = document.createElement(tag);
  setAttrs(el, attrs);
  append(el, children);
  return el;
}

/** svg builds an SVG element, with the same attrs as h. */
function svg(tag, attrs = {}, ...children) {
  const el = document.createElementNS('http://www.w3.org/2000/svg', tag);
  setAttrs(el, attrs);
  append(el, children);
  return el;
}

function setAttrs(el, attrs) {
  for (const [k, v] of Object.entries(attrs || {})) {
    if (v === null || v === undefined || v === false) continue;
    if (k === 'class') el.setAttribute('class', v);
    else if (k === 'text') el.textContent = v;
    else if (k === 'html') el.innerHTML = v;
    else if (k === 'style' && typeof v === 'object') Object.assign(el.style, v);
    else if (k.startsWith('on') && typeof v === 'function') el.addEventListener(k.slice(2).toLowerCase(), v);
    else el.setAttribute(k, v === true ? '' : v);
  }
}

function append(el, children) {
  for (const c of children.flat(Infinity)) {
    if (c === null || c === undefined || c === false) continue;
    el.appendChild(typeof c === 'string' || typeof c === 'number' ? document.createTextNode(String(c)) : c);
  }
}

/** Fmt formats numbers like report/ does in tables.md. */
const Fmt = {
  decimal(v, prec) {
    return v.toLocaleString('en-US', { minimumFractionDigits: prec, maximumFractionDigits: prec });
  },
  ns(ns) {
    if (!(ns > 0)) return '—';
    if (ns >= 1e9) return Fmt.decimal(ns / 1e9, 2) + ' s';
    if (ns >= 1e6) return Fmt.decimal(ns / 1e6, 2) + ' ms';
    if (ns >= 1e5) return Fmt.decimal(ns / 1e3, 0) + ' µs';
    if (ns >= 1e3) return Fmt.decimal(ns / 1e3, 1) + ' µs';
    return Fmt.decimal(ns, 0) + ' ns';
  },
  mib(b) {
    if (!(b > 0)) return '—';
    const v = b / 1048576;
    if (v >= 100) return Fmt.decimal(v, 0) + ' MiB';
    if (v >= 1) return Fmt.decimal(v, 1) + ' MiB';
    return Fmt.decimal(v, 2) + ' MiB';
  },
  bytes(b) {
    if (!(b > 0)) return '—';
    if (b >= 1000 * 1048576) return Fmt.decimal(b / 1073741824, 1) + ' GB';
    if (b >= 1000 * 1024) return Fmt.decimal(b / 1048576, 1) + ' MB';
    if (b >= 1024) return Fmt.decimal(b / 1024, 0) + ' KB';
    return Fmt.decimal(b, 0) + ' B';
  },
  pct(r) {
    return Fmt.decimal(r * 100, 0) + '%';
  },
  ratio(r) {
    return '×' + (r >= 10 ? Fmt.decimal(r, 0) : Fmt.decimal(r, 1));
  },
  int(n) {
    return Fmt.decimal(n, 0);
  },
  /** Tick labels on a log time axis (powers of ten of nanoseconds). */
  nsTick(ns) {
    const units = [[1e9, 's'], [1e6, 'ms'], [1e3, 'µs'], [1, 'ns']];
    for (const [d, u] of units) if (ns >= d) return Fmt.decimal(ns / d, 0) + ' ' + u;
    return ns + ' ns';
  },
  /** Tick labels on a log memory axis (powers of two of MiB). */
  mibTick(b) {
    const v = b / 1048576;
    return v >= 1024 ? Fmt.decimal(v / 1024, 0) + ' GiB' : Fmt.decimal(v, v < 1 ? 1 : 0) + ' MiB';
  },
};

/** Scale maps a domain to a range, linearly or on a log scale. */
const Scale = {
  log(d0, d1, r0, r1) {
    const a = Math.log(d0), b = Math.log(d1);
    const f = (v) => r0 + ((Math.log(Math.max(v, d0 / 1e3)) - a) / (b - a)) * (r1 - r0);
    f.domain = [d0, d1];
    return f;
  },
  linear(d0, d1, r0, r1) {
    const f = (v) => r0 + ((v - d0) / (d1 - d0)) * (r1 - r0);
    f.domain = [d0, d1];
    return f;
  },
  /** Powers of base that bracket [min, max]. */
  powers(min, max, base) {
    const lo = Math.floor(Math.log(min) / Math.log(base) + 1e-9);
    const hi = Math.ceil(Math.log(max) / Math.log(base) - 1e-9);
    const out = [];
    for (let e = lo; e <= hi; e++) out.push(Math.pow(base, e));
    return out;
  },
};

/** Data loads the JSON files written by site/main.go. */
const Data = {
  async load(name) {
    const res = await fetch(`data/${name}.json`, { cache: 'no-cache' });
    if (!res.ok) throw new Error(`data/${name}.json: ${res.status}`);
    return res.json();
  },
};

/** Store wraps localStorage, which can be missing or throw. */
const Store = {
  get(k) {
    try {
      return window.localStorage.getItem(k);
    } catch {
      return null;
    }
  },
  set(k, v) {
    try {
      if (v === null) window.localStorage.removeItem(k);
      else window.localStorage.setItem(k, v);
    } catch {
      /* private mode: the preference just isn't remembered */
    }
  },
};

/** Theme: auto (follows the OS), light or dark, on <html data-theme>. */
const Theme = {
  modes: ['auto', 'light', 'dark'],
  get() {
    const t = Store.get('theme');
    return this.modes.includes(t) ? t : 'auto';
  },
  apply(mode) {
    if (mode === 'auto') document.documentElement.removeAttribute('data-theme');
    else document.documentElement.setAttribute('data-theme', mode);
    Store.set('theme', mode === 'auto' ? null : mode);
    document.dispatchEvent(new CustomEvent('themechange', { detail: mode }));
  },
  next() {
    const m = this.modes[(this.modes.indexOf(this.get()) + 1) % this.modes.length];
    this.apply(m);
    return m;
  },
};
Theme.apply(Theme.get());

/**
 * onResize calls fn(width) now and whenever the element's width changes,
 * so charts are drawn at the size they are shown, text unstretched.
 */
function onResize(el, fn) {
  let last = -1;
  const run = () => {
    const w = Math.round(el.clientWidth);
    if (w > 0 && w !== last) {
      last = w;
      fn(w);
    }
  };
  if ('ResizeObserver' in window) new ResizeObserver(() => requestAnimationFrame(run)).observe(el);
  else window.addEventListener('resize', run);
  run();
}

function median(v) {
  if (!v.length) return 0;
  const s = [...v].sort((a, b) => a - b);
  const n = s.length;
  return n % 2 ? s[(n - 1) / 2] : (s[n / 2 - 1] + s[n / 2]) / 2;
}

function fileLink(path, text) {
  return h('a', { class: 'link', href: `${REPO}/blob/main/${path}`, text: text || path });
}

/** Page boots the shared chrome, then hands the data to the page's App. */
const Page = {
  async boot({ active, data, app }) {
    const header = document.getElementById('site-header');
    const footer = document.getElementById('site-footer');
    SiteHeader.render(header, { active });
    const main = document.getElementById('app');
    try {
      const d = await Data.load(data);
      SiteFooter.render(footer, { run: d.run });
      app.init(main, d);
    } catch (err) {
      main.replaceChildren(
        h('div', { class: 'card p-6 text-sm text-ink-2' },
          h('p', { class: 'font-semibold text-ink', text: 'The data could not be loaded.' }),
          h('p', { class: 'mt-2', text: String(err.message || err) }),
          h('p', { class: 'mt-2' }, 'Build it with ', h('code', { text: 'cd site && go run . -serve :8000' }), '.')),
      );
      console.error(err);
    }
  },
};
