import { hvaccontroller, genRequest } from '$lib/hvac';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/occupancy`, genRequest());
	const item = await res.json();
	if (!item.Recurring) item.Recurring = new Array();

	const sys = await fetch(`${hvaccontroller}/api/v1/system`, genRequest());
	const system = await sys.json();
	item.Rooms = system.Rooms;

	return item;
}
