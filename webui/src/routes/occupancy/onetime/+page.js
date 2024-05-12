import { hvaccontroller } from '$lib/hvac.js';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/occupancy`);
	const item = await res.json();
	return item;
}
