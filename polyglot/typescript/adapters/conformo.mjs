// Conformo (@conformo/toolkit and its packages): pdf-lib + pdfjs-dist for the
// PDF, own XML reader, official Schematron (CEN 1.3.17, KoSIT XRechnung 3.0.2,
// Peppol, BR-FR) compiled to SEF and run by SaxonJS, XSD by xmllint-wasm.
import { readFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { PDFDocument } from 'pdf-lib';
import { detectXmlFormat, extractEmbeddedXml } from '@conformo/parse';
import { validate as validateXml } from '@conformo/validate';
import { finalizePDFA } from '@conformo/pdf';
import { serialize, validate as validateAny } from '@conformo/toolkit';

const isError = (r) => r.severity === 'error' || r.severity === 'fatal';
const count = (r) => {
	const errs = r.results.filter(isError);
	const xsd = r.xsd?.errors ?? [];
	return { errors: errs.length + xsd.length, rules: [...xsd.map(() => 'XSD'), ...errs.map((e) => e.ruleId)] };
};

// The sRGB profile createPdf passes to finalizePDFA, read once.
const toolkitDir = join(dirname(fileURLToPath(import.meta.url)), '..', 'node_modules', '@conformo', 'toolkit');
let icc;
const iccProfile = () => (icc ??= readFileSync(join(toolkitDir, 'assets', 'sRGB.icc')));
const createDate = new Date('2026-11-16T12:00:00Z');

export default {
	lib: '@conformo/toolkit',
	pkg: '@conformo/toolkit',
	notes: {
		extract: '@conformo/parse extractEmbeddedXml',
		parse: '@conformo/parse detectXmlFormat (CII/UBL parsing into the EN 16931 model of @conformo/core)',
		validate:
			'@conformo/validate validate(xml) (defaults): CII D16B XSD (xmllint-wasm) + CEN EN 16931 CII 1.3.17 schematron, and KoSIT XRechnung 3.0.2 if BT-24 declares it, compiled to SEF and run by SaxonJS; error/fatal severities after the level overrides of the KoSIT configuration; BR-FR not enabled (country not supplied)',
		write: 'detectXmlFormat outside the measurement, then @conformo/toolkit serialize(invoice, "cii")',
		embed:
			'createPdf recipe applied to the host: PDFDocument.load (pdf-lib) + @conformo/pdf finalizePDFA (en16931 profile, sRGB ICC of the toolkit, fixed date) + save({ useObjectStreams: false })',
		check: '@conformo/toolkit validate(PDF bytes): extraction + XSD + schematrons + container notices (XMP, file name; they do not affect valid)',
	},
	extract(s) {
		return extractEmbeddedXml(s.bytes).then((r) => r.xml);
	},
	parse(s) {
		return detectXmlFormat(s.text).invoice;
	},
	async validate(s) {
		return count(await validateXml(s.text));
	},
	write: {
		model: (s) => detectXmlFormat(s.text).invoice,
		write: (m) => serialize(m, 'cii'),
	},
	async embed(host, xml) {
		const doc = await PDFDocument.load(host.bytes);
		await finalizePDFA(doc, {
			xml: xml.text,
			profile: 'en16931',
			title: 'Invoice 471102',
			author: 'Lieferant GmbH',
			producer: 'Conformo',
			creatorTool: 'Conformo',
			iccProfile: iccProfile(),
			createDate,
		});
		return doc.save({ useObjectStreams: false });
	},
	async check(s) {
		return count(await validateAny(s.bytes));
	},
};
