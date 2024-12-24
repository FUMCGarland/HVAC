<script>
	import { onMount } from 'svelte';
	import { UsersGroupOutline, CloseCircleOutline } from 'flowbite-svelte-icons';
	import { invalidateAll } from '$app/navigation';
	import {
		Badge,
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Heading,
		A,
		P
	} from 'flowbite-svelte';
	import { setOccupancyManual } from '$lib/hvac';

	export let data;

	// how do we sort the data for display
	const sortBy = { col: 'Name', ascending: true };

	// refresh every 30 seconds
	onMount(() => {
		console.log(data);
		tablesort(sortBy.col);

		const interval = setInterval(async () => {
			await invalidateAll();
			tablesort(sortBy.col);
			tablesort(sortBy.col); // twice to keep the ascending correct
		}, 30000);

		return () => {
			clearInterval(interval);
		};
	});

	$: tablesort = (column) => {
		if (sortBy.col == column) {
			sortBy.ascending = !sortBy.ascending;
		} else {
			sortBy.col = column;
			sortBy.ascending = true;
		}

		let sm = sortBy.ascending ? 1 : -1;

		let sortcallback = (a, b) =>
			a[column] < b[column] ? -1 * sm : a[column] > b[column] ? 1 * sm : 0;
		data.Rooms = data.Rooms.sort(sortcallback);
	};

	function zoneName(zoneID) {
		if (zoneID == 0) return "no zone";
		const z = data.Zones.filter((z) => z.ID == zoneID);
		return z[0].Name;
	}

	function toggleOccupied(room, state) {
		console.log('toggle occupied', room, state);
		setOccupancyManual(room, state);
	}

	function zoneAvgTemp(room) {
		// console.log(room);
		var zoneID = 0;
		if (data.SystemMode == 0) {
			zoneID = room.HeatZone;
		} else {
			zoneID = room.CoolZone;
		}
		const zone = data.Zones.find((z) => z.ID == zoneID);
		return Math.round(zone.AverageTemp);
	}

	function zoneRunning(zone) {
		// if a zone has blowers, it is running if all the blowers in the zone are running
		const blowers = data.Blowers.filter((b) => b.Zone == zone);
		if (blowers.length != 0) {
			const running = blowers.filter((b) => b.Running);
			return blowers.length == running.length && blowers.length != 0;
		}

		// if a zone does not have blowers, it is running if the pump is running
		// zone has loop, loop has pump, pump has running
		const loop = data.Loops.filter((l) => l.RadiantZone == zone);
		if (loop.length != 1) return; // no blowers, no radiant... nothing to run
		const pump = data.Pumps.filter((p) => p.Loop == loop[0].ID);
		return pump[0].Running;
	}
</script>

<Heading tag="h2">Rooms</Heading>
<P
	>Legend:
	<Badge color="green">+/- 3 degrees of target</Badge>
	<Badge color="blue">3-10 degrees cooler</Badge>
	<Badge color="yellow">3-10 degrees warmer</Badge>
	<Badge color="purple">&gt; 10 degrees cooler</Badge>
	<Badge color="red">&gt; 10 degrees warmer</Badge>
</P>
<Table>
	<TableHead>
		<TableHeadCell on:click={tablesort('Name')}>Name</TableHeadCell>
		<TableHeadCell on:click={tablesort('Occupied')}>Occupied</TableHeadCell>
		<TableHeadCell on:click={tablesort('Zone')}>Zone</TableHeadCell>
		<TableHeadCell on:click={tablesort('Temperature')}>Zone Avg</TableHeadCell>
		<TableHeadCell on:click={tablesort('Temperature')}>Room Temp</TableHeadCell>
		<TableHeadCell on:click={tablesort('Battery')}>Battery</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Rooms as room}
			<TableBodyRow>
				<TableBodyCell><A href="/room/{room.ID}">{room.Name}</A></TableBodyCell>
				<TableBodyCell>
					{#if !room.Occupied}
						<A on:click={toggleOccupied(room.ID, true)}><CloseCircleOutline /></A>
					{/if}
					{#if room.Occupied}
						<A on:click={toggleOccupied(room.ID, false)}><UsersGroupOutline /></A>
					{/if}
				</TableBodyCell>
				{#if data.SystemMode == 1}
				<TableBodyCell
					><A href="/zone/{room.CoolZone}">
						{#if zoneRunning(room.CoolZone)}
							<Badge color="green">Running</Badge>
						{/if}
						{zoneName(room.CoolZone)}</A
					></TableBodyCell
				>
				{/if}
				{#if data.SystemMode == 0}
				<TableBodyCell
					><A href="/zone/{room.HeatZone}">
						{#if zoneRunning(room.HeatZone)}
							<Badge color="green">Running</Badge>
						{/if}
						{zoneName(room.HeatZone)}</A
					></TableBodyCell
				>
				{/if}
				<TableBodyCell>{zoneAvgTemp(room)}</TableBodyCell>
				{#if room.Temperature == 0}
					<TableBodyCell>&nbsp;</TableBodyCell>
				{/if}
				{#if room.Temperature != 0 && room.Temperature < room.Targets.Min - 10}
					<TableBodyCell><Badge color="purple">{room.Temperature}</Badge></TableBodyCell>
				{/if}
				{#if room.Temperature >= room.Targets.Min - 10 && room.Temperature < room.Targets.Min}
					<TableBodyCell><Badge color="blue">{room.Temperature}</Badge></TableBodyCell>
				{/if}
				{#if room.Temperature >= room.Targets.Min && room.Temperature < room.Targets.Max}
					<TableBodyCell><Badge color="green">{room.Temperature}</Badge></TableBodyCell>
				{/if}
				{#if room.Temperature >= room.Targets.Max && room.Temperature < room.Targets.Max + 10}
					<TableBodyCell><Badge color="yellow">{room.Temperature}</Badge></TableBodyCell>
				{/if}
				{#if room.Temperature >= room.Targets.Max + 10}
					<TableBodyCell><Badge color="red">{room.Temperature}</Badge></TableBodyCell>
				{/if}

				{#if room.ShellyID != ""}
				{#if room.Battery == 101}
					<TableBodyCell><Badge color="green">Powered</Badge></TableBodyCell>
				{/if}
				{#if room.Battery <= 100 && room.Battery > 45}
					<TableBodyCell><Badge color="green">{room.Battery}</Badge></TableBodyCell>
				{/if}
				{#if room.Battery <= 45 && room.Battery > 15}
					<TableBodyCell><Badge color="yellow">{room.Battery}</Badge></TableBodyCell>
				{/if}
				{#if room.Battery <= 15}
					<TableBodyCell><Badge color="red">{room.Battery}</Badge></TableBodyCell>
				{/if}
				{/if}
				{#if room.ShellyID == ""}
					<TableBodyCell>(none)</TableBodyCell>
				{/if}
			</TableBodyRow>
		{/each}
	</TableBody>
</Table>
