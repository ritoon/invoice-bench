// @fin.cx/einvoice: @xmldom/xmldom + xpath for the XML, pdf-lib for the PDF,
// EN 16931 / XRechnung / Factur-X rules written in TypeScript. Its Schematron
// layer needs files that the postinstall script downloads from GitHub (master
// branches, unpinned): that script is not run here (--ignore-scripts), and
// EInvoice.validate() does not use that layer anyway.
import { EInvoice, PDFEmbedder, PDFExtractor, ValidationLevel } from '@fin.cx/einvoice';

const extractor = new PDFExtractor();
const embedder = new PDFEmbedder();
// The optional validators behind feature flags (EN 16931 business rules on the
// decoded model, code lists) are added: the most complete set the library runs.
const VALIDATE_OPTIONS = { featureFlags: ['EN16931_BUSINESS_RULES', 'CODE_LIST_VALIDATION'] };

export default {
	lib: '@fin.cx/einvoice',
	pkg: '@fin.cx/einvoice',
	notes: {
		extract: 'new PDFExtractor().extractXml (une instance réutilisée ; ses console.log sont neutralisés)',
		parse: 'EInvoice.fromXml',
		validate:
			'EInvoice.fromXml puis validate(BUSINESS, { featureFlags: [EN16931_BUSINESS_RULES, CODE_LIST_VALIDATION] }) : règles EN 16931 / XRechnung / Factur-X réécrites en TypeScript (pas de XSD ; la couche schematron dépend de fichiers téléchargés au postinstall et n’est pas utilisée par validate)',
		write: 'EInvoice.fromXml hors mesure, puis exportXml("facturx")',
		embed: 'new PDFEmbedder().embedXml(pdf, xml, "factur-x.xml", …) : pdf-lib attach + save, sans XMP ni métadonnées PDF/A',
	},
	unsupported: {
		check: 'pas de contrôle natif d’un PDF hybride (EInvoice.fromPdf extrait et lit, la validation reste à appeler à part)',
	},
	async extract(s) {
		const r = await extractor.extractXml(s.bytes);
		if (!r.success) throw new Error(r.error?.message ?? 'no XML found');
		return r.xml;
	},
	parse(s) {
		return EInvoice.fromXml(s.text);
	},
	async validate(s) {
		const inv = await EInvoice.fromXml(s.text);
		const r = await inv.validate(ValidationLevel.BUSINESS, VALIDATE_OPTIONS);
		return { errors: r.errors.length, rules: r.errors.map((e) => e.code) };
	},
	write: {
		model: (s) => EInvoice.fromXml(s.text),
		write: (inv) => inv.exportXml('facturx'),
	},
	async embed(host, xml) {
		const r = await embedder.embedXml(host.bytes, xml.text, 'factur-x.xml', 'Factur-X XML invoice');
		if (!r.success) throw new Error(r.error?.message ?? 'embed failed');
		return r.data;
	},
};
