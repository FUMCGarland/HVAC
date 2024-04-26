import { hvaccontroller } from '$lib/hvac.js';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/schedule`);
	const data = await res.json();
	const s = data.List.filter((sz) => sz.ID == params.id)[0];

	const system = await fetch(`${hvaccontroller}/api/v1/system`);
	s.System = await system.json();
	return s;
}
