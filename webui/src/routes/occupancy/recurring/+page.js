import { hvaccontroller } from '$lib/hvac.js';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/occupancy`);
	const item = await res.json();
	if (!item.Recurring) item.Recurring = new Array();

	const sys = await fetch(`${hvaccontroller}/api/v1/system`);
	const system = await sys.json();
	item.Rooms = system.Rooms;

	return item;
}
