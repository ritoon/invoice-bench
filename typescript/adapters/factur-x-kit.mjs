// factur-x-kit: CII reader and writer for MINIMUM … EN16931 (EXTENDED falls
// back to the EN16931 mapping; no XRechnung), business rules written as zod
// refinements, PDF generation with pdf-lib. XML from Buffer (fromXML).
import { FacturX } from 'factur-x-kit';

const ruleOf = (e) => /\bBR-[A-Z0-9-]*\d+\b/.exec(e.message)?.[0] ?? e.message.slice(0, 40);

export default {
	lib: 'factur-x-kit',
	pkg: 'factur-x-kit',
	notes: {
		parse: 'FacturX.fromXML(bytes)',
		validate:
			'FacturX.fromXML(bytes) then validate(): zod schema of the detected profile and BR-* rules rewritten in TypeScript (MINIMUM to EN16931); no XSD or schematron',
		write: 'FacturX.fromXML outside the measurement, then getXML()',
	},
	unsupported: {
		extract:
			'no extraction alone: FacturXPdf is not exported and FacturX.fromPDF chains extraction and parsing of the XML (factur-x.xml only)',
		embed: 'no embedding of an existing XML: getPDF({ existingNonConformantPdf }) rewrites the XML from the model',
		check: 'no native check of a hybrid PDF',
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
