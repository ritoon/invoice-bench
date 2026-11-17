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
			'validateInvoice(octets) (défauts) : XSD CII D16B (validateur XSD maison) + schematron CEN EN 16931 CII 1.3.16, et KoSIT XRechnung 2.6.0 si BT-24 le déclare, avec les sévérités de la configuration KoSIT ; moteur XPath 2.0 maison, sans Java',
		check: 'validateInvoice(octets du PDF) : extraction, comparaison du profil XMP, XSD et schematrons sur le XML',
	},
	unsupported: {
		parse: 'pas de modèle de facture (lecture réservée à « Summand Pro »)',
		write: 'pas d’écriture (réservée à « Summand Pro »)',
		embed: 'pas d’injection dans un PDF',
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
