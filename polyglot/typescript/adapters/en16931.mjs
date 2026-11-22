// en16931: runs the CEF EN 16931 Schematron artefacts (1.3.16, vendored
// unmodified) with fontoxpath-exact-decimal on slimdom. Validation only.
import { SchematronRunner, loadCiiSchematron, loadUblSchematron } from 'en16931';
import { isUbl } from '../lib/corpus.mjs';

// Loading and compiling the Schematron is the documented one-time setup
// (`new SchematronRunner(loadCiiSchematron())`); it is done lazily on the first
// call, so it lands in first_call_ns and not in the timed series.
let runner, ublRunner;

export default {
	lib: 'en16931',
	pkg: 'en16931',
	notes: {
		validate:
			'new SchematronRunner(loadCiiSchematron()) once (on the first call), then runner.validate(xml): CEN EN 16931 CII 1.3.16 schematron alone (no XSD, no XRechnung), the UBL one (loadUblSchematron) for a UBL document; the result does not distinguish severities, every fired assertion counts',
	},
	unsupported: {
		extract: 'no PDF reading',
		parse: 'no invoice model',
		write: 'no writing',
		embed: 'no embedding into a PDF',
		check: 'no PDF reading',
	},
	validate(s) {
		let r;
		if (isUbl(s.text)) {
			ublRunner ??= new SchematronRunner(loadUblSchematron());
			r = ublRunner.validate(s.text);
		} else {
			runner ??= new SchematronRunner(loadCiiSchematron());
			r = runner.validate(s.text);
		}
		return { errors: r.fired.length, rules: r.fired.map((f) => f.assertId) };
	},
};
