import { hvaccontroller, genRequest } from '$lib/hvac';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`, genRequest());
	const item = await res.json();
	if (!item.Chillers) item.Chillers = new Array();
	return item;
}
