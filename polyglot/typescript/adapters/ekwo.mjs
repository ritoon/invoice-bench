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
			'readCii(xml).violations : anomalies de lecture (dates, IBAN…) + 8 contrôles arithmétiques BR-CO-10..17 (checkReceived) ; pas de XSD ni de schematron',
		write: 'generateCiiXml(invoice, { profile: "en16931" }) ; le modèle lu par readCii n’est pas celui de l’écriture, d’où un modèle Invoice à la main (une seule note, pas de GTIN d’article)',
		embed: '@ekwo-ai/factur-x/pdf embedFacturX(pdf, xml, { profile: "en16931" }) : pdf-lib, pièce jointe + XMP',
	},
	unsupported: {
		check: 'pas de contrôle natif d’un PDF hybride',
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
			smallOnly(bench, 'le modèle de readCii ne s’écrit pas');
			return ekwoInvoice();
		},
		write: (inv) => generateCiiXml(inv, { profile: 'en16931' }),
		handBuilt: true,
	},
	embed(host, xml) {
		return embedFacturX(host.bytes, xml.text, { profile: 'en16931' });
	},
};
