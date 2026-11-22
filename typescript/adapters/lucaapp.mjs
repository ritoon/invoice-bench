// @lucaapp/zugferd-ts: pure TypeScript, xmllint-wasm for the XSD,
// node-schematron (fontoxpath) for the CEN EN 16931 Schematron 1.3.16.
import { detectProfile, embedXml, extractXml, parseCiiXml, serializeToXml, validateUblXml, validateXml } from '@lucaapp/zugferd-ts';
import { isUbl } from '../lib/corpus.mjs';

export default {
	lib: '@lucaapp/zugferd-ts',
	pkg: '@lucaapp/zugferd-ts',
	notes: {
		extract: 'extractXml',
		parse: 'parseCiiXml',
		validate:
			'validateXml (defaults): Factur-X 1.08 XSD of the detected profile (xmllint-wasm) + CEN EN 16931 CII 1.3.16 schematron (node-schematron, pure JS), error-severity issues; the library classifies every failed assertion as error (the CEN flag="warning" attributes are not read); for a UBL document, validateUblXml: UBL 2.1 XSD only (the library has no UBL schematron)',
		write: 'parseCiiXml then detectProfile outside the measurement, serializeToXml(model, profile)',
		embed: 'embedXml(pdf, xml) (defaults: factur-x.xml, PDF/A-3b layer)',
	},
	unsupported: {
		check: 'no native complete check: validatePdfA checks only the container (21 rules, PDFA_RULES), validateXml only the XML',
	},
	extract(s) {
		return extractXml(s.bytes);
	},
	parse(s) {
		return parseCiiXml(s.text);
	},
	async validate(s) {
		if (isUbl(s.text)) {
			const u = await validateUblXml(s.text);
			return { errors: u.errors.length, rules: u.errors.map((e) => `XSD: ${e.message}`) };
		}
		const r = await validateXml(s.text);
		const errs = r.issues.filter((i) => i.severity === 'error');
		return { errors: errs.length, rules: errs.map((i) => i.ruleId ?? i.source) };
	},
	write: {
		model(s) {
			return { inv: parseCiiXml(s.text), profile: detectProfile(s.text) };
		},
		write(m) {
			return serializeToXml(m.inv, m.profile);
		},
	},
	embed(host, xml) {
		return embedXml(host.bytes, xml.text);
	},
};
