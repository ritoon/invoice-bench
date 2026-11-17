// Hand-built models of corpus/XML-Rechnung/CII/EN16931_Einfach.cii.xml, one
// per input format, for the libraries that can write CII but cannot read it
// (PROTOCOL.md, Write_*: "à défaut, à la main à partir du contenu de
// EN16931_Einfach"). Same invoice everywhere: 471102 of 2018-03-05, two lines
// (20 x 9.90 at 19 %, 50 x 5.50 at 7 %), 473.00 net, 56.87 VAT, 529.87 due.
import { Unsupported } from './harness.mjs';

export const HAND_NOTE = 'modèle construit à la main';

// Write_Lines has no hand model: only the Small invoice is rebuilt by hand.
export function smallOnly(bench, why = 'pas de lecteur CII dans la librairie') {
	if (bench !== 'Write_Small') throw new Unsupported(`modèle à la main seulement pour Write_Small (${why})`);
}

const NOTE_REG =
	'Lieferant GmbH\t\t\t\t\nLieferantenstraße 20\t\t\t\t\n80333 München\t\t\t\t\nDeutschland\t\t\t\t\nGeschäftsführer: Hans Muster\nHandelsregisternummer: H A 123\n      ';
const NOTE_1 = 'Rechnung gemäß Bestellung vom 01.03.2018.';
const TERMS = 'Zahlbar innerhalb 30 Tagen netto bis 04.04.2018, 3% Skonto innerhalb 10 Tagen bis 15.03.2018';

const LINES = [
	{ id: '1', gtin: '4012345001235', sellerId: 'TB100A4', name: 'Trennblätter A4', qty: 20, price: 9.9, rate: 19, total: 198 },
	{ id: '2', gtin: '4000050986428', sellerId: 'ARNR2', name: 'Joghurt Banane', qty: 50, price: 5.5, rate: 7, total: 275 },
];

// @stackforge-eu/factur-x: FacturXInvoiceInput (buildXml).
export function stackforgeInput() {
	return {
		document: {
			id: '471102',
			issueDate: '2018-03-05',
			typeCode: '380',
			notes: [{ content: NOTE_1 }, { content: NOTE_REG, subjectCode: 'REG' }],
		},
		seller: {
			name: 'Lieferant GmbH',
			id: '549910',
			globalId: { value: '4000001123452', schemeID: '0088' },
			address: { line1: 'Lieferantenstraße 20', city: 'München', postalCode: '80333', country: 'DE' },
			taxRegistrations: [
				{ id: '201/113/40209', schemeId: 'FC' },
				{ id: 'DE123456789', schemeId: 'VA' },
			],
		},
		buyer: {
			name: 'Kunden AG Mitte',
			id: 'GE2020211',
			address: { line1: 'Kundenstraße 15', city: 'Frankfurt', postalCode: '69876', country: 'DE' },
		},
		lines: LINES.map((l) => ({
			id: l.id,
			name: l.name,
			quantity: l.qty,
			unitCode: 'H87',
			unitPrice: l.price,
			grossUnitPrice: l.price,
			lineTotal: l.total,
			vatCategoryCode: 'S',
			vatRatePercent: l.rate,
			standardIdentifier: { value: l.gtin, schemeID: '0160' },
			sellerAssignedId: l.sellerId,
		})),
		delivery: { date: '2018-03-05' },
		payment: { termsDescription: TERMS },
		vatBreakdown: [
			{ categoryCode: 'S', ratePercent: 7, taxableAmount: 275, taxAmount: 19.25 },
			{ categoryCode: 'S', ratePercent: 19, taxableAmount: 198, taxAmount: 37.62 },
		],
		totals: {
			lineTotal: 473,
			chargeTotal: 0,
			allowanceTotal: 0,
			taxBasisTotal: 473,
			taxTotal: 56.87,
			grandTotal: 529.87,
			prepaidAmount: 0,
			duePayableAmount: 529.87,
			currency: 'EUR',
		},
	};
}

