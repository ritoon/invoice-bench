// Reference validator: SaxonC-HE 13.1.1 (Saxonica, MPL-2.0), Saxon compiled
// to a native library by GraalVM (no JVM), running the official stylesheets,
// as polyglot/python/adapters/saxonche_schematron.py does with saxonche
// (the Python binding of the same SaxonC):
//
// - CEN EN 16931 1.3.17 (ConnectingEurope/eInvoicing-EN16931):
//   xslt/EN16931-CII-validation.xslt for a CII invoice,
//   xslt/EN16931-UBL-validation.xslt for a UBL invoice or credit note;
// - KoSIT XRechnung 3.0.2 schematron 2.6.0 (itplr-kosit/xrechnung-schematron,
//   schematron/cii/XRechnung-CII-validation.xsl or
//   schematron/ubl/XRechnung-UBL-validation.xsl) in addition when BT-24 names
//   XRechnung, as the scenarios of validator-configuration-xrechnung.
//
// The stylesheets of a syntax are compiled on its first call (in
// first_call_ns for CII, the syntax of the bench), then reused. Counted as the
// KoSIT validator does: every assertion flagged fatal or error, with the
// custom levels of the scenarios of validator-configuration-xrechnung
// v2026-08-31 for XRechnung (see level()). The SVRL report is read with
// pugixml inside the timed call, as the Python adapter reads it with lxml. No
// XSD (not part of the CEN artefacts). The CII path is the one of the bench;
// the UBL path serves the conformance cases (../../cases).
#include <saxonc/SaxonApiException.h>
#include <saxonc/SaxonProcessor.h>
#include <saxonc/XdmNode.h>
#include <saxonc/Xslt30Processor.h>
#include <saxonc/XsltExecutable.h>

#include <pugixml.hpp>

#include <cstring>
#include <map>
#include <memory>
#include <stdexcept>
#include <string>
#include <vector>

#include "harness.hpp"

