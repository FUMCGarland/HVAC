import { hvaccontroller } from '$lib/hvac.js';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`);
	const item = await res.json();
	const pump = item.Pumps.find((pump) => pump.ID == params.id);
	return pump;
}
