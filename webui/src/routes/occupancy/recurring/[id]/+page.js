import { hvaccontroller, genRequest } from '$lib/hvac';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/occupancy`, genRequest());
	const all = await res.json();

	const filt = all.Recurring.filter((r) => params.id == r.ID);
	const item = filt[0];

	const sys = await fetch(`${hvaccontroller}/api/v1/system`, genRequest());
	const system = await sys.json();
	item.SystemRooms = system.Rooms;

	return item;
}
