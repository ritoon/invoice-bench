// node survey.mjs <lib>  →  lines of survey_typescript.jsonl (validate, parse, extract)
//
// Same surveys as bench/ecosystem/common/harness.go: valid samples a validator
// flags (false positives, out of 26), CII samples a parser reads (out of 26),
// hybrid PDFs an invoice XML is extracted from (out of 104).
import { pathToFileURL } from 'node:url';
import * as corpus from './lib/corpus.mjs';
import { emit, errMessage, LANG, pkgVersion, silenceConsole } from './lib/harness.mjs';
import { loadAdapter } from './adapters/index.mjs';

function line(adapter, op, ok, total, failures, note) {
	return { lang: LANG, lib: adapter.lib, pkg: adapter.pkg, version: pkgVersion(adapter.pkg), op, ok, total, failures, note: note ?? '' };
}

export async function surveyValidate(adapter) {
	let flagged = 0, total = 0;
	const failures = [];
	for (const f of corpus.ciiFiles()) {
		const s = corpus.readSample(f);
		total++;
		try {
			const r = await adapter.validate(s);
			if (r.errors > 0) {
				flagged++;
				failures.push({ file: s.name, errors: r.errors, rules: [...new Set(r.rules ?? [])] });
			}
		} catch (e) {
			flagged++;
			failures.push({ file: s.name, error: errMessage(e) });
		}
	}
	return line(adapter, 'validate', flagged, total, failures, adapter.notes?.validate);
}

export async function surveyParse(adapter) {
	let ok = 0, total = 0;
	const failures = [];
	for (const f of corpus.ciiFiles()) {
		const s = corpus.readSample(f);
		total++;
		try {
			const m = await adapter.parse(s);
			if (m == null) throw new Error('nothing returned');
			ok++;
		} catch (e) {
			failures.push({ file: s.name, error: errMessage(e) });
		}
	}
	return line(adapter, 'parse', ok, total, failures, adapter.notes?.parse);
}

export async function surveyExtract(adapter) {
	let ok = 0, total = 0;
	const failures = [];
	for (const f of corpus.pdfFiles()) {
		const s = corpus.readSample(f);
		total++;
		try {
			const x = await adapter.extract(s);
			if (!corpus.looksLikeCII(x)) throw new Error('no CrossIndustryInvoice in the extracted XML');
			ok++;
		} catch (e) {
			failures.push({ file: s.path.slice(corpus.CORPUS.length + 1), error: errMessage(e) });
		}
	}
	return line(adapter, 'extract', ok, total, failures, adapter.notes?.extract);
}

if (import.meta.url === pathToFileURL(process.argv[1]).href) {
	const [lib] = process.argv.slice(2);
	silenceConsole();
	// A stray rejected promise from a library must not end the survey: it is
	// reported on stderr and the files go on (each call has its own try/catch).
	process.on('unhandledRejection', (e) => process.stderr.write(`${lib}: unhandled rejection: ${errMessage(e)}\n`));
	const adapter = await loadAdapter(lib);
	if (typeof adapter.validate === 'function') emit(await surveyValidate(adapter));
	if (typeof adapter.parse === 'function') emit(await surveyParse(adapter));
	if (typeof adapter.extract === 'function') emit(await surveyExtract(adapter));
	process.exit(0);
}
