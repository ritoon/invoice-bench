// Shared inputs, identical to bench/ecosystem/common/common.go: same files,
// same lists, same sanity check.
import { readFileSync, readdirSync, statSync } from 'node:fs';
import { basename, dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const here = dirname(fileURLToPath(import.meta.url));
export const ECOSYSTEM = resolve(here, '..', '..', '..', 'ecosystem');
export const CORPUS = process.env.FX_CORPUS || join(ECOSYSTEM, 'corpus');

export function load(...rel) {
	const path = join(CORPUS, ...rel);
	const bytes = readFileSync(path);
	return { name: basename(path), path, bytes, text: bytes.toString('utf8') };
}

function loadHost(name) {
	const path = join(ECOSYSTEM, name);
	const bytes = readFileSync(path);
	return { name, path, bytes, text: '' };
}

// PDF inputs (Factur-X / ZUGFeRD 2.x hybrids).
export const pdfSmall = () => load('XML-Rechnung', 'FX', 'EN16931_Einfach.pdf');
export const pdfLarge = () => load('XML-Rechnung', 'FX', 'EN16931_Elektron.pdf');
export const pdfMinimum = () => load('ZUGFeRDv2', 'correct', 'FNFE-factur-x-examples', 'Facture_FR_MINIMUM.pdf');
export const pdfExtended = () => load('ZUGFeRDv2', 'correct', 'intarsys', 'EXTENDED', 'zugferd_2p0_EXTENDED_Warenrechnung.pdf');

// XML inputs (CII D16B).
export const xmlSmall = () => load('XML-Rechnung', 'CII', 'EN16931_Einfach.cii.xml');
export const xmlLines = () => load('XML-Rechnung', 'CII', 'EN16931_Rabatte.cii.xml');
export const xmlLarge = () => load('XML-Rechnung', 'CII', 'XRECHNUNG_Elektron.cii.xml');

// Host PDFs for embedding (built by bench/ecosystem/run.sh setup).
export const hostPDF = () => loadHost('host.pdf');
export const hostPDFClassic = () => loadHost('host_classic.pdf');

// LooksLikeCII is the sanity check applied to an extracted XML.
export function looksLikeCII(x) {
	if (x == null) return false;
	const s = typeof x === 'string' ? x : Buffer.from(x.buffer ?? x, x.byteOffset ?? 0, x.byteLength ?? x.length).toString('latin1');
	return s.includes('CrossIndustryInvoice') || s.includes('CrossIndustryDocument');
}

// The 26 valid CII samples (false-positive survey and parse survey).
export function ciiFiles() {
	const dir = join(CORPUS, 'XML-Rechnung', 'CII');
	return readdirSync(dir)
		.filter((f) => f.endsWith('.cii.xml') && f !== 'not_validating_full_invoice_based_onTest_EeISI_300_CENfullmodel.cii.xml')
		.sort()
		.map((f) => join(dir, f));
}

// The 104 valid hybrid PDFs: ZUGFeRDv2/correct (recursive) and XML-Rechnung/FX.
export function pdfFiles() {
	const out = [];
	const walk = (d) => {
		for (const e of readdirSync(d).sort()) {
			const p = join(d, e);
			if (statSync(p).isDirectory()) walk(p);
			else if (e.endsWith('.pdf')) out.push(p);
		}
	};
	walk(join(CORPUS, 'ZUGFeRDv2', 'correct'));
	const fx = join(CORPUS, 'XML-Rechnung', 'FX');
	for (const e of readdirSync(fx).sort()) if (e.endsWith('.pdf')) out.push(join(fx, e));
	return out;
}

export function readSample(path) {
	const bytes = readFileSync(path);
	return { name: basename(path), path, bytes, text: bytes.toString('utf8') };
}
