import { durationMult, hvaccontroller } from '$lib/hvac.js';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`);
	const item = await res.json();
	const blower = item.Blowers.find((blower) => blower.ID == params.id);

	blower.Runtime = Math.floor(blower.Runtime / durationMult);
	blower.FilterTime = Math.floor(blower.FilterTime / durationMult);

	return blower;
}
