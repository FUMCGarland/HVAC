export function lowestFreeID(arr) {
	if (arr.len == 0) {
		return 0;
	}

	const ids = arr.map((x) => x.ID).sort();
	let i = 1;
	ids.forEach((x) => {
		if (x != i) return i;
		i++;
	});
	return i;
}

export function getAccessLevel() {
	const raw = localStorage.getItem('jwt');
	if (!raw) return 0;

	const token = JSON.parse(window.atob(raw.split('.')[1]).toString());
	return token.level;
}

export function toLocaleTimestring(timestring) {
	const hm = timestring.split(':');
	let t = new Date();
	t.setUTCHours(Number(hm[0]));
	t.setUTCMinutes(Number(hm[1]));
	t.setUTCSeconds(0);
	return t.toLocaleTimeString();
}

export function toZTimestring(timestring) {
	const hm = timestring.split(':');

	let hour = Number(hm[0]);
	if (hm[2] !== undefined && hm[2].endsWith('PM')) {
		hour = hour + 12;
	}

	let t = new Date();
	t.setHours(hour);
	t.setMinutes(Number(hm[1]));
	t.setSeconds(0);

	return t.getUTCHours() + ':' + hm[1]; // saves writing padding code, but assumes the user isn't doing something stupid
}
