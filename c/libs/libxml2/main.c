/*
 * Floor: libxml2 2.15.4, the XML library of the C world (and of lxml), with no
 * knowledge of Factur-X. Parse = DOM (xmlCtxtReadMemory), Write = DOM
 * serialization (xmlDocDumpMemory), ParseValidate = DOM + XML Schema
 * validation against the UN/CEFACT CII D16B schema (xmlSchemaValidateDoc).
 * The schema only: no EN 16931 business rule. A UBL invoice or credit note
 * (the conformance cases, ../cases) is validated against the UBL 2.1
 * schema of the same configuration instead, chosen by the root element.
 */
#include <libxml/parser.h>
#include <libxml/tree.h>
#include <libxml/xmlerror.h>
#include <libxml/xmlschemas.h>
#include <libxml/xmlversion.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "harness.h"

static const fx_lib LIB = {"libxml2", "libxml2", LIBXML_DOTTED_VERSION};

/* The CII D16B and UBL 2.1 schemas of the KoSIT XRechnung 3.0.2
 * configuration (2026-08-31), fetched by run.sh setup. */
static const char *const XSD_REL[3] = {
	"kosit-config/resources/cii/16b/xsd/CrossIndustryInvoice_100pD16B.xsd",
	"kosit-config/resources/ubl/2.1/xsd/maindoc/UBL-Invoice-2.1.xsd",
	"kosit-config/resources/ubl/2.1/xsd/maindoc/UBL-CreditNote-2.1.xsd",
};
#define MAX_MESSAGES 20

typedef struct {
	char *msgs[MAX_MESSAGES];
	size_t n, errors;
} errlist;

static void collect(void *user, const xmlError *err) {
	errlist *l = user;
	if (!err || err->level < XML_ERR_ERROR)
		return;
	l->errors++;
	if (l->n < MAX_MESSAGES) {
		char buf[256];
		snprintf(buf, sizeof buf, "XSD:%d: %s", err->line, err->message ? err->message : "?");
		size_t k = strlen(buf);
		while (k && (buf[k - 1] == '\n' || buf[k - 1] == ' '))
			buf[--k] = 0;
		l->msgs[l->n++] = strdup(buf);
	}
}

static void errlist_free(errlist *l) {
	for (size_t i = 0; i < l->n; i++)
		free(l->msgs[i]);
	free(l);
}

/* XML -> DOM, the parser errors kept in `l` (or ignored when NULL). */
static xmlDocPtr read_doc(const unsigned char *xml, size_t len, errlist *l, fx_err *err) {
	xmlParserCtxtPtr p = xmlNewParserCtxt();
	if (!p) {
		fx_fail(err, "xmlNewParserCtxt: out of memory");
		return NULL;
	}
	errlist local = {0};
	xmlCtxtSetErrorHandler(p, collect, l ? (void *)l : (void *)&local);
	xmlDocPtr doc = xmlCtxtReadMemory(p, (const char *)xml, (int)len, "invoice.xml", NULL, XML_PARSE_NONET);
	if (!doc) {
		errlist *src = l ? l : &local;
		fx_fail(err, "xmlCtxtReadMemory: %s", src->n ? src->msgs[0] : "not well-formed");
	}
	for (size_t i = 0; i < local.n; i++)
		free(local.msgs[i]);
	xmlFreeParserCtxt(p);
	return doc;
}

static int parse(void *ctx, const unsigned char *xml, size_t len, void **model, fx_err *err) {
	(void)ctx;
	xmlDocPtr doc = read_doc(xml, len, NULL, err);
	if (!doc)
		return -1;
	*model = doc;
	return 0;
}

static void free_doc(void *ctx, void *m) {
	(void)ctx;
	xmlFreeDoc(m);
}

static void release_xml(void *p) { xmlFree(p); }

static int write_doc(void *ctx, void *model, fx_bytes *out, fx_err *err) {
	(void)ctx;
	xmlChar *mem = NULL;
	int size = 0;
	xmlDocDumpMemory(model, &mem, &size);
	if (!mem)
		return fx_fail(err, "xmlDocDumpMemory failed");
	out->data = mem;
	out->len = (size_t)size;
	out->owner = mem;
	out->release = release_xml;
	return 0;
}

