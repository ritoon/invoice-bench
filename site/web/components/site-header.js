/**
 * SiteHeader — wordmark, the two pages, link to the repository, theme switch.
 * Usage: SiteHeader.render(container, { active: 'perf' | 'compliance' })
 */
const SiteHeader = {
  pages: [
    { id: 'perf', href: 'index.html', text: 'CPU & memory' },
    { id: 'compliance', href: 'compliance.html', text: 'Compliance' },
  ],

  themeLabel(mode) {
    return { auto: 'Theme: auto', light: 'Theme: light', dark: 'Theme: dark' }[mode];
  },

  themeIcon(mode) {
    const s = (tag, a) => svg(tag, a);
    const icon = svg('svg', { width: 16, height: 16, viewBox: '0 0 16 16', 'aria-hidden': 'true', fill: 'none', stroke: 'currentColor', 'stroke-width': 1.5, 'stroke-linecap': 'round' });
    if (mode === 'light') {
      icon.append(s('circle', { cx: 8, cy: 8, r: 3 }));
      for (let i = 0; i < 8; i++) {
        const a = (i * Math.PI) / 4;
        icon.append(s('line', { x1: 8 + Math.cos(a) * 5.2, y1: 8 + Math.sin(a) * 5.2, x2: 8 + Math.cos(a) * 6.8, y2: 8 + Math.sin(a) * 6.8 }));
      }
    } else if (mode === 'dark') {
      icon.append(s('path', { d: 'M13 9.5A5.5 5.5 0 1 1 6.5 3a4.5 4.5 0 0 0 6.5 6.5Z' }));
    } else {
      icon.append(s('circle', { cx: 8, cy: 8, r: 5.5 }), svg('path', { d: 'M8 2.5v11A5.5 5.5 0 0 0 8 2.5Z', fill: 'currentColor', stroke: 'none' }));
    }
    return icon;
  },

  create({ active }) {
    const themeBtn = h('button', {
      type: 'button',
      class: 'inline-flex items-center gap-2 rounded-full px-3 py-1.5 text-sm text-ink-2 hover:bg-surface-2 hover:text-ink transition-colors',
    });
    const paint = (mode) => {
      themeBtn.replaceChildren(this.themeIcon(mode), h('span', { class: 'hidden sm:inline', text: this.themeLabel(mode) }));
      themeBtn.setAttribute('aria-label', `${this.themeLabel(mode)}. Switch theme`);
    };
    themeBtn.addEventListener('click', () => paint(Theme.next()));
    paint(Theme.get());

    return h('div', { class: 'mx-auto max-w-site px-4 sm:px-6' },
      h('div', { class: 'flex flex-wrap items-center gap-x-6 gap-y-3 py-4' },
        h('a', { href: 'index.html', class: 'flex items-baseline gap-2 mr-auto' },
          h('span', { class: 'text-base font-semibold tracking-tight text-ink', text: 'invoice-bench' }),
          h('span', { class: 'hidden md:inline text-sm text-muted', text: 'Factur-X / ZUGFeRD / EN 16931 libraries, nine languages' })),
        h('nav', { class: 'order-3 w-full sm:order-none sm:w-auto', 'aria-label': 'Pages' },
          h('ul', { class: 'flex gap-1' },
            this.pages.map((p) => h('li', {},
              h('a', { href: p.href, class: 'nav-tab inline-block', 'aria-current': p.id === active ? 'page' : null, text: p.text }))))),
        h('div', { class: 'flex items-center gap-1' },
          themeBtn,
          h('a', { href: REPO, class: 'inline-flex items-center gap-2 rounded-full px-3 py-1.5 text-sm text-ink-2 hover:bg-surface-2 hover:text-ink transition-colors' },
            this.githubIcon(), h('span', { class: 'hidden sm:inline', text: 'Repository' })))));
  },

  githubIcon() {
    return svg('svg', { width: 16, height: 16, viewBox: '0 0 16 16', fill: 'currentColor', 'aria-hidden': 'true' },
      svg('path', { d: 'M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.01 8.01 0 0 0 16 8c0-4.42-3.58-8-8-8Z' }));
  },

  render(container, props) {
    container.replaceChildren(this.create(props));
    return container;
  },
};
