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
		write: 'zugferd({ profile: EN16931, strict: false }).create(saisie) hors mesure, puis toXML() ; strict: false car la validation XSD de la librairie appelle Java',
	},
	unsupported: {
		extract: 'pas de lecture de PDF',
		parse: 'pas de lecteur CII',
		validate: 'seule validation : XSD Factur-X en mode strict, par xsd-schema-validator qui exige Java (absent) ; pas de schematron',
		embed: 'pas d’injection d’un XML existant : embedInPdf attache le XML que la librairie génère depuis sa saisie',
		check: 'pas de lecture de PDF',
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
