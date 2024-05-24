import { invalidateAll } from '$app/navigation';
import { toast } from '@zerodevx/svelte-toast';
import { goto } from '$app/navigation';

export const hvaccontroller = import.meta.env.VITE_API_URL;
export const durationMult = 60000000000; // TODO Josh wants hours not minutes

export function genRequest() {
	const jwtstring = localStorage.getItem('jwt');
	if (!jwtstring) {
		goto('/login');
		return;
	}

	const request = {
		mode: 'cors',
		credentials: 'include',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + jwtstring
		}
	};
	return request;
}

// TODO: these are inconsistent, pass in object{} and JSON.stringify() in the body:
// as in the putSchedule
export async function setSystemControlMode(m) {
	const cmd = `{ "ControlMode": ${m} }`;
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + localStorage.getItem('jwt'),
			'Content-Type': 'application/json'
		},
		body: cmd
	};

	const response = await fetch(`${hvaccontroller}/api/v1/system/control`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function setSystemMode(m) {
	const cmd = `{ "Mode": ${m} }`;
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + localStorage.getItem('jwt'),
			'Content-Type': 'application/json'
		},
		body: cmd
	};

	const response = await fetch(`${hvaccontroller}/api/v1/system/mode`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

// these can be consolidated into one with small wrappers, but that's work for later
export async function blowerStart(id, minutes = 60, source = 'manual') {
	if (minutes > 600 || minutes < 30) {
		toast.push('Blower runtime out-of-range (min 30, max 600)');
		return;
	}
	const goduration = minutes * durationMult;

	const cmd = `{ "TargetState": true, "RunTime": ${goduration}, "Source": "${source}" }`;
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + localStorage.getItem('jwt'),
			'Content-Type': 'application/json'
		},
		body: cmd
	};

	const response = await fetch(`${hvaccontroller}/api/v1/blower/${id}/target`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log('server returned ', response.status, payload.error);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function blowerStop(id, source = 'manual') {
	const cmd = `{ "TargetState": false, "RunTime": 0, "Source": "${source}" }`;
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + localStorage.getItem('jwt'),
			'Content-Type': 'application/json'
		},
		body: cmd
	};

	const response = await fetch(`${hvaccontroller}/api/v1/blower/${id}/target`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function pumpStart(id, minutes = 60, source = 'manual') {
	if (minutes > 600 || minutes < 30) {
		toast.push('Pump runtime out-of-range (min 30, max 600)');
		return;
	}
	const goduration = minutes * durationMult;

	const cmd = `{ "TargetState": true, "RunTime": ${goduration}, "Source": "${source}" }`;
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + localStorage.getItem('jwt'),
			'Content-Type': 'application/json'
		},
		body: cmd
	};

	const response = await fetch(`${hvaccontroller}/api/v1/pump/${id}/target`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function pumpStop(id, source = 'manual') {
	const cmd = `{ "TargetState": false, "RunTime": 0, "Source": "${source}" }`;
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + localStorage.getItem('jwt'),
			'Content-Type': 'application/json'
		},
		body: cmd
	};

	const response = await fetch(`${hvaccontroller}/api/v1/pump/${id}/target`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function updateZoneTargets(id, cmd) {
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + localStorage.getItem('jwt'),
			'Content-Type': 'application/json'
		},
		body: cmd
	};

	const response = await fetch(`${hvaccontroller}/api/v1/zone/${id}/temps`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log(payload);
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function postSchedule(cmd) {
	cmd.Runtime = cmd.Runtime * durationMult;

	const request = {
		method: 'POST',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + localStorage.getItem('jwt'),
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(cmd)
	};

	const response = await fetch(`${hvaccontroller}/api/v1/schedule`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log(payload);
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function deleteSchedule(id) {
	const request = {
		method: 'DELETE',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin'
	};

	const response = await fetch(`${hvaccontroller}/api/v1/sched/${id}`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log(payload);
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function putSchedule(cmd) {
	cmd.Runtime = cmd.Runtime * durationMult;

	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + localStorage.getItem('jwt'),
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(cmd)
	};

	const response = await fetch(`${hvaccontroller}/api/v1/sched/${cmd.ID}`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log(payload);
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function zoneStart(id, minutes = 60, source = 'manual') {
	if (minutes > 600 || minutes < 30) {
		toast.push('runtime out-of-range (min 30, max 600)');
		return;
	}
	const goduration = minutes * durationMult;

	const cmd = `{ "TargetState": true, "RunTime": ${goduration}, "Source": "${source}" }`;
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + localStorage.getItem('jwt'),
			'Content-Type': 'application/json'
		},
		body: cmd
	};

	const response = await fetch(`${hvaccontroller}/api/v1/zone/${id}/target`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function zoneStop(id, source = 'manual') {
	const cmd = `{ "TargetState": false, "RunTime": 0, "Source": "${source}" }`;
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + localStorage.getItem('jwt'),
			'Content-Type': 'application/json'
		},
		body: cmd
	};

	const response = await fetch(`${hvaccontroller}/api/v1/zone/${id}/target`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function chillerStop(id, source = 'manual') {
	const cmd = `{ "TargetState": false, "RunTime": 0, "Source": "${source}" }`;
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			Authorization: 'Bearer ' + localStorage.getItem('jwt'),
			'Content-Type': 'application/json'
		},
		body: cmd
	};

	const response = await fetch(`${hvaccontroller}/api/v1/chiller/${id}/target`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}
