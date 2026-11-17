//! zugferd 0.3.0 (FaultierSP): ZUGFeRD XML generation from a builder, levels
//! MINIMUM, BASIC WL and BASIC. No parser, no PDF, no XML validation.
//!
//! Write_Small serialises a model built once, by hand, from the content of
//! EN16931_Einfach at the BASIC level (the highest one the crate generates).

use chrono::NaiveDate;
use harness::{Lib, Op};
use zugferd::components::enums::identifier_scheme_code::IdentifierSchemeCode;
use zugferd::components::enums::unit_code::UnitCode;
use zugferd::{
    ApplicableTradeTax, AssociatedDocumentLineDocument, BilledQuantity, CountryCode, CurrencyCode, GlobalID,
    IncludedSupplyChainTradeLineItem, Invoice, InvoiceBuilder, InvoiceTypeCode, NetPriceProductTradePrice,
    SpecificationLevel, SpecifiedLineTradeAgreement, SpecifiedLineTradeDelivery, SpecifiedLineTradeSettlement,
    SpecifiedTradeProduct, SpecifiedTradeSettlementLineMonetarySummation, VATCategoryCode,
};

const LIB: Lib = Lib { lib: "zugferd", pkg: "zugferd", version: "0.3.0" };

fn line(
    id: &'static str,
    gtin: &'static str,
    name: &'static str,
    price: f64,
    qty: f64,
    rate: f64,
    total: f64,
) -> IncludedSupplyChainTradeLineItem<'static> {
    IncludedSupplyChainTradeLineItem {
        associated_document_line_document: AssociatedDocumentLineDocument { line_id: id, included_note: None },
        specified_trade_product: SpecifiedTradeProduct {
            global_id: Some(GlobalID { value: gtin, scheme_id: IdentifierSchemeCode::GTIN }),
            name,
            description: None,
        },
        specified_line_trade_agreement: SpecifiedLineTradeAgreement {
            gross_price_product_trade_price: None,
            net_price_product_trade_price: NetPriceProductTradePrice { charge_amount: price },
        },
        specified_line_trade_delivery: SpecifiedLineTradeDelivery {
            billed_quantity: BilledQuantity { value: qty, unit_code: UnitCode::Piece },
        },
        specified_line_trade_settlement: SpecifiedLineTradeSettlement {
            applicable_trade_tax: ApplicableTradeTax {
                calculated_amount: None,
                type_code: "VAT",
                category_code: VATCategoryCode::StandardRate,
                basis_amount: None,
                rate_applicable_percent: Some(rate),
                exemption_reason: None,
                exemption_reason_code: None,
            },
            specified_trade_settlement_line_monetary_summation: SpecifiedTradeSettlementLineMonetarySummation {
                line_total_amount: total,
            },
            specified_trade_allowance_charge: Vec::new(),
        },
    }
}

/// EN16931_Einfach (invoice 471102 of 2018-03-05), by hand. The builder
/// requires a buyer order reference (BT-13) and a payment due date that the
/// sample does not carry: "-" and the date of its payment terms are used.
fn einfach() -> Result<Invoice<'static>, String> {
    let d = |y, m, dd| NaiveDate::from_ymd_opt(y, m, dd).expect("date");
    let mut b = InvoiceBuilder::new();
    b.set_invoice_type_code(InvoiceTypeCode::CommercialInvoice)
        .set_invoice_nr("471102")
        .set_date_of_issue(d(2018, 3, 5))
        .set_invoice_notes(vec!["Rechnung gemäß Bestellung vom 01.03.2018."])
        .set_sellers_name("Lieferant GmbH")
        .set_sellers_postal_trade_address_postcode_code("80333")
        .set_sellers_postal_trade_address_line_one("Lieferantenstraße 20")
        .set_sellers_postal_trade_address_city_name("München")
        .set_sellers_postal_trade_address_country_code(CountryCode::Germany)
        .set_sellers_specified_tax_registration("DE123456789")
        .set_buyers_name("Kunden AG Mitte")
        .set_buyers_postal_trade_address_postcode_code("69876")
        .set_buyers_postal_trade_address_line_one("Kundenstraße 15")
        .set_buyers_postal_trade_address_city_name("Frankfurt")
        .set_buyers_postal_trade_address_country_code(CountryCode::Germany)
        .set_buyers_order_specified_document("-")
        .set_occurrence_date(d(2018, 3, 5))
        .set_invoice_currency_code(CurrencyCode::Euro)
        .set_specified_trade_payment_terms_due_date(d(2018, 4, 4))
        .add_applicable_trade_tax(ApplicableTradeTax {
            calculated_amount: Some(19.25),
            basis_amount: Some(275.0),
            rate_applicable_percent: Some(7.0),
            category_code: VATCategoryCode::StandardRate,
            ..Default::default()
        })
        .add_applicable_trade_tax(ApplicableTradeTax {
            calculated_amount: Some(37.62),
            basis_amount: Some(198.0),
            rate_applicable_percent: Some(19.0),
            category_code: VATCategoryCode::StandardRate,
            ..Default::default()
        })
        .set_monetary_summation_line_total_amount(473.0)
        .set_monetary_summation_charge_total_amount(0.0)
        .set_monetary_summation_allowance_total_amount(0.0)
        .set_monetary_summation_tax_basis_total_amount(473.0)
        .set_monetary_summation_tax_total_amount(56.87)
        .set_monetary_summation_grand_total_amount(529.87)
        .set_monetary_summation_due_payable_amount(529.87)
        .add_supply_chain_trade_line_item(line("1", "4012345001235", "Trennblätter A4", 9.9, 20.0, 19.0, 198.0))
        .add_supply_chain_trade_line_item(line("2", "4000050986428", "Joghurt Banane", 5.5, 50.0, 7.0, 275.0));
    b.build(SpecificationLevel::Basic)
}

fn main() {
    harness::run(LIB, |h| {
        h.unsupported(Op::Extract, "pas de lecture de PDF");
        h.unsupported(Op::Parse, "pas de lecteur XML (génération seulement)");
        h.unsupported(Op::Validate, "pas de lecteur XML ; validate_business_rules ne s'applique qu'au modèle construit");
        h.skip("Write_Lines", "modèle construit à la main seulement pour EN16931_Einfach (pas de lecteur XML)");
        h.write(
            "modèle construit à la main (contenu de EN16931_Einfach, niveau BASIC : le crate ne génère pas EN 16931) ; Invoice::to_xml_string (serde + quick-xml)",
            |_x| einfach(),
            |inv: &Invoice<'static>| inv.to_xml_string(),
        );
        h.unsupported(Op::Embed, "pas de PDF");
        h.unsupported(Op::CheckPdf, "pas de PDF");
    });
}
