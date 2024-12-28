import { hvaccontroller, genRequest } from '$lib/hvac';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/occupancy/internal`, genRequest());
	const item = await res.json();
	console.log(item);

	return { data: item };
}
