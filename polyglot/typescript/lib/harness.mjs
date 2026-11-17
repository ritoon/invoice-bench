// Timing and survey harness, a port of bench/ecosystem/common/harness.go to
// the rules of ../PROTOCOL.md: check the result once before timing (that call
// is first_call_ns), "unsupported" when the library cannot do the operation,
// "fail" when it throws or returns a wrong result, then N found by steps
// 1, 2, 5, 10, 20, 50… until one series lasts BENCHTIME seconds, then COUNT
// timed series of N calls on the monotonic clock. One thread, no forced GC.
import { readFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const here = dirname(fileURLToPath(import.meta.url));

export const LANG = 'typescript';
export const COUNT = Number(process.env.COUNT || 3);
export const BENCHTIME = Number(process.env.BENCHTIME || 1);

// Several libraries log on every call (fin.cx prints four lines per
// extraction, stackforge warns about the PDF/A output intent). Logging is not
// the operation being measured and would time the terminal, so the console is
// silenced in the measuring processes; results go straight to stdout.
const out = process.stdout.write.bind(process.stdout);
export function silenceConsole() {
	for (const k of ['log', 'info', 'warn', 'debug', 'error', 'trace', 'dir', 'table']) console[k] = () => {};
}
export function emit(obj) {
	out(JSON.stringify(obj) + '\n');
}

// Read from node_modules directly: some packages do not export package.json.
export function pkgVersion(pkg) {
	return JSON.parse(readFileSync(join(here, '..', 'node_modules', pkg, 'package.json'), 'utf8')).version;
}

const isThenable = (v) => v != null && typeof v.then === 'function';
const now = process.hrtime.bigint;
export let sink; // keeps results alive so that no call is optimised away

export function errMessage(e) {
	if (e == null) return String(e);
	const m = e instanceof Error ? `${e.name === 'Error' ? '' : e.name + ': '}${e.message}` : String(e);
	return m.replace(/\s+/g, ' ').trim().slice(0, 400);
}

export class Unsupported extends Error {}
export class Wrong extends Error {}

// series runs n calls and returns the elapsed nanoseconds.
async function series(fn, n, isAsync) {
	let r;
	const t0 = now();
	if (isAsync) for (let i = 0; i < n; i++) r = await fn();
	else for (let i = 0; i < n; i++) r = fn();
	const dt = now() - t0;
	sink = r;
	return Number(dt);
}

const STEPS = [1, 2, 5];
function* steps() {
	for (let m = 1; ; m *= 10) for (const s of STEPS) yield s * m;
}

// measure times fn. check validates the first result (throws Wrong).
export async function measure(fn, check) {
	let first, r0;
	const t0 = now();
	let res = fn();
	const isAsync = isThenable(res);
	if (isAsync) res = await res;
	first = Number(now() - t0);
	r0 = res;
	if (check) await check(r0);
	sink = r0;

	const target = BENCHTIME * 1e9;
	let N = 1;
	for (const n of steps()) {
		const dt = await series(fn, n, isAsync);
		N = n;
		if (dt >= target) break;
	}
	const ns = [];
	for (let i = 0; i < COUNT; i++) {
		const dt = await series(fn, N, isAsync);
		ns.push(Math.round((dt / N) * 10) / 10);
	}
	return { ns_per_op: ns, iters: N, first_call_ns: first };
}

// record builds one bench_typescript.jsonl line.
export function record(adapter, bench, fields) {
	return {
		lang: LANG,
		lib: adapter.lib,
		pkg: adapter.pkg,
		version: pkgVersion(adapter.pkg),
		bench,
		status: fields.status,
		ns_per_op: fields.ns_per_op ?? [],
		iters: fields.iters ?? 0,
		first_call_ns: fields.first_call_ns ?? null,
		note: fields.note ?? '',
	};
}
