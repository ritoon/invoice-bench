// Floor: pugixml 1.16 (MIT), the fast in-memory DOM of the C++ world, with no
// knowledge of Factur-X and no validation. Parse = xml_document::load_buffer
// (parse_default), Write = xml_document::save (format_default) into a string.
#include <pugixml.hpp>

#include <memory>
#include <stdexcept>
#include <string>

#include "harness.hpp"

namespace {

using Doc = std::unique_ptr<pugi::xml_document>;

Doc parse(fx::Bytes xml) {
	auto d = std::make_unique<pugi::xml_document>();
	pugi::xml_parse_result r = d->load_buffer(xml.data(), xml.size());
	if (!r)
		throw std::runtime_error(std::string("load_buffer: ") + r.description() + " at offset " +
					 std::to_string(r.offset));
	return d;
}

struct StringWriter : pugi::xml_writer {
	std::string out;
	void write(const void *data, size_t size) override { out.append(static_cast<const char *>(data), size); }
};

std::string write(const Doc &d) {
	StringWriter w;
	w.out.reserve(32 * 1024);
	d->save(w);
	return std::move(w.out);
}

} // namespace

int main(int argc, char **argv) {
	fx_lib lib{"pugixml", "pugixml", "1.16"};
	static_assert(PUGIXML_VERSION == 1160, "adapter written for pugixml 1.16");
	return fx::run(lib, argc, argv, [](fx_harness *h) {
		auto pa = [](fx::Bytes x) { return parse(x); };
		auto wr = [](const Doc &d) { return write(d); };
		fx_unsupported(h, FX_EXTRACT, "XML library only");
		fx::parse(h, "floor: pugixml DOM (load_buffer, parse_default), not an invoice model", pa);
		fx_unsupported(h, FX_VALIDATE, "no validation (neither XSD nor rules)");
		fx::write(h, "floor: serialization of the pugixml DOM (save, format_default), not an invoice model", pa, wr);
		fx_unsupported(h, FX_EMBED, "XML library only");
		fx_unsupported(h, FX_CHECK_PDF, "XML library only");
	});
}
