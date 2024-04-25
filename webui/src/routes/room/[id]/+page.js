import { hvaccontroller } from '$lib/hvac.js';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`);
	const data = await res.json();

	const room = data.Rooms.filter((r) => r.ID == params.id);
    if (room.length == 1) {
        return room[0];
    } else {
        return {
            ID: params.id,
            Name: "404 Room not found",
            Zone: 0,
            Temperature: 0
        }
    }
}
