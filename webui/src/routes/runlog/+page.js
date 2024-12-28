import { hvaccontroller, genRequest } from '$lib/hvac';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/runlog`, genRequest());
	const item = await res.text();
	return { data: item };
}
