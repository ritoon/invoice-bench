// @kontor-mcp/core: the library behind the Kontor MCP server. @xmldom/xmldom,
// xmllint-wasm for the XSD, official CEN 1.3.16 and KoSIT XRechnung
// Schematron compiled to SEF (in @kontor-mcp/rules) run by SaxonJS.
import { extractEmbeddedXml, isErrorLevel, parseInvoice, validateInvoice } from '@kontor-mcp/core';

export default {
	lib: '@kontor-mcp/core',
	pkg: '@kontor-mcp/core',
	notes: {
		extract: 'extractEmbeddedXml',
		parse: 'parseInvoice(octets)',
		validate:
			'validateInvoice(octets) (défauts) : XSD CII D16B (xmllint-wasm) + schematron CEN EN 16931 CII 1.3.16, plus KoSIT XRechnung 3.0.2 (schematron 2.5.0) si BT-24 le déclare, compilés en SEF et exécutés par SaxonJS ; erreurs error/fatal hors couche « plausibility » (règles maison KONTOR-PLAUS-*, que la librairie exclut elle-même de valid)',
	},
	unsupported: {
		write: 'pas de sérialisation CII publique : modelToCii n’est pas exporté (exports limités à « . »), seul convertInvoice enchaîne lecture, validation et écriture',
		embed: 'pas d’injection dans un PDF existant : generateZugferdPdf dessine son propre PDF à partir d’une saisie',
		check: 'pas de contrôle natif d’un PDF hybride : validateInvoice ne prend que du XML (detectInvoicePdf extrait, sans valider)',
	},
	async extract(s) {
		return (await extractEmbeddedXml(s.bytes)).xml;
	},
	parse(s) {
		return parseInvoice(s.bytes);
	},
	async validate(s) {
		const r = await validateInvoice(s.bytes);
		const errs = r.findings.filter((f) => f.source !== 'plausibility' && isErrorLevel(f));
		return { errors: errs.length, rules: errs.map((f) => f.ruleId) };
	},
};
