<script>
	import {
		Badge,
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Heading,
		A
	} from 'flowbite-svelte';

	export let data;
	let sortBy = { col: 'Name', ascending: true };

	$: sort = (column) => {
		if (sortBy.col == column) {
			sortBy.ascending = !sortBy.ascending;
		} else {
			sortBy.col = column;
			sortBy.ascending = true;
		}

		let sm = sortBy.ascending ? 1 : -1;

		let sort = (a, b) => (a[column] < b[column] ? -1 * sm : a[column] > b[column] ? 1 * sm : 0);

		data.Rooms = data.Rooms.sort(sort);
	};

	function zoneName(zoneID) {
		const z = data.Zones.filter((z) => z.ID == zoneID);
		return z[0].Name;
	}

	data.Rooms.forEach((r) => {
		r.Targets = roomZoneTargets(r);
	});

	function roomZoneTargets(room) {
		const d = data.Zones.filter((z) => {
			return z.ID == room.Zone;
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
</script>

<Heading tag="h2">Rooms</Heading>
<Table>
	<TableHead>
		<TableHeadCell on:click={sort('Name')}>Name</TableHeadCell>
		<TableHeadCell on:click={sort('Occupied')}>Occupied</TableHeadCell>
		<TableHeadCell on:click={sort('Zone')}>Zone</TableHeadCell>
		<TableHeadCell on:click={sort('Temperature')}>Temperature</TableHeadCell>
		<TableHeadCell on:click={sort('Humidity')}>Humidity</TableHeadCell>
		<TableHeadCell on:click={sort('Battery')}>Battery</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Rooms as room}
			<TableBodyRow>
				<TableBodyCell><A href="/room/{room.ID}">{room.Name}</A></TableBodyCell>
				<TableBodyCell>{room.Occupied}</TableBodyCell>
				<TableBodyCell><A href="/zone/{room.Zone}">{zoneName(room.Zone)}</A></TableBodyCell>

				{#if room.Temperature == 0}
					<TableBodyCell>&nbsp;</TableBodyCell>
				{/if}
				{#if room.Temperature != 0 && room.Temperature < room.Targets.Min - 10}
					<TableBodyCell><Badge color="red">{room.Temperature}</Badge></TableBodyCell>
				{/if}
				{#if room.Temperature >= room.Targets.Min - 10 && room.Temperature < room.Targets.Min - 5}
					<TableBodyCell><Badge color="yellow">{room.Temperature}</Badge></TableBodyCell>
				{/if}
				{#if room.Temperature >= room.Targets.Min - 5 && room.Temperature < room.Targets.Max + 5}
					<TableBodyCell><Badge color="green">{room.Temperature}</Badge></TableBodyCell>
				{/if}
				{#if room.Temperature >= room.Targets.Max + 5 && room.Temperature < room.Targets.Max + 10}
					<TableBodyCell><Badge color="yellow">{room.Temperature}</Badge></TableBodyCell>
				{/if}
				{#if room.Temperature >= room.Targets.Max + 10}
					<TableBodyCell><Badge color="red">{room.Temperature}</Badge></TableBodyCell>
				{/if}

				{#if room.Humidity < 1}
					<TableBodyCell>&nbsp;</TableBodyCell>
				{/if}
				{#if room.Humidity > 0 && room.Humidity < 20}
					<TableBodyCell><Badge color="red">{room.Humidity}</Badge></TableBodyCell>
				{/if}
				{#if room.Humidity >= 20 && room.Humidty <= 30}
					<TableBodyCell><Badge color="yellow">{room.Humidity}</Badge></TableBodyCell>
				{/if}
				{#if room.Humidity > 30 && room.Humidity <= 60}
					<TableBodyCell><Badge color="green">{room.Humidity}</Badge></TableBodyCell>
				{/if}
				{#if room.Humidity > 60 && room.Humidity <= 80}
					<TableBodyCell><Badge color="yellow">{room.Humidity}</Badge></TableBodyCell>
				{/if}
				{#if room.Humidity > 80}
					<TableBodyCell><Badge color="red">{room.Humidity}</Badge></TableBodyCell>
				{/if}

				{#if room.Battery > 45}
					<TableBodyCell><Badge color="green">{room.Battery}</Badge></TableBodyCell>
				{/if}
				{#if room.Battery <= 45 && room.Battery > 15}
					<TableBodyCell><Badge color="yellow">{room.Battery}</Badge></TableBodyCell>
				{/if}
				{#if room.Battery <= 15}
					<TableBodyCell><Badge color="red">{room.Battery}</Badge></TableBodyCell>
				{/if}
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
