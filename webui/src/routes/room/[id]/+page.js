import { hvaccontroller, genRequest } from '$lib/hvac';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`, genRequest());
	const data = await res.json();

	const room = data.Rooms.filter((r) => r.ID == params.id);
	if (room.length == 1) {
		room[0].LastUpdate = new Date(Date.parse(room[0].LastUpdate)).toLocaleString();
		return room[0];
	}
	console.log(room);
	return {
		ID: params.id,
		Name: '404 Room not found',
		Zone: 0,
		Temperature: 0
	};
}
