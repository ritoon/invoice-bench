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
			'new SchematronRunner(loadCiiSchematron()) une fois (au premier appel), puis runner.validate(xml) : schematron CEN EN 16931 CII 1.3.16 seul (pas de XSD, pas de XRechnung) ; le résultat ne distingue pas les sévérités, toute assertion déclenchée compte',
	},
	unsupported: {
		extract: 'pas de lecture de PDF',
		parse: 'pas de modèle de facture',
		write: 'pas d’écriture',
		embed: 'pas d’injection dans un PDF',
		check: 'pas de lecture de PDF',
	},
	validate(s) {
		runner ??= new SchematronRunner(loadCiiSchematron());
		const r = runner.validate(s.text);
		return { errors: r.fired.length, rules: r.fired.map((f) => f.assertId) };
	},
};
