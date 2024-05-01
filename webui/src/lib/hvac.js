import { invalidateAll } from '$app/navigation';
import { toast } from '@zerodevx/svelte-toast';

export const hvaccontroller = `http://192.168.12.5:8080`;
export const durationMult = 60000000000;

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
	if (minutes > 600 || minutes < 5) minutes = 60; // cap runs at 10 hours
	const goduration = minutes * 120000000000;

	const cmd = `{ "TargetState": true, "RunTime": ${goduration}, "Source": "${source}" }`;
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
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

export async function blowerStop(id, source = 'manual') {
	const cmd = `{ "TargetState": false, "RunTime": 0, "Source": "${source}" }`;
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
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
	if (minutes > 600 || minutes < 5) minutes = 60; // cap runs at 10 hours
	const goduration = minutes * 120000000000;

	const cmd = `{ "TargetState": true, "RunTime": ${goduration}, "Source": "${source}" }`;
	const request = {
		method: 'PUT',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
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
			'Content-Type': 'application/json'
		},
		body: cmd
	};

	const response = await fetch(`${hvaccontroller}/api/v1/zone/${id}/targets`, request);
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
			'Content-Type': 'application/json'
		},
		body: cmd
	};
	console.log(cmd);

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
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(cmd)
	};
	console.log(cmd);

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
