/**
 * SiteFooter — where the figures come from: run dates, machine, toolchains.
 * Usage: SiteFooter.render(container, { run })
 */
const SiteFooter = {
  create({ run }) {
    const env = new Map(run.env);
    const keep = ['cpu', 'os', 'go', 'node', 'python', 'rust', 'java', 'dotnet', 'cc', 'php'];
    return h('div', { class: 'mx-auto max-w-site px-4 sm:px-6 py-10 text-xs text-muted' },
      h('div', { class: 'grid gap-6 md:grid-cols-[1fr_2fr]' },
        h('div', { class: 'space-y-2 prose-site' },
          h('p', {}, 'Raw records: ',
            fileLink(`results/${run.date}`, `results/${run.date}`), ' and ',
            fileLink(`go/results/${run.goDate}`, `go/results/${run.goDate}`), '.'),
          h('p', {}, 'Protocol: ', fileLink('PROTOCOL.md', 'PROTOCOL.md'), '. Built by ', fileLink('site/main.go', 'site/main.go'),
            ' from the committed records', run.generated ? ` on ${run.generated.slice(0, 10)}` : '', '.'),
          h('p', {}, 'Code under the MIT license. The libraries keep their own licenses.')),
        h('dl', { class: 'grid grid-cols-[auto_1fr] gap-x-4 gap-y-1' },
          keep.filter((k) => env.has(k)).map((k) => [
            h('dt', { class: 'font-medium text-ink-2', text: k }),
            h('dd', { class: 'num break-words', text: env.get(k) }),
          ]))));
  },

  render(container, props) {
    container.replaceChildren(this.create(props));
    return container;
  },
};
