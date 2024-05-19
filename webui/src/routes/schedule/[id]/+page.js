import { hvaccontroller, durationMult, genRequest } from '$lib/hvac';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/schedule`, genRequest());
	const data = await res.json();
	const s = data.List.filter((sz) => sz.ID == params.id)[0];

	const system = await fetch(`${hvaccontroller}/api/v1/system`, genRequest());
	s.System = await system.json();
	s.RunTime = s.RunTime / durationMult; // work in minutes
	return s;
}
