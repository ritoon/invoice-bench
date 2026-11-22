<?php

/*
 * tiime/factur-x: rewrite of atgp/factur-x for PHP 8.3 (extraction with
 * smalot/pdfparser, injection with FPDI, XSD Factur-X 1.07.3). The Schematron
 * and XSLT files under schematron/ are shipped but never used by src/. The
 * extractor only looks for an attachment named factur-x.xml.
 */

declare(strict_types=1);

use Tiime\FacturX\ProfileExtractor;
use Tiime\FacturX\Reader;
use Tiime\FacturX\Writer;
use Tiime\FacturX\XsdProcessor;
use Tiime\FacturX\XsdValidationException;

return [
    'pkg' => 'tiime/factur-x',
    'notes' => [
        'extract' => 'Reader::extractXML($pdf, validateXsd: false); only an attachment named factur-x.xml is looked for',
        'validate' => 'ProfileExtractor + XsdProcessor: Factur-X 1.07.3 XSD of the profile (schematron shipped but never called)',
        'embed' => 'Writer::generate, default options (XSD included, FPDI: the PDF is rewritten, PDF/A-3 XMP)',
    ],
    'unsupported' => [
        'parse' => 'no model (the model is tiime/cross-industry-invoice, measured separately)',
        'write' => 'no model',
        'check_pdf' => 'no native check of a hybrid PDF',
    ],
    'survey_notes' => [
        'validate' => 'Factur-X 1.07.3 XSD only',
        'extract' => 'Reader::extractXML(validateXsd: false), factur-x.xml attachment only',
    ],

    'extract' => static fn (string $pdf, ?string $path): string => (new Reader())->extractXML($pdf, false),

    'validate' => static function (string $xml, ?string $path): array {
        $doc = new \DOMDocument();
        $doc->loadXML($xml);
        $profile = ProfileExtractor::process($doc);
        try {
            (new XsdProcessor($profile))->validate($xml);
        } catch (XsdValidationException $e) {
            $errors = $e->xmlErrors;

            return [max(1, \count($errors)), array_map(static fn ($er) => 'XSD: ' . mb_substr(trim($er instanceof \LibXMLError ? $er->message : (string) $er), 0, 80), $errors ?: ['?'])];
        }

        return [0, []];
    },

    'embed' => static fn (string $host, string $xml): string => (new Writer())->generate($host, $xml),
];
