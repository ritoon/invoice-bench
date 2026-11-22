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
			'validateXsd(xml, Profile.EN16931): Factur-X EN16931 XSD alone (libxml2-wasm; schema re-read and compiled on every call by the library); no EN 16931 schematron',
		write: 'buildXml(FacturXInvoiceInput input, Profile.EN16931), without a CII reader',
		embed: 'embedFacturX({ pdf, xml, profile: EN16931 }): pdf-lib, AFRelationship attachment + PDF/A-3b XMP',
	},
	unsupported: {
		parse: 'no CII reader (extractXml returns the text, without a model)',
		check: 'no native check of a hybrid PDF',
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
