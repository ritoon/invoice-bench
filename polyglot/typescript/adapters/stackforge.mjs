// @stackforge-eu/factur-x: generation from its own input object (buildXml),
// PDF/A-3 embedding with pdf-lib, extraction, XSD by libxml2-wasm. No CII
// reader, no Schematron.
import { Profile, buildXml, embedFacturX, extractXml, validateXsd } from '@stackforge-eu/factur-x';
import { smallOnly, stackforgeInput } from '../lib/handmodels.mjs';

export default {
	lib: '@stackforge-eu/factur-x',
	pkg: '@stackforge-eu/factur-x',
	notes: {
		extract: 'extractXml',
		validate:
			'validateXsd(xml, Profile.EN16931) : XSD Factur-X EN16931 seul (libxml2-wasm ; schéma relu et compilé à chaque appel par la librairie) ; pas de schematron EN 16931',
		write: 'buildXml(saisie FacturXInvoiceInput, Profile.EN16931), sans lecteur CII',
		embed: 'embedFacturX({ pdf, xml, profile: EN16931 }) : pdf-lib, pièce jointe AFRelationship + XMP PDF/A-3b',
	},
	unsupported: {
		parse: 'pas de lecteur CII (extractXml rend le texte, sans modèle)',
		check: 'pas de contrôle natif d’un PDF hybride',
	},
	async extract(s) {
		return (await extractXml(s.bytes)).xml;
	},
	async validate(s) {
		const r = await validateXsd(s.text, Profile.EN16931);
		return { errors: r.errors.length, rules: r.errors.map(() => 'XSD') };
	},
	write: {
		model(s, bench) {
			smallOnly(bench);
			return stackforgeInput();
		},
		write: (input) => buildXml(input, Profile.EN16931),
		handBuilt: true,
	},
	async embed(host, xml) {
		return (await embedFacturX({ pdf: host.bytes, xml: xml.text, profile: Profile.EN16931 })).pdf;
	},
};
