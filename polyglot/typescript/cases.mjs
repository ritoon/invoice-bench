// node cases.mjs <lib>  →  lines of cases_typescript.jsonl
//
// Validation of the conformance cases of ../../cases (cases.tsv, or CASES_DIR):
// one line per suite, {suite, results: [{case, errors, rules[, error]}], note},
// as the other harnesses write them. No line for a library without validate.
import { readFileSync } from 'node:fs';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';
import * as corpus from './lib/corpus.mjs';
import { emit, errMessage, LANG, pkgVersion, silenceConsole } from './lib/harness.mjs';
import { loadAdapter } from './adapters/index.mjs';

const here = dirname(fileURLToPath(import.meta.url));
export const CASES = resolve(process.env.CASES_DIR || join(here, '..', '..', 'cases'));

// The documents of cases.tsv by suite, in its order (columns suite and file).
export function caseSuites() {
	const out = new Map();
	for (const line of readFileSync(join(CASES, 'cases.tsv'), 'utf8').split('\n').slice(1)) {
		const [suite, file] = line.split('\t');
		if (!suite?.trim() || !file) continue;
		if (!out.has(suite)) out.set(suite, []);
		out.get(suite).push(file);
	}
	return out;
}

export async function casesValidate(adapter) {
	const lines = [];
	for (const [suite, files] of caseSuites()) {
		let rejected = 0;
		const results = [];
		for (const file of files) {
			const s = corpus.readSample(join(CASES, file));
			try {
				const r = await adapter.validate(s);
				if (r.errors > 0) rejected++;
				results.push({ case: file, errors: r.errors, rules: [...new Set(r.rules ?? [])].sort() });
			} catch (e) {
				rejected++;
				results.push({ case: file, errors: 0, rules: [], error: errMessage(e) });
			}
		}
		process.stderr.write(`   CASES ${suite}: ${rejected}/${files.length} rejected\n`);
		lines.push({ lang: LANG, lib: adapter.lib, pkg: adapter.pkg, version: pkgVersion(adapter.pkg), suite, results, note: adapter.notes?.validate ?? '' });
	}
	return lines;
}

if (import.meta.url === pathToFileURL(process.argv[1]).href) {
	const [lib] = process.argv.slice(2);
	silenceConsole();
	// A stray rejected promise from a library must not end the run: it is
	// reported on stderr and the cases go on (each call has its own try/catch).
	process.on('unhandledRejection', (e) => process.stderr.write(`${lib}: unhandled rejection: ${errMessage(e)}\n`));
	const adapter = await loadAdapter(lib);
	if (typeof adapter.validate === 'function') for (const l of await casesValidate(adapter)) emit(l);
	process.exit(0);
}
