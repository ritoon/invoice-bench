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
			'check({ xml, level: "en16931", schematron: true }) : XSD Factur-X 1.09 EN16931 (libxml2-wasm) + schematron officiel Factur-X 1.09 EN16931 compilé en SEF, exécuté par SaxonJS ; le résultat ne porte pas de sévérité, toute assertion en échec compte ; niveau forcé à en16931 car la détection automatique lève « Unknown level: kosit » sur les XRECHNUNG',
		write: 'xmlToInvoice hors mesure, puis invoiceToXml + toString (le XmlDocument libxml2-wasm est libéré après chaque appel)',
		embed: 'generate({ pdf, xml }) (options par défaut)',
	},
	unsupported: {
		check: 'pas de contrôle natif d’un PDF hybride (check ne prend que du XML)',
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
