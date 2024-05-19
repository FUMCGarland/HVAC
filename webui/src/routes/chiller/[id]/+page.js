import { durationMult, hvaccontroller, genRequest } from '$lib/hvac';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`, genRequest());
	const item = await res.json();
	const chiller = item.Chillers.find((chiller) => chiller.ID == params.id);
	chiller.Runtime = Math.floor(chiller.Runtime / durationMult);

	console.log(chiller);
	return chiller;
}
