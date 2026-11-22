// en16931: runs the CEF EN 16931 Schematron artefacts (1.3.16, vendored
// unmodified) with fontoxpath-exact-decimal on slimdom. Validation only.
import { SchematronRunner, loadCiiSchematron } from 'en16931';

// Loading and compiling the Schematron is the documented one-time setup
// (`new SchematronRunner(loadCiiSchematron())`); it is done lazily on the first
// call, so it lands in first_call_ns and not in the timed series.
let runner;

export default {
	lib: 'en16931',
	pkg: 'en16931',
	notes: {
		validate:
			'new SchematronRunner(loadCiiSchematron()) once (on the first call), then runner.validate(xml): CEN EN 16931 CII 1.3.16 schematron alone (no XSD, no XRechnung); the result does not distinguish severities, every fired assertion counts',
	},
	unsupported: {
		extract: 'no PDF reading',
		parse: 'no invoice model',
		write: 'no writing',
		embed: 'no embedding into a PDF',
		check: 'no PDF reading',
	},
	validate(s) {
		runner ??= new SchematronRunner(loadCiiSchematron());
		const r = runner.validate(s.text);
		return { errors: r.fired.length, rules: r.fired.map((f) => f.assertId) };
	},
};