// @ekwo-ai/factur-x: Invoice (generateCiiXml computes the totals and the VAT
// breakdown itself; its model has one free-text note and no item GTIN).
export function ekwoInvoice() {
	return {
		number: '471102',
		type: 'invoice',
		issueDate: '2018-03-05',
		deliveryDate: '2018-03-05',
		currency: 'EUR',
		seller: {
			name: 'Lieferant GmbH',
			address: { line1: 'Lieferantenstraße 20', postalCode: '80333', city: 'München', country: 'DE' },
			vatId: 'DE123456789',
			legalId: { scheme: '0088', value: '4000001123452' },
		},
		buyer: {
			name: 'Kunden AG Mitte',
			address: { line1: 'Kundenstraße 15', postalCode: '69876', city: 'Frankfurt', country: 'DE' },
		},
		lines: LINES.map((l) => ({
			id: l.id,
			name: l.name,
			quantity: l.qty,
			unitCode: 'H87',
			unitPrice: l.price,
			vatRate: l.rate,
			vatCategory: 'S',
			sellerItemId: l.sellerId,
		})),
		note: NOTE_1,
		payment: { terms: TERMS },
	};
}

// @e-invoice-eu/core: the library's UBL-shaped JSON (ajv schema), rendered to
// CII by InvoiceService.generate(…, { format: 'CII' }).
export function eInvoiceEuInvoice() {
	const eur = (k, v) => ({ [k]: v, [`${k}@currencyID`]: 'EUR' });
	const vat = (rate) => ({ 'cbc:ID': 'S', 'cbc:Percent': rate, 'cac:TaxScheme': { 'cbc:ID': 'VAT' } });
	return {
		'ubl:Invoice': {
			'cbc:CustomizationID': 'urn:cen.eu:en16931:2017',
			'cbc:ID': '471102',
			'cbc:IssueDate': '2018-03-05',
			'cbc:InvoiceTypeCode': '380',
			'cbc:Note': [NOTE_1, NOTE_REG],
			'cbc:DocumentCurrencyCode': 'EUR',
			'cac:AccountingSupplierParty': {
				'cac:Party': {
					'cac:PartyIdentification': [{ 'cbc:ID': '549910' }, { 'cbc:ID': '4000001123452', 'cbc:ID@schemeID': '0088' }],
					'cac:PostalAddress': {
						'cbc:StreetName': 'Lieferantenstraße 20',
						'cbc:CityName': 'München',
						'cbc:PostalZone': '80333',
						'cac:Country': { 'cbc:IdentificationCode': 'DE' },
					},
					'cac:PartyTaxScheme': [
						{ 'cbc:CompanyID': 'DE123456789', 'cac:TaxScheme': { 'cbc:ID': 'VAT' } },
						{ 'cbc:CompanyID': '201/113/40209', 'cac:TaxScheme': { 'cbc:ID': 'FC' } },
					],
					'cac:PartyLegalEntity': { 'cbc:RegistrationName': 'Lieferant GmbH' },
				},
			},
			'cac:AccountingCustomerParty': {
				'cac:Party': {
					'cac:PartyIdentification': { 'cbc:ID': 'GE2020211' },
					'cac:PostalAddress': {
						'cbc:StreetName': 'Kundenstraße 15',
						'cbc:CityName': 'Frankfurt',
						'cbc:PostalZone': '69876',
						'cac:Country': { 'cbc:IdentificationCode': 'DE' },
					},
					'cac:PartyLegalEntity': { 'cbc:RegistrationName': 'Kunden AG Mitte' },
				},
			},
			'cac:Delivery': { 'cbc:ActualDeliveryDate': '2018-03-05' },
			'cac:PaymentTerms': { 'cbc:Note': TERMS },
			'cac:TaxTotal': [
				{
					...eur('cbc:TaxAmount', '56.87'),
					'cac:TaxSubtotal': [
						{ ...eur('cbc:TaxableAmount', '275.00'), ...eur('cbc:TaxAmount', '19.25'), 'cac:TaxCategory': vat('7.00') },
						{ ...eur('cbc:TaxableAmount', '198.00'), ...eur('cbc:TaxAmount', '37.62'), 'cac:TaxCategory': vat('19.00') },
					],
				},
			],
			'cac:LegalMonetaryTotal': {
				...eur('cbc:LineExtensionAmount', '473.00'),
				...eur('cbc:TaxExclusiveAmount', '473.00'),
				...eur('cbc:TaxInclusiveAmount', '529.87'),
				...eur('cbc:AllowanceTotalAmount', '0.00'),
				...eur('cbc:ChargeTotalAmount', '0.00'),
				...eur('cbc:PrepaidAmount', '0.00'),
				...eur('cbc:PayableAmount', '529.87'),
			},
			'cac:InvoiceLine': LINES.map((l) => ({
				'cbc:ID': l.id,
				'cbc:InvoicedQuantity': l.qty.toFixed(4),
				'cbc:InvoicedQuantity@unitCode': 'H87',
				...eur('cbc:LineExtensionAmount', l.total.toFixed(2)),
				'cac:Item': {
					'cbc:Name': l.name,
					'cac:SellersItemIdentification': { 'cbc:ID': l.sellerId },
					'cac:StandardItemIdentification': { 'cbc:ID': l.gtin, 'cbc:ID@schemeID': '0160' },
					'cac:ClassifiedTaxCategory': vat(l.rate.toFixed(2)),
				},
				'cac:Price': eur('cbc:PriceAmount', l.price.toFixed(4)),
			})),
		},
	};
}

