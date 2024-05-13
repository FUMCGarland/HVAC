export function lowestFreeID(arr) {
	if (arr.len == 0) {
		return 0;
	}

	const ids = arr.map((x) => x.ID).sort();
	var i = 1;
	ids.forEach((x) => {
		if (x != i) return i;
		i++;
	});
	return i;
}
