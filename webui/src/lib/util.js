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
	const token = JSON.parse(window.atob(raw.split('.')[1]).toString());
	console.log(token);
	return token.level;
}