// node-zugferd: the EN16931 profile's input object (invoicer.create).
export function nodeZugferdData() {
	return {
		number: '471102',
		typeCode: '380',
		issueDate: new Date(Date.UTC(2018, 2, 5)),
		includedNote: [{ content: NOTE_1 }, { content: NOTE_REG, subjectCode: 'REG' }],
		transaction: {
			line: LINES.map((l) => ({
				identifier: l.id,
				tradeProduct: {
					globalIdentifier: { value: l.gtin, schemeIdentifier: '0160' },
					sellerAssignedID: l.sellerId,
					name: l.name,
				},
				tradeAgreement: {
					grossTradePrice: { chargeAmount: l.price.toFixed(4) },
					netTradePrice: { chargeAmount: l.price.toFixed(4) },
				},
				tradeDelivery: { billedQuantity: { amount: l.qty.toFixed(4), unitMeasureCode: 'H87' } },
				tradeSettlement: {
					tradeTax: { typeCode: 'VAT', categoryCode: 'S', rateApplicablePercent: l.rate.toFixed(2) },
					monetarySummation: { lineTotalAmount: l.total.toFixed(2) },
				},
			})),
			tradeAgreement: {
				seller: {
					identifier: ['549910'],
					globalIdentifier: { value: '4000001123452', schemeIdentifier: '0088' },
					name: 'Lieferant GmbH',
					postalAddress: { postCode: '80333', line1: 'Lieferantenstraße 20', city: 'München', countryCode: 'DE' },
					taxRegistration: { localIdentifier: '201/113/40209', vatIdentifier: 'DE123456789' },
				},
				buyer: {
					identifier: 'GE2020211',
					name: 'Kunden AG Mitte',
					postalAddress: { postCode: '69876', line1: 'Kundenstraße 15', city: 'Frankfurt', countryCode: 'DE' },
				},
			},
			tradeDelivery: { information: { deliveryDate: new Date(Date.UTC(2018, 2, 5)) } },
			tradeSettlement: {
				currencyCode: 'EUR',
				vatBreakdown: [
					{ calculatedAmount: '19.25', typeCode: 'VAT', basisAmount: '275.00', categoryCode: 'S', rateApplicablePercent: '7.00' },
					{ calculatedAmount: '37.62', typeCode: 'VAT', basisAmount: '198.00', categoryCode: 'S', rateApplicablePercent: '19.00' },
				],
				paymentTerms: { description: TERMS },
				monetarySummation: {
					lineTotalAmount: '473.00',
					chargeTotalAmount: '0.00',
					allowanceTotalAmount: '0.00',
					taxBasisTotalAmount: '473.00',
					taxTotal: { amount: '56.87', currencyCode: 'EUR' },
					grandTotalAmount: '529.87',
					paidAmount: '0.00',
					duePayableAmount: '529.87',
				},
			},
		},
	};
}
