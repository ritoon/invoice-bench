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
		write: 'new InvoiceService(logger).generate(UBL-style JSON of the library, { format: "CII", lang: "de-de" }), ajv check of the JSON included',
	},
	unsupported: {
		extract: 'no PDF reading',
		parse: 'no XML reader (JSON or spreadsheet input only)',
		validate: 'no XML validation (only the input JSON is checked, by an ajv schema)',
		embed: 'no embedding of an existing XML: the Factur-X-* formats generate the XML from the JSON, then attach it to a supplied PDF (or one rendered by LibreOffice)',
		check: 'no PDF reading',
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
