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

	export let data;

	// how do we sort the data for display
	const sortBy = { col: 'Name', ascending: true };

	// refresh every 30 seconds
	onMount(() => {
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
		const z = data.Zones.filter((z) => z.ID == zoneID);
		return z[0].Name;
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
		<TableHeadCell>Target Temp</TableHeadCell>
		<TableHeadCell on:click={tablesort('Temperature')}>Temperature</TableHeadCell>
		<TableHeadCell on:click={tablesort('Humidity')}>Humidity</TableHeadCell>
		<TableHeadCell on:click={tablesort('Battery')}>Battery</TableHeadCell>
	</TableHead>
	<TableBody>
		{#each data.Rooms as room}
			<TableBodyRow>
				<TableBodyCell><A href="/room/{room.ID}">{room.Name}</A></TableBodyCell>
				<TableBodyCell>
					{#if !room.Occupied}
						<CloseCircleOutline />
					{/if}
					{#if room.Occupied}
						<UsersGroupOutline />
					{/if}
				</TableBodyCell>
				<TableBodyCell><A href="/zone/{room.Zone}">{zoneName(room.Zone)}</A></TableBodyCell>
				<TableBodyCell><A href="/zone/{room.Zone}">{room.Targets.Min + 3}</A></TableBodyCell>

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
