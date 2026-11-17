// @lucaapp/zugferd-ts: pure TypeScript, xmllint-wasm for the XSD,
// node-schematron (fontoxpath) for the CEN EN 16931 Schematron 1.3.16.
import { detectProfile, embedXml, extractXml, parseCiiXml, serializeToXml, validateXml } from '@lucaapp/zugferd-ts';

export default {
	lib: '@lucaapp/zugferd-ts',
	pkg: '@lucaapp/zugferd-ts',
	notes: {
		extract: 'extractXml',
		parse: 'parseCiiXml',
		validate:
			'validateXml (défauts) : XSD Factur-X 1.08 du profil détecté (xmllint-wasm) + schematron CEN EN 16931 CII 1.3.16 (node-schematron, JS pur), issues de sévérité error ; la librairie classe toute assertion en échec en error (les flag="warning" du CEN ne sont pas lus)',
		write: 'parseCiiXml puis detectProfile hors mesure, serializeToXml(modèle, profil)',
		embed: 'embedXml(pdf, xml) (défauts : factur-x.xml, surcouche PDF/A-3b)',
	},
	unsupported: {
		check: 'pas de contrôle complet natif : validatePdfA ne contrôle que le conteneur (21 règles, PDFA_RULES), validateXml que le XML',
	},
	extract(s) {
		return extractXml(s.bytes);
	},
	parse(s) {
		return parseCiiXml(s.text);
	},
	async validate(s) {
		const r = await validateXml(s.text);
		const errs = r.issues.filter((i) => i.severity === 'error');
		return { errors: errs.length, rules: errs.map((i) => i.ruleId ?? i.source) };
	},
	write: {
		model(s) {
			return { inv: parseCiiXml(s.text), profile: detectProfile(s.text) };
		},
		write(m) {
			return serializeToXml(m.inv, m.profile);
		},
	},
	embed(host, xml) {
		return embedXml(host.bytes, xml.text);
	},
};
