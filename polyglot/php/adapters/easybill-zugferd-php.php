<?php

/*
 * easybill/zugferd-php: JMS-serialized object model of ZUGFeRD 1 and of
 * ZUGFeRD 2 / Factur-X (CII), Reader and Builder. Validation is the XSD
 * Factur-X 1.08.0 chosen by the caller (Validator::validateAgainstXsd takes
 * the schema file): the adapter maps BT-24 to the profile's schema, EN16931
 * for XRechnung. No PDF.
 */

declare(strict_types=1);

use Easybill\ZUGFeRD2\Builder;
use Easybill\ZUGFeRD2\Model\CrossIndustryInvoice;
use Easybill\ZUGFeRD2\Reader;
use Easybill\ZUGFeRD2\Validator;

$schemaFor = static function (string $xml): string {
    $id = preg_match('~GuidelineSpecifiedDocumentContextParameter>\s*<(?:\w+:)?ID>([^<]+)<~', $xml, $m) ? strtolower($m[1]) : '';

    return match (true) {
        str_contains($id, 'extended') => Validator::SCHEMA_EXTENDED,
        str_contains($id, 'basicwl') => Validator::SCHEMA_BASIC_WL,
        str_contains($id, 'basic') => Validator::SCHEMA_BASIC,
        str_contains($id, 'minimum') => Validator::SCHEMA_MINIMUM,
        default => Validator::SCHEMA_EN16931,
    };
};

return [
    'pkg' => 'easybill/zugferd-php',
    'notes' => [
        'parse' => 'Easybill\\ZUGFeRD2\\Reader::transform (JMS Serializer)',
        'validate' => 'Validator::validateAgainstXsd: Factur-X 1.08.0 XSD, schema of the profile chosen by the adapter from BT-24 (EN16931 for XRechnung)',
        'write' => 'Easybill\\ZUGFeRD2\\Builder::transform of the model read by Reader',
    ],
    'unsupported' => [
        'extract' => 'no PDF',
        'embed' => 'no PDF',
        'check_pdf' => 'no PDF',
    ],
    'survey_notes' => [
        'validate' => 'Factur-X 1.08.0 XSD only',
        'parse' => 'Reader::transform',
    ],

    'parse' => static fn (string $xml, ?string $path): CrossIndustryInvoice => Reader::create()->transform($xml),

    'validate' => static function (string $xml, ?string $path) use ($schemaFor): array {
        $err = (new Validator())->validateAgainstXsd($xml, $schemaFor($xml));
        if ($err === null) {
            return [0, []];
        }
        $lines = array_values(array_filter(array_map('trim', explode("\n", $err))));

        return [max(1, \count($lines)), array_map(static fn ($l) => 'XSD: ' . mb_substr($l, 0, 80), $lines)];
    },

    'write' => static fn (CrossIndustryInvoice $m): string => Builder::create()->transform($m),
];
