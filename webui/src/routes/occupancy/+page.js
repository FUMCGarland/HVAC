import { hvaccontroller, genRequest } from '$lib/hvac';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/occupancy`, genRequest());
	const item = await res.json();
	if (!item.Recurring) item.Recurring = new Array();
	if (!item.OneTime) item.OneTime = new Array();

	return item;
}