namespace {

// Stylesheets by syntax: CEN, then KoSIT XRechnung.
const char *const REL[2][2] = {
	{"cen-1.3.17/xslt/EN16931-CII-validation.xslt", "xrechnung-2.6.0/schematron/cii/XRechnung-CII-validation.xsl"},
	{"cen-ubl-1.3.17/xslt/EN16931-UBL-validation.xslt",
	 "xrechnung-2.6.0/schematron/ubl/XRechnung-UBL-validation.xsl"},
};

struct Sheets {
	XsltExecutable *cen = nullptr, *xr = nullptr;
	std::string error;
};

struct State {
	SaxonProcessor *proc = nullptr;
	Xslt30Processor *xslt = nullptr;
	Sheets syn[2]; // 0 CII, 1 UBL
};
State S;

// The level of a rule under the scenario of validator-configuration-xrechnung
// v2026-08-31 that BT-24 selects (CIUS, extension or CVD), its flag otherwise:
// the customLevel elements of scenarios.xml, where they change whether the
// rule is an error.
std::string level(const std::string &rule, const std::string &flag, const std::string &spec,
		  const std::string &root) {
	if (spec.find("xrechnung") == std::string::npos)
		return flag;
	bool ext = spec.find(":extension:") != std::string::npos, cvd = spec.find(":cvd_") != std::string::npos;
	auto in = [&](std::initializer_list<const char *> l) {
		for (const char *r : l)
			if (rule == r)
				return true;
		return false;
	};
	if (in({"BR-CL-21", "BR-CL-23", "CII-SR-475", "CII-SR-476"}))
		return "warning";
	if (in({"CII-SR-452", "CII-SR-453", "CII-SR-454", "CII-SR-465", "CII-SR-466"}))
		return "error";
	if (rule == "BR-CL-13")
		return cvd ? "information" : flag;
	if (in({"BR-CL-10", "BR-CL-11", "BR-CL-24", "BR-CL-25", "BR-CL-26"}))
		return ext ? "information" : flag;
	if (in({"BR-CO-16", "UBL-CR-470"})) // UBL extension scenario only
		return ext && root != "CrossIndustryInvoice" ? "information" : flag;
	if (rule == "UBL-CR-646")
		return ext ? "information" : cvd || root != "CreditNote" ? "error" : flag;
	return flag;
}

std::string saxon_message(SaxonApiException &e) {
	const char *m = e.getMessageWithErrorCode();
	if (!m || !*m)
		m = e.getMessage();
	return m ? m : "SaxonApiException";
}

Sheets &compile(int syn) {
	Sheets &sh = S.syn[syn];
	if (sh.cen)
		return sh;
	if (!sh.error.empty())
		throw std::runtime_error(sh.error);
	try {
		for (const char *rel : REL[syn]) {
			std::string p = fx::artefact(rel);
			if (FILE *f = std::fopen(p.c_str(), "rb"))
				std::fclose(f);
			else
				throw std::runtime_error(p + " missing: run ./run.sh setup");
		}
		if (!S.proc) {
			S.proc = new SaxonProcessor(false);
			S.xslt = S.proc->newXslt30Processor();
			// Base URI of the static context: SaxonC wants an absolute path.
			S.xslt->setcwd(fx_artefacts());
		}
		sh.cen = S.xslt->compileFromFile(fx::artefact(REL[syn][0]).c_str());
		sh.xr = S.xslt->compileFromFile(fx::artefact(REL[syn][1]).c_str());
		if (!sh.cen || !sh.xr)
			throw std::runtime_error("compileFromFile returned null");
	} catch (SaxonApiException &e) {
		sh.cen = sh.xr = nullptr;
		sh.error = "compile: " + saxon_message(e);
		throw std::runtime_error(sh.error);
	} catch (const std::exception &e) {
		sh.cen = sh.xr = nullptr;
		sh.error = e.what();
		throw;
	}
	return sh;
}

// Bytes read to find the root element and BT-24: the official samples open
// with long licence comments.
const size_t HEAD = 65536;

std::string_view local(std::string_view qname) {
	size_t c = qname.find(':');
	return c == std::string_view::npos ? qname : qname.substr(c + 1);
}

// Local name of the root element ("" if none in the first HEAD bytes).
std::string root_name(fx::Bytes xml) {
	xml = xml.substr(0, HEAD);
	for (size_t p = xml.find('<'); p != fx::Bytes::npos; p = xml.find('<', p + 1)) {
		if (p + 1 >= xml.size())
			break;
		char c = xml[p + 1];
		if (c == '?' || c == '!') {
			if (xml.compare(p, 4, "<!--") == 0) {
				size_t e = xml.find("-->", p + 4);
				if (e == fx::Bytes::npos)
					break;
				p = e;
			}
			continue;
		}
		size_t e = xml.find_first_of(" \t\r\n/>", p + 1);
		return std::string(local(xml.substr(p + 1, (e == fx::Bytes::npos ? xml.size() : e) - p - 1)));
	}
	return "";
}

// Text of the first element named `name` (any prefix) after `from`.
std::string element_text(fx::Bytes xml, size_t from, std::string_view name) {
	for (size_t p = xml.find('<', from); p != fx::Bytes::npos; p = xml.find('<', p + 1)) {
		size_t gt = xml.find('>', p);
		if (gt == fx::Bytes::npos)
			return "";
		fx::Bytes tag = xml.substr(p + 1, gt - p - 1);
		if (local(tag) == name) {
			size_t lt = xml.find('<', gt);
			fx::Bytes t = xml.substr(gt + 1, lt == fx::Bytes::npos ? fx::Bytes::npos : lt - gt - 1);
			size_t a = t.find_first_not_of(" \t\r\n"), b = t.find_last_not_of(" \t\r\n");
			return a == fx::Bytes::npos ? "" : std::string(t.substr(a, b - a + 1));
		}
	}
	return "";
}

// BT-24, the specification identifier (same test as the Python adapter's
// regex): CII GuidelineSpecifiedDocumentContextParameter/ID, UBL
// CustomizationID.
std::string spec_id(fx::Bytes xml, const std::string &root) {
	xml = xml.substr(0, HEAD);
	if (root == "CrossIndustryInvoice") {
		static const char tag[] = "GuidelineSpecifiedDocumentContextParameter>";
		size_t at = xml.find(tag);
		return at == fx::Bytes::npos ? "" : element_text(xml, at + sizeof tag - 1, "ID");
	}
	size_t at = xml.find("CustomizationID>");
	if (at == fx::Bytes::npos)
		return "";
	size_t lt = xml.rfind('<', at);
	return lt == fx::Bytes::npos ? "" : element_text(xml, lt, "CustomizationID");
}

// Rules of the SVRL failed assertions and successful reports flagged fatal or
// error (after the levels of the KoSIT scenario).
void errors(const char *svrl, const std::string &spec, const std::string &root, std::vector<std::string> &out) {
	pugi::xml_document d;
	if (!d.load_string(svrl))
		throw std::runtime_error("unreadable SVRL report");
	struct Walker : pugi::xml_tree_walker {
		const std::string *spec, *root;
		std::vector<std::string> *out;
		bool for_each(pugi::xml_node &n) override {
			if (n.type() != pugi::node_element)
				return true;
			std::string_view l = local(n.name());
			if (l != "failed-assert" && l != "successful-report")
				return true;
			std::string rule = n.attribute("id").as_string("?");
			std::string flag = n.attribute("flag").as_string("error");
			for (auto &c : flag)
				c = (char)std::tolower((unsigned char)c);
			flag = level(rule, flag, *spec, *root);
			if (flag == "fatal" || flag == "error")
				out->push_back(rule);
			return true;
		}
	} w;
	w.spec = &spec;
	w.root = &root;
	w.out = &out;
	d.traverse(w);
}

struct Report {
	std::vector<std::string> rules;
};

Report validate(fx::Bytes xml) {
	if (xml.size() >= 3 && std::memcmp(xml.data(), "\xEF\xBB\xBF", 3) == 0)
		xml.remove_prefix(3);
	std::string root = root_name(xml);
	Sheets &sh = compile(root == "Invoice" || root == "CreditNote" ? 1 : 0);
	std::string spec = spec_id(xml, root);
	bool xr = spec.find("xrechnung") != std::string::npos;
	std::string text(xml);
	Report rep;
	XdmNode *node = nullptr;
	const char *svrl = nullptr;
	try {
		node = S.proc->parseXmlFromString(text.c_str(), "UTF-8");
		svrl = sh.cen->transformToString(node);
		errors(svrl ? svrl : "", spec, root, rep.rules);
		SaxonProcessor::deleteString(svrl);
		svrl = nullptr;
		if (xr) {
			svrl = sh.xr->transformToString(node);
			errors(svrl ? svrl : "", spec, root, rep.rules);
			SaxonProcessor::deleteString(svrl);
			svrl = nullptr;
		}
	} catch (SaxonApiException &e) {
		SaxonProcessor::deleteString(svrl);
		delete node;
		throw std::runtime_error(saxon_message(e));
	} catch (...) {
		SaxonProcessor::deleteString(svrl);
		delete node;
		throw;
	}
	delete node;
	return rep;
}

} // namespace

