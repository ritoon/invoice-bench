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
			'validate(bytes): EN 16931 rules + XRechnung CIUS (profile read from BT-24) + BR-CL codes, rewritten in TypeScript (no schematron, no XSD); fatal-severity errors',
		write: 'parseCiiInvoice outside the measurement, then generateCii(invoice)',
		check: 'validate(PDF bytes): extraction + rules on the XML + container findings (AW-PDF-*, warnings only); fatal-severity errors',
	},
	unsupported: {
		embed: 'no embedding into a PDF (the library reads PDFs, does not write them)',
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
