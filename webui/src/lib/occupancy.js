import { invalidateAll } from '$app/navigation';
import { toast } from '@zerodevx/svelte-toast';
import { hvaccontroller, durationMult } from './hvac.js';

export async function postRecurringOccupancy(cmd) {
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
		body: JSON.stringify(cmd)
	};

	const response = await fetch(`${hvaccontroller}/api/v1/occupancy/recurring`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log(payload);
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function deleteRecurringOccupancy(id) {
	const request = {
		method: 'DELETE',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin'
	};

	const response = await fetch(`${hvaccontroller}/api/v1/occupancy/recurring/${id}`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log(payload);
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function putRecurringOccupancy(cmd) {
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

	const response = await fetch(`${hvaccontroller}/api/v1/occupancy/recurring/${cmd.ID}`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log(payload);
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function postOneTimeOccupancy(cmd) {
	const request = {
		method: 'POST',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(cmd)
	};

	const response = await fetch(`${hvaccontroller}/api/v1/occupancy/onetime`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log(payload);
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function deleteOneTimeOccupancy(id) {
	const request = {
		method: 'DELETE',
		mode: 'cors',
		credentials: 'include',
		redirect: 'manual',
		referrerPolicy: 'origin'
	};

	const response = await fetch(`${hvaccontroller}/api/v1/occupancy/onetime/${id}`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log(payload);
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}

export async function putOneTimeOccupancy(cmd) {
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

	const response = await fetch(`${hvaccontroller}/api/v1/occupancy/onetime/${cmd.ID}`, request);
	const payload = await response.json();

	if (response.status != 200) {
		console.log(payload);
		console.log('server returned ', response.status);
		toast.push('Server Responded with: ' + response.status + ': ' + payload.error);
		return;
	}
	invalidateAll();
}
