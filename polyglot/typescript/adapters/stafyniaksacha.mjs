// @stafyniaksacha/facturx: port of akretion/factur-x (Python). pdf-lib for
// the PDF, libxml2-wasm for the XSD, SaxonJS for the Factur-X Schematron.
import { check, extract, generate, invoiceToXml, xmlToInvoice } from '@stafyniaksacha/facturx';

// The Schematron results carry no flag (warning / fatal): every failed assert
// counts. Code-list asserts (FACTUR-X_EN16931_codedb.xml) have no id either.
const ruleOf = (e) => e.id ?? `codedb ${/\*:(\w+)\[[^\]]*\]\[\d+\]$/.exec(e.location ?? '')?.[1] ?? ''}`.trim();

export default {
	lib: '@stafyniaksacha/facturx',
	pkg: '@stafyniaksacha/facturx',
	notes: {
		extract: 'extract({ pdf })',
		parse: 'xmlToInvoice',
		validate:
			'check({ xml, level: "en16931", schematron: true }): Factur-X 1.09 EN16931 XSD (libxml2-wasm) + official Factur-X 1.09 EN16931 schematron compiled to SEF, run by SaxonJS; the result carries no severity, every failed assertion counts; level forced to en16931 because automatic detection throws "Unknown level: kosit" on the XRECHNUNG files',
		write: 'xmlToInvoice outside the measurement, then invoiceToXml + toString (the libxml2-wasm XmlDocument is freed after each call)',
		embed: 'generate({ pdf, xml }) (default options)',
	},
	unsupported: {
		check: 'no native check of a hybrid PDF (check only takes XML)',
	},
	async extract(s) {
		return (await extract({ pdf: s.bytes })).xml;
	},
	parse(s) {
		return xmlToInvoice(s.text);
	},
	async validate(s) {
		const r = await check({ xml: s.text, level: 'en16931', schematron: true });
		const sch = r.schematronErrors ?? [];
		return {
			errors: r.errors.length + sch.length,
			rules: [...r.errors.map(() => 'XSD'), ...sch.map(ruleOf)],
		};
	},
	write: {
		model: (s) => xmlToInvoice(s.text),
		async write(m) {
			const doc = await invoiceToXml(m);
			try {
				return doc.toString();
			} finally {
				doc.dispose();
			}
		},
	},
	embed(host, xml) {
		return generate({ pdf: host.bytes, xml: xml.text });
	},
};
