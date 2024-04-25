import { hvaccontroller } from '$lib/hvac.js';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`);
	const data = await res.json();

	const ld = {};

	const loop = data.Loops.filter((l) => l.ID == params.id);
	if (loop.length == 1) {
		ld.Loop = loop[0];
	}

	const pump = data.Pumps.filter((p) => p.Loop == params.id);
	if (pump.length == 1) {
		ld.Pump = pump[0];
	}
	ld.Blowers = data.Blowers.filter((b) => b.HotLoop == params.id || b.ColdLoop == params.id);

	return ld;
}
