import { durationMult, hvaccontroller, genRequest } from '$lib/hvac';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`, genRequest());
	const item = await res.json();
	const pump = item.Pumps.find((pump) => pump.ID == params.id);
	pump.Runtime = Math.floor(pump.Runtime / durationMult);

	console.log(pump);
	return pump;
}
