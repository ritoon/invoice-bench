import { pathToFileURL } from 'node:url';

// Library ids accepted by run.sh, bench.mjs and survey.mjs, in campaign order.
export const ADAPTERS = {
	stackforge: './stackforge.mjs',
	'e-invoice-eu': './e-invoice-eu.mjs',
	'node-zugferd': './node-zugferd.mjs',
	stafyniaksacha: './stafyniaksacha.mjs',
	'factur-x-kit': './factur-x-kit.mjs',
	lucaapp: './lucaapp.mjs',
	attestwire: './attestwire.mjs',
	ekwo: './ekwo.mjs',
	conformo: './conformo.mjs',
	kontor: './kontor.mjs',
	summand: './summand.mjs',
	en16931: './en16931.mjs',
	fincx: './fincx.mjs',
};

export async function loadAdapter(id) {
	const path = ADAPTERS[id];
	if (!path) throw new Error(`unknown library "${id}" (known: ${Object.keys(ADAPTERS).join(' ')})`);
	return (await import(path)).default;
}

// node adapters/index.mjs  →  the ids, space-separated (used by run.sh)
if (import.meta.url === pathToFileURL(process.argv[1]).href) console.log(Object.keys(ADAPTERS).join(' '));
