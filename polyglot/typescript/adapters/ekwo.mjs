// @ekwo-ai/factur-x: CII generation from a plain invoice object, a CII reader
// with eight arithmetic checks (BR-CO-10..17), PDF/A-3 embedding and
// extraction with pdf-lib.
import { generateCiiXml, readCii } from '@ekwo-ai/factur-x';
import { embedFacturX, extractFacturX } from '@ekwo-ai/factur-x/pdf';
import { ekwoInvoice, smallOnly } from '../lib/handmodels.mjs';

export default {
	lib: '@ekwo-ai/factur-x',
	pkg: '@ekwo-ai/factur-x',
	notes: {
		extract: '@ekwo-ai/factur-x/pdf extractFacturX',
		parse: 'readCii(xml)',
		validate:
			'readCii(xml).violations: parsing anomalies (dates, IBAN…) + 8 arithmetic checks BR-CO-10..17 (checkReceived); no XSD or schematron',
		write: 'generateCiiXml(invoice, { profile: "en16931" }); the model read by readCii is not the one writing takes, hence a hand-built Invoice model (one note only, no item GTIN)',
		embed: '@ekwo-ai/factur-x/pdf embedFacturX(pdf, xml, { profile: "en16931" }): pdf-lib, attachment + XMP',
	},
	unsupported: {
		check: 'no native check of a hybrid PDF',
	},
	async extract(s) {
		const r = await extractFacturX(s.bytes);
		if (!r) throw new Error('no embedded invoice');
		return r.xml;
	},
	parse(s) {
		return readCii(s.text);
	},
	validate(s) {
		const r = readCii(s.text);
		return { errors: r.violations.length, rules: r.violations.map((v) => v.code) };
	},
	write: {
		model(s, bench) {
			smallOnly(bench, 'the readCii model cannot be written');
			return ekwoInvoice();
		},
		write: (inv) => generateCiiXml(inv, { profile: 'en16931' }),
		handBuilt: true,
	},
	embed(host, xml) {
		return embedFacturX(host.bytes, xml.text, { profile: 'en16931' });
	},
};