/* The compiled schemas, each built on its first validation (the CII one's
 * cost is in first_call_ns of ParseValidate_Small), then reused. */
static xmlSchemaPtr schemas[3];
static char schema_errors[3][1024];

static xmlSchemaPtr compiled(int which, fx_err *err) {
	xmlSchemaPtr schema = schemas[which];
	char *schema_error = schema_errors[which];
	if (schema)
		return schema;
	if (schema_error[0]) {
		fx_fail(err, "%s", schema_error);
		return NULL;
	}
	char path[4096];
	snprintf(path, sizeof path, "%s/%s", fx_artefacts(), XSD_REL[which]);
	xmlSchemaParserCtxtPtr pc = xmlSchemaNewParserCtxt(path);
	errlist l = {0};
	if (pc) {
		xmlSchemaSetParserStructuredErrors(pc, collect, &l);
		schema = xmlSchemaParse(pc);
		xmlSchemaFreeParserCtxt(pc);
	}
	schemas[which] = schema;
	if (!schema)
		snprintf(schema_error, sizeof schema_errors[which], "cannot compile %s (run ./run.sh setup): %s", path,
			 l.n ? l.msgs[0] : "missing file");
	for (size_t i = 0; i < l.n; i++)
		free(l.msgs[i]);
	if (!schema)
		fx_fail(err, "%s", schema_error);
	return schema;
}

static int validate(void *ctx, const unsigned char *xml, size_t len, void **report, fx_err *err) {
	(void)ctx;
	errlist *l = calloc(1, sizeof *l);
	xmlDocPtr doc = read_doc(xml, len, l, err);
	if (!doc) {
		errlist_free(l);
		return -1;
	}
	xmlNodePtr root = xmlDocGetRootElement(doc);
	const char *name = root ? (const char *)root->name : "";
	xmlSchemaPtr s = compiled(strcmp(name, "Invoice") == 0 ? 1 : strcmp(name, "CreditNote") == 0 ? 2 : 0, err);
	if (!s) {
		xmlFreeDoc(doc);
		errlist_free(l);
		return -1;
	}
	xmlSchemaValidCtxtPtr vc = xmlSchemaNewValidCtxt(s);
	xmlSchemaSetValidStructuredErrors(vc, collect, l);
	int r = xmlSchemaValidateDoc(vc, doc);
	xmlSchemaFreeValidCtxt(vc);
	xmlFreeDoc(doc);
	if (r < 0) {
		errlist_free(l);
		return fx_fail(err, "xmlSchemaValidateDoc: internal error %d", r);
	}
	if (r > 0 && l->errors == 0)
		l->errors = 1;
	*report = l;
	return 0;
}

static void summarize(void *ctx, void *report, fx_findings *out) {
	(void)ctx;
	errlist *l = report;
	out->errors = l->errors;
	for (size_t i = 0; i < l->n; i++)
		fx_findings_add(out, l->msgs[i]);
}

static void free_report(void *ctx, void *r) {
	(void)ctx;
	errlist_free(r);
}

static void body(fx_harness *h, void *ctx) {
	(void)ctx;
	fx_unsupported(h, FX_EXTRACT, "XML library only");
	fx_parse(h, "floor: libxml2 DOM (xmlCtxtReadMemory, XML_PARSE_NONET), not an invoice model", parse,
		 free_doc, NULL);
	fx_validate(h,
		    "floor: DOM + xmlSchemaValidateDoc, UN/CEFACT CII D16B XSD (UBL 2.1 for a UBL document; KoSIT "
		    "XRechnung 3.0.2 configuration of 2026-08-31), schema compiled on the first call; schema only, no "
		    "EN 16931 business rules",
		    validate, summarize, free_report, NULL);
	fx_write(h, "floor: serialization of the libxml2 DOM (xmlDocDumpMemory), not an invoice model", parse,
		 free_doc, write_doc, NULL);
	fx_unsupported(h, FX_EMBED, "XML library only");
	fx_unsupported(h, FX_CHECK_PDF, "XML library only");
}

int main(int argc, char **argv) {
	xmlInitParser();
	int r = fx_run("c", LIB, argc, argv, body, NULL);
	for (int i = 0; i < 3; i++)
		if (schemas[i])
			xmlSchemaFree(schemas[i]);
	return r;
}
