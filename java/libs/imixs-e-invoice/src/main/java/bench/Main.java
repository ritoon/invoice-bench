package bench;

import java.io.ByteArrayInputStream;

import bench.harness.Harness;
import bench.harness.Harness.Lib;
import bench.harness.Harness.Op;
import org.imixs.einvoice.EInvoiceModel;
import org.imixs.einvoice.EInvoiceModelFactory;

/**
 * Imixs e-invoice 1.0.3 (imixs/e-invoice): a read/edit wrapper over the DOM
 * of a CII, UBL or KSeF invoice (getters and setters on the main fields),
 * serialised back with getContent. No validation, no PDF.
 */
public final class Main {

    static final Lib LIB = new Lib("imixs-e-invoice", "org.imixs.util:imixs-e-invoice", "1.0.3");

    static EInvoiceModel parse(byte[] xml) throws Exception {
        return EInvoiceModelFactory.read(new ByteArrayInputStream(xml));
    }

    public static void main(String[] args) {
        Harness.run(LIB, args, h -> {
            h.unsupported(Op.EXTRACT, "no PDF parsing");
            h.parse("EInvoiceModelFactory.read(InputStream) → EInvoiceModelCII (full DOM, then parsing of the header, "
                    + "party and line fields by walking the nodes)", Main::parse);
            h.unsupported(Op.VALIDATE, "no validation");
            h.write("EInvoiceModel.getContent(): reserializes the (indented) DOM that the model wraps, model read "
                    + "by EInvoiceModelFactory.read", Main::parse, EInvoiceModel::getContent);
            h.unsupported(Op.EMBED, "no PDF writing");
            h.unsupported(Op.CHECK_PDF, "no PDF parsing");
        });
    }
}
