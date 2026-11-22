<?php

/*
 * intragile/php-einvoicing: EN 16931 model with CII and UBL readers/writers,
 * a Factur-X PDF reader and an incremental-update PDF embedder written from
 * scratch (no FPDI), and validation by the official Schematron rules
 * pre-compiled to PHP (packages intragile/einvoicing-rules-*), run by its own
 * XPath 2.0 evaluator, after the XSD.
 *
 * Validation: XRechnungValidator::create() for BT-24 = EN 16931 or XRechnung
 * (the KoSIT scenario table: XSD CII D16B, CEN EN 16931 1.3.16 as pinned by
 * the KoSIT configuration, plus XRechnung schematron 2.6.0 and the KoSIT
 * custom levels for XRechnung); FacturXValidator::create() for the other
 * Factur-X profiles (Factur-X 1.09.2 schematron + XSD). Both are built once,
 * outside the timing. Errors = fatal findings + schema errors + XPath dynamic
 * errors, as ValidationReport::isValid.
 */

declare(strict_types=1);

use IntrAgile\EInvoicing\Model\Invoice;
use IntrAgile\EInvoicing\Syntax\Cii\CiiReader;
use IntrAgile\EInvoicing\Syntax\Cii\CiiWriter;
use IntrAgile\EInvoicing\Syntax\FacturX\FacturXPdfReader;
use IntrAgile\EInvoicing\Syntax\FacturX\FacturXPdfWriter;
use IntrAgile\EInvoicing\Validation\DynamicError;
use IntrAgile\EInvoicing\Validation\FacturX\FacturXValidator;
use IntrAgile\EInvoicing\Validation\Finding;
use IntrAgile\EInvoicing\Validation\SchemaError;
use IntrAgile\EInvoicing\Validation\XRechnung\XRechnungValidator;

$validators = [];
$pick = static function (string $xml) use (&$validators): XRechnungValidator|FacturXValidator {
    $id = preg_match('~GuidelineSpecifiedDocumentContextParameter>\s*<(?:\w+:)?ID>([^<]+)<~', $xml, $m) ? trim($m[1]) : '';
    // A UBL document (conformance cases) goes to the KoSIT scenarios, which
    // cover UBL Invoice and CreditNote; FacturXValidator reads CII only.
    $kosit = $id === 'urn:cen.eu:en16931:2017' || str_contains($id, 'xrechnung')
        || \in_array(\Bench\invoice_root($xml), ['Invoice', 'CreditNote'], true);

    return $kosit
        ? ($validators['kosit'] ??= XRechnungValidator::create())
        : ($validators['facturx'] ??= FacturXValidator::create());
};

return [
    'pkg' => 'intragile/php-einvoicing',
    'notes' => [
        'extract' => 'FacturXPdfReader::read()->xml (PDF reader written in the library)',
        'parse' => 'CiiReader::read()->invoice (EN 16931 model)',
        'validate' => 'XRechnungValidator (KoSIT scenarios: CII D16B XSD + CEN 1.3.16 schematron precompiled to PHP, + XRechnung 2.6.0 and KoSIT levels if BT-24 is XRechnung; the UBL scenarios for a UBL document); FacturXValidator for the other Factur-X profiles; errors = fatal + XSD + dynamic errors',
        'write' => 'CiiWriter::write of the model read by CiiReader',
        'embed' => 'FacturXPdfWriter::embed (incremental update, without FPDI; does not convert to PDF/A)',
    ],
    'unsupported' => [
        'check_pdf' => 'no PDF/A check (only warnings when embedding)',
    ],
    'survey_notes' => [
        'validate' => 'KoSIT scenarios: XSD + CEN 1.3.16 (+ XRechnung 2.6.0 and KoSIT levels), rules precompiled to PHP',
        'parse' => 'CiiReader::read',
        'extract' => 'FacturXPdfReader::read',
    ],

    'extract' => static fn (string $pdf, ?string $path): string => (new FacturXPdfReader())->read($pdf)->xml,

    'parse' => static fn (string $xml, ?string $path): Invoice => (new CiiReader())->read($xml)->invoice,

    'validate' => static function (string $xml, ?string $path) use ($pick): array {
        $r = $pick($xml)->validate($xml);
        $rules = [
            ...array_map(static fn (SchemaError $e) => 'XSD: ' . mb_substr($e->message, 0, 80), $r->schemaErrors),
            ...array_map(static fn (DynamicError $e) => 'XPATH: ' . ($e->ruleId ?? $e->errorCode), $r->dynamicErrors),
            ...array_map(static fn (Finding $f) => $f->ruleId, $r->fatal()),
        ];

        return [\count($rules), $rules];
    },

    'write' => static fn (Invoice $inv): string => (new CiiWriter())->write($inv),

    'embed' => static fn (string $host, string $xml): string => (new FacturXPdfWriter())->embed($host, $xml)->pdf,
];
