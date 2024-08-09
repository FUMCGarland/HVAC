import { durationMult, hvaccontroller, genRequest } from '$lib/hvac';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`, genRequest());
	const item = await res.json();

	const sched = await fetch(`${hvaccontroller}/api/v1/schedule`, genRequest());
	const s = await sched.json();
	item.Schedule = s.List;
	item.Schedule.forEach((s) => {
		s.RunTime = s.RunTime / durationMult;
	});
	item.Schedule.sort((a, b) => {
		Number(a.ID) - Number(b.ID);
	});
	return item;
}
