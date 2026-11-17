// @attestwire/en16931: zero dependencies, own XML parser, own PDF reader and
// DEFLATE decoder, EN 16931 / XRechnung / Peppol rules written by hand in TS.
import { extractFacturX, generateCii, parseCiiInvoice, validate } from '@attestwire/en16931';

// result.errors holds the fatal findings; warnings and information are apart.
const count = (r) => ({ errors: r.errors.length, rules: r.errors.map((e) => e.rule) });

export default {
	lib: '@attestwire/en16931',
	pkg: '@attestwire/en16931',
	notes: {
		extract: 'extractFacturX',
		parse: 'parseCiiInvoice',
		validate:
			'validate(octets) : règles EN 16931 + CIUS XRechnung (profil lu dans BT-24) + codes BR-CL, réécrites en TypeScript (pas de schematron, pas de XSD) ; erreurs de sévérité fatal',
		write: 'parseCiiInvoice hors mesure, puis generateCii(invoice)',
		check: 'validate(octets du PDF) : extraction + règles sur le XML + constats du conteneur (AW-PDF-*, avertissements seulement) ; erreurs de sévérité fatal',
	},
	unsupported: {
		embed: 'pas d’injection dans un PDF (la librairie lit les PDF, ne les écrit pas)',
	},
	extract(s) {
		return extractFacturX(s.bytes).xml;
	},
	parse(s) {
		return parseCiiInvoice(s.text);
	},
	validate(s) {
		return count(validate(s.bytes));
	},
	write: {
		model: (s) => parseCiiInvoice(s.text).invoice,
		write: (m) => generateCii(m),
	},
	check(s) {
		return count(validate(s.bytes));
	},
};
