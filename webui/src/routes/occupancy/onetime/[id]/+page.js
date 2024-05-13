import { hvaccontroller } from '$lib/hvac.js';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/occupancy`);
	const all = await res.json();

	const filt = all.OneTime.filter((r) => params.id == r.ID);
	const item = filt[0];

	const sys = await fetch(`${hvaccontroller}/api/v1/system`);
	const system = await sys.json();
	item.SystemRooms = system.Rooms;

	return item;
}
