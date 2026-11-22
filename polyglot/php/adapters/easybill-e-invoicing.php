<?php

/*
 * easybill/e-invoicing: JMS-serialized models of CII and UBL (EN 16931 and its
 * CIUS), Reader and Transformer. No validation, no PDF.
 */

declare(strict_types=1);

use easybill\eInvoicing\Reader;
use easybill\eInvoicing\Transformer;

return [
    'pkg' => 'easybill/e-invoicing',
    'notes' => [
        'parse' => 'Reader::read (JMS Serializer; comments removed by regex, DOM reloaded); error if ReaderResult::isError()',
        'write' => 'Transformer::transformToXml of the model read by Reader',
    ],
    'unsupported' => [
        'extract' => 'no PDF',
        'validate' => 'no validation',
        'embed' => 'no PDF',
        'check_pdf' => 'no PDF',
    ],
    'survey_notes' => [
        'parse' => 'Reader::read',
    ],

    'parse' => static function (string $xml, ?string $path): object {
        $r = Reader::create()->read($xml);
        if ($r->isError()) {
            throw $r->getError();
        }

        return $r->getDocument();
    },

    'write' => static fn (object $m): string => Transformer::create()->transformToXml($m),
];
