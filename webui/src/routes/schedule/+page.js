import { durationMult, hvaccontroller } from '$lib/hvac.js';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`);
	const item = await res.json();

	const sched = await fetch(`${hvaccontroller}/api/v1/schedule`);
	const s = await sched.json();
	item.Schedule = s.List;

	item.Schedule.forEach((s) => {
		s.RunTime = s.RunTime / durationMult;
	});
	return item;
}
