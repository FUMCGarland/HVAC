import { hvaccontroller } from '$lib/hvac.js';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`);
	const item = await res.json();
	const blower = item.Blowers.find((blower) => blower.ID == params.id);
	return blower;
}
