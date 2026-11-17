// Conformo (@conformo/toolkit and its packages): pdf-lib + pdfjs-dist for the
// PDF, own XML reader, official Schematron (CEN 1.3.17, KoSIT XRechnung 3.0.2,
// Peppol, BR-FR) compiled to SEF and run by SaxonJS, XSD by xmllint-wasm.
import { readFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import { PDFDocument } from 'pdf-lib';
import { detectXmlFormat, extractEmbeddedXml } from '@conformo/parse';
import { validate as validateXml } from '@conformo/validate';
import { finalizePDFA } from '@conformo/pdf';
import { serialize, validate as validateAny } from '@conformo/toolkit';

const isError = (r) => r.severity === 'error' || r.severity === 'fatal';
const count = (r) => {
	const errs = r.results.filter(isError);
	const xsd = r.xsd?.errors ?? [];
	return { errors: errs.length + xsd.length, rules: [...xsd.map(() => 'XSD'), ...errs.map((e) => e.ruleId)] };
};

// The sRGB profile createPdf passes to finalizePDFA, read once.
const toolkitDir = join(dirname(fileURLToPath(import.meta.url)), '..', 'node_modules', '@conformo', 'toolkit');
let icc;
const iccProfile = () => (icc ??= readFileSync(join(toolkitDir, 'assets', 'sRGB.icc')));
const createDate = new Date('2026-11-16T12:00:00Z');

export default {
	lib: '@conformo/toolkit',
	pkg: '@conformo/toolkit',
	notes: {
		extract: '@conformo/parse extractEmbeddedXml',
		parse: '@conformo/parse detectXmlFormat (lecture CII/UBL vers le modèle EN 16931 de @conformo/core)',
		validate:
			'@conformo/validate validate(xml) (défauts) : XSD CII D16B (xmllint-wasm) + schematron CEN EN 16931 CII 1.3.17, et KoSIT XRechnung 3.0.2 si BT-24 le déclare, compilés en SEF et exécutés par SaxonJS ; sévérités error/fatal après les re-niveaux de la configuration KoSIT ; BR-FR non activé (country non fourni)',
		write: 'detectXmlFormat hors mesure, puis @conformo/toolkit serialize(invoice, "cii")',
		embed:
			'recette de createPdf appliquée à l’hôte : PDFDocument.load (pdf-lib) + @conformo/pdf finalizePDFA (profil en16931, ICC sRGB du toolkit, date fixe) + save({ useObjectStreams: false })',
		check: '@conformo/toolkit validate(octets du PDF) : extraction + XSD + schematrons + notices du conteneur (XMP, nom de fichier ; n’affectent pas valid)',
	},
	extract(s) {
		return extractEmbeddedXml(s.bytes).then((r) => r.xml);
	},
	parse(s) {
		return detectXmlFormat(s.text).invoice;
	},
	async validate(s) {
		return count(await validateXml(s.text));
	},
	write: {
		model: (s) => detectXmlFormat(s.text).invoice,
		write: (m) => serialize(m, 'cii'),
	},
	async embed(host, xml) {
		const doc = await PDFDocument.load(host.bytes);
		await finalizePDFA(doc, {
			xml: xml.text,
			profile: 'en16931',
			title: 'Invoice 471102',
			author: 'Lieferant GmbH',
			producer: 'Conformo',
			creatorTool: 'Conformo',
			iccProfile: iccProfile(),
			createDate,
		});
		return doc.save({ useObjectStreams: false });
	},
	async check(s) {
		return count(await validateAny(s.bytes));
	},
};
