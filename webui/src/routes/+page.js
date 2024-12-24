import { hvaccontroller, genRequest } from '$lib/hvac';

export async function load({ fetch, params }) {
	const res = await fetch(`${hvaccontroller}/api/v1/system`, genRequest());
	const item = await res.json();

	item.Rooms.forEach((r) => {
		if (r.LastUpdate == '0001-01-01T00:00:00Z') {
			r.Temperature = 0;
		} else {
			const lastupdate = Date.parse(r.LastUpdate);
			const fourHoursAgo = new Date();
			fourHoursAgo.setHours(fourHoursAgo.getHours() - 4);
			if (lastupdate < fourHoursAgo) { // older than 4 hours
				r.Temperature = 0;
			}
		}
		r.Temperature = Math.round(r.Temperature);
		r.Targets = roomZoneTargets(item, r);
	});

	return item;
}

function roomZoneTargets(data, room) {
	const d = data.Zones.filter((z) => {
		if (data.SystemMode == 0) return z.ID == room.HeatZone;
		return z.ID == room.CoolZone;
	});
	const rz = d[0];
	if (data.SystemMode == 1) {
		if (room.Occupied) {
			return { Min: rz.Targets.CoolingOccupiedTemp - 3, Max: rz.Targets.CoolingOccupiedTemp + 3 };
		} else {
			return {
				Min: rz.Targets.CoolingUnoccupiedTemp - 3,
				Max: rz.Targets.CoolingUnoccupiedTemp + 3
			};
		}
	} else {
		if (room.Occupied) {
			return { Min: rz.Targets.HeatingOccupiedTemp - 3, Max: rz.Targets.HeatingOccupiedTemp + 3 };
		} else {
			return {
				Min: rz.Targets.HeatingUnoccupiedTemp - 3,
				Max: rz.Targets.HeatingUnoccupiedTemp + 3
			};
		}
	}
}
