<?php

/*
 * daniel-jorg-schuppelius/php-erechnung-toolkit (AGPL-3.0-or-later): German
 * e-invoicing toolkit (XRechnung/ZUGFeRD model, UBL and CII parser and
 * generator, order and GAEB formats, Peppol). PDF features need the optional
 * package daniel-jorg-schuppelius/php-pdf-toolkit (AGPL too, installed here),
 * whose ZugferdReader runs poppler's pdfdetach (-list then -savefile to a
 * private temp dir).
 *
 * Validation offline: CiiSchemaValidator (XSD CII D16B) on the XML, then
 * DocumentRulesValidator on the parsed document (profile, supporting
 * documents, periods, Leitweg-ID: BR-29, BR-30, BR-52, BR-CL-24, BR-CO-19/20,
 * BR-DE-22 and TK-* checks). The EN 16931 / XRechnung Schematron goes through
 * KositValidator, which runs the bundled KoSIT jar with Java: not measured.
 * Embed: not measured: ZugferdPdfGenerator renders its own PDF (no host) and
 * ZugferdPdfAugmenter only adds supporting documents to a ZUGFeRD PDF.
 */

declare(strict_types=1);

use ERechnungToolkit\Entities\Document;
use ERechnungToolkit\Generators\ERechnungGenerator;
use ERechnungToolkit\Parsers\ERechnungParser;
use ERechnungToolkit\Parsers\ZugferdPdfParser;
use ERechnungToolkit\Validators\CiiSchemaValidator;
use ERechnungToolkit\Validators\DocumentRulesValidator;
use ERechnungToolkit\Validators\ValidationMessage;

return [
    'pkg' => 'daniel-jorg-schuppelius/php-erechnung-toolkit',
    'notes' => [
        'extract' => 'ZugferdPdfParser::extractXml($path) via php-pdf-toolkit (poppler\'s pdfdetach: -list then -savefile)',
        'parse' => 'ERechnungParser::parse',
        'validate' => 'CiiSchemaValidator (CII D16B XSD) then DocumentRulesValidator on the parsed document (profile, attachments, periods, Leitweg-ID: 7 BR rules and TK-* checks); EN 16931/XRechnung schematron through KositValidator (Java) not measured',
        'write' => 'ERechnungGenerator::generateCii of the document read by ERechnungParser',
    ],
    'unsupported' => [
        'embed' => 'no embedding into a host PDF (ZugferdPdfGenerator renders its own PDF; ZugferdPdfAugmenter only adds supporting documents)',
        'check_pdf' => 'no PDF check',
    ],
    'survey_notes' => [
        'validate' => 'CII D16B XSD + a few rules of the toolkit (BR-29, BR-30, BR-52, BR-CL-24, BR-CO-19/20, BR-DE-22); KoSIT requires Java',
        'parse' => 'ERechnungParser::parse',
        'extract' => 'ZugferdPdfParser::extractXml (pdfdetach)',
    ],

    'extract' => static function (string $pdf, ?string $path): string {
        $p = new ZugferdPdfParser();
        $xml = $path !== null ? $p->extractXml($path) : $p->extractXmlFromString($pdf);
        if ($xml === null) {
            throw new \RuntimeException('no XML invoice extracted (pdfdetach unavailable or no attachment)');
        }

        return $xml;
    },

    'parse' => static fn (string $xml, ?string $path): Document => (new ERechnungParser())->parse($xml),

    'validate' => static function (string $xml, ?string $path): array {
        $rules = array_map(static fn (string $e) => 'XSD: ' . mb_substr($e, 0, 80), (new CiiSchemaValidator())->validate($xml));
        $doc = (new ERechnungParser())->parse($xml);
        foreach ((new DocumentRulesValidator())->validate($doc)->getErrors() as $m) {
            /** @var ValidationMessage $m */
            $rules[] = $m->getCode();
        }

        return [\count($rules), $rules];
    },

    'write' => static fn (Document $d): string => (new ERechnungGenerator())->generateCii($d),
];
