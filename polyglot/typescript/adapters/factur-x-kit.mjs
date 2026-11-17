// factur-x-kit: CII reader and writer for MINIMUM … EN16931 (EXTENDED falls
// back to the EN16931 mapping; no XRechnung), business rules written as zod
// refinements, PDF generation with pdf-lib. XML from Buffer (fromXML).
import { FacturX } from 'factur-x-kit';

const ruleOf = (e) => /\bBR-[A-Z0-9-]*\d+\b/.exec(e.message)?.[0] ?? e.message.slice(0, 40);

export default {
	lib: 'factur-x-kit',
	pkg: 'factur-x-kit',
	notes: {
		parse: 'FacturX.fromXML(octets)',
		validate:
			'FacturX.fromXML(octets) puis validate() : schéma zod du profil détecté et règles BR-* réécrites en TypeScript (MINIMUM à EN16931) ; pas de XSD ni de schematron',
		write: 'FacturX.fromXML hors mesure, puis getXML()',
	},
	unsupported: {
		extract:
			'pas d’extraction seule : FacturXPdf n’est pas exporté et FacturX.fromPDF enchaîne extraction et lecture du XML (factur-x.xml uniquement)',
		embed: 'pas d’injection d’un XML existant : getPDF({ existingNonConformantPdf }) réécrit le XML depuis le modèle',
		check: 'pas de contrôle natif d’un PDF hybride',
	},
	parse(s) {
		return FacturX.fromXML(s.bytes);
	},
	async validate(s) {
		const fx = await FacturX.fromXML(s.bytes);
		const r = fx.validate();
		const errs = r.valid ? [] : (r.errors ?? [{ message: 'invalid' }]);
		return { errors: errs.length, rules: errs.map(ruleOf) };
	},
	write: {
		model: (s) => FacturX.fromXML(s.bytes),
		write: (fx) => fx.getXML(),
	},
};
