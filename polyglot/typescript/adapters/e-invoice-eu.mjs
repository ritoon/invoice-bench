// @e-invoice-eu/core: generation only, from its UBL-shaped JSON model (checked
// by an ajv JSON schema) to UBL, CII, XRechnung or Factur-X. No reader, no XML
// validation; Factur-X output needs a PDF input or LibreOffice.
import { InvoiceService } from '@e-invoice-eu/core';
import { eInvoiceEuInvoice, smallOnly } from '../lib/handmodels.mjs';

const logger = { log() {}, warn() {}, error() {} };
const service = new InvoiceService(logger);

export default {
	lib: '@e-invoice-eu/core',
	pkg: '@e-invoice-eu/core',
	notes: {
		write: 'new InvoiceService(logger).generate(json façon UBL de la librairie, { format: "CII", lang: "de-de" }), contrôle ajv du JSON inclus',
	},
	unsupported: {
		extract: 'pas de lecture de PDF',
		parse: 'pas de lecteur XML (entrée JSON ou tableur seulement)',
		validate: 'pas de validation XML (seul le JSON d’entrée est contrôlé par un schéma ajv)',
		embed: 'pas d’injection d’un XML existant : les formats Factur-X-* génèrent le XML depuis le JSON puis l’attachent à un PDF fourni (ou rendu par LibreOffice)',
		check: 'pas de lecture de PDF',
	},
	write: {
		model(s, bench) {
			smallOnly(bench);
			return eInvoiceEuInvoice();
		},
		write: (inv) => service.generate(inv, { format: 'CII', lang: 'de-de' }),
		handBuilt: true,
	},
};
