import { hvaccontroller } from '$lib/hvac.js';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`);
	const item = await res.json();
	const zone = item.Zones.find((zone) => zone.ID == params.id);

	zone.Rooms = item.Rooms.filter((room) => room.Zone == params.id);
	zone.Blowers = item.Blowers.filter((blower) => blower.Zone == params.id);
	return zone;
}
