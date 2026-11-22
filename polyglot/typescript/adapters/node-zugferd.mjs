// node-zugferd: generation per profile from its own input object, embedding
// into a PDF it is given. strict validation runs the Factur-X XSD through
// xsd-schema-validator, which needs Java: it is turned off here (strict: false).
import { zugferd } from 'node-zugferd';
import { EN16931 } from 'node-zugferd/profile/en16931';
import { nodeZugferdData, smallOnly } from '../lib/handmodels.mjs';

const invoicer = zugferd({ profile: EN16931, strict: false, logger: false });

export default {
	lib: 'node-zugferd',
	pkg: 'node-zugferd',
	notes: {
		write: 'zugferd({ profile: EN16931, strict: false }).create(input) outside the measurement, then toXML(); strict: false because the XSD validation of the library calls Java',
	},
	unsupported: {
		extract: 'no PDF reading',
		parse: 'no CII reader',
		validate: 'only validation: Factur-X XSD in strict mode, by xsd-schema-validator, which requires Java (not installed); no schematron',
		embed: 'no embedding of an existing XML: embedInPdf attaches the XML that the library generates from its input',
		check: 'no PDF reading',
	},
	write: {
		model(s, bench) {
			smallOnly(bench);
			return invoicer.create(nodeZugferdData());
		},
		write: (doc) => doc.toXML(),
		handBuilt: true,
	},
};
