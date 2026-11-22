// @kontor-mcp/core: the library behind the Kontor MCP server. @xmldom/xmldom,
// xmllint-wasm for the XSD, official CEN 1.3.16 and KoSIT XRechnung
// Schematron compiled to SEF (in @kontor-mcp/rules) run by SaxonJS.
import { extractEmbeddedXml, isErrorLevel, parseInvoice, validateInvoice } from '@kontor-mcp/core';

export default {
	lib: '@kontor-mcp/core',
	pkg: '@kontor-mcp/core',
	notes: {
		extract: 'extractEmbeddedXml',
		parse: 'parseInvoice(bytes)',
		validate:
			'validateInvoice(bytes) (defaults): CII D16B XSD (xmllint-wasm) + CEN EN 16931 CII 1.3.16 schematron, plus KoSIT XRechnung 3.0.2 (schematron 2.5.0) if BT-24 declares it, compiled to SEF and run by SaxonJS; error/fatal errors outside the "plausibility" layer (in-house KONTOR-PLAUS-* rules, which the library itself excludes from valid)',
	},
	unsupported: {
		write: 'no public CII serialization: modelToCii is not exported (exports limited to "."), only convertInvoice chains parsing, validation and writing',
		embed: 'no embedding into an existing PDF: generateZugferdPdf draws its own PDF from an input structure',
		check: 'no native check of a hybrid PDF: validateInvoice only takes XML (detectInvoicePdf extracts, without validating)',
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
