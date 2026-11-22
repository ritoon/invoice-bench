// @sweberdev/summand: validator only, zero dependencies, own XML Schema
// validator and own XPath 2.0 / Schematron engine running the official CEN
// and KoSIT rule sets.
import { extractInvoiceXml, validateInvoice } from '@sweberdev/summand';

const count = (r) => ({ errors: r.errors.length, rules: r.errors.map((e) => e.id) });

export default {
	lib: '@sweberdev/summand',
	pkg: '@sweberdev/summand',
	notes: {
		extract: 'extractInvoiceXml',
		validate:
			'validateInvoice(bytes) (defaults): CII D16B XSD (in-house XSD validator) + CEN EN 16931 CII 1.3.16 schematron, and KoSIT XRechnung 2.6.0 if BT-24 declares it, with the severities of the KoSIT configuration; in-house XPath 2.0 engine, no Java',
		check: 'validateInvoice(PDF bytes): extraction, XMP profile comparison, XSD and schematrons on the XML',
	},
	unsupported: {
		parse: 'no invoice model (parsing reserved for "Summand Pro")',
		write: 'no writing (reserved for "Summand Pro")',
		embed: 'no embedding into a PDF',
	},
	extract(s) {
		const r = extractInvoiceXml(s.bytes);
		if (!r) throw new Error('no embedded invoice');
		return r.xml;
	},
	validate(s) {
		return count(validateInvoice(s.bytes));
	},
	check(s) {
		return count(validateInvoice(s.bytes));
	},
};
