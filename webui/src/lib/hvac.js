import { invalidateAll } from '$app/navigation';
import { toast } from '@zerodevx/svelte-toast';

export const hvaccontroller = `http://192.168.12.5:8080`;

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
export async function blowerStart(id, duration = '60', source = 'manual') {
	if (duration > 600) duration = 60; // cap runs at 10 hours
	const d = duration * 60;

	const cmd = `{ "TargetState": true, "RunTime": ${d}, "Source": "${source}" }`;
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

export async function pumpStart(id, duration = '60', source = 'manual') {
	if (duration > 600) duration = 60; // cap runs at 10 hours
	const d = duration * 60;

	const cmd = `{ "TargetState": true, "RunTime": ${d}, "Source": "${source}" }`;
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
