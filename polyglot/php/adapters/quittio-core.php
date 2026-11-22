<?php

/*
 * quittio/core (+ quittio/country-de): invoicing framework (EN 16931 data
 * model, calculation, CII/UBL writers and readers, ZUGFeRD/Factur-X PDF/A-3
 * rendering with mPDF, validation levels L1 to L7). The German country pack
 * adds the XRechnung profiles and the BR-DE rules; it is registered here.
 *
 * Validation offline = the default levels of Quittio::create(): L1 structure,
 * L2 mandatory fields of the profile, L3 business rules (EN 16931 rules coded
 * in PHP, BR-DE of the pack), L4 XSD CII D16B, L7 plausibility (IBAN, VAT ids,
 * Leitweg-ID). L5 (Schematron) is the KoSIT validator (Java jar or daemon) and
 * L6 is veraPDF (Docker or binary): neither is configured.
 *
 * Embed: MpdfPostProcessor::process, the step of the render pipeline that
 * imports the pages of a PDF with FPDI (via mPDF) and attaches the XML as a
 * PDF/A-3 associated file; mPDF's cache goes to the bench temp dir.
 */

declare(strict_types=1);

use Quittio\Country\DE\CountryDE;
use Quittio\Core\Pdf\EmbeddedXml;
use Quittio\Core\Pdf\PdfEmbeddedFileExtractor;
use Quittio\Core\Pdf\PostProcessing\MpdfPostProcessor;
use Quittio\Core\Quittio;
use Quittio\Core\Reading\ReadResult;
use Quittio\Core\Validation\ValidationMessage;
use Quittio\Core\Xml\Cii\CiiWriter;

$q = null;
$get = static function () use (&$q): Quittio {
    return $q ??= Quittio::create(config: ['temp_dir' => \Bench\tmp_dir() . '/quittio'], countryPacks: [new CountryDE()]);
};

return [
    'pkg' => 'quittio/core',
    'notes' => [
        'extract' => 'PdfEmbeddedFileExtractor::extractFile()->xml',
        'parse' => 'Quittio::read (CII reader + Germany pack; recomputes the totals)',
        'validate' => 'Quittio::validateFile, default offline levels: L1-L3 PHP rules (EN 16931, BR-DE of the country-de pack), L4 CII D16B XSD, L7 plausibility; L5 schematron (KoSIT, Java) and L6 veraPDF not configured; errors = ValidationResult::errors()',
        'write' => 'CiiWriter::serialize of the data read by Quittio::read, profile recognised while parsing',
        'embed' => 'MpdfPostProcessor::process (pages re-imported by FPDI through mPDF, XML as a PDF/A-3 associated file)',
    ],
    'unsupported' => [
        'check_pdf' => 'L6 veraPDF only (Docker or external binary)',
    ],
    'survey_notes' => [
        'validate' => 'levels L1-L4 and L7 (EN 16931 rules in PHP, BR-DE, XSD); no schematron (KoSIT requires Java)',
        'parse' => 'Quittio::read',
        'extract' => 'PdfEmbeddedFileExtractor::extractFile',
    ],

    'extract' => static function (string $pdf, ?string $path): string {
        $f = (new PdfEmbeddedFileExtractor())->extractFile($pdf);
        if ($f === null) {
            throw new \RuntimeException('no invoice attachment');
        }

        return $f->xml;
    },

    'parse' => static fn (string $xml, ?string $path): ReadResult => $get()->read($xml),

    'validate' => static function (string $xml, ?string $path) use ($get): array {
        $errors = $get()->validateFile($xml)->errors();

        return [\count($errors), array_map(static fn (ValidationMessage $m) => $m->code, $errors)];
    },

    'write' => static function (ReadResult $r) use ($get): string {
        $profiles = $get()->profiles();
        $profile = $profiles->get($r->profile !== null && $profiles->has($r->profile) ? $r->profile : Quittio::DEFAULT_PROFILE);

        return (new CiiWriter())->serialize($r->data, $profile);
    },

    'embed' => static fn (string $host, string $xml): string => (new MpdfPostProcessor(\Bench\tmp_dir() . '/quittio-mpdf'))->process($host, new EmbeddedXml($xml), 'Invoice'),
];