int main(int argc, char **argv) {
	fx_lib lib{"saxonc-he", "SaxonC-HE", "13.1.1"};
	return fx::run(lib, argc, argv, [](fx_harness *h) {
		auto va = [](fx::Bytes x) { return validate(x); };
		auto su = [](const Report &r) { return fx::Findings{r.rules.size(), r.rules}; };
		fx_unsupported(h, FX_EXTRACT, "validator only");
		fx_unsupported(h, FX_PARSE, "validator only");
		fx::validate(h,
			     "official CEN EN 16931 1.3.17 XSLT, CII or UBL by root (+ KoSIT XRechnung 3.0.2 schematron "
			     "2.6.0 if BT-24 is XRechnung) via SaxonC-HE 13.1.1, compiled on the first call; fatal/error "
			     "errors with the levels of the KoSIT scenarios",
			     va, su);
		fx_unsupported(h, FX_WRITE, "validator only");
		fx_unsupported(h, FX_EMBED, "validator only");
		fx_unsupported(h, FX_CHECK_PDF, "validator only");
		// Same thread as the GraalVM isolate that SaxonProcessor created.
		for (Sheets &sh : S.syn) {
			delete sh.cen;
			delete sh.xr;
		}
		delete S.xslt;
		delete S.proc;
		if (S.proc)
			SaxonProcessor::release();
	});
}
