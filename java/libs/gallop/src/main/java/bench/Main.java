package bench;

import java.math.BigDecimal;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;

import bench.harness.Harness;
import bench.harness.Harness.Lib;
import bench.harness.Harness.Op;
import de.codebarista.gallop.EInvoiceWriter;
import de.codebarista.gallop.model.DeliveryInformation;
import de.codebarista.gallop.model.Invoice;
import de.codebarista.gallop.model.InvoiceNote;
import de.codebarista.gallop.model.Item;
import de.codebarista.gallop.model.PaymentInstructions;
import de.codebarista.gallop.model.PostalAddress;
import de.codebarista.gallop.model.SellerOrBuyer;
import de.codebarista.gallop.model.TaxCategory;
import de.codebarista.gallop.model.Vat;

/**
 * gallop 3.0.0 (codebarista/gallop): a CII writer for XRechnung, ZUGFeRD and
 * Factur-X EN 16931 from a fluent model, on the JDK's DOM. Writer only: no
 * reader, no validation, no PDF.
 */
public final class Main {

    static final Lib LIB = new Lib("gallop", "de.codebarista:gallop", "3.0.0");

    static BigDecimal d(String s) {
        return new BigDecimal(s);
    }

    static OffsetDateTime day(int y, int m, int dd) {
        return OffsetDateTime.of(y, m, dd, 0, 0, 0, 0, ZoneOffset.UTC);
    }

    static PostalAddress address(String line, String zip, String city) {
        return PostalAddress.create().addressLineOne(line).zipCode(zip).city(city).countryIsoCode("DE");
    }

    static Item line(long id, String sellerId, String name, long qty, String price, String rate, String total) {
        return Item.create().id(id).sellerAssignedId(sellerId).name(name).quantity(qty).unitCode("H87")
                .unitPrice(d(price)).itemTotalNetAmount(d(total))
                .vat(Vat.create().rate(d(rate)).category(TaxCategory.STANDARD_RATE));
    }

    /** The content of EN16931_Einfach, built by hand (gallop cannot read
     *  XML). Quantities are whole numbers in gallop's model, which fits this
     *  invoice; the global trade item numbers (GTIN) have no field. */
    static Invoice einfach() {
        return Invoice.create()
                .documentId("471102")
                .issueDate(day(2018, 3, 5))
                .documentTypeCode("380")
                .currency("EUR")
                .addInvoiceNote(InvoiceNote.create().note("Rechnung gemäß Bestellung vom 01.03.2018."))
                .addInvoiceNote(InvoiceNote.create().note("Lieferant GmbH\nLieferantenstraße 20\n80333 München\n"
                        + "Deutschland\nGeschäftsführer: Hans Muster\nHandelsregisternummer: H A 123"))
                .seller(SellerOrBuyer.create().name("Lieferant GmbH").vatId("DE123456789")
                        .sellerTaxRegistrationIdentifier("201/113/40209")
                        .address(address("Lieferantenstraße 20", "80333", "München")))
                .buyer(SellerOrBuyer.create().name("Kunden AG Mitte")
                        .address(address("Kundenstraße 15", "69876", "Frankfurt")))
                // gallop writes a ship-to party (BG-13) with the delivery
                // date, and BR-57 then wants its country: the buyer's address.
                .deliveryInfo(DeliveryInformation.create().actualDeliveryDate(day(2018, 3, 5))
                        .deliveryAddress(address("Kundenstraße 15", "69876", "Frankfurt")))
                // gallop always writes a payment means group (BG-16), which
                // EN16931_Einfach does not have: code 1 ("instrument not
                // defined") keeps the output valid (BR-57, BR-CL-16).
                .paymentInstructions(PaymentInstructions.create().meansType("1").paymentTerms(
                        "Zahlbar innerhalb 30 Tagen netto bis 04.04.2018, 3% Skonto innerhalb 10 Tagen bis 15.03.2018"))
                .addItem(line(1, "TB100A4", "Trennblätter A4", 20, "9.90", "19.00", "198.00"))
                .addItem(line(2, "ARNR2", "Joghurt Banane", 50, "5.50", "7.00", "275.00"))
                .addVatTotal(Vat.create().rate(d("7.00")).category(TaxCategory.STANDARD_RATE)
                        .taxableAmount(d("275.00")).taxAmount(d("19.25")))
                .addVatTotal(Vat.create().rate(d("19.00")).category(TaxCategory.STANDARD_RATE)
                        .taxableAmount(d("198.00")).taxAmount(d("37.62")))
                .lineTotalAmount(d("473.00"))
                .chargeTotalAmount(d("0.00"))
                .allowanceTotalAmount(d("0.00"))
                .taxBasisTotalAmount(d("473.00"))
                .taxTotalAmount(d("56.87"))
                .grandTotalAmount(d("529.87"))
                .paidAmount(d("0.00"))
                .duePayableAmount(d("529.87"));
    }

    public static void main(String[] args) {
        Harness.run(LIB, args, h -> {
            h.unsupported(Op.EXTRACT, "no PDF parsing");
            h.unsupported(Op.PARSE, "no XML reader (generation only)");
            h.unsupported(Op.VALIDATE, "no validation");
            h.skip("Write_Lines", "hand-built model only for EN16931_Einfach (no XML reader)");
            h.write("hand-built model (content of EN16931_Einfach, without the GTINs or the REG code of the "
                    + "note: no field; payment means code 1 and delivery address = buyer's address "
                    + "added, since gallop always writes the BG-16 and BG-13 groups; output conforms to the CEN "
                    + "1.3.17 rules); EInvoiceWriter.generateFacturXXML (Factur-X EN 16931, JDK DOM)",
                    x -> einfach(), EInvoiceWriter::generateFacturXXML);
            h.unsupported(Op.EMBED, "no PDF (XML only)");
            h.unsupported(Op.CHECK_PDF, "no PDF");
        });
    }
}
