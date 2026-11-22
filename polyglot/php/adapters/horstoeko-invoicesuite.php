<?php

/*
 * horstoeko/invoicesuite: successor of horstoeko/zugferd by the same author,
 * one DTO for many formats (ZUGFeRD/Factur-X, XRechnung CII and UBL, Peppol,
 * FatturaPA, ...). It declares a Composer conflict with horstoeko/zugferd, so
 * it lives in its own environment (invoicesuite/vendor).
 *
 * Offline validation is XSD only (InvoiceSuiteXsdDocumentValidator, schema of
 * the format provider recognised in the content). The two other validators are
 * not measured: InvoiceSuiteKositDocumentValidator downloads the KoSIT jar and
 * runs Java, InvoiceSuiteDocuflairDocumentValidator posts the invoice to a
 * remote web service.
 */

declare(strict_types=1);

use horstoeko\invoicesuite\InvoiceSuiteDocumentReader;
use horstoeko\invoicesuite\InvoiceSuitePdfDocumentBuilder;
use horstoeko\invoicesuite\InvoiceSuitePdfDocumentReader;
use horstoeko\invoicesuite\utils\InvoiceSuiteMessageBagItem;
use horstoeko\invoicesuite\validators\InvoiceSuiteXsdDocumentValidator;

return [
    'pkg' => 'horstoeko/invoicesuite',
    'notes' => [
        'extract' => 'InvoiceSuitePdfDocumentReader::createFromContent (prinsfrank/pdfparser; the reader also deserializes the XML it finds) then getInvoiceDocumentAttachment()->getAttachmentContent()',
        'parse' => 'InvoiceSuiteDocumentReader::createFromContent (format provider guessed, JMS Serializer)',
        'validate' => 'InvoiceSuiteXsdDocumentValidator::createFromContent()->validate(): XSD of the recognised format provider (Factur-X EN16931 for XRechnung CII); KoSIT (Java) and Docuflair (remote service) validators not measured',
        'write' => 'InvoiceSuiteDocumentReader::copyToBuilder()->getContent() (DTO then serialization)',
        'embed' => 'InvoiceSuitePdfDocumentBuilder::createFromDocumentContentAndPdfContent()->generatePdfDocumentAndGetContent() (FPDI/TCPDF)',
    ],
    'unsupported' => [
        'check_pdf' => 'no PDF/A check',
    ],
    'survey_notes' => [
        'validate' => 'XSD only (KoSIT requires Java, Docuflair is a remote service)',
        'parse' => 'InvoiceSuiteDocumentReader::createFromContent',
        'extract' => 'InvoiceSuitePdfDocumentReader::createFromContent (fails if the XML found matches no format provider)',
    ],

    'extract' => static function (string $pdf, ?string $path): string {
        $att = InvoiceSuitePdfDocumentReader::createFromContent($pdf)->getInvoiceDocumentAttachment();
        if ($att === null) {
            throw new \RuntimeException('no invoice attachment');
        }

        return $att->getAttachmentContent();
    },

    'parse' => static fn (string $xml, ?string $path): InvoiceSuiteDocumentReader => InvoiceSuiteDocumentReader::createFromContent($xml),

    'validate' => static function (string $xml, ?string $path): array {
        $v = InvoiceSuiteXsdDocumentValidator::createFromContent($xml)->validate();
        $msgs = array_merge($v->getErrorMessagesInMessageBag(), $v->getInternalErrorMessagesInMessageBag());

        return [\count($msgs), array_map(static fn (InvoiceSuiteMessageBagItem $m) => 'XSD: ' . mb_substr($m->getMessageContent(), 0, 80), $msgs)];
    },

    'write' => static fn (InvoiceSuiteDocumentReader $r): string => $r->copyToBuilder()->getContent(),

    'embed' => static fn (string $host, string $xml): string => InvoiceSuitePdfDocumentBuilder::createFromDocumentContentAndPdfContent($xml, $host)->generatePdfDocumentAndGetContent(),
];
